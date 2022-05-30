package mq

import (
	"awesome/anet"
	"container/list"
	"log"
	"sync"
	"time"
)

/**
 *	可靠队列已发送map,用来实现统一的超时管理,sequenceId--->AmqMessage
 * 		所有已发布的消息都存储在这里等待超时，并删除。
 **/
var reliableWaitMap sync.Map

func initCore() {
	anet.sequenceIdLocker = new(sync.Mutex)
	//reliableMsgCache = make([]*AmqMessage,defaultAmqMsgCacheCapacity)
	reliableMsgCache = NewSafeList()
	//unreliableMsgCache = make([]*AmqMessage,defaultAmqMsgCacheCapacity)
	unreliableMsgCache = NewSafeList()
	go reliableLoop()
	go unreliableLoop()
	go timeoutLoop()
	go initDebug()
}

/***************
* todo,待严格的单元测试
* 可靠的队列的消费循环。
* 	reliableMsgCache 链表中pop一条数据，
*	将数据尝试写出给订阅者
*	随机模式如果写出失败，则将消息放到链表尾部，等待下次重试。如果是写出成功，从List删除。
*	全部扇出模式，扇出给所有订阅者，如果没有任何订阅者写出成功，则放入队尾，等待重试，如果是部分成功或者全部成功，不再重试，从list删除。
*	不管消息有没有成功写出，都会添加到reliableWaitMap中，由timeoutLoop检测是否有超时。
* 	如果是正常的被消费了，则reliableCallback回调函数将对应的消息从reliableWaitMap中删除。
******************/
func reliableLoop() {
	var header *list.Element = nil
	for {
		//log.Println("链表reliableMsgCache的长度", reliableMsgCache.Len())
		if item := reliableMsgCache.Front(); item != nil {

			if header == nil {
				header = item
			} else if header == item {
				header = nil
				//log.Println("遍历结束，休眠等待下一次")
				time.Sleep(defaultLoopInterval * time.Millisecond)
			}
			msg := item.Value.(*AmqMessage)
			//log.Printf("读取到可靠消息队列，进行处理 %v, %p", msg.msg, msg.msg)
			reliableWaitMap.Store(msg.msg.SequenceID, msg)
			result := processReliable(msg)
			if result == 0 {
				//已写出，删除之，并转储到reliableWaitMap,在timeoutLoop中轮候超时，或者cb中正常返回。
				log.Println("已写出，删除之，并转储到reliableWaitMap,在timeoutLoop中轮候超时，或者cb中正常返回。")
				reliableMsgCache.Remove(item)
			} else if result == 1 {
				//log.Println("返回1，没有订阅者，把消息从队列头移动到尾部。在下一次轮训的时候处理。创建时间:", msg.createTimestampMillisecond, "当前时间:", time.Now().Unix())
				//返回1，没有订阅者，把消息从队列头移动到尾部。在下一次轮训的时候处理。
				reliableMsgCache.MoveToBack(item)
			} else if result == -2 {
				//超时，丢弃
				log.Printf("超时，丢弃,topic:%d, cmd:%d", msg.msg.ReserveLow, msg.msg.Cmd)
				reliableMsgCache.Remove(item)
			} else if result == -1 {
				//有订阅者，但是写出失败，移至队尾部
				log.Println("有订阅者，但是写出失败，移至队尾部")
				reliableMsgCache.MoveToBack(item)
			}
		} else {
			time.Sleep(defaultNullLoopInterval * time.Millisecond)
		}
	}
}

/************
 * 队列超时检测。
 * todo,需要测试是否有泄露的情况。
 ****************/
func timeoutLoop() {
	for {
		time.Sleep(defaultLoopInterval * time.Millisecond) //100毫秒检测一次超时。1秒钟检测10次。
		var len = 0
		reliableWaitMap.Range(func(key, value interface{}) bool {
			msg := value.(*AmqMessage)
			if msg != nil {
				if msg.createTimestampMillisecond+defaultTimeoutMillisecond < time.Now().UnixMilli() {
					//超时
					log.Printf("超时，解散任务:%v+", msg.msg)
					msg.response(AmqAckTypeTimeout, msg.msg)
					reliableWaitMap.Delete(key)
				}
			} else {
				log.Println("超时队列转*AmqMessage失败", value)
			}
			len += 1
			return true
		})

		//log.Println("reliableWaitMap len", len)

	}
}

func unreliableLoop() {
	for {
		if item := unreliableMsgCache.Front(); item != nil {
			msg := item.Value.(*AmqMessage)
			log.Printf("读取到不可靠的数据 %v ", msg.msg)
			switch msg.msg.ReserveLow {
			case AmqCmdDefUnreliable2All:
				transUnreliableToAll(AMQTopic(msg.msg.Cmd), msg.msg)
			case AmqCmdDefUnreliable2RandomOne:
				transUnreliableToRandomOne(AMQTopic(msg.msg.Cmd), msg.msg)
			}
			unreliableMsgCache.Remove(item)
		} else {
			time.Sleep(defaultLoopInterval * time.Millisecond)
		}
	}
}

/*********
 * 可靠的消息队列的处理逻辑
 *		如果成功写出，返回0，
 *		如果遍历完，没有可写的订阅者的订阅者，返回1,
 *		如果有订阅者，但是写出失败，返回-1。
 * 		如果超时，返回-2
 ***********/
func processReliable(msg *AmqMessage) (result int) {
	//log.Println("processReliable处理", msg)
	if msg.createTimestampMillisecond+defaultTimeoutMillisecond < time.Now().UnixMilli() {
		//消息已经超时，返回超时
		log.Println("消息超时，丢弃数据")
		msg.response(AmqAckTypeTimeout, msg.msg)
		return -2
	}
	if msg.msg.ReserveLow == AmqCmdDefReliable2RandomOne {
		result = transReliableToRandomOne(AMQTopic(msg.msg.Cmd), msg.msg, reliableCallback)
		//log.Println("AmqCmdDefReliable2RandomOne=====================》", msg.msg)
	} else if msg.msg.ReserveLow == AmqCmdDefReliable2SpecOne {
		log.Println("AmqCmdDefReliable2SpecOne=====================》", msg.msg)
		result = transReliableToSpecOne(AMQTopic(msg.msg.Cmd), msg.msg, reliableCallback)
	}
	return result
}

/***************
 * 可靠队列的回调
 *		reliableWaitMap中查找注册的message，并转发给对应发布者。然后删除map中的内容。
 *		如果未查找到，则说明已经被超时机制超时了。记录错误，并丢弃消息。
 ******************/
func reliableCallback(pack *anet.PackHead) {
	log.Println("mq收到订阅者的回包", string(pack.Body))
	if v, ok := reliableWaitMap.Load(pack.SequenceID); ok {
		if msg, isOk := v.(*AmqMessage); isOk {
			//删除
			//log.Println("mq中间件收到订阅者的回包，删除waitMap的节点", pack.SequenceID)
			reliableWaitMap.Delete(pack.SequenceID)
			//回包会回填mq的客户端发布的MQ的sequenceId,需要注意Delete与response的先后顺序的差异会导致SequenceID被污染。
			msg.response(AmqAckTypeSuccess, pack)
		}
	} else {
		log.Println("mq的proxy节点收到的回包已经超时")
	}
}

func initDebug() {
	go func() {
		for true {
			time.Sleep(1 * time.Second)
			log.Println("reliableMsgCache", reliableMsgCache.Len())
			log.Println("unreliableMsgCache", unreliableMsgCache.Len())
		}
	}()
}

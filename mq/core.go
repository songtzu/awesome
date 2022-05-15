package mq

import (
	"awesome/anet"
	"container/list"
	"log"
	"sync"
	"time"
)

func init() {
	sequenceIdLocker = new(sync.Mutex)
	//reliableMsgCache = make([]*AmqMessage,defaultAmqMsgCacheCapacity)
	reliableMsgCache = NewSafeList()
	//unreliableMsgCache = make([]*AmqMessage,defaultAmqMsgCacheCapacity)
	unreliableMsgCache = NewSafeList()
	go reliableLoop()
	go unreliableLoop()
	go timeoutLoop()
}

/***************
* 可靠的队列的消费循环。
*
******************/
func reliableLoop() {
	var header *list.Element = nil
	for {
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
				log.Println("返回1，没有订阅者，把消息从队列头移动到尾部。在下一次轮训的时候处理。创建时间:", msg.createTimestampMillisecond, "当前时间:", time.Now().Unix())
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
			time.Sleep(defaultLoopInterval * time.Millisecond)
			//log.Println("reliableLoop 休眠")

		}
	}
}

/************
 * 队列超时检测。
 ****************/
func timeoutLoop() {
	for {
		time.Sleep(defaultLoopInterval * time.Millisecond) //100毫秒检测一次超时。1秒钟检测10次。
		reliableWaitMap.Range(func(key, value interface{}) bool {
			msg := value.(*AmqMessage)
			if msg != nil {
				if msg.createTimestampMillisecond+defaultTimeoutMillisecond < getMillisecondTimestamp() {
					//超时
					log.Printf("超时，解散任务:%v+", msg.msg)
					msg.response(AmqAckTypeTimeout, msg.msg)

					reliableWaitMap.Delete(key)
				}
			} else {
				log.Println("超时队列转*AmqMessage失败", value)
			}
			return true
		})

	}
}

func unreliableLoop() {
	for {
		if item := unreliableMsgCache.Front(); item != nil {
			msg := item.Value.(*AmqMessage)
			log.Printf("读取到不可靠的数据 %v ", msg.msg)
			switch msg.msg.ReserveLow {
			case AmqCmdDefUnreliable2All:
				mqUnreliable2All(AMQTopic(msg.msg.Cmd), msg.msg)
			case AmqCmdDefUnreliable2RandomOne:
				mq2UnreliableRandomOne(AMQTopic(msg.msg.Cmd), msg.msg)
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
	log.Println("processReliable处理", msg)
	if msg.createTimestampMillisecond+defaultTimeoutMillisecond < getMillisecondTimestamp() {
		//消息已经超时，返回超时
		log.Println("消息超时，丢弃数据")
		msg.response(AmqAckTypeTimeout, msg.msg)
		return -2
	}
	if msg.msg.ReserveLow == AmqCmdDefReliable2RandomOne {
		log.Println("AmqCmdDefReliable2RandomOne=====================》", msg.msg)
		result = transReliableToRandomOne(AMQTopic(msg.msg.Cmd), msg.msg, reliableCallback)
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

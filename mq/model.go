package mq

import (
	"awesome/anet"
	"container/list"
	"fmt"
	"log"
	"math"
	"sync"
	"time"
)

type AmqMessage struct {
	originalSequenceId         uint32 //发布者传过来的包序。保存记录起来。
	msg                        *anet.PackHead
	sourceConn                 *anet.Connection
	createTimestampMillisecond int64
	pushedSubscriberIds        []int //xmqSubImpl.id，推送过的订阅者ID。用于记录推送失败后，推送给其他订阅者。
}

func getMillisecondTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

type SafeList struct {
	mutex sync.RWMutex
	list  *list.List
}

func NewSafeList() *SafeList {
	return &SafeList{list: list.New()}
}

func (s *SafeList) PushBack(v interface{}) {
	s.mutex.Lock()
	s.list.PushBack(v)
	s.mutex.Unlock()
}

func (s *SafeList) Front() *list.Element {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.list.Front()
}

func (s *SafeList) Len() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.list.Len()
}

func (s *SafeList) Remove(e *list.Element) interface{} {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.list.Remove(e)
}

func (s *SafeList) MoveToBack(e *list.Element) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.list.MoveToBack(e)
}

const defaultAmqMsgCacheCapacity = 1000
const defaultLoopInterval = 1
const defaultTimeoutDelay = defaultLoopInterval * 2

const defaultTimeoutMillisecond = 30000

//var reliableMsgCache []*AmqMessage
var reliableMsgCache *SafeList

//var unreliableMsgCache []*AmqMessage
var unreliableMsgCache *SafeList

/**
 *	可靠队列已发送map,用来实现统一的超时管理,sequenceId--->AmqMessage
 * 		所有已发布的消息都存储在这里等待超时，并删除。
 **/
var reliableWaitMap sync.Map

var sequenceId uint32 = 1000

var sequenceIdLocker *sync.Mutex

//生成序列号
func genSequenceId() uint32 {
	sequenceIdLocker.Lock()
	if sequenceId >= math.MaxUint32 {
		sequenceId = 1000
	}
	sequenceId += 1
	sequenceIdLocker.Unlock()
	return sequenceId
}
func (a *AmqMessage) OnTimeOut() {
	log.Println(fmt.Sprintf(" %s time out", string(a.msg.Body)))
}

func (a *AmqMessage) response(ackType AmqAckType, pack *anet.PackHead) {
	pack.ReserveHigh = ackType
	pack.SequenceID = a.originalSequenceId //回填sequenceId
	a.sourceConn.WriteMessage(pack)
	//log.Println("写回数据给发布者,",n,err,string(pack.Body))

	//log.Println("response:", n, err)
}

func (a *AmqMessage) OnFillRemove() {
	log.Println(fmt.Sprintf(" %s fill remove", string(a.msg.Body)))
}
func pushReliableMsg(msg *anet.PackHead, source *anet.Connection) {
	var originId = msg.SequenceID

	msg.SequenceID = genSequenceId()

	//reliableMsgCache = append(reliableMsgCache,&AmqMessage{sourceConn:source,msg:msg,createTimestamp:time.Now().Unix()})
	reliableMsgCache.PushBack(&AmqMessage{sourceConn: source, msg: msg, createTimestampMillisecond: getMillisecondTimestamp(), originalSequenceId: originId})
	//log.Printf("%d,可靠的消息队列长度  %d", originId, reliableMsgCache.Len())
}

func pushUnreliableMsgCache(msg *anet.PackHead, source *anet.Connection) {
	var originId = msg.SequenceID
	msg.SequenceID = genSequenceId()
	//unreliableMsgCache = append(reliableMsgCache,&AmqMessage{sourceConn:source,msg:msg,createTimestamp:time.Now().Unix()})
	unreliableMsgCache.PushBack(&AmqMessage{sourceConn: source, msg: msg, createTimestampMillisecond: getMillisecondTimestamp(), originalSequenceId: originId})
}

/***************
 * 可靠队列的回调
 *		reliableWaitMap中查找注册的message，并转发给对应发布者。然后删除map中的内容。
 *		如果未查找到，则说明已经被超时机制超时了。记录错误，并丢弃消息。
 ******************/
func reliableCallback(pack *anet.PackHead) {
	//log.Println("mq收到订阅者的回包", string(pack.Body))
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
				//log.Println("已写出，删除之，并转储到reliableWaitMap,在timeoutLoop中轮候超时，或者cb中正常返回。")
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
	//log.Println("processReliable处理", msg)
	if msg.createTimestampMillisecond+defaultTimeoutMillisecond < getMillisecondTimestamp() {
		//消息已经超时，返回超时
		log.Println("消息超时，丢弃数据")
		msg.response(AmqAckTypeTimeout, msg.msg)
		return -2
	}
	if msg.msg.ReserveLow == AmqCmdDefReliable2RandomOne {
		//log.Println("AmqCmdDefReliable2RandomOne=====================》", msg.msg)
		result = mq2ReliableRandomOne(AMQTopic(msg.msg.Cmd), msg.msg, reliableCallback)
	} else if msg.msg.ReserveLow == AmqCmdDefReliable2SpecOne {
		//log.Println("AmqCmdDefReliable2SpecOne=====================》", msg.msg)
		result = mqReliable2All(AMQTopic(msg.msg.Cmd), msg.msg, reliableCallback)
	}
	return result
}

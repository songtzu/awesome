package mq

import (
	"awesome/anet"
	"container/list"
	"fmt"
	"log"
	"sync"
	"time"
)

type AmqMessage struct {
	originalSequenceId         uint32 //发布者传过来的包序。保存记录起来。
	msg                        *anet.PackHead
	sourceConn                 *anet.Connection    //发布者的conn句柄。保存
	srcChan                    chan *anet.PackHead //发布者的chan，此句柄仅用于http模式的发布者。sourceConn为空的时候，srcChan有值。
	createTimestampMillisecond int64
	pushedSubscriberIds        []int //xmqSubImpl.id，推送过的订阅者ID。用于记录推送失败后，推送给其他订阅者。
}

func getMillisecondTimestamp() int64 {
	return time.Now().UnixMilli()
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

var reliableMsgCache *SafeList

var unreliableMsgCache *SafeList

func (a *AmqMessage) OnTimeOut() {
	log.Printf(" %s time out", string(a.msg.Body))
}

func (a *AmqMessage) response(ackType AmqAckType, pack *anet.PackHead) {
	pack.ReserveHigh = ackType
	pack.SequenceID = a.originalSequenceId //回填sequenceId
	if a.sourceConn != nil {
		a.sourceConn.WriteMessage(pack)
	} else {
		a.srcChan <- pack
	}

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
	reliableMsgCache.PushBack(&AmqMessage{sourceConn: source, msg: msg, createTimestampMillisecond: time.Now().UnixMilli(), originalSequenceId: originId})
	//log.Printf("%d,可靠的消息队列长度  %d", originId, reliableMsgCache.Len())
}

func pushReliableMsgFromHttpSvr(msg *anet.PackHead, srcChan chan *anet.PackHead) {
	var originId = msg.SequenceID

	msg.SequenceID = genSequenceId()

	//reliableMsgCache = append(reliableMsgCache,&AmqMessage{sourceConn:source,msg:msg,createTimestamp:time.Now().Unix()})
	reliableMsgCache.PushBack(&AmqMessage{srcChan: srcChan,
		msg: msg, createTimestampMillisecond: time.Now().UnixMilli(),
		sourceConn:         nil,
		originalSequenceId: originId})
	//log.Printf("%d,可靠的消息队列长度  %d", originId, reliableMsgCache.Len())
}

func pushUnreliableMsgCache(msg *anet.PackHead, source *anet.Connection) {
	var originId = msg.SequenceID
	msg.SequenceID = genSequenceId()
	//unreliableMsgCache = append(reliableMsgCache,&AmqMessage{sourceConn:source,msg:msg,createTimestamp:time.Now().Unix()})
	unreliableMsgCache.PushBack(&AmqMessage{sourceConn: source, msg: msg, createTimestampMillisecond: time.Now().UnixMilli(), originalSequenceId: originId})
}

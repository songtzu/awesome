package mq

import (
	"awesome/anet"
	"container/list"
	"sync"
	"time"
)

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

func (s *SafeList) Front() (front *list.Element) {
	s.mutex.RLock()
	front = s.list.Front()
	s.mutex.RUnlock()
	return front
}

func (s *SafeList) Len() (len int) {
	if s == nil {
		return 0
	}
	s.mutex.RLock()
	len = s.list.Len()
	s.mutex.RUnlock()
	return len
}

func (s *SafeList) Remove(e *list.Element) (v interface{}) {
	s.mutex.Lock()
	v = s.list.Remove(e)
	s.mutex.Unlock()
	return v
}

func (s *SafeList) MoveToBack(e *list.Element) {
	s.mutex.Lock()
	s.list.MoveToBack(e)
	s.mutex.Unlock()
}

var reliableMsgCache *SafeList

var unreliableMsgCache *SafeList

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

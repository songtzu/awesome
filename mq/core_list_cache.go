package mq

import (
	"awesome/anet"
	"container/list"
	"sync"
	"time"
)

type SafeList struct {
	mutex       sync.RWMutex
	list        *list.List
	wakeChan    chan int
	PushCount   int64
	RemoveCount int64
}

func NewSafeList() *SafeList {
	return &SafeList{list: list.New(), wakeChan: make(chan int, 1)}
}

func (s *SafeList) PushBack(v interface{}) {

	s.mutex.Lock()
	//log.Printf("PushBack len:%d, pushCount:%d", len(s.wakeChan), s.PushCount)
	s.list.PushBack(v)
	s.PushCount += 1
	if len(s.wakeChan) == 0 {
		s.wakeChan <- 0
	}
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
	s.RemoveCount += 1
	//log.Printf("pushCount:%d, removeCount:%d", s.PushCount, s.RemoveCount)
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
	msg.SequenceID = anet.AllocateNewSequenceId()
	reliableMsgCache.PushBack(&AmqMessage{sourceConn: source, msg: msg, createTimestampMillisecond: time.Now().UnixMilli(), originalSequenceId: originId})
}

func pushReliableMsgFromHttpSvr(msg *anet.PackHead, srcChan chan *anet.PackHead) {
	var originId = msg.SequenceID
	msg.SequenceID = anet.AllocateNewSequenceId()
	reliableMsgCache.PushBack(&AmqMessage{srcChan: srcChan,
		msg: msg, createTimestampMillisecond: time.Now().UnixMilli(),
		sourceConn:         nil,
		originalSequenceId: originId})
}

func pushUnreliableMsgCache(msg *anet.PackHead, source *anet.Connection) {
	var originId = msg.SequenceID
	msg.SequenceID = anet.AllocateNewSequenceId()
	unreliableMsgCache.PushBack(&AmqMessage{sourceConn: source, msg: msg, createTimestampMillisecond: time.Now().UnixMilli(), originalSequenceId: originId})
}

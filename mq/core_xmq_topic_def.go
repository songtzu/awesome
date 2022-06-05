package mq

import (
	"sync"
)

type xmqTopic struct {
	topic AMQTopic
	subs  []*xmqSubImpl //此topic的订阅者Implement对象。
	nodes chan *AmqMessage
	sync.Mutex
}

func (x *xmqTopic) removeSubImpl(imp *xmqSubImpl) {
	x.Lock()
	for i := 0; i < len(x.subs); {
		if x.subs[i].id == imp.id {
			x.subs = append(x.subs[:i], x.subs[i+1:]...)
			return
		}
	}
	x.Unlock()
}

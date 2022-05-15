package mq

import (
	"time"
)

type xmqSub struct {
	topic AMQTopic
	subs  []*xmqSubImpl //此topic的订阅者Implement对象。
	nodes chan *AmqMessage
}

func (x *xmqSub) publish(node *AmqMessage) {
	if node.msg.Cmd == AMQCmdDefPub || node.msg.Cmd == AmqCmdDefUnreliable2All {
		x.unreliable2All(node)
	} else if node.msg.Cmd == AmqCmdDefReliable2RandomOne {
		x.reliable2RandomOne(node)
	}
}
func (x *xmqSub) startTransSub() {
	for {
		if len(x.subs) == 0 {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		select {
		case msg := <-x.nodes:
			x.publish(msg)
		}
	}
}

//
//func (x *xmqSub) enqueue(node *AmqMessage) {
//	log.Println("message into queue", string(node.msg.Body))
//	x.nodes <- node
//}

//选择一个订阅者处理，如果超时或者未返回，则轮转下一个订阅者处理。或返回超时给发布者。
func (x *xmqSub) reliable2RandomOne(node *AmqMessage) {

}

//不可靠发布（无回包），所有订阅者都会收到此message.
//		仅仅使用类似日志，统计类业务，不同的sub订阅同一个topic，加工整理成不同的统计结果。
//		业务不依赖此消息队列是否有启动进程处理业务。
func (x *xmqSub) unreliable2All(node *AmqMessage) {
	for _, sub := range x.subs {
		sub.conn.WriteMessage(node.msg)
	}
}

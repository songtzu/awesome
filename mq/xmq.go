package mq

import (
	"awesome/anet"
	"sync"
)

var xmqInstance *Xmq = nil

/*
 * received message from transPub and publish to xsub.
 * 		contains a nodes of server accept connection from apub.
 *			and a nodes of server accept connection from xsub.
 */

type Xmq struct {
	//topicMap	map[AMQTopic]*xmqSub		//this is deal with the real sub, which would sub some topics.
	topicMap sync.Map //this is deal with the real sub, which would sub some topics.map:topic--->xmqSub
}

/*NewXmq
 * xPubBindAddress:真实的发布者连接此地址
 * xSubBindAddress：真是的订阅者连接此地址。
 */
func NewXmq(xPubBindAddress string, xSubBindAddress string) (xmq *Xmq) {
	xmq = &Xmq{}
	if xmqInstance != nil {
		return xmqInstance
	}
	go initCore()
	xmqInstance_ := &Xmq{
		//conn:c,
		//topicMap:make(map[AMQTopic]*xmqSub),
	}
	xmqInstance_.startSub(xSubBindAddress)
	xmqInstance_.startPub(xPubBindAddress)
	xmqInstance = xmqInstance_
	return xmqInstance
}

//发布者的代理，发布者连接此服务
func (x *Xmq) startPub(xPubBindAddress string) {
	impl := &xmqPubImpl{}
	go anet.StartTcpSvr(xPubBindAddress, impl)
}

//订阅者的代理，订阅者连接此服务。
func (x *Xmq) startSub(xSubBindAddress string) {
	impl := &xmqSubImpl{}
	go anet.StartTcpSvr(xSubBindAddress, impl)
}

func (x *Xmq) newXmqSub(topic AMQTopic, si *xmqSubImpl) *xmqSub {
	sub := &xmqSub{topic: topic, nodes: make(chan *AmqMessage, defaultAMQChanSize), subs: []*xmqSubImpl{}}
	if si != nil {
		sub.subs = append(sub.subs, si)
	}
	go sub.startTransSub()
	xmqInstance.topicMap.Store(topic, sub)
	return sub
}

func (x *Xmq) subTopics(topics []AMQTopic, a *xmqSubImpl) {
	for _, t := range topics {
		if v, ok := x.topicMap.Load(t); ok {
			s := v.(*xmqSub)
			s.subs = append(s.subs, a)
		} else {
			x.newXmqSub(t, a)
		}
	}

}

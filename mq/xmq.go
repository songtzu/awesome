package mq

import (
	"awesome/anet"
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

var xmqInstance *Xmq = nil

/*
 * received message from transPub and publish to xsub.
 * 		contains a nodes of server accept connection from apub.
 *			and a nodes of server accept connection from xsub.
 */

type Xmq struct {
	//topicMap	map[AMQTopic]*xmqSub		//this is deal with the real sub, which would sub some topics.
	topicMap sync.Map //this is deal with the real sub, which would sub some topics.
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

func (x *Xmq) enqueuePub2SubChan(node *AmqMessage) {
	log.Println(fmt.Sprintf("enqueue pub 2 sub chan, message is %s", string(node.msg.Body)))
	if v, ok := x.topicMap.Load(AMQTopic(node.msg.ReserveLow)); ok {
		s := v.(*xmqSub)
		s.enqueue(node)
	} else {
		x.newXmqSub(AMQTopic(node.msg.ReserveLow), nil).enqueue(node)
	}
}

func (x *Xmq) newXmqSub(topic AMQTopic, si *xmqSubImpl) *xmqSub {
	sub := &xmqSub{nodes: make(chan *AmqMessage, defaultAMQChanSize), subs: []*xmqSubImpl{}}
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
func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

/***************
 * 	消息转发给随机一个，如果write返回了错误码，则说明此节点有故障，更换一个节点写出。只要写出，就不管是否被处理，等待上层业务的超时机制。
 *		如果成功写出，返回0，
 *		如果遍历完，没有可写的订阅者的订阅者，返回1,
 *		如果有订阅者，但是写出失败，返回-1。
 *****************/
func mq2ReliableRandomOne(topic AMQTopic, pack *anet.PackHead, cb anet.DefNetIOCallback) int {
	if v, ok := xmqInstance.topicMap.Load(topic); ok {
		s := v.(*xmqSub)
		rand.Seed(time.Now().UnixNano())
		x := rand.Intn(100)
		startIndex := x % len(s.subs)
		//log.Println("mq2ReliableRandomOne-----------", pack)
		for i := 0; i < len(s.subs); i++ {
			if _, err := s.subs[(i+startIndex)%len(s.subs)].conn.WriteMessageWithCallback(pack, cb); err == nil {
				return 0
			}
		}
		if len(s.subs) == 0 {
			return 1
		}
		return -1
	} else {
		//没有订阅者。
		return 1
	}

}

/***************
 * 	消息转发给所有订阅者。
 *		如果成功写出，返回0，
 *		如果遍历完，没有可写的订阅者的订阅者，返回1,
 *		如果有订阅者，但是写出失败，返回-1。
 *****************/
func mqReliable2All(topic AMQTopic, pack *anet.PackHead, cb anet.DefNetIOCallback) int {
	var isAllUnreachable = true
	if v, ok := xmqInstance.topicMap.Load(topic); ok {
		s := v.(*xmqSub)

		//for _,sub:=range s.subs {
		//	if _,err:=sub.conn.WriteMessageWithCallback(pack,cb);err==nil{
		//		isAllUnreachable = false
		//	}
		//}

		// todo go 程太多
		var wg sync.WaitGroup
		var ctx,cancel = context.WithCancel(context.Background())
		defer cancel()

		f := func(sub *xmqSubImpl) <- chan error{
				_,err := sub.conn.WriteMessageWithCallback(pack, cb)
				var c = make(chan error,1)
				c <- err
				close(c)
				return c
		}

		for _, sub := range s.subs {

			wg.Add(1)
			go func(sub *xmqSubImpl) {
				defer wg.Done()

				select {
				case <- ctx.Done():
					return

				case err := <- f(sub):
					if err == nil {
						isAllUnreachable = false
						cancel()
					}
				}
			}(sub)
		}

		wg.Wait()
	} else {
		//没有订阅者。
		return 1
	}
	if isAllUnreachable {
		//所有订阅者均离线
		return 1
	}
	return 0
}


/***************
 * 	消息转发给随机一个，如果write返回了错误码，则说明此节点有故障，更换一个节点写出。只要写出，就不管是否被处理
 *****************/
func mq2UnreliableRandomOne(topic AMQTopic, pack *anet.PackHead) int {
	if v, ok := xmqInstance.topicMap.Load(topic); ok {
		s := v.(*xmqSub)
		rand.Seed(time.Now().UnixNano())
		x := rand.Intn(100)
		startIndex := x % len(s.subs)
		log.Println("mq2ReliableRandomOne-----------", pack)
		for i := 0; i < len(s.subs); i++ {
			if _, err := s.subs[(i+startIndex)%len(s.subs)].conn.WriteMessage(pack); err == nil {
				return 0
			}
		}
		if len(s.subs) == 0 {
			return 1
		}
		return -1
	} else {
		//没有订阅者。
		return 1
	}

}

/***************
 * 	消息转发给所有订阅者。
 *		如果成功写出，返回0，
 *		如果遍历完，没有可写的订阅者的订阅者，返回1,
 *		如果有订阅者，但是写出失败，返回-1。
 *****************/
func mqUnreliable2All(topic AMQTopic, pack *anet.PackHead) int {
	if v, ok := xmqInstance.topicMap.Load(topic); ok {
		s := v.(*xmqSub)

		//for _,sub:=range s.subs {
		//	if _,err:=sub.conn.WriteMessageWithCallback(pack,cb);err==nil{
		//		isAllUnreachable = false
		//	}
		//}
		for _, sub := range s.subs {
			if _,err := sub.conn.WriteMessage(pack);err!=nil {
				log.Printf("err:%v",err)
			}

		}
		return 0
	}
	return 1
}
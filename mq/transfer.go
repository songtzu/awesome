package mq

import (
	"awesome/anet"
	"log"
	"math/rand"
	"time"
)

/***************
 *
 * 	消息转发给所有订阅者。
 *		如果成功写出，返回0，
 *		如果遍历完，没有可写的订阅者的订阅者，返回1,
 *		如果有订阅者，但是全部写出失败，返回-1。
 * 		如果有订阅者，部分写出成功部分失败，返回0
 *	由于未记录成功写出的订阅者，部分成功的重试机制会导致上一次写出成功的消费者，重复消费。
 *	todo,考虑记录部分成功的消费者，这样，可以实现重试。
 *****************/
func transReliableToSpecOne(topic AMQTopic, pack *anet.PackHead, cb anet.DefNetIOCallback) int {
	var isAllUnreachable = true
	if v, ok := xmqInstance.topicMap.Load(topic); ok {
		s := v.(*xmqSub)
		for _, sub := range s.subs {

			if _, err := sub.conn.WriteMessageWithCallback(pack, cb); err == nil {
				isAllUnreachable = false
			}

		}

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
 * 	消息转发给随机一个，如果write返回了错误码，则说明此节点有故障，更换一个节点写出。只要写出，就不管是否被处理，等待上层业务的超时机制。
 *		如果成功写出，返回0，
 *		如果遍历完，没有可写的订阅者的订阅者，返回1,
 *		如果有订阅者，但是写出失败，返回-1。
 *****************/
func transReliableToRandomOne(topic AMQTopic, pack *anet.PackHead, cb anet.DefNetIOCallback) int {
	if v, ok := xmqInstance.topicMap.Load(topic); ok {
		s := v.(*xmqSub)
		rand.Seed(time.Now().UnixNano())
		x := rand.Intn(100)
		startIndex := x % len(s.subs)
		//log.Println("transReliableToRandomOne-----------", pack)
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
 * 	消息转发给随机一个，如果write返回了错误码，则说明此节点有故障，更换一个节点写出。只要写出，就不管是否被处理
 *****************/
func mq2UnreliableRandomOne(topic AMQTopic, pack *anet.PackHead) int {
	if v, ok := xmqInstance.topicMap.Load(topic); ok {
		s := v.(*xmqSub)
		rand.Seed(time.Now().UnixNano())
		x := rand.Intn(100)
		startIndex := x % len(s.subs)
		log.Println("transReliableToRandomOne-----------", pack)
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
			if _, err := sub.conn.WriteMessage(pack); err != nil {
				log.Printf("err:%v", err)
			}

		}
		return 0
	}
	return 1
}

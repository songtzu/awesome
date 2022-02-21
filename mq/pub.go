package mq

import (
	"awesome/anet"
	"fmt"
	"log"
	"strings"
)

/*
 * pub is a tcp request client.
 *
 */
var pubInstance *AMQPub = nil

type AMQPub struct {
	//conn   *net.Connection
	topicMap map[AMQTopic][]*pubImpl
}

func GetAMQPubInstance(bindAddress string) *AMQPub {
	if pubInstance != nil {
		return pubInstance
	}
	impl := &pubImpl{cb: recieveCallback}

	go anet.StartTcpSvr(bindAddress, impl)
	//c:=net.NewTcpClientConnect(bindAddress,impl,1)
	pub := &AMQPub{
		//conn:c,
		topicMap: make(map[AMQTopic][]*pubImpl),
	}
	pubInstance = pub
	return pubInstance
}
func recieveCallback(pack *anet.PackHead) {

}
func printMap() {
	for k, v := range pubInstance.topicMap {
		for _, c := range v {
			log.Println(fmt.Sprintf("订阅主题%v，订阅者%d", k, c.id))
		}
	}
}
func (a *AMQPub) PubUnreliable2AllMessage(data []byte, cmd anet.PackHeadCmd) (err error) {
	msg := &anet.PackHead{ReserveLow: AmqCmdDefUnreliable2All, Cmd: cmd,
		SequenceID: 0,
		Length:     uint32(len(data)), Body: data}
	str := string(data)
	printMap()
	for k, v := range pubInstance.topicMap {
		if strings.HasPrefix(str, string(k)) {
			for _, c := range v {
				log.Println(fmt.Sprintf("发布的消息内容%s,订阅的主题%v,订阅者ID%d", str, k, c.id))
				//log.Println(c.id, k)
				log.Println(c.conn.WriteMessage(msg))
			}
		}
	}
	return nil
}

/********
 * 发布给
 *********/
func (a *AMQPub) PubReliable2RandomOneMessage(data []byte, cmd anet.PackHeadCmd) (err error) {
	msg := &anet.PackHead{ReserveLow: AmqCmdDefReliable2RandomOne, Cmd: cmd,
		SequenceID: 0,
		Length:     uint32(len(data)), Body: data}
	str := string(data)
	printMap()
	for k, v := range pubInstance.topicMap {
		if strings.HasPrefix(str, string(k)) {
			for _, c := range v {
				log.Println(fmt.Sprintf("发布的消息内容%s,订阅的主题%v,订阅者ID%d", str, k, c.id))
				//log.Println(c.id, k)
				log.Println(c.conn.WriteMessage(msg))
			}
		}
	}
	return nil
}

func (a *AMQPub) hasSubed(v []*pubImpl, conn *pubImpl) bool {
	for _, i := range v {
		if i.id == conn.id {
			return true
		}
	}
	return false
}

func (a *AMQPub) subATopic(topic AMQTopic, conn *pubImpl) (isOk bool) {
	if v, ok := pubInstance.topicMap[topic]; ok {
		log.Println("存在该ID", ok, v, topic)
		if a.hasSubed(v, conn) {
			log.Println("重复订阅", conn.id)
			//check conn has sub this topic before, to avoid multi subs.
			return true
		}
		v = append(v, conn)
	} else {
		log.Println("新的消息订阅", topic, "订阅者ID", conn.id)
		subs := []*pubImpl{conn}
		pubInstance.topicMap[topic] = subs
	}

	return true
}

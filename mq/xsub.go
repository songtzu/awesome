package mq

import (
	"awesome/anet"
	"encoding/json"
)

type AMQXSub struct {
	conn *anet.Connection
}

func NewXSub(bindAddress string, cb AMQCallback) *AMQXSub {
	impl := &subImpl{cb: cb}
	c := anet.NewTcpClientConnect(bindAddress, impl, 1000, true)
	sub := &AMQXSub{
		conn: c,
	}
	return sub
}

//
//func (a *AMQXSub) MessagePub(data []byte) ( err error) {
//	msg := &net.PackHead{Cmd: AMQCmdDefPub,
//		SequenceID: 0,
//		Length: uint32(len(data)), Body: data}
//	str:=string(data)
//	printMap()
//	for k,v:=range pubInstance.topicMap{
//		if strings.HasPrefix(str,string(k)){
//
//			for _,c:=range v{
//				log.Println(fmt.Sprintf("发布的消息内容%s,订阅的主题%s,订阅者ID%d",str,k, c.id))
//				//log.Println(c.id, k)
//				log.Println(c.conn.WriteMessage(msg))
//			}
//		}
//	}
//	return nil
//}
func (a *AMQXSub) TopicSubscription(topics []AMQTopic) (n int, err error) {
	t := &AMQProtocolSubTopic{Topics: topics}
	bin, _ := json.Marshal(t)
	//log.Println(string(bin), "size:",len(bin))
	msg := &anet.PackHead{Cmd: AMQCmdDefSubTopic,
		SequenceID: 0,
		Length:     uint32(len(bin)), Body: bin}

	return a.conn.WriteMessage(msg)
}

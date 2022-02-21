package mq

import (
	"awesome/anet"

	"encoding/json"
)

type AMQSub struct {
	conn *anet.Connection
}

func NewAMQSub(bindAddress string, cb AMQCallback) *AMQSub {
	impl := &subImpl{cb: cb}
	c := anet.NewTcpClientConnect(bindAddress, impl, 1000, true)
	sub := &AMQSub{
		conn: c,
	}
	return sub
}

//func (a *AMQSub) MessagePub(data []byte) (n int, err error){
//	return a.conn.WriteBytes(data)
//}

func (a *AMQSub) TopicSubscription(topics []AMQTopic) (n int, err error) {
	t := &AMQProtocolSubTopic{Topics: topics}
	bin, _ := json.Marshal(t)
	//log.Println(string(bin), "size:",len(bin))
	msg := &anet.PackHead{Cmd: AMQCmdDefSubTopic,
		SequenceID: 0,
		Length:     uint32(len(bin)), Body: bin}

	return a.conn.WriteMessage(msg)
}

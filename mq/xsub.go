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

func (a *AMQXSub) TopicSubscription(topics []AMQTopic) (n int, err error) {
	t := &AMQProtocolSubTopic{Topics: topics}
	bin, _ := json.Marshal(t)
	//log.Println(string(bin), "size:",len(bin))
	msg := &anet.PackHead{Cmd: AMQCmdDefSubTopic,
		SequenceID: 0,
		Length:     uint32(len(bin)), Body: bin}

	return a.conn.WriteMessage(msg)
}

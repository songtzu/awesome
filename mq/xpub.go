package mq

import (
	"awesome/anet"
)

type AMQXPub struct {
	conn *anet.Connection
}

func NewXPub(bindAddress string) *AMQXPub {
	impl := &xPubImpl{}
	c := anet.NewTcpClientConnect(bindAddress, impl, 1000, true)
	sub := &AMQXPub{
		conn: c,
	}
	return sub
}

/**************
 * 无回包订阅模型
 * 	pub
 ****************/
func (a *AMQXPub) MessagePub(topic AMQTopic, data []byte) (err error) {
	msg := &anet.PackHead{Cmd: AMQCmdDefPub,
		SequenceID: 0,
		Length:     uint32(len(data)), Body: data, ReserveLow: uint32(topic)}
	//str:=string(data)
	//log.Println(str)
	a.conn.WriteMessage(msg)
	return nil
}

package mq

import (
	"awesome/anet"
	"encoding/json"
	"fmt"
	"log"
)

type xmqSubImpl struct {
	//reliableCallback AMQCallback
	conn *anet.Connection
	id   int
}

func (a *xmqSubImpl) IOnInit(connection *anet.Connection) {

}

func (a *xmqSubImpl) subTopic(pack *anet.PackHead) {
	msg := &AMQProtocolSubTopic{}

	if err := json.Unmarshal(pack.Body, msg); err != nil {
		log.Println("error when recieve sub topic action", err)
	} else {
		log.Println(fmt.Sprintf("订阅消息,订阅者id%d, topic %s", a.id, string(pack.Body)))

		xmqInstance.subTopics(msg.Topics, a)

	}
}
func (a *xmqSubImpl) IOnProcessPack(pack *anet.PackHead) {
	if pack.Cmd == AMQCmdDefSubTopic {
		a.subTopic(pack)
	} else {
		//proxy组件，应该转发给proxy检测回包
		log.Println("收到其他消息", pack)
	}
}

/*
 * this interface SHOULD NOT CALL close.
 */
func (a *xmqSubImpl) IOnClose(err error) (tryReconnect bool) {
	return true
}

//func (a *xmqSubImpl) IWrite(msg interface{}, ph *net.PackHead){
//
//}

func (a *xmqSubImpl) IOnConnect(isOk bool) {

}

func (a *xmqSubImpl) IOnNewConnection(connection *anet.Connection) {
	log.Println("new connection")
	//fork:=&pubImpl{reliableCallback:a.reliableCallback, id:newConnId(), conn:connection}
	a.conn = connection
	a.id = newConnId()

	//test:=&net.PackHead{Cmd:AMQCmdDefPub,Length:uint32(len([]byte("hello"))),Body:[]byte("hello")}
	//a.conn.WriteMessage(test)
	//a.id = newConnId()
	//a.conn = connection
}

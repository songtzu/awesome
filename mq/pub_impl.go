package mq

import (
	"awesome/anet"
	"encoding/json"
	"fmt"
	"log"
)

var pubConnId = 0

func newConnId() int {
	pubConnId++
	return pubConnId
}

type pubImpl struct {
	cb   AMQCallback
	conn *anet.Connection
	id   int
}

func (a *pubImpl) subTopic(pack *anet.PackHead) {
	msg := &AMQProtocolSubTopic{}

	if err := json.Unmarshal(pack.Body, msg); err != nil {
		log.Println("error when recieve sub topic action", err)
	} else {
		log.Println(fmt.Sprintf("订阅消息,订阅者id%d, topic %s", a.id, string(pack.Body)))
		for _, topic := range msg.Topics {
			pubInstance.subATopic(topic, a)
		}

	}
}

func (a *pubImpl) IOnInit(connection *anet.Connection) {

}

func (a *pubImpl) IOnProcessPack(pack *anet.PackHead) {

	/*
	 * subscriber sub some topics.
	 */
	//log.Println("订阅消息",pack.Cmd, string(pack.Body), len(pack.Body), pack.Length)
	if pack.Cmd == AMQCmdDefSubTopic {
		a.subTopic(pack)
	}
}

/*
 * this interface SHOULD NOT CALL close.
 */
func (a *pubImpl) IOnClose(err error) (tryReconnect bool) {
	return true
}

//func (a *pubImpl) IWrite(msg interface{}, ph *net.PackHead){
//
//}

func (a *pubImpl) IOnConnect(isOk bool) {

}

func (a *pubImpl) IOnNewConnection(connection *anet.Connection) {
	log.Println("new connection")
	//fork:=&pubImpl{reliableCallback:a.reliableCallback, id:newConnId(), conn:connection}
	a.conn = connection
	a.id = newConnId()

	//test := &anet.PackHead{Cmd: AMQCmdDefPub, Length: uint32(len([]byte("hello"))), Body: []byte("hello")}
	//a.conn.WriteMessage(test)
	//a.id = newConnId()
	//a.conn = connection
}

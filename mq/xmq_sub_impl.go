package mq

import (
	"awesome/anet"
	"encoding/json"
	"log"
)

type xmqSubImpl struct {
	//reliableCallback AMQCallback
	conn *anet.Connection
	id   int
}

func (a *xmqSubImpl) IOnInit(connection *anet.Connection) {

}

func (a *xmqSubImpl) subTopic(pack *anet.PackHead) (ack *AMQProtocolSubTopicAck) {
	msg := &AMQProtocolSubTopic{}
	ack = &AMQProtocolSubTopicAck{Status: 0, Message: "ok"}
	if err := json.Unmarshal(pack.Body, msg); err != nil {
		log.Printf("error :%s when recieve sub topic action", err.Error())
		ack.Status = -1
		ack.Message = "failed to unmarshal json"
	} else {
		log.Printf("xmqSubImpl订阅消息,订阅者id%d, topic %v", a.id, msg)

		xmqInstance.subTopics(msg.Topics, a)
	}
	return ack
}
func (a *xmqSubImpl) IOnProcessPack(pack *anet.PackHead, connection *anet.Connection) {
	log.Printf("xmqSubImpl,pack:%v", pack)
	if pack.Cmd == AMQCmdDefSubTopic {
		ack := a.subTopic(pack)
		connection.WriteJsonObj(ack, AMQCmdDefSubTopicAck, pack.SequenceID)
	} else {
		//proxy组件，应该转发给proxy检测回包
		log.Println("xmqSubImpl收到其他消息", pack)
	}
}

/*IOnClose
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

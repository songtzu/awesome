package mq

import (
	"awesome/anet"
	"encoding/json"
	"log"
)

type xmqSubImpl struct {
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
	xmqInstance.printSub()
	return ack
}
func (a *xmqSubImpl) IOnProcessPack(pack *anet.PackHead, connection *anet.Connection) {
	//log.Printf("xmqSubImpl,pack:%v", pack)
	if pack.Cmd == AMQCmdDefSubTopic {
		ack := a.subTopic(pack)
		connection.WriteJsonObj(ack, AMQCmdDefSubTopicAck, pack.SequenceID)
	} else {
		//proxy组件，应该转发给proxy检测回包
		//log.Println("xmqSubImpl收到其他消息", pack)
		reliableCallback(pack)
	}
}

/*IOnClose
 * this interface SHOULD NOT CALL close.
 */
func (a *xmqSubImpl) IOnClose(err error) (tryReconnect bool) {
	log.Println("xmqSubImpl IOnClose")
	xmqInstance.unSubTopics(a)
	xmqInstance.printSub()
	return true
}

func (a *xmqSubImpl) IOnConnect(isOk bool) {

}

func (a *xmqSubImpl) IOnNewConnection(connection *anet.Connection) {
	log.Println("new connection")
	a.conn = connection
	a.id = anet.GenNewId()
}

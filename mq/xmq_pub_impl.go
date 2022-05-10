package mq

import (
	"awesome/anet"
	"log"
)

type xmqPubImpl struct {
	//reliableCallback AMQCallback
	conn *anet.Connection
	id   int
}

func (a *xmqPubImpl) IOnInit(connection *anet.Connection) {

}

//收到真实发布者发布的消息
func (a *xmqPubImpl) transPub(pack *anet.PackHead) {
	xmqInstance.enqueuePub2SubChan(&AmqMessage{msg: pack, sourceConn: a.conn, createTimestampMillisecond: getMillisecondTimestamp()})
}

func (a *xmqPubImpl) IOnProcessPack(pack *anet.PackHead, connection *anet.Connection) {
	log.Println("xmqPubImpl..IOnProcessPack.", string(pack.Body), pack)
	if pack.ReserveLow == AMQCmdDefPub || pack.ReserveLow == AmqCmdDefUnreliable2All || pack.ReserveLow == AmqCmdDefUnreliable2RandomOne {
		//a.transPub(pack)
		log.Println("不可靠消息发布", pack.Cmd, string(pack.Body))
		pushUnreliableMsgCache(pack, a.conn)
	} else if pack.ReserveLow == AmqCmdDefReliable2RandomOne || pack.ReserveLow == AmqCmdDefReliable2SpecOne {
		log.Println("可靠的消息发布", pack.Cmd, string(pack.Body))
		//log.Println(pack.SequenceID,string(pack.Body))
		pushReliableMsg(pack, a.conn)
	}
}

/*
 * this interface SHOULD NOT CALL close.
 */
func (a *xmqPubImpl) IOnClose(err error) (tryReconnect bool) {
	return true
}

//func (a *xmqPubImpl) IWrite(msg interface{}, ph *net.PackHead){
//
//}

func (a *xmqPubImpl) IOnConnect(isOk bool) {

}

func (a *xmqPubImpl) IOnNewConnection(connection *anet.Connection) {
	log.Println("new connection")
	//fork:=&pubImpl{reliableCallback:a.reliableCallback, id:newConnId(), conn:connection}
	a.conn = connection
	a.id = newConnId()

	//test := &anet.PackHead{Cmd: AMQCmdDefPub, Length: uint32(len([]byte("hello"))), Body: []byte("hello")}
	//a.conn.WriteMessage(test)
	//a.id = newConnId()
	//a.conn = connection
}

//
//func (a *xmqPubImpl) modelReliable2RandomOne(pack *anet.PackHead) {
//	xmqInstance.enqueuePub2SubChan(&AmqMessage{msg:pack,sourceConn:a.conn,createTimestamp:time.Now().Unix()})
//}

package mq

import (
	"awesome/anet"
	"log"
)

type xPubImpl struct {
	cb AMQCallback
	//cbNewConn aMQNewConnCallback
	conn *anet.Connection
}

func (a *xPubImpl) IOnInit(connection *anet.Connection) {

}

func (a *xPubImpl) IOnProcessPack(pack *anet.PackHead, connection *anet.Connection) {
	if a.cb != nil {
		a.cb(pack)
	}
}

/*
 * this interface SHOULD NOT CALL close.
 */
func (a *xPubImpl) IOnClose(err error) (tryReconnect bool) {
	log.Println("IOnClose订阅连接关闭")
	return true
}

//func (a *xPubImpl) IWrite(msg interface{}, ph *net.PackHead){
//
//}

func (a *xPubImpl) IOnConnect(isOk bool) {

}

func (a *xPubImpl) IOnNewConnection(connection *anet.Connection) {

}

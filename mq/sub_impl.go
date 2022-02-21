package mq

import (
	"awesome/anet"
	"log"
)

//type aMQNewConnCallback func(conn *net.Connection)
type subImpl struct {
	cb AMQCallback
	//cbNewConn aMQNewConnCallback
	conn *anet.Connection
}

func (a *subImpl) IOnInit(connection *anet.Connection) {

}

func (a *subImpl) IOnProcessPack(pack *anet.PackHead) {
	if a.cb != nil {
		a.cb(pack)
	}
}

/*
 * this interface SHOULD NOT CALL close.
 */
func (a *subImpl) IOnClose(err error) (tryReconnect bool) {
	log.Println("IOnClose订阅连接关闭")
	return true
}

//func (a *subImpl) IWrite(msg interface{}, ph *net.PackHead){
//
//}

func (a *subImpl) IOnConnect(isOk bool) {

}

func (a *subImpl) IOnNewConnection(connection *anet.Connection) {

}

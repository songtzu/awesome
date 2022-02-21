package test

import (
	"testing"
	"awesome/anet"
	"fmt"
	"time"
)

func TestFramework(t *testing.T) {
	iml:=&clientImpl{}
	c:= anet.NewTcpClientConnect("127.0.0.1:19999",iml,1000)
	pack:=&anet.PackHead{Cmd: 1,Body:[]byte("hello testing. by jack")}
	c.WriteMessage(pack)
	time.Sleep(1*time.Minute)
}



type clientImpl struct {
	//cbNewConn aMQNewConnCallback
	conn *anet.Connection
}

func (a *clientImpl) IOnInit(connection *anet.Connection){

}

func (a *clientImpl) IOnProcessPack(pack *anet.PackHead) {

}
/*
 * this interface SHOULD NOT CALL close.
 */
func (a *clientImpl) IOnClose(err error)(tryReconnect bool ){
	fmt.Println("IOnClose订阅连接关闭")
	return true
}
//func (a *subImpl) IWrite(msg interface{}, ph *net.PackHead){
//
//}


func (a *clientImpl) IOnConnect(isOk bool){

}


func (a *clientImpl) IOnNewConnection(connection *anet.Connection){

}

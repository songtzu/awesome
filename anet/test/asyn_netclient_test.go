package test

import (
	"awesome/anet"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type AsynPingMessage struct {
	IsOk      bool   `json:"isOk"`
	Timestamp int64  `json:"timestamp"`
	Address   string `json:"address"`
}

type asynClientImpl struct {
	//cbNewConn aMQNewConnCallback
	conn *anet.Connection
}

func (a *asynClientImpl) IOnInit(connection *anet.Connection) {

}

func (a *asynClientImpl) IOnProcessPack(pack *anet.PackHead, connection *anet.Connection) {

}

/*
 * this interface SHOULD NOT CALL close.
 */
func (a *asynClientImpl) IOnClose(err error) (tryReconnect bool) {
	fmt.Println("IOnClose订阅连接关闭")
	return true
}

//func (a *subImpl) IWrite(msg interface{}, ph *net.PackHead){
//
//}

func (a *asynClientImpl) IOnConnect(isOk bool) {
	fmt.Println("=============建立链接的回调")

}

func (a *asynClientImpl) IOnNewConnection(connection *anet.Connection) {
	fmt.Println("=============建立链接的回调")
	a.conn = connection
}

func TestAsynNetClient(t *testing.T) {
	go asynRun()
	time.Sleep(1 * time.Minute)
}
func asynRun() {
	imp := &asynClientImpl{}
	imp.conn = anet.NewNetClient("ws://127.0.0.1:19999/", imp, 1000, true)
	ping := AsynPingMessage{IsOk: true, Timestamp: time.Now().Unix(), Address: "ws://127.0.0.1:19999/"}
	bin, _ := json.Marshal(ping)
	pack := &anet.PackHead{Cmd: 1, Body: bin}
	p, isTimeout := imp.conn.WriteMessageWaitResponseWithinTimeLimit(pack, 1000)

	fmt.Println("是否超时:", isTimeout, "回包内容", p)

	time.Sleep(1 * time.Minute)
}

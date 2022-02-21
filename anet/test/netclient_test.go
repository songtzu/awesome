package test

import (
	"awesome/anet"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type PingMessage struct {
	IsOk bool `json:"isOk"`
	Timestamp int64 `json:"timestamp"`
	Address string `json:"address"`
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
	fmt.Println("=============建立链接的回调")

}


func (a *clientImpl) IOnNewConnection(connection *anet.Connection){
	fmt.Println("=============建立链接的回调")
	a.conn = connection
}


func TestSynNetClient(t *testing.T)  {
	go run()
	time.Sleep(1*time.Minute)
}
func run()  {
	imp := &clientImpl{}
	imp.conn = anet.NewNetClient("ws://127.0.0.1:19999/", imp, 1000)
	ping := PingMessage{IsOk: true, Timestamp: time.Now().Unix(), Address: "ws://127.0.0.1:19999/"}
	bin, _ := json.Marshal(ping)
	pack := &anet.PackHead{Cmd: 1, Body: bin}
	_, err := imp.conn.WriteMessageWithCallback(pack, func(msg *anet.PackHead) {
		ack := &PingMessage{}
		if unmarshal_err := json.Unmarshal(msg.Body, ack); unmarshal_err == nil {
			fmt.Println("返回的数据",ack)
		}else {
			fmt.Println("错误",unmarshal_err)
		}

	})
	if err != nil {
		fmt.Println("测试未通过",err)
	}
	time.Sleep(1*time.Minute)
}
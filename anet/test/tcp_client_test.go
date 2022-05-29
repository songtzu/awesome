package test

import (
	"awesome/anet"
	"fmt"
	"log"
	"testing"
	"time"
)

type PingMessage struct {
	IsOk      bool   `json:"isOk"`
	Timestamp int64  `json:"timestamp"`
	Address   string `json:"address"`
}

type TCPClientImpl struct {
	//cbNewConn aMQNewConnCallback
	conn *anet.Connection
}

func (a *TCPClientImpl) IOnInit(connection *anet.Connection) {

}

func (a *TCPClientImpl) IOnProcessPack(pack *anet.PackHead, connection *anet.Connection) {

}

/*
 * this interface SHOULD NOT CALL close.
 */
func (a *TCPClientImpl) IOnClose(err error) (tryReconnect bool) {
	fmt.Println("IOnClose订阅连接关闭")
	return true
}

//func (a *subImpl) IWrite(msg interface{}, ph *net.PackHead){
//
//}

func (a *TCPClientImpl) IOnConnect(isOk bool) {
	fmt.Println("=============建立链接的回调")

}

func (a *TCPClientImpl) IOnNewConnection(connection *anet.Connection) {
	fmt.Println("=============建立链接的回调")
	a.conn = connection
}

func TestTCPClient(t *testing.T) {
	imp := &TCPClientImpl{}
	imp.conn = anet.NewNetClient("tcp://127.0.0.1:19999/", imp, 1000, true)
	go runTcpClientUsingCb(imp)
	time.Sleep(1 * time.Minute)
}
func runTcpClientUsingCb(imp *TCPClientImpl) {
	totalCount := 1000
	failCount := 0
	passCount := 0
	start := time.Now()
	for i := 0; i < totalCount; i++ {
		str := fmt.Sprintf("发送第:%d次数据", i)
		pack := &anet.PackHead{Cmd: 1, Body: []byte(str)}
		if _, err := imp.conn.WriteMessageWithCallback(pack, func(msg *anet.PackHead) {
			passCount += 1
		}); err != nil {
			failCount += 1
		}
	}
	timeCost := time.Now().Sub(start).Milliseconds()
	log.Printf("totalCount:%d,failCount:%d,passCount:%d, timeCost:%d, avg:%d", totalCount, failCount, passCount, timeCost, timeCost/int64(totalCount))
}

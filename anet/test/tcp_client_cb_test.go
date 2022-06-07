package test

import (
	"awesome/anet"
	"fmt"
	"log"
	"testing"
	"time"
)

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

/*
 * test result:
 * 	totalCount:10000,failCount:0,passCount:9999, timeCost:1519, avg:0.151900
 */
func TestTCPClient(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	go startClientCb()
	//go printResult()
	time.Sleep(1 * time.Minute)
}

func startClientCb() {
	tcpClientCbTestInfo.SetTotalCount = tcpClientCbTestInfo.SetCountEachThread * tcpClientCbTestInfo.ThreadCount
	for i := 0; i < tcpClientCbTestInfo.ThreadCount; i++ {
		imp := &TCPClientImpl{}
		imp.conn = anet.NewNetClient("tcp://127.0.0.1:19999", imp, 1000, true)
		go runTcpClientUsingCb(imp)
	}
}

var tcpClientCbTestInfo = &anet.TestInfo{ThreadCount: 4, SetCountEachThread: 5000000}

func cb(msg *anet.PackHead) {

	//lock.Lock()
	tcpClientCbTestInfo.Lock()
	tcpClientCbTestInfo.PassCount += 1
	tcpClientCbTestInfo.CurrentTotalCount += 1
	if tcpClientCbTestInfo.CurrentTotalCount == tcpClientCbTestInfo.SetTotalCount {
		log.Println("執行完成")
		tcpClientCbTestInfo.TimeCost = time.Now().Sub(tcpClientCbTestInfo.Start).Milliseconds()
		log.Printf("thead:%d,totalCount:%d,failCount:%d,passCount:%d, timeCost:%d, avg:%f", tcpClientCbTestInfo.ThreadCount,
			tcpClientCbTestInfo.CurrentTotalCount, tcpClientCbTestInfo.FailCount, tcpClientCbTestInfo.PassCount, tcpClientCbTestInfo.TimeCost, float64(tcpClientCbTestInfo.TimeCost)/float64(tcpClientCbTestInfo.CurrentTotalCount))
	}
	tcpClientCbTestInfo.Unlock()
	//lock.Unlock()
}
func runTcpClientUsingCb(imp *TCPClientImpl) {
	tcpClientCbTestInfo.Start = time.Now()
	for i := 0; i < tcpClientCbTestInfo.SetCountEachThread; i++ {
		str := fmt.Sprintf("发送第:%d次数据", i)
		pack := &anet.PackHead{Cmd: 1, Body: []byte(str)}
		pack.ReserveLow = uint32(i)
		if _, err := imp.conn.WriteMessageWithCallback(pack, cb); err != nil {
			tcpClientCbTestInfo.Lock()
			tcpClientCbTestInfo.FailCount += 1
			tcpClientCbTestInfo.Unlock()
		}
	}

}

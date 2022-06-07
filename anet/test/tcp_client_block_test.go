package test

import (
	"awesome/anet"
	"fmt"
	"log"
	"testing"
	"time"
)

type BlockTCPClientImpl struct {
	//cbNewConn aMQNewConnCallback
	conn *anet.Connection
}

func (a *BlockTCPClientImpl) IOnInit(connection *anet.Connection) {

}

func (a *BlockTCPClientImpl) IOnProcessPack(pack *anet.PackHead, connection *anet.Connection) {

}

/*
 * this interface SHOULD NOT CALL close.
 */
func (a *BlockTCPClientImpl) IOnClose(err error) (tryReconnect bool) {
	fmt.Println("IOnClose订阅连接关闭")
	return true
}

func (a *BlockTCPClientImpl) IOnConnect(isOk bool) {
	fmt.Println("=============建立链接的回调")

}

func (a *BlockTCPClientImpl) IOnNewConnection(connection *anet.Connection) {
	fmt.Println("=============建立链接的回调")
	a.conn = connection
}

/*
 * test tcpClient:
 * 	totalCount:10000,failCount:0,passCount:9999, timeCost:1519, avg:0.151900
 */
func TestBlockTCPClient(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	for i := 0; i < tcpClientBlockResult.ThreadCount; i++ {
		go worker()
	}
	time.Sleep(3 * time.Minute)
}

func worker() {
	imp := &TCPClientImpl{}
	imp.conn = anet.NewNetClient("tcp://127.0.0.1:19999", imp, 1000, true)
	tcpClientBlockResult.SetTotalCount = tcpClientBlockResult.ThreadCount * tcpClientBlockResult.SetCountEachThread
	tcpClientBlockResult.Start = time.Now()
	runBlockTask(imp)
}

var tcpClientBlockResult = &anet.TestInfo{Start: time.Now(), CurrentTotalCount: 0, ThreadCount: 20, SetCountEachThread: 100000}

func runBlockTask(imp *TCPClientImpl) {
	for i := 0; i < tcpClientBlockResult.SetCountEachThread; i++ {
		str := fmt.Sprintf("发送第:%d次数据", i)
		pack := &anet.PackHead{Cmd: 1, Body: []byte(str)}
		pack.ReserveLow = uint32(i)

		if _, isTimeout := imp.conn.WriteMessageWaitResponseWithinTimeLimit(pack, 100); isTimeout {
			tcpClientBlockResult.Lock()
			tcpClientBlockResult.FailCount += 1
			tcpClientBlockResult.CurrentTotalCount += 1
			tcpClientBlockResult.Unlock()
		} else {
			tcpClientBlockResult.Lock()
			tcpClientBlockResult.CurrentTotalCount += 1
			tcpClientBlockResult.PassCount += 1
			tcpClientBlockResult.Unlock()
		}
	}
	if tcpClientBlockResult.CurrentTotalCount == tcpClientBlockResult.SetTotalCount {
		tcpClientBlockResult.TimeCost = time.Now().Sub(tcpClientBlockResult.Start).Milliseconds()

		log.Printf("thead:%d,totalCount:%d,failCount:%d,passCount:%d, timeCost:%d, avg:%f", tcpClientBlockResult.ThreadCount, tcpClientBlockResult.CurrentTotalCount,
			tcpClientBlockResult.FailCount, tcpClientBlockResult.PassCount, tcpClientBlockResult.TimeCost, float64(tcpClientBlockResult.TimeCost)/float64(tcpClientBlockResult.CurrentTotalCount))
	}

}

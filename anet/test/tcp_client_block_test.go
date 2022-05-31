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

//func (a *subImpl) IWrite(msg interface{}, ph *net.PackHead){
//
//}

func (a *BlockTCPClientImpl) IOnConnect(isOk bool) {
	fmt.Println("=============建立链接的回调")

}

func (a *BlockTCPClientImpl) IOnNewConnection(connection *anet.Connection) {
	fmt.Println("=============建立链接的回调")
	a.conn = connection
}

/*
 * test result:
 * 	totalCount:10000,failCount:0,passCount:9999, timeCost:1519, avg:0.151900
 */
func TestBlockTCPClient(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	for i := 0; i < pSize; i++ {
		go worker()
	}
	go calcResult()
	//go printResult()
	time.Sleep(1 * time.Minute)
}
func worker() {
	imp := &TCPClientImpl{}
	imp.conn = anet.NewNetClient("tcp://127.0.0.1:19999", imp, 1000, true)
	runBlockTask(imp, 10000)
}

func calcResult() {
	total := &testDef{}
	for i := 0; i < pSize; i++ {
		item := <-resultChan
		total.TotalCount += item.TotalCount
		total.PassCount += item.PassCount
		total.TimeCost += item.TimeCost
		total.FailCount += item.FailCount
	}
	log.Printf("totalCount:%d,failCount:%d,passCount:%d, timeCost:%d, avg:%f", total.TotalCount,
		total.FailCount, total.PassCount, total.TimeCost, float64(total.TimeCost)/float64(total.TotalCount))
}

const pSize = 2

var resultChan = make(chan *testDef)

type testDef struct {
	TotalCount int
	FailCount  int
	PassCount  int
	TimeCost   int64
	Start      time.Time
}

//var (
//	blockTotalCount       = 10000
//	blockFailCount        = 0
//	blockPassCount        = 0
//	blockTimeCost   int64 = 0
//	blockStart      time.Time
//	blockLock       sync.Mutex
//)

func runBlockTask(imp *TCPClientImpl, c int) {
	result := &testDef{Start: time.Now(), TotalCount: c}
	for i := 0; i < c; i++ {
		str := fmt.Sprintf("发送第:%d次数据", i)
		pack := &anet.PackHead{Cmd: 1, Body: []byte(str)}
		pack.ReserveLow = uint32(i)
		if _, isTimeout := imp.conn.WriteMessageWaitResponseWithinTimeLimit(pack, 100); isTimeout {
			result.FailCount += 1
		} else {
			result.PassCount += 1
		}
	}
	result.TimeCost = time.Now().Sub(result.Start).Milliseconds()

	//log.Printf("totalCount:%d,failCount:%d,passCount:%d, timeCost:%d, avg:%f", blockTotalCount,
	//	blockFailCount, blockPassCount, blockTimeCost, float64(blockTimeCost)/float64(blockTotalCount))
	resultChan <- result
}

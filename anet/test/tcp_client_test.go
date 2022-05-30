package test

import (
	"awesome/anet"
	"fmt"
	"log"
	"sync"
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

	imp := &TCPClientImpl{}
	imp.conn = anet.NewNetClient("tcp://127.0.0.1:19999", imp, 1000, true)
	go runTcpClientUsingCb(imp)
	//go printResult()
	time.Sleep(1 * time.Minute)
}

var (
	totalCount       = 10000
	failCount        = 0
	passCount        = 0
	timeCost   int64 = 0
	start      time.Time
	lock       sync.Mutex
)

func cb(msg *anet.PackHead) {
	//log.Println("成功的回調")
	if msg.ReserveLow == uint32(totalCount-1) {
		log.Println("執行完成")
		timeCost = time.Now().Sub(start).Milliseconds()
		log.Printf("totalCount:%d,failCount:%d,passCount:%d, timeCost:%d, avg:%f", totalCount, failCount, passCount, timeCost, float64(timeCost)/float64(totalCount))
	}
	//lock.Lock()
	passCount += 1
	//lock.Unlock()
}
func runTcpClientUsingCb(imp *TCPClientImpl) {
	start = time.Now()
	for i := 0; i < totalCount; i++ {
		str := fmt.Sprintf("发送第:%d次数据", i)
		pack := &anet.PackHead{Cmd: 1, Body: []byte(str)}
		pack.ReserveLow = uint32(i)
		if _, err := imp.conn.WriteMessageWithCallback(pack, cb); err != nil {
			failCount += 1
		}
	}

}


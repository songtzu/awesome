package test

import (
	"awesome/anet"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestTcpSvr(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	impl := &svrImplement{}
	go anet.StartTcpSvr("127.0.0.1:19999", impl)
	time.Sleep(1 * time.Minute)
}

type svrImplement struct {
	//reliableCallback AMQCallback
	conn *anet.Connection
	id   int
}

func (a *svrImplement) IOnInit(connection *anet.Connection) {

}

func (a *svrImplement) IOnProcessPack(pack *anet.PackHead, connection *anet.Connection) {
	//log.Println("xmqPubImpl..IOnProcessPack.", string(pack.Body), pack)
	str := fmt.Sprintf("yes, we got:%s", string(pack.Body))
	pack.Body = []byte(str)
	connection.WriteMessage(pack)
}

/*
 * this interface SHOULD NOT CALL close.
 */
func (a *svrImplement) IOnClose(err error) (tryReconnect bool) {
	return true
}

func (a *svrImplement) IOnConnect(isOk bool) {

}

func (a *svrImplement) IOnNewConnection(connection *anet.Connection) {
	log.Println("new connection")
	a.conn = connection
	a.id = anet.GenNewId()

}

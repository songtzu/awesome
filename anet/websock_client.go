package anet

import (
	"golang.org/x/net/websocket"
	"log"
	"time"
)

func NewWebsockClientConnect( bindAddress string, iConn InterfaceNet,timeoutMillisecond int, endlessReconnect bool) *Connection {

	var c = &Connection{
		connectTimeoutMillisecond: timeoutMillisecond,
		connectAddr:               bindAddress,
		netProtocol:               NetProtocolTypeWebSock,
		iConn:                     iConn,
		state:                     ConnectionStateConnected,
		connectionType:            connectionTypeClient,
	}
	if err:=c.websockConnect(endlessReconnect);err!=nil{
		c.state = ConnectionStateClosed
		return nil
	}
	c.iConn.IOnConnect(true)
	return c
}
//
//func (c *Connection) websockReadLoop() {
//	if n, err = ws.Read(msg); err != nil {
//		fmt.Println(err)
//	}
//}

func (c *Connection)websockConnect(retryWhenFailed bool) (err error)  {
	//origin := "http://118.10.30.11/"
	//url := "ws://127.0.0.1:19999/ping"
	//origin:=""
	log.Println(c.connectAddr)
	ws, err := websocket.Dial(c.connectAddr, "", c.connectAddr)
	if err != nil {
		//alog.Fatal(err)
		log.Println(err,c.connectAddr)
		if retryWhenFailed{
			time.Sleep(1*time.Second)
			return c.websockConnect(retryWhenFailed)
		}
		return err
	}
	c.connWebSock = ws
	go webSockReadLoop(c)
	//for index:=0;index<100;index++{
	//	pack:=&net.PackHead{Body:[]byte("=================="+tag+"===================>hello, world!"+strconv.Itoa(index)),Cmd:uint32(1)}
	//	bin,_:=pack.SerializePackHead()
	//	if _, err := ws.Write(bin); err != nil {
	//		fmt.Println(err)
	//	}
	//	fmt.Println("ping发送",string(pack.Body))
	//	var msg = make([]byte, 512)
	//	var n int
	//	if n, err = ws.Read(msg); err != nil {
	//		fmt.Println(err)
	//	}
	//	fmt.Printf("Received: %s.\n", msg[:n])
	//	time.Sleep(1000*time.Millisecond)
	//}
	return nil
}

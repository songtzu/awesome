package anet

import (
	"log"
	"net"
	"time"

	"errors"
)

//
//
//func NewTcpClientConnection( iConn InterfaceNet ) *Connection {
//
//	var c = &Connection{
//		netProtocol: netProtocolTypeTCP,
//		iConn:    iConn,
//	}
//	return c
//}



func NewTcpClientConnect( bindAddress string, iConn InterfaceNet,timeoutMillisecond int, endlessReconnect bool) *Connection {

	var c = &Connection{
		connectTimeoutMillisecond:timeoutMillisecond,
		connectAddr:bindAddress,
		netProtocol: NetProtocolTypeTCP,
		iConn:    iConn,
		state:ConnectionStateConnected,
		connectionType:connectionTypeClient,
	}

	c.iConn.IOnInit(c)
	if err:=c.tcpConnect(endlessReconnect);err!=nil{
		c.state = ConnectionStateClosed
		return nil
	}


	c.iConn.IOnConnect(true)
	return c
}

func (c *Connection) tcpConnect(retryWhenFailed bool) (err error) {
	rw, err := net.DialTimeout("tcp", c.connectAddr, time.Duration(c.connectTimeoutMillisecond)*time.Millisecond)
	if err != nil {

		msg:= "failed tcpConnect to " +  c.connectAddr +", with err:%s" + err.Error()
		log.Println(msg)
		if retryWhenFailed{
			time.Sleep(3*time.Second)
			return c.tcpConnect(retryWhenFailed)
		}
		return errors.New(msg)
	}else {
		c.connTcp = rw
	}
	c.startTcpClient()
	return nil
}

func (c *Connection)startTcpClient() bool {
	go c.tcpReadLoop()
	return true
}


//func (c *Connection)WriteBytes(cmd int,data []byte) (n int ,err error) {
//	pack:=&PackHead{magicNum,cmd,0,len(data),0,0,data}
//	len,err:=c.connTcp.Write(data)
//	c.WriteBytes()
//	return len,err
//}

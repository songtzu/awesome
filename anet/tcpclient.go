package anet

import (
	"fmt"
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

func NewTcpClientConnect(bindAddress string, iConn InterfaceNet, timeoutMillisecond int, endlessReconnect bool) *Connection {

	var c = &Connection{
		connectTimeoutMillisecond: timeoutMillisecond,
		connectAddr:               bindAddress,
		netProtocol:               NetProtocolTypeTCP,
		iConn:                     iConn,
		state:                     ConnectionStateConnected,
		connectionType:            connectionTypeClient,
	}

	c.iConn.IOnInit(c)
	if err := c.tcpConnect(endlessReconnect); err != nil {
		c.state = ConnectionStateClosed
		return nil
	}

	//c.iConn.IOnConnect(true)
	return c
}

/*
 * 程序启动第一次连接，断开的原因肯定是配置或服务器未启动，此时连接失败，不应当重试，retryWhenFailed为false
 *	断开连接，重连，应当不停重拾。retryWhenFailed为true
 */
func (c *Connection) tcpConnect(retryWhenFailed bool) (err error) {
	c.retryCount += 1
	log.Println("Connection, tcpConnect, ", retryWhenFailed)
	rw, err := net.DialTimeout("tcp", c.connectAddr, time.Duration(c.connectTimeoutMillisecond)*time.Millisecond)
	if err != nil {
		msg := fmt.Sprintf("failed tcpConnect to :%s,  with err:%s, retryCount:%d", c.connectAddr, err.Error(), c.retryCount)
		log.Println(msg)
		if retryWhenFailed {
			time.Sleep(3 * time.Second)
			return c.tcpConnect(retryWhenFailed)
		}
		//第一次启动，失败，直接返回错误，程序应当退出。
		return errors.New(msg)
	} else {
		c.connTcp = rw
		c.connSucceedCount += 1
	}
	//
	c.startTcpClient()
	c.state = ConnectionStateConnected
	c.iConn.IOnConnect(c.connSucceedCount > 1)
	return nil
}

func (c *Connection) startTcpClient() bool {
	log.Println("Connection, startTcpClient")
	go c.tcpReadLoop()
	return true
}

//func (c *Connection)WriteBytes(cmd int,data []byte) (n int ,err error) {
//	pack:=&PackHead{magicNum,cmd,0,len(data),0,0,data}
//	len,err:=c.connTcp.Write(data)
//	c.WriteBytes()
//	return len,err
//}

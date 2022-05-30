package anet

import "log"

/*
 * 提供ws和tcp协议细节的屏蔽封装客户端。
 */

func NewNetClient(bindAddress string, iConn InterfaceNet, timeoutMillisecond int, endlessReconnect bool) *Connection {
	addr, protoType := ParseNetIpFromAddressWithProtocol(bindAddress)
	log.Println(addr, protoType)
	if protoType == NetProtocolTypeTCP {
		return NewTcpClientConnect(addr, iConn, timeoutMillisecond, endlessReconnect)
	} else if protoType == NetProtocolTypeWebSock {
		log.Println("ws协议")
		return NewWebsockClientConnect(bindAddress, iConn, timeoutMillisecond, endlessReconnect)
	}
	return nil
}

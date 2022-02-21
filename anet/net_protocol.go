package anet

import (
	"fmt"
	"log"
	"strings"
)

const (
	ProtoTcp = "tcp://"
	ProtoWebSock = "ws://"
)

type NetProtocolType uint32
const (
	NetProtocolTypeTCP       NetProtocolType = 1
	NetProtocolTypeWebSock   NetProtocolType = 2
	NetProtocolTypeUndefined NetProtocolType = 2
)
func ParseNetIpFromAddressWithProtocol(address string) (_addr string, protocolType NetProtocolType) {
	if strings.HasPrefix(address,ProtoTcp){
		_addr = strings.Replace(address,ProtoTcp,"",1)
		return _addr,NetProtocolTypeTCP
	}else if strings.HasPrefix(address,ProtoWebSock){

		_addr = strings.Replace(address,ProtoWebSock,"",1)
		return _addr , NetProtocolTypeWebSock
	}
	log.Println(fmt.Sprintf("address must be something like '%s', or '%s' we got %s",ProtoTcp,ProtoWebSock,address))
	return "", NetProtocolTypeUndefined
}
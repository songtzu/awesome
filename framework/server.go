package framework

import (
	"awesome/alog"
	"awesome/config"
	"awesome/anet"
)

func StartSvr(instance IFramework)  {
	InitFrameworkInstance(instance)
	impl:=&PlayerImpl{}
	alog.Debug("=======",config.GetConfig().Server.BindAddress)
	addr,protoType:= anet.ParseNetIpFromAddressWithProtocol(config.GetConfig().Server.BindAddress)
	if protoType== anet.NetProtocolTypeTCP{
		anet.StartTcpSvr(addr,impl)
	}else if protoType== anet.NetProtocolTypeWebSock{
		anet.StartWebSocket(addr,impl)
	}
}



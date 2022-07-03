package framework

import (
	"awesome/alog"
	"awesome/anet"
	"awesome/config"
	"log"
	"os"
)

func StartSvr(instance IFramework)  {
	InitFrameworkInstance(instance)
	impl:=&PlayerImpl{}
	alog.Debug("=======",config.GetConfig().Server.BindAddress)
	if config.GetConfig().Server.IsConnStart{
		addr,protoType:= anet.ParseNetIpFromAddressWithProtocol(config.GetConfig().Server.BindAddress)
		if protoType== anet.NetProtocolTypeTCP{
			anet.StartTcpSvr(addr,impl)
		}else if protoType== anet.NetProtocolTypeWebSock{
			anet.StartWebSocket(addr,impl)
		}
	}

}

func StartHttp(instance IFramework)  {
	log.Println("===",config.GetConfig().Server.IsHttpStart)
	if config.GetConfig().Server.IsHttpStart{
		err := StartEchoServer( config.GetConfig().Server.HttpAddress )
		if err!=nil{
			log.Printf("http server ,add:%s启动失败:%s", config.GetConfig().Server.HttpAddress, err.Error())
			os.Exit(-10)
		}
		instance.OnRegisterHttpRouters(echoInstance)
		log.Printf("http server start, %s",config.GetConfig().Server.HttpAddress)
	}
}

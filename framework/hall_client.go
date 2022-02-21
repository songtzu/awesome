package framework

import (
	"awesome/alog"
	"awesome/anet"
	"awesome/config"
	"awesome/message"
	"encoding/json"
)


var hallImplInstance *hallClientImpl

func StartHallSession() {
	alog.Debug("创建对大厅的会话")
	go startUpHallSession()
}
func startUpHallSession() bool {
	url := config.GetConfig().Hall.HallAddress
	hallImplInstance = &hallClientImpl{}
	alog.Debug("链接的URL ", url)
	hallImplInstance.conn = anet.NewNetClient(url, hallImplInstance, 1000, true)
	if hallImplInstance.conn == nil {
		alog.Err("无法和网关建立连接")
		return true
	}
	return registSelf()
}
func registSelf() bool {
	regist := message.MsgRegistService{Status: message.SeviceStatusNormal,
		ServerId: config.GetConfig().Server.ServerID,
		AppId:config.GetConfig().Server.AppID,
		BindAddress:config.GetConfig().Server.BindAddress,
		Version:config.GetConfig().Server.Version}

	bin, _ := json.Marshal(regist)
	pack := &anet.PackHead{Cmd: message.InnerCmdRegistService, Body: bin}
	isTimeout, _ := hallImplInstance.conn.WriteMessageWaitResponseWithinTimeLimit(pack, 5000)

	if isTimeout {
		alog.Err("注册服务")
		return false
	}
	return true
}

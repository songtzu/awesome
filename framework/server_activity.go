package framework

import (
	"awesome/alog"
	"awesome/anet"
	"awesome/config"
	"awesome/localcache"
	"awesome/message"
	"encoding/json"
	"fmt"
)

func activity(pack *anet.PackHead) ( *anet.PackHead) {

	respMsg:=&message.MsgProxyActivityAck{}
	msg:=&message.MsgProxyActivity{}
	if err:=json.Unmarshal(pack.Body,msg); err!=nil{
		alog.Err(fmt.Sprintf("proxy激活请求解析失败,原始数据%s,错误提示:%s",string(pack.Body),err.Error()))
		respMsg.Status = 1
		bin,_:=json.Marshal(respMsg)
		pack.Body = bin

		return pack
	}
	if len(msg.ProxyAddress) ==0{
		alog.Err(fmt.Sprintf("proxy激活请求解析失败,原始数据%s,网关地址为空 ",string(pack.Body) ))
		respMsg.Status = 2
		bin,_:=json.Marshal(respMsg)
		pack.Body = bin
		return pack
	}
	localcache.InsertProxyInfo(msg.ProxyAddress)
	respMsg.Status = 0
	respMsg.AppId = config.GetConfig().Server.AppID
	respMsg.BindAddress = config.GetConfig().Server.BindAddress
	respMsg.ServerId =config.GetConfig().Server.ServerID
	respMsg.Version =config.GetConfig().Server.Version
	bin,_:=json.Marshal(respMsg)
	pack.Body = bin
	alog.Debug("正常完成服务激活")
	return pack
}
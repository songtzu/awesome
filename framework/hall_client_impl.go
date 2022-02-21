package framework

import (
	"awesome/alog"
	"awesome/anet"
	"awesome/config"
	"awesome/message"
	"encoding/json"
	"time"
)

type hallClientImpl struct {
	conn *anet.Connection
}

func (a *hallClientImpl) IOnInit(connection *anet.Connection) {
	a.conn = connection
}

func (a *hallClientImpl) IOnProcessPack(pack *anet.PackHead) {
	if pack.Cmd == message.InnerCmdPingAck {
		return
	}

	alog.Info("大厅消息派发", pack.Cmd)

}

/*
 * this interface SHOULD NOT CALL close.
 */
func (a *hallClientImpl) IOnClose(err error) (tryReconnect bool) {
	alog.Info("大厅的会话中断")
	return true
}

func (a *hallClientImpl) IOnConnect(isOk bool) {
	go a.ping()
}

func (a *hallClientImpl) IOnNewConnection(connection *anet.Connection) {
	alog.Info("创建连接")

}

//不停的ping，如果有心跳超时，则重连
func (a *hallClientImpl) ping() {
	for {
		ping := &message.MsgPing{IsOk: true, Timestamp: time.Now().Unix()}
		bin, _ := json.Marshal(ping)
		pack := &anet.PackHead{Cmd: message.InnerCmdPing, Body: bin}
		isTimeOut, _ := a.conn.WriteMessageWaitResponseWithinTimeLimit(pack, 5000)
		//alog.Info(ack)
		if isTimeOut {
			alog.Info("大厅会话的心跳消息超时，即将重连")
			a.conn.CloseConnWithoutRecon(nil)
			return
		}
		time.Sleep(time.Duration(config.GetConfig().Hall.HallHeartBeatInterval) * time.Second)
	}

}

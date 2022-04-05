package framework

import (
	"awesome/config"
	"google.golang.org/protobuf/proto"

	"awesome/alog"
	"awesome/anet"
	"awesome/defs"
	"sync"
)

const newestCachedMessageSize = 100

type PlayerImpl struct {
	mux sync.Mutex

	conn *anet.Connection
	//房间的句柄
	room *Room
	/*
	 * cache newest message dispatch to player.
	 * 		it should be used for logger & some other usages.
	 */
	newestCachedMessages []defs.TypeCmd
	userId               defs.TypeUserId
	userData             interface{} //
}

var undefinedUserId defs.TypeUserId = 0

func newUndefinedUserId() defs.TypeUserId {
	undefinedUserId--
	return undefinedUserId
}
func (p *PlayerImpl) IOnInit(connection *anet.Connection) {

}

func (p *PlayerImpl) IOnProcessPack(pack *anet.PackHead) {
	p.recordNewestMessage(pack)
	userMsg := &UserMessage{pack: pack, user: p}
	if pack.Cmd == config.GetConfig().ActiveCmd {
		alog.Trace("激活服务器")
		if _, err := p.conn.WriteMessage(activity(pack)); err != nil {
			alog.Err("服务激活请求写回数据失败", err.Error())
		}
		return
	}
	if config.GetConfig().IsCmdInsideIgnoreList(pack.Cmd) {
		//alog.Debug("系统配置心跳命令，忽略")
		_, err := p.conn.WriteMessage(pack)
		if err != nil {
			alog.Info("写数据失败", p.userId)
		}
		return
	}
	if p.room != nil {
		//写入chan，
		alog.Debug("normal room redirect.")
		p.room.enqueueMessage(userMsg)
	} else {
		alog.Debug("玩家未关联房间，进入创建房间的流程")
		//user not associate with a room. redirect it to special room
		specialRoomInstance.enqueueMessage(userMsg)
	}
	//alog.Trace(p.conn)
	//alog.Trace("原样回包",string(pack.Body))
	//if p.conn.
	//p.conn.WriteMessage(pack)
}

func (p *PlayerImpl) SendBinary(body []byte, cmd int) (len int, err error) {
	return p.conn.WriteBytes(body, uint32(cmd))
}

func (p *PlayerImpl) SendPackage(pack *anet.PackHead) (len int, err error) {
	//alog.Debug("写出二进制")
	return p.conn.WriteMessage(pack)
}
func (p *PlayerImpl) SendPackageWithCallback(pack *anet.PackHead, cb anet.DefNetIOCallback) (len int, err error) {
	//alog.Debug("写出二进制")
	return p.conn.WriteMessageWithCallback(pack, cb)
}

func (p *PlayerImpl) IOnConnect(isOk bool) {

}

/*
 * this interface SHOULD NOT CALL close.
 */
func (p *PlayerImpl) IOnClose(err error) (tryReconnect bool) {
	return false
}
func (p *PlayerImpl) IOnNewConnection(connection *anet.Connection) {
	//alog.Debug("IOnNewConnection回调")
	//connection.PrintNetProtocol()
	p.conn = connection
	p.userId = newUndefinedUserId()
	p.newestCachedMessages = []defs.TypeCmd{}
	p.room = nil
}

func (p *PlayerImpl) SetUserData(extension interface{}) {
	p.userData = extension
}

func (p *PlayerImpl) GetUserData() interface{} {
	return p.userData
}

func (p *PlayerImpl) SetUserId(userId defs.TypeUserId) {
	if p.room != nil {
		//player has been associate with a room, in this case, update playerMap
		if p.room.players.playerExist(p.userId) {
			p.room.players.playerDelete(p.userId)
		}
	}
	p.userId = userId
	p.room.players.playerSet(p.userId, p)
}

func (p *PlayerImpl) GetUserId() (userId defs.TypeUserId) {
	userId = p.userId
	return
}

func (p *PlayerImpl) recordNewestMessage(pack *anet.PackHead) {
	if len(p.newestCachedMessages) > newestCachedMessageSize {
		p.newestCachedMessages = p.newestCachedMessages[1:]
	}
	p.newestCachedMessages = append(p.newestCachedMessages, defs.TypeCmd(pack.Cmd))
}

func (p *PlayerImpl) newestMessage() defs.TypeCmd {
	if len(p.newestCachedMessages) == 0 {
		return 0
	}
	return p.newestCachedMessages[len(p.newestCachedMessages)-1]
}

func (p *PlayerImpl) AddToRoom(roomCode defs.RoomCode) (isOk bool) {
	if room := roomMapGet(roomCode); room != nil {
		room.AddPlayerToRoom(p)
		return true
	}
	return false
}

func (p *PlayerImpl) SendMsg(cmd int, msg interface{}, seq uint32) (err error) {

	bin, err := proto.Marshal(msg.(proto.Message))
	if err != nil {
		alog.Err("proto marshal failed", msg)
		return
	}

	p.conn.WriteSeqBytes(bin, uint32(cmd),seq)
	return
}

package framework

import (
	"awesome/anet"
	"awesome/defs"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

func GetRoomData(roomId uint64) (extension interface{}) {
	r := roomMapGet(defs.RoomCode(roomId))
	if r == nil {
		return nil
	}
	return r.GetRoomData()
}

func GetRoomList() (list []interface{}) {
	list = make([]interface{}, 0)
	roomMap.Range(func(roomId, room interface{}) bool {
		if r, ok := room.(*Room); ok && r != nil {
			list = append(list, r.GetRoomData())
		} else {
			log.Println("roomId:%v not Room   val:%v", roomId, room)
		}
		return true
	})
	return
}

func SendUserMsg(player *PlayerImpl, cmd int, msg interface{}) error {
	if player == nil {
		return errors.New(fmt.Sprintf("player is nil"))
	}
	player.SendMsg(cmd, msg)
	return nil
}

func SendMsg(uid defs.TypeUserId, cmd int, msg interface{}) error {
	u := UserMapGet(uid)
	if u == nil {
		return fmt.Errorf("userid %d not found ", uid)
	}
	u.SendMsg(cmd, msg)
	return nil
}

func SendUserMsgWithId(inviteCode defs.RoomCode, uid defs.TypeUserId, cmd int, msg interface{}) error {
	r := roomMapGet(inviteCode)
	if r == nil {
		return fmt.Errorf("SendUserMsgWithId: room not found, code=%d", inviteCode)
	}

	player := r.GetPlayerById(uid)
	if player == nil {
		return fmt.Errorf("SendUserMsgWithId: user not found, code=%d, user=%d", inviteCode, uid)
	}
	player.SendMsg(cmd, msg)
	return nil
}

func GetRoomClients(inviteCode int) (map[int]interface{}, error) {
	r := roomMapGet(defs.RoomCode(inviteCode))
	if r == nil {
		return nil, errors.New(fmt.Sprintf("not fount %v room", inviteCode))
	}
	var data = make(map[int]interface{}, 16)
	r.Players.Range(func(k, v interface{}) bool {
		data[int(k.(defs.TypeUserId))] = v.(*PlayerImpl).GetUserData()
		return true
	})

	return data, nil
}

func RoomDelClient(roomid defs.RoomCode, userid defs.TypeUserId) int {
	room := roomMapGet(roomid)
	if room == nil {
		log.Println("DelClient is err:%d", -1)
		return -1
	}
	room.DelPlayFromRoomById(userid)
	return 0
}

func RoomAddClient(roomid defs.RoomCode, oldUid, newUid defs.TypeUserId, data interface{}) int {
	//glog.Infoln("RoomAddClient", roomid, oldUid, newUid)
	room := roomMapGet(roomid)
	if room == nil {
		log.Println("room not found:%d", roomid, oldUid, newUid)
		return -1
	}

	//todo, 惰断开情景，把原uid-->playerinfo里面的conn.iConn置空,避免后面EOF又抛出短线消息处理.
	deadSession := UserMapGet(newUid)
	if deadSession!=nil{
		deadSession.conn.ResetIConnToNil()
	}

	session := UserMapGet(oldUid)
	if session == nil {
		log.Println("get user session fail:", roomid, oldUid, newUid)
		return -1
	}

	//更新
	userSessionUpdate(oldUid, newUid, session)
	session.room = room
	session.userId = newUid
	session.userData = data
	room.AddPlayerToRoom(session)
	//glog.Infoln("添加进房间的地址：", session.userId, session)
	return 0
}

func RoomSetClientData(inviteCode, uid int, data interface{}) int {
	room := roomMapGet(defs.RoomCode(inviteCode))
	if room == nil {
		log.Println("room not found:%d", inviteCode)
		return -1
	}
	val, ok := room.Players.Load(uid)
	if !ok {
		log.Println("player not found:%d, %d", inviteCode, uid)
		return -1
	}
	val.(*PlayerImpl).userData = data
	return 0
}

func RoomGetClientData(roomId defs.RoomCode, uid int) (interface{}, error) {
	if roomId <= 0 || uid <= 0 {
		return nil, fmt.Errorf("invitecode or uid <=0 (uid:%d, inviteCode:%d)", uid, roomId)
	}

	r := roomMapGet(defs.RoomCode(roomId))
	if r == nil {
		return nil, fmt.Errorf("inviteCode=%d not found", roomId)
	}
	//glog.Infoln("RoomGetClientData:", r.RoomCode, uid, r)
	player := r.GetPlayerById(defs.TypeUserId(uid))
	if player == nil {
		return nil, fmt.Errorf("inviteCode=%d,user=%d not found", roomId, uid)
	}
	return player.userData, nil
}

func CloseUserSession(uid defs.TypeUserId) {
	defer UserMapDelete(uid)
	session := UserMapGet(uid)
	if session == nil {
		return
	}
	if session.conn != nil {
		session.conn.CloseConnWithoutRecon(nil)
	}
}

func SendHallMsgSyncTimeout(cmd uint32, msg interface{}, timeout time.Duration) (*anet.PackHead, error) {
	//return sendHallMsgSyncTimeout(cmd, msg, timeout)
	//todo, jack
	return nil,nil
}

func SendHallMsgAsy(cmd uint32, msg interface{}) {
	//sendHallMsgAsy(cmd, msg)
	//todo, jack
}

// GeneralMapGet ---------全局map操作
func GeneralMapGet(key string) interface{} {
	return GlobalMapGet(key)
}
func GeneralMapSet(key string, val interface{}) {
	GlobalMapSet(key, val)
}
func GeneralMapDelete(key string) {
	GlobalMapDelete(key)
}


func DeleteAndCloseRoom(inviteCode defs.RoomCode) error {
	r := roomMapGet(inviteCode)
	if r == nil {
		return fmt.Errorf("not fount %d room", inviteCode)
	}
	r.Close()
	return nil
}

//更新用户session
func userSessionUpdate(oldUid, newUid defs.TypeUserId, user *PlayerImpl) {
	UserMapDelete(oldUid)
	UserMapSet(newUid, user)
}




func SerializePackWithPB(ph *anet.PackHead, msg interface{}) (err error) {
	// TODO 加密
	switch v := msg.(type) {
	case []byte:
		ph.Body = v
		return
	case proto.Message:
		if ph.Body, err = proto.Marshal(v); err == nil {
			return
		} else {
			log.Println(fmt.Sprintf("proto marshal cmd: %d error: %v", ph.Cmd, err))
			return err
		}
	case types.Nil:
		ph.Body = nil
		return
	default:
		return errors.New("tcp_conn: error msg type")
	}
	//return errors.New("unkonw error while serialize proto")
}
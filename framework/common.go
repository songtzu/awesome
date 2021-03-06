package framework

import (
	"awesome/anet"
	"awesome/defs"
	"code.google.com/p/go.tools/go/types"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
)

func GetRoomData(roomId defs.RoomCode) (extension interface{}) {
	r := roomMapGet(roomId)
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
			log.Printf("roomId:%v not Room val:%v", roomId, room)
		}
		return true
	})
	return
}

func SendUserMsg(player *PlayerImpl, cmd defs.TypeCmd, msg interface{}) error {
	if player == nil {
		return fmt.Errorf("player is nil while send cmd:%d", cmd)
	}
	return player.SendMsg(cmd, msg, 0)
}

func SendBinaryMsg(player *PlayerImpl, cmd defs.TypeCmd, binary []byte, sequenceId uint32) (err error) {
	if player == nil {
		return fmt.Errorf("player is nil while send cmd:%d", cmd)
	}
	_, err = player.SendBinary(binary, cmd, sequenceId)
	return err
}

func SendMsg(uid defs.TypeUserId, cmd defs.TypeCmd, msg interface{}) error {
	u := UserMapGet(uid)
	if u == nil {
		return fmt.Errorf("userid %d not found ", uid)
	}
	return u.SendMsg(cmd, msg, 0)
}

func SendUserMsgWithId(inviteCode defs.RoomCode, uid defs.TypeUserId, cmd defs.TypeCmd, msg interface{}) error {
	r := roomMapGet(inviteCode)
	if r == nil {
		return fmt.Errorf("SendUserMsgWithId: room not found, code=%d", inviteCode)
	}

	player := r.GetPlayerById(uid)
	if player == nil {
		return fmt.Errorf("SendUserMsgWithId: user not found, code=%d, user=%d", inviteCode, uid)
	}
	return player.SendMsg(cmd, msg, 0)
}

func GetRoomClients(inviteCode int) (map[int]interface{}, error) {
	r := roomMapGet(defs.RoomCode(inviteCode))
	if r == nil {
		return nil, errors.New(fmt.Sprintf("not fount %v room", inviteCode))
	}
	var data = make(map[int]interface{}, 16)
	r.players.Range(func(k, v interface{}) bool {
		data[int(k.(defs.TypeUserId))] = v.(*PlayerImpl).GetUserData()
		return true
	})

	return data, nil
}

func RoomDeletePlayer(roomId defs.RoomCode, userid defs.TypeUserId) int {
	room := roomMapGet(roomId)
	if room == nil {
		log.Printf("room:%d,deleted userid :%d err:%d", roomId, userid, -1)
		return -1
	}
	room.DelPlayFromRoomById(userid)
	return 0
}

func RoomAddPlayer(roomCode defs.RoomCode, oldUid, newUid defs.TypeUserId, data interface{}) int {
	room := roomMapGet(roomCode)
	if room == nil {
		log.Printf("room not found:%d, oldUid:%d, newUid:%d", roomCode, oldUid, newUid)
		return -1
	}

	//todo, ????????????????????????uid-->playerinfo?????????conn.iConn??????,????????????EOF???????????????????????????.
	deadSession := UserMapGet(newUid)
	if deadSession != nil {
		deadSession.conn.ResetIConnToNil()
	}

	session := UserMapGet(oldUid)
	if session == nil {
		log.Println("get user session fail:", roomCode, oldUid, newUid)
		return -1
	}

	//??????
	userSessionUpdate(oldUid, newUid, session)
	session.room = room
	session.userId = newUid
	session.userData = data
	room.AddPlayerToRoom(session)
	return 0
}

func RoomSetClientData(inviteCode, uid int, data interface{}) int {
	room := roomMapGet(defs.RoomCode(inviteCode))
	if room == nil {
		log.Printf("room not found:%d", inviteCode)
		return -1
	}
	val, ok := room.players.Load(uid)
	if !ok {
		log.Printf("player not found:%d, %d", inviteCode, uid)
		return -1
	}
	val.(*PlayerImpl).userData = data
	return 0
}

func RoomGetPlayerData(roomId defs.RoomCode, uid int) (interface{}, error) {
	if roomId <= 0 || uid <= 0 {
		return nil, fmt.Errorf("invitecode or uid <=0 (uid:%d, inviteCode:%d)", uid, roomId)
	}

	r := roomMapGet(defs.RoomCode(roomId))
	if r == nil {
		return nil, fmt.Errorf("inviteCode=%d not found", roomId)
	}
	player := r.GetPlayerById(defs.TypeUserId(uid))
	if player == nil {
		return nil, fmt.Errorf("inviteCode=%d,user=%d not found", roomId, uid)
	}
	return player.userData, nil
}

func ClosePlayerSession(uid defs.TypeUserId) {
	defer UserMapDelete(uid)
	session := UserMapGet(uid)
	if session == nil {
		return
	}
	if session.conn != nil {
		session.conn.CloseConnWithoutRecon(nil)
	}
}

// GeneralMapGet ---------??????map??????
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

//????????????session
func userSessionUpdate(oldUid, newUid defs.TypeUserId, user *PlayerImpl) {
	UserMapDelete(oldUid)
	UserMapStore(newUid, user)
}

func SerializePackWithPB(ph *anet.PackHead, msg interface{}) (err error) {
	// TODO ??????
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

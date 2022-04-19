package framework

import (
	"awesome/alog"
	"awesome/defs"
	"fmt"
	"log"
)

const (
	RoomCodeNil = -1
	RoomCodeMatch = -999
)
/*
 * 特殊的房间
 * 		写入没有关联房间的玩家的消息。
 */
type specialRoom struct {
	*Room
}

func (s *specialRoom) matchWorker(userId defs.TypeUserId, matchData *MatchRule, msg *UserMessage) {
	r := roomMapGet(RoomCodeMatch)
	if r == nil {
		extension, err := frameworkInterfaceInstance.OnCreateRoom(msg.pack)
		if err != nil {
			alog.Err("error:", defs.GetError(defs.ErrorDefFailedCreateRoom))
			return
		}
		r, _ = createRoom(RoomCodeMatch, extension)
		msg.user.room = r
		//redirect to normal message chan.
		r.workerChan <- msg
	} else {
		alog.Debug("空房间，仍然在map找到该房间号", RoomCodeMatch, "重定向到该房间的chan中")
		r.workerChan <- msg
	}
}

func (s *specialRoom) specialWorkerForNilRoom(msg *UserMessage) {
	alog.Debug("create room:%d cmd:%d start a worker", s.RoomCode, msg.pack.Cmd)
	if matchData, userId:=frameworkInterfaceInstance.OnParseMatch(msg.pack);matchData!=nil{
		log.Printf("user:%d进入匹配流程:%d，deadline:%d", userId, matchData.MatchNum, matchData.DeadlineTimestamp)
		s.matchWorker(userId, matchData, msg)
		return
	}

	roomCode, userId, err := frameworkInterfaceInstance.OnParseUser(msg.pack)
	if err != nil {
		log.Printf("OnParseRoomCodeAndUser--->error:%s, userId:%d, roomCode:%d", defs.GetError(defs.ErrorDefFailedParseRoomCode), userId, roomCode)
		return
	}

	makeUserConnReady(msg.user, userId)
	r := roomMapGet(roomCode)
	if r == nil {
		extension, err := frameworkInterfaceInstance.OnCreateRoom(msg.pack)
		if err != nil {
			alog.Err("error:", defs.GetError(defs.ErrorDefFailedCreateRoom))
			return
		}
		r, _ = createRoom(roomCode, extension)
		msg.user.room = r
		//redirect to normal message chan.
		r.workerChan <- msg
	} else {
		alog.Debug("空房间，仍然在map找到该房间号", roomCode, "重定向到该房间的chan中")
		r.workerChan <- msg
	}

	return
}

func (s *specialRoom) specialRoomWorker() {
	s.taskWorker(func() {
		for {
			select {
			case msg := <-s.workerChan:
				alog.Debug("special room got message for normal..")
				if msg.user.room == nil {
					recoverWorker(func() {
						s.specialWorkerForNilRoom(msg)
					})
				} else {
					alog.Debug("非空房间，不处理，等房间chan执行")
				}

			case msg := <-s.sysMsg:
				alog.Debug("special room got message for system message.")
				//如果是空房間的話 系統消息不派發到房間，如果設置了超時回調會派發
				//消息到達空房間
				if msg.Cmd == SystemMessageDefTimeOut && msg.DealHandle != nil {
					recoverWorker(func() {
						//msg.DealHandle.(data.HandleTimeout)(msg.Msg, c.GetRoomData())
					})
				}
			}
		}
	})
}

var specialRoomInstance *specialRoom = nil

func init() {
	specialRoomInstance = &specialRoom{NewRoom(RoomCodeNil)}
	go specialRoomInstance.specialRoomWorker()
}

func createRoom(roomCode defs.RoomCode, extension interface{}) (*Room, error) {
	alog.Debug(fmt.Sprintf("创建一个新房间[%d],userData:%v", roomCode, extension))
	r := NewRoom(roomCode)
	r.roomData = extension
	r.recoverWorker()
	roomMapSet(roomCode, r)
	return r, nil
}

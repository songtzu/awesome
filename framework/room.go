package framework

import (
	"awesome/defs"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"sync"
	"sync/atomic"
	"time"
)



const NilRoomHelperInviteCode = 0
const CreateRoomHelperInviteCode = -1
const defaultWorkerSize = 512
const defaultSysMsgSize = 1024
const writeTimeOut = 5 * time.Second
const specialChanSize = 1024 * 10

const (
	RoomStatusRunning  int32 = 0
	RoomStatusStopped int32 = 1
)

// 房间数据设置为消息派发方式获取(机器人可通过方法获取)
// 不应该出现同时获取房间数据的情况(目前获取数据lua有提供的方法 用于机器人，关闭机器人功能时，要把方法屏蔽)
type Room struct {
	//房间数据
	roomData interface{} `json:"userData"` // 扩展项，房间数据

	//players map[int]*ClientSession `json:"-"` // 房间玩家,座位号作为map的键值
	players *Players `json:"players"`

	RoomCode defs.RoomCode `json:"inviteCode"` 			//房间邀请馬，此可作为房间标志

	workerChan    chan *UserMessage    			// 用户产生的消息
	sysMsg   chan *SystemMessage 				// 系统产生的消息

	runFlag int32
	mutex   sync.RWMutex
}

//var LogicEngine =
func (r *Room) GetRoomData() interface{} {
	return r.roomData
}
func (r *Room) SetRoomData(data interface{})  {
	r.roomData = data
}
// inviteCode <= 0 时为系统预设房间(创建房间，空房间)
func NewRoom(roomCode defs.RoomCode) *Room {

	room := &Room{
		RoomCode: roomCode,
		players:  &Players{},
		sysMsg:   make(chan *SystemMessage, defaultSysMsgSize),
	}
	if roomCode <= 0 { // 系统预设房间
		room.workerChan = make(chan *UserMessage, specialChanSize)
	} else {
		room.workerChan = make(chan *UserMessage, defaultWorkerSize)
	}
	return room
}



func (r *Room) AddPlayerToRoom(player *PlayerImpl) bool {
	r.players.Store(player.userId, player)
	player.room = r
	return true
}

func(r *Room) DelPlayerFromRoom(player *PlayerImpl) bool {
	player.room = nil
	r.players.Delete(player.userId)
	return true
}

func (r *Room) DelPlayFromRoomById(userId defs.TypeUserId) bool {
	playerImp, ok := r.players.playerGet(userId)
	if !ok || playerImp == nil {
		return false
	}
	playerImp.room = nil
	r.players.Delete(userId)
	return true
}


func (r *Room)GetPlayerById(userId defs.TypeUserId) (value *PlayerImpl) {
	playerImp, ok := r.players.playerGet(userId)
	if ok {
		return nil
	}
	value = playerImp
	return
}


func Broadcasts(roomCode defs.RoomCode, cmd int, ph proto.Message) error{
	room := roomMapGet(roomCode)
	if room == nil {
		return errors.New(fmt.Sprintf("not fount %v room", roomCode))
	}

	room.players.Range(func(uid, p interface{}) bool {
		p.(*PlayerImpl).SendMsg(cmd, ph, 0)
		return true
	})

	return nil
}

func Broadcast_(roomCode defs.RoomCode, cmd int, exculdUid []int, ph proto.Message) error {
	room := roomMapGet(roomCode)
	if room == nil {
		return errors.New(fmt.Sprintf("not fount %v room", roomCode))
	}

	room.players.Range(func(uid, p interface{}) bool {
		if !isExistArray(exculdUid, uid.(int)) {
			p.(*PlayerImpl).SendMsg(cmd, ph, 0)
		}
		return true
	})

	return nil
}


func isExistArray(src []int, des int) bool {
	for _, v := range src {
		if v == des {
			return true
		}
	}
	return false
}



func (r *Room) Close() {
	if !atomic.CompareAndSwapInt32(&r.runFlag, RoomStatusRunning, RoomStatusStopped) {
		log.Printf("room:%d is closed already :%d", r.RoomCode, r.runFlag)
		return
	}
	// 关闭相应定时器

	DeleteRoomEvents(r.RoomCode)
	//utils.RemoveCheats(r.InviteCode)

	close(r.sysMsg)
	close(r.workerChan)

	//go r.writeFile()
	for m := range r.workerChan {
		if m != nil {
			//todo:这个是应该要打开的，处理未完消息 2019.02.18
			//CreateRoomHelper.WriteMsg(m)
		} else {
			log.Println("msg is null", r.RoomCode)
		}
	}

	//读出所有的离线消息
	for range r.sysMsg {

	}

	log.Printf("inviteCode:%d 关闭?", r.RoomCode)
	roomMapDelete(r.RoomCode)
}
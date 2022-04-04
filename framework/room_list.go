package framework

import (
	"awesome/anet"
	"awesome/defs"
	"log"
	"sync"
)

var roomMap sync.Map



//func roomMap() map[interface{}]interface{} {
//	var m = map[interface{}]interface{}{}
//	roomMap.Range(func(k, v interface{}) bool {
//		m[k] = v
//		return true
//	})
//	return m
//}

func roomMapGet(roomCode defs.RoomCode) *Room {
	_room, _ := roomMap.Load(roomCode)
	if _room != nil {
		return _room.(*Room)
	}
	return nil
}

func RoomMapGet(code defs.RoomCode)*Room {
	return roomMapGet(code)
}

func roomMapSet(roomCode defs.RoomCode, room *Room) {
	roomMap.Store(roomCode, room)
}

func roomMapDelete(roomCode defs.RoomCode) (result int) {
	roomMap.Delete(roomCode)

	return 0
}

func roomMapCheck(roomCode defs.RoomCode) (result bool) {
	return check(roomCode)
}

func check(roomCode defs.RoomCode) bool {
	_, result := roomMap.Load(roomCode)
	return result
}



func BroadcastToAllRooms(cmd int, msg interface{}) {
	head := &anet.PackHead{Cmd: uint32(cmd)}
	err := SerializePackWithPB(head, msg)
	if err != nil {
		log.Println("broadcastAllRoom SerializePackWithPB error:", err)
		return
	}

	roomMap.Range(func(inviteCode, room interface{}) bool {
		if r, ok := room.(*Room); ok && r != nil {
			userMsg := &UserMessage{pack: head, user: nil}
			r.enqueueMessage(userMsg)
		} else {
			log.Printf("inviteCode:%v not Room val:%v", inviteCode, room)
		}
		return true
	})
}
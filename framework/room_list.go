package framework

import (
	"sync"
	"awesome/defs"
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

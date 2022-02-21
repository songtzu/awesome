package framework

import (
	"awesome/defs"

	"awesome/alog"
	"fmt"

)


func (p *PlayerImpl) LogInfo(args ...interface{}) {
	var roomCode defs.RoomCode =-1
	if p.room!=nil{
		roomCode= p.room.RoomCode
	}
	alog.Info(fmt.Sprintf("%d,%d,%d,", roomCode, p.userId, p.newestMessage()),fmt.Sprint( args...))
}
/*
 *
 */
func (p *PlayerImpl) LogErr(args ...interface{}) {
	//fmt.Println(args)
	var roomCode defs.RoomCode =-1
	if p.room!=nil{
		roomCode= p.room.RoomCode
	}
	alog.Err(fmt.Sprintf("%d,%d,%d,", roomCode, p.userId, p.newestMessage()),fmt.Sprint( args...))
}

func (p *PlayerImpl) LogWarn(args ...interface{}) {
	var roomCode defs.RoomCode =-1
	if p.room!=nil{
		roomCode= p.room.RoomCode
	}
	alog.Warn(fmt.Sprintf("%d,%d,%d,", roomCode, p.userId, p.newestMessage()),fmt.Sprint( args...))
}


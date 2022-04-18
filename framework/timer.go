package framework

import (
	"awesome/alog"
	"awesome/defs"
	"awesome/timer"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type TypeTimeTaskCallBack = func(key string,roomExtension interface{} )
var timerCbMap=make(map[string]TypeTimeTaskCallBack)

func timerCallback(innerKey string ){
	if roomCode,key:=parseTimerInnerKey(innerKey);roomCode!=0{
		room:=roomMapGet(roomCode)
		cb,ok:=timerCbMap[innerKey]
		msg:=&SystemMessage{SystemMessageDefTimer,key,nil}
		if ok{
			msg.DealHandle = cb
		}

		if room==nil{
			//nil room, redirect to special room chan.
			specialRoomInstance.enqueueSystemMessage(msg)
		}else {
			//room found for roomCode.
			room.enqueueSystemMessage(msg)
		}
	}else {
		alog.Err("bad error when try to parse timer innner key during callback")
	}
}

//// 设置时间
//func AddTimeTaskWithCallback(key string,interval int64,cb_ TypeTimeTaskCallBack) error{
//	timer.SetTimeTaskWithCallback(key,interval,cb)
//	return nil
//}

func parseTimerInnerKey(innerKey string) (roomCode defs.RoomCode, key string) {
	s := strings.Split(innerKey, ":")
	if len(s)>=2{
		if v,err:=strconv.Atoi(s[0]);err==nil{
			return defs.RoomCode(v),s[1]
		}
		return 0,""
	}
	return 0,""
}
//AddRoomTimeTaskWithCallback 设置时间
func AddRoomTimeTaskWithCallback( roomCode int,key string,interval int64,cb_ TypeTimeTaskCallBack) error{
	innerKey:=fmt.Sprintf("%d:%s", roomCode,key)
	if cb_!=nil{
		timerCbMap[innerKey] = cb_
	}
	timer.SetTimeTaskWithCallback(innerKey,interval,timerCallback)
	return nil
}


func TimerKeyGen(roomCode defs.RoomCode, event string) string  {
	return fmt.Sprintf("%d:%s", roomCode, event)
}

func TimerKeySplit(key string) ( defs.RoomCode,  string) {
	i := strings.Index(key,":")
	if i>0{
		invite, err := strconv.Atoi(key[0:i])
		if err==nil{
			if len(key) > i{
				return defs.RoomCode(invite),key[i+1:]
			}else {
				return defs.RoomCode(invite),""
			}
			//return  invite,key[i:]
		}else {
			alog.Info("拆分定时器key，无法解析出房间号")
			return 0,""
		}

	}
	alog.Err("拆分定时器key失败",key)
	return 0,""
}


func RemoveTimerEvent(roomCode defs.RoomCode, key string) {
	innerKey := TimerKeyGen(roomCode, key)
	timer.DeleteTimer(innerKey)
}

func DeleteRoomEvents(code defs.RoomCode)  {
	partialKey := TimerKeyGen(code,"")
	c:=timer.DeleteEventsByPartialMatch(partialKey)
	log.Printf("delete room:%d time event count:%d", code, c)
}
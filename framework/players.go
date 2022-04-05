package framework

import (
	"awesome/anet"
	"awesome/defs"
	"log"
	"sync"
)

type Players struct {
	sync.Map
}

//var _players	sync.Map		//this is deal with the real sub, which would sub some topics.

func (p *Players)playerGet(userId defs.TypeUserId) (value *PlayerImpl, ok bool) {
	if v,ok:= p.Load(userId);ok{
		if p,ok:=v.(*PlayerImpl);ok{
			return p,true
		}else {
			return nil,false
		}
	}
	return nil,false
}

func (p *Players)playerSet(userId defs.TypeUserId,v *PlayerImpl )  {
	p.Store(userId,v)
}


func (p *Players)playerDelete(userId defs.TypeUserId) {
	p.Delete(userId)
}

func (p *Players)playerExist(userId defs.TypeUserId) (  found bool) {
	if v,ok:= p.Load(userId);ok{
		if _,ok:=v.(*PlayerImpl);ok{
			return true
		}else {
			return false
		}
	}
	return false
}



var globalPlayerMap = Players{}


func UserMapGet(userid defs.TypeUserId) (user *PlayerImpl) {

	if u,ok := globalPlayerMap.Load(userid);ok{
		return u.(*PlayerImpl)
	}else {
		return nil
	}
}

func PushPackToUser(userid defs.TypeUserId, head *anet.PackHead) (isSucceed  bool) {
	if usrImp := UserMapGet(userid); usrImp != nil {
		if usrImp.conn != nil {
			if _,err:=usrImp.conn.WriteMessage(head);err==nil{
				return true
			}
		}
	}
	return false
}

func PushBinToUser(userid defs.TypeUserId, bin []byte) (isSucceed  bool) {
	if usrImp := UserMapGet(userid); usrImp != nil {
		if usrImp.conn != nil {
			if _,err:=usrImp.conn.WriteBinary(bin);err==nil{
				return true
			}
		}
	}
	return false
}


func UserMapStore(userid defs.TypeUserId, user *PlayerImpl) {
	globalPlayerMap.Store(userid, user)
}

func UserMapDelete(userid defs.TypeUserId) (result int) {
	result = 0
	if _,ok:=globalPlayerMap.Load(userid);ok{
		result = 1
	}
	globalPlayerMap.Delete(userid)
	return result
}

func UserMapCheck(userid defs.TypeUserId) (result bool) {
	if _,ok:=globalPlayerMap.Load(userid);ok{
		return true
	}
	return false
}

func UserMapUpdateUID(oldUid, newUid defs.TypeUserId) (result bool) {
	if oldUid == newUid {
		return true
	}
	user := UserMapGet(oldUid)
	if user != nil {
		globalPlayerMap.Delete(oldUid)

		user.userId = newUid
		globalPlayerMap.Store( newUid, user)
		return true
	} else {
		log.Println("- not find userid")
		return false
	}


}
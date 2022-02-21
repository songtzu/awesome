package framework

import (
	"sync"
	"awesome/defs"
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

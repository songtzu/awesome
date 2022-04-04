package framework

import "sync"

var generalMap sync.Map

func GlobalMapSet(key string,data interface{}) {
	generalMap.Store(key,data)
}

func GlobalMapGet(key string) interface{} {
	v,ok := generalMap.Load(key)
	if !ok {
		return nil
	}
	return v
}

func GlobalMapDelete(key string) {
	generalMap.Delete(key)
}

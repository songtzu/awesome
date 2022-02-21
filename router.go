package awesome

import (
	"awesome/anet"
	"reflect"
)

var routers = make(map[anet.PackHeadCmd]*router, 50)

type router struct {
	fun reflect.Value
	msg reflect.Type
}

// cmd 	消息id
// f   	处理消息的函数 如: login(session *server.ClientSession, req *protocol.CLogin) (resp proto.Message, err error)
// msg 	消息对应的protobuf请求包类型
func RegisterMessageRouter(cmd anet.PackHeadCmd, cb interface{}, msg interface{}) {
	_, ok := routers[cmd]
	if ok {
		//logdebug("cmd has registered before", cmd)
		return
	}

	if reflect.TypeOf(cb).Kind() != reflect.Func {
		//logerr("cb should be a function", cmd )
		return
	}

	if reflect.TypeOf(msg).Kind() != reflect.Struct {
		//logerr("msg should be a struct", cmd)
		return
	}
	routers[cmd] = &router{
		fun: reflect.ValueOf(cb),
		msg: reflect.TypeOf(msg),
	}
}

func GetFunc(cmd anet.PackHeadCmd) reflect.Value {
	return routers[cmd].fun
}

func Exist(cmd anet.PackHeadCmd) bool {
	_, ok := routers[cmd]
	return ok
}

func GetProto(cmd anet.PackHeadCmd) reflect.Type {
	return routers[cmd].msg
}

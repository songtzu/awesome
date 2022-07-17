package framework

import (
	"awesome/defs"
	"log"
	"reflect"
)

var cmdRouterMaps = make(map[uint32]*router, 50)

type router struct {
	fun reflect.Value
	msg reflect.Type
}

const (
	CmdRouteRespTypeProtobuf = 0
	CmdRouteRespTypeJsonObj = 1
	CmdRouteRespTypeBinary = 2
)

//RegisterCmdCallbackFunc
//cmd 	消息id
//f   	处理消息的函数 如: login(session *server.ClientSession, req *protocol.UserLogin) (resp proto.Message, CmdRouteRespType, cmd int)
// msg 	消息对应的protobuf请求包类型
func RegisterCmdCallbackFunc(cmd defs.TypeCmd, f interface{}, msg interface{}) {
	_, ok := cmdRouterMaps[cmd]
	if ok {
		log.Printf("cmd:%d重复注册", cmd)
		return
	}

	if reflect.TypeOf(f).Kind() != reflect.Func {
		log.Printf("cmd:%d处理不是函数,其类型型：%s ", cmd, reflect.TypeOf(f).Kind().String())
		return
	}

	if reflect.TypeOf(msg).Kind() != reflect.Struct {
		log.Printf("cmd:%d关联类型:%s不是一个结构体", cmd, reflect.TypeOf(msg).Kind().String() )
		return
	}
	cmdRouterMaps[cmd] = &router{
		fun: reflect.ValueOf(f),
		msg: reflect.TypeOf(msg),
	}
}

func getCmdRouterFunc(cmd defs.TypeCmd) reflect.Value {
	return cmdRouterMaps[cmd].fun
}

func isCmdRouterExist(cmd defs.TypeCmd) bool {
	_, ok := cmdRouterMaps[cmd]
	return ok
}

func getCmdRouterProto(cmd defs.TypeCmd) reflect.Type {
	return cmdRouterMaps[cmd].msg
}
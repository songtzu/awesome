package framework

import (
	"log"
	"reflect"
)

var cmdRouterMaps = make(map[uint32]*router, 50)

type router struct {
	fun reflect.Value
	msg reflect.Type
}


//RegisterCmdCallbackFunc
//cmd 	消息id
//f   	处理消息的函数 如: login(session *server.ClientSession, req *protocol.UserLogin) (resp proto.Message, err error)
// msg 	消息对应的protobuf请求包类型
func RegisterCmdCallbackFunc(cmd uint32, f interface{}, msg interface{}) {
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
	cmdRouterMaps[uint32(cmd)] = &router{
		fun: reflect.ValueOf(f),
		msg: reflect.TypeOf(msg),
	}
}

func GetFunc(cmd uint32) reflect.Value {
	return cmdRouterMaps[cmd].fun
}

func Exist(cmd uint32) bool {
	_, ok := cmdRouterMaps[cmd]
	return ok
}

func GetProto(cmd uint32) reflect.Type {
	return cmdRouterMaps[cmd].msg
}
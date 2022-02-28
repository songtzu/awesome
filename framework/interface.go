package framework

import (
	"awesome/anet"
	"awesome/defs"
	"fmt"
	"github.com/labstack/echo"
	"sync/atomic"
)
type Echo *echo.Echo

type IFramework interface {



	/*
	 * 注册http
	 */
	OnRegisterHttpRouters(e Echo)

	/*
	 *	初始化结束的回调
	 */
	OnInit()()
	/*
	 * 消息派发到logic实现接口
	 * room结构为logic创建的结构 通过OnCreateRoom时返回
	 */
	OnDispatchLogicMessage( roomCode defs.RoomCode, room *Room, user *PlayerImpl, msg *anet.PackHead) (err error)


	/*
	 * OnParseRoomCode接口解析的房间号为新房间的时候，调用此接口
	 *		此接口要求不处理业务逻辑，仅仅需要解析生成房间的缓存数据。
	 * 		逻辑层在此接口解析并创建房间缓存数据，
	 *		返回房间号和房间数据！
	 * 		返回错误时,消息将会被重定向至OnError
	 * 		返回成功，返回的房间数据会被缓存，并创建消息队列，此业务消息会被重定向到OnDispatchLogicMessage接口！
	 */
	OnCreateRoom(msg *anet.PackHead) ( extension interface{}, error error)

	/*
	 *  系统消息派发到 logic
	 * 		当前系统消息为超时消息和用户离线消息
	 */
	OnDispatchSystemMessage(room interface{}, msg *anet.PackHead) (err error)


	/*
	 * 错误重定向至此接口（）
	 */
	OnError(msg *anet.PackHead)


	/*
	 * 解析房间id
	 *  1、流程说明
	 * 			用户连接服务器，此时并未通过认证授权，服务器逻辑层和框架层无法得知用户所属房间。
	 *			也无法知道此用户是断线重连还是新开房间的新玩家。
	 * 			在用户提交OnParseRoomCode 接口所能解析的数据包解析出房间号之前，框架无法，也不应该派发具体的业务消息到正常的房间chan。
	 * 用于新用户加入查找房间时从协议中获取到房间RoomCode
	 */
	OnParseRoomCode(msg *anet.PackHead) (roomCode defs.RoomCode,err error)
}

var frameworkInterfaceInstance IFramework
var logicEngineInitFlag int32 = 0

func GetLogicEngine() IFramework {
	return frameworkInterfaceInstance
}

func InitFrameworkInstance(engine IFramework) error {

	if engine == nil {
		return fmt.Errorf(defs.GetError(defs.ErrorDefUnImplementInterface).Error())
	}

	if atomic.CompareAndSwapInt32(&logicEngineInitFlag, 0, 1) {
		frameworkInterfaceInstance = engine
		frameworkInterfaceInstance.OnInit()
	} else {
		return fmt.Errorf(defs.GetError(defs.ErrorDefMultiImplementInterface).Error())
	}
	return nil
}


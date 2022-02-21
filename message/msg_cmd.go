package message

const (
	//游戏服向大厅发起ping保活。
	InnerCmdPing    = 1
	InnerCmdPingAck = 2
	//游戏服向大厅注册服务
	InnerCmdRegistService    = 11
	InnerCmdRegistServiceAck = 12
	//大厅向游戏服创建房间
	InnerCmdRoomCreate    = 21
	InnerCmdRoomCreateAck = 22

	/*
	 * 大厅向游戏服发起，鉴于此接口走客户端通道（为了避免也业务层的CMD冲突，不在此定义），此CMD保留在配置文件中。
	 */
	//InnerCmdProxyActivity = 31
	//InnerCmdProxyActivityAck = 32

	//游戏服向大厅销毁房间
	InnerCmdApplyRoomDestroy    = 31
	InnerCmdApplyRoomDestroyAck = 32

	// 大厅向游戏服主动销毁
	InnerCmdRoomDestroy = 41
	InnerCmdRoomDestroyAck = 42

	// 大厅通知玩家进游戏
	InnerCmdPlayerEnter = 51
	InnerCmdPlayerEnterAck = 52

	// 游戏服通知大厅玩家已经离开
	InnerCmdApplyPlayerLeave = 61
	InnerCmdApplyPlayerLeaveAck = 62

	// 大厅向游戏服玩家离开游戏
	InnerCmdPlayerLeave = 71
	InnerCmdPlayerLeaveAck = 72
)

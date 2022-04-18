package framework

import "awesome/defs"

/*
 *1、用户连接，先回调，让逻辑层返回是否是匹配模式。
	如果是匹配模式，则不执行OnParseRoomCode回调，以及OnCreateRoom回调。
	如果不是匹配模式，
 * 解析场次信息的时候，返回
 * 需要能够在超时的时候返回超时批次的玩家。
 *	成功匹配到足够的人数，返回匹配成功的玩家列表。
 * 如果是正常的游戏场次，也需要从自动匹配的场次中获取玩家，如何实现？
 */

const ROOMCODEMATCH = -999

type matchPlayer struct {
	Round int			//场次
	Player *PlayerImpl	//玩家信息
	Deadline int64		//截止时间，精度为秒。
}

type roomForMatch struct {
	Players []*matchPlayer
}


func AddPlayerToMatchQueue()  {

}

//InvitePlayerFromMatchQueue 如需要给指定场次的游戏邀请count个玩家，调用此方法。
func InvitePlayerFromMatchQueue(code defs.RoomCode, count int)  {

}
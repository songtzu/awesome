package framework

import (
	"awesome/defs"
	"log"
)

/*
 *1、用户连接，先回调，让逻辑层返回是否是匹配模式。
	如果是匹配模式，则不执行OnParseRoomCode回调，以及OnCreateRoom回调。
	如果不是匹配模式，
 * 解析场次信息的时候，返回
 * 需要能够在超时的时候返回超时批次的玩家。
 *	成功匹配到足够的人数，返回匹配成功的玩家列表。
 * 如果是正常的游戏场次，也需要从自动匹配的场次中获取玩家，如何实现？
 */

//匹配事件的事件间隔
const match_event_interval = 1000

type matchPlayer struct {
	Player *PlayerImpl	//玩家信息
	MatchRuleData *MatchRule
}

type MatchRule struct {
	MatchNum int			//需要匹配的人数
	DeadlineTimestamp int64	//截止时间，精度为秒。
}
type roomForMatch struct {
	matchTaskMap map[int][]*matchPlayer	//每个匹配场次规则都单独定义一个map k-v，每次定时任务都检查
}

func newRoomContainerForMatch( firstRule *MatchRule , firstPlayerImp *PlayerImpl) (room *roomForMatch ){
	room = &roomForMatch{}
	room.matchTaskMap = make(map[int][]*matchPlayer)
	matchTask := &matchPlayer{MatchRuleData: firstRule, Player: firstPlayerImp}
	s:=make([]*matchPlayer,0)
	s = append(s, matchTask)
	room.matchTaskMap[firstRule.MatchNum] = s
	return room
}

func AddPlayerToMatchQueue()  {

}

//InvitePlayerFromMatchQueue 如需要给指定场次的游戏邀请count个玩家，调用此方法。
func InvitePlayerFromMatchQueue(code defs.RoomCode, count int)  {

}

func startMatchTimeTask()  {
	AddRoomTimeTaskWithCallback(RoomCodeMatch,"match_event", match_event_interval,matchCallback)
}

func matchCallback(key string, extension interface{})  {
	log.Printf("匹配的容器房间的定时任务:%s",key)
}
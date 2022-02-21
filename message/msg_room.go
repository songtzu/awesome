package message

/*****************************************
 * 大厅和游戏服务器通信，关于房间的协议。
 * 		此文件用于内网的tcp通信
 *
 ***************************************/

/*
 * 大厅向游戏服务器发起此请求，创建规则为extend类型的房间，
 * 	此协议内容复制自roomInfo，为了不给awesome框架引入pb的内容，兼容为json版本
 */
type RoomInfo struct {
	RoomCode uint32 `json:"roomCode,omitempty"`
	RoomId   int64  `json:"roomId,omitempty"`
	//游戏名。斗地主，跑得快等大类型。
	AppId     int    `json:"appId,omitempty"`
	MinPlayer uint32 `json:"minPlayer,omitempty"`
	MaxPlayer uint32 `json:"maxPlayer,omitempty"`
	//游戏类型，
	GameType        int     `json:"gameType,omitempty"`
	CreateTimeStamp uint64  `json:"createTimeStamp,omitempty"`
	GameStatus      int     `json:"gameStatus,omitempty"`
	GamePlayerNo    uint32  `json:"gamePlayerNo,omitempty"`
	CreatorId       int32   `json:"creatorId,omitempty"`
	ObserverLimit   int32   `json:"observerLimit,omitempty"`
	GameServerId    int64   `json:"gameServerId,omitempty"`
	UsersInRoom     []int32 `json:"usersInRoom,omitempty"`
	//游戏规则的细节。json扩展。
	ExtendGameRule string `json:"extendGameRule,omitempty"`
}
type RoomCreate struct {
	RoomInfo RoomInfo `json:"roomInfo"`
}
type RoomCreateAck struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

/*
 * 游戏服向大厅发送的销毁房间的协议
 */
type RoomDestroy struct {
}

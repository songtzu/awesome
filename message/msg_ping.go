package message

type MsgPing struct {
	IsOk bool `json:"isOk"`
	Timestamp int64 `json:"timestamp"`
	Address string `json:"address"`
	PlayerCount uint32 `json:"player_count"`
}


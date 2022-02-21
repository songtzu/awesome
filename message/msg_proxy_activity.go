package message

type MsgProxyActivity struct {
	Timestamp int64
	ProxyAddress string
}


type MsgProxyActivityAck struct {
	Status int
	ServerId int
	AppId int
	BindAddress string
	Version string
}

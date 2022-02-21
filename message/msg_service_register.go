package message

type DefServiceStatuType int

const (
	SeviceStatusNormal  DefServiceStatuType = 0
	SeviceStatusErrorService  DefServiceStatuType = 1
	SeviceStatusOutOfService  DefServiceStatuType = 2
)
type MsgRegistService struct {
	//0,正常，1异常服务，2，停服
	Status DefServiceStatuType
	ServerId int
	AppId int
	BindAddress string
	Version string
}

type MsgRegistServiceAck struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
} 
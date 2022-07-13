package framework

import "github.com/labstack/echo"

const (
	ErrorOk                        = 0
	ErrorUnDefined                 = 1
	ErrorDataBaseError             = 2
	ErrorPermissionNotAllowed      = 3
	ErrorExisted                   = 4
	ErrorParameterNeeded           = 5
	ErrorServiceNotFound           = 6
	ErrorUserNotFound              = 6
	ErrorDataNotFound              = 7
	ErrorFieldCannotBeEmpty        = 11
	ErrorDeviceHasRegisteredBefore = 12

	ErrorCanNotDeviceLogin    = 13
	ErrorParameterFormatError = 14 //数据格式错误
	ErrorServiceUnavailable   = 15 //服务不可用
	ErrorUserNotAuthored      = 20 //未登录用户。
	ErrorDataFalsify          = 21 // 篡改过的数据
	ErrorUserInMatch 			= 30	//用户正在匹配中
	ErrorPermissions           = 41 // 权限不够
)
type Echo = *echo.Echo

type EchoCtx = echo.Context

type GeneralResponse struct {
	Status  int         `json:"status" bson:"status"`
	Message string      `json:"message" bson:"message"`
	Data    interface{} `json:"data"`
}
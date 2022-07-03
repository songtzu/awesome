package awe_util

import (
	"github.com/labstack/echo"
)

func Register(e *echo.Echo, path string, get, post, patch, delete echo.HandlerFunc) {
	e.GET(path, get)
	e.POST(path, post)
	e.PATCH(path, patch)
	e.DELETE(path, delete)
}

type HandlerWithSession func(ctx echo.Context, data string) error

//var headerInstance interface{} = nil
//
//func RegisHeaderType( header interface{} )  {
//	headerInstance = header
//}

type GeneralResponse struct {
	Status  int         `json:"status" bson:"status"`
	Message string      `json:"message" bson:"message"`
	Data    interface{} `json:"data"`
}
//
//func EchoHandlerWithSession(handlerWithSession HandlerWithSession) func(ctx echo.Context) error {
//	return func(ctx echo.Context) error {
//
//		if str, err := HeaderAuthorInfo(ctx); err != nil {
//			return ctx.JSON(http.StatusOK, &GeneralResponse{Status: -1, Message: err.Error()})
//		} else {
//			return handlerWithSession(ctx, str)
//		}
//	}
//}

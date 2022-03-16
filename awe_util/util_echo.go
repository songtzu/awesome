package awe_util

import (
	"github.com/labstack/echo"
	"net/http"
)

func Regist(e *echo.Echo, path string, get, post, patch, delete echo.HandlerFunc) {
	e.GET(path, get)
	e.POST(path, post)
	e.PATCH(path, patch)
	e.DELETE(path, delete)
}

type HandlerWithSession func(ctx echo.Context, user *model.User) error

func EchoHandlerWithSession(handlerWithSession HandlerWithSession) func(ctx echo.Context) error {

	return func(ctx echo.Context) error {
		//if(ctx.Request().URL.Path=="/api/user" || strings.Contains(ctx.Request().URL.Path,"test")) {
		//
		//}

		if userInfo, err := HeaderAuthorInfo(ctx); err != nil {
			return ctx.JSON(http.StatusOK, &defs.GeneralResponse{Status: defs.ErrorPermissionNotAllowed, Message: err.Error()})
		} else {
			return handlerWithSession(ctx, userInfo)
		}

	}
	//return func(c *echo.Context) error { return handler(c, id,role) }
}

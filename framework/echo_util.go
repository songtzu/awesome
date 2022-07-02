package framework

import (
	"awesome/db"
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"sync"
)


var sessionMap sync.Map //session的map,string(userid)--->*model.UserSession

type HandlerWithSession func(ctx echo.Context, userSession []byte) error

func EchoHandlerWithSession(handlerWithSession HandlerWithSession) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		if userInfo, err := HeaderAuthorInfo(ctx); err != nil {
			return ctx.JSON(http.StatusOK, GeneralResponse{Status: ErrorPermissionNotAllowed, Message: err.Error()})
		} else {
			return handlerWithSession(ctx, userInfo)
		}
	}
}

//HeaderAuthorInfo 检查会话信息。
func HeaderAuthorInfo(c echo.Context) (user []byte, err error) {
	var token = c.Request().Header.Get("token")
	if len(token) == 0 {
		return nil, errors.New(fmt.Sprintf("未收取的请求token:%s ", token))
	}
	if v, ok := sessionMap.Load(token); ok {
		if sd, ok := v.(string); !ok {
			return nil, errors.New("session错误")
		}  else {
			//log.Println("内存缓存命中session", id)
			return []byte(sd), nil
		}
	}
	//session找不到，去redis中查询。
	//var sd = &model.UserSession{}
	//if err = db.RedisHMGet(defs.RedisHMSession, token, sd); err != nil {
	//	return nil, err
	//}
	//log.Println("进程重启过，从redis命中token", sd.UserData.Id)
	err = db.GetRedisClient().
	if err!=nil{
		log.Println("toke过期或者被伪造", sd.Token, token)
		return nil, errors.New("toke过期或者被伪造")
	}
	StoreSessionToSession(token, sd)
	return &sd.UserData, nil
}


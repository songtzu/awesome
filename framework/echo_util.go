package framework

import (
	"awesome/awe_util"
	"awesome/db"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)


var sessionMap sync.Map //session的map,string(userid)--->*model.UserSession

const redisKeyFormatToken = "http_token:%s"

const redisKeyFormatRequestCache = "cache:%s_%s" //cache:url:token

const SessionTokenSalt = "AiWj8720DWdW9AcJo"
//请求的缓存时长，默认30秒。
var requestCacheSeconds = 30

var cacheUrlMap sync.Map	//需要缓存的url类型

func SetRequestCacheSeconds(s int)  {
	requestCacheSeconds = s
}

type HandlerWithSession func(ctx EchoCtx, userSession string) error
//type HandlerEcho func(ctx EchoCtx ) error
type HandlerEcho  = echo.HandlerFunc

func EchoHandlerWithSession(handlerWithSession HandlerWithSession) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		if userInfo, err := HeaderAuthorInfo(ctx); err != nil {
			return ctx.JSON(http.StatusOK, GeneralResponse{Status: ErrorPermissionNotAllowed, Message: err.Error()})
		}else {
			if cache, ok := isRequestCached(ctx); ok{
				log.Println("命中缓存", cache)
				ctx.Response().Header().Set("cached", "true")
				return ctx.JSONBlob(http.StatusOK, []byte(cache))

			}else{
				log.Println("未命中缓存",ctx.Request().URL)
				return handlerWithSession(ctx, userInfo)
			}
		}
	}
}

func echoHandlerWrap(he HandlerEcho) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		if cache, ok := isRequestCached(ctx); ok {
			log.Println("命中缓存",cache)
			ctx.Response().Header().Set("cached","true")
			return ctx.JSON(http.StatusOK, cache)
		} else {
			return he(ctx)
		}
	}
}

func isRequestCached(c echo.Context) (result string, ok bool) {

	if _,ok= cacheUrlMap.Load(c.Request().URL.String());!ok{
		log.Printf("url不允许缓存:%s",c.Request().URL )
		return "",false
	}
	var token = c.Request().Header.Get("token")
	var err error = nil
	if result, err = db.RedisKeyGetStr(fmt.Sprintf(redisKeyFormatRequestCache,c.Request().URL,token));err==nil{
		ok = true
	}else {
		ok = false
	}

	log.Printf("路径:%s,是否:%v命中缓存,结果:%s",c.Request().URL, ok, result)
	return result,ok
}

//HeaderAuthorInfo 检查会话信息。
func HeaderAuthorInfo(c echo.Context) (user string, err error) {
	var token = c.Request().Header.Get("token")
	cookieToken,err := c.Cookie("token")
	if len(token) == 0 && err==nil{
		token = cookieToken.Value
	}
	if len(token) == 0 {
		return "", errors.New(fmt.Sprintf("未收取的请求token:%s ", token))
	}
	if v, ok := sessionMap.Load(token); ok {
		if sd, ok := v.(string); !ok {
			return "", errors.New("session错误")
		}  else {
			//log.Println("内存缓存命中session", id)
			return sd, nil
		}
	}

	v,err := db.RedisKeyGetStr(fmt.Sprintf(redisKeyFormatToken,token))
	if err!=nil{
		log.Printf("toke:%s过期或者被伪造err:%s", token,err.Error())
		return "", errors.New("toke过期或者被伪造")
	}
	sessionMap.Store(token,v)
	return v, nil
}

func SessionSet(c echo.Context, user string, ttl time.Duration) (token string, err error) {
	s, err := session.Get("session", c)
	if err != nil {
		return "", err
	}
	s.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	token =  awe_util.PasswordMd5(user+strconv.Itoa(int(time.Now().Unix()))+SessionTokenSalt)
	s.Values["token"] = token
	err = s.Save(c.Request(), c.Response())
	c.Response().Header().Set("token", token)

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: false,
		MaxAge:   3600,
	}
	//c.Response().Header().Set("Set-Cookie", cookie.String())
	c.SetCookie(cookie)
	sessionMap.Store(token, user)
	//保存到redis中
	err = db.RedisKeySetStr(fmt.Sprintf(redisKeyFormatToken,token), user, ttl)
	if err!=nil{
		log.Printf("session保存失败:%s,生成token:%s, user:%s, ttl:%d", err.Error(),  token, user,ttl)
	}else{
		log.Printf("session保存成功,生成token:%s, user:%s, ttl:%d",  token, user,ttl)
	}

	return token, err
}
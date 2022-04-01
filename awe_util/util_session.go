package awe_util

import (
	"awesome/db"
	"errors"
	"github.com/labstack/echo"
	"log"
	"sync"

	"github.com/labstack/echo-contrib/session"
	//"net/http"
	//"defs"
	"github.com/gorilla/sessions"
)

const RedisSessionHashMap  = "userSession"
var sessionMap sync.Map //session的map,string(userid)--->*model.UserSession




/*SessionSetter
 *把session保存两份，一份保存到内存，一份保存到redis，客户端的cookie只保留session对应的key
 ********/
func SessionSetter(c echo.Context, token string, data interface{}) ( err error) {
	s, err := session.Get("session", c)
	if err != nil {
		return err
	}
	s.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	s.Values["token"] = token
	s.Save(c.Request(), c.Response())
	sessionMap.Store(token, data)
	//save session key to redis.
	db.RedisHMSet(RedisSessionHashMap, token, data)
	log.Printf("store redis session to redis key :%s, value:%v", token, data)
	return  nil
}

/*SessionFill
 * 把session保存两份，一份保存到内存，一份保存到redis，客户端的cookie只保留session对应的key
 ********/
func SessionFill(c echo.Context, token string, id int) error {
	s, err := session.Get("session", c)
	if err != nil {
		log.Println("找不到ssession")
		return err
	}
	s.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	s.Values["token"] = token
	s.Values["id"] = id
	s.Save(c.Request(), c.Response())

	return nil
}

//HeaderAuthorInfo 检查会话信息。
func HeaderAuthorInfo(c echo.Context) (str string, err error) {
	//s, err := session.Get("session", c)

	token := c.Request().Header.Get("token")
	if len(token) == 0   {
		return "",errors.New( "request illegal " )
	}
	ok:=false
	if _,ok = sessionMap.Load(token); ok {
		return "",nil
	}
	//session找不到，去redis中查询。
	if str,err = db.RedisHMGetStr(RedisSessionHashMap, token); err != nil {
		return str,err
	}

	return str,nil
}


//SessionParser 用http头的token从内存中获取会话信息。
func SessionParser(c echo.Context, data interface{}) (  err error) {
	//s, err := session.Get("session", c)

	token := c.Request().Header.Get("token")
	if len(token) == 0   {
		return errors.New( "request illegal " )
	}
	ok:=false
	if _,ok = sessionMap.Load(token); ok {
		return nil
	}
	//session找不到，去redis中查询。
	if err = db.RedisHMGet(RedisSessionHashMap, token, data); err != nil {
		return err
	}

	return nil
}



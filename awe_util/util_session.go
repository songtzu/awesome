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


func getSessionFromMemAndRedis(name string, c echo.Context) {

}

//
//func SetSessionKey(c echo.Context, id int64, email string, usrName string, extend string, role string) (err error) {
//	s, err := session.Get("session", c)
//	if err != nil {
//		return err
//	}
//	s.Options = &sessions.Options{
//		Path:     "/",
//		MaxAge:   86400 * 7,
//		HttpOnly: true,
//	}
//	log.Println("SetSessionKey", id, email)
//
//	s.Values["id"] = strconv.Itoa(int(id))
//	s.Values["email"] = email
//	s.Values["username"] = usrName
//	s.Values["timestamp"] = strconv.Itoa(int(time.Now().Unix()))
//	s.Values["role"] = role
//
//	s.Values["token"] = PasswordMd5(strconv.Itoa(int(id)) + email + usrName + strconv.Itoa(int(time.Now().Unix())) + role)
//	s.Values["extend"] = extend
//	log.Println("------", s.Values["id"])
//	s.Save(c.Request(), c.Response())
//	return nil
//}

/*SessionSet
 * 把session保存两份，一份保存到内存，一份保存到redis，客户端的cookie只保留session对应的key
 ********/
func SessionSet(c echo.Context, token string, data interface{}) ( err error) {
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
func HeaderAuthorInfo(c echo.Context, data interface{}) (err error) {
	//s, err := session.Get("session", c)

	token := c.Request().Header.Get("token")
	if len(token) == 0   {
		return errors.New( "request illgal " )
	}
	ok:=false
	if data,ok = sessionMap.Load(token); ok {
		return nil
	}
	//session找不到，去redis中查询。
	if err = db.RedisHMGet(RedisSessionHashMap, token, data); err != nil {
		return err
	}

	return nil
}

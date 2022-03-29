package awe_util

import (
	"awesome/db"
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/labstack/echo-contrib/session"
	//"net/http"
	//"defs"
	"github.com/gorilla/sessions"
)

const RedisSessionHashMap  = "userSession"
var sessionMap sync.Map //session的map,string(userid)--->*model.UserSession

// Logger returns a middleware that logs HTTP requests.
func SessionChecker() echo.MiddlewareFunc {
	log.Println("Session===>MiddlewareFunc========>")
	return sessionBody()
}

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

/********
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

/********
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

//检查会话信息。
func HeaderAuthorInfo(c echo.Context, data interface{}) (err error) {
	//s, err := session.Get("session", c)

	token := c.Request().Header.Get("token")
	if len(token) == 0   {
		return   errors.New(fmt.Sprintf("未收取的请求token:%s,id:%v", token, id))
	}

	if v, ok := sessionMap.Load(token); ok {
		if sd, ok := v.(*model.UserSession); !ok {
			return nil, errors.New("session错误")
		} else if sd.Token != token {
			log.Println("toke过期或者被伪造", sd.Token, token)
			return nil, errors.New("toke过期或者被伪造")
		} else {
			//log.Println("内存缓存命中session", id)
			return &(sd.UserData), nil
		}
	}
	//session找不到，去redis中查询。
	var sd = &model.UserSession{}
	if err = db.RedisHMGet(defs.RedisHMSession, id, sd); err != nil {
		return nil, err
	}
	if sd.Token != token {
		log.Println("toke过期或者被伪造", sd.Token, token)
		return nil, errors.New("toke过期或者被伪造")
	}
	fmt.Println(sessionMap.Load(id))
	fmt.Println("进程重启过，从redis命中token", sd.UserData.Id)
	//sessionMap.Store(strconv.Itoa(id), sd) //
	StoreSessionToSession(id, sd)
	return &sd.UserData, nil
}

//检查会话信息。
func SessionGet(c echo.Context) (user *model.User, err error) {
	s, err := session.Get("session", c)
	if err != nil {
		return nil, err
	}
	var token = s.Values["token"]
	if token == nil || s.Values["id"] == nil {
		return nil, errors.New(fmt.Sprintf("未收取的请求token:%s,id:%v", token, s.Values["id"]))
	}
	log.Println(reflect.TypeOf(s.Values["id"]), s.Values["id"])
	var id int = 0
	if reflect.TypeOf(s.Values["id"]).Kind() == reflect.String {
		var sid = s.Values["id"].(string)
		if id, err = strconv.Atoi(sid); err != nil {
			return nil, err
		}
	} else {
		id = int(s.Values["id"].(int64))
	}

	if v, ok := sessionMap.Load(strconv.Itoa(id)); ok {
		if sd, ok := v.(*model.UserSession); !ok {
			return nil, errors.New("session错误")
		} else if sd.Token != token {
			log.Println("toke过期或者被伪造", sd.Token, token)
			return nil, errors.New("toke过期或者被伪造")
		} else {
			//log.Println("内存缓存命中session", id)
			return &(sd.UserData), nil
		}
	}
	//session找不到，去redis中查询。
	var sd = &model.UserSession{}
	if err = db.RedisHMGet(defs.RedisHMSession, strconv.Itoa(id), sd); err != nil {
		return nil, err
	}
	if sd.Token != token {
		log.Println("toke过期或者被伪造", sd.Token, token)
		return nil, errors.New("toke过期或者被伪造")
	}
	fmt.Println(sessionMap.Load(id))
	fmt.Println("进程重启过，从redis命中token", sd.UserData.Id)
	//sessionMap.Store(strconv.Itoa(id), sd) //
	StoreSessionToSession(strconv.Itoa(id), sd)
	return &sd.UserData, nil
}
func StoreSessionToSession(id string, sd *model.UserSession) {
	sessionMap.Store(id, sd) //
}
func GetUserInfoBySession(ctx echo.Context) (user *model.User) {
	// 从内存中获取玩家

	return nil
}

func ParseSessionByUid(ctx echo.Context) (uid int64, email, usrName, extend string, _role int64) {
	s, err := session.Get("session", ctx)
	if err != nil {
		log.Println("====", err)
		return 0, "", "", "", 0
	}
	if s == nil {
		log.Println("会话内容为空")
		return 0, "", "", "", 0
	}
	if s.Values["id"] == nil {
		return 0, "", "", "", 0
	}
	//log.Println("ParseSessionByUid=============>", s.Values["id"])
	str_uid := s.Values["id"].(string)
	//bson.UnmarshalJSON([]byte(str_uid),&uid)
	id, _ := strconv.Atoi(str_uid)
	role, _ := strconv.Atoi(s.Values["role"].(string))
	//log.Println("the ======== uid is =========", id, str_uid)
	//obj:=s.Values["id"].(bson.ObjectId)
	//log.Println("ParseSessionByUid====>", uid)
	return int64(id), s.Values["email"].(string), s.Values["username"].(string), s.Values["extend"].(string), int64(role)
}

/***
 * TODO,先不做session合法性校验，以后扩展使用redis保存的session
 ***/
func sessionBody() echo.MiddlewareFunc {
	// Defaults
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			//url := strings.Replace(c.Request().URL.Path, "/", "", -1)
			//url = strings.ToLower(url)
			//log.Println(url)
			if ((c.Request().URL.Path == "user" || c.Request().URL.Path == "/api/user") && (c.Request().Method == "GET" || c.Request().Method == "POST")) || strings.Contains(c.Request().URL.Path, "test") || c.Request().URL.Path == "fileupload" {
				//USER请求的GET、POST方法不校验session。
				log.Println("方法为GET/POST一种。")
			} else {
				//glog.Infoln("需要检查session", url,"method", c.Request().Method)
				s, err := session.Get("session", c)
				log.Println(s, err)
				if err != nil {
					log.Println("====", err)
					return c.JSON(http.StatusOK, &defs.GeneralResponse{Status: defs.ErrorPermissionNotAllowed, Message: "未授权访问"})
				}
				if s == nil {
					log.Println("会话内容为空")
					return next(c)
				}

				if s.Values["id"] == nil || s.Values["token"] == nil || s.Values["username"] == nil || s.Values["email"] == nil || s.Values["timestamp"] == nil || s.Values["role"] == nil {
					log.Println("访问失败", s.Values["id"])
					return c.JSON(http.StatusOK, &defs.GeneralResponse{Status: defs.ErrorPermissionNotAllowed, Message: "未授权访问"})
				}
				checkSum := PasswordMd5(s.Values["id"].(string) + s.Values["email"].(string) + s.Values["username"].(string) + s.Values["timestamp"].(string) + s.Values["role"].(string))
				//log.Println("======checkSum======", checkSum)
				if s.Values["token"] != checkSum {
					log.Println("篡改授权", s)
					return c.JSON(http.StatusOK, &defs.GeneralResponse{Status: defs.ErrorPermissionNotAllowed, Message: "安全警告，使用篡改的授权"})
				}
			}
			//glog.Infoln("执行后续操作")
			return next(c)

		}
	}
}

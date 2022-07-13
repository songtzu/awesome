package framework

import (
	"awesome/db"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"log"
	"time"
)
func allowOrigin(origin string) (bool, error) {
	// In this example we use a regular expression but we can imagine various
	// kind of custom logic. For example, an external datasource could be used
	// to maintain the list of allowed origins.
	return true, nil
}

var echoInstance *echo.Echo


func StartEchoServer(address string) (err error) {

	echoInstance = echo.New()
	echoInstance.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	echoInstance.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOriginFunc: allowOrigin,
			//AllowOrigins: []string{"*"},
			/*********
			 * https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Access-Control-Allow-Headers
			 */
			AllowHeaders: []string{"Content-Type", "Accept-Language", "Access-Token", "Authorization", "id", "token", "Set-Cookie"},
			AllowCredentials: true,
			ExposeHeaders: []string{"Content-Type", "Access-Token", "Authorization", "Access-Control-Max-Age", "id", "token", "Set-Cookie"},
			MaxAge:        8640000,
		}))
	// Middleware
	//e.Pre(middleware.RemoveTrailingSlash())
	echoInstance.Use(middleware.Logger())
	echoInstance.Use(middleware.Recover())
	echoInstance.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		if v,ok:=cacheUrlMap.Load(c.Request().URL.String());ok{
			token := c.Request().Header.Get("token")
			//ctx.Response().Header().Set("cached","true")
			if c.Response().Header().Get("cached") != "true"{
				db.RedisKeySetStr(fmt.Sprintf(redisKeyFormatRequestCache,c.Request().URL.String(),token),string(resBody),time.Duration(v.(int))*time.Second)
			}else{
				log.Println("本次命中缓存，不再更新缓存",c.Request().URL)
			}
		}
		log.Println(c.Request().Method, c.Request().URL,"\n", string(reqBody),"\n", string(resBody))

	}))
	go func() {
		err = echoInstance.Start(address)
	}()
	time.Sleep(10*time.Millisecond)
	return err
}

func RegisterBatchHttpHandle( path string, get, post, patch, delete HandlerWithSession) (err error) {
	log.Println("注册路由")
	if echoInstance==nil{
		log.Println("http server not start")
		return errors.New("http server not start")
	}
	echoInstance.GET(path, EchoHandlerWithSession(get))
	echoInstance.POST(path, EchoHandlerWithSession(post))
	echoInstance.PATCH(path,EchoHandlerWithSession(patch))
	echoInstance.DELETE(path,EchoHandlerWithSession(delete))
	log.Printf("注册:%s,Get:%v, Post:%v, patch:%v, delete:%v",path, get, post, patch, delete)
	return nil
}

func RegisterHttpGetWithSessionHandle( path string, handle HandlerWithSession, cacheSeconds int ) (err error) {
	if echoInstance==nil{
		log.Println("http server not start")
		return errors.New("http server not start")
	}
	if cacheSeconds >0 {
		cacheUrlMap.Store(path,cacheSeconds)
		log.Println(path,len(path))
		log.Println(cacheUrlMap.Load(path))
	}
	echoInstance.GET(path,EchoHandlerWithSession(handle))
	return nil
}

func RegisterHttpPostWithSessionHandle( path string, handle HandlerWithSession) (err error) {
	if echoInstance==nil{
		log.Println("http server not start")
		return errors.New("http server not start")
	}
	echoInstance.POST(path,EchoHandlerWithSession(handle))
	return nil
}

func RegisterHttpDeleteWithSessionHandle( path string, handle HandlerWithSession) (err error) {
	if echoInstance==nil{
		log.Println("http server not start")
		return errors.New("http server not start")
	}
	echoInstance.DELETE(path,EchoHandlerWithSession(handle))
	return nil
}

func RegisterHttpPatchWithSessionHandle( path string, handle  HandlerWithSession) (err error) {
	if echoInstance==nil{
		log.Println("http server not start")
		return errors.New("http server not start")
	}
	echoInstance.PATCH(path,EchoHandlerWithSession(handle))
	return nil
}





func RegisterHttpGetHandle( path string, handle HandlerEcho, cacheSeconds int ) (err error) {
	if echoInstance==nil{
		log.Println("http server not start")
		return errors.New("http server not start")
	}
	if cacheSeconds >0 {
		cacheUrlMap.Store(path,cacheSeconds)
	}

	echoInstance.GET(path,echoHandlerWrap(handle))
	return nil
}

func RegisterHttpPostHandle( path string, handle HandlerEcho) (err error) {
	if echoInstance==nil{
		log.Println("http server not start")
		return errors.New("http server not start")
	}
	echoInstance.POST(path,handle)
	return nil
}

func RegisterHttpDeleteHandle( path string, handle HandlerEcho) (err error) {
	if echoInstance==nil{
		log.Println("http server not start")
		return errors.New("http server not start")
	}
	echoInstance.DELETE(path, handle)
	return nil
}

func RegisterHttpPatchHandle( path string, handle HandlerEcho) (err error) {
	if echoInstance==nil{
		log.Println("http server not start")
		return errors.New("http server not start")
	}
	echoInstance.PATCH(path,handle)
	return nil
}


package framework

import (
	"errors"
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
	//e.Use(standard.WrapMiddleware(cors.New(cors.Options{
	//	AllowedOrigins: []string{"http://localhost"},
	//}).Handler))
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

func RegisterHttpGetWithSessionHandle( path string, handle HandlerWithSession) (err error) {
	if echoInstance==nil{
		log.Println("http server not start")
		return errors.New("http server not start")
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





func RegisterHttpGetHandle( path string, handle HandlerEcho) (err error) {
	if echoInstance==nil{
		log.Println("http server not start")
		return errors.New("http server not start")
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


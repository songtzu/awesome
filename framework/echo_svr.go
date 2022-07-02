package framework

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"log"
)
func allowOrigin(origin string) (bool, error) {
	// In this example we use a regular expression but we can imagine various
	// kind of custom logic. For example, an external datasource could be used
	// to maintain the list of allowed origins.
	return true, nil
}

var echoInstance *echo.Echo


func StartEchoServer(address string)  {

	echoInstance := echo.New()
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
		log.Println(c.Request().Method, c.Request().URL, string(reqBody), string(resBody))
	}))
	echoInstance.Logger.Fatal(echoInstance.Start(address))
}

func RegisterHttpHandle( path string, get, post, patch, delete echo.HandlerFunc) (err error) {
	if echoInstance==nil{
		return fmt.Errorf("http server not start,")
	}
	echoInstance.GET(path, get)
	echoInstance.POST(path, post)
	echoInstance.PATCH(path,patch)
	echoInstance.DELETE(path,delete)
	return nil
}
//
//func ()  {
//
//}
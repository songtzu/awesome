package ws

import (
	"testing"
	"net/http"
	"golang.org/x/net/websocket"
	"fmt"

	"time"

	"awesome/anet"
	"strconv"
)


func svrConnHandler(conn *websocket.Conn,tag string, delay int) {
	request := make([]byte, 128)
	defer conn.Close()
	for {
		readLen, err := conn.Read(request)

		if err!=nil{
			fmt.Println("read error",err,tag)
		}
		fmt.Println("tag",tag)
		//socket被关闭了
		if readLen == 0 {
			fmt.Println("Client connection close!")
			break
		} else {
			//输出接收到的信息
			fmt.Println("读取数据",tag)
			fmt.Println(string(request[:readLen]))
			go writeConn(conn,delay,tag)
			//time.Sleep(time.Second)
			////发送
			//time.Sleep(time.Duration(delay)*time.Second)
			//conn.Write([]byte(tag+" World !"))
		}

		request = make([]byte, 128)
	}
}

func writeConn(conn *websocket.Conn,delay int, tag string)  {
	//发送
	time.Sleep(time.Second)
	time.Sleep(time.Duration(delay)*time.Second)
	conn.Write([]byte(tag+" World !"))
}

func pingHandler(conn *websocket.Conn) {
	//time.Sleep(1*time.Second)
	svrConnHandler(conn,"pingHandler",0)
}

func echoHandler(conn *websocket.Conn)  {
	//time.Sleep(5*time.Second)
	fmt.Println("延迟5秒的接口")
	go svrConnHandler(conn,"echoHandler",5)
}

func rootHandler(conn *websocket.Conn)  {
	//time.Sleep(10*time.Second)
	fmt.Println("延迟10秒的接口,root")
	 svrConnHandler(conn,"rootHandler",10)
}
func startSvr()  {
	fmt.Println("Func inside.")

	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/ping",websocket.Handler(pingHandler))
	http.Handle("/", websocket.Handler(rootHandler))

	err := http.ListenAndServe("127.0.0.1:19999", nil)
	if err!=nil{
		fmt.Errorf("websock启动失败%v", err)
	}
	fmt.Println("Func finish.")
}
func TestWebSocketSvr(t *testing.T) {
	startSvr()
	time.Sleep(5*time.Minute)
}

func TestWebSockClientPing(t *testing.T)  {
	clientPing("Jack")
}
func TestWebSockClientPing2(t *testing.T)  {
	clientPing("Rose")
}


func TestWebSockClientEcho(t *testing.T)  {
	clientEcho()
}

func TestWebSockClientRoot(t *testing.T)  {
	clientRoot()
}

func TestWebSockClientMultiThread(t *testing.T)  {
	go clientPing("multi")
	go 	clientEcho()
	go clientRoot()
	time.Sleep(1*time.Minute)
}

func TestWebSockClientMultiThread2(t *testing.T)  {

	go clientRoot()
	go clientRoot()
	for index:=0;index<=10 ;index++  {
		fmt.Println("休眠",index)
		time.Sleep(1*time.Second)
	}
	time.Sleep(1*time.Minute)
}


func clientPing(tag string)  {
	origin := "http://127.0.0.1/"
	url := "ws://127.0.0.1:19999/ping"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		//alog.Fatal(err)
		fmt.Println(err)
	}

	for index:=0;index<100;index++{
		pack:=&anet.PackHead{Body: []byte("=================="+tag+"===================>hello, world!"+strconv.Itoa(index)),Cmd:uint32(1)}
		bin,_:=pack.SerializePackHead()
		if _, err := ws.Write(bin); err != nil {
			fmt.Println(err)
		}
		fmt.Println("ping发送",string(pack.Body))
		var msg = make([]byte, 512)
		var n int
		if n, err = ws.Read(msg); err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Received: %s.\n", msg[:n])
		time.Sleep(1000*time.Millisecond)
	}

}

func clientEcho()  {
	origin := "http://127.0.0.1/"
	url := "ws://127.0.0.1:19999/echo"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		fmt.Println(err)
	}
	if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
		fmt.Println(err)
	}
	fmt.Println("echog发送")
	var msg = make([]byte, 512)
	var n int
	if n, err = ws.Read(msg); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Received: %s.\n", msg[:n])
}

func clientRoot()  {
	origin := "http://127.0.0.1/"
	url := "ws://127.0.0.1:19999/"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("root发送")
	if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
		fmt.Println(err)
	}
	var msg = make([]byte, 512)
	var n int
	if n, err = ws.Read(msg); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Received: %s.\n", msg[:n])
}
package anet

import (
	"log"
	"net/http"
	"runtime/debug"

	"golang.org/x/net/websocket"
)

var netCallbackInstance InterfaceNet

func webSockHandler(conn *websocket.Conn) {
	conn.PayloadType = websocket.BinaryFrame
	//logdebug("ws根消息的回调处理，webSockHandler")
	forkInstance := forkNewInstanceOfInterface(netCallbackInstance)
	con := NewWebSockConnection(conn, forkInstance)
	forkInstance.IOnNewConnection(con)
	//logdebug("判断是否是tcp协议",con.netProtocol== NetProtocolTypeTCP)
	webSockReadLoop(con)
}

func StartWebSocket(addr string, instance InterfaceNet) {
	netCallbackInstance = instance
	http.Handle("/", websocket.Handler(webSockHandler))
	log.Println("ws StartWebSocket启动.")

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Println("websock start up failed, err:%v", err)
	}
}

func webSockReadLoop(c *Connection) {
	var errInfo error
	conn := c.connWebSock
	request := make([]byte, 1024)
	//defer conn.Close()
	defer func() {
		log.Println("connection tcpReadLoop exit")
		if err := recover(); err != nil {
			log.Println(err, string(debug.Stack()))
		}
		go c.CloseConnWithoutRecon(errInfo)
	}()
	var remains = 0

	for {
		readLen, err := conn.Read(request)
		remains += readLen
		if err != nil {
			//logerr("read error",err)
		}
		//socket被关闭了
		if readLen == 0 {
			log.Println("websocket client connection close!", c.connWebSock.LocalAddr().String(), c.connWebSock.RemoteAddr().String())
			break
		} else {
			if msg, parsedSize := parsePackHead(request, remains); msg != nil {
				//logdebug("收到一个消息包")
				c.dispatchMsg(msg)
				//bin,_:=msg.SerializePackHead()
				//conn.Write(bin)
				remains -= parsedSize
				copy(request, request[parsedSize:parsedSize+remains])
			} else {
				//包比较大，开辟更大的缓存读数据
				if newBufferSizeWhenBufferNotEnough*len(request) <= maxPackSize {
					log.Println("windowBuffer空间不足，", len(request))
					var largeBuffer = make([]byte, newBufferSizeWhenBufferNotEnough*len(request))
					copy(largeBuffer, request[0:remains])
					request = largeBuffer
				} else {
					//一次数据超限，建议关闭连接
					remains = 0
				}
			}
			////输出接收到的信息
			//fmt.Println(string(request[:readLen]))
			//
			////发送
			//conn.Write([]byte("World !"))
		}
		//request = make([]byte, 128)
	}
}

//
//func (u *userConn) processPack(pack *PackHead) {
//
//	if u.room!=nil{
//		//写入chan，
//		u.room.enqueueMessage(pack)
//	}else {
//		//todo, 房间不存在，写入单独的chan
//		specialRoomInstance.enqueueMessage(pack)
//	}
//
//}
//
//func (u *userConn) writeMessage(head PackHead)  {
//	data,_:=head.serializePackHead()
//	u.conn.Write(data)
//}

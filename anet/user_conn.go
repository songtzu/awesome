package anet

//import (
//	"golang.org/x/net/websocket"
//	"fmt"
//	"awesome/defs"
//)
const newBufferSizeWhenBufferNotEnough  = 10
//最大协议长度100k
const maxPackSize  = 500*1024

var _initialUserId int= -1

func getInitUserId() int  {
	_initialUserId--
	return _initialUserId
}
/*
 * websock.Conn的子类
 */
//type userConn struct {
//	conn   *websocket.Conn
//	userId defs.TypeUserId
//	//房间的句柄
//	room *Room
//}
//func NewPlayer(session *ClientSession,extension interface{}) *Player {
//	return &Player{
//		Session:session,
//		extension:extension,
//	}
//}
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

//func (u *userConn) writeMessage(head PackHead)  {
//	data,_:=head.serializePackHead()
//	u.conn.Write(data)
//}
//func (u *userConn)readLoop() {
//	conn:= u.conn
//	request := make([]byte, 1024)
//	defer conn.Close()
//	var remains = 0
//
//	for {
//		readLen, err := conn.Read(request)
//		remains += readLen
//		if err!=nil{
//			fmt.Println("read error",err)
//		}
//		//socket被关闭了
//		if readLen == 0 {
//			fmt.Println("Client connection close!")
//			break
//		} else {
//			if msg,parsedSize:=parsePackHead(request);msg!=nil{
//				u.processPack(msg)
//				remains -= parsedSize
//				copy(request, request[parsedSize:remains])
//			}else {
//				//包比较大，开辟更大的缓存读数据
//				if newBufferSizeWhenBufferNotEnough*len(request)<=maxPackSize{
//					var largeBuffer   = make([]byte, newBufferSizeWhenBufferNotEnough*len(request))
//					copy(largeBuffer, request[0:remains])
//					request = largeBuffer
//				}else {
//					//一次数据超限，建议关闭连接
//					remains = 0
//				}
//			}
//			////输出接收到的信息
//			//fmt.Println(string(request[:readLen]))
//			//
//			////发送
//			//conn.Write([]byte("World !"))
//		}
//		//request = make([]byte, 128)
//	}
//}

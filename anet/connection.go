package anet

import (
	"golang.org/x/net/websocket"
	"log"
	"net"
	"reflect"
	"sync"
	"time"

	"sync/atomic"

	"errors"

	"runtime/debug"
)

const (
	ConnectionStateConnected = 1
	ConnectionStateClosed    = 2
)

type connectionTyp uint32

const (
	connectionTypeClient connectionTyp = 1
	connectionTypeServer connectionTyp = 2
)

type router struct {
	fun reflect.Value
	msg reflect.Type
}

//抛出的连接对象
type Connection struct {
	connectTimeoutMillisecond int
	connectAddr               string
	state                     int32
	//是客户端还是服务端
	connectionType connectionTyp
	//callback awesome.IFramework
	//address string
	//callback InterfaceNet
	netProtocol NetProtocolType
	/*
	 * if address protocol is tcp
	 *
	 */
	connTcp net.Conn

	/*
	 * if address protocol is ws
	 * then connWebSock should not be nil
	 */
	connWebSock *websocket.Conn
	iConn       InterfaceNet
	//写数据的锁
	mutex sync.RWMutex

	//内置的callback map
	routers map[PackHeadCmd]*router
}

func NewTcpConnection(conn net.Conn, iconn InterfaceNet) *Connection {

	var c = &Connection{
		netProtocol:    NetProtocolTypeTCP,
		connTcp:        conn,
		iConn:          iconn,
		state:          ConnectionStateConnected,
		connectionType: connectionTypeServer,
	}
	return c
}

func NewWebSockConnection(conn *websocket.Conn, iconn InterfaceNet) *Connection {
	con := &Connection{
		netProtocol:    NetProtocolTypeWebSock,
		connWebSock:    conn,
		iConn:          iconn,
		state:          ConnectionStateConnected,
		connectionType: connectionTypeServer,
	}
	return con
}

func (c *Connection) CloseConnWithoutRecon(err error) {
	if atomic.CompareAndSwapInt32(&c.state, ConnectionStateConnected, ConnectionStateClosed) {
		//alog.Err("close conn ")
		if c.netProtocol == NetProtocolTypeTCP {
			err := c.connTcp.Close()
			if err != nil {
				log.Println("error when close Tcp Socket conn", err)
			}
		} else if c.netProtocol == NetProtocolTypeWebSock {
			err := c.connWebSock.Close()
			if err != nil {
				log.Println("error when close WebSocket conn", err)
			}
		}
		if c.iConn != nil {
			//alog.Trace("c.connectionType",c.connectionType)
			if c.iConn.IOnClose(err) && c.connectionType == connectionTypeClient {
				//implement return defines reconnect.
				c.tcpConnect(true)
			}
		} else {
			log.Println("c.iConn interface is not set.")
		}
	}
}
//var (
//
//	outfile, _ = os.Create(os.Args[0]+"_log.txt") // update path for your needs
//	l      = log.New(outfile, "", 0)
//)
func (c *Connection) dispatchMsg(pack *PackHead) {
	//l.Println("dispatchMsg:", pack.SequenceID, string(pack.Body), pack.Cmd)
	//fmt.Println("dispatchMsg", pack, c.connectionType)
	//if c.connectionType==connectionTypeClient {

	if cb := popCallback(pack); cb != nil {
		//log.Println("注册回调的包")
		cb(pack)
		return
	}
	//}
	if pack == nil {
		//alog.Err("msg to dispatch is nil")
		log.Println("msg to dispatch is nil")
		return
	}
	if c.iConn == nil {
		log.Println("callback instance is nil")
		return
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println("panic .", err, string(debug.Stack()))
		}
	}()

	//logdebug(reflect.TypeOf(c.iConn))
	c.iConn.IOnProcessPack(pack)
}

func (c *Connection) WriteMessage(msg *PackHead) (n int, err error) {
	data, _ := msg.SerializePackHead()
	if c == nil {
		log.Println("链接不存在")
		return 0, errors.New("往空链接写入数据")
		//alog.Err("nil connection ", reflect.TypeOf(c.iConn), string(debug.Stack()))
	}
	if c.state != ConnectionStateConnected {
		return 0, errors.New("try to write to a con which not est")
	}
	if c.netProtocol == NetProtocolTypeTCP {
		//logdebug("tcp写出数据")
		return c.connTcp.Write(data)
	} else {
		n, err = c.connWebSock.Write(data)
		//logdebug("ws 协议写出数据",n,err)
		return
	}
}

/*
 * 此方法仅仅用于网络链接的客户端
 */
func (c *Connection) WriteMessageWithCallback(msg *PackHead, cb DefNetIOCallback) (n int, err error) {
	if msg.SequenceID <= 0 {
		msg.SequenceID = allocateNewSequenceId()
	}
	registCallback(msg, cb)
	data, _ := msg.SerializePackHead()
	if c == nil {
		log.Println("链接不存在====")
		return 0, errors.New("往空链接写入数据")
		//alog.Err("nil connection ", reflect.TypeOf(c.iConn), string(debug.Stack()))
	}
	if c.state != ConnectionStateConnected {
		return 0, errors.New("try to write to a con which not est")
	}
	if c.netProtocol == NetProtocolTypeTCP {
		//log.Println("tcp写出数据")
		return c.connTcp.Write(data)
	} else {
		n, err = c.connWebSock.Write(data)
		log.Println("ws 协议写出数据", n, err)
		return
	}
}

/*
 * 此方法仅仅用于网络链接的客户端，服务器端未考虑是否会有bug
 */
func (c *Connection) WriteMessageWaitResponseWithinTimeLimit(msg *PackHead, timeLimitMillisecond int64) (ackMsg *PackHead, isTimeOut bool) {
	msg.SequenceID = allocateNewSequenceId()
	var evtChan = make(chan *PackHead, 1)
	registCallbackWithinTimeLimit(msg, nil, timeLimitMillisecond, evtChan)
	data, _ := msg.SerializePackHead()
	if c == nil {
		log.Println("链接不存在,往空链接写入数据")
		return nil, true
		//alog.Err("nil connection ", reflect.TypeOf(c.iConn), string(debug.Stack()))
	}
	if c.state != ConnectionStateConnected {
		log.Println("连接为空")
		return nil, true
	}
	if c.netProtocol == NetProtocolTypeTCP {

		if n, err := c.connTcp.Write(data); err != nil {
			log.Println("tcp写出数据失败，写出长度", n, "错误内容", err)
		}
		//logdebug("ws 协议写出数据",n,err)

	} else {
		if n, err := c.connWebSock.Write(data); err != nil {
			log.Println("websock写出数据失败，写出长度", n, "错误内容", err)
		}
	}
	select {
	case ackMsg := <-evtChan:
		close(evtChan)
		//log.Println("设置超时的网络IO正常返回", ackMsg)
		return ackMsg, false
	case <-time.After(time.Millisecond * time.Duration(timeLimitMillisecond)):
		//log.Println("设置了超时的网络IO超时返回", msg)
		return nil, true
	}

}

func (c *Connection) WriteBytes(bin []byte, cmd uint32) (n int, err error) {
	msg := &PackHead{magicNum, cmd, 0, uint32(len(bin)), 0, 0, bin}
	return c.WriteMessage(msg)
}

func (c *Connection) PrintNetProtocol() {
	log.Println("是否是TCP", c.netProtocol == NetProtocolTypeTCP)
}

func (c *Connection) WriteSeqBytes(bin []byte, cmd uint32, seq uint32) (n int, err error) {
	msg := &PackHead{magicNum, cmd, seq, uint32(len(bin)), 0, 0, bin}
	return c.WriteMessage(msg)
}

// cmd 	消息id
// f   	处理消息的函数 如: login(session *server.ClientSession, req *protocol.CLogin) (resp proto.Message, err error)
// msg 	消息对应的protobuf请求包类型
func (c *Connection) CbMessageRouter(cmd PackHeadCmd, cb interface{}, msg interface{}) {
	_, ok := c.routers[cmd]
	if ok {
		//logdebug("cmd has registered before", cmd)
		return
	}

	if reflect.TypeOf(cb).Kind() != reflect.Func {
		//logerr("cb should be a function", cmd )
		return
	}

	if reflect.TypeOf(msg).Kind() != reflect.Struct {
		//logerr("msg should be a struct", cmd)
		return
	}
	c.routers[cmd] = &router{
		fun: reflect.ValueOf(cb),
		msg: reflect.TypeOf(msg),
	}
}

func (c *Connection) CbGetFunc(cmd PackHeadCmd) reflect.Value {
	return c.routers[cmd].fun
}

func (c *Connection) CbExist(cmd PackHeadCmd) bool {
	_, ok := c.routers[cmd]
	return ok
}

func (c *Connection) CbGetProto(cmd PackHeadCmd) reflect.Type {
	return c.routers[cmd].msg
}


func (c *Connection) WriteBinary(bin []byte) (n int, err error) {
	if c == nil {
		log.Println("链接不存在")
		return 0, errors.New("往空链接写入数据")
	}
	if c.state != ConnectionStateConnected {
		return 0, errors.New("try to write to a con which not est")
	}
	if c.netProtocol == NetProtocolTypeTCP {
		return c.connTcp.Write(bin)
	} else {
		n, err = c.connWebSock.Write(bin)
		return
	}
}
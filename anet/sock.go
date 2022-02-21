package anet

import (
	"log"
	"net"

	"bufio"

	"errors"
	"fmt"
	"runtime/debug"
)

const initBufferSize = 10240

func StartTcpSvr(address string, instance InterfaceNet) {

	//defer srv.OnClose()
	log.Println(address)
	l, e := net.Listen("tcp", address)
	if e != nil {
		log.Println("[TCPServer] listen error: %v", e)
		panic(e.Error())
	}

	defer l.Close()

	for {
		rw, e := l.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				continue
			}
			log.Println("[TCPServer] accept error: %v", e)
			return
		}
		//c := CreateConnection(address)
		//srv.OnCreateConnection(c)
		//c:=Connection{netProtocol: netProtocolTypeTCP,connTcp:rw, callback:instance}
		//c.startTcp(rw)
		//todo, fork instance and run.
		forkInstance := forkNewInstanceOfInterface(instance)

		c := NewTcpConnection(rw, forkInstance)

		forkInstance.IOnNewConnection(c)
		//newThingInterf.
		//forkInstance:=reflect.New(reflect.TypeOf(instance))
		//v:=forkInstance.(InterfaceNet)
		//instance.IOnNewConnection(c)
		c.startTcp()
	}
}

func (c *Connection) startTcp() bool {
	go c.tcpReadLoop()
	return true
}

func (c *Connection) tcpReadLoop() {
	defer func() {

		if err := recover(); err != nil {
			log.Println(err, string(debug.Stack()))
		}
	}()
	/*
	 * cache buffer is elastic allocate, set max buffer if needed.
	 */
	var cacheBuffer = make([]byte, initBufferSize)
	var remains int = 0
	ioReader := bufio.NewReader(c.connTcp)
	for {
		//c.connTcp.SetReadDeadline(time.Now().Add(time.Second * 10))
		lengthRead, err := ioReader.Read(cacheBuffer[remains:])
		//logdebug("tcp读取到数据",lengthRead, string(cacheBuffer))
		if err != nil {
			log.Println("tcp connection close.", err)
			c.CloseConnWithoutRecon(err)
			return
		}
		remains += lengthRead
		for {
			var pack *PackHead
			var packLength int
			if pack, packLength = parsePackHead(cacheBuffer, remains); packLength > 0 {
				c.dispatchMsg(pack)
				/*
				 * rm buffers which has been read.
				 */
				//fmt.Println(len(cacheBuffer), packLength, remains,lengthRead)
				remains -= packLength
				copy(cacheBuffer, cacheBuffer[packLength:packLength+remains])
				//logdebug(remains, packLength, "裁剪之后的内容:", string(cacheBuffer))

			} else if packLength < 0 {
				/*
				 * todo,
				 * consider close con due to bad magic number.
				 */
				log.Println("connection closed")
				c.CloseConnWithoutRecon(errors.New(fmt.Sprintf("data error, packLength:%d", packLength)))
				return
			} else if packLength == 0 {
				/*
				 * buffer size not enough for a full pack.
				 */
				if remains > len(cacheBuffer)/2 {
					var largeBuffer = make([]byte, 10*len(cacheBuffer))
					cacheBuffer = largeBuffer
				}
				break
			}
		}

	}
}

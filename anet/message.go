package anet

import (
	"encoding/binary"
	"fmt"
	"log"
)

/*
 * 魔数常量
 */
const magicNum uint32 = 20170103
const packHeadLength = 24
const packMaxSize = 50 * 1024

type PackHeadCmd = uint32
type PackHeadSequenceId = uint32
type PackHeadReserveLow = uint32
type PackHeadReserveHigh = uint32

/*
 * websocket版私有协议体定义
 *
 */
type PackHead struct {
	/*
		 * [0,4)
		 * 魔数，包的确认字。
		 * 		tcp协议必需此字节，
				websock按道理是无需此字段设计，沿用惯例，并保持框架对两种协议的兼容设计上考虑，保留此字段
	*/
	magicNum uint32
	/*
	 * [4,8)
	 * cmd enum, 命令字枚举
	 */
	Cmd uint32
	/*
	 * [8,12)
	 * 客户端自增，原样返回给客户端（如果服务器使用websocket)
	 */
	SequenceID uint32
	/*
	 * [12, 16),鉴于websocket本身是应用层协议，已经处理了包的概念，忽略此字段
	 */
	Length uint32
	/*
	 * [16,20)	高位保留字段
	 *  mq中，High 字段用来存储消息发布的status，0表示成功，1表示超时
	 */
	ReserveHigh uint32
	/*
	 *  [20,24)	低位保留字段, mq中存储mq的业务模型。
	 */
	ReserveLow uint32
	/*
	 * [24, 24+length )
	 */
	Body []byte
}

func parsePackHead(data []byte, bufferSize int) (msg *PackHead, length int) {
	/*
	 * 鉴于websock本身处理了包，出现包体长度错误的情况为非法请求
	 */
	//fmt.Println(fmt.Sprintf("data len:%d, bufferSize:%d", len(data), bufferSize))
	if len(data) < packHeadLength || bufferSize < packHeadLength {
		//logdebug(fmt.Errorf("包体长度错误实际长度：%d,包头定义长度,缓存长度%d",len(data) , packHeadLength,bufferSize))
		return nil, 0
	}
	msg = &PackHead{}
	msg.magicNum = binary.BigEndian.Uint32(data[:4])
	if msg.magicNum != magicNum {
		log.Println("魔数不匹配", msg.magicNum, "正确的魔数", magicNum, data)
		return nil, bufferSize
	}
	//else {
	//	logdebug("正确的魔数",data)
	//}
	msg.Cmd = binary.BigEndian.Uint32(data[4:8])
	msg.SequenceID = binary.BigEndian.Uint32(data[8:12])
	msg.Length = binary.BigEndian.Uint32(data[12:16])
	if msg.Length >= packMaxSize {
		log.Println(fmt.Sprintf("pack size %d too large, it should not big than %d.", msg.Length, packMaxSize))
		return nil, bufferSize
	}
	if int(msg.Length)+packHeadLength > bufferSize {
		log.Println("数据未读完", msg.Length, "bufferSize", bufferSize)
		return nil, 0
	}
	msg.ReserveHigh = binary.BigEndian.Uint32(data[16:20])
	msg.ReserveLow = binary.BigEndian.Uint32(data[20:24])

	var body = make([]byte, msg.Length)
	copy(body, data[packHeadLength:packHeadLength+msg.Length])

	msg.Body = body
	if len(msg.Body) != int(msg.Length) {
		log.Println("数据长度不一致")
	}
	return msg, packHeadLength + int(msg.Length)
}

func (this *PackHead) SerializePackHead() (data []byte, length int) {
	if this == nil {
		return []byte{}, 0
	}
	this.Length = uint32(len(this.Body))
	data = make([]byte, this.Length+packHeadLength)
	binary.BigEndian.PutUint32(data[:4], magicNum)
	binary.BigEndian.PutUint32(data[4:8], this.Cmd)
	binary.BigEndian.PutUint32(data[8:12], this.SequenceID)
	binary.BigEndian.PutUint32(data[12:16], this.Length)
	binary.BigEndian.PutUint32(data[16:20], this.ReserveHigh)
	binary.BigEndian.PutUint32(data[20:24], this.ReserveLow)
	copy(data[packHeadLength:], this.Body)
	length = len(data)
	return
}

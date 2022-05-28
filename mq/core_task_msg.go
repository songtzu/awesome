package mq

import (
	"awesome/anet"
)

type AmqMessage struct {
	originalSequenceId         uint32 //发布者传过来的包序。保存记录起来。
	msg                        *anet.PackHead
	sourceConn                 *anet.Connection    //发布者的conn句柄。保存
	srcChan                    chan *anet.PackHead //发布者的chan，此句柄仅用于http模式的发布者。sourceConn为空的时候，srcChan有值。
	createTimestampMillisecond int64
	pushedSubscriberIds        []int //xmqSubImpl.id，推送过的订阅者ID。用于记录推送失败后，推送给其他订阅者。
}

func (a *AmqMessage) response(ackType AmqAckType, pack *anet.PackHead) {
	pack.ReserveHigh = ackType
	pack.SequenceID = a.originalSequenceId //回填sequenceId
	if a.sourceConn != nil {
		a.sourceConn.WriteMessage(pack)
	} else {
		a.srcChan <- pack
	}
	//log.Println("写回数据给发布者,",n,err,string(pack.Body))
	//log.Println("response:", n, err)
}

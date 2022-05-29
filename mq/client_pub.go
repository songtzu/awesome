package mq

import (
	"awesome/anet"
	"errors"
	"log"
)

/*AmqClientPublisher
 * pub is a tcp request client.
 *
 */
//var clientPublisherInstance *AmqClientPublisher = nil
type AmqClientPublisher struct {
	conn *anet.Connection
	id   int
}

func NewClientPublish(url string) (instance *AmqClientPublisher, err error) {
	instance = &AmqClientPublisher{}
	c := anet.NewTcpClientConnect(url, instance, 1000, true)
	instance.conn = c
	if c == nil {
		return nil, errors.New("连接失败")
	}
	return instance, nil
}

/*PubUnreliableToAll
 * 	发布给所有的订阅者，无回包。
 *********/
func (a *AmqClientPublisher) PubUnreliableToAll(data []byte, cmd anet.PackHeadCmd) (err error) {
	msg := &anet.PackHead{ReserveLow: AmqCmdDefUnreliable2All, Cmd: cmd,
		SequenceID: 0,
		Length:     uint32(len(data)), Body: data}
	//str:=string(data)
	_, err = a.conn.WriteMessage(msg)
	return err
}

/*PubReliableToRandomOne
 * 随机选择一个订阅者发布，等待订阅者回包
 *		或者超时。
 *********/
func (a *AmqClientPublisher) PubReliableToRandomOne(data []byte, cmd anet.PackHeadCmd) (pack *anet.PackHead, isTimeout bool) {
	msg := &anet.PackHead{ReserveLow: AmqCmdDefReliable2RandomOne, Cmd: cmd,
		SequenceID: 0,
		Length:     uint32(len(data)), Body: data}
	return a.conn.WriteMessageWaitResponseWithinTimeLimit(msg, defaultTimeoutMillisecond+defaultTimeoutDelay)

}

/*PubReliableToSpecOne
 * 消息会扇出给所有订阅者，但是只有特定的一个或者零个订阅者能处理此消息。等待订阅者回包
 *		或者超时。
 *********/
func (a *AmqClientPublisher) PubReliableToSpecOne(data []byte, cmd anet.PackHeadCmd) (pack *anet.PackHead, isTimeout bool) {
	msg := &anet.PackHead{ReserveLow: AmqCmdDefReliable2SpecOne, Cmd: cmd,
		SequenceID: 0,
		Length:     uint32(len(data)), Body: data}
	return a.conn.WriteMessageWaitResponseWithinTimeLimit(msg, defaultTimeoutMillisecond+defaultTimeoutDelay)
}

/*PubUnreliableToRandomOne
 * 随机选择一个订阅者发布，无超时或回包。错误仅代表msg是否写出到内核。不代表xmq节点是否成功处理。也不代表消费者是否处理。
 *		或者超时。
 *********/
func (a *AmqClientPublisher) PubUnreliableToRandomOne(data []byte, cmd anet.PackHeadCmd) (err error) {
	msg := &anet.PackHead{ReserveLow: AmqCmdDefReliable2RandomOne, Cmd: cmd,
		SequenceID: 0,
		Length:     uint32(len(data)), Body: data}
	_, err = a.conn.WriteMessage(msg)
	return err
}

func (a *AmqClientPublisher) IOnInit(connection *anet.Connection) {

}

func (a *AmqClientPublisher) IOnProcessPack(pack *anet.PackHead, connection *anet.Connection) {

}

/*IOnClose
 * this interface SHOULD NOT CALL close.
 */
func (a *AmqClientPublisher) IOnClose(err error) (tryReconnect bool) {
	return true
}

func (a *AmqClientPublisher) IOnConnect(isOk bool) {

}

func (a *AmqClientPublisher) IOnNewConnection(connection *anet.Connection) {
	log.Println("new connection")
	//fork:=&pubImpl{reliableCallback:a.reliableCallback, id:newConnId(), conn:connection}
	a.conn = connection
	a.id = newConnId()

	//test := &anet.PackHead{Cmd: AMQCmdDefPub, Length: uint32(len([]byte("hello"))), Body: []byte("hello")}
	//a.conn.WriteMessage(test)
	//a.id = newConnId()
	//a.conn = connection
}

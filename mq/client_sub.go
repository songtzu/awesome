package mq

import (
	"awesome/anet"
	"encoding/json"
	"google.golang.org/protobuf/proto"
	"log"
	"reflect"
)

type AmqClientSubscriber struct {
	conn        *anet.Connection
	cb          AMQCallback
	lastMessage anet.PackHead
}

func NewClientSubscriber(bindAddress string, cb AMQCallback) *AmqClientSubscriber {
	impl := &AmqClientSubscriber{cb: cb}
	c := anet.NewTcpClientConnect(bindAddress, impl, 1000, true)
	impl.conn = c
	return impl
}

func (a *AmqClientSubscriber) TopicSubscription(topics []AMQTopic) (n int, err error) {
	t := &AMQProtocolSubTopic{Topics: topics}
	bin, _ := json.Marshal(t)
	//log.Println(string(bin), "size:",len(bin))
	msg := &anet.PackHead{Cmd: AMQCmdDefSubTopic,
		SequenceID: 0,
		Length:     uint32(len(bin)), Body: bin}

	return a.conn.WriteMessage(msg)
}

func (a *AmqClientSubscriber) IOnInit(connection *anet.Connection) {

}

func (a *AmqClientSubscriber) IOnProcessPack(pack *anet.PackHead) {
	//log.Println("网络层收到数据",string(pack.Body),pack.SequenceID)
	if a.conn.CbExist(pack.Cmd) {
		hd := a.conn.CbGetFunc(pack.Cmd)
		t := a.conn.CbGetProto(pack.Cmd)
		v := reflect.New(t)
		if err := proto.Unmarshal(pack.Body, v.Interface().(proto.Message)); err == nil {
			hd.Call([]reflect.Value{reflect.ValueOf(a), v})
		} else {
			log.Panic("protocol  unmarshal fail: ", err, pack.Cmd)
		}
	}
	if a.cb != nil {
		a.lastMessage = *pack
		a.cb(pack)
	}
}

/*
 * this interface SHOULD NOT CALL close.
 */
func (a *AmqClientSubscriber) IOnClose(err error) (tryReconnect bool) {
	log.Println("IOnClose订阅连接关闭")
	return true
}

//func (a *subImpl) IWrite(msg interface{}, ph *net.PackHead){
//
//}

func (a *AmqClientSubscriber) IOnConnect(isOk bool) {

}

func (a *AmqClientSubscriber) IOnNewConnection(connection *anet.Connection) {

}

func (a *AmqClientSubscriber) Response(msg []byte) error {
	a.lastMessage.Body = msg
	_, err := a.conn.WriteMessage(&a.lastMessage)
	return err
}

func (a *AmqClientSubscriber) RegistCallback(msg []byte) error {

	a.lastMessage.Body = msg
	_, err := a.conn.WriteMessage(&a.lastMessage)
	return err
}

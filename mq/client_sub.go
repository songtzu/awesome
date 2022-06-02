package mq

import (
	"awesome/anet"
	"encoding/json"
	"errors"
	"google.golang.org/protobuf/proto"
	"log"
	"reflect"
)

type AmqClientSubscriber struct {
	conn        *anet.Connection
	cb          anet.DefNetIOCallbackWithArrResponse
	lastMessage *anet.PackHead
	topics      []AMQTopic
}

func NewClientSubscriber(bindAddress string, cb anet.DefNetIOCallbackWithArrResponse) *AmqClientSubscriber {
	impl := &AmqClientSubscriber{cb: cb, topics: make([]AMQTopic, 0)}
	c := anet.NewTcpClientConnect(bindAddress, impl, 1000, true)
	impl.conn = c
	return impl
}

func (a *AmqClientSubscriber) TopicSubscription(topics []AMQTopic) (err error) {
	t := &AMQProtocolSubTopic{Topics: topics}
	bin, _ := json.Marshal(t)
	//log.Println(string(bin), "size:",len(bin))
	msg := &anet.PackHead{Cmd: AMQCmdDefSubTopic,
		SequenceID: 0,
		Length:     uint32(len(bin)), Body: bin}

	pack, timeout := a.conn.WriteMessageWaitResponseWithinTimeLimit(msg, 200)
	if timeout {
		log.Println("订阅超时了")
		return errors.New("sub topic failed")
	}
	ack := &AMQProtocolSubTopicAck{}
	if err = json.Unmarshal(pack.Body, ack); err != nil {
		log.Println(err)
		return err
	}
	if ack.Status == 0 {
		for _, v := range topics {
			if !anet.Contains(a.topics, v) {
				a.topics = append(a.topics, v)
			}
		}
	}
	log.Println("订阅结果", ack)
	return nil
}

func (a *AmqClientSubscriber) IOnInit(connection *anet.Connection) {

}

func (a *AmqClientSubscriber) IOnProcessPack(pack *anet.PackHead, connection *anet.Connection) {
	log.Printf("网络层收到数据%s,len:%d,seqId:%d", string(pack.Body), len(pack.Body), pack.SequenceID)

	if a.conn.CbExist(pack.Cmd) {
		hd := a.conn.CbGetFunc(pack.Cmd)
		t := a.conn.CbGetProto(pack.Cmd)
		v := reflect.New(t)
		if err := proto.Unmarshal(pack.Body, v.Interface().(proto.Message)); err == nil {
			hd.Call([]reflect.Value{reflect.ValueOf(a), v})
		} else {
			log.Panic("protocol unmarshal fail: ", err, pack.Cmd)
		}
	}
	if a.cb != nil {
		log.Println("IOnProcessPack执行预先注册的回调...")
		a.lastMessage = pack
		arr, cmd := a.cb(pack)
		if len(arr) > 0 {
			connection.WriteBytes(arr, cmd, pack.SequenceID)
		}
	}
}

/*
 * this interface SHOULD NOT CALL close.
 */
func (a *AmqClientSubscriber) IOnClose(err error) (tryReconnect bool) {
	log.Println("IOnClose订阅连接关闭")
	//a.conn.TCPReConnect()
	return true
}

func (a *AmqClientSubscriber) IOnConnect(isReconnect bool) {
	log.Println("AmqClientSubscriber,isReconnect", isReconnect)
	if isReconnect {
		a.TopicSubscription(a.topics)
	}
}

func (a *AmqClientSubscriber) IOnNewConnection(connection *anet.Connection) {
	log.Println("IOnNewConnection")
}

//
//func (a *AmqClientSubscriber) Response(msg []byte) error {
//	a.lastMessage.Body = msg
//	_, err := a.conn.WriteMessage(&a.lastMessage)
//	return err
//}

func (a *AmqClientSubscriber) RegisterCallback(msg []byte) error {
	log.Println("AmqClientSubscriber===>RegistCallback")
	a.lastMessage.Body = msg
	_, err := a.conn.WriteMessage(a.lastMessage)
	return err
}

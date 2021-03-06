package test

import (
	"awesome/anet"
	"awesome/mq"
	"log"
	"sync"
	"testing"
	"time"
)

var instance *mq.AmqClientSubscriber

func testSubClientCb(pack *anet.PackHead) {
	log.Println("订阅者，收到订阅消息", pack)
	instance.Response([]byte("yes we got it."))
}

func TestSubClient(t *testing.T) {
	log.Println("创建订阅者的客户端")
	instance = mq.NewClientSubscriber("127.0.0.1:19999", testSubClientCb)
	instance.TopicSubscription([]mq.AMQTopic{1000, 1001, 1002})
	time.Sleep(10 * time.Minute)
}

func TestSyncMapDeleteInRange(t *testing.T) {
	var m sync.Map
	//Store
	m.Store(1, "1a")
	m.Store(2, "2b")
	m.Store(3, "3b")
	m.Store(4, "4b")
	m.Store(5, "5b")
	m.Store(6, "6b")
	m.Store(7, "7b")
	m.Store(8, "8b")
	m.Store(9, "9b")
	m.Range(func(key, value interface{}) bool {
		k := key.(int)
		if k == 5 || k == 2 {
			m.Delete(key)
		}
		return true
	})

	log.Println(m.Load(5))
	log.Println(m.Load(3))

}

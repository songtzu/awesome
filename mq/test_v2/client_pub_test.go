package test_v2

import (
	"awesome/anet"
	"awesome/mq"
	"fmt"
	"log"
	"testing"
	"time"
)

//var instancePub *mq.AmqClientPublisher

var clientPubTest = &anet.TestInfo{Start: time.Now(), CurrentTotalCount: 0, ThreadCount: 64, SetCountEachThread: 50000}

func TestClientPubReliable2RandomOneMessageNormal(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	clientPubTest.UpdateTotalSetCount()

	for i := 0; i < clientPubTest.ThreadCount; i++ {
		go worker()
	}
	log.Println("======运行结束")
	time.Sleep(1 * time.Hour)
}

func worker() {
	if pub, err := mq.NewClientPublish("127.0.0.1:8888"); err == nil {
		publishToRandomOne(pub)

	} else {
		log.Println(err, "运行错误")
	}
}

func publishToRandomOne(pub *mq.AmqClientPublisher) {
	for i := 0; i < clientPubTest.SetCountEachThread; i++ {
		str := fmt.Sprintf("发送第:%d次请求", i)
		_, isTimeout := pub.PubReliableToRandomOne([]byte(str), 1001)
		if isTimeout {
			log.Println("超时")
			clientPubTest.AddFailCount()
		} else {
			clientPubTest.AddPassCount()
		}
	}
	log.Println("执行到这里了")
	clientPubTest.TryPrint()

}

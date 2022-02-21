package mq

import (
	"awesome/anet"
	"log"
	"testing"
	"time"
)

func testForSubCb(message *anet.PackHead) {
	log.Println("head:", message)
	log.Println("sub callback", string(message.Body))
}

func TestXSub(t *testing.T) {
	log.Println("AMQ SUB START")
	instance := NewXSub("127.0.0.1:19999", testForSubCb)
	instance.TopicSubscription([]AMQTopic{1001, 1002})
	time.Sleep(1 * time.Second)
	time.Sleep(5 * time.Minute)
}

func TestXSub2(t *testing.T) {
	log.Println("AMQ SUB START")
	instance := NewXSub("127.0.0.1:19999", testForSubCb)
	instance.TopicSubscription([]AMQTopic{1003})
	time.Sleep(1 * time.Minute)
}

package mq

import (
	"awesome/anet"
	"log"
	"testing"
	"time"
)

func test_cb(message *anet.PackHead) {
	log.Println("head:", message)
	log.Println("sub callback", string(message.Body))
}

func TestAMQSUB(t *testing.T) {
	log.Println("AMQ SUB START")
	sub := NewAMQSub(":7777", test_cb)
	time.Sleep(1 * time.Second)
	sub.TopicSubscription([]AMQTopic{1001, 1002})
	time.Sleep(1 * time.Minute)
}

func TestAMQSUB2(t *testing.T) {
	log.Println("AMQ SUB START")
	sub := NewAMQSub(":7777", reliableCallback)
	time.Sleep(1 * time.Second)
	sub.TopicSubscription([]AMQTopic{1000})
	time.Sleep(1 * time.Minute)
}

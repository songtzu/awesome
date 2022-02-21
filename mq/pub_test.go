package mq

import (
	"encoding/json"
	"log"
	"testing"
	"time"
)

func TestAMQPub(t *testing.T) {
	log.Println("AMQ SUB START")
	mq := GetAMQPubInstance(":7777")
	log.Println("begin to sleep.")
	time.Sleep(15 * time.Second)
	log.Println("public message")
	mq.PubReliable2RandomOneMessage([]byte("1001,this is test content for AMQ"), 10)
	time.Sleep(1 * time.Hour)
}

func TestSubStruct(t *testing.T) {
	msg := &AMQProtocolSubTopic{}
	str := `{"topics":[1001,1002]}`
	bin := []byte(str)
	if err := json.Unmarshal(bin, msg); err != nil {
		log.Println("error when recieve sub topic action", err)
	} else {
		log.Println(msg.Topics)
	}
}

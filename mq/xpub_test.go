package mq

import (
	"log"

	"testing"
	"time"
)

func TestXPub(t *testing.T) {
	log.Println("AMQ SUB START")
	instance := NewXPub("127.0.0.1:18888")
	instance.MessagePub(1001, []byte("hello world, this is topic about 1001"))
	time.Sleep(1 * time.Second)
	time.Sleep(1 * time.Minute)
}

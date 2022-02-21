package mq

import (
	"log"
	"testing"
	"time"
)

func TestXmq(t *testing.T) {
	log.Println("AMQ SUB START")
	NewXmq("127.0.0.1:18888", "127.0.0.1:19999")

	time.Sleep(10 * time.Minute)
}

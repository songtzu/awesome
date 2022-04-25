package test_v2

import (
	"awesome/mq"
	"fmt"
	"log"
	"testing"
	"time"
)

const (
	xPublicAddress = "127.0.0.1:8888"
	xSubscribeAddress = "127.0.0.1:9999"
)

func TestMqSvr(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("message queue test& dev code")
	mq.NewXmq(xPublicAddress,xSubscribeAddress)
	go mq.StartHttpForMQ(":6666",xPublicAddress)
	for ; ;  {
		time.Sleep(1*time.Hour)
	}

}

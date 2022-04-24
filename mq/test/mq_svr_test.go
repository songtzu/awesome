package test

import (
	"awesome/mq"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestMqSvr(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("message queue test& dev code")
	mq.NewXmq("127.0.0.1:8888","127.0.0.1:9999")
	for ; ;  {
		time.Sleep(1*time.Hour)
	}

}

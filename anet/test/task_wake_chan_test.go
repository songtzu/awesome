package test

import (
	"log"
	"testing"
	"time"
)

func TestChan(t *testing.T) {
	var c chan int = make(chan int, 1)
	go taskChanWrite(c)
	go taskChanRead(c)
	time.Sleep(1 * time.Minute)
}

func taskChanWrite(vc chan int) {
	for i := 0; i < 1000000000000; i++ {
		if len(vc) == 0 {
			vc <- i
		}

		time.Sleep(100 * time.Millisecond)
		log.Println("write:", i)
	}
}

func taskChanRead(vc chan int) {
	for true {
		i := <-vc
		log.Println("read:", i)
		time.Sleep(1 * time.Second)
	}
}

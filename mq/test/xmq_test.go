package test

import (
	"awesome/anet"
	"awesome/mq"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var (
	subAddr       = ":19999"
	pubAddr       = ":18888"
	subId   int32 = -1
)

func TestPubReliable2RandomOneMessage(t *testing.T) {

	var wg sync.WaitGroup
	var recvStatistic = map[int32]int{}

	var success int32 = 0
	var timeoutNum int32 = 0
	newPubSend1001 := func(n ...int) {
		pub, err := mq.NewClientPublish(pubAddr)
		if err != nil {
			log.Panic(err)
		}
		var c = 1
		if len(n) > 0 {
			c = n[0]
		}
		if c == 0 {
			c = 1
		}
		for i := 0; i < c; i++ {
			_, isTimeout := pub.PubReliableToRandomOne([]byte("one client have"), 1001)
			if isTimeout {
				atomic.AddInt32(&timeoutNum, 1)
			} else {
				atomic.AddInt32(&success, 1)
			}
		}
	}

	Recv1001 := func() {
		var instance *mq.AmqClientSubscriber
		id := atomic.AddInt32(&subId, 1)
		instance = mq.NewClientSubscriber(subAddr, func(head *anet.PackHead) {

			if head.ReserveLow > 0 {
				err := instance.Response([]byte(fmt.Sprintf("%d %d 回包", id, head.Cmd)))
				if err != nil {
					log.Printf("err: %v", err)
				}
			}
			recvStatistic[id]++

			wg.Done()
		})
		instance.TopicSubscription([]mq.AMQTopic{1001})
	}

	wg.Add(1)
	go func() {
		newPubSend1001()
	}()

	go func() {
		Recv1001()
	}()

	wg.Wait()
	log.Println("---------------------------")
	// 多个消费者
	wg.Add(1)
	go func() {
		newPubSend1001()
	}()

	go func() {
		Recv1001()
	}()
	go func() {
		Recv1001()
	}()
	wg.Wait()
	log.Printf("多消费ok")

	now := time.Now()

	// 测试随机
	wg.Add(10 * 1000)
	go func() {
		for i := 0; i < 10; i++ {
			go newPubSend1001(1000)
		}
	}()
	wg.Wait()
	for k, v := range recvStatistic {
		log.Printf("消费id : %d   消费数:%d  ", k, v)
	}

	log.Printf("10*1000 消息 10个生产者  3个消费   sub:%s", time.Now().Sub(now))
	if atomic.LoadInt32(&success) == 10*1000 {
		log.Println("test success")
	} else {
		log.Printf("test failed   success %d  timeout %d", success, timeoutNum)
	}

}

func init() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func TestPubReliable2SpecOneMessage(t *testing.T) {

	var (
		wg sync.WaitGroup
	)

	const (
		subCount = 5
		pubCount = 10
		message  = 10
	)

	type s struct {
		se int32
		re int32
	}

	var recvStatistic = [subCount]s{}

	var totalRecv int32 = 0

	rand.Seed(time.Now().UnixNano())

	newPubSend1001 := func(n ...int) {
		pub, err := mq.NewClientPublish(pubAddr)
		if err != nil {
			log.Panic(err)
		}
		var c = 1
		if len(n) > 0 {
			c = n[0]
		}
		if c == 0 {
			c = 1
		}
		for i := 0; i < c; i++ {
			var xfid = rand.Intn(subCount)
			atomic.AddInt32(&recvStatistic[xfid].se, 1)

			result, isTimeout := pub.PubReliableToSpecOne([]byte(fmt.Sprintf("%d", xfid)), 1001)
			if isTimeout {
				log.Printf("is timeout %d", xfid)
			} else {
				log.Printf("消费成功  resp: %+v", string(result.Body))
			}
		}
	}

	Recv1001 := func() {
		var instance *mq.AmqClientSubscriber
		id := atomic.AddInt32(&subId, 1)
		instance = mq.NewClientSubscriber(subAddr, func(head *anet.PackHead) {

			atomic.AddInt32(&totalRecv, 1)

			if string(head.Body) == fmt.Sprintf("%d", id) {
				err := instance.Response([]byte(fmt.Sprintf("%d %d 回包", id, head.Cmd)))
				if err != nil {
					log.Printf("err: %v", err)
				}
				atomic.AddInt32(&recvStatistic[int(head.Body[0]-'0')].re, 1)
				wg.Done()
			} else {
				log.Println(id, "收到", string(head.Body))
			}

		})
		instance.TopicSubscription([]mq.AMQTopic{1001})
	}

	// 10个消费者 一共10*1000个消息  随机处理
	now := time.Now()

	go func() {
		for i := 0; i < subCount; i++ {
			go Recv1001()
		}
	}()
	// 测试随机
	wg.Add(pubCount * message)
	go func() {
		for i := 0; i < pubCount; i++ {
			go newPubSend1001(message)
		}
	}()

	wg.Wait()
	for k, v := range recvStatistic {
		log.Printf("消费id : %d   消费数:%d  ", k, v)
	}

	log.Printf("%d*%d 消息 %d个生产者  %d个消费   sub:%s", pubCount, message, pubCount, subCount, time.Now().Sub(now))
	log.Printf("total message %d", atomic.LoadInt32(&totalRecv))
}

func TestPubUnreliable2AllMessage(t *testing.T) {
	send := int32(0)
	recn := int32(0)
	newPubSend1001 := func(n ...int) {
		pub, err := mq.NewClientPublish(pubAddr)
		if err != nil {
			log.Panic(err)
		}
		var c = 1
		if len(n) > 0 {
			c = n[0]
		}
		if c == 0 {
			c = 1
		}
		for i := 0; i < c; i++ {

			atomic.AddInt32(&send, 1)
			err := pub.PubUnreliableToAll([]byte(fmt.Sprintf("%d", 1)), 1001)
			if err != nil {
				log.Printf("send 2 all message %v", err)
			}
		}
	}
	Recv1001 := func() {
		instance = mq.NewClientSubscriber(subAddr, func(head *anet.PackHead) {
			atomic.AddInt32(&recn, 1)
		})
		n, err := instance.TopicSubscription([]mq.AMQTopic{1001})
		if err != nil {
			log.Printf("n %d err %v", n, err)
		}
	}

	go newPubSend1001(100)

	for i := 0; i < 10; i++ {
		go Recv1001()
	}

	var ti = time.NewTimer(time.Minute)
	for {
		select {
		case <-time.After(time.Second):
			if atomic.LoadInt32(&send) == 100 && atomic.LoadInt32(&recn) == 100*10 {
				log.Printf("test success")
				return
			}
			log.Printf("%d %d", atomic.LoadInt32(&send), atomic.LoadInt32(&recn))
		case <-ti.C:
			log.Printf("test failed 1 minute not success")
			return

		}
	}

}

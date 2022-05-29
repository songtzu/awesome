package test_v2

import (
	"awesome/mq"
	"fmt"
	"log"
	"testing"
	"time"
)

var instancePub *mq.AmqClientPublisher

/*************
 * 测试 有回包的随机选择订阅者的模式，有订阅者。
 ********************/
func TestClientPubReliable2RandomOneMessageNormal(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error
	if instancePub, err = mq.NewClientPublish("127.0.0.1:8888"); err == nil {
		publishToRandomOne()
		//time.Sleep(1 * time.Second)
		//
		//result, isTimeout = instancePub.PubReliableToRandomOne([]byte("hello world+++++++222"), 1001)
		//if isTimeout {
		//	log.Println("超时")
		//} else {
		//	log.Printf("222===>%+v\n", result)
		//	log.Println(string(result.Body))
		//}

	} else {
		log.Println(err, "运行错误")
	}
	log.Println("======运行结束")
	time.Sleep(1 * time.Hour)
}

func publishToRandomOne() {
	var timeoutCount = 0
	var rightCount = 0
	var totalCount = 1000
	start := time.Now()
	for i := 0; i < totalCount; i++ {
		str := fmt.Sprintf("发送第:%d次请求", i)
		result, isTimeout := instancePub.PubReliableToRandomOne([]byte(str), 1001)
		if isTimeout {
			log.Println("超时")
			timeoutCount += 1
		} else {
			//log.Printf("111===>%+v\n", result)
			rightCount += 1
			log.Println(string(result.Body))
		}
	}
	costTime := time.Now().Sub(start).Milliseconds()
	log.Printf("totalCount:%d, timeoutCount:%d, rightCount:%d, time cost ms:%d, avg:%d", totalCount, timeoutCount, rightCount, costTime, costTime/int64(totalCount))
}

func publishToAll() {
	err := instancePub.PubUnreliableToAll([]byte("hello world+++++++111"), 1001)
	if err != nil {
		log.Printf("publishToAll:%s", err.Error())
	}
}

package test_v2

import (
	"awesome/mq"
	"fmt"
	"log"
	"testing"
	"time"
)

var _instance *mq.AmqClientPublisher

/*************
 * 测试 有回包的随机选择订阅者的模式，有订阅者。
 ********************/
func TestClientPubToAll(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error
	if _instance, err = mq.NewClientPublish("127.0.0.1:8888"); err == nil {
		publishToAll2()
	} else {
		log.Println(err, "运行错误")
	}
	log.Println("======运行结束")
	time.Sleep(1 * time.Hour)
}

const totalCount = 1000

func publishToAll2() {
	var timeoutCount = 0
	var rightCount = 0
	start := time.Now()
	for i := 0; i < totalCount; i++ {
		str := fmt.Sprintf("发送第:%d次请求", i)
		_instance.PubUnreliableToAll([]byte(str), 1001)
		log.Println(str)
	}
	costTime := time.Now().Sub(start).Milliseconds()
	log.Printf("totalCount:%d, timeoutCount:%d, rightCount:%d, time cost ms:%d, avg:%f", totalCount, timeoutCount, rightCount, costTime, float64(costTime)/float64(totalCount))
}

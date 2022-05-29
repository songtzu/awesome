package test

import (
	"awesome/mq"
	"log"
	"testing"
	"time"
)

var instancePub *mq.AmqClientPublisher

/*************
 * 测试 有回包的随机选择订阅者的模式，有订阅者。
 ********************/
func TestClientPubReliable2RandomOneMessageNormal(t *testing.T) {
	var err error
	if instancePub, err = mq.NewClientPublish("127.0.0.1:18888"); err == nil {
		result, isTimeout := instancePub.PubReliableToRandomOne([]byte("hello world+++++++222"), 1001)
		if isTimeout {
			log.Println("超时")
		} else {
			log.Printf("%+v\n", result)
			log.Println(string(result.Body))
		}

	} else {
		log.Println(err, "运行错误")
	}
	log.Println("======运行结束")
	time.Sleep(1 * time.Hour)
}

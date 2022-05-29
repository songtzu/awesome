package test

import (
	"awesome/anet"
	"awesome/mq"
	"log"
	"testing"
	"time"
)

/*************
 * 测试 有回包的随机选择订阅者的模式，无订阅者。
 ********************/
func TestClientPubReliable2RandomOneMessageNoSubTimeout(t *testing.T) {
	var err error
	if instancePub, err = mq.NewClientPublish("127.0.0.1:18888"); err == nil {
		result, isTimeout := instancePub.PubReliableToRandomOne([]byte("hello world---------------11"), 10001)
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

func test_no_sub_timeout_sub_client_cb(pack *anet.PackHead) {
	log.Println("订阅者，收到订阅消息", pack)
	instance.Response([]byte("yes we got it."))
}

func TestClientSub_No_Sub_Timeout(t *testing.T) {
	log.Println("创建订阅者的客户端", time.Now().UnixNano())
	instance = mq.NewClientSubscriber("127.0.0.1:19999", test_no_sub_timeout_sub_client_cb)
	instance.TopicSubscription([]mq.AMQTopic{1000, 1001, 1002})
	time.Sleep(10 * time.Minute)
}

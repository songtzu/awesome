package test

import (
	"awesome/anet"
	"awesome/mq"
	"log"
	"testing"
	"time"
)

/*************
 * 测试 有回包的随机选择订阅者的模式，订阅者,但是订阅者处理超时。
 ********************/
func TestClientPubReliable2RandomOneMessageSubTimeout(t *testing.T) {
	var err error
	if instancePub, err = mq.NewClientPublish("127.0.0.1:18888"); err == nil {
		result, isTimeout := instancePub.PubReliableToRandomOne([]byte("hello world123444555"), 1001)
		if isTimeout {
			log.Println("mq客户端自超时判断")
		} else if result != nil && result.ReserveHigh == mq.AmqAckTypeTimeout {
			log.Printf("mq的中间件超时%+v\n", result)
			//log.Println(string(result.Body))
		} else {
			log.Printf("mq正常返回%v+", result)
		}

	} else {
		log.Println(err, "运行错误")
	}
	log.Println("======运行结束")
	time.Sleep(1 * time.Hour)
}

func testRandomOneSubTimeoutSubClientCB(pack *anet.PackHead) {
	log.Println("订阅者，收到订阅消息", pack)
	time.Sleep(4 * time.Second)
	instance.Response([]byte("yes we got it."))
}

func TestClientSubReliable2RandomOneSubTimeout(t *testing.T) {
	log.Println("创建订阅者的客户端", time.Now().UnixNano())
	instance = mq.NewClientSubscriber("127.0.0.1:19999", testRandomOneSubTimeoutSubClientCB)
	instance.TopicSubscription([]mq.AMQTopic{1000, 1001, 1002})
	time.Sleep(10 * time.Minute)
}

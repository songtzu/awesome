package test

import (
	"awesome/anet"
	"awesome/mq"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)


var readLock sync.RWMutex
var readCount = 0
var writeLock sync.Mutex
var writeCount = 0

/*************
 * 测试 有回包的随机选择订阅者的模式，订阅者,但是订阅者处理超时。
 ********************/
func TestClientPubReliable2RandomOneMessageSubBenchmark(t *testing.T) {
	//f, err := os.OpenFile("publisher_log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	//if err != nil {
	//	log.Fatalf("error opening file: %v", err)
	//}
	//defer f.Close()
	//
	//log.SetOutput(f)
	log.Println("This is a test log entry")
	for i := 0; i< 1000 ; i++  {
		go randomBenchmarkWorker(i)
	}
	go randomBenchmarkPrint()
	time.Sleep(1*time.Hour)
}

func randomBenchmarkPrint()  {
	for ; ;  {
		time.Sleep(1*time.Second)
		var read = 0
		var write = 0
		writeLock.Lock()
		write = writeCount
		writeLock.Unlock()

		readLock.Lock()
		read = readCount
		readLock.Unlock()

		log.Println("写出:",write,",读入:",read)
	}
}

func randomBenchmarkWorker(index int)  {
	var err error
	var cmd uint32 = 1001
	if index%2==0{
		cmd = 1002
	}
	if instancePub, err = mq.NewClientPublish("127.0.0.1:18888"); err == nil {
		for i:=0;i<1000 ; i++ {

			//time.Sleep(1*time.Millisecond)
			str:=fmt.Sprintf("发布数据,go程:%d,标号:%d",index,i)
			result, isTimeout := instancePub.PubReliable2RandomOneMessage([]byte(str), cmd)
			writeLock.Lock()
			writeCount++
			//log.Println("writeCount:",writeCount)
			writeLock.Unlock()
			if isTimeout  {
				log.Println("mq客户端自超时判断",i,"消息内容:",str)
			} else if result!=nil && result.ReserveHigh==mq.AmqAckTypeTimeout{
				log.Printf("中间件超时%+v\n", result)
				//log.Println(string(result.Body))
			}else {
				//log.Printf("mq正常返回%v+",result)
				readLock.Lock()
				readCount++
				readLock.Unlock()
			}
		}
	} else {
		log.Println(err, "运行错误")
	}
}


func testRandomOneSubBenchmarkSubClientCB(pack *anet.PackHead) {
	//log.Println("订阅者，收到订阅消息", string(pack.Body),pack.SequenceID)
	//time.Sleep(4*time.Second)
	instance.Response([]byte(fmt.Sprintf("订阅者回复,%s",string(pack.Body))))
}

func TestClientSubReliable2RandomOneSubOneBenchmark(t *testing.T) {
	log.Println("创建订阅者的客户端",time.Now().UnixNano())
	instance = mq.NewClientSubscriber("127.0.0.1:19999", testRandomOneSubBenchmarkSubClientCB)
	instance.TopicSubscription([]mq.AMQTopic{1000, 1001})
	time.Sleep(10 * time.Minute)
}


func TestClientSubReliable2RandomOneSubTwoBenchmark(t *testing.T) {
	log.Println("创建订阅者的客户端",time.Now().UnixNano())
	instance = mq.NewClientSubscriber("127.0.0.1:19999", testRandomOneSubBenchmarkSubClientCB)
	instance.TopicSubscription([]mq.AMQTopic{  1002})
	time.Sleep(10 * time.Minute)
}
package mq

import (
	"log"
	"net/http"
	"time"
)

/*StartHttpForMQ
 * xmq自身创建一个发布者，与xmq连接，在http服务器与xmq之间做桥接。
 */
func StartHttpForMQ(httpAddress string, pubAddress string)  {
	time.Sleep(10*time.Millisecond)
	log.Println("StartHttpForMQ",httpAddress)
	startBridgePublishClient(pubAddress)
	http.HandleFunc("/api/show_status", showStatus)
	http.HandleFunc("/api/publish", publishDefaultMessage)//AMQCmdDefPub
	//http.HandleFunc("/api/publish/unreliable_all", publishDefaultMessage)//AmqCmdDefUnreliable2All
	//http.HandleFunc("/api/publish/unreliable_rand_one", publishDefaultMessage)//AmqCmdDefUnreliable2All
	//http.HandleFunc("/api/publish/reliable_rand_one", publishDefaultMessage)//AmqCmdDefUnreliable2All
	//http.HandleFunc("/api/publish/reliable_spec_one", publishDefaultMessage)//AmqCmdDefUnreliable2All
	err:=http.ListenAndServe(httpAddress,nil)
	log.Printf("web server start result :%v",err)
}
var bridgePublishClient *AmqClientPublisher

func startBridgePublishClient(publishAddr string) (err error) {
	bridgePublishClient, err = NewClientPublish(publishAddr)
	return err
}


type generalResponse struct {
	Code int
	Message string
	Data interface{}
}

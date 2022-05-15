package mq

import (
	"fmt"
	"os"

	//"fmt"
	"log"
	"net/http"
	"time"
)

/*StartHttpForMQ
 * xmq自身创建一个发布者，与xmq连接，在http服务器与xmq之间做桥接。
 */
func StartHttpForMQ(httpAddress string, pubAddress string) {
	time.Sleep(100 * time.Millisecond)
	log.Println("StartHttpForMQ", httpAddress)
	err := startBridgePublishClient(pubAddress)
	if err != nil {
		log.Println("启动错误", err)
		os.Exit(-20)
	}
	http.HandleFunc("/", indexHandler)
	//http://127.0.0.1:9876/api/show_status
	http.HandleFunc("/api/show_status", showStatus)
	//http://127.0.0.1:9876/api/publish?cmd=1001&action=12&body=%E5%BE%88%E9%95%BF%E7%9A%84%E7%94%B5%E5%BD%B1
	http.HandleFunc("/api/publish", publishDefaultMessage) //AMQCmdDefPub
	//http.HandleFunc("/api/publish/unreliable_all", publishDefaultMessage)//AmqCmdDefUnreliable2All
	//http.HandleFunc("/api/publish/unreliable_rand_one", publishDefaultMessage)//AmqCmdDefUnreliable2All
	//http.HandleFunc("/api/publish/reliable_rand_one", publishDefaultMessage)//AmqCmdDefUnreliable2All
	//http.HandleFunc("/api/publish/reliable_spec_one", publishDefaultMessage)//AmqCmdDefUnreliable2All
	err = http.ListenAndServe(httpAddress, nil)
	log.Printf("web server start result :%v", err)
}

var bridgePublishClient *AmqClientPublisher

func startBridgePublishClient(publishAddr string) (err error) {
	bridgePublishClient, err = NewClientPublish(publishAddr)
	return err
}

type generalResponse struct {
	Code    int
	Message string
	Data    interface{}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,POST")

	fmt.Fprintf(w, "hello world")
}

func Sss() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8000", nil)
}

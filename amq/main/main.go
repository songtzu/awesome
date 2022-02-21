package main

import (
	"awesome/mq"
	"fmt"
	"log"
	"os"
	"time"
)

var logger *log.Logger
var file *os.File
var err error

func main()  {
	//file, err = os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 666)
	//if err != nil {
	//	logger.Fatal(err)
	//}
	//logger = log.New(file, "", log.LstdFlags)
	//logger.SetPrefix("Test- ") // 设置日志前缀
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//os.Stdout = file
	//os.Stderr = file
	//defer file.Close()

	fmt.Println("message queue test& dev code")
	mq.NewXmq("127.0.0.1:18888","127.0.0.1:19999")
	for ; ;  {
		time.Sleep(1*time.Hour)
	}

}

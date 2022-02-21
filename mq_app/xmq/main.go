package main

import "awesome/mq_app/xmq/app"

/***********
 * 启动一个xpub，处理真实的publisher的消息
 * 	启动一个xsub, 处理真实的subscriber的消息
 *************/
func main() {
	app := app.NewApp()
	app.Run()
}

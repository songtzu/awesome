package app

import (
	"awesome/app"
	"awesome/alog"
	"awesome/mq"
	"awesome/mq_app/xsub/msghandler"
)

type XsubApp struct {
	app.App
	sub *mq.AMQXSub
	dispatch msghandler.IMessageDispatch
}

func(this *XsubApp) OnInit() {
	alog.Info("XsubApp init ...")
	this.dispatch = &msghandler.MsgHandler{}
	this.sub = mq.NewXSub( "127.0.0.1:19999",
		this.dispatch.OnDispatchMessage)
}

func(this *XsubApp) OnStart() {
	alog.Info("XsubApp start ...")
	this.sub.TopicSubscription([]mq.AMQTopic{1001,1002})
}

func(this *XsubApp) OnStop() {
	alog.Info("XsubApp stop ...")

}


func NewApp() *XsubApp {
	this := &XsubApp{}
	this.Derived = this
	this.Init()
	return this
}
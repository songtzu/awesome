package app

import (
	"awesome/app"
	"awesome/alog"
	"awesome/mq"
)


type XmqApp struct {
	app.App
}

func(this *XmqApp) OnInit() {
	alog.Info("XmqApp init ...")

	mq.NewXmq(xargs.Xmq.pubAddr, xargs.Xmq.subAddr)
}

func(this *XmqApp) OnStart() {
	alog.Info("XmqApp start ...")

}

func(this *XmqApp) OnStop() {
	alog.Info("XmqApp stop ...")
}


func NewApp() *XmqApp {
	this := &XmqApp{}
	this.Derived = this
	this.Args = xargs
	this.Init()
	return this
}
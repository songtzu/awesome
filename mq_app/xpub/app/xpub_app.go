package app

import (
	"awesome/app"
	"awesome/alog"
	"awesome/mq"
)

type XpubApp struct {
	app.App
	pub *mq.AMQXPub
}

func(this *XpubApp) OnStart() {
	alog.Info("XpubApp start ...")

	this.pub.MessagePub(1001,[]byte("hello world, this is topic about 1001"))
}

func(this *XpubApp) OnStop() {
	alog.Info("XpubApp stop ...")
}
func(this *XpubApp) OnInit() {
	alog.Info("XpubApp onInit ...")
}


func NewApp() *XpubApp {
	this := &XpubApp{}
	this.Derived = this
	this.SetStatus(app.SERVER_READY)
	this.DebugPort = 60021
	this.SetAppId(app.Str2AppId("0.0.2.1"))
	this.pub = mq.NewXPub("127.0.0.1:18888" )
	return this
}
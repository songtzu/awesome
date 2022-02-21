package app

import (
	"awesome/app"
	"flag"
)

var (
	xargs *Args = &Args{}
)

type Args struct {
	app.ArgsBase
	Xsub ArgsXSub
}

type ArgsXSub struct {
	Addr string
}

func (this *Args) OnInit() {
	flag.StringVar(&xargs.Xsub.Addr,"addr", "", "Subscriber addr")
}

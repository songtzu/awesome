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
	Xpub ArgsXPub
}

type ArgsXPub struct {
	Addr string
}

func (this *Args) OnInit() {
	flag.StringVar(&xargs.Xpub.Addr,"addr", "", "Publisher addr")
}

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
	Xmq ArgsXmq		`json:"xmq"`
}

type ArgsXmq struct {
	pubAddr string 	`json:"pub_addr"`
	subAddr string	`json:"sub_addr"`
}

func (this *Args) OnInit() {
	flag.StringVar(&xargs.Xmq.pubAddr,"pubAddr", "", "Publisher addr")
	flag.StringVar(&xargs.Xmq.subAddr,"subAddr", "", "Subscriber addr")
}

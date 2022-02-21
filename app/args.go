package app

import (
	"awesome/alog"
	"io/ioutil"
	"hjson-go"
	"encoding/json"
	"flag"
)

type IArgs interface {
	IArgsBase
	OnInit()
}

type IArgsBase interface {
	Init(derived IArgs)
	GetBase() *ArgsBase
}


type ArgsBase struct {
	Common ArgsCommon	`json:"common"`
	Redis ArgsRedis `json:"redis"`
}

type ArgsRedis struct {
	Ip 			string 	`json:"ip"`
	Port 		uint32 `json:"port"`
	Account 	string `json:"account"`
	Password 	string `json:"password"`
}

type ArgsCommon struct {
	Version 	string	`json:"version"`
	AppId 		string 	`json:"app_id"`
	Addr  	 	string	`json:"addr"`
	DebugPort 	uint32 	`json:"debug_port"`// prof 端口

	//CompileTime *string // 编译时间
	//GitHash 	*string // 当前githash
}

func (this *ArgsBase) Init(derived IArgs) {
	parseArgs(derived)

	err := readConf(derived)
	if err != nil {
		panic(err)
	}

	this.requireF()
	alog.Info("config final:", derived)
}

func readConf(derived IArgs) error {
	alog.Info("配置文件地址:", *confPath)

	bs, err := ioutil.ReadFile(*confPath)
	if err != nil {
		alog.Err(err)
		return err
	}

	var mm= map[string]interface{}{}
	err = hjson.Unmarshal(bs, &mm)
	if err != nil {
		alog.Err(err)
		return err
	}


	alog.Info("111111111", mm)

	bsj, err := json.Marshal(mm)
	if err != nil {
		alog.Err(err)
		return err
	}
	alog.Info("2222222", string(bsj), derived)
	err = json.Unmarshal(bsj, derived)
	if err != nil {
		alog.Err(err)
		return err
	}
	alog.Info("读取配置文件成功  文件地址:", *confPath, ", 配置内容:%v",  derived)

	return nil
}

var (
	confPath *string
	version *string
	appId *string
	addr *string
	debugPort *int
)

func parseArgs(derived IArgs) {
	flag.Parse()
	confPath = flag.String("confPath", "app.json", "config path")

	version = flag.String("version", "", "app version")
	appId = flag.String("appId", "", "app id string")
	addr = flag.String("addr", "", "app addr ")
	debugPort = flag.Int("debugPort", 0, "debug port")

	alog.Info("flag parse config path: ", *confPath)
}

func(this *ArgsBase) requireF() {
	if version != nil && *version != "" {
		this.Common.Version = *version
	}

	if appId != nil && *appId != "" {
		this.Common.AppId = *appId
	}

	if addr != nil && *addr != "" {
		this.Common.Addr = *addr
	}

	if debugPort != nil && *debugPort != 0 {
		this.Common.DebugPort = uint32(*debugPort)
	}
}

func (this *ArgsBase) GetBase() *ArgsBase {
	return this
}

var (
	xargs *ArgsBase
)

func SetArgs(args *ArgsBase) {
	xargs = args
}

func GetArgs() *ArgsBase {
	return xargs
}


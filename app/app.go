package app

import (
	"runtime"
	"net/http"
	"fmt"
	"os"
	"os/signal"
	"reflect"
)

type AppInterface interface {
	OnInit()
	OnStart()
	OnStop()
}

const (
	SERVER_READY = iota
	SERVER_START
	SERVER_STOP
)


type App struct {
	appId     APPID
	Derived   AppInterface
	Args      IArgs
	status    int
	signal    chan os.Signal
	DebugPort uint32
}

func (this *App) Init() {
	// parse args
	this.initArgs()

	this.SetAppId(Str2AppId(GetArgs().Common.AppId))

	this.SetStatus(SERVER_READY)
	this.DebugPort = GetArgs().Common.DebugPort

	this.Derived.OnInit()

}

func(this *App) Run() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.GC()

	//注册退出回调
	RegisterProcessExit(func() {
		this.status = SERVER_STOP
		this.Derived.OnStop()
	})

	// prob
	this.runProf(this.DebugPort)

	this.signal = make(chan os.Signal)
	signal.Notify(this.signal, os.Kill, os.Interrupt)

	this.status = SERVER_START
	this.Derived.OnStart()

	for this.status == SERVER_START {
		select {
			case <- this.signal:
				OnProcessExit()
		}
	}
}

func (this *App) runProf(port uint32) {
	if this.DebugPort != 0 {
		go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}
}

func (this *App) GetAppId() APPID {
	return this.appId
}

func (this *App) GetStatus() int {
	return this.status
}

func (this *App) initArgs() {
	if this.Args == nil {
		return
	}

	this.Args.Init(this.Args)
	f := reflect.ValueOf(this.Args).MethodByName("OnInit")
	if f.IsValid() {
		f.Call([]reflect.Value{})
	}

	SetArgs(this.Args.GetBase())
}

func (this *App) SetStatus(status int) {
	this.status = status
}
func (this *App) SetAppId(id APPID) {
	this.appId = id
}
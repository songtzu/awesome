package alog

import (
	"fmt"
	"runtime"
	"strings"

	"awesome/defs"
	"os"
	"path/filepath"
	"time"
)

var (
	pid      = os.Getpid()
	program  = filepath.Base(os.Args[0])
)

func Debug(args ...interface{}) {
	traceOut("D", args...)
}

func Trace(args ...interface{}) {
	traceOut("T", args...)
}
func traceOut(level string, args ...interface{}) {
	file, line := header(  0)
	//str:=fmt.Sprintf("%s,%s.%d,%d,%d,%d,%s:%d,%s",level,
	//	time.Now().Format(defs.GolangTimeBase),time.Now().Nanosecond()/1000,
	//	config.GetConfig().Server.ServerID, config.GetConfig().Server.AppID,
	//	pid, file, line, fmt.Sprint(args...))
	str:=fmt.Sprintf("%s,%s.%d,%d,%d,%d,%s:%d,%s",level,
		time.Now().Format(defs.GolangTimeBase),time.Now().Nanosecond()/1000,
		0, 0,
		pid, file, line, fmt.Sprint(args...))
	fmt.Println(str)
}
func Info(args ...interface{}) {
	traceOut("I", args...)
}
func Err(args ...interface{}) {
	traceOut("E", args...)
}

func Warn(args ...interface{}) {
	traceOut("W", args...)
}

func Falnf(fm string,args ...interface{}) {
	traceOut("F",fmt.Sprintf(fm,args...))
	os.Exit(1)
}

func Faln(args ...interface{}) {
	traceOut("F", args...)
	os.Exit(1)
}

func Errf(fm string,args ...interface{}) {
	traceOut("E",fmt.Sprintf(fm,args...))
}

func Infof(fm string,args ...interface{}) {
	traceOut("I",fmt.Sprintf(fm,args...))
}

/*
 * Level,Time,ServerId,AppId,Pid,File,Line,RoomCode,UserId,Cmd,MESSAGE
 */
func writeOut(level string, args ...interface{}) {
	file, line := header(  1)
	//str:=fmt.Sprintf("%s,%s.%d,%d,%d,%d,%s:%d,%s",level,
	//	time.Now().Format(defs.GolangTimeBase),time.Now().Nanosecond()/1000,
	//		config.GetConfig().Server.ServerID, config.GetConfig().Server.AppID,
	//			pid, file, line, fmt.Sprint(args...))
	str:=fmt.Sprintf("%s,%s.%d,%d,%d,%d,%s:%d,%s",level,
		time.Now().Format(defs.GolangTimeBase),time.Now().Nanosecond()/1000,
		0,0,
		pid, file, line, fmt.Sprint(args...))
	fmt.Println(str)
}


func header( depth int) ( string, int) {
	_, file, line, ok := runtime.Caller(3 + depth)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return file, line
}

func init() {

}


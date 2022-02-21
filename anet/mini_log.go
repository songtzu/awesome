package anet
//
//import (
//	"fmt"
//	"os"
//	"path/filepath"
//	"runtime"
//	"strings"
//	"time"
//)
//
//func logdebug(args ...interface{}) {
//	traceOut("D", args...)
//}
//
//func traceOut(level string, args ...interface{}) {
//	file, line := header(  0)
//	str:=fmt.Sprintf("%s,%s.%d,%d,%d,%d,%s:%d,%s",level,
//		time.Now().Format("2006-01-02 15:04:05"),time.Now().Nanosecond()/1000,
//		-1, -2,
//		pid, file, line, fmt.Sprint(args...))
//	fmt.Println(str)
//}
//func loginfo(args ...interface{}) {
//	traceOut("I", args...)
//}
//func logerr(args ...interface{}) {
//	traceOut("E", args...)
//}
//
//func logwarn(args ...interface{}) {
//	traceOut("W", args...)
//}
//
//func header( depth int) ( string, int) {
//	_, file, line, ok := runtime.Caller(3 + depth)
//	if !ok {
//		file = "???"
//		line = 1
//	} else {
//		slash := strings.LastIndex(file, "/")
//		if slash >= 0 {
//			file = file[slash+1:]
//		}
//	}
//	return file, line
//}
///*
// * Level,Time,ServerId,AppId,Pid,File,Line,RoomCode,UserId,Cmd,MESSAGE
// */
//func writeOut(level string, args ...interface{}) {
//	file, line := header(  3)
//	str:=fmt.Sprintf("%s,%s.%d,%d,%d,%d,%s:%d,%s",level,
//		time.Now().Format("2006-01-02 15:04:05"),time.Now().Nanosecond()/1000,
//		-1, -2,
//		pid, file, line, fmt.Sprint(args...))
//	fmt.Println(str)
//}
//
//var (
//	pid      = os.Getpid()
//	program  = filepath.Base(os.Args[0])
//)
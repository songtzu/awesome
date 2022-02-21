package timer

import (
	"testing"
	"fmt"
	"time"
	"strconv"
)
var count = 0
func timerCb(key string ) {
	count++
	if count>999990{
		fmt.Println(count,"执行回调" + key + "当前时间" +strconv.Itoa(int(getMillisecond())))
	}
	//fmt.Println(count,"执行回调" + key + "当前时间" +strconv.Itoa(int(getMillisecond())))
}
func TestTimer(t *testing.T) {
	fmt.Println(getMillisecond())
	for index:=0;index <1000000;index ++{
		execTime:=int(getMillisecond()) + index/100
		key:=strconv.Itoa(index) +":::"+ strconv.Itoa(int(execTime))
		SetTimeTaskWithCallback(key,int64(index/100),timerCb)
	}

	fmt.Println("当前时间", time.Now().UnixNano())
	time.Sleep(50*time.Second)

}


func TestTimer2(t *testing.T) {
	execTime:=int(getMillisecond() + 1000)
	key:=strconv.Itoa(int(execTime))
	fmt.Println("执行时间" , execTime,"当前实际",getMillisecond())
	SetTimeTaskWithCallback(key,1000,timerCb)

	//fmt.Println("当前时间", time.Now().UnixNano())
	time.Sleep(50*time.Second)

}
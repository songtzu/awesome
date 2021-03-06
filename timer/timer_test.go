package timer

import (
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"
)
var count = 0
func timerCb(key string ) {
	fmt.Println("timerCb", key)
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

func TestTimer4(t *testing.T) {
	execTime:=int(getMillisecond() + 1000)
	key:=strconv.Itoa(int(execTime))
	fmt.Println("执行时间" , execTime,"当前实际",getMillisecond())
	SetTimeTaskWithCallback(key,5000,timerCb)
	DeleteTimer(key)
	time.Sleep(50*time.Second)
}

func TestTimer3(t *testing.T) {
	node1 := Node{executeTime:123}
	node2 := Node{executeTime:13}
	node3 := Node{executeTime:16}
	linkTimeTask.insertWithSort(&node1)
	linkTimeTask.insertWithSort(&node2)
	linkTimeTask.insertWithSort(&node3)
	node4 := Node{executeTime:234}
	linkTimeTask.insertWithSort(&node4)
	tail := linkTimeTask.last
	for tail != nil {
		fmt.Println("tail:", tail)
		tail = tail.prev
	}

	head := linkTimeTask.head
	for head != nil {
		fmt.Println("head:",head)
		head = head.next
	}

	node5 := Node{executeTime:28}
	linkTimeTask.insertWithSort(&node5)
	head = linkTimeTask.head

	for head != nil {
		fmt.Println(head)
		head = head.next
	}

}

func TestTimer5(t *testing.T) {
	log.Printf("启动时的长度，node list长度:%d, timer map 长度:%d",linkTimeTask.length, len(timerMap))

	removeKey:= ""
	for i:=0;i<10;i++ {
		key:=strconv.Itoa(i)
		if i==3{
			removeKey = key
		}
		fmt.Printf("设置序号:%d, 当前实际:%d \n", i , getMillisecond())
		log.Printf("序号：%d, node list长度:%d, timer map 长度:%d",i ,linkTimeTask.length, len(timerMap))

		SetTimeTaskWithCallback(key,5000,timerCb)
	}

	DeleteTimer(removeKey)
	log.Printf("删除key:%s", removeKey)
	time.Sleep(7*time.Second)

	log.Printf("node list长度:%d, timer map 长度:%d , %v",linkTimeTask.length, len(timerMap), timerMap)

	time.Sleep(7*time.Second)

}

func TestTimer6(t *testing.T)  {
	log.Printf("启动时的长度，node list长度:%d, timer map 长度:%d",linkTimeTask.length, len(timerMap))

	for i:=100;i<120;i++ {
		key:=strconv.Itoa(i)

		fmt.Printf("设置序号:%d, 当前实际:%d \n", i , getMillisecond())
		log.Printf("序号：%d, node list长度:%d, timer map 长度:%d",i ,linkTimeTask.length, len(timerMap))

		SetTimeTaskWithCallback(key,5000,timerCb)
	}
	log.Printf("DeleteEventsByPartialMatch之前node list长度:%d, timer map 长度:%d , %v",linkTimeTask.length, len(timerMap), timerMap)

	DeleteEventsByPartialMatch("11")

	log.Printf("DeleteEventsByPartialMatch之后node list长度:%d, timer map 长度:%d , %v",linkTimeTask.length, len(timerMap), timerMap)
	time.Sleep(7*time.Second)
	log.Printf("退出前的node list长度:%d, timer map 长度:%d , %v",linkTimeTask.length, len(timerMap), timerMap)

}

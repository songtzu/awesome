package timer

import (
	"time"
	"sync"
	"fmt"

)

var linkTimeTask = LinkedList{}
var timerMap = make(map[string]*Node) //保存待执行的计时器，方便按链表节点指针地址直接删除定时器
var timerMapMutex sync.RWMutex

type DefTimerCallback = func(key string )

type timerDef struct {
	key string
	createTime int64
	//execTime int64
	cb interface{}
}

func getMillisecond() int64  {
	return time.Now().UnixNano()/1000000
}
/**
 * key:定时器键名
 * delayMillisecond：定时器执行时间,刻度精确到毫秒。理论时间。
 * cb 回调
 */
func SetTimeTaskWithCallback(key_ string, delayMillisecond int64, cb_ DefTimerCallback)  {

	if v,ok:=timerMap[key_];ok{
		//linkTimeTask.erase(v)
		v.removeTaskByKey(key_)
	}
	createTime:=getMillisecond()
	executeTime:=createTime + delayMillisecond
	task := &timerDef{key:key_,createTime:createTime ,cb:cb_}
	if node:=findTaskNode(executeTime);node!=nil{
		//存在该时间单元（毫秒）的node，插入该node的元素（slice）中。
		//fmt.Println("存在该执行单元")
		node.insertTask(task)
		timerMap[key_] = node
	}else {
		//不存在，新建node。
		node:=Node{  executeTime:executeTime}
		node.insertTask(task)
		linkTimeTask.insertWithSort(&node)
		//fmt.Println("不存在node，创建node，有序插入",node.executeTime)
		timerMap[key_] = &node
	}

}

func findTaskNode(executeTime int64) *Node {
	node := linkTimeTask.last
	for ;;{
		if node==nil{
			return nil
		}
		if node.executeTime == executeTime{
			return node
		}
		node = node.prev
	}
	return nil
}

func run() {
	//每一毫秒激活一次，从head开始处理双向链表已到期的任务
	for ; ; {

		for ; ; {
			node := linkTimeTask.GetHead()
			now := getMillisecond()
			if node!=nil{
				//fmt.Println("执行非空任务",node.executeTime,node.executeTime <= now, now)
			}
			if node != nil && node.executeTime <= now {
				//到期的任务
				for _, task := range node.data {
					if cb, ok := task.cb.(DefTimerCallback); ok {
						//fmt.Println("任务到期", task.key)
						cb(task.key)
					}
				}
				linkTimeTask.erase(node)

				//fmt.Println(len(node.data))
				//node = node.next
			} else {
				break
			}
		}
		time.Sleep(1 * time.Millisecond)
	}
	fmt.Println("timewheel执行退出")
}
func init()  {
	go run()
}
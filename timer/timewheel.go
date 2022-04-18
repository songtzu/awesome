package timer

import (
	"log"
	"strings"
	"sync"
	"time"
)

var linkTimeTask = LinkedList{}
var timeLinkRWMutex sync.RWMutex
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
/*SetTimeTaskWithCallback
 * key:定时器键名
 * delayMillisecond：定时器执行时间,刻度精确到毫秒。理论时间。
 * cb 回调
 */
func SetTimeTaskWithCallback(key_ string, delayMillisecond int64, cb_ DefTimerCallback)  {
	DeleteTimer(key_)
	createTime:=getMillisecond()
	executeTime:=createTime + delayMillisecond
	task := &timerDef{key:key_,createTime:createTime ,cb:cb_}
	//fmt.Println("执行时间:", key_, executeTime)
	if node:=findTaskNode(executeTime);node!=nil{
		//存在该时间单元（毫秒）的node，插入该node的元素（slice）中。
		//fmt.Println("存在该执行单元")
		log.Printf("存在该执行单元:%d",executeTime)
		node.insertTask(task)
		timerMapMutex.Lock()
		timerMap[key_] = node
		timerMapMutex.Unlock()
	}else {
		//不存在，新建node。
		log.Printf("不存在此node，新建:%d",executeTime)

		node:=Node{  executeTime:executeTime}
		node.insertTask(task)
		timeLinkRWMutex.Lock()
		linkTimeTask.insertWithSort(&node)
		timeLinkRWMutex.Unlock()
		timerMapMutex.Lock()
		timerMap[key_] = &node
		timerMapMutex.Unlock()
	}
}

func DeleteTimer(key string) {
	timerMapMutex.Lock()
	defer timerMapMutex.Unlock()

	if v,ok:=timerMap[key];ok{
		v.removeTaskByKey(key)
		delete(timerMap,key)
	}

	//deleteTimer(key)
}

func findTaskNode(executeTime int64) *Node {
	timeLinkRWMutex.RLock()
	defer timeLinkRWMutex.RUnlock()

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
			timeLinkRWMutex.Lock()
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
						deleteTimerFromMap(task.key)
					}
				}
				linkTimeTask.erase(node)

			} else {
				timeLinkRWMutex.Unlock()
				break
			}
			timeLinkRWMutex.Unlock()
		}
		time.Sleep(10 * time.Millisecond)
	}
}


func init()  {
	go run()
}

/*DeleteEventsByPartialMatch
 * 部分匹配key查询删除。
 */
func DeleteEventsByPartialMatch(partialKey string) (count int) {
	count = 0
	timerMapMutex.Lock()
	for key,y := range timerMap {
		if strings.Contains(key,partialKey){
			count += 1
			delete(timerMap, key)
			y.removeTaskByKey(key)
		}

	}
	timerMapMutex.Unlock()
	return count
}
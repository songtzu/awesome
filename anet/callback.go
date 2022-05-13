package anet

import (
	"log"
	"math"
	"sync"
	"time"
)

type DefNetIOCallback = func(msg *PackHead)

const minDelayTimeMillisecond = 100

var netIOCallbackMap sync.Map //map:seq---->*netIORegistCallback

func registCallback(head *PackHead, cb DefNetIOCallback) {
	//netIOCallbackMap.Store(head.SequenceID,cb)
	//log.Printf("package:%v+, 注册回调函数:%p", head, cb)
	registCallbackWithinTimeLimit(head, cb, 0, nil)
}
func registCallbackWithinTimeLimit(head *PackHead, cb DefNetIOCallback, delayMillisecond int64, evtChan chan *PackHead) {
	createTime := time.Now().UnixMilli()
	if delayMillisecond <= minDelayTimeMillisecond {
		delayMillisecond = minDelayTimeMillisecond + 1
	}
	register := &netIORegistCallback{cb: cb, createTime: createTime, isTimeout: false,
		deadline: createTime + delayMillisecond, eventChan: evtChan}
	netIOCallbackMap.Store(head.SequenceID, register)
}

type netIORegistCallback struct {
	cb DefNetIOCallback
	//精确到毫秒
	createTime int64
	//精度为毫秒,如果与createTime相等，则为无超时限制。
	deadline int64
	//如果是设置了超时的回调接口，接收到数据的时候，写入此chan
	eventChan chan *PackHead
	isTimeout bool
}

/*
 * return nil if not found.
 */
func popCallback(head *PackHead) (isProcessed bool) {
	log.Printf("popCallback, reservHigh:%d, pack:%d,", head.ReserveHigh, head.SequenceID)
	if v, ok := netIOCallbackMap.Load(head.SequenceID); ok {
		//var cb DefNetIOCallback
		log.Printf("popCallback,111, cmd%d, SequenceID:%d", head.Cmd, head.SequenceID)

		if register, ok := v.(*netIORegistCallback); ok {
			log.Printf("popCallback,222, cmd%d", head.Cmd)

			netIOCallbackMap.Delete(head.SequenceID)
			if register.cb != nil {
				//log.Printf("popCallback,333, cmd%d, cb:%v", head.Cmd, reflect.TypeOf(register.cb))
				register.cb(head)
				return true
			}
			log.Printf("popCallback,444, cmd%d", head.Cmd)

			if register.eventChan != nil {
				currentTime := time.Now().UnixMilli()
				if register.deadline >= currentTime {
					//没超时的任务
					//logdebug("设置超时时间的任务，正常返回")
					tmp := make([]byte, len(head.Body))
					//log.Println("popCallback===>", string(head.Body))
					copy(tmp, head.Body)
					head.Body = tmp
					//log.Println("popCallback===2222===>", string(head.Body))
					register.eventChan <- head
					return true
				} else {
					//超时任务
					log.Println("超时任务,当前时间", currentTime, "设置的超时时间：", register.deadline, "创建时间", register.createTime)
				}
			}
		} else {
			log.Println("type convert error for net callback")
		}
		return true
	} else {
		log.Println("popCallback not ok", head.SequenceID, string(head.Body))
	}

	return false
}

const startIndexForSequenceId = 1000000

var autoIncreaseSequenceId uint32 = startIndexForSequenceId
var autoIncreaseSequenceIdLocker = new(sync.Mutex)

func allocateNewSequenceId() uint32 {
	autoIncreaseSequenceIdLocker.Lock()
	if autoIncreaseSequenceId > math.MaxUint32-1 {
		autoIncreaseSequenceId = startIndexForSequenceId
	}
	autoIncreaseSequenceId++
	defer autoIncreaseSequenceIdLocker.Unlock()
	return autoIncreaseSequenceId
}

func init() {
	go taskCheckTimeout()
}

func taskCheckTimeout() {
	for {
		netIOCallbackMap.Range(func(key, value interface{}) bool {
			item := value.(*netIORegistCallback)
			if item.deadline < time.Now().UnixMilli() {
				item.isTimeout = true
				item.eventChan <- nil
				netIOCallbackMap.Delete(key)
				log.Println("检测到超时了", key)
				return false
			}
			return true
		})
		time.Sleep(minDelayTimeMillisecond * time.Millisecond)
	}
}

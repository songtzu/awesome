package anet

import (
	"log"
	"math"
	"sync"
	"time"
)

type DefNetIOCallback = func(msg *PackHead)

var netIOCallbackMap sync.Map

func registCallback(head *PackHead, cb DefNetIOCallback) {
	//netIOCallbackMap.Store(head.SequenceID,cb)
	//log.Printf("package:%v+, 注册回调函数:%p", head, cb)
	registCallbackWithinTimeLimit(head, cb, 0, nil)
}
func registCallbackWithinTimeLimit(head *PackHead, cb DefNetIOCallback, delayMillisecond int64, evtChan chan *PackHead) {
	createTime := time.Now().UnixMilli()
	regist := &netIORegistCallback{cb: cb, createTime: createTime, deadline: createTime + delayMillisecond, eventChan: evtChan}
	netIOCallbackMap.Store(head.SequenceID, regist)
}

type netIORegistCallback struct {
	cb DefNetIOCallback
	//精确到毫秒
	createTime int64
	//精度为毫秒,如果与createTime相等，则为无超时限制。
	deadline int64
	//如果是设置了超时的回调接口，接收到数据的时候，写入此chan
	eventChan chan *PackHead
}

/*
 * return nil if not found.
 */
func popCallback(head *PackHead) DefNetIOCallback {
	log.Printf("popCallback, reservHigh:%d, pack:%v,", head.ReserveHigh, head)
	if v, ok := netIOCallbackMap.Load(head.SequenceID); ok {
		//var cb DefNetIOCallback
		log.Printf("popCallback,111, cmd%d, SequenceID:%d", head.Cmd, head.SequenceID)

		if regist, ok := v.(*netIORegistCallback); ok {
			log.Printf("popCallback,222, cmd%d", head.Cmd)

			netIOCallbackMap.Delete(head.SequenceID)
			if regist.cb != nil {
				log.Printf("popCallback,333, cmd%d, cb:%v", head.Cmd, regist.cb)
				return regist.cb
			}
			log.Printf("popCallback,444, cmd%d", head.Cmd)

			if regist.eventChan != nil {
				currentTime := time.Now().UnixNano() / int64(time.Millisecond)
				if regist.deadline >= currentTime {
					//没超时的任务
					//logdebug("设置超时时间的任务，正常返回")
					tmp := make([]byte, len(head.Body))
					log.Println("popCallback===>", string(head.Body))
					copy(tmp, head.Body)
					head.Body = tmp
					log.Println("popCallback===2222===>", string(head.Body))
					regist.eventChan <- head
				} else {
					//超时任务
					log.Println("超时任务,当前时间", currentTime, "设置的超时时间：", regist.deadline, "创建时间", regist.createTime)
				}
			}
		} else {
			log.Println("type convert error for net callback")
		}

	} else {
		log.Println("popCallback not ok", head.SequenceID, string(head.Body))
	}

	//netIOCallbackMap.Range(func(key, value interface{}) bool {
	//	fmt.Println(key,value)
	//	return true
	//})

	return nil
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

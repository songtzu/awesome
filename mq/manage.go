package mq

import (
	"container/heap"
	"container/list"
	"fmt"
	"sync"
	"time"
)

var (
	defaultTimeout = 3 * time.Second
)

type amqMessageModelManage struct {
	freeList     *list.List                 // 等待处理的消息
	usingMap     map[IMessage]*messageData  // 正在处理的消息
	itemStateMap map[IMessage]*list.Element // 用于查询元素在链表中的位置
	timeoutHeap  *minHeap                   // 堆由于处理超时的消息，消息
	mutex        sync.RWMutex

	maxSize     int // 最大消息数量
	maxConsumer int // 最大等待消费者

	consumer        int           // 当前等待消息的消费者
	consumerBuf     chan IMessage // 等待消息的消费者从此channel中取数据
	timeoutDuration time.Duration
	
	once sync.Once
}

type messState struct {
	// state 有5种状态
	//   1.状态0表示准备中
	//   2.状态1表示因队列满而标记删除（heap中lazy删除）
	//   3.状态2表示发送中
	//   4.状态3表示超时，此状态只是一个临时状态，此数据会被马上删除
	//   5.状态4表示已经被成功处理，运行检查超时时清理掉（heap中lazy删除）
	//
	state int16

	// 存活到的时间
	alive time.Time
}

func (m *messState) Less(dst *messageData) bool {
	return m.alive.Before(dst.alive)
}

type messageData struct {
	IMessage

	messState
}

const (
	state_ready = iota
	state_del
	state_sending
	state_timeout
	state_ok
)

type minHeap struct {
	list []*messageData
}

type IMessage interface {
	// 超时处理
	OnTimeOut()
	OnFillRemove()
}

type IMessageStateData interface {
	GetIMessage() IMessage
}

func (m *minHeap) Len() int {
	return len(m.list)
}

func (m *minHeap) Less(i, j int) bool {
	if m.list[i].alive.IsZero() {
		return false
	}
	return m.list[i].Less(m.list[j])
}

func (m *minHeap) Swap(i, j int) {
	m.list[i], m.list[j] = m.list[j], m.list[i]
}

func (m *minHeap) Push(x interface{}) {
	m.list = append(m.list, x.(*messageData))
}

func (m *minHeap) Pop() interface{} {
	n := len(m.list)
	x := m.list[n-1]
	m.list = m.list[0 : n-1]
	return x
}

func (a *amqMessageModelManage) timeoutLoop() {
	var timer = time.NewTicker(1 * time.Millisecond)
	defer timer.Stop()

	for {
		select {
		case now := <-timer.C:
			{

				a.mutex.Lock()
				for a.timeoutHeap.Len() > 0 {
					var min = a.timeoutHeap.list[0]
					if min.state == state_ok || min.state == state_del {
						heap.Pop(a.timeoutHeap)
						continue
					}

					// 无超时
					if min.alive.IsZero() {
						break
					}
					if min.alive.Before(now) {
						heap.Pop(a.timeoutHeap)

						if min.state == state_ready {
							//去掉freeList
							v, _ := a.itemStateMap[min.IMessage]

							a.freeList.Remove(v)
							delete(a.itemStateMap, min.IMessage)

						} else {
							delete(a.usingMap, min.IMessage)
						}

						min.state = state_timeout

						go min.OnTimeOut()

						continue
					}

					break

				}
				a.mutex.Unlock()
			}
		}
	}
}

func (a *amqMessageModelManage) usingLock(mess *messageData) {
	a.usingMap[mess.IMessage] = mess
}

func (a *amqMessageModelManage) deleteFirstLock() {
	first := a.freeList.Front()
	a.freeList.Remove(first)

	ms := first.Value.(*messageData)

	delete(a.itemStateMap, ms.IMessage)
	ms.state = state_del

	go ms.IMessage.OnFillRemove()
}

func (a *amqMessageModelManage) pushFreeLock(data *messageData) {
	item := a.freeList.PushBack(data)
	a.itemStateMap[data.IMessage] = item
}

// 添加消息(默认超时时间)
func (a *amqMessageModelManage) Add(mess IMessage) {
	a.AddWithTimeout(mess, a.timeoutDuration)
}

// 添加无超时的消息
func (a *amqMessageModelManage) AddNoTimeOut(mess IMessage) {
	a.AddWithTimeout(mess, 0)
}

// 添加超时消息，超时时间<=0时添加为无超时消息
func (a *amqMessageModelManage) AddWithTimeout(mess IMessage, timeoutDur time.Duration) {
	var alive time.Time
	if timeoutDur > 0 {
		alive = time.Now().Add(timeoutDur)
	}
	var data = &messageData{
		IMessage: mess,
		messState: messState{
			state: 0,
			alive: alive,
		},
	}

	// 添加节点
	a.mutex.Lock()
	defer a.mutex.Unlock()

	heap.Push(a.timeoutHeap, data)
	if a.consumer > 0 {
		data.state = state_sending
		a.consumer--
		a.usingLock(data)
		a.consumerBuf <- mess

	} else {
		if a.freeList.Len() >= a.maxSize {
			a.deleteFirstLock()

		}
		a.pushFreeLock(data)
	}
}

// 消息处理完成，清理掉
func (a *amqMessageModelManage) Success(mess IMessage) bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	data, ok := a.usingMap[mess]
	if !ok {
		return false
	}

	delete(a.usingMap, mess)

	data.state = state_ok
	return true
}

// 获取一个消息,没有消息时阻塞等待
func (a *amqMessageModelManage) Get() IMessage {
	a.mutex.Lock()

	//a.messageList.Front()
	item := a.freeList.Front()

	if item == nil {
		a.consumer++
		a.mutex.Unlock()

		return <-a.consumerBuf
	}

	var d = item.Value.(*messageData)
	d.state = state_sending
	a.usingMap[d.IMessage] = d
	a.freeList.Remove(item)
	delete(a.itemStateMap, d.IMessage)

	a.mutex.Unlock()

	return d.IMessage
}

func (a *amqMessageModelManage) SetMaxSize(maxSize int) {
	a.mutex.Lock()
	a.maxSize = maxSize
	a.mutex.Unlock()
}

func (a *amqMessageModelManage) Start() {
	a.once.Do(func() {
		go a.timeoutLoop()
	})
}

func (a *amqMessageModelManage) FreeLen() int {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	var count = len(a.itemStateMap)
	if a.freeList.Len() != count {
		panic(fmt.Sprintf("length not vaild  %d %d", count, a.freeList.Len()))
	}

	return count
}

func (a *amqMessageModelManage) UsingLen() int {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	return len(a.usingMap)

}

func (a *amqMessageModelManage) HeapLen() int {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	return a.timeoutHeap.Len()
}

func NewAmqMessageModelManage(maxSize int, maxSummer int, timeoutDuration time.Duration) *amqMessageModelManage {
	if timeoutDuration <= 0 {
		timeoutDuration = defaultTimeout
	}

	return &amqMessageModelManage{
		freeList:        list.New(),
		timeoutHeap:     &minHeap{list: make([]*messageData, 0, 1024)},
		itemStateMap:    make(map[IMessage]*list.Element),
		usingMap:        make(map[IMessage]*messageData),
		maxSize:         maxSize,
		maxConsumer:     maxSummer,
		consumerBuf:     make(chan IMessage, maxSummer),
		timeoutDuration: timeoutDuration,
	}
}

var AmqMessageManage = NewAmqMessageModelManage(2048, 20, 3*time.Second)

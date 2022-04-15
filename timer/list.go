package timer

import (
	"log"
	"strings"
)


type LinkedList struct {
	head   *Node
	last   *Node
	length uint
}

func NewLinkedList() *LinkedList {
	var list *LinkedList = new(LinkedList)
	list.head = nil
	list.last = nil
	list.length = 0
	return list
}

/**
 * 获取表头
 */
func (this LinkedList) GetHead() *Node {
	return this.head
}

/**
 * 获取表尾
 */
func (this LinkedList) GetLast() *Node {
	return this.last
}

func (this LinkedList) Length() uint {
	return this.length
}
/*
 * 有序插入
 * 		根据executeTime排序插入
 */
func (this *LinkedList) insertWithSort(node *Node)  {
	//fmt.Println("有序插入")
	tail := this.last
	for ;;{
		if tail!=nil{
			if tail.executeTime<=node.executeTime{
				//从尾部开始查找最后的执行节点
				node.prev = tail
				node.next = tail.next
				tail.next = node
				if tail==this.last{
					//尾部追加，跟新尾节点
					this.last = node
					//fmt.Println("尾部追加新节点")
				}
				return
			}else{
				//不是该节点，继续遍历前节点
				tail = tail.prev
			}
		}else {
			//首节点或者空表
			if this.head!=nil && this.head.next!=nil{
				//插入位置为首节点
				node.next = this.head.next
				//fmt.Println("非空首节点")
			}
			//else {
				//插入空表
				//fmt.Println("空表插入",node.executeTime)
			//}
			this.head = node
			this.head.prev = nil
			this.last = node

			return
		}
	}
}
func (this *LinkedList) pushBack(node Node) *Node {
	node.next = nil
	if nil == this.head { //空表
		this.head = &node
		this.head.prev = nil
		this.last = this.head
	} else {
		node.prev = this.last
		this.last.next = &node
		this.last = this.last.next
	}
	log.Printf("insert %d %v\n", this.length, this.last.Data())
	this.length++
	return this.last
}

func (this *LinkedList) erase(node *Node) (isOk bool){
	if nil == node {
		return false
	}

	if node == this.head && node == this.last {
		this.head = nil
		this.last = nil
		this.length = 0
		//fmt.Println("仅有的节点删除")
	} else {
		//fmt.Println("还有剩余节点")
		if node == this.head {
			this.head = this.head.next
			if nil != this.head {
				this.head.prev = nil
			}
		} else if node == this.last {
			node.prev.next = nil
			this.last = node.prev
		} else {
			node.prev.next = node.next
			node.next.prev = node.prev
		}
	}
	this.length--
	return true
}

func deleteTimer(key string) {
	deleteTimerFromMap(key).remove()
}
func deleteTimerFromMap(key string) *Node{
	key = strings.TrimSpace(key)

	timerMapMutex.Lock()
	node:= timerMap[key]
	delete(timerMap,key)
	timerMapMutex.Unlock()
	return node
}

func deleteRoomAllTimer(roomCode string) map[string]*Node {
	var nodes = make(map[string]*Node,16)
	roomCode = strings.TrimSpace(roomCode)
	var lenInviteCode = len(roomCode)
	timerMapMutex.Lock()
	for key,y := range timerMap {
		if len(key) >= lenInviteCode && key[:lenInviteCode] == roomCode {
			nodes[key] = y
			delete(timerMap,key)
		}
	}
	timerMapMutex.Unlock()
	return nodes
}

/*
func upsert(key string,node *Node) {
	timerMapMutex.Lock()
	nodeT := timerMap[key]
	timerMap[key] = node
	timerMapMutex.Unlock()
	nodeT.remove()
	if nodeT != nil {
		if t,ok  := nodeT.Data().(timerDef);ok {
			EffectiveManager.DeleteRoomId(key[:6],t.Id)
		}
	}
}
*/

func (n *Node)_remove() bool {
	if nil == n {
		return false
	} else if nil == n.prev { //该元素处于表头，不删除，默认表头不存元素
		return false
	} else if nil == n.next { //该元素处于表尾
		n.prev.next = nil
		n.prev = nil
	} else {
		n.next.prev = n.prev
		n.prev.next = n.next
		n.prev = nil
		n.next = nil
	}
	return true
}
/*
func (n *Node) removeLock() bool{
	rwmutex.Lock()
	ok := n._remove()
	rwmutex.Unlock()
	return ok
}
*/
func (n *Node)remove() bool {
	return n._remove()
}

func (this *Node) insertHead(node Node) *Node { //从表头插入
	if nil == this || nil != this.prev { //为空，或者不是表头(表头的prev为空)
		return nil
	} else {
		if nil != this.next {
			this.next.prev = &node
			node.next = this.next
		}
		this.next = &node
		node.prev = this
	}
	return &node
}

func (this *Node) Next() (node *Node) {
	return this.next
}

func (this *Node) Prev() (node *Node) {
	return this.prev
}

func (this *Node) Data() (data interface{}) {
	return this.data
}

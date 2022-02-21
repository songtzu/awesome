package timer



//定义节点
type Node struct {
	//data interface{}
	data []*timerDef
	executeTime int64
	prev *Node
	next *Node
}

func (this *Node) insertTask(task *timerDef)  {
	if this.data == nil{
		//fmt.Println("插入空")
		this.data = []*timerDef{}
		this.data = append(this.data, task)
	}else {
		this.data = append(this.data, task)
	}
	//fmt.Println(len(this.data))
}

//从slice中删除节点。
func (this *Node) removeTaskByKey(taskKey string) (hasDel bool) {
	if this.data!=nil{
		for index,item:=range this.data{
			if item.key==taskKey{
				this.data=append(this.data[:index], this.data[index+1:]...)
				return true
			}
		}
	}
	return false

}
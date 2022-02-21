package test

import (
	"reflect"
	"fmt"
	"testing"
	"runtime/debug"

	"awesome/alog"
)


type Str struct {
	Field string
}
type InterfaceServer interface {
	ISet(f string)
	IPrint()
}
func (s *Str) ISet(f string) {
	fmt.Println("s==nil",s==nil)
	s.Field = f
}
func (s *Str) IPrint() {
	fmt.Println("IPrint,s==nil", s==nil)
	fmt.Println("====", s.Field)
}
func call(i InterfaceServer)  {
	defer func() {
		if err := recover();err != nil {
			alog.Err("panic .",err,string(debug.Stack()))
		}
	}()
	fmt.Println("~~~~~~~~~~~~~", reflect.TypeOf(i).Elem())
	fmt.Println("---------222--------", reflect.TypeOf(reflect.ValueOf(i).Interface()))
	fmt.Println("-----------------", reflect.TypeOf(i))

	fmt.Println("++++++++++++",reflect.TypeOf( reflect.New(reflect.TypeOf(i)).Interface() ))
	//n:=reflect.New(reflect.TypeOf(i)).Interface()

	p:=reflect.TypeOf(i)
	pp:=reflect.New(p)
	elm,oook:=pp.Elem().Interface().(InterfaceServer)
	s,ook:=elm.(*Str)
	if ook{
		fmt.Println("直接访问",s==nil, ook)
		//fmt.Println("变量",s.Field)
	}else {
		fmt.Println("转结构体失败")
	}

	fmt.Println("复制接口的新实例",reflect.TypeOf(elm), elm==nil,oook)
	//turn2p:= reflect.ValueOf(pp).Interface()
	//fmt.Println("+++++++++++++turn2p+++++++++",reflect.TypeOf(pp),reflect.TypeOf(turn2p))
	//fmt.Println("++++++++++++",reflect.TypeOf(  pp ))


	//cp,ok:=reflect.New(reflect.TypeOf(i)).Interface().(InterfaceServer)
	//fmt.Println("====",ok)
	//p:=(cp).(InterfaceServer)
	//fmt.Println(reflect.TypeOf(cp))
	i.ISet("i instance")
	j:=i
	j.ISet("jjj sample")
	fmt.Println("=========test if elm is nil:", elm==nil)
	elm.IPrint()
	elm.ISet("cp instance")
	j.IPrint()
	i.IPrint()
	elm.IPrint()
}
//func call2(i InterfaceServer)  {
//	typeOfStruc:= reflect.TypeOf(i).Elem()
//	originObj:=i.(*typeOfStruc)
//}



func TestRf(t *testing.T) {


	s := &Str{Field: "astr"}
	call(s)
	a := interface{}(s)

	v := reflect.Indirect(reflect.ValueOf(a))
	b := v.Interface()

	s.Field = "changed Field"

	fmt.Println(a, b)
}

package test

import (
	"testing"
	"fmt"
	"reflect"
)

type IUser interface {
	Name() string
	SetName(name string)
}

type Admin struct {
	name string
}

func (a *Admin) Name() string {
	return a.name
}

func (a *Admin) SetName(name string) {
	a.name = name
}

func TestUser(t *testing.T)  {
	var user1 IUser
	user1 = &Admin{name:"user1"}
	var user2 IUser
	padmin := user1.(*Admin) // Obtain *Admin pointer
	admin2 := *padmin        // Make a copy of the Admin struct
	user2 = &admin2          // Wrap its address in another IUser
	user2.SetName("user2")
	fmt.Printf("User2's name: %s\n", user2.Name()) // The name will be changed as "user2"
	fmt.Printf("User1's name: %s\n", user1.Name())  // The name will be changed as "user2" too, How to make the user1 name does not change？
}

func passType(t reflect.Type) IUser {
	fmt.Println(t)
	v:=reflect.New(t).Interface().(IUser)
	v.SetName("jack")
	fmt.Println(reflect.TypeOf(v))
	fmt.Println(v.Name())
	return v
}


func passInterfaceType(t IUser) IUser {
	k:=reflect.TypeOf(t).Elem()
	u:=passType(k)
	return u
	//fmt.Println(t)
	//v:=reflect.New(t).Interface().(IUser)
	//v.SetName("jack")
	//fmt.Println(v.Name())
}

func passInterfaceType2ReturnInterface(t IUser) IUser {
	k:=reflect.TypeOf(t).Elem()

	v:=reflect.New(k).Interface().(IUser)
	v.SetName("jack")
	fmt.Println(reflect.TypeOf(v))
	fmt.Println(v.Name())
	return v

	//fmt.Println(t)
	//v:=reflect.New(t).Interface().(IUser)
	//v.SetName("jack")
	//fmt.Println(v.Name())
}


func TestReflect(t *testing.T) {
	u1:=passType(reflect.TypeOf(Admin{}))
	u2:=passInterfaceType(&Admin{})
	u3:=passInterfaceType2ReturnInterface(&Admin{})
	fmt.Println("获取到接口实例",u1,u2)
	u1.SetName("u1")
	u2.SetName("u2")
	u3.SetName("u3")
	fmt.Println(u1.Name())
	fmt.Println(u2.Name())
	fmt.Println(u3.Name())
}


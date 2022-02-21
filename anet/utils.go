package anet

import (
	"log"
	"reflect"
)

func forkNewInstanceOfInterface(t InterfaceNet) InterfaceNet {
	k:=reflect.TypeOf(t).Elem()

	v,ok:=reflect.New(k).Interface().(InterfaceNet)
	if ok{
		return v
	}else {
		log.Println("error when fork a instance of interface")
		return nil
	}

}


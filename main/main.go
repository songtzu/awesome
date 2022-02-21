package main

import (
	"awesome/framework"
	"logic"
	"time"
)

func main()  {
	var instance = &logic.AwesomeImplement{}
	go framework.StartSvr(instance)
	go framework.StartHallSession()
	for{
		time.Sleep(1*time.Hour)
	}
}

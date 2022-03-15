package main

import (
	"awesome/framework"
	"log"
	"logic"
	"os"
	"time"
)

func main()  {
	var instance = &logic.AwesomeImplement{}
	if err:=framework.InitDatabase();err!=nil{
		log.Fatalln("exit due to db error",err)
		os.Exit(-5)
	}
	go framework.StartSvr(instance)

	instance.OnInit()

	for{
		time.Sleep(1*time.Hour)
	}
}

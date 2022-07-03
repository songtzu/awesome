package db

import (
	"log"
	"testing"
	"time"
)

type sample struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Pwd  string `json:"pwd"`
}

func TestRedisSet(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := NewRedisPool("redis://:foobared@127.0.0.1:6379/1")
	log.Println("redis创建的结果",err)
	var s2 = &sample{Name: "rose", Age: 20, Pwd: "pwd from rose"}
	arr := make([]sample, 0)
	arr = append(arr, *s2)
	arr = append(arr, *s2)


	log.Println("====", err, s2)

	err = RedisKeySetStr("test_hm:1", "test", 100*time.Second )
	log.Printf("RedisSetKeyStr保存结果:%v",err)
	err = RedisKeySetObj("test_hm:2", s2, 100*time.Second )
	log.Printf("RedisSetKeyObj保存结果:%v",err)

	txt,err := RedisKeyGet("test_hm"  ).Text()
	log.Println("读取记录",txt,err)

}
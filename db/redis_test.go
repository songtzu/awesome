package db

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis/redis"
	"log"
	"reflect"
	"strings"
	"testing"
)

type sample struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Pwd  string `json:"pwd"`
}

func TestRedisSetBytes(t *testing.T) {

	NewRedisPool("redis://foo@127.0.0.1:6379/1")

	_, err := GetRedisInstance().Do("SET", "go_key", "redigo")
	if err != nil {
		fmt.Println("err while setting:", err)
	}

	var s = &sample{Age: 1, Name: "jack", Pwd: "this is pwd"}
	arr, err := json.Marshal(s)
	if err != nil {
		fmt.Println("json to bytes error", err)
	}
	fmt.Println(string(arr))
	_, err = GetRedisInstance().Do("HMSET", "test_hm", "test", string(arr))
	if err != nil {
		fmt.Println("hmset err", err)
	}
	var s2 = sample{}
	r, err2 := redis.Values(GetRedisInstance().Do("HMGET", "test_hm", "test"))
	if err2 != nil {
		fmt.Println("hmset err", err)
	}
	redis.ScanStruct(r, s2)
	fmt.Println(s2)

}

func TestRedisSetStruct(t *testing.T) {

	NewRedisPool("redis://foo@127.0.0.1:6379/1")

	_, err := GetRedisInstance().Do("SET", "go_key", "redigo")
	if err != nil {
		fmt.Println("err while setting:", err)
	}

	var s = sample{Age: 1, Name: "jack", Pwd: "this is pwd"}

	_, err = GetRedisInstance().Do("HMSET", "test_hm", "test", s)
	if err != nil {
		fmt.Println("hmset err", err)
	}
	var s2 = sample{}
	r, err2 := redis.Values(GetRedisInstance().Do("HMGET", "test_hm", "test"))
	if err2 != nil {
		fmt.Println("hmset err", err)
	}
	redis.ScanStruct(r, s2)
	fmt.Println(s2)

}

func TestRedisGET(t *testing.T) {

	NewRedisPool("redis://foo@127.0.0.1:6379/1")

	var s2 = &sample{}
	//dr,e:=GetRedisInstance().Do("HMGET","test_hm","test")
	//fmt.Println(dr)
	//json.Unmarshal((dr.([]byte)),s2)
	//fmt.Println(s2)
	//bytes,err :=redis.Bytes(GetRedisInstance().Do("HMGET","test_hm","test"))
	//if err!=nil{
	//	fmt.Println("结构错误",bytes,err)
	//}
	//json.Unmarshal(bytes,s2)
	r, err2 := redis.Values(GetRedisInstance().Do("HMGET", "test_hm", "test"))
	if err2 != nil {
		fmt.Println("hmset err", err2)
	}
	e3 := redis.ScanStruct(r, s2)
	fmt.Println("扫描", s2, r, e3)

}

func TestRedisGETString(t *testing.T) {

	NewRedisPool("redis://foo@127.0.0.1:6379/1")

	var s2 = &sample{}
	dr, e := GetRedisInstance().Do("HMGET", "test_hm", "test")
	fmt.Println(dr, reflect.TypeOf(dr))
	if e != nil {
		fmt.Println("get错误", e)
	}
	//c.Send("HMGET", "test_hm","test")
	//c.Flush()
	//var vb []byte
	vb, err := redis.ByteSlices(dr, e)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	fmt.Println("vb===>", vb)

	//bytes,err :=redis.Bytes(dr,e)
	//if err!=nil{
	//	fmt.Println("结构错误",bytes,err)
	//}
	err = json.Unmarshal(vb[0], s2)
	if err != nil {
		fmt.Println("结构错误", err)
	} else {
		fmt.Println(s2.Pwd, "name:", s2.Name, s2.Age)
	}

}

func TestRedisGETUtil(t *testing.T) {
	NewRedisPool("redis://foo@127.0.0.1:6379/1")
	var s2 = &sample{}

	err := RedisHMGet("test_hm", "test", s2)

	log.Println("====", err, s2)

}

func TestRedisSETUtil(t *testing.T) {
	NewRedisPool("redis://foo@127.0.0.1:6379/1")
	var s2 = &sample{Name: "rose", Age: 20, Pwd: "pwd from rose"}
	arr := make([]sample, 0)
	arr = append(arr, *s2)
	arr = append(arr, *s2)

	err := RedisHMSet("test_hm", "test5", arr)

	log.Println("====", err, s2)

}

func TestRedisOther(t *testing.T) {
	NewRedisPool("redis://foo@127.0.0.1:6379/1")

	var p1, p2 struct {
		Title  string `redis:"title"`
		Author string `redis:"author"`
		Body   string `redis:"body"`
	}

	p1.Title = "Example"
	p1.Author = "Gary"
	p1.Body = "Hello"

	if _, err := redisPool.Get().Do("HMSET", redis.Args{}.Add("id1").AddFlat(&p1)...); err != nil {
		fmt.Println(err)
		return
	}

	m := map[string]string{
		"title":  "Example2",
		"author": "Steve",
		"body":   "Map",
	}

	if _, err := redisPool.Get().Do("HMSET", redis.Args{}.Add("id2").AddFlat(m)...); err != nil {
		fmt.Println(err)
		return
	}

	for _, id := range []string{"id1", "id2"} {

		v, err := redis.Values(redisPool.Get().Do("HGETALL", id))
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := redis.ScanStruct(v, &p2); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("====》%+v\n", p2)
	}
}

func TestRedisGETALL(t *testing.T) {
	NewRedisPool("redis://foo@127.0.0.1:6379/2")
	var s1 = &sample{Name: "jack", Age: 20, Pwd: "pwd from rose"}
	var s2 = &sample{Name: "rose", Age: 20, Pwd: "pwd from rose"}
	RedisHMSet("test", s1.Name, s1)

	RedisHMSet("test", s2.Name, s2)

	arr, _ := redis.Values(redisPool.Get().Do("HVALS", "test"))
	arrStr, _ := redis.Strings(redisPool.Get().Do("HVALS", "test"))
	log.Println(strings.Join(arrStr, ","))
	log.Println("arrStr===>", arrStr, len(arrStr), arrStr[0])
	log.Println("Arr===>", arr, len(arr))
	for idx := 0; idx < len(arr); idx++ {
		v := &sample{}
		log.Println(idx, "---->", arr[idx])
		itemArr, err := redis.Bytes(arr[idx], nil)
		err = json.Unmarshal(itemArr, v)
		if err != nil {
			log.Println("redis hmget unmarshal error", err, v)
		}
		log.Println(v)
	}

}
func createSlice(p interface{}, bin []byte) interface{} {

	elemType := reflect.TypeOf(p)
	newInstance := reflect.New(elemType.Elem())

	log.Println(elemType, bin)
	elemSlice := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(elemType)), 0, 10)

	log.Println(newInstance)
	err := json.Unmarshal(bin, &newInstance)
	log.Println("err===>", err, newInstance)
	elemSlice = reflect.Append(elemSlice, newInstance)
	log.Println(json.Marshal(elemSlice))
	return elemSlice
}
func TestSlice(t *testing.T) {
	var s2 = &sample{Name: "rose", Age: 20, Pwd: "pwd from rose"}
	bin, _ := json.Marshal(s2)
	log.Println("bin-->", string(bin))
	s := createSlice(sample{}, bin)
	log.Println("++++", s)
}

type Test struct {
	Name string `json:"name,omitempty"`
}

func create(a interface{}) {
	b := []byte(`[{"name": "go"}]`)

	err := json.Unmarshal(b, &a)
	fmt.Println(err, a)
}

func TestDynamicJson(t *testing.T) {
	l := []Test{}
	create(&l)

}

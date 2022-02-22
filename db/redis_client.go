package db

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis/redis"
	"log"
	"strings"
	"time"
)

var redisPool *redis.Pool

/*
 * "redis://<user>:<pass>@localhost:6379/<db>"
 * example:"redis://<pass>@localhost:6379/<db>"
 */
func NewRedisPool(server string) {
	redisPool = &redis.Pool{
		MaxIdle:     30,
		MaxActive:   1000,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	return
}

func GetRedisInstance() (con redis.Conn) {
	return redisPool.Get()
}

/**
 * 只有存在内容才会返回nil
 *		如果不存在此字段或者是redis异常，均返回错误。
 ***/
func RedisHMGet(key, field string, v interface{}) error {
	con := redisPool.Get()
	if con == nil {
		return errors.New("invalid redis con fd")
	}
	arr, err := redis.ByteSlices(con.Do("HMGET", key, field))
	if err != nil {
		log.Println("redis hmget error", err.Error(), key, field)
	}

	if len(arr) == 0 || len(arr[0]) == 0 {
		//log.Println("hmget返回为空")
		return errors.New("不存在此字段内容" + key + field)
	}

	err = json.Unmarshal(arr[0], v)
	if err != nil {
		log.Println("redis hmget unmarshal error", err, key, field, v, len(arr[0]))
	}
	con.Close()
	return err

}

func RedisHMSet(key, field string, v interface{}) error {
	arr, err := json.Marshal(v)
	if err != nil {
		log.Println("json to bytes error", err)
		return errors.New("json序列化错误" + err.Error())
	}
	con := redisPool.Get()
	if con == nil {
		return errors.New("invalid redis con fd")
	}
	_, err = con.Do("HMSET", key, field, string(arr))

	if err != nil {
		log.Println("redis hmset marshal error", err, key, field, v)
	}
	con.Close()
	return err

}

func RedisHMGetAllByKey(key string, v interface{}) error {
	//vals,_:=redisPool.Get().Do("HVALS", key)
	//redis.ScanStruct(vals,v)
	con := redisPool.Get()
	if con == nil {
		return errors.New("invalid redis con fd")
	}
	arr, err := redis.Strings(con.Do("HVALS", key))
	if err != nil {
		log.Println("redis hmget error", err.Error(), key)
	}
	rowStr := "[" + strings.Join(arr, ",") + "]"
	if len(arr) == 0 || len(arr[0]) == 0 {
		//log.Println("hmget返回为空")
		return errors.New("不存在此字段内容" + key)
	}

	err = json.Unmarshal([]byte(rowStr), v)
	if err != nil {
		log.Println("redis hmget unmarshal error", err, key, v, len(arr[0]))
	}
	con.Close()

	return err

}

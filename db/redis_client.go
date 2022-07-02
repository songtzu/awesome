package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"log"

)

var redisClient *redis.Client
var ctx context.Context


//NewRedisPool
//"redis://<user>:<pass>@localhost:6379/<db>"
//"redis://<user>:<pass>@localhost:6379/<db>"
//因为redis没有用户名，因此，正确的形式是"redis://<pass>@localhost:6379/<db>"
func NewRedisPool(server string) (err error) {
	ctx = context.Background()

	opt, err := redis.ParseURL(server)
	if err != nil {
		log.Printf("redis url:%s format error:%s",server, err.Error())
		return err
	}

	redisClient = redis.NewClient(opt)

	return nil
}

func GetRedisClient() (*redis.Client) {
	return redisClient
}

func RedisKeyGet(key string) (cmd *redis.Cmd) {
	return redisClient.Do(ctx,"get",key)
}

func RedisKeySetStr(key string,v string, ttl time.Duration) (err error) {
	return redisClient.Set(ctx,key,v,ttl).Err()
}

func RedisKeySetObj(key string,v interface{}, ttl time.Duration) (err error) {
	if bin,err:=json.Marshal(v);err!=nil{
		log.Printf("redis set key obj failed :%s",err.Error())
		return err
	}else{
		return redisClient.Set(ctx,key,string(bin),ttl).Err()
	}
}

func RedisKeyGetStr(key string ) (v string,err error) {
	return redisClient.Get(ctx,key ).Result()
}

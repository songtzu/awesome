package app

import (
	"strings"
	"strconv"
	"sync"
)

/*
	0.0.0.0
	level：所处游戏服的层级 登录服 日志服 0层级 大厅服 游戏服1层级
	group: [0x00,0xFF]的整数,服务分组,比如所有game为一个组,hall一个组
	typeId: game下面的不同玩法
	instance:instance实例
*/

type APPID uint32
func (this APPID) GetLevel() uint32 {
	return (uint32(this)) >> 24
}

func (this APPID) GetGroup() uint32 {
	return (uint32(this)& 0x00ff0000) >> 16
}

func (this APPID)GetType() uint32 {
	return (uint32(this)& 0x0000ff00) >> 8
}

func (this APPID)GetInstanceId() uint32 {
	return (uint32(this)& 0x000000ff)
}

// appid 字符串转数字
func Str2AppId(strId string) APPID {
	arrId := strings.Split(strId, ".")
	if(len(arrId) != 4){
		return 0
	}

	level,_ := strconv.ParseInt(arrId[0], 10, 8)
	group,_ := strconv.ParseInt(arrId[1], 10, 8)
	typeId,_ := strconv.ParseInt(arrId[2], 10, 8)
	instance,_ := strconv.ParseInt(arrId[3], 10, 8)

	var appId uint32
	appId = uint32(level * 256 * 256 *256 + group * 256 * 256 + typeId * 256 + instance)
	return APPID(appId)
}
//日志中大量使用 缓存一份
type appIdCache struct{
	cache map[APPID]string
	wmutex sync.RWMutex
}

var idCache appIdCache
const (
	max_id_cache = 1000
)
// 数字转字符串 日志常用，
func AppId2Str(appid APPID) string{

	idCache.wmutex.RLock()
	strId, ok := idCache.cache[appid]
	idCache.wmutex.RUnlock()
	if ok {
		return strId
	}

	uintapp := uint32(appid)
	str0 := string(uintapp & 0xff000000 >> 24)
	str1 := string(uintapp & 0x00ff0000 >> 16)
	str2 := string(uintapp & 0x0000ff00 >> 8)
	str3 := string(uintapp & 0xff0000ff)

	strId = str0 + "." + str1 + "." + str2 + "." + str3
	if len(idCache.cache) < max_id_cache {
		idCache.wmutex.Lock()
		idCache.cache[appid] = strId
		idCache.wmutex.Unlock()
	}

	return strId
}

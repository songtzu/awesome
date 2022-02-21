package localcache

import (
	"awesome/alog"
	"sync"
	"time"
)

/*
 * 重定向到此服务的
		网关地址列表
 */
var ProxyAddressMap sync.Map

type ProxyInfoDetail struct {
	ActivityTime int64
	PlayerNum int
	Address string
}

func InsertProxyInfo(proxyAddress string)  {
	if _,ok:=ProxyAddressMap.Load(proxyAddress);ok{
		alog.Info("重复激活",proxyAddress)
		ProxyAddressMap.Store(proxyAddress,time.Now().Unix())
	}

}
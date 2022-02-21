package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"hjson-go"

	"fmt"
	"path/filepath"
	"os"
	"strings"
)

type CfgHall struct {
	HallAddress  string `json:"hallAddress"`  // 大厅地址
	//心跳间隔,毫秒
	HallHeartBeatInterval int `json:"hallHeartBeatInterval"`

	HallConnectTimeout int `json:"hall_connect_timeout"`
}
type CfgServer struct {
	/*
	 * accept tcp and websock.
	 * 		tcp://127.0.0.1:80
	 *		ws://127.0.0.1:80
	 */
	BindAddress    string `json:"bindAddress"`
	ServerID    int    `json:"serverId"`     // 本游戏服的id
	Version     string `json:"version"` // 版本
	AppID       int    `json:"appId"`  //
}
type CfgDeamon struct {
	SmtpConfig *configSmtp `json:"smtp"`
}
/**
 * 配置此字段，则房间层缓存均写入redis中。
 * 		此配置用于一个重大版本更新发布时，渡过版本不稳定期。
 */
type CfgCache struct {
	redisAddress string `json:"redisAddress"`
	redisDB string `json:"redisDb"`
}

// 配置文件读取
type Config struct {

	Server CfgServer             `json:"server"`
	Hall CfgHall                 `json:"hall"`
	Deamon CfgDeamon             `json:"deamon"`
	Additional        interface{} `json:"additional"`

	IgnoreCmd         []uint32          `json:"ignoreCmd"`           // 忽略的cmd

	IgnoreCmdMap map[uint32]bool
	//激活包的CMD，其他CMD不可以与此CMD冲突
	ActiveCmd uint32 `json:"activeCmd"`
}

type AlarmPhoneInfo struct {
	Phone    string `json:"phone"` //逗号隔开
	Key      string `json:"key"`
	TestSend bool   `json:"test_send"`
}

type Encrypt struct {
	AesKey string `json:"aes_key"`
	Open   bool   `json:"open"`
}
type configSmtp struct {
	User          string   `json:"user"`
	Password      string   `json:"password"`
	Host          string   `json:"host"`
	ToMailAddress []string `json:"mailAddressList"`

	CustomInter int `json:"custom_inter"` //hour
}

func (c *Config) IsCmdInsideIgnoreList(cmd uint32) bool {
	 _, ok := c.IgnoreCmdMap[cmd];
	 return ok
}
var appConfig *Config = nil

var configData []byte

func GetConfig() *Config {
	return appConfig
}

func GetAdditionalConfig() ([]byte, error) {
	if appConfig == nil {
	}
	return json.Marshal(appConfig.Additional)
}

var (
	CfgFilePath *string
	bindAddress *string
	webPort     *string
	serverId    *int
	serverIp    *string
	serverType  *int
)

func parseAges() {
	CfgFilePath = flag.String("configPath", "./config.json", "config path")
	bindAddress = flag.String("bindAddress", "", "bind address with protocol")
	serverId = flag.Int("sid", 0, "game server id")
	serverIp = flag.String("serverIp", "", "game server ip")
	serverType = flag.Int("serverType", 0, "game server type")
	flag.Parse()
}
func getEnvPath() string {
	paths:=os.Getenv("GOPATH")
	arr:=strings.Split(paths,";")
	for index,item:=range arr {
		fmt.Println(index,item)
		if strings.Contains(item,"awesome"){
			return item
		}

	}
	return ""
}
func init() {
	parseAges()
	appConfig = &Config{}
	app := filepath.Base(os.Args[0])
	fmt.Println(app)

	if strings.Contains(app, "Test") {
		cfgPath:=getEnvPath() + `\bin\config.json`
		log.Println("测试配置路径",cfgPath)
		hjson.ParseHjson(cfgPath, appConfig)
	}else {
		hjson.ParseHjson(*CfgFilePath, appConfig)
	}
	bin,_:=json.Marshal(appConfig)
	log.Println("config file path:", *CfgFilePath)
	log.Println("config:",string(bin))


	appConfig.IgnoreCmdMap = make(map[uint32]bool,0)
	for _,v:=range appConfig.IgnoreCmd {
		appConfig.IgnoreCmdMap[v] = true
	}
	return

}

func getLocalIp() string {
	addrSlice, err := net.InterfaceAddrs()
	if nil != err {
		fmt.Println("本地ip获取失败")
	}
	for _, addr := range addrSlice {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() {
				return ipnet.IP.String()
			}
		}
	}
	fmt.Println("本地ip获取失败")
	return ""
}

func GetExternal(address string) (string, error) {
	resp, err := http.Get(address)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content), nil
}

func ShowMap(m map[string]interface{}, f func(sfmt string, args ...interface{}), mapName string) {
	f(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>%20s\n", mapName)
	showMap(m, f, 0)
	f("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<%20s\n", "end")
}
func showMap(m map[string]interface{}, f func(sfmt string, args ...interface{}), level int) {
	for k, v := range m {

		if mm, ok := v.(map[string]interface{}); ok {
			f("%"+strconv.Itoa(level*5)+"v %-"+strconv.Itoa(50-level*5)+"v \n", "", k)
			showMap(mm, f, level+1)
		} else {
			f("%"+strconv.Itoa(level*5)+"v %-"+strconv.Itoa(50-level*5)+"v : %-20v\n", "", k, v)
		}
	}
}


func GetExtensionConfig() ([]byte,error) {

	return json.Marshal(appConfig.Additional)
}
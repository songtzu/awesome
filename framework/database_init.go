package framework

import (
	"awesome/config"
	"awesome/db"
)

func InitDatabase() (err error) {
	if config.GetConfig().Database.IsStart{
		if err = db.InitDB(config.GetConfig().Database.Address, config.GetConfig().Database.Protocol,config.GetConfig().Database.IsDebug);err!=nil{
			return err
		}
	}
	if config.GetConfig().Redis.IsStart {
		if err = db.NewRedisPool(config.GetConfig().Redis.Address);err!=nil{
			return err
		}
	}
	return nil
}

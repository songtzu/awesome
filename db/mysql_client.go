package db

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
)

//var database_addr string
//postgres/mysql
func InitMySql(address string, protocol string, isDebug bool) (err error) {
	err = nil
	if db == nil {
		client, err := NewClient(protocol, address)
		if err != nil {
			log.Printf("%s connect failed, url:%s",protocol, address)
		}
		db = client
		db.ShowSQL(isDebug)
	}
	return err
}
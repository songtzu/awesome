package db

import (
	_ "github.com/lib/pq"
	"log"
)

//var database_addr string
//postgres/mysql
func InitPg(address string, protocol string, isDebug bool) (err error) {
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
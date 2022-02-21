package test

import (
	"testing"
 
	"hjson-go"
)
//innner
type Info struct {
	Class string `json:"class"`
	Music string `json:"music"`
} 
//
type Sample struct {
	UserName string `json:"userName"`
	Age int `json:"age"`
	Info Info `json:"info"`

} 

func TestFromFile(t *testing.T) {
	s:=&Sample{}
	 hjson.ParseHjson("./sample.hjson",s)

}
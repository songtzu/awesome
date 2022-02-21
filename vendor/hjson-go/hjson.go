package hjson

import (
	"io/ioutil"
	"encoding/json"

	"fmt"

)

func getContent(file string) []byte {
	if data, err := ioutil.ReadFile(file); err != nil {
		panic(err)
	} else {
		return data
	}
}
func ParseHjson(filePath string, data interface{})  {
 	var output map[string]interface{}

	bin:=getContent(filePath)
 	err:=Unmarshal(bin,&output)
	if err!=nil{
		fmt.Println(err)
	}

	b,err:=json.Marshal(&output)
	json.Unmarshal(b,data)
	fmt.Println(data)

}
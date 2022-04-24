package mq

import (
	"encoding/json"
	"log"
	"net/http"
)

func StartHttpForMQ(address string)  {
	http.HandleFunc("/api/show_status", showStatus)
	http.HandleFunc("/api/pub_msg", publishMessage)
	http.ListenAndServe(address,nil)
}

type generalResponse struct {
	Code int
	Message string
	Data interface{}
}

func showStatus(w http.ResponseWriter, r *http.Request){
	log.Printf("HandleFunc")
	w.Write([]byte(string("HandleFunc")))
}

func publishMessage(w http.ResponseWriter, r *http.Request){
	resp:=&generalResponse{}
	if r.Method != "POST"{
		resp.Code = -10
		resp.Message = "仅限POST方法,不可用" +r.Method
	}else {

	}
	bin,_:=json.Marshal(resp)
	log.Printf("HandleFunc")
	w.Write(bin)
}

package mq

import (
	"awesome/anet"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

/*
 * mq的状态统计范畴：
 *	当前topic列表，
 *		每个topic的订阅者（id,订阅时间，）。
 *		每个topic的总订阅数。
 *		每个topic，最近100个消息。
 */
func showStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,POST")
	log.Printf("HandleFunc")

	arr := make([]*xmqSubImpl, 0)
	xmqInstance.topicMap.Range(func(key, value interface{}) bool {
		item := value.(*xmqTopic)
		log.Println(key, item)
		arr = append(arr, item.subs...)
		return true
	})
	b, err := json.Marshal(arr)
	if err != nil {
		log.Println("api/status,", err.Error())
	}
	log.Println(string(b))
	w.Write(b)
}

func publishDefaultMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,POST")
	log.Printf("通过http发布消息")
	resp := &generalResponse{Data: "默认回复"}
	var pub *publishReq = nil
	var err error
	if r.Method == "POST" {
		pub, err = parsePublishPost(r)
	} else if r.Method == "GET" {
		pub, err = parsePublishGet(r)
	}
	if err != nil {
		resp.Code = -1
		resp.Message = "error" + err.Error()
	} else {
		//var ack *anet.PackHead
		//var istimeout = true
		log.Println(pub.Action, string(pub.Body), pub.Cmd)
		switch pub.Action {
		case AmqCmdDefUnreliable2All:
			bridgePublishClient.PubUnreliableToAll([]byte(pub.Body), anet.PackHeadCmd(pub.Cmd))
		case AmqCmdDefUnreliable2RandomOne:
			bridgePublishClient.PubUnreliableToRandomOne([]byte(pub.Body), anet.PackHeadCmd(pub.Cmd))
		case AmqCmdDefReliable2RandomOne:
			log.Println("类型", pub.Action)
			//if ack, isTimeout := bridgePublishClient.PubReliableToRandomOne([]byte(pub.Body), anet.PackHeadCmd(pub.Cmd)); isTimeout {
			//	resp.Code = 911
			//	resp.Message = "time out"
			//} else {
			//	resp.Code = 0
			//	resp.Message = "ok"
			//	resp.Data = string(ack.Body)
			//	log.Println("正常的返回", resp.Data)
			//}
			var msg = &anet.PackHead{Cmd: uint32(pub.Cmd), ReserveLow: AmqCmdDefReliable2RandomOne, SequenceID: 0, Length: uint32(len(pub.Body)), Body: []byte(pub.Body)}
			var evtChan = make(chan *anet.PackHead, 1)

			pushReliableMsgFromHttpSvr(msg, evtChan)
			select {
			case ackMsg := <-evtChan:
				close(evtChan)
				log.Println("设置超时的网络IO正常返回", ackMsg)
				resp.Code = 0
				resp.Message = "ok"
				resp.Data = string(ackMsg.Body)
				bin, _ := json.Marshal(resp)
				w.Write(bin)
			case <-time.After(time.Millisecond * time.Duration(3000)):
				//log.Println("设置了超时的网络IO超时返回", msg, timeLimitMillisecond)
				//return nil, true
				resp.Code = 911
				resp.Message = "time out"
				bin, _ := json.Marshal(resp)
				w.Write(bin)
			}
			//
		case AmqCmdDefReliable2SpecOne:
			err = bridgePublishClient.PubUnreliableToAll([]byte(pub.Body), anet.PackHeadCmd(pub.Cmd))
		}
	}
	bin, _ := json.Marshal(resp)
	log.Printf("HandleFunc")
	w.Write(bin)

}

type publishReq struct {
	Body   string `json:"body"`
	Cmd    int    `json:"cmd"`
	Action int    `json:"action"`
}

func parsePublishPost(r *http.Request) (pub *publishReq, err error) {
	bin, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	} else {
		pub = &publishReq{}
		if err = json.Unmarshal(bin, pub); err != nil {
			return pub, err
		} else {
			return pub, nil
		}
	}
}

func parsePublishGet(r *http.Request) (pub *publishReq, err error) {
	pub = &publishReq{}
	v := r.URL.Query()
	body := v.Get("body")
	cmdStr := v.Get("cmd")
	cmd := -1
	actionStr := v.Get("action")
	action := 0
	if cmd, err = strconv.Atoi(cmdStr); err != nil {
		return nil, err
	}
	if action, err = strconv.Atoi(actionStr); err != nil {
		return nil, err
	}
	pub.Cmd = cmd
	pub.Body = body
	pub.Action = action
	return pub, nil

}

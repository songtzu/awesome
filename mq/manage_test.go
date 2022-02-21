package mq

import (
	"awesome/anet"
	"log"
	"testing"
	"time"
)

func TestNewAmqMessageModelManage(t *testing.T) {

	AmqMessageManage.Start()
	AmqMessageManage.SetMaxSize(3)

	var m1 = &AmqMessage{
		createTimestamp: time.Now().Unix(),
		msg: &anet.PackHead{
			Body: []byte("被挤掉的"),
		},
	}
	var m2 = &AmqMessage{
		createTimestamp: time.Now().Unix(),
		msg: &anet.PackHead{
			Body: []byte("正常执行"),
		},
	}

	var m3 = &AmqMessage{
		createTimestamp: time.Now().Unix(),
		msg: &anet.PackHead{
			Body: []byte("超时1"),
		},
	}
	var m5 = &AmqMessage{
		createTimestamp: time.Now().Unix() + 1,
		msg: &anet.PackHead{
			Body: []byte("超时2"),
		},
	}
	var m4 *AmqMessage
	AmqMessageManage.Add(m1)
	AmqMessageManage.Add(m2)
	AmqMessageManage.Add(m3)
	AmqMessageManage.Add(m5)

	go func() {
		time.Sleep(5 * time.Second)
		m4 = &AmqMessage{
			createTimestamp: time.Now().Unix(),
			msg: &anet.PackHead{
				Body: []byte("阻塞获取"),
			},
		}
		AmqMessageManage.Add(m4)
	}()

	log.Println("想处理的消息:", string(AmqMessageManage.Get().(*AmqMessage).msg.Body))
	log.Printf("ok:%v \n", 2 == AmqMessageManage.FreeLen() && 1 == AmqMessageManage.UsingLen())
	Success(m2)
	log.Printf("ok:%v \n", 2 == AmqMessageManage.FreeLen() && 0 == AmqMessageManage.UsingLen())
	Success(m2)
	log.Printf("ok:%v \n", 2 == AmqMessageManage.FreeLen())
	time.Sleep(5 * time.Second)
	log.Printf("ok:%v \n", 0 == AmqMessageManage.FreeLen())
	Success(m3)
	Success(m5)

	// 删除
	log.Println("想处理的消息:", string(AmqMessageManage.Get().(*AmqMessage).msg.Body))
	Success(m4)
	Success(m4)
	time.Sleep(time.Millisecond * 15)
	log.Printf("ok:%v \n", 0 == AmqMessageManage.FreeLen() && 0 == AmqMessageManage.HeapLen() && 0 == AmqMessageManage.UsingLen())
	time.Sleep(time.Hour)
}

func Success(model *AmqMessage) {
	AmqMessageManage.Start()
	AmqMessageManage.Add(&AmqMessage{
		createTimestamp: time.Now().Unix(),
		msg: &anet.PackHead{
			Body: []byte("3秒超时"),
		},
	})
	AmqMessageManage.AddNoTimeOut(&AmqMessage{
		createTimestamp: time.Now().Unix(),
		msg: &anet.PackHead{
			Body: []byte("无超时"),
		},
	})
	AmqMessageManage.AddWithTimeout(&AmqMessage{
		createTimestamp: time.Now().Unix(),
		msg: &anet.PackHead{
			Body: []byte("2秒超时"),
		},
	},time.Second * 2)
	<- chan bool(nil)
}

func TestAmqMessageModelManage_Add(t *testing.T) {
	AmqMessageManage.Start()
	log.Println("start ...")
	AmqMessageManage.Add(&AmqMessage{
		createTimestamp: time.Now().Unix(),
		msg: &anet.PackHead{
			Body: []byte("3秒超时"),
		},
	})
	AmqMessageManage.AddNoTimeOut(&AmqMessage{
		createTimestamp: time.Now().Unix(),
		msg: &anet.PackHead{
			Body: []byte("无超时"),
		},
	})
	AmqMessageManage.AddWithTimeout(&AmqMessage{
		createTimestamp: time.Now().Unix(),
		msg: &anet.PackHead{
			Body: []byte("2秒超时"),
		},
	},time.Second * 2)
	<- chan bool(nil)
}

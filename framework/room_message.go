package framework

import (
	"runtime/debug"

	"awesome/alog"
)

func (r *Room) enqueueMessage(msg *UserMessage) {
	alog.Info("enqueue message into room chan")

	r.workerChan <- msg
}

func (r *Room) enqueueSystemMessage(msg *SystemMessage) {
	r.sysMsg <- msg
}

func (r *Room) taskWorker(f func()) {
	alog.Debug("room task worker start")
	go func() {
		defer func() {
			if err := recover(); err != nil {
				alog.Err("framework panic", err, string(debug.Stack()))
			}
		}()

		f()
	}()
}

func (r *Room) recoverWorker() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				alog.Err("panic catch by recoverWorker,", err, string(debug.Stack()))
			}
		}()

		r.roomWorkerLoop()
	}()
}

//系统房间消息处理
func (r *Room) roomSystemMsgEntry(msg *SystemMessage) error {
	if msg.Cmd == SystemMessageDefTimer {
		msg.DealHandle.(TypeTimeTaskCallBack)(msg.Msg, r.GetRoomData())
	}
	return nil
}

func (r *Room) roomWorkerLoop() {
	var err error

	for {
		select {
		case msg := <-r.sysMsg:
			if msg == nil {
				return
			}

			alog.Debug("room roomWorkerLoop awake. got a system message.")
			recoverWorker(func() {
				err = r.roomSystemMsgEntry(msg)
			})

		case msg := <-r.workerChan: // client post message.
			if msg == nil {
				return
			}
			alog.Debug("room worker, got a normal message ", msg.pack, string(msg.pack.Body))
			recoverWorker(func() {
				err := frameworkInterfaceInstance.OnDispatchLogicMessage(r.RoomCode, r, msg.user, msg.pack)
				if err != nil {
					alog.Debug(string(debug.Stack()))
				}
			})
		}

		if err != nil {
			alog.Err("err:", err)
		}
	}
}

func recoverWorker(f func(), panicInfo ...string) {
	func() {
		defer func() {
			if err := recover(); err != nil {
				alog.Err("panic？。。。", err, string(debug.Stack()), panicInfo)
			}
		}()

		f()
	}()
}

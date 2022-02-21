package msghandler

import (
	"awesome/anet"
	"fmt"
)

type MsgHandler struct {

}

func (this *MsgHandler) OnDispatchMessage(msg *anet.PackHead) {
	fmt.Println("head:",msg)
	fmt.Println("sub callback",string(msg.Body))
}
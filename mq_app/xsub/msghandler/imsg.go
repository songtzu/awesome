package msghandler

import "awesome/anet"

type IMessageDispatch interface {

	OnDispatchMessage(message *anet.PackHead)

}

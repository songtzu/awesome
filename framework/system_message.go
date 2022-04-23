package framework

import "awesome/anet"

const (
	SystemMessageDefUndefined =0
	SystemMessageDefError = 1
	SystemMessageDefOffline = 2
	SystemMessageDefTimeOut =3
	SystemMessageDefTimer =4
	SystemMessageMatchEvent = 10
)

type SystemMessage struct {
	Cmd        int
	Msg        string
	DealHandle interface{}
}
//
type UserMessage struct {
	pack *anet.PackHead
	user *PlayerImpl
}

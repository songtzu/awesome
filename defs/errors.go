package defs

import (
	"errors"
	"fmt"
)

//var (
//	ErrNotFound = errors.New("not found"),
//	Err
//)

const (
	//未知错误
	ErrorDefUndefinedError = 1
	ErrorDefDefault = 2
	ErrorDefDataNotFound = 3
	ErrorDefUnImplementInterface = 4
	ErrorDefMultiImplementInterface = 5
	//logic specified error
	ErrorDefFailedParseRoomCode = 101
	ErrorDefFailedCreateRoom = 102
	ErrorDefUnMatchedRoomCodeDuringParseAndCreate = 103
)

var errMsg = map[int]error{
	ErrorDefUndefinedError:errors.New("undefined errors"),
	ErrorDefDefault:errors.New("default error"),
	ErrorDefDataNotFound:errors.New("data item not found"),
	ErrorDefUnImplementInterface:errors.New("interface not implemented"),
	ErrorDefMultiImplementInterface:errors.New("multi implemented interface"),
	ErrorDefFailedParseRoomCode:errors.New("bad error when try to call parse roomCode for a new user"),
	ErrorDefFailedCreateRoom:errors.New("bad interface callback when framework try to call logic implement to create a room"),
	ErrorDefUnMatchedRoomCodeDuringParseAndCreate:errors.New("got different roomCode during parse and create"),
}

func GetError(err int) error {
	if v,ok:=errMsg[err];ok{
		return v
	}
	return errors.New(fmt.Sprintf("undefined error for error code:%d",err))
}
package mq

import (
	"math"
	"sync"
)

func Contains[T comparable](slice []T, target T) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}

	return false
}

var sequenceId uint32 = 100000
var sequenceIdLocker *sync.Mutex

//生成序列号
func genSequenceId() uint32 {
	sequenceIdLocker.Lock()
	if sequenceId >= math.MaxUint32 {
		sequenceId = 1000
	}
	sequenceId += 1
	sequenceIdLocker.Unlock()
	return sequenceId
}

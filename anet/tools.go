package anet

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

//var sequenceId uint32 = 100000
//var sequenceIdLocker *sync.Mutex

type sequenceIdGen struct {
	SequenceId uint32
	sync.Mutex
}

const startIndexForSequenceId = 1000000

var sequenceIdGenInstance = &sequenceIdGen{
	SequenceId: startIndexForSequenceId,
}

/*AllocateNewSequenceId
 *生成序列号
 */
func AllocateNewSequenceId() (id uint32) {
	sequenceIdGenInstance.Lock()
	if sequenceIdGenInstance.SequenceId >= math.MaxUint32 {
		sequenceIdGenInstance.SequenceId = startIndexForSequenceId
	}
	sequenceIdGenInstance.SequenceId += 1
	id = sequenceIdGenInstance.SequenceId
	sequenceIdGenInstance.Unlock()
	return id
}

type IdGen struct {
	Id int
	sync.Mutex
}

func GenNewId() (result int) {
	idGen.Lock()
	idGen.Id += 1
	result = idGen.Id
	idGen.Unlock()
	return result
}

var idGen = &IdGen{Id: 0}

package anet

import (
	"sync"
	"time"
)

type TestInfo struct {
	TotalCount int
	FailCount  int
	PassCount  int
	TimeCost   int64
	Start      time.Time
	ThreadCount int
	SetCountEachThread int
	SetTotalCount int
	sync.Mutex
}

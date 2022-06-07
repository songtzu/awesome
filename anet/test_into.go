package anet

import (
	"log"
	"sync"
	"time"
)

type TestInfo struct {
	CurrentTotalCount  int
	FailCount          int
	PassCount          int
	TimeCost           int64
	Start              time.Time
	ThreadCount        int
	SetCountEachThread int
	SetTotalCount      int
	sync.Mutex
}

func (t *TestInfo) UpdateTotalSetCount() {
	t.SetTotalCount = t.SetCountEachThread * t.ThreadCount
}

func (t *TestInfo) tickEnd() {
	t.TimeCost = time.Now().Sub(t.Start).Milliseconds()
}

func (t *TestInfo) AddPassCount() {
	t.Lock()
	t.PassCount += 1
	t.CurrentTotalCount += 1
	t.Unlock()
}

func (t *TestInfo) AddFailCount() {
	t.Lock()
	t.FailCount += 1
	t.CurrentTotalCount += 1
	t.Unlock()
}

func (t *TestInfo) TryPrint() bool {
	if t.CurrentTotalCount == t.SetTotalCount {
		t.PrintTestResult()
		return true
	}
	return false
}

func (t *TestInfo) PrintTestResult() {
	t.tickEnd()
	log.Printf("thead:%d,totalCount:%d,failCount:%d,passCount:%d, timeCost:%d, avg:%f", t.ThreadCount, t.CurrentTotalCount,
		t.FailCount, t.PassCount, t.TimeCost, float64(t.TimeCost)/float64(t.CurrentTotalCount))
}

package anet

import "time"

type TestInfo struct {
	TotalCount int
	FailCount  int
	PassCount  int
	TimeCost   int64
	Start      time.Time
}

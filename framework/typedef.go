package framework

type RoomID int64
type CMD int32

type SeatNumber int

type MatchRule struct {
	IsMatch bool
	MatchNum int
	DeadlineTimestamp int64

}
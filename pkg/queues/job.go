package queues

import "time"

// JobRequest type
type JobRequest struct {
	Relevance        float64
	Delay            time.Duration
	TimeInSystem     time.Duration
	TimeInProcessing time.Duration
	ArrivalTime      time.Time
	DispatchTime     time.Time
	Id               int
	IsFinished       bool
}

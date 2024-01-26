package ratelimiter

import (
	"time"
)

type Rate struct {
	NumberOfRequests uint
	Interval         time.Time
}

// ???
func getIntervalInEpochSeconds() uint64 {
	return uint64(time.Now().Unix())
}

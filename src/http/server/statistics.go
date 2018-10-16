package server

import (
	"sync"
)

type Statistics struct {
	sync.RWMutex            `json:"-"`
	RequestCount            int64 `json:"total"`
	TotalTimeForAllRequests int64 `json:"-"`
	AverageRequestTime      int64 `json:"average"`
}

func (stats *Statistics) IncrementRequestCount() {
	stats.Lock()
	defer stats.Unlock()

	stats.RequestCount++
}

func (stats *Statistics) AddDuration(duration int64) {
	stats.Lock()
	defer stats.Unlock()

	stats.TotalTimeForAllRequests += duration
}

func (stats *Statistics) UpdateAverageRequestDuration() {
	stats.Lock()
	defer stats.Unlock()

	if stats.RequestCount < 1 {
		stats.AverageRequestTime = 0
		return
	}

	stats.AverageRequestTime = stats.TotalTimeForAllRequests / stats.RequestCount
}
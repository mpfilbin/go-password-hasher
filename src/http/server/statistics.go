package server

import (
	"sync"
	"time"
)


type Statistics struct {
	sync.RWMutex `json:"-"`
	RequestCount       uint64 `json:"total"`
	RequestDurations   []time.Duration `json:"-"`
	AverageRequestTime int64 `json:"average"`
}

func (stats *Statistics) IncrementRequestCount() {
	stats.Lock()
	defer stats.Unlock()

	stats.RequestCount++
}

func (stats *Statistics) AddDuration(duration time.Duration) {
	stats.Lock()
	stats.RequestDurations = append(stats.RequestDurations, duration)
	stats.Unlock()
}

func (stats *Statistics) UpdateAverageRequestDuration() {
	stats.Lock()
	defer stats.Unlock()

	var totalTimeForAllRequests int64
	requestCount := len(stats.RequestDurations)

	if requestCount < 1 {
		stats.AverageRequestTime = 0
		stats.Unlock()
		return
	}

	for i := 0; i < requestCount; i++ {
		totalTimeForAllRequests += stats.RequestDurations[i].Nanoseconds() * 1000
	}

	stats.AverageRequestTime = totalTimeForAllRequests / int64(requestCount)
}
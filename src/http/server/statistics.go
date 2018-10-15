package server

import (
	"sync"
	"time"
)


type Statistics struct {
	sync.RWMutex `json:"-"`
	RequestCount       int64 `json:"total"`
	RequestDurations   []int64 `json:"-"`
	AverageRequestTime int64 `json:"average"`
}

func (stats *Statistics) IncrementRequestCount() {
	stats.Lock()
	defer stats.Unlock()

	stats.RequestCount++
}

func (stats *Statistics) AddDuration(duration int64) {
	stats.Lock()
	defer stats.Unlock()

	stats.RequestDurations = append(stats.RequestDurations, duration)
}

func (stats *Statistics) UpdateAverageRequestDuration() {
	stats.Lock()
	defer stats.Unlock()

	var totalTimeForAllRequests int64
	requestCount := len(stats.RequestDurations)

	if requestCount < 1 {
		stats.AverageRequestTime = 0
		return
	}

	for i := 0; i < requestCount; i++ {
		totalTimeForAllRequests += stats.RequestDurations[i] / int64(time.Microsecond)
	}

	stats.AverageRequestTime = totalTimeForAllRequests / int64(requestCount)
}
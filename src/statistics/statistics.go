package statistics

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type statistics struct {
	sync.RWMutex `json:"-"`
	RequestCount       uint64 `json:"total"`
	RequestDurations   []time.Duration `json:"-"`
	AverageRequestTime int64 `json:"average"`
}

func (stats *statistics) IncrementRequestCount() {
	stats.Lock()
	stats.RequestCount++
	stats.Unlock()
}

func (stats *statistics) AddDuration(duration time.Duration) {
	stats.Lock()
	stats.RequestDurations = append(stats.RequestDurations, duration)
	stats.Unlock()
}

func (stats *statistics) UpdateAverageRequestDuration() {
	stats.Lock()
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

	stats.Unlock()

	stats.AverageRequestTime = totalTimeForAllRequests / int64(requestCount)
}

type requestHandler func(w http.ResponseWriter, r *http.Request)


var sharedStatistics = &statistics{}


func TrackTiming(handler requestHandler) func(http.ResponseWriter, *http.Request){

	return func (response http.ResponseWriter, request *http.Request) {
		log.Printf("Received %v %v", request.Method, request.URL.Path)
		sharedStatistics.IncrementRequestCount()

		start := time.Now()
		handler(response, request)
		duration := time.Since(start)

		sharedStatistics.AddDuration(duration)
		log.Printf("Request handled in %s milliseconds", duration)

	}

}

func Report(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		sharedStatistics.UpdateAverageRequestDuration()
		jsonContent, err := json.Marshal(sharedStatistics)

		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		response.Header().Set("Content-Type", "application/json")
		response.Write(jsonContent)
	default:
		http.NotFound(response, request)
	}

}
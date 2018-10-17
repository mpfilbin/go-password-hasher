package server

import (
	"fmt"
)

func ExampleStatistics_IncrementRequestCount() {
	stats := &statistics{}

	stats.IncrementRequestCount()

	fmt.Println(stats.RequestCount)
	// Output: 1
}

func ExampleStatistics_AddDuration() {
	stats := &statistics{}

	stats.AddDuration(5)
	stats.AddDuration(6)
	stats.AddDuration(7)

	fmt.Println(stats.TotalTimeForAllRequests)
	// Output: 18
}

func ExampleStatistics_UpdateAverageRequestDuration() {
	stats := &statistics{}

	stats.IncrementRequestCount()
	stats.AddDuration(10)
	stats.IncrementRequestCount()
	stats.AddDuration(16)

	stats.UpdateAverageRequestDuration()

	fmt.Println(stats.AverageRequestTime)
	// Output: 13
}

func ExampleStatistics_UpdateAverageRequestDurationWithZeroRequests() {
	stats := &statistics{}
	stats.UpdateAverageRequestDuration()
	fmt.Println(stats.AverageRequestTime)
	// Output: 0
}
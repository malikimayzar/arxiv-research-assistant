package client

import (
	"fmt"
	"time"

	"github.com/sony/gobreaker"
)

func newCircuitBreaker(name string) *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        name,
		MaxRequests: 3,                // max request saat half-open
		Interval:    60 * time.Second, // window untuk counting failures
		Timeout:     30 * time.Second, // berapa lama open sebelum coba half-open
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Trip kalau 5+ failure atau failure rate > 60%
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 5 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			fmt.Printf("⚡ Circuit breaker [%s]: %s → %s\n", name, from, to)
		},
	})
}

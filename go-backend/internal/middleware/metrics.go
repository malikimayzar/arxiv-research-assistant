package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 5, 10, 30, 60, 120, 300},
		},
		[]string{"method", "path"},
	)

	queryFaithfulness = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "query_faithfulness_score",
			Help: "Average faithfulness score of last queries",
		},
	)

	activeGoroutines = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_goroutines",
			Help: "Number of active goroutines",
		},
	)
)

func Metrics() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Response().StatusCode())
		path := c.Route().Path

		httpRequestsTotal.WithLabelValues(c.Method(), path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Method(), path).Observe(duration)

		return err
	}
}

func SetFaithfulness(score float64) {
	queryFaithfulness.Set(score)
}

func SetActiveGoroutines(n float64) {
	activeGoroutines.Set(n)
}

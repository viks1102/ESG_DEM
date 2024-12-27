package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Metric to track HTTP request duration latency
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of HTTP request durations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Metric to track total HTTP requests by status code
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// Metric to track the number of active goroutines
	goroutines = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "goroutines_count",
			Help: "Number of active goroutines",
		},
	)
)

func initMetrics() {
	// TASK: What metrics would be useful for this type of service?
	// Registering the metrics
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(goroutines)
}

// Track request duration
func trackRequestDuration(method, path string, start time.Time) {
	duration := time.Since(start).Seconds()
	requestDuration.WithLabelValues(method, path).Observe(duration)
}

// Track HTTP request count
func trackRequestCount(method, path string, status int) {
	httpRequestsTotal.WithLabelValues(method, path, fmt.Sprintf("%d", status)).Inc()
}

// Update active goroutine count
func updateGoroutineCount() {
	goroutines.Set(float64(runtime.NumGoroutine()))
}


func startMetricsServer(cfg ConfigMetrics, wg *sync.WaitGroup) {
	defer wg.Done()

	port := fmt.Sprintf(":%v", cfg.Port)
	slog.With(slog.Any("port", port)).Info("metrics server started")

	// Update goroutine count periodically
	go func() {
		for {
			updateGoroutineCount()
			time.Sleep(10 * time.Second) // Update every 10 seconds
		}
	}()

	http.Handle(cfg.Path, promhttp.Handler())
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Printf("error with listen and serve %v", err.Error())
	}
}

/* Changes made:
	1.	Metrics Definitions:
	•	Added requestDuration, httpRequestsTotal, and goroutines as the metrics you want to track.
	•	Registered the metrics using prometheus.MustRegister() in the initMetrics() function.
	2.	Request Duration Metric:
	•	trackRequestDuration() is used to track how long each request takes (latency).
	•	It calculates the duration using time.Since(start).Seconds() and observes it in the http_request_duration_seconds histogram.
	3.	HTTP Request Count by Status Code:
	•	trackRequestCount() tracks the total number of requests by method, path, and status code (using http_requests_total).
	4.	Active Goroutines Metric:
	•	updateGoroutineCount() periodically updates the number of active goroutines (goroutines_count).
	•	A goroutine is spawned to update the count every 10 seconds.
	5.	Server Handling:
	•	The startMetricsServer() function starts the metrics server as before, but now it also handles the periodic update of active goroutines.

	*/
package http

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "Total number of HTTP requests grouped by route, method, and status code.",
		},
		[]string{"route", "method", "status_code"},
	)

	httpLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_latency_seconds",
			Help:    "HTTP request latency in seconds grouped by route, method, and status code.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"route", "method", "status_code"},
	)

	httpErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP responses with status code >= 500.",
		},
		[]string{"route", "method", "status_code"},
	)
)

func RegisterMetrics(registry *prometheus.Registry) {
	registry.MustRegister(httpRequestTotal, httpLatency, httpErrors)
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}

		statusCode := strconv.Itoa(c.Writer.Status())
		labels := prometheus.Labels{
			"route":       route,
			"method":      c.Request.Method,
			"status_code": statusCode,
		}

		httpRequestTotal.With(labels).Inc()
		httpLatency.With(labels).Observe(time.Since(start).Seconds())
		if c.Writer.Status() >= 500 {
			httpErrors.With(labels).Inc()
		}
	}
}

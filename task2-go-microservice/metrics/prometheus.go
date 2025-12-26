package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Request duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func init() {
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(requestDuration)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		totalRequests.WithLabelValues(r.Method, r.URL.Path).Inc()
		next.ServeHTTP(w, r)
		duration := time.Since(start).Seconds()
		requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}

func Handler() http.HandlerFunc {
	return promhttp.Handler().ServeHTTP
}

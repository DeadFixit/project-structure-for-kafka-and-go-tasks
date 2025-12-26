package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "analytics_http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path"},
	)
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "analytics_http_request_duration_seconds",
			Help:    "Request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
	AnomaliesTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "analytics_anomalies_total",
		Help: "Anomalies detected by z-score",
	})
)

func init() {
	prometheus.MustRegister(RequestsTotal)
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(AnomaliesTotal)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		RequestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()
		next.ServeHTTP(w, r)
		RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
	})
}

func Handler() http.Handler {
	return promhttp.Handler()
}

package utils

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

var (
	onceLimiter sync.Once
	limiter     *rate.Limiter
)

func getLimiter() *rate.Limiter {
	onceLimiter.Do(func() {
		limiter = rate.NewLimiter(rate.Limit(1000), 5000)
	})
	return limiter
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !getLimiter().Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

package middleware

import (
	"log"
	"net/http"
	"time"
)

// MetricsMiddleware logs request processing time and basic metrics
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request start
		log.Printf("Request started: %s %s", r.Method, r.URL.Path)

		// Call next handler
		next.ServeHTTP(w, r)

		// Log request completion with timing
		duration := time.Since(start)
		log.Printf("Request processed: %s %s in %v", r.Method, r.URL.Path, duration)
	})
}

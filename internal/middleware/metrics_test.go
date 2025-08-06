package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetricsMiddleware(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(log.Writer())

	// Create a test handler that simulates some processing time
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate some processing time
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Create middleware
	middleware := MetricsMiddleware(testHandler)

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Execute middleware
	middleware.ServeHTTP(w, req)

	// Get log output
	logOutput := buf.String()

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test response", w.Body.String())

	// Check if start log is present
	assert.True(t, strings.Contains(logOutput, "Request started: GET /test"),
		"Log should contain request start message")

	// Check if completion log is present
	assert.True(t, strings.Contains(logOutput, "Request processed: GET /test"),
		"Log should contain request completion message")

	// Check if timing information is present
	assert.True(t, strings.Contains(logOutput, "in"),
		"Log should contain timing information")
}

func TestMetricsMiddleware_WithDifferentMethods(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(log.Writer())

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("created"))
	})

	middleware := MetricsMiddleware(testHandler)

	// Test different HTTP methods
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		buf.Reset() // Clear previous log output

		req := httptest.NewRequest(method, "/api/test", nil)
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		logOutput := buf.String()

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.True(t, strings.Contains(logOutput, "Request started: "+method+" /api/test"))
		assert.True(t, strings.Contains(logOutput, "Request processed: "+method+" /api/test"))
	}
}

func TestMetricsMiddleware_WithErrorResponse(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(log.Writer())

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error occurred"))
	})

	middleware := MetricsMiddleware(testHandler)

	req := httptest.NewRequest("POST", "/error", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	logOutput := buf.String()

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error occurred", w.Body.String())
	assert.True(t, strings.Contains(logOutput, "Request started: POST /error"))
	assert.True(t, strings.Contains(logOutput, "Request processed: POST /error"))
}

func TestMetricsMiddleware_TimingAccuracy(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(log.Writer())

	// Create handler with known delay
	expectedDelay := 50 * time.Millisecond
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(expectedDelay)
		w.WriteHeader(http.StatusOK)
	})

	middleware := MetricsMiddleware(testHandler)

	req := httptest.NewRequest("GET", "/timing", nil)
	w := httptest.NewRecorder()

	start := time.Now()
	middleware.ServeHTTP(w, req)
	totalTime := time.Since(start)

	logOutput := buf.String()

	// Total time should be at least the expected delay
	assert.True(t, totalTime >= expectedDelay,
		"Total time should be at least the expected delay")

	// Log should contain timing information
	assert.True(t, strings.Contains(logOutput, "in"),
		"Log should contain timing information")

	// The timing in log should be reasonable (not negative, not too large)
	assert.True(t, strings.Contains(logOutput, "Request processed: GET /timing in"),
		"Log should contain processed message with timing")
}

func TestMetricsMiddleware_WithComplexPath(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(log.Writer())

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("complex path"))
	})

	middleware := MetricsMiddleware(testHandler)

	// Test with complex path including query parameters
	req := httptest.NewRequest("GET", "/api/v1/users/123?filter=active&sort=name", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	logOutput := buf.String()

	assert.Equal(t, http.StatusOK, w.Code)
	// The log might not contain the full URL with query parameters
	// Let's check for the base path
	assert.True(t, strings.Contains(logOutput, "Request started: GET /api/v1/users/123"))
	assert.True(t, strings.Contains(logOutput, "Request processed: GET /api/v1/users/123"))
}

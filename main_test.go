package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	// Create request
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Execute
	healthCheckHandler(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	// Parse response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check response structure
	assert.Equal(t, "healthy", response["status"])
	assert.NotEmpty(t, response["time"])

	// Validate time format (RFC3339)
	_, err = time.Parse(time.RFC3339, response["time"])
	assert.NoError(t, err, "Time should be in RFC3339 format")
}

func TestHealthCheckHandler_ContentType(t *testing.T) {
	// Create request
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Execute
	healthCheckHandler(w, req)

	// Assertions
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestHealthCheckHandler_ResponseStructure(t *testing.T) {
	// Create request
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Execute
	healthCheckHandler(w, req)

	// Parse response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check all expected fields are present
	assert.Contains(t, response, "status")
	assert.Contains(t, response, "time")

	// Check values
	assert.Equal(t, "healthy", response["status"])
	assert.NotEmpty(t, response["time"])
}

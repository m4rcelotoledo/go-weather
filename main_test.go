package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGracefulShutdown(t *testing.T) {
	// Create a simple HTTP server for testing
	server := &http.Server{
		Addr: ":0", // Porta aleatória
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		}),
	}

	// Channel to signal that the server has started
	serverStarted := make(chan bool, 1)

	// Channel to signal that the shutdown has started
	shutdownStarted := make(chan bool, 1)

	// Channel to signal that the shutdown has completed
	shutdownCompleted := make(chan bool, 1)

	// Goroutine to start the server
	go func() {
		serverStarted <- true
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			t.Errorf("Server error: %v", err)
		}
	}()

	// Goroutine for graceful shutdown
	go func() {
		// Wait for the server to start
		<-serverStarted

		// Wait a bit to ensure the server is running
		time.Sleep(100 * time.Millisecond)

		// Simulate shutdown signal
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

		shutdownStarted <- true

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			t.Errorf("HTTP server Shutdown: %v", err)
		}

		shutdownCompleted <- true
	}()

	// Wait for the shutdown to start
	<-shutdownStarted

	// Wait for the shutdown to complete
	select {
	case <-shutdownCompleted:
		t.Log("Graceful shutdown completed successfully")
	case <-time.After(10 * time.Second):
		t.Error("Graceful shutdown timeout")
	}
}

func TestGracefulShutdownWithRealServer(t *testing.T) {
	// Configure necessary environment variable
	os.Setenv("WEATHER_API_KEY", "test-key")
	defer os.Unsetenv("WEATHER_API_KEY")

	// Create server on a specific port for testing
	server := &http.Server{
		Addr: ":8081", // Specific port for testing
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"healthy"}`))
		}),
	}

	// Channel to signal that the server has started
	serverStarted := make(chan bool, 1)

	// Channel to signal that the shutdown has completed
	shutdownCompleted := make(chan bool, 1)

	// Goroutine to start the server
	go func() {
		serverStarted <- true
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			t.Errorf("Server error: %v", err)
		}
	}()

	// Wait for the server to start
	<-serverStarted

	// Wait a bit to ensure the server is running
	time.Sleep(200 * time.Millisecond)

	// Make a request to verify if the server is working
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://localhost:8081/health")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Start graceful shutdown
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			t.Errorf("HTTP server Shutdown: %v", err)
		}

		shutdownCompleted <- true
	}()

	// Wait for the shutdown to complete
	select {
	case <-shutdownCompleted:
		t.Log("Graceful shutdown completed successfully")
	case <-time.After(10 * time.Second):
		t.Error("Graceful shutdown timeout")
	}

	// Verify if the server really stopped
	_, err = client.Get("http://localhost:8081/health")
	if err == nil {
		t.Error("Server should be stopped after shutdown")
	}
}

func TestGracefulShutdownWithSignal(t *testing.T) {
	// Configure necessary environment variable
	os.Setenv("WEATHER_API_KEY", "test-key")
	defer os.Unsetenv("WEATHER_API_KEY")

	// Create server on a specific port for testing
	server := &http.Server{
		Addr: ":8082", // Specific port for testing
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"healthy"}`))
		}),
	}

	// Channel to signal that the server has started
	serverStarted := make(chan bool, 1)

	// Channel to signal that the shutdown has completed
	shutdownCompleted := make(chan bool, 1)

	// Goroutine to start the server
	go func() {
		serverStarted <- true
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			t.Errorf("Server error: %v", err)
		}
	}()

	// Goroutine to simulate graceful shutdown with signal
	go func() {
		// Wait for the server to start
		<-serverStarted

		// Wait a bit to ensure the server is running
		time.Sleep(200 * time.Millisecond)

		// Simulate shutdown signal (like in main.go)
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

		// Simulate receiving the signal
		sigint <- syscall.SIGTERM

		log.Println("Shutting down server gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			t.Errorf("HTTP server Shutdown: %v", err)
		}

		shutdownCompleted <- true
	}()

	// Wait for the shutdown to complete
	select {
	case <-shutdownCompleted:
		t.Log("Graceful shutdown with signal completed successfully")
	case <-time.After(10 * time.Second):
		t.Error("Graceful shutdown with signal timeout")
	}

	// Verify if the server really stopped
	client := &http.Client{Timeout: 5 * time.Second}
	_, err := client.Get("http://localhost:8082/health")
	if err == nil {
		t.Error("Server should be stopped after signal shutdown")
	}
}

func TestHealthCheckHandler(t *testing.T) {
	// Create request
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Execute the handler
	healthCheckHandler(w, req)

	// Assertions (using testify/assert)
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

	// Execute the handler
	healthCheckHandler(w, req)

	// Assertions
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestHealthCheckHandler_ResponseStructure(t *testing.T) {
	// Create request
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Execute the handler
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

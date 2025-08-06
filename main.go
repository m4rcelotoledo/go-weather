package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-weather/internal/clients/viacep"
	"go-weather/internal/clients/weather"
	"go-weather/internal/config"
	"go-weather/internal/handlers"
	metrics "go-weather/internal/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// Health check handler
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"status": "healthy",
		"time":   time.Now().UTC().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Check if API key is configured
	if cfg.WeatherAPIKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable is required")
	}

	// Create clients
	viaCEPClient := viacep.NewClient()
	weatherClient := weather.NewClient()

	// Create handler
	weatherHandler := handlers.NewWeatherHandler(viaCEPClient, weatherClient, cfg)

	// Configure router
	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(metrics.MetricsMiddleware)

	// Define routes
	r.Route("/", func(r chi.Router) {
		r.Get("/health", healthCheckHandler)
		r.Get("/cep/{cep}", weatherHandler.GetWeatherByCEP)
	})

	// Create server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}()

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	log.Println("Server stopped")
}

package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"go-weather/internal/clients/viacep"
	"go-weather/internal/clients/weather"
	"go-weather/internal/config"

	"github.com/go-chi/chi/v5"
)

type WeatherHandler struct {
	viaCEPClient  viacep.Client
	weatherClient weather.Client
	config        *config.Config
}

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func NewWeatherHandler(viaCEPClient viacep.Client, weatherClient weather.Client, cfg *config.Config) *WeatherHandler {
	return &WeatherHandler{
		viaCEPClient:  viaCEPClient,
		weatherClient: weatherClient,
		config:        cfg,
	}
}

func (h *WeatherHandler) GetWeatherByCEP(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")

	// Log the incoming request
	log.Printf("Processing CEP: %s", cep)

	// Validate CEP (8 numeric digits)
	if !isValidCEP(cep) {
		log.Printf("Invalid CEP format: %s", cep)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode("invalid zipcode")
		return
	}

	// Find city by CEP
	address, err := h.viaCEPClient.GetAddressByCEP(cep)
	if err != nil {
		log.Printf("CEP not found: %s - Error: %v", cep, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("can not find zipcode")
		return
	}

	// Log found city
	log.Printf("Found city: %s for CEP: %s", address.Localidade, cep)

	// Check if API key is configured
	if h.config.WeatherAPIKey == "" {
		log.Printf("Weather API key not configured")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Weather API key not configured")
		return
	}

	// Get current weather
	weatherData, err := h.weatherClient.GetCurrentWeather(address.Localidade, h.config.WeatherAPIKey)
	if err != nil {
		log.Printf("Error fetching weather data for %s: %v", address.Localidade, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error fetching weather data")
		return
	}

	// Calculate temperatures
	tempC := weatherData.Current.TempC
	tempF := tempC*1.8 + 32
	tempK := tempC + 273

	// Log temperature data
	log.Printf("Temperature retrieved: %.1f°C for %s", tempC, address.Localidade)

	response := WeatherResponse{
		TempC: roundToOneDecimal(tempC),
		TempF: roundToOneDecimal(tempF),
		TempK: roundToOneDecimal(tempK),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func isValidCEP(cep string) bool {
	// Remove espaços e hífens
	cleaned := strings.ReplaceAll(strings.ReplaceAll(cep, "-", ""), " ", "")
	matched, _ := regexp.MatchString(`^\d{8}$`, cleaned)
	return matched
}

func roundToOneDecimal(value float64) float64 {
	rounded, _ := strconv.ParseFloat(strconv.FormatFloat(value, 'f', 1, 64), 64)
	return rounded
}

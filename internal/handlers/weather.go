package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

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

	// Validate CEP (8 numeric digits)
	if !isValidCEP(cep) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode("invalid zipcode")
		return
	}

	// Find city by CEP
	address, err := h.viaCEPClient.GetAddressByCEP(cep)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("can not find zipcode")
		return
	}

	// Check if API key is configured
	if h.config.WeatherAPIKey == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Weather API key not configured")
		return
	}

	// Get current weather
	weatherData, err := h.weatherClient.GetCurrentWeather(address.Localidade, h.config.WeatherAPIKey)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error fetching weather data")
		return
	}

	// Calculate temperatures
	tempC := weatherData.Current.TempC
	tempF := tempC*1.8 + 32
	tempK := tempC + 273

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
	// CEP must have exactly 8 numeric digits
	matched, _ := regexp.MatchString(`^\d{8}$`, cep)
	return matched
}

func roundToOneDecimal(value float64) float64 {
	rounded, _ := strconv.ParseFloat(strconv.FormatFloat(value, 'f', 1, 64), 64)
	return rounded
}

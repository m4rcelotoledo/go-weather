package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-weather/internal/clients/viacep"
	"go-weather/internal/clients/weather"
	"go-weather/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock clients
type MockViaCEPClient struct {
	mock.Mock
}

func (m *MockViaCEPClient) GetAddressByCEP(cep string) (*viacep.Response, error) {
	args := m.Called(cep)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*viacep.Response), args.Error(1)
}

type MockWeatherClient struct {
	mock.Mock
}

func (m *MockWeatherClient) GetCurrentWeather(city, apiKey string) (*weather.Response, error) {
	args := m.Called(city, apiKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*weather.Response), args.Error(1)
}

func TestWeatherHandler_GetWeatherByCEP_Success(t *testing.T) {
	// Setup
	mockViaCEP := new(MockViaCEPClient)
	mockWeather := new(MockWeatherClient)
	cfg := &config.Config{
		WeatherAPIKey: "test-api-key",
		Port:          "8080",
	}
	handler := NewWeatherHandler(mockViaCEP, mockWeather, cfg)

	// Mock responses
	viaCEPResponse := &viacep.Response{
		CEP:        "01001-000",
		Localidade: "São Paulo",
		Erro:       false,
	}
	weatherResponse := &weather.Response{
		Current: struct {
			TempC float64 `json:"temp_c"`
		}{TempC: 25.5},
	}

	mockViaCEP.On("GetAddressByCEP", "01001000").Return(viaCEPResponse, nil)
	mockWeather.On("GetCurrentWeather", "São Paulo", "test-api-key").Return(weatherResponse, nil)

	// Create request
	req := httptest.NewRequest("GET", "/01001000", nil)
	w := httptest.NewRecorder()

	// Create chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("cep", "01001000")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Execute
	handler.GetWeatherByCEP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response WeatherResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 25.5, response.TempC)
	assert.Equal(t, 77.9, response.TempF)
	assert.Equal(t, 298.5, response.TempK)

	mockViaCEP.AssertExpectations(t)
	mockWeather.AssertExpectations(t)
}

func TestWeatherHandler_GetWeatherByCEP_InvalidCEP(t *testing.T) {
	// Setup
	mockViaCEP := new(MockViaCEPClient)
	mockWeather := new(MockWeatherClient)
	cfg := &config.Config{
		WeatherAPIKey: "test-api-key",
		Port:          "8080",
	}
	handler := NewWeatherHandler(mockViaCEP, mockWeather, cfg)

	// Create request with invalid CEP
	req := httptest.NewRequest("GET", "/1234567", nil)
	w := httptest.NewRecorder()

	// Create chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("cep", "1234567")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Execute
	handler.GetWeatherByCEP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid zipcode", response)
}

func TestWeatherHandler_GetWeatherByCEP_CEPNotFound(t *testing.T) {
	// Setup
	mockViaCEP := new(MockViaCEPClient)
	mockWeather := new(MockWeatherClient)
	cfg := &config.Config{
		WeatherAPIKey: "test-api-key",
		Port:          "8080",
	}
	handler := NewWeatherHandler(mockViaCEP, mockWeather, cfg)

	// Mock error response
	mockViaCEP.On("GetAddressByCEP", "00000000").Return(nil, errors.New("CEP not found"))

	// Create request
	req := httptest.NewRequest("GET", "/00000000", nil)
	w := httptest.NewRecorder()

	// Create chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("cep", "00000000")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Execute
	handler.GetWeatherByCEP(w, req)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "can not find zipcode", response)

	mockViaCEP.AssertExpectations(t)
}

func TestWeatherHandler_GetWeatherByCEP_MissingAPIKey(t *testing.T) {
	// Setup
	mockViaCEP := new(MockViaCEPClient)
	mockWeather := new(MockWeatherClient)
	cfg := &config.Config{
		WeatherAPIKey: "", // Empty API key
		Port:          "8080",
	}
	handler := NewWeatherHandler(mockViaCEP, mockWeather, cfg)

	// Mock ViaCEP response
	viaCEPResponse := &viacep.Response{
		CEP:        "01001-000",
		Localidade: "São Paulo",
		Erro:       false,
	}
	mockViaCEP.On("GetAddressByCEP", "01001000").Return(viaCEPResponse, nil)

	// Create request
	req := httptest.NewRequest("GET", "/01001000", nil)
	w := httptest.NewRecorder()

	// Create chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("cep", "01001000")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Execute
	handler.GetWeatherByCEP(w, req)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Weather API key not configured", response)

	mockViaCEP.AssertExpectations(t)
}

func TestIsValidCEP(t *testing.T) {
	tests := []struct {
		cep      string
		expected bool
	}{
		{"01001000", true},
		{"12345678", true},
		{"1234567", false},   // 7 digits
		{"123456789", false}, // 9 digits
		{"abcd1234", false},  // contains letters
		{"", false},          // empty
		{"01001-000", true},  // contains dash (now supported)
	}

	for _, test := range tests {
		result := isValidCEP(test.cep)
		assert.Equal(t, test.expected, result, "CEP: %s", test.cep)
	}
}

func TestRoundToOneDecimal(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{25.56, 25.6},
		{25.54, 25.5},
		{25.0, 25.0},
		{25.99, 26.0},
	}

	for _, test := range tests {
		result := roundToOneDecimal(test.input)
		assert.Equal(t, test.expected, result, "Input: %f", test.input)
	}
}

// TestNewWeatherHandler tests the constructor
func TestNewWeatherHandler(t *testing.T) {
	// Create mock clients
	mockViaCEP := new(MockViaCEPClient)
	mockWeather := new(MockWeatherClient)
	cfg := &config.Config{
		WeatherAPIKey: "test-api-key",
		Port:          "8080",
	}

	// Create handler
	handler := NewWeatherHandler(mockViaCEP, mockWeather, cfg)

	// Assertions
	assert.NotNil(t, handler)
	assert.Equal(t, mockViaCEP, handler.viaCEPClient)
	assert.Equal(t, mockWeather, handler.weatherClient)
	assert.Equal(t, cfg, handler.config)
}

// TestWeatherHandler_GetWeatherByCEP_WeatherAPIError tests weather API error scenario
func TestWeatherHandler_GetWeatherByCEP_WeatherAPIError(t *testing.T) {
	// Setup
	mockViaCEP := new(MockViaCEPClient)
	mockWeather := new(MockWeatherClient)
	cfg := &config.Config{
		WeatherAPIKey: "test-api-key",
		Port:          "8080",
	}
	handler := NewWeatherHandler(mockViaCEP, mockWeather, cfg)

	// Mock responses
	viaCEPResponse := &viacep.Response{
		CEP:        "01001-000",
		Localidade: "São Paulo",
		Erro:       false,
	}

	mockViaCEP.On("GetAddressByCEP", "01001000").Return(viaCEPResponse, nil)
	mockWeather.On("GetCurrentWeather", "São Paulo", "test-api-key").Return(nil, errors.New("weather API error"))

	// Create request
	req := httptest.NewRequest("GET", "/01001000", nil)
	w := httptest.NewRecorder()

	// Create chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("cep", "01001000")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Execute
	handler.GetWeatherByCEP(w, req)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Error fetching weather data", response)

	mockViaCEP.AssertExpectations(t)
	mockWeather.AssertExpectations(t)
}

// TestIsValidCEP_EdgeCases tests edge cases for CEP validation
func TestIsValidCEP_EdgeCases(t *testing.T) {
	tests := []struct {
		cep      string
		expected bool
		desc     string
	}{
		{"01001-000", true, "CEP with dash"},
		{"01001 000", true, "CEP with space"},
		{"01001000", true, "CEP without formatting"},
		{"1234567", false, "7 digits"},
		{"123456789", false, "9 digits"},
		{"abcd1234", false, "contains letters"},
		{"", false, "empty string"},
		{"01001-00", false, "6 digits with dash"},
		{"01001-0000", false, "9 digits with dash"},
		{"01001-00a", false, "contains letter"},
		{"01001-00@", false, "contains special character"},
		{"01001-00.", false, "contains dot"},
		{"01001-00/", false, "contains slash"},
		{"01001-00\\", false, "contains backslash"},
		{"01001-00 ", false, "ends with space but only 6 digits"},
		{" 01001-000", true, "starts with space (removed by function)"},
		{"01001-000 ", true, "ends with space (removed by function)"},
		{" 01001-000 ", true, "starts and ends with space (removed by function)"},
	}

	for _, test := range tests {
		result := isValidCEP(test.cep)
		assert.Equal(t, test.expected, result, "CEP: %s - %s", test.cep, test.desc)
	}
}

// TestRoundToOneDecimal_EdgeCases tests edge cases for rounding
func TestRoundToOneDecimal_EdgeCases(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
		desc     string
	}{
		{0.0, 0.0, "zero"},
		{-25.56, -25.6, "negative number"},
		{25.55, 25.6, "round up"},
		{25.54, 25.5, "round down"},
		{25.50, 25.5, "exact half"},
		{25.51, 25.5, "just above half"},
		{25.49, 25.5, "just below half"},
		{999.99, 1000.0, "large number"},
		{-999.99, -1000.0, "large negative number"},
		{0.01, 0.0, "very small positive"},
		{-0.01, 0.0, "very small negative"},
	}

	for _, test := range tests {
		result := roundToOneDecimal(test.input)
		assert.Equal(t, test.expected, result, "Input: %f - %s", test.input, test.desc)
	}
}

// TestWeatherHandler_GetWeatherByCEP_DifferentCEPFormats tests CEP processing with different formats
func TestWeatherHandler_GetWeatherByCEP_DifferentCEPFormats(t *testing.T) {
	// Setup
	mockViaCEP := new(MockViaCEPClient)
	mockWeather := new(MockWeatherClient)
	cfg := &config.Config{
		WeatherAPIKey: "test-api-key",
		Port:          "8080",
	}
	handler := NewWeatherHandler(mockViaCEP, mockWeather, cfg)

	// Test different CEP formats
	cepFormats := []string{
		"01001000",  // Without dash
		"01001-000", // With dash
		"01001 000", // With space
	}

	for _, cepFormat := range cepFormats {
		// Mock responses
		viaCEPResponse := &viacep.Response{
			CEP:        "01001-000",
			Localidade: "São Paulo",
			Erro:       false,
		}
		weatherResponse := &weather.Response{
			Current: struct {
				TempC float64 `json:"temp_c"`
			}{TempC: 25.5},
		}

		mockViaCEP.On("GetAddressByCEP", cepFormat).Return(viaCEPResponse, nil)
		mockWeather.On("GetCurrentWeather", "São Paulo", "test-api-key").Return(weatherResponse, nil)

		// Create request
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		// Create chi context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("cep", cepFormat)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Execute
		handler.GetWeatherByCEP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code, "CEP format: %s", cepFormat)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response WeatherResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 25.5, response.TempC, "CEP format: %s", cepFormat)
		assert.Equal(t, 77.9, response.TempF, "CEP format: %s", cepFormat)
		assert.Equal(t, 298.5, response.TempK, "CEP format: %s", cepFormat)
	}

	mockViaCEP.AssertExpectations(t)
	mockWeather.AssertExpectations(t)
}

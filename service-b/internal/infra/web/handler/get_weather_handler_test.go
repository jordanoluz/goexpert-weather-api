package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidZipCode(t *testing.T) {
	tests := []struct {
		name     string
		zipcode  string
		expected bool
	}{
		{"valid zip code", "12345678", true},
		{"invalid zip code (letters)", "abcd1234", false},
		{"invalid zip code (short)", "1234", false},
		{"invalid zip code (empty)", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidZipCode(tt.zipcode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRemoveAccents(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"without accents", "CityName", "CityName"},
		{"with accents", "São Paulo", "Sao Paulo"},
		{"with mixed accents", "Curitiba, Paraná", "Curitiba, Parana"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeAccents(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRoundToTwoDecimal(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected float64
	}{
		{"rounds up", 1.236, 1.24},
		{"rounds down", 1.234, 1.23},
		{"exact two decimals", 1.23, 1.23},
		{"zero value", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := roundToTwoDecimal(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetWeatherHandler(t *testing.T) {
	tests := []struct {
		name           string
		zipcode        string
		expectedStatus int
		expectedBody   GetWeatherResponse
	}{
		{
			name:           "valid request",
			zipcode:        "93010001",
			expectedStatus: http.StatusOK,
			expectedBody:   GetWeatherResponse{},
		},
		{
			name:           "invalid zip code",
			zipcode:        "1234",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "non-existent city",
			zipcode:        "87654321",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/weather?zipcode="+tt.zipcode, nil)
			rec := httptest.NewRecorder()

			GetWeatherHandler(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK {
				var result GetWeatherResponse
				err := json.NewDecoder(rec.Body).Decode(&result)
				assert.NoError(t, err)
				assert.NotNil(t, tt.expectedBody, result)
			}
		})
	}
}

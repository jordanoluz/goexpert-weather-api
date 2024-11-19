package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type PostWeatherRequest struct {
	ZipCode string `json:"cep"`
}

type PostWeatherResponse struct {
	City                  string  `json:"city"`
	TemperatureCelcius    float64 `json:"temp_C"`
	TemperatureFahrenheit float64 `json:"temp_F"`
	TemperatureKelvin     float64 `json:"temp_K"`
}

func PostWeatherHandler(w http.ResponseWriter, r *http.Request) {
	var postWeatherRequest PostWeatherRequest
	if err := json.NewDecoder(r.Body).Decode(&postWeatherRequest); err != nil {
		http.Error(w, "failed to decode post weather request", http.StatusUnprocessableEntity)
		return
	}

	if !isValidZipCode(postWeatherRequest.ZipCode) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	weatherResponse, err := fetchWeatherApi(postWeatherRequest.ZipCode)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(weatherResponse); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func fetchWeatherApi(zipcode string) (*PostWeatherResponse, error) {
	url := fmt.Sprintf("https://goexpert-weather-api-1036920645078.southamerica-east1.run.app/weather?zipcode=%s", zipcode)

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch weather: %s", res.Status)
	}

	var weatherResponse PostWeatherResponse
	if err := json.NewDecoder(res.Body).Decode(&weatherResponse); err != nil {
		return nil, fmt.Errorf("failed to decode weather response: %w", err)
	}

	return &weatherResponse, nil
}

func isValidZipCode(zipcode string) bool {
	match, err := regexp.MatchString(`^\d{8}$`, zipcode)
	if err != nil {
		return false
	}

	return match
}

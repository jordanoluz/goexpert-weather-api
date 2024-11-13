package handler

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type WeatherResponse struct {
	TemperatureCelcius    float64 `json:"temp_C"`
	TemperatureFahrenheit float64 `json:"temp_F"`
	TemperatureKelvin     float64 `json:"temp_K"`
}

type ViaCepResponse struct {
	City string `json:"localidade"`
}

type WeatherApiResponse struct {
	Current WeatherApiCurrentResponse `json:"current"`
}

type WeatherApiCurrentResponse struct {
	TemperatureCelcius float64 `json:"temp_c"`
}

const WeatherApiKey = "28238955cc184dffb22235923241111"

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	zipcode := r.URL.Query().Get("zipcode")

	if !isValidZipCode(zipcode) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	city, err := fetchCity(zipcode)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	weatherApiResponse, err := fetchWeather(city)
	if err != nil {
		http.Error(w, "failed to fetch temperature", http.StatusInternalServerError)
		return
	}

	temperatureCelcius := weatherApiResponse.Current.TemperatureCelcius

	weatherResponse := WeatherResponse{
		TemperatureCelcius:    roundToTwoDecimal(temperatureCelcius),
		TemperatureFahrenheit: roundToTwoDecimal(temperatureCelcius*1.8 + 32),
		TemperatureKelvin:     roundToTwoDecimal(temperatureCelcius + 273),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(weatherResponse); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func isValidZipCode(zipcode string) bool {
	match, err := regexp.MatchString(`^\d{8}$`, zipcode)
	if err != nil {
		return false
	}

	return match
}

func fetchCity(zipcode string) (string, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", zipcode)

	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch city: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch city: %s", res.Status)
	}

	var viaCepResponse ViaCepResponse
	if err := json.NewDecoder(res.Body).Decode(&viaCepResponse); err != nil {
		return "", fmt.Errorf("failed to decode city response: %w", err)
	}

	if viaCepResponse.City == "" {
		return "", fmt.Errorf("city not found for zip code '%s'", zipcode)
	}

	return viaCepResponse.City, nil
}

func fetchWeather(city string) (*WeatherApiResponse, error) {
	city = removeAccents(city)
	city = strings.ReplaceAll(city, " ", "%20")

	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", WeatherApiKey, city)

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch weather: %s", res.Status)
	}

	var weatherApiResponse WeatherApiResponse
	if err := json.NewDecoder(res.Body).Decode(&weatherApiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode weather response: %w", err)
	}

	return &weatherApiResponse, nil
}

func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

func roundToTwoDecimal(value float64) float64 {
	return math.Round(value*100) / 100
}

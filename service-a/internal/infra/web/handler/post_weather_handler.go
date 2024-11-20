package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	internalOtel "github.com/jordanoluz/goexpert-weather-api/otel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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
	ctx, span := internalOtel.Tracer.Start(r.Context(), "post-weather-handler")
	defer span.End()

	var postWeatherRequest PostWeatherRequest
	if err := json.NewDecoder(r.Body).Decode(&postWeatherRequest); err != nil {
		http.Error(w, "failed to decode post weather request", http.StatusUnprocessableEntity)
		return
	}

	if !isValidZipCode(postWeatherRequest.ZipCode) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	weatherResponse, err := fetchWeatherApi(ctx, postWeatherRequest.ZipCode)
	if err != nil {
		log.Println(err)
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(weatherResponse); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func fetchWeatherApi(ctx context.Context, zipcode string) (*PostWeatherResponse, error) {
	ctx, span := internalOtel.Tracer.Start(ctx, "fetch-weather-api")
	defer span.End()

	url := fmt.Sprintf("%s/weather?zipcode=%s", os.Getenv("SERVICE_B_URL"), zipcode)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request with context: %w", err)
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	res, err := http.DefaultClient.Do(req)
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

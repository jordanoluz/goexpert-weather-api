package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jordanoluz/goexpert-weather-api/service-a/internal/infra/web/handler"
)

const ApiPort = 8181

func main() {
	http.HandleFunc("POST /weather", handler.PostWeatherHandler)

	log.Printf("listening and serving on port: %d", ApiPort)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", ApiPort), nil); err != nil {
		log.Fatalf("failed to listen and serve: %v", err)
	}
}

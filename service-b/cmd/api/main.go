package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jordanoluz/goexpert-weather-api/service-b/internal/infra/web/handler"
)

const ApiPort = 8080

func main() {
	http.HandleFunc("GET /weather", handler.GetWeatherHandler)

	log.Printf("listening and serving on port: %d", ApiPort)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", ApiPort), nil); err != nil {
		log.Fatalf("failed to listen and serve: %v", err)
	}
}

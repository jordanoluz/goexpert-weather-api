package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"

	"github.com/jordanoluz/goexpert-weather-api/otel"
	"github.com/jordanoluz/goexpert-weather-api/service-b/internal/infra/web/handler"
)

const ApiPort = 8080

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	tracerProvider, err := otel.InitTracerProvider()
	if err != nil {
		log.Fatalf("failed to initialize tracer provider: %v", err)
	}

	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown tracer provider: %v", err)
		}
	}()

	http.HandleFunc("GET /weather", handler.GetWeatherHandler)

	log.Printf("listening and serving on port: %d", ApiPort)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", ApiPort), nil); err != nil {
		log.Fatalf("failed to listen and serve: %v", err)
	}
}

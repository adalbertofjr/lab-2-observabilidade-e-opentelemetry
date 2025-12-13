package main

import (
	"fmt"
	"net/http"

	"github.com/adalbertofjr/lab-2-go-service-a-otel/internal/infra/api"
	"github.com/adalbertofjr/lab-2-go-service-a-otel/internal/infra/gateway"
	"github.com/adalbertofjr/lab-2-go-service-a-otel/internal/usecase/weather"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	startServer()
}

func startServer() {
	weatherGateway := gateway.NewWeatherAPI()
	weatherUseCase := weather.NewWeatherUseCase(weatherGateway)
	weatherHandler := api.NewWeatherHandler(weatherUseCase)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.HandleFunc("/", weatherHandler.GetCurrentWeather)

	fmt.Println("Starting web server on port", ":8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}

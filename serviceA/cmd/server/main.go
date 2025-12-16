package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/adalbertofjr/lab-2-go-service-a-otel/internal/infra/api"
	"github.com/adalbertofjr/lab-2-go-service-a-otel/internal/infra/gateway"
	"github.com/adalbertofjr/lab-2-go-service-a-otel/internal/usecase/weather"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := initProvider()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatalf("Failed to shutdown TracerProvider: %v", err)
		}
	}()

	tracer := otel.Tracer("serviceA-tracer")

	startServer(tracer)
}

func startServer(tracer trace.Tracer) {
	weatherGateway := gateway.NewWeatherAPI()
	weatherUseCase := weather.NewWeatherUseCase(weatherGateway, tracer)
	weatherHandler := api.NewWeatherHandler(weatherUseCase, tracer)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.HandleFunc("/", weatherHandler.GetCurrentWeather)

	fmt.Println("Starting web server on port", ":8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}

func initProvider() (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("ServiceA"),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to create Resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		"localhost:4317",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to create gRPC connection to collector: %w", err)
	}

	tracerExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("Failed to create the collector trace exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(tracerExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider.Shutdown, nil
}

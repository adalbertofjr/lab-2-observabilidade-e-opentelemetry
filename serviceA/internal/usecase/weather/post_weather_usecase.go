package weather

import (
	"context"

	"github.com/adalbertofjr/lab-2-go-service-a-otel/internal/domain/entity"
	domainGateway "github.com/adalbertofjr/lab-2-go-service-a-otel/internal/domain/gateway"
	"github.com/adalbertofjr/lab-2-go-service-a-otel/pkg/utility"
	"go.opentelemetry.io/otel/trace"
)

type WeatherUseCase struct {
	weatherGateway domainGateway.WeatherGateway
	tracer         trace.Tracer
}

func NewWeatherUseCase(gateway domainGateway.WeatherGateway, tracer trace.Tracer) *WeatherUseCase {
	return &WeatherUseCase{weatherGateway: gateway, tracer: tracer}
}

func (w *WeatherUseCase) GetCurrentWeather(ctx context.Context, cep string) (*entity.Weather, error) {
	// Span para validação do CEP
	ctx, spanValidateCep := w.tracer.Start(ctx, "validate_cep")
	cepFormated, err := utility.CEPFormatter(cep)
	spanValidateCep.End()
	if err != nil {
		return nil, err
	}

	// Span para chamada ao gateway
	ctx, spanGetWeather := w.tracer.Start(ctx, "call_service_b")
	weatherData, err := w.weatherGateway.GetCurrentWeather(cepFormated)
	spanGetWeather.End()
	if err != nil {
		return nil, err
	}

	currentWeather := entity.NewWeather(
		weatherData.City,
		weatherData.Temp_c,
		weatherData.Temp_f,
		weatherData.Temp_k,
	)

	return currentWeather, nil
}

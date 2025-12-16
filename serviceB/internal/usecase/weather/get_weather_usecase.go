package usecase

import (
	"context"

	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/domain/entity"
	domainGateway "github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/domain/gateway"
	internalerror "github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/infra/internal_error"

	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/pkg/utility"
	"go.opentelemetry.io/otel/trace"
)

type WeatherUseCase struct {
	weatherGateway domainGateway.WeatherGateway
	tracer         trace.Tracer
}

func NewWeatherUseCase(gateway domainGateway.WeatherGateway, tracer trace.Tracer) *WeatherUseCase {
	return &WeatherUseCase{weatherGateway: gateway, tracer: tracer}
}

func (w *WeatherUseCase) GetCurrentWeather(ctx context.Context, cep string) (*entity.Weather, *internalerror.InternalError) {
	// Span para validação do CEP
	ctx, spanValidateCep := w.tracer.Start(ctx, "validate_cep")
	cepFormated, err := utility.CEPFormatter(cep)
	if err != nil {
		return nil, internalerror.CEPInvalidError()
	}
	spanValidateCep.End()

	ctx, spanFetchWeatherData := w.tracer.Start(ctx, "fetch_weather_data")
	weatherData, err := w.weatherGateway.GetCurrentWeather(ctx, cepFormated)
	if err != nil {
		return nil, internalerror.CEPNotFoundError()
	}
	spanFetchWeatherData.End()

	currentWeather := entity.NewWeather(
		weatherData.City,
		weatherData.Temp_c)

	return currentWeather, nil
}

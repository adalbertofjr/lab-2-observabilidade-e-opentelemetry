package gateway

import (
	"context"

	"github.com/adalbertofjr/lab-2-go-service-a-otel/internal/domain/entity"
)

type WeatherGateway interface {
	GetCurrentWeather(ctx context.Context, cep string) (*entity.Weather, error)
}

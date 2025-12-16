package gateway

import (
	"context"

	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/domain/entity"
)

type WeatherGateway interface {
	GetCurrentWeather(ctx context.Context, cep string) (*entity.Weather, error)
}

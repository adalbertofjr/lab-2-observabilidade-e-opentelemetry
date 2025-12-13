package gateway

import (
	"github.com/adalbertofjr/lab-2-go-service-a-otel/internal/domain/entity"
)

type WeatherGateway interface {
	GetCurrentWeather(cep string) (*entity.Weather, error)
}

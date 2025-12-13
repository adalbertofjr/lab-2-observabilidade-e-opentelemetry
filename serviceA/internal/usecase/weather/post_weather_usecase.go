package weather

import (
	"github.com/adalbertofjr/lab-2-go-service-a-otel/internal/domain/entity"
	domainGateway "github.com/adalbertofjr/lab-2-go-service-a-otel/internal/domain/gateway"
	"github.com/adalbertofjr/lab-2-go-service-a-otel/pkg/utility"
)

type WeatherUseCase struct {
	weatherGateway domainGateway.WeatherGateway
}

func NewWeatherUseCase(gateway domainGateway.WeatherGateway) *WeatherUseCase {
	return &WeatherUseCase{weatherGateway: gateway}
}

func (w *WeatherUseCase) GetCurrentWeather(cep string) (*entity.Weather, error) {
	cepFormated, err := utility.CEPFormatter(cep)
	if err != nil {
		return nil, err
	}

	weatherData, err := w.weatherGateway.GetCurrentWeather(cepFormated)
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

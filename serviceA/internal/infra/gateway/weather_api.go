package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/adalbertofjr/lab-2-go-service-a-otel/internal/domain/entity"
)

type WeatherAPI struct {
}

type WeatherAPIResponse struct {
	City   string  `json:"city"`
	Temp_c float64 `json:"temp_c"`
	Temp_f float64 `json:"temp_f"`
	Temp_k float64 `json:"temp_k"`
}

func NewWeatherAPI() *WeatherAPI {
	return &WeatherAPI{}
}

func (w *WeatherAPI) GetCurrentWeather(cep string) (*entity.Weather, error) {
	url := fmt.Sprintf("http://localhost:8000/?cep=%s", url.QueryEscape(cep))
	client := http.Client{}
	resp, err := client.Post(url, "application/json", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get location data: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var weatherResponse WeatherAPIResponse
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		return nil, err
	}

	weatherData := entity.NewWeather(
		weatherResponse.City,
		weatherResponse.Temp_c,
		weatherResponse.Temp_f,
		weatherResponse.Temp_k,
	)

	return weatherData, nil
}

package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/adalbertofjr/lab-2-go-service-a-otel/internal/domain/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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

func (w *WeatherAPI) GetCurrentWeather(ctx context.Context, cep string) (*entity.Weather, error) {
	url := fmt.Sprintf("http://localhost:8000/?cep=%s", url.QueryEscape(cep))

	// Cria a requisição com contexto
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, err
	}

	// Injeta os headers de propagação de trace
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := http.Client{}
	resp, err := client.Do(req)
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

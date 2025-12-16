package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/domain/entity"
	"go.opentelemetry.io/otel/trace"
)

type WeatherAPI struct {
	APIKey string
	tracer trace.Tracer
}

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
}

type WeatherAPIResponse struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}

type Location struct {
	Name string `json:"name"`
}

type Current struct {
	Temp_c float64 `json:"temp_c"`
	Temp_f float64 `json:"temp_f"`
	Temp_k float64 `json:"temp_k"`
}

func NewWeatherAPI(apikey string, tracer trace.Tracer) *WeatherAPI {
	return &WeatherAPI{APIKey: apikey, tracer: tracer}
}

func (w *WeatherAPI) getLocation(ctx context.Context, cep string) (*ViaCEPResponse, error) {
	ctx, spanFetchCepLocation := w.tracer.Start(ctx, "fetch_cep_location")
	defer spanFetchCepLocation.End()

	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", url.QueryEscape(cep))
	client := http.Client{}
	resp, err := client.Get(url)
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

	var location ViaCEPResponse
	err = json.Unmarshal(body, &location)
	if err != nil {
		return nil, err
	}

	return &location, nil
}

func (w *WeatherAPI) GetCurrentWeather(ctx context.Context, cep string) (*entity.Weather, error) {
	location, err := w.getLocation(ctx, cep)
	if err != nil {
		return nil, err
	}

	ctx, spanFetchCurrentWeather := w.tracer.Start(ctx, "fetch_current_weather")
	defer spanFetchCurrentWeather.End()
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", w.APIKey, url.QueryEscape(location.Localidade))
	client := http.Client{}
	resp, err := client.Get(url)
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
		weatherResponse.Location.Name,
		weatherResponse.Current.Temp_c,
	)

	return weatherData, nil
}

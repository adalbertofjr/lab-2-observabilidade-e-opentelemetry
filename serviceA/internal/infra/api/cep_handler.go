package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/adalbertofjr/lab-2-go-service-a-otel/api/dto"
	"github.com/adalbertofjr/lab-2-go-service-a-otel/internal/domain/entity"
)

type WeatherUseCaseInterface interface {
	GetCurrentWeather(cep string) (*entity.Weather, error)
}

type WeatherHandler struct {
	usecase WeatherUseCaseInterface
}

func NewWeatherHandler(useCase WeatherUseCaseInterface) *WeatherHandler {
	return &WeatherHandler{usecase: useCase}
}

type CEP struct {
	CEP string `json:"cep"`
}

func (c *WeatherHandler) GetCurrentWeather(w http.ResponseWriter, r *http.Request) {
	req := r.Body
	defer req.Close()

	body, err := io.ReadAll(req)
	if err != nil {
		panic(err)
	}

	var cep CEP
	err = json.Unmarshal(body, &cep)
	if err != nil {
		panic(err)
	}

	currentWeather, err := c.usecase.GetCurrentWeather(cep.CEP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	weatherDTO := dto.NewWeatherDTO(
		currentWeather.City,
		currentWeather.Temp_c,
		currentWeather.Temp_f,
		currentWeather.Temp_k,
	)

	weatherCurrentJSON, jsonErr := json.Marshal(weatherDTO)
	if jsonErr != nil {
		http.Error(w, "Error marshalling location data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(weatherCurrentJSON))
}

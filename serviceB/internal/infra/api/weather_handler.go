package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/domain/entity"
	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/infra/api/dto"
	internalerror "github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/infra/internal_error"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type WeatherUseCaseInterface interface {
	GetCurrentWeather(ctx context.Context, cep string) (*entity.Weather, *internalerror.InternalError)
}

type WeatherHandler struct {
	usecase WeatherUseCaseInterface
	tracer  trace.Tracer
}

func NewWeatherHandler(useCase WeatherUseCaseInterface, tracer trace.Tracer) *WeatherHandler {
	return &WeatherHandler{usecase: useCase, tracer: tracer}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	cep := r.URL.Query().Get("cep")
	ctx, spanStart := h.tracer.Start(ctx, "Get /?cep="+cep)
	defer spanStart.End()

	weatherCurrent, err := h.usecase.GetCurrentWeather(ctx, cep)
	if err != nil {
		http.Error(w, err.MSG, err.Code)
		return
	}

	weatherDTO := dto.NewWeatherDTO(
		weatherCurrent.City,
		weatherCurrent.Temp_c,
		weatherCurrent.Temp_f,
		weatherCurrent.Temp_k,
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

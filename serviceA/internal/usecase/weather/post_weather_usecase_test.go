package weather

import (
	"errors"
	"testing"

	"github.com/adalbertofjr/lab-2-go-service-a-otel/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWeatherGateway is a mock implementation of the WeatherGateway interface
type MockWeatherGateway struct {
	mock.Mock
}

func (m *MockWeatherGateway) GetCurrentWeather(cep string) (*entity.Weather, error) {
	args := m.Called(cep)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Weather), args.Error(1)
}

func TestWeatherUseCase_GetCurrentWeather_Success(t *testing.T) {
	// Arrange
	mockGateway := new(MockWeatherGateway)
	usecase := NewWeatherUseCase(mockGateway)
	cep := "12345678"
	expectedWeather := &entity.Weather{
		City:   "Test City",
		Temp_c: 25.0,
		Temp_f: 77.0,
		Temp_k: 298.15,
	}

	mockGateway.On("GetCurrentWeather", "12345678").Return(expectedWeather, nil)

	// Act
	weather, err := usecase.GetCurrentWeather(cep)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, weather)
	assert.Equal(t, expectedWeather.City, weather.City)
	assert.Equal(t, expectedWeather.Temp_c, weather.Temp_c)
	mockGateway.AssertExpectations(t)
}

func TestWeatherUseCase_GetCurrentWeather_InvalidCEP(t *testing.T) {
	// Arrange
	mockGateway := new(MockWeatherGateway)
	usecase := NewWeatherUseCase(mockGateway)
	cep := "12345" // Invalid CEP

	// Act
	weather, err := usecase.GetCurrentWeather(cep)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, weather)
	assert.Equal(t, "invalid zipcode", err.Error())
	mockGateway.AssertNotCalled(t, "GetCurrentWeather", mock.Anything)
}

func TestWeatherUseCase_GetCurrentWeather_GatewayError(t *testing.T) {
	// Arrange
	mockGateway := new(MockWeatherGateway)
	usecase := NewWeatherUseCase(mockGateway)
	cep := "87654321"
	gatewayError := errors.New("gateway failed")

	mockGateway.On("GetCurrentWeather", "87654321").Return(nil, gatewayError)

	// Act
	weather, err := usecase.GetCurrentWeather(cep)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, weather)
	assert.Equal(t, gatewayError, err)
	mockGateway.AssertExpectations(t)
}

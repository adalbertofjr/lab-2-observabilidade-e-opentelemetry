package entity

type Weather struct {
	City   string
	Temp_c float64
	Temp_f float64
	Temp_k float64
}

func NewWeather(city string, tempC float64, tempF float64, tempK float64) *Weather {
	weather := Weather{
		City:   city,
		Temp_c: tempC,
		Temp_f: tempF,
		Temp_k: tempK,
	}

	return &weather
}

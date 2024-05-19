package repo

import (
	"log"
	"os"
	"testing"

	"github.com/aleroxac/goexpert-weather-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestGetWeather(t *testing.T) {
	repo := NewWeatherRepository()

	via_cep_api_key := os.Getenv("VIA_CEP_API_KEY")
	if via_cep_api_key == "" {
		log.Fatal("Please provide the environment variable VIA_CEP_API_KEY and try again.")
	}

	weather_res, err := repo.Get("São Paulo", via_cep_api_key)
	assert.NoError(t, err)
	assert.Contains(t, string(weather_res), "São Paulo")
	assert.Contains(t, string(weather_res), "temp")
	assert.Contains(t, string(weather_res), "temp_min")
	assert.Contains(t, string(weather_res), "temp_max")
}

func TestConvertToWeatherResponse(t *testing.T) {
	weather_json := []byte(
		`{
			"coord": {
				"lon": -46.6361,
				"lat": -23.5475
			},
			"weather": [
				{
					"id": 803,
					"main": "Clouds",
					"description": "broken clouds",
					"icon": "04d"
				}
			],
			"base": "stations",
			"main": {
				"temp": 21.1,
				"feels_like": 21.35,
				"temp_min": 19.75,
				"temp_max": 24.14,
				"pressure": 1024,
				"humidity": 80
			},
			"visibility": 10000,
			"wind": {
				"speed": 3.6,
				"deg": 140
			},
			"clouds": {
				"all": 75
			},
			"dt": 1716126286,
			"sys": {
				"type": 1,
				"id": 8394,
				"country": "BR",
				"sunrise": 1716111334,
				"sunset": 1716150647
			},
			"timezone": -10800,
			"id": 3448439,
			"name": "São Paulo",
			"cod": 200
		}`,
	)
	repo := NewWeatherRepository()
	weather_res, err := repo.ConvertToWeatherResponse(weather_json)
	assert.NoError(t, err)
	assert.IsType(t, weather_res, &entity.WeatherResponse{})
}

func TestConvertToWeather(t *testing.T) {
	weather_json := []byte(
		`{
			"coord": {
				"lon": -46.6361,
				"lat": -23.5475
			},
			"weather": [
				{
					"id": 803,
					"main": "Clouds",
					"description": "broken clouds",
					"icon": "04d"
				}
			],
			"base": "stations",
			"main": {
				"temp": 21.1,
				"feels_like": 21.35,
				"temp_min": 19.75,
				"temp_max": 24.14,
				"pressure": 1024,
				"humidity": 80
			},
			"visibility": 10000,
			"wind": {
				"speed": 3.6,
				"deg": 140
			},
			"clouds": {
				"all": 75
			},
			"dt": 1716126286,
			"sys": {
				"type": 1,
				"id": 8394,
				"country": "BR",
				"sunrise": 1716111334,
				"sunset": 1716150647
			},
			"timezone": -10800,
			"id": 3448439,
			"name": "São Paulo",
			"cod": 200
		}`,
	)
	repo := NewWeatherRepository()
	weather_res, _ := repo.ConvertToWeatherResponse(weather_json)

	weather, err := repo.ConvertToWeather(weather_res)
	assert.NoError(t, err)
	assert.IsType(t, weather, &entity.Weather{})
	assert.Equal(t, weather.Fahrenheit, weather.Celcius*1.8+32)
	assert.Equal(t, weather.Kelvin, weather.Celcius+273.15)
}

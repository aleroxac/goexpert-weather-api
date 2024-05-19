package usecase_test

import (
	"log"
	"os"
	"testing"

	"github.com/aleroxac/goexpert-weather-api/internal/infra/web"
	"github.com/aleroxac/goexpert-weather-api/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestGetWeather(t *testing.T) {
	via_cep_api_key := os.Getenv("VIA_CEP_API_KEY")
	if via_cep_api_key == "" {
		log.Fatal("Please provide the environment variable VIA_CEP_API_KEY and try again.")
	}

	t.Run("valid weather", func(t *testing.T) {
		get_weather_dto := usecase.WeatherInputDTO{
			Localidade: "São Paulo",
			ApiKey:     via_cep_api_key,
		}
		webCEPHandler := web.NewWebCEPHandler()

		getWeather := usecase.NewGetWeatherUseCase(webCEPHandler.WeatherRepository)
		weather_output, err := getWeather.Execute(get_weather_dto)
		assert.NoError(t, err)
		assert.IsType(t, &weather_output, &usecase.WeatherOutputDTO{})
	})

	t.Run("missing Localidade field", func(t *testing.T) {
		get_weather_dto := usecase.WeatherInputDTO{
			ApiKey: via_cep_api_key,
		}
		webCEPHandler := web.NewWebCEPHandler()

		getWeather := usecase.NewGetWeatherUseCase(webCEPHandler.WeatherRepository)
		weather_output, err := getWeather.Execute(get_weather_dto)
		assert.EqualError(t, err, "missing input: Localidade")
		assert.IsType(t, &weather_output, &usecase.WeatherOutputDTO{})
	})

	t.Run("missing ApiKey field", func(t *testing.T) {
		get_weather_dto := usecase.WeatherInputDTO{
			Localidade: "São Paulo",
		}
		webCEPHandler := web.NewWebCEPHandler()

		getWeather := usecase.NewGetWeatherUseCase(webCEPHandler.WeatherRepository)
		weather_output, err := getWeather.Execute(get_weather_dto)
		assert.EqualError(t, err, "missing input: ApiKey")
		assert.IsType(t, &weather_output, &usecase.WeatherOutputDTO{})
	})

	t.Run("fail to get weather", func(t *testing.T) {
		get_weather_dto := usecase.WeatherInputDTO{
			Localidade: "goexpert",
			ApiKey:     via_cep_api_key,
		}
		webCEPHandler := web.NewWebCEPHandler()

		getWeather := usecase.NewGetWeatherUseCase(webCEPHandler.WeatherRepository)
		weather_output, err := getWeather.Execute(get_weather_dto)
		assert.EqualError(t, err, "fail to get weather")
		assert.Equal(t, weather_output.Celcius, float64(0))
		assert.Equal(t, weather_output.Fahrenheit, float64(0))
		assert.Equal(t, weather_output.Kelvin, float64(0))
	})
}

package usecase_test

import (
	"testing"

	"github.com/aleroxac/goexpert-weather-api/internal/infra/web"
	"github.com/aleroxac/goexpert-weather-api/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestValidateCEP(t *testing.T) {
	get_weather_dto := usecase.ValidateCEPInputDTO{
		CEP: "01001001",
	}
	webCEPHandler := web.NewWebCEPHandler()

	validate_cep := usecase.NewValidateCEPUseCase(webCEPHandler.CEPRepository)
	weather_output := validate_cep.Execute(get_weather_dto)
	assert.True(t, weather_output)
}

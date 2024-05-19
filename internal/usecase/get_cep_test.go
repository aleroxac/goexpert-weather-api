package usecase_test

import (
	"testing"

	"github.com/aleroxac/goexpert-weather-api/internal/infra/web"
	"github.com/aleroxac/goexpert-weather-api/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestGetCEPUseCase(t *testing.T) {
	t.Run("valid cep", func(t *testing.T) {
		cep_address := "01001001"
		webCEPHandler := web.NewWebCEPHandler()

		get_cep_dto := usecase.CEPInputDTO{
			CEP: cep_address,
		}
		getCEP := usecase.NewGetCEPUseCase(webCEPHandler.CEPRepository)
		cep_output, err := getCEP.Execute(get_cep_dto)
		assert.NoError(t, err)
		assert.Equal(t, cep_output.CEP, "01001-001")
		assert.Equal(t, cep_output.Logradouro, "Praça da Sé")
		assert.Equal(t, cep_output.Complemento, "lado par")
		assert.Equal(t, cep_output.Bairro, "Sé")
		assert.Equal(t, cep_output.Localidade, "São Paulo")
		assert.Equal(t, cep_output.UF, "SP")
		assert.Equal(t, cep_output.IBGE, "3550308")
		assert.Equal(t, cep_output.GIA, "1004")
		assert.Equal(t, cep_output.DDD, "11")
		assert.Equal(t, cep_output.SIAFI, "7107")
	})

	t.Run("invalid cep", func(t *testing.T) {
		cep_address := "0100100"
		webCEPHandler := web.NewWebCEPHandler()

		get_cep_dto := usecase.CEPInputDTO{
			CEP: cep_address,
		}
		getCEP := usecase.NewGetCEPUseCase(webCEPHandler.CEPRepository)
		cep_output, err := getCEP.Execute(get_cep_dto)
		assert.Error(t, err)
		assert.Equal(t, cep_output.CEP, "")
		assert.Equal(t, cep_output.Logradouro, "")
		assert.Equal(t, cep_output.Complemento, "")
		assert.Equal(t, cep_output.Bairro, "")
		assert.Equal(t, cep_output.Localidade, "")
		assert.Equal(t, cep_output.UF, "")
		assert.Equal(t, cep_output.IBGE, "")
		assert.Equal(t, cep_output.GIA, "")
		assert.Equal(t, cep_output.DDD, "")
		assert.Equal(t, cep_output.SIAFI, "")
	})
}

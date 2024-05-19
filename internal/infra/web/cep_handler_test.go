package web

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aleroxac/goexpert-weather-api/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestCEPHandler(t *testing.T) {
	via_cep_api_key := os.Getenv("VIA_CEP_API_KEY")
	if via_cep_api_key == "" {
		log.Fatal("Please provide the environment variable VIA_CEP_API_KEY and try again.")
	}

	router := chi.NewRouter()
	router.Get("/cep/{cep}", NewWebCEPHandler().Get)

	req, err := http.NewRequest("GET", "/cep/01001001", nil)
	req.Header.Add("VIA_CEP_API_KEY", via_cep_api_key)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 200)

	var weather *usecase.WeatherOutputDTO
	err = json.Unmarshal(rr.Body.Bytes(), &weather)
	assert.NoError(t, err)

	assert.IsType(t, weather, &usecase.WeatherOutputDTO{})
	assert.NotEmpty(t, weather.Celcius)
	assert.NotEmpty(t, weather.Fahrenheit)
	assert.NotEmpty(t, weather.Kelvin)
}

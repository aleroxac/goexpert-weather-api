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
	open_weathermap_api_key := os.Getenv("OPEN_WEATHERMAP_API_KEY")
	if open_weathermap_api_key == "" {
		log.Fatal("Please provide the environment variable OPEN_WEATHERMAP_API_KEY and try again.")
	}

	router := chi.NewRouter()
	router.Get("/cep/{cep}", NewWebCEPHandler().Get)

	req, err := http.NewRequest("GET", "/cep/01001001", nil)
	req.Header.Add("OPEN_WEATHERMAP_API_KEY", open_weathermap_api_key)
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

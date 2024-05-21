package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aleroxac/goexpert-weather-api/internal/entity"
	"github.com/aleroxac/goexpert-weather-api/internal/infra/repo"
	"github.com/aleroxac/goexpert-weather-api/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type WebCEPHandler struct {
	CEPRepository     entity.CEPRepositoryInterface
	WeatherRepository entity.WeatherRepositoryInterface
}

func NewWebCEPHandler() *WebCEPHandler {
	return &WebCEPHandler{
		CEPRepository:     repo.NewCEPRepository(),
		WeatherRepository: repo.NewWeatherRepository(),
	}
}

func (h *WebCEPHandler) Get(w http.ResponseWriter, r *http.Request) {
	cep_address := chi.URLParam(r, "cep")

	open_weathermap_api_key := r.Header.Get("OPEN_WEATHERMAP_API_KEY")
	if open_weathermap_api_key == "" {
		http.Error(w, "Please, provide the OPEN_WEATHERMAP_API_KEY header", http.StatusBadRequest)
		log.Println("Please, provide the OPEN_WEATHERMAP_API_KEY header")
		return
	}

	// CEP FLOW
	validate_cep_dto := usecase.ValidateCEPInputDTO{
		CEP: cep_address,
	}

	validateCEP := usecase.NewValidateCEPUseCase(h.CEPRepository)
	is_valid := validateCEP.Execute(validate_cep_dto)
	if !is_valid {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	get_cep_dto := usecase.CEPInputDTO{
		CEP: cep_address,
	}

	getCEP := usecase.NewGetCEPUseCase(h.CEPRepository)
	cep_output, err := getCEP.Execute(get_cep_dto)

	if err != nil {
		http.Error(w, "error getting cep", http.StatusInternalServerError)
		return
	}

	if cep_output.Localidade == "" {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	// WEATHER FLOW
	weather_dto := usecase.WeatherInputDTO{
		Localidade: cep_output.Localidade,
		ApiKey:     open_weathermap_api_key,
	}
	getWeather := usecase.NewGetWeatherUseCase(h.WeatherRepository)
	weather_output, err := getWeather.Execute(weather_dto)
	if err != nil || (weather_output.Celcius == 0 && weather_output.Fahrenheit == 0 && weather_output.Kelvin == 0) {
		http.Error(w, "error getting weather", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(weather_output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

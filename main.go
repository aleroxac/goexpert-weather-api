package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

const (
	TIMEOUT_API = 2 * time.Second
	API_PORT    = 8080
)

var (
	WEATHER_API_KEY = os.Getenv("WEATHER_API_KEY")
)

type CEP struct {
	CEP string
}

type Location struct {
	CEP         string
	Logradouro  string
	Complemento string
	Bairro      string
	Localidade  string
	UF          string
	IBGE        string
	GIA         string
	DDD         string
	SIAFI       string
}

type Weather struct {
	Celcius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

type WeatherDetails struct {
	Temp float64
}

type WeatherResponse struct {
	Main WeatherDetails
}

func New(celcius, fahrenheit, kelvin float64) *Weather {
	return &Weather{
		Celcius:    celcius,
		Fahrenheit: fahrenheit,
		Kelvin:     kelvin,
	}
}

func (c *CEP) Validate() bool {
	check, _ := regexp.MatchString("^[0-9]{8}$", c.CEP)
	return (len(c.CEP) == 8 && c.CEP != "" && check)
}

func (c *CEP) GetCEP() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_API)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://viacep.com.br/ws/%s/json", c.CEP), nil)
	if err != nil {
		log.Fatalf("Fail to create the request: %v", err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Fail to make the request: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	ctx_err := ctx.Err()
	if ctx_err != nil {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			log.Fatalf("Max timeout reached: %v", err)
			return nil, err
		}
	}

	resp_json, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Fail to read the response: %v", err)
		return nil, err
	}

	return resp_json, nil
}

func (c *CEP) ConvertCEPToLocation(cep_response []byte) (*Location, error) {
	var location Location
	err := json.Unmarshal(cep_response, &location)
	if err != nil {
		log.Fatalf("Fail to decode the response: %v", err)
		return nil, err
	}
	return &location, nil
}

func (l *Location) GetWeather(api_key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_API)
	defer cancel()
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf(
			"http://api.openweathermap.org/data/2.5/weather?q=%v&appid=%v&units=metric&temperature.unit=celsius",
			l.Localidade,
			api_key,
		),
		nil,
	)
	if err != nil {
		log.Fatalf("Fail to create the request: %v", err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Fail to make the request: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	ctx_err := ctx.Err()
	if ctx_err != nil {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			log.Fatalf("Max timeout reached: %v", err)
			return nil, err
		}
	}

	resp_json, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Fail to read the response: %v", err)
		return nil, err
	}

	return resp_json, nil
}

func (w *Weather) ConvertWeatherDetails(weather_response []byte) (*WeatherResponse, error) {
	var weather WeatherResponse
	err := json.Unmarshal(weather_response, &weather)
	if err != nil {
		log.Fatalf("Fail to decode the response: %v", err)
		return nil, err
	}
	return &weather, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	cep := CEP{
		CEP: r.FormValue("cep"),
	}
	if !cep.Validate() {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}
	cep_response, err := cep.GetCEP()
	if err != nil {
		http.Error(w, "error getting cep", http.StatusInternalServerError)
		return
	}
	cep_info, err := cep.ConvertCEPToLocation(cep_response)
	if err != nil {
		http.Error(w, "error converting cep", http.StatusInternalServerError)
		return
	}
	if cep_info.Localidade == "" {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	weather_res, err := cep_info.GetWeather(WEATHER_API_KEY)
	if err != nil {
		http.Error(w, "error getting weather", http.StatusInternalServerError)
		return
	}
	var weather Weather
	weather_celcius, err := weather.ConvertWeatherDetails(weather_res)
	if err != nil {
		http.Error(w, "error converting weather", http.StatusInternalServerError)
		return
	}
	weather.Celcius, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", weather_celcius.Main.Temp), 2)
	weather.Fahrenheit, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", weather_celcius.Main.Temp*1.8+32), 2)
	weather.Kelvin, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", weather_celcius.Main.Temp+273.15), 2)

	w.Header().Set("Content-Type", "application/json")
	json_resp, err := json.Marshal(weather)
	if err != nil {
		log.Fatalf("Fail to encode response: %v", err)
		http.Error(w, fmt.Sprintf("Fail to encode response: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(json_resp)
}

func main() {
	if WEATHER_API_KEY == "" {
		log.Fatalf("WEATHER_API_KEY is empty")
		return
	}

	http.HandleFunc("/weather", handler)
	log.Printf("Listening on %d", API_PORT)
	http.ListenAndServe(":8080", nil)
}

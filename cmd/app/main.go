package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aleroxac/goexpert-weather-api/internal/infra/repo"
	"github.com/aleroxac/goexpert-weather-api/internal/infra/web"
	"github.com/aleroxac/goexpert-weather-api/internal/infra/web/webserver"
)

func ConfigureServer() *webserver.WebServer {
	webserver := webserver.NewWebServer(":8080")

	cepRepo := repo.NewCEPRepository()
	weatherRepo := repo.NewWeatherRepository(&http.Client{})

	webCEPHandler := web.NewWebCEPHandlerWithDeps(cepRepo, weatherRepo, os.Getenv("OPEN_WEATHERMAP_API_KEY"))
	webStatusHandler := web.NewWebStatusHandler()

	webserver.AddHandler("GET /cep/{cep}", webCEPHandler.Get)
	webserver.AddHandler("GET /status", webStatusHandler.Get)

	return webserver
}

func main() {
	webserver := ConfigureServer()
	fmt.Println("Starting web server on port", ":8080")
	webserver.Start()
}

package main

import (
	"fmt"

	"github.com/aleroxac/goexpert-weather-api/internal/infra/web"
	"github.com/aleroxac/goexpert-weather-api/internal/infra/web/webserver"
)

func ConfigureServer() *webserver.WebServer {
	webserver := webserver.NewWebServer(":8080")

	webCEPHandler := web.NewWebCEPHandler()
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

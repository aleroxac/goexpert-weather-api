package main

import (
	"fmt"

	"github.com/aleroxac/goexpert-weather-api/internal/infra/web"
	"github.com/aleroxac/goexpert-weather-api/internal/infra/web/webserver"
)

func main() {
	// ----- WEBSERVER
	webserver := webserver.NewWebServer(":8080")
	webCEPHandler := web.NewWebCEPHandler()
	webserver.AddHandler("GET /cep/{cep}", webCEPHandler.Get)
	fmt.Println("Starting web server on port", ":8080")
	webserver.Start()
}

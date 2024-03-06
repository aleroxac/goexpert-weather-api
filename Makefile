conf ?= .env
include $(conf)
export $(shell sed 's/=.*//' $(conf))



## ---------- UTILS
.PHONY: help
help: ## Show this menu
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: clean
clean: ## Clean all temp files
	@rm -f coverage.*



## ----- MAIN
build: ## build the container image
	@docker build -t gcr.io/aleroxac/goexpert-cloudrun:v1 .

push: ## push the container image to image registry
	@docker push gcr.io/aleroxac/goexpert-cloudrun:v1

deploy: push ## deploy the application to cloud-run
	@gcloud run deploy goexpert-weather-api

serve: ## run the server
	@go run main.go

run: ## make some api requests
	@echo -e "\n-------------------- 422 --------------------"; curl -s "http://localhost:8080/weather?cep=1234567"
	@echo -e "\n-------------------- 402 --------------------"; curl -s "http://localhost:8080/weather?cep=12345678"
	@echo -e "\n-------------------- 200 --------------------"; curl -s "http://localhost:8080/weather?cep=13330250"

test: ## run tests
	@go test -v ./... -coverprofile=coverage.out
	@go tool cover -html coverage.out -o coverage.html

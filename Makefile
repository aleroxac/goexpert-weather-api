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
	@docker build -t gcr.io/aleroxac/goexpert-weather-api:v1 .

push: ## push the container image to image registry
	@docker push gcr.io/aleroxac/goexpert-weather-api:v1

serve: ## run the server
	@cd cmd/app && go run main.go
	@cd -

run: ## make some api requests
	@echo -e "\n-------------------- 422 --------------------"; curl -s "http://localhost:8080/cep/1234567" -H "VIA_CEP_API_KEY: ${VIA_CEP_API_KEY}"
	@echo -e "\n-------------------- 402 --------------------"; curl -s "http://localhost:8080/cep/12345678" -H "VIA_CEP_API_KEY: ${VIA_CEP_API_KEY}"
	@echo -e "\n-------------------- 200 --------------------"; curl -s "http://localhost:8080/cep/13330250" -H "VIA_CEP_API_KEY: ${VIA_CEP_API_KEY}"

test: ## run tests
	@go test -v ./... -coverprofile=coverage.out
	@go tool cover -html coverage.out -o coverage.html

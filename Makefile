.PHONY: build run docker-up docker-down db-up test tidy help

# Variables
APP_NAME=server
CMD_PATH=cmd/server/main.go
CONFIG_PATH=config/local.yml

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	go build -o bin/$(APP_NAME) $(CMD_PATH)

run: ## Run the application locally
	CONFIG_PATH=$(CONFIG_PATH) go run $(CMD_PATH)

docker-up: ## Start all services with docker compose
	docker compose up -d --build

docker-down: ## Stop all services
	docker compose down

db-up: ## Start only the database service
	docker compose up -d db

test: ## Run all tests
	go test -v ./...

tidy: ## Run go mod tidy
	go mod tidy

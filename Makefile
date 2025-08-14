# Vyve Backend Makefile
.PHONY: help dev prod test migrate seed clean docker-build docker-push deploy logs

# Variables
DOCKER_REGISTRY ?= 
IMAGE_NAME ?= vyve-api
VERSION ?= latest
DB_URL ?= postgres://vyve:vyve@localhost:5432/vyve_dev?sslmode=disable
MIGRATE_PATH = ./migrations

# Colors for output
GREEN := \033[0;32m
RED := \033[0;31m
NC := \033[0m # No Color

help: ## Show this help message
	@echo "Vyve Backend - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""

# Development Commands
dev: ## Start development environment with hot reload
	@echo "$(GREEN)Starting development environment...$(NC)"
	docker compose -f docker-compose.dev.yml up --build

dev-down: ## Stop development environment
	@echo "$(RED)Stopping development environment...$(NC)"
	docker compose -f docker-compose.dev.yml down

dev-logs: ## Show development logs
	docker compose -f docker-compose.dev.yml logs -f

dev-clean: ## Clean development environment (removes volumes)
	@echo "$(RED)Cleaning development environment...$(NC)"
	docker compose -f docker-compose.dev.yml down -v
	rm -rf ./data

# Production Commands
prod: ## Build production Docker image
	@echo "$(GREEN)Building production image...$(NC)"
	docker build -t $(IMAGE_NAME):$(VERSION) -f Dockerfile.prod .

prod-up: ## Start production environment locally
	@echo "$(GREEN)Starting production environment...$(NC)"
	docker compose -f docker-compose.prod.yml up -d

prod-down: ## Stop production environment
	@echo "$(RED)Stopping production environment...$(NC)"
	docker compose -f docker-compose.prod.yml down

# Testing
test: ## Run all tests
	@echo "$(GREEN)Running tests...$(NC)"
	go test -v -cover ./...

test-integration: ## Run integration tests
	@echo "$(GREEN)Running integration tests...$(NC)"
	go test -v -tags=integration ./...

test-coverage: ## Run tests with coverage report
	@echo "$(GREEN)Generating coverage report...$(NC)"
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Database Operations
migrate-up: ## Run database migrations up
	@echo "$(GREEN)Running migrations up...$(NC)"
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up

migrate-down: ## Run database migrations down (1 step)
	@echo "$(RED)Rolling back migration...$(NC)"
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down 1

migrate-create: ## Create new migration (usage: make migrate-create name=create_users_table)
	@echo "$(GREEN)Creating migration: $(name)$(NC)"
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

migrate-force: ## Force migration version (usage: make migrate-force version=1)
	@echo "$(RED)Forcing migration version to $(version)...$(NC)"
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" force $(version)

seed: ## Seed database with sample data
	@echo "$(GREEN)Seeding database...$(NC)"
	go run cmd/seed/main.go

# Build & Deployment
build: ## Build Go binary
	@echo "$(GREEN)Building binary...$(NC)"
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api cmd/api/main.go

docker-build: ## Build Docker image
	@echo "$(GREEN)Building Docker image...$(NC)"
	docker build -t $(IMAGE_NAME):$(VERSION) .

docker-push: ## Push Docker image to registry
	@echo "$(GREEN)Pushing Docker image to registry...$(NC)"
	docker tag $(IMAGE_NAME):$(VERSION) $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(VERSION)
	docker push $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(VERSION)

deploy-ec2: ## Deploy to EC2 instance
	@echo "$(GREEN)Deploying to EC2...$(NC)"
	./scripts/deploy-ec2.sh

deploy-ecs: ## Deploy to ECS
	@echo "$(GREEN)Deploying to ECS...$(NC)"
	./scripts/deploy-ecs.sh

# Code Quality
lint: ## Run linter
	@echo "$(GREEN)Running linter...$(NC)"
	golangci-lint run ./...

fmt: ## Format code
	@echo "$(GREEN)Formatting code...$(NC)"
	go fmt ./...
	gofmt -s -w .

vet: ## Run go vet
	@echo "$(GREEN)Running go vet...$(NC)"
	go vet ./...

mod-tidy: ## Tidy go modules
	@echo "$(GREEN)Tidying modules...$(NC)"
	go mod tidy

# Development Tools
install-tools: ## Install development tools
	@echo "$(GREEN)Installing development tools...$(NC)"
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/air-verse/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

swagger: ## Generate Swagger documentation
	@echo "$(GREEN)Generating Swagger docs...$(NC)"
	swag init -g cmd/api/main.go -o ./docs

# Monitoring
logs: ## Show production logs
	docker compose -f docker-compose.prod.yml logs -f api

logs-tail: ## Show last 100 lines of logs
	docker compose -f docker-compose.prod.yml logs --tail=100 api

health: ## Check API health
	@echo "$(GREEN)Checking API health...$(NC)"
	curl -f http://localhost:8080/health || echo "$(RED)API is not healthy$(NC)"

# Utility
clean: ## Clean build artifacts and temporary files
	@echo "$(RED)Cleaning build artifacts...$(NC)"
	rm -rf bin/ tmp/ coverage.* *.out
	go clean -cache

env-example: ## Copy .env.example to .env
	@echo "$(GREEN)Creating .env file...$(NC)"
	cp .env.example .env

psql: ## Connect to development PostgreSQL
	docker compose -f docker-compose.dev.yml exec postgres psql -U vyve -d vyve_dev

redis-cli: ## Connect to development Redis
	docker compose -f docker-compose.dev.yml exec redis redis-cli

backup-db: ## Backup database
	@echo "$(GREEN)Backing up database...$(NC)"
	docker compose -f docker-compose.dev.yml exec postgres pg_dump -U vyve vyve_dev > backup_$(shell date +%Y%m%d_%H%M%S).sql

# Default target
.DEFAULT_GOAL := help
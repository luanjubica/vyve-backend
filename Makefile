.PHONY: build run test clean docker-build docker-run setup dev lint format

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=vyve-backend

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v cmd/server/main.go

# Run the application
run:
	$(GOCMD) run cmd/server/main.go

# Run with hot reload (requires air)
dev:
	air -c .air.toml

# Test the application
test:
	$(GOTEST) -v ./...

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Download dependencies
deps:
	$(GOCMD) mod download
	$(GOCMD) mod tidy

# Setup development environment
setup:
	./scripts/setup.sh

# Docker commands
docker-build:
	docker build -t vyve-backend .

docker-run:
	docker-compose up --build

docker-stop:
	docker-compose down

# Format code
format:
	$(GOCMD) fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Security check (requires gosec)
security:
	gosec ./...

# Run all checks
check: format lint test

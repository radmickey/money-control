.PHONY: all build run test clean proto docker-build docker-up docker-down

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOFMT=gofmt

# Protobuf
PROTOC=protoc

# Docker
DOCKER_COMPOSE=docker compose

# Service names
SERVICES=auth accounts transactions assets currency insights gateway

all: help

# Build all services
build:
	@echo "Building all services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd backend && $(GOBUILD) -o bin/$$service ./services/$$service; \
	done

# Run a specific service (usage: make run SERVICE=auth)
run:
	@echo "Running $(SERVICE) service..."
	cd backend && $(GOCMD) run ./services/$(SERVICE)

# Run the gateway
run-gateway:
	@echo "Running API Gateway..."
	cd backend && $(GOCMD) run ./services/gateway

# Run tests
test:
	@echo "Running tests..."
	cd backend && $(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	cd backend && $(GOTEST) -v -coverprofile=coverage.out ./...
	cd backend && $(GOCMD) tool cover -html=coverage.out -o coverage.html

# Generate protobuf files
proto:
	@echo "Generating protobuf files..."
	cd backend && $(PROTOC) --go_out=. --go-grpc_out=. proto/*.proto

# Format code
fmt:
	@echo "Formatting code..."
	cd backend && $(GOFMT) -s -w .

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	cd backend && $(GOMOD) download
	cd backend && $(GOMOD) tidy

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf backend/bin
	rm -rf backend/coverage.out
	rm -rf backend/coverage.html

# Docker commands
docker-build:
	@echo "Building Docker images..."
	$(DOCKER_COMPOSE) build

docker-up:
	@echo "Starting Docker containers..."
	$(DOCKER_COMPOSE) up -d

docker-up-logs:
	@echo "Starting Docker containers with logs..."
	$(DOCKER_COMPOSE) up

docker-down:
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) down

docker-logs:
	@echo "Viewing Docker logs..."
	$(DOCKER_COMPOSE) logs -f

docker-ps:
	@echo "Listing Docker containers..."
	$(DOCKER_COMPOSE) ps

docker-clean:
	@echo "Cleaning Docker resources..."
	$(DOCKER_COMPOSE) down -v --rmi all

# Database migrations (for local development)
migrate-up:
	@echo "Running migrations..."
	# Add migration command here

migrate-down:
	@echo "Rolling back migrations..."
	# Add rollback command here

# Development helpers
dev-auth:
	@echo "Starting auth service in development mode..."
	cd backend && DEBUG=true DATABASE_URL="postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable" REDIS_URL="redis://localhost:6379" go run ./services/auth

dev-gateway:
	@echo "Starting gateway in development mode..."
	cd backend && DEBUG=true go run ./services/gateway

# Frontend commands
web-install:
	@echo "Installing web dependencies..."
	cd frontend/web && npm install

web-dev:
	@echo "Starting web dev server..."
	cd frontend/web && npm start

web-build:
	@echo "Building web app..."
	cd frontend/web && npm run build

mobile-install:
	@echo "Installing mobile dependencies..."
	cd frontend/mobile && npm install

mobile-ios:
	@echo "Starting iOS app..."
	cd frontend/mobile && npx expo start --ios

mobile-android:
	@echo "Starting Android app..."
	cd frontend/mobile && npx expo start --android

# Help
help:
	@echo "Available commands:"
	@echo "  make build          - Build all services"
	@echo "  make run SERVICE=x  - Run a specific service"
	@echo "  make run-gateway    - Run the API gateway"
	@echo "  make test           - Run tests"
	@echo "  make proto          - Generate protobuf files"
	@echo "  make docker-up      - Start all Docker containers"
	@echo "  make docker-down    - Stop all Docker containers"
	@echo "  make docker-logs    - View Docker logs"
	@echo "  make web-dev        - Start web development server"
	@echo "  make mobile-ios     - Start iOS app"
	@echo "  make mobile-android - Start Android app"


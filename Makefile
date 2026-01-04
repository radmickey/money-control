.PHONY: start stop restart build logs clean status help

# Default target
help:
	@echo "Money Control - Available commands:"
	@echo ""
	@echo "  make start     - Start all services"
	@echo "  make stop      - Stop all services"
	@echo "  make restart   - Restart all services"
	@echo "  make build     - Build and start all services"
	@echo "  make logs      - Show logs (follow mode)"
	@echo "  make status    - Show services status"
	@echo "  make clean     - Stop and remove all containers, volumes"
	@echo "  make db-reset  - Reset all databases (WARNING: deletes all data)"
	@echo ""

# Start all services
start:
	@echo "🚀 Starting Money Control..."
	docker compose up -d
	@echo "✅ Services started! Open http://localhost:3000"

# Stop all services
stop:
	@echo "🛑 Stopping Money Control..."
	docker compose down
	@echo "✅ Services stopped"

# Restart all services
restart: stop start

# Build and start all services
build:
	@echo "🔨 Building and starting Money Control..."
	docker compose up -d --build
	@echo "✅ Services built and started! Open http://localhost:3000"

# Show logs
logs:
	docker compose logs -f

# Show logs for specific service
logs-%:
	docker compose logs -f $*

# Show services status
status:
	@echo "📊 Services status:"
	@docker compose ps

# Clean everything (containers + volumes)
clean:
	@echo "🧹 Cleaning up..."
	docker compose down -v --remove-orphans
	@echo "✅ Cleanup complete"

# Reset databases (WARNING: deletes all data)
db-reset:
	@echo "⚠️  WARNING: This will delete all data!"
	@read -p "Are you sure? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	docker compose down -v
	docker compose up -d
	@echo "✅ Databases reset"

# Development: rebuild specific service
rebuild-%:
	docker compose build $*
	docker compose up -d $*

.PHONY: build down rebuild logs dev-db-up dev-db-down dev-db-logs dev-db-clean tidy run

# Docker commands
build:
	@echo "Building Docker images..."
	docker-compose up -d

down:
	@echo "Stopping Docker containers and cleaning up volumes and images..."
	docker-compose down -v --rmi all

rebuild:
	@echo "Rebuilding and restarting Docker containers..."
	docker-compose up -d --build

logs:
	@echo "Viewing Docker container logs (press Ctrl+C to exit)..."
	docker-compose logs -f

# Development Database Commands
dev-db-up:
	@echo "Starting development database container..."
	docker-compose -f docker-compose.dev.yml up -d

dev-db-down:
	@echo "Stopping development database container..."
	docker-compose -f docker-compose.dev.yml down

dev-db-logs:
	@echo "Viewing development database container logs (press Ctrl+C to exit)..."
	docker-compose -f docker-compose.dev.yml logs -f

dev-db-clean:
	@echo "Cleaning up development database containers, volumes, and images..."
	docker-compose -f docker-compose.dev.yml down -v --rmi all

# Go commands
tidy:
	@echo "Tidying Go modules..."
	go mod tidy

run:
	@echo "Running Go application locally..."
	go run main.go
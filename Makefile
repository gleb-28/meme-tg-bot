.PHONY: build down rebuild logs tidy run

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


# Go commands
tidy:
	@echo "Tidying Go modules..."
	go mod tidy

run:
	@echo "Running Go application locally..."
	go run cmd/bot/main.go
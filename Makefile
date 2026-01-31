.PHONY: build tidy run


build:
	@echo "Building Go application..."
	go build -o meme-tg-bot cmd/bot/main.go

tidy:
	@echo "Tidying Go modules..."
	go mod tidy

run:
	@echo "Running Go application locally..."
	go run cmd/bot/main.go
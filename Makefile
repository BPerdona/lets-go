include .env
export

default: run

run: format test clean build
	@echo "Running application..."
	@go run ./cmd/lets-go

build:
	@echo "Building application..."
	@go build -o bin/lets-go ./cmd/lets-go

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin

update:
	@echo "Updating dependencies..."
	@go mod tidy

format:
	@echo "Formatting code..."
	@go fmt ./...

test:
	@echo "Running tests..."
	@go test ./...

tcp-connect:
	telnet localhost $(SERVER_PORT)
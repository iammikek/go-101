.PHONY: help test test-coverage lint build run docker-up docker-down install-tools

help:
	@echo "Available commands:"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make lint          - Run linters"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application locally"
	@echo "  make docker-up     - Start Docker containers"
	@echo "  make docker-down   - Stop Docker containers"
	@echo "  make install-tools - Install development tools"

test:
	go test -v -race ./internal/... ./tests/...

test-coverage:
	go test -v -race -coverprofile=coverage.out ./internal/... ./tests/...
	go tool cover -func=coverage.out

lint:
	golangci-lint run ./...

build:
	go build -v -o go-101 .

run:
	go run .

docker-up:
	docker compose up --build

docker-down:
	docker compose down

install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install -tags sqlite github.com/golang-migrate/migrate/v4/cmd/migrate@latest

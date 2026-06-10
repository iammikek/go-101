.PHONY: help test test-coverage test-coverage-report lint build run docker-up docker-down migrate-install migrate-up migrate-down migrate-new install-tools

help:
	@echo "Available commands:"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make lint          - Run linters"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application locally"
	@echo "  make docker-up     - Start Docker containers"
	@echo "  make docker-down   - Stop Docker containers"
	@echo "  make install-tools - Install required tools"

test:
	go test -v -race ./...

test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-coverage-report:
	go test -v -race -coverprofile=coverage.out ./...
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

migrate-install:
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

migrate-up:
	migrate -path ./migrations -database "sqlite3://app.db" up

migrate-down:
	migrate -path ./migrations -database "sqlite3://app.db" down

migrate-new:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir ./migrations -seq $$name

install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install -tags sqlite github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b ./bin

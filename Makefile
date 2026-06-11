.PHONY: help test test-local test-coverage lint lint-local build run docker-up docker-down docker-shell install-tools

DOCKER_COMPOSE = docker compose
DEV = $(DOCKER_COMPOSE) run --rm dev

help:
	@echo "Available commands:"
	@echo "  make test          - Run tests (Docker)"
	@echo "  make test-local    - Run tests (local Go)"
	@echo "  make test-coverage - Run tests with coverage report (Docker)"
	@echo "  make lint          - Run linters (Docker)"
	@echo "  make lint-local    - Run linters (local Go)"
	@echo "  make build         - Build the application (Docker)"
	@echo "  make run           - Run the application locally"
	@echo "  make docker-up     - Start Docker containers"
	@echo "  make docker-down   - Stop Docker containers"
	@echo "  make docker-shell  - Open a shell in the dev container"
	@echo "  make install-tools - Install development tools (local Go)"

test:
	$(DEV) go test -v -race ./internal/... ./tests/...

test-local:
	go test -v -race ./internal/... ./tests/...

test-coverage:
	$(DEV) sh -c 'go test -v -race -coverprofile=coverage.out ./internal/... ./tests/... && go tool cover -func=coverage.out'

lint:
	$(DEV) golangci-lint run ./...

lint-local:
	golangci-lint run ./...

build:
	$(DEV) go build -v -o go-101 .

run:
	go run .

docker-up:
	$(DOCKER_COMPOSE) up --build api

docker-down:
	$(DOCKER_COMPOSE) down

docker-shell:
	$(DOCKER_COMPOSE) run --rm dev bash

install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install -tags sqlite github.com/golang-migrate/migrate/v4/cmd/migrate@latest

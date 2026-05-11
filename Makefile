SHELL := /bin/bash

GO       ?= go
PKG      ?= ./...
APP_NAME ?= tourismania-api

DB_URL ?= postgres://root:qwerty123@localhost:5432/tourismania?sslmode=disable

.PHONY: help tidy build run test test-unit test-integration test-app lint \
        migrate-up migrate-down migrate-new sqlc swag jwt-keys docker-up docker-down

help:
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  %-18s %s\n", $$1, $$2}'

tidy: ## go mod tidy
	$(GO) mod tidy

build: ## Build server + cli binaries into ./bin
	mkdir -p bin
	$(GO) build -trimpath -o bin/server ./cmd/server
	$(GO) build -trimpath -o bin/cli    ./cmd/cli

run: ## Run the HTTP server
	$(GO) run ./cmd/server

test: ## Run all tests
	$(GO) test -race ./tests/...

test-unit: ; $(GO) test -race ./tests/unit/...
test-integration: ; $(GO) test -race ./tests/integration/...
test-app: ; $(GO) test -race ./tests/application/...

migrate-up: ## Apply all migrations
	migrate -path ./migrations -database "$(DB_URL)" up

migrate-down: ## Roll back last migration
	migrate -path ./migrations -database "$(DB_URL)" down 1

migrate-new: ## Create new migration: make migrate-new name=foo
	migrate create -ext sql -dir ./migrations -seq $(name)

sqlc: ## Regenerate sqlc code
	sqlc generate

swag: ## Generate OpenAPI docs into ./docs
	swag init -g cmd/server/main.go -o docs

jwt-keys: ## Generate RSA keypair for JWT (RS256)
	openssl genpkey -algorithm RSA -out config/jwt/private.pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -pubout -in config/jwt/private.pem -out config/jwt/public.pem

docker-up: ## docker-compose up
	docker-compose up -d

docker-down: ## docker-compose down
	docker-compose down

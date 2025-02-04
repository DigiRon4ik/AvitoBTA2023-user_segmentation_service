.DEFAULT_GOAL := run
.PHONY: run lint

lint:
	@golangci-lint run

run: lint
	@go run cmd/app/main.go

up:
	@docker-compose up -d

down:
	@docker-compose down

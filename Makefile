.DEFAULT_GOAL := run
.PHONY: run lint

lint:
	@golangci-lint run

run: lint
	@go run cmd/app/main.go

.PHONY: docs build run tidy

docs:
	$(GOPATH)/bin/swag init -g cmd/main.go -o docs --parseInternal

build:
	go build ./...

run:
	go run ./cmd

tidy:
	go mod tidy

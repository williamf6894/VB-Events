.PHONY: docs build run tidy localdb test seed

SWAG := $(shell go env GOPATH)/bin/swag

seed:
	go run ./cmd/seed

test:
	go test ./... -v

docs:
	$(SWAG) init -g cmd/main.go -o docs --parseInternal

build:
	go build ./...

run:
	go run ./cmd

tidy:
	go mod tidy

localdb:
	../testing/start-local-db.sh

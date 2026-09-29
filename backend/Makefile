.PHONY: docs build run tidy localdb

docs:
	$(GOPATH)/bin/swag init -g cmd/main.go -o docs --parseInternal

build:
	go build ./...

run:
	go run ./cmd

tidy:
	go mod tidy

localdb:
	../testing/start-local-db.sh

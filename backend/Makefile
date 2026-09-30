.PHONY: docs build run tidy localdb test seed api-test

SWAG := $(shell go env GOPATH)/bin/swag

seed:
	go run ./cmd/seed

test:
	go test ./... -v

api-test:
	npx -y newman run ../testing/postman/vb-events.postman_collection.json \
		-e ../testing/postman/vb-events.postman_environment.json

docs:
	$(SWAG) init -g cmd/main.go -o docs --parseInternal

build:
	go build ./...

run:
	go run ./cmd

fmt:
	go fmt ./...

tidy:
	go mod tidy

localdb:
	../testing/start-local-db.sh

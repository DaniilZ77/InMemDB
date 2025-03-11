.PHONY: server client lint test build

server:
	CONFIG_PATH=./config/server.yaml go run ./cmd/db

client:
	go run ./cmd/client --address="127.0.0.1:3223"

lint:
	golangci-lint run

test:
	go test -v ./...

build_server:
	go build -o bin/db ./cmd/db

build_client:
	go build -o bin/client ./cmd/client

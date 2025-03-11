server:
	CONFIG_PATH=./config/local.yaml go run ./cmd/db

lint:
	golangci-lint run

test:
	go test -v ./...

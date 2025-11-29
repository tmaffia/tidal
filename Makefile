.PHONY: test lint all

all: test lint

test:
	go test -v ./...

lint:
	go run -modfile=tools/go.mod github.com/golangci/golangci-lint/cmd/golangci-lint run

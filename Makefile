.PHONY: test lint all

all: test lint

test:
	go test -v ./...

lint:
	go tool -modfile=tools/go.mod golangci-lint run

.PHONY: all deps test lint

all: deps lint test

deps:
	go get -t -v ./...

test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic

lint:
	golangci-lint run ./...

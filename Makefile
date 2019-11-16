SHELL := /bin/bash -o pipefail

.PHONY: check-linter
check-linter:
	@which golangci-lint >/dev/null || (echo "ERROR: golangci-lint not found" && false)

.PHONY: lint
lint: check-linter
	golangci-lint run ./...

.PHONY: test
test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic

.PHONY: cover
cover:
	go tool cover -html=coverage.txt

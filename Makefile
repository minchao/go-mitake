SHELL := /usr/bin/env bash -o errexit -o pipefail -o nounset

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-23s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: all
all: lint test ## Run linters and tests

.PHONY: lint
lint: install.golangci-lint ## Run linters
	golangci-lint run ./...

.PHONY: test
test: ## Run unit tests
	go test -v -race -coverprofile=coverage.txt -covermode=atomic

.PHONY: cover
cover: test ## Run unit tests and open coverage report in browser
	go tool cover -html=coverage.txt

## Dependency Tools

.PHONY: install.golangci-lint
install.golangci-lint:
ifeq (, $(shell which golangci-lint))
	$(error The 'golangci-lint' command not found, install it from https://golangci-lint.run/usage/install/ or via 'brew install golangci-lint')
endif

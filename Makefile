.PHONY: help
help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the provider
	go generate ./...
	go build -o ${BINARY}

.PHONY: test
test: ## Run tests
	go test -cover -race ./...

.PHONY: cover
cover: ## Generate Go coverage report
	@echo "mode: count" > coverage.out
	@go test -coverprofile coverage.tmp ./...
	@tail -n +2 coverage.tmp >> coverage.out
	@go tool cover -html=coverage.out

.PHONY: mocks
mocks: ## Generate mocks for unit tests
	mockery

# variables for local build that's used for testing and development
HOSTNAME=border0.com
NAMESPACE=border0
NAME=border0
BINARY=terraform-provider-${NAME}
VERSION=0.1.0
OS_ARCH?=$(shell go env GOOS)_$(shell go env GOARCH)

.PHONY: install
install: build ## Install the provider in the terraform plugins directory
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

.PHONY: docs
docs: ## Generate terraform provider docs
	go generate ./...

.PHONY: release
release:
	goreleaser release --clean --snapshot --skip=publish,sign

.PHONY: fmt
fmt:
	@terraform fmt --recursive .

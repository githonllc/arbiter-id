.PHONY: help build test test-race test-coverage lint fmt vet clean examples benchmark deps

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: ## Build the project
	@echo "Building..."
	@go build -v ./...

# Test targets
test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	@go test -v -race ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Code quality targets
lint: ## Run golangci-lint
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		$(shell go env GOPATH)/bin/golangci-lint run; \
	fi

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

# Dependency management
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod verify

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	@go mod tidy

# Examples
examples: ## Run examples
	@echo "Running examples..."
	@cd examples/simple && go run main.go

# Clean targets
clean: ## Clean build artifacts
	@echo "Cleaning..."
	@go clean ./...
	@rm -f coverage.out coverage.html

# CI targets
ci: deps vet lint test-race ## Run CI pipeline
	@echo "CI pipeline completed successfully"

# Development targets
dev: fmt vet test ## Run development checks
	@echo "Development checks completed"

# Release preparation
pre-release: clean deps fmt vet lint test-race test-coverage benchmark ## Prepare for release
	@echo "Pre-release checks completed"
.PHONY: help build test lint clean run docker-build docker-push

# Variables
IMAGE_TAG ?= latest
REGISTRY ?= kubeops

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the monolithic application
	@echo "Building KubeOps..."
	go build -o bin/kubeops ./cmd/server
	@echo "Build complete!"

run: ## Run the application locally
	@echo "Starting KubeOps..."
	go run ./cmd/server/main.go

test: ## Run all tests
	@echo "Running tests..."
	go test -v ./...

lint: ## Run golangci-lint
	@echo "Running linter..."
	golangci-lint run

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	@echo "Clean complete!"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(REGISTRY)/kubeops:$(IMAGE_TAG) -f Dockerfile .
	@echo "Docker build complete!"

docker-push: ## Push Docker image
	@echo "Pushing Docker image..."
	docker push $(REGISTRY)/kubeops:$(IMAGE_TAG)
	@echo "Docker push complete!"

mod-tidy: ## Tidy go modules
	go mod tidy

mod-download: ## Download go modules
	go mod download

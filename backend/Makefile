.PHONY: help build test lint clean docker-build docker-push dev-deploy

# Variables
IMAGE_TAG ?= latest
REGISTRY ?= kubeops
SERVICES = api-gateway kube-manager ai-inspector devops-service logging-service monitoring-service

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build all backend services
	@echo "Building all services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd $$service && go build -o ../../bin/$$service ./cmd/server && cd ..; \
	done
	@echo "Build complete!"

build-%: ## Build specific service (e.g., make build-api-gateway)
	@echo "Building $*..."
	@cd $* && go build -o ../bin/$* ./cmd/server
	@echo "Build complete!"

test: ## Run all tests
	@echo "Running tests..."
	go test -v ./...

test-%: ## Run tests for specific service
	@echo "Running tests for $*..."
	cd $* && go test -v ./...

lint: ## Run golangci-lint
	@echo "Running linter..."
	golangci-lint run

lint-%: ## Run linter for specific service
	@echo "Running linter for $*..."
	golangci-lint run ./$*/...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	@echo "Clean complete!"

run-api-gateway: ## Run API Gateway locally
	@echo "Starting API Gateway..."
	cd api-gateway && go run cmd/server/main.go

run-kube-manager: ## Run Kube Manager locally
	@echo "Starting Kube Manager..."
	cd kube-manager && go run cmd/server/main.go

docker-build: ## Build all Docker images
	@echo "Building Docker images..."
	@for service in $(SERVICES); do \
		echo "Building Docker image for $$service..."; \
		docker build --build-arg SERVICE=$$service -t $(REGISTRY)/$$service:$(IMAGE_TAG) -f Dockerfile .; \
	done
	@echo "Docker build complete!"

docker-build-%: ## Build specific Docker image
	@echo "Building Docker image for $*..."
	docker build --build-arg SERVICE=$* -t $(REGISTRY)/$*:$(IMAGE_TAG) -f Dockerfile .
	@echo "Docker build complete!"

docker-push: ## Push all Docker images
	@echo "Pushing Docker images..."
	@for service in $(SERVICES); do \
		echo "Pushing $$service..."; \
		docker push $(REGISTRY)/$$service:$(IMAGE_TAG); \
	done
	@echo "Docker push complete!"

mod-tidy: ## Tidy go modules
	go mod tidy

mod-download: ## Download go modules
	go mod download

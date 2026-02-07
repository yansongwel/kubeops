#!/bin/bash
# KubeOps Development Scripts

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Help function
show_help() {
    echo "KubeOps Development Scripts"
    echo ""
    echo "Usage: ./scripts/dev.sh [command]"
    echo ""
    echo "Commands:"
    echo "  setup          - Set up development environment"
    echo "  kind           - Create local kind cluster"
    echo "  deploy         - Deploy to local cluster"
    echo "  destroy        - Destroy local cluster"
    echo "  logs           - Show logs from all services"
    echo "  test           - Run all tests"
    echo "  lint           - Run linting"
    echo "  build          - Build all services"
    echo ""
}

# Setup development environment
setup() {
    echo -e "${GREEN}Setting up development environment...${NC}"

    # Check prerequisites
    command -v go >/dev/null 2>&1 || { echo -e "${RED}Go is required but not installed${NC}"; exit 1; }
    command -v node >/dev/null 2>&1 || { echo -e "${RED}Node.js is required but not installed${NC}"; exit 1; }
    command -v docker >/dev/null 2>&1 || { echo -e "${RED}Docker is required but not installed${NC}"; exit 1; }
    command -v kubectl >/dev/null 2>&1 || { echo -e "${RED}kubectl is required but not installed${NC}"; exit 1; }
    command -v helm >/dev/null 2>&1 || { echo -e "${RED}helm is required but not installed${NC}"; exit 1; }
    command -v kind >/dev/null 2>&1 || { echo -e "${RED}kind is required but not installed${NC}"; exit 1; }

    echo -e "${GREEN}Installing Go dependencies...${NC}"
    cd backend
    go mod download
    go install github.com/cosmtrek/air@latest
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    cd ..

    echo -e "${GREEN}Installing Node dependencies...${NC}"
    cd frontend
    npm install
    cd ..

    echo -e "${GREEN}Installing buf (protobuf)...${NC}"
    go install github.com/bufbuild/buf/cmd/buf@latest

    echo -e "${GREEN}✓ Development environment setup complete!${NC}"
}

# Create kind cluster
kind() {
    echo -e "${GREEN}Creating kind cluster...${NC}"

    kind create cluster --name kubeops --config scripts/kind-config.yaml

    echo -e "${GREEN}✓ Kind cluster created${NC}"
    echo -e "${YELLOW}Note: Don't forget to configure kubectl context${NC}"
}

# Deploy to local cluster
deploy() {
    echo -e "${GREEN}Deploying to local cluster...${NC}"

    # Create namespaces
    kubectl create namespace kubeops --dry-run=client -o yaml | kubectl apply -f -
    kubectl create namespace logging --dry-run=client -o yaml | kubectl apply -f -
    kubectl create namespace monitoring --dry-run=client -o yaml | kubectl apply -f -

    # Install dependencies
    helm repo add bitnami https://charts.bitnami.com/bitnami
    helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
    helm repo add grafana https://grafana.github.io/helm-charts
    helm repo update

    # Deploy PostgreSQL
    helm upgrade --install postgresql bitnami/postgresql \
        --namespace kubeops \
        --set auth.password=postgres \
        --set persistence.enabled=false

    # Deploy Redis
    helm upgrade --install redis bitnami/redis \
        --namespace kubeops \
        --set architecture=standalone \
        --set auth.enabled=false \
        --set persistence.enabled=false

    # Build and deploy services
    echo -e "${YELLOW}Building Docker images...${NC}"
    docker build -t kubeops/api-gateway:dev -f backend/api-gateway/Dockerfile backend/
    docker build -t kubeops/kube-manager:dev -f backend/kube-manager/Dockerfile backend/
    docker build -t kubeops/frontend:dev -f frontend/Dockerfile frontend/

    # Load images into kind
    kind load docker-image kubeops/api-gateway:dev --name kubeops
    kind load docker-image kubeops/kube-manager:dev --name kubeops
    kind load docker-image kubeops/frontend:dev --name kubeops

    # Deploy with Helm
    helm upgrade --install kubeops deploy/helm/kubeops \
        --namespace kubeops \
        --set image.tag=dev \
        --values deploy/examples/values-dev.yaml

    echo -e "${GREEN}✓ Deployment complete${NC}"
}

# Destroy kind cluster
destroy() {
    echo -e "${YELLOW}Destroying kind cluster...${NC}"
    kind delete cluster --name kubeops
    echo -e "${GREEN}✓ Cluster destroyed${NC}"
}

# Show logs
logs() {
    kubectl logs -n kubeops -l app.kubernetes.io/name=kubeops --all-follow=true
}

# Run tests
test() {
    echo -e "${GREEN}Running tests...${NC}"

    # Backend tests
    echo -e "${YELLOW}Running backend tests...${NC}"
    cd backend
    go test -v ./...

    # Frontend tests
    echo -e "${YELLOW}Running frontend tests...${NC}"
    cd ../frontend
    npm run test

    echo -e "${GREEN}✓ All tests passed${NC}"
}

# Run linting
lint() {
    echo -e "${GREEN}Running linters...${NC}"

    # Backend linting
    echo -e "${YELLOW}Running golangci-lint...${NC}"
    cd backend
    golangci-lint run

    # Frontend linting
    echo -e "${YELLOW}Running ESLint...${NC}"
    cd ../frontend
    npm run lint

    echo -e "${GREEN}✓ Linting complete${NC}"
}

# Build all services
build() {
    echo -e "${GREEN}Building all services...${NC}"

    # Backend services
    echo -e "${YELLOW}Building backend services...${NC}"
    cd backend
    go build -o bin/api-gateway ./api-gateway/cmd/server
    go build -o bin/kube-manager ./kube-manager/cmd/server

    # Frontend
    echo -e "${YELLOW}Building frontend...${NC}"
    cd ../frontend
    npm run build

    echo -e "${GREEN}✓ Build complete${NC}"
}

# Main script
case "$1" in
    setup)
        setup
        ;;
    kind)
        kind
        ;;
    deploy)
        deploy
        ;;
    destroy)
        destroy
        ;;
    logs)
        logs
        ;;
    test)
        test
        ;;
    lint)
        lint
        ;;
    build)
        build
        ;;
    *)
        show_help
        ;;
esac

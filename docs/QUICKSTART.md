# Quick Start Guide

This guide will help you get KubeOps up and running quickly.

## Prerequisites

Ensure you have the following installed:

- **Go** 1.21+ - [Download](https://golang.org/dl/)
- **Node.js** 18+ - [Download](https://nodejs.org/)
- **Docker** - [Download](https://docker.com/get-docker)
- **kubectl** - [Install guide](https://kubernetes.io/docs/tasks/tools/)
- **helm** 3.x - [Install guide](https://helm.sh/docs/intro/install/)
- **kind** (for local K8s) - [Install guide](https://kind.sigs.k8s.io/docs/user/quick-start/)

## Quick Start (5 Minutes)

### Option 1: Local Development (Recommended for First Time)

#### Step 1: Clone and Setup

```bash
# Clone the repository
git clone https://github.com/your-org/kubeops.git
cd kubeops

# Install dependencies
./scripts/dev.sh setup
```

#### Step 2: Start Infrastructure

```bash
# Start PostgreSQL and Redis
docker-compose -f deploy/docker-compose-dev.yaml up -d
```

#### Step 3: Run Monolith Backend and Frontend

```bash
# Terminal 1 - Monolith Backend
cd backend
go run cmd/server/main.go

# Terminal 2 - Frontend
cd frontend
npm install
npm run dev
```

#### Step 4: Access Dashboard

Open your browser and navigate to:
```
http://localhost:5173
```

Default credentials (for development):
- Username: `admin`
- Password: `admin`

### Option 2: Deploy to Local Kind Cluster

#### Step 1: Create Kind Cluster

```bash
./scripts/dev.sh kind
```

#### Step 2: Deploy KubeOps

```bash
# Set kubectl context
kubectl cluster-info --context kind-kubeops

# Deploy with Helm
helm install kubeops deploy/helm/kubeops \
  --namespace kubeops \
  --create-namespace \
  -f deploy/examples/values-dev.yaml

# Wait for pods to be ready
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=kubeops -n kubeops --timeout=300s
```

#### Step 3: Access Dashboard

```bash
# Port forward to access the dashboard
kubectl port-forward -n kubeops svc/kubeops 8080:80

# Or port forward in background
kubectl port-forward -n kubeops svc/kubeops 8080:80 &
```

Open your browser:
```
http://localhost:8080
```

## Verify Installation

### Check Backend Health

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "service": "kubeops"
}
```

### Check Kubernetes Connection

```bash
# List namespaces via API
curl http://localhost:8080/api/v1/namespaces
```

### Check Frontend

Open browser and verify dashboard loads.

## Next Steps

### Explore the Dashboard

1. **Dashboard**: Overview of cluster health
2. **Clusters**: Manage your K8s clusters
3. **Workloads**: View and manage deployments, pods, etc.
4. **AI Inspector**: AI-powered cluster analysis
5. **DevOps**: Jenkins and ArgoCD integration
6. **Logs**: View and search logs
7. **Monitoring**: Metrics and dashboards

### Configure Integrations

#### AI Provider (Optional)

Edit `deploy/helm/kubeops/values.yaml`:

```yaml
config:
  ai:
    provider: ollama  # or openai
    openai:
      apiKey: "your-api-key"
    ollama:
      url: http://ollama:11434
```

#### Logging Stack

Deploy ELK stack:
```bash
helm install logging deploy/helm/logging/logging-elk \
  --namespace logging \
  --create-namespace
```

Or Loki stack:
```bash
helm install logging deploy/helm/logging/logging-loki \
  --namespace logging \
  --create-namespace
```

#### Monitoring Stack

Deploy Prometheus:
```bash
helm install monitoring deploy/helm/monitoring/monitoring-prometheus \
  --namespace monitoring \
  --create-namespace
```

Or VictoriaMetrics:
```bash
helm install monitoring deploy/helm/monitoring/monitoring-victoriametrics \
  --namespace monitoring \
  --create-namespace
```

## Development Workflow

### Run Tests

```bash
# Backend tests
cd backend
make test

# Frontend tests
cd frontend
npm run test
```

### Build for Production

```bash
# Build backend
cd backend
make build

# Build frontend
cd frontend
npm run build
```

### Create Docker Images

```bash
# Build monolith backend image
cd backend
make docker-build
```

## Troubleshooting

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```

### Cannot Connect to Kubernetes

```bash
# Check kubeconfig
echo $KUBECONFIG

# Verify cluster access
kubectl cluster-info

# Check context
kubectl config current-context
```

### Database Connection Error

```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Check logs
docker logs kubeops-postgres

# Restart if needed
docker-compose -f deploy/docker-compose-dev.yaml restart postgres
```

### Frontend Build Failed

```bash
# Clean and reinstall
cd frontend
rm -rf node_modules package-lock.json
npm install
npm run dev
```

## Uninstall

### Stop Local Development

```bash
# Stop services
pkill -f "go run"

# Stop databases
docker-compose -f deploy/docker-compose-dev.yaml down
```

### Remove from Kind Cluster

```bash
# Uninstall Helm release
helm uninstall kubeops -n kubeops

# Delete namespace
kubectl delete namespace kubeops logging monitoring

# Delete kind cluster
kind delete cluster --name kubeops
```

## Resources

- [Architecture Documentation](docs/architecture/README.md)
- [API Documentation](docs/api/README.md) (coming soon)
- [Development Guide](docs/development/README.md) (coming soon)
- [Deployment Guide](docs/deployment/README.md) (coming soon)
- [Contributing Guide](CONTRIBUTING.md)

## Getting Help

- **Issues**: [GitHub Issues](https://github.com/your-org/kubeops/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-org/kubeops/discussions)
- **Documentation**: [Full Docs](https://docs.kubeops.io)

## What's Next?

- 📚 Read the [Architecture Guide](docs/architecture/README.md)
- 🔧 Configure [AI Provider](docs/configuration/ai-provider.md)
- 🚀 Set up [CI/CD Integration](docs/configuration/cicd.md)
- 📊 Deploy [Logging Stack](docs/deployment/logging.md)
- 📈 Deploy [Monitoring Stack](docs/deployment/monitoring.md)

---

Happy managing! 🎉

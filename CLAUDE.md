# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**KubeOps** is an enterprise-grade Kubernetes management platform designed to provide comprehensive K8s resource management, AI-powered inspection, DevOps/CI/CD integration, logging, and monitoring capabilities.

### Vision
- Manage all K8s resources and CRDs with extensibility
- AI-powered cluster inspection and health analysis
- Integrated DevOps pipeline (Jenkins + ArgoCD)
- Flexible logging architecture (ELK or Vector+Loki stack)
- Flexible monitoring architecture (Prometheus or VictoriaMetrics)
- Open-source, adaptable for all enterprise needs

## Technology Stack

### Backend (Go)
- **Language**: Go 1.21+
- **API Framework**: Gin (REST) + gRPC (internal service communication)
- **K8s Client**: controller-runtime + client-go
- **Database**: PostgreSQL (primary) + Redis (cache)
- **Message Queue**: NATS (optional, for async tasks)

### Frontend (Vue)
- **Framework**: Vue 3 with Composition API
- **Language**: TypeScript
- **Build Tool**: Vite
- **UI Framework**: Element Plus or Ant Design Vue
- **State Management**: Pinia
- **Routing**: Vue Router

### Infrastructure
- **Container Runtime**: Docker / Containerd
- **Orchestration**: Kubernetes 1.28+
- **Deployment**: Helm charts
- **Service Mesh**: Istio (optional)
- **Ingress**: Nginx / Traefik

## Architecture

### Microservices Layout

```
KubeOps/
├── backend/
│   ├── api-gateway/          # API Gateway (auth, routing, rate limiting)
│   ├── kube-manager/         # K8s resource management service
│   ├── ai-inspector/         # AI inspection service
│   ├── devops-service/       # CI/CD integration (Jenkins + ArgoCD)
│   ├── logging-service/      # Logging platform integration
│   └── monitoring-service/   # Monitoring platform integration
│   ├── common/               # Shared libraries (go modules)
│   └── proto/                # gRPC protocol definitions
├── frontend/                 # Vue 3 Dashboard
├── deploy/                   # Kubernetes manifests & Helm charts
│   ├── helm/
│   │   ├── kubeops/
│   │   ├── logging/
│   │   │   ├── elk/         # ELK stack
│   │   │   └── loki/        # Vector+Loki stack
│   │   └── monitoring/
│   │       ├── prometheus/  # Prometheus+Grafana
│   │       └── victoriametrics/ # Prometheus+VictoriaMetrics
│   ├── base/                # Base infrastructure (PostgreSQL, Redis, etc.)
│   └── examples/            # Example configurations
├── docs/                    # Architecture docs, API docs, guides
├── scripts/                 # Development scripts
└── tests/                   # E2E tests, integration tests
```

### Core Services

#### 1. API Gateway (`api-gateway/`)
- **Purpose**: Single entry point, authentication, authorization, routing
- **Key Responsibilities**:
  - JWT/OAuth authentication
  - RBAC authorization
  - Request routing to backend services
  - Rate limiting and throttling
  - API versioning
- **Port**: 8080 (HTTP), 9090 (gRPC)

#### 2. Kube Manager (`kube-manager/`)
- **Purpose**: Core Kubernetes resource management
- **Key Responsibilities**:
  - CRUD operations for all K8s resources (Deployments, Services, ConfigMaps, etc.)
  - CRD support and extensibility
  - Real-time resource watch and event streaming
  - Resource validation and policies
  - Multi-cluster support
- **K8s Integration**: Uses controller-runtime for efficient resource operations

#### 3. AI Inspector (`ai-inspector/`)
- **Purpose**: AI-powered cluster health inspection and analysis
- **Key Responsibilities**:
  - Anomaly detection in resource usage
  - Configuration analysis and recommendations
  - Security vulnerability scanning
  - Cost optimization suggestions
  - Predictive analytics (capacity planning)
  - Log analysis (NLP-based)
- **AI/ML Integration**:
  - Can integrate with OpenAI API, local LLMs (Ollama), or custom models
  - Prometheus metrics as input
  - Log data from logging service

#### 4. DevOps Service (`devops-service/`)
- **Purpose**: CI/CD pipeline integration
- **Integrations**:
  - Jenkins (job creation, triggering, status monitoring)
  - ArgoCD (GitOps deployment, sync status, application management)
- **Key Responsibilities**:
  - Pipeline management and visualization
  - Code quality integration (SonarQube, tests)
  - Deployment rollback and promotion
  - Environment management (dev, staging, prod)

#### 5. Logging Service (`logging-service/`)
- **Purpose**: Unified logging platform integration
- **Supported Stacks** (pluggable):
  - **Option 1**: Vector + Fluentd + Elasticsearch + Kibana (ELK)
  - **Option 2**: Vector + VictoriaLogs + Loki + Grafana + Prometheus
  - **Extensible**: Graylog, Watchalert, etc.
- **Key Responsibilities**:
  - Log collection and aggregation
  - Log search and filtering
  - Log retention and archival policies
  - Integration with AI Inspector for analysis

#### 6. Monitoring Service (`monitoring-service/`)
- **Purpose**: Monitoring and alerting integration
- **Supported Stacks** (pluggable):
  - **Option 1**: Prometheus + Grafana
  - **Option 2**: Prometheus + VictoriaMetrics (cluster)
- **Key Responsibilities**:
  - Metrics collection and querying
  - Dashboard management
  - Alert rule management
  - Notification integration (Email, Slack, PagerDuty, etc.)
  - SLO/SLI tracking

### Service Communication

- **External API**: RESTful (HTTP/JSON) via API Gateway
- **Internal Services**: gRPC for high-performance inter-service communication
- **Events**: NATS for async messaging (optional)
- **Service Discovery**: Kubernetes native (DNS) or Consul

## Development Workflow

### Prerequisites
```bash
# Install required tools
- Go 1.21+
- Node.js 18+
- Docker/Podman
- kubectl + helm
- kind (local K8s) or access to a K8s cluster
- Make (optional, for build scripts)
- buf (for protobuf code generation)
```

### Building the Project

#### Backend (Go)
```bash
# Build all services
cd backend
make build

# Build specific service
cd backend/kube-manager
go build -o kube-manager ./cmd/server

# Run tests
cd backend
make test

# Run tests for specific package
cd backend/kube-manager
go test ./pkg/...
```

#### Frontend (Vue)
```bash
cd frontend
npm install
npm run dev          # Development server
npm run build        # Production build
npm run test         # Run tests
npm run lint         # Lint code
```

### Running Locally

#### Option 1: Using kind (local K8s cluster)
```bash
# Create local cluster
scripts/create-kind-cluster.sh

# Deploy to local cluster
helm install kubeops deploy/helm/kubeops --namespace kubeops --create-namespace

# Or use make
make dev-deploy
```

#### Option 2: Local development (services outside K8s)
```bash
# Start databases (PostgreSQL, Redis)
docker-compose -f deploy/docker-compose-dev.yaml up -d

# Run backend services (each terminal)
cd backend/api-gateway && go run cmd/server/main.go
cd backend/kube-manager && go run cmd/server/main.go
# ... other services

# Run frontend
cd frontend && npm run dev
```

### Testing

#### Unit Tests
```bash
# Backend
cd backend
make test-unit

# Frontend
cd frontend
npm run test:unit
```

#### Integration Tests
```bash
# Backend
cd backend
make test-integration

# Frontend
cd frontend
npm run test:integration
```

#### E2E Tests
```bash
cd tests/e2e
npm run test
```

### Deploying to Kubernetes

```bash
# Using Helm
helm upgrade --install kubeops deploy/helm/kubeops \
  --namespace kubeops \
  --create-namespace \
  --set logging.stack=elk \  # or 'loki'
  --set monitoring.stack=prometheus \  # or 'victoriametrics'
  -f deploy/examples/values-dev.yaml

# Deploy logging stack
helm install logging deploy/helm/logging/logging-elk --namespace logging

# Deploy monitoring stack
helm install monitoring deploy/helm/monitoring/monitoring-prometheus --namespace monitoring
```

## Configuration

### Environment Variables

Each service uses environment variables for configuration. See `deploy/helm/kubeops/values.yaml` for all configurable options.

Key configurations:
- `KUBECONFIG`: Path to kubeconfig (for development)
- `DATABASE_URL`: PostgreSQL connection string
- `REDIS_URL`: Redis connection string
- `JENKINS_URL`: Jenkins server URL
- `ARGOCD_URL`: ArgoCD server URL

### Pluggable Components

The platform supports multiple implementations for:

1. **Logging Backend**: Switch between ELK, Loki, or custom
   - Config: `logging.stack.type = "elk" | "loki"`
2. **Monitoring Backend**: Switch between Prometheus, VictoriaMetrics
   - Config: `monitoring.stack.type = "prometheus" | "victoriametrics"`
3. **AI Provider**: Switch between OpenAI, local LLM, or custom
   - Config: `ai.provider = "openai" | "ollama" | "custom"`

## Key Design Principles

1. **Extensibility First**: All integrations (logging, monitoring, AI) are plugin-based
2. **Multi-Tenant Ready**: Architecture supports multi-tenancy from day one
3. **Cloud Native**: Built for Kubernetes, uses K8s primitives
4. **API First**: Everything is accessible via API
5. **GitOps Friendly**: Configuration via Git, infrastructure as code
6. **Security**: RBAC, network policies, secrets management

## Database Schema

### PostgreSQL Tables
- `tenants` - Multi-tenant support
- `users` - User accounts and authentication
- `roles` - RBAC roles
- `permissions` - Fine-grained permissions
- `clusters` - Managed K8s clusters
- `pipelines` - CI/CD pipelines
- `applications` - Application metadata
- `inspection_reports` - AI inspection history
- `alerts` - Monitoring alerts
- `audit_logs` - Audit trail

## API Conventions

### REST API Design
- **Versioning**: `/api/v1/`, `/api/v2/`
- **Resource Naming**: Plural nouns (e.g., `/api/v1/clusters`, `/api/v1/deployments`)
- **HTTP Methods**: GET (list/get), POST (create), PUT/PATCH (update), DELETE (delete)
- **Response Format**: JSON with consistent envelope
  ```json
  {
    "data": { ... },
    "metadata": { "total": 100, "page": 1 }
  }
  ```

### gRPC Conventions
- **Proto Files**: Located in `backend/proto/`
- **Code Generation**: Use `buf` to generate Go code
- **Service Naming**: `[ServiceName]Service`
- **Package Structure**: Follow Google API design guide

## AI Inspector Notes

The AI Inspector service analyzes cluster health and provides recommendations:

### Data Sources
1. **Kubernetes API**: Resource configurations, events, status
2. **Prometheus Metrics**: CPU, memory, network, custom metrics
3. **Logs**: From logging service
4. **Historical Data**: Trends and patterns over time

### Analysis Types
1. **Resource Optimization**: Unused resources, over-provisioning
2. **Security**: Vulnerable images, RBAC issues, network policies
3. **Reliability**: Missing limits/requests, single points of failure
4. **Cost**: Cost reduction opportunities
5. **Compliance**: Best practice violations, security standards

### Integration Points
- Can call external AI APIs (OpenAI, Claude, etc.)
- Can run local models (Ollama, custom models)
- Results stored in PostgreSQL for historical analysis

## Contributing Guidelines

### Code Style
- **Go**: Follow `gofmt`, use `golangci-lint`
- **Vue**: Follow `vue-eslint`, use TypeScript strictly
- **Commits**: Conventional Commits format (`feat:`, `fix:`, `docs:`, etc.)

### Adding New Features
1. Create feature branch from `develop`
2. Implement with tests
3. Update documentation
4. Submit PR with description and testing notes

### Adding New Integrations
Follow the plugin pattern in `backend/common/plugin`:
1. Implement plugin interface
2. Add configuration options
3. Register plugin in service initialization
4. Add documentation

## Troubleshooting

### Common Issues

1. **Cannot connect to Kubernetes cluster**
   - Check `KUBECONFIG` environment variable
   - Verify service account has necessary permissions
   - Check RBAC rules in `deploy/base/rbac/`

2. **Services cannot communicate**
   - Verify Kubernetes DNS is working
   - Check network policies
   - Review service discovery configuration

3. **Database connection errors**
   - Check PostgreSQL is running: `kubectl get pods -n kubeops`
   - Verify connection string in secrets
   - Review database logs

4. **Frontend build fails**
   - Delete `node_modules/` and `package-lock.json`
   - Run `npm install`
   - Ensure Node.js version is 18+

### Debug Mode

Enable debug logging:
```bash
# Backend
export LOG_LEVEL=debug

# Frontend
npm run dev -- --debug
```

## Resources

- **Internal Documentation**: `docs/`
- **API Documentation**: Auto-generated from OpenAPI spec (Swagger UI at `/docs`)
- **Architecture Decisions**: `docs/architecture/adr-*.md`
- **Kubernetes Docs**: https://kubernetes.io/docs/
- **Vue Docs**: https://vuejs.org/guide/

## Project Status

This is a greenfield project. Core features are being implemented in phases:

- [x] Project architecture and setup
- [ ] API Gateway
- [ ] Kube Manager (core K8s resource management)
- [ ] AI Inspector
- [ ] DevOps Service (Jenkins + ArgoCD)
- [ ] Logging Service (ELK + Loki options)
- [ ] Monitoring Service (Prometheus + VictoriaMetrics options)
- [ ] Frontend Dashboard
- [ ] Multi-cluster support
- [ ] Documentation and examples

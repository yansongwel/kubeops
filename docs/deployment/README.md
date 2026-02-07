# 部署指南

本文档说明如何部署 KubeOps 的后端与前端。

## 前置要求

- Docker 或 Podman
- Kubernetes 集群（部署到 K8s 时）
- helm 3.x（部署到 K8s 时）

## 前端 Docker 部署

### 构建镜像

```bash
cd kubeops/frontend
docker build -t kubeops-frontend:latest -f frontend/Dockerfile frontend
```

### 运行容器

```bash
docker run --rm -p 8080:80 --name kubeops-frontend \
  --add-host host.docker.internal:host-gateway \
  -e API_GATEWAY=host.docker.internal:8080 \
  kubeops-frontend:latest
```

### 访问地址

```
http://localhost:8080
```

### 说明

- 前端默认通过 `/api` 代理后端，可通过 `API_GATEWAY` 环境变量自定义后端地址。
- 若后端也在 Docker 中运行，可加入同一网络，并设置 `API_GATEWAY=api-gateway:8080`。
- 若需要自定义 API 前缀，可在构建前设置 `VITE_API_BASE_URL` 环境变量。

## 后端 Docker 部署

### 构建镜像

```bash
cd kubeops/backend
docker build -t kubeops-backend:latest -f backend/Dockerfile .
```

### 运行容器

```bash
docker run --rm -p 8080:8080 --name kubeops-backend \
  -v ~/.kube/config:/root/.kube/config:ro \
  -e KUBECONFIG=/root/.kube/config \
  kubeops-backend:latest \
  --postgres-host 192.168.33.100 \
  --postgres-port 5432 \
  --postgres-user kubeops \
  --postgres-password kubeops \
  --postgres-db kubeops \
  --redis-addr 192.168.33.100:6379
```

### 说明

- 后端启动时优先使用集群内配置；本地 Docker 运行时推荐挂载 `~/.kube/config`。
- 你的 kubeconfig 需指向 `https://192.168.33.100:6443`（或集群实际 API Server 地址）。
- 可使用参数或环境变量配置连接信息与端口。
- 后端使用以下环境变量：
  - PostgreSQL：`POSTGRES_HOST`、`POSTGRES_PORT`、`POSTGRES_USER`、`POSTGRES_PASSWORD`、`POSTGRES_DB`、`POSTGRES_SSLMODE`
  - Redis：`REDIS_ADDR`、`REDIS_PASSWORD`、`REDIS_DB`
  - 端口：`PORT`

### 环境变量示例（与你当前环境一致）

```bash
docker run --rm -p 8080:8080 --name kubeops-backend \
  -v ~/.kube/config:/root/.kube/config:ro \
  -e KUBECONFIG=/root/.kube/config \
  -e POSTGRES_HOST=192.168.33.100 \
  -e POSTGRES_PORT=5432 \
  -e POSTGRES_USER=kubeops \
  -e POSTGRES_PASSWORD=kubeops \
  -e POSTGRES_DB=kubeops \
  -e REDIS_ADDR=192.168.33.100:6379 \
  -e REDIS_PASSWORD=kubeops \
  kubeops-backend:latest
```

### 帮助与参数列表

```bash
./kubeops --help
```

## 后端部署到 Kubernetes

```bash
kubectl create namespace kubeops

helm install kubeops deploy/helm/kubeops \
  --namespace kubeops \
  --set logging.stack=elk \
  --set monitoring.stack=prometheus
```

或使用自定义配置：

```bash
helm install kubeops deploy/helm/kubeops \
  --namespace kubeops \
  -f deploy/examples/values-production.yaml
```

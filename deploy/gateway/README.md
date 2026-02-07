# API 网关配置

此目录包含 APISIX 和 Higress API 网关的配置。

## 支持的 API 网关

### APISIX（推荐）

Apache APISIX 是一个云原生、高性能、可扩展的 API 网关。

**特性**：
- 动态路由和负载均衡
- 认证授权（JWT、Keycloak、LDAP、OpenID Connect）
- 限流熔断
- 灰度发布
- 可观测性（Prometheus、SkyWalking）
- 高性能（基于 OpenResty + LuaJIT）
- 丰富的插件生态

### Higress

Higress 是阿里云开源的云原生 API 网关，深度集成 Istio。

**特性**：
- 支持 Ingress 和 Gateway API
- 与 Istio 深度集成
- 支持 Dubbo、gRPC、HTTP
- WAF 防护
- 插件市场

## 快速开始

### 使用 APISIX

#### 1. 安装 APISIX

```bash
# 使用 Helm 安装
helm repo add apisix https://charts.apiseven.com
helm repo update

# 安装 APISIX 和 Dashboard
helm install apisix apisix/apisix \
  --namespace ingress-apisix \
  --create-namespace \
  --set gateway.type=LoadBalancer \
  --set admin.enabled=true \
  --set dashboard.enabled=true

# 等待 Pod 就绪
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=apisix -n ingress-apisix --timeout=300s
```

#### 2. 配置 APISIX 路由

```bash
# 应用 KubeOps 路由配置
kubectl apply -f deploy/apisix/routes.yaml
```

#### 3. 访问 Dashboard

```bash
# 端口转发
kubectl port-forward -n ingress-apisix svc/apisix-dashboard 8080:8080

# 访问 http://localhost:8080
# 默认用户名：admin
# 默认密码：admin
```

### 使用 Higress

#### 1. 安装 Higress

```bash
# 添加 Higress Helm 仓库
helm repo add higress https://higress.io/helm-charts
helm repo update

# 安装 Higress
helm install higress higress/higress \
  --namespace higress-system \
  --create-namespace \
  --set global.enableIstioAPI=true

# 等待 Pod 就绪
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=higress -n higress-system --timeout=300s
```

#### 2. 配置 Higress 路由

```bash
# 应用 KubeOps 路由配置
kubectl apply -f deploy/higress/routes.yaml
```

#### 3. 访问控制台

```bash
# 端口转发
kubectl port-forward -n higress-system svc/higress-console 8080:8080

# 访问 http://localhost:8080
```

## 配置对比

| 特性 | APISIX | Higress |
|------|--------|---------|
| 性能 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| 易用性 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| Istio 集成 | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| 插件生态 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| 国内支持 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

## APISIX 路由配置

### 基础路由

```yaml
apiVersion: apisix.apache.org/v2
kind: ApisixRoute
metadata:
  name: api-route
  namespace: kubeops
spec:
  http:
  - name: api
    match:
      paths:
      - /api/*
    backends:
    - serviceName: api-gateway
      servicePort: 8080
    plugins:
    - name: prometheus
      enable: true
    - name: limit-req
      enable: true
      config:
        rate: 100
        burst: 200
```

### 认证插件

```yaml
apiVersion: apisix.apache.org/v2
kind: ApisixRoute
metadata:
  name: api-auth
  namespace: kubeops
spec:
  http:
  - name: api-with-auth
    match:
      paths:
      - /api/v1/*
    backends:
    - serviceName: api-gateway
      servicePort: 8080
    plugins:
    - name: jwt-auth
      enable: true
      config:
        key: my-secret-key
    - name: openid-connect
      enable: true
      config:
        discovery: https://keycloak.example.com/.well-known/openid-configuration
        client_id: kubeops
        client_secret: xxxxx
```

### 灰度发布

```yaml
apiVersion: apisix.apache.org/v2
kind: ApisixRoute
metadata:
  name: canary-release
  namespace: kubeops
spec:
  http:
  - name: canary
    match:
      paths:
      - /api/v1/*
      headers:
        x-canary:
          exact: "true"
    backends:
    - serviceName: api-gateway-v2
      servicePort: 8080
      weight: 10
    - serviceName: api-gateway-v1
      servicePort: 8080
      weight: 90
```

## Higress 路由配置

### 基础 Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: kubeops-ingress
  namespace: kubeops
  annotations:
    # Higress 注解
    higress.io/destination: api-gateway.kubeops.svc.cluster.local:8080
    # 限流
    higress.io/rate-limit: "100"
    # 超时
    higress.io/proxy-timeout: "30s"
spec:
  ingressClassName: higress
  rules:
  - host: kubeops.local
    http:
      paths:
      - path: /api
        pathType: Prefix
        backend:
          service:
            name: api-gateway
            port:
              number: 8080
```

### 灰度发布

```yaml
apiVersion: networking.higress.io/v1
kind: McpBridge
metadata:
  name: kubeops-mcp
  namespace: kubeops
spec:
  istioResources:
  - apiVersion: networking.istio.io/v1beta1
    kind: VirtualService
    namespace: kubeops
---
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: kubeops-canary
  namespace: kubeops
spec:
  hosts:
  - "api.kubeops.local"
  http:
  - match:
    - headers:
        x-canary:
          exact: "true"
    route:
    - destination:
        host: api-gateway
        subset: v2
  - route:
    - destination:
        host: api-gateway
        subset: v1
      weight: 90
    - destination:
        host: api-gateway
        subset: v2
      weight: 10
```

## 插件配置

### 限流插件

```yaml
# APISIX
apiVersion: apisix.apache.org/v2
kind: ApisixPluginConfig
metadata:
  name: rate-limit
  namespace: kubeops
spec:
  plugins:
  - name: limit-req
    config:
      rate: 100
      burst: 200
      key_type: var
      key: remote_addr
      rejected_code: 429
      rejected_msg: "请求过于频繁，请稍后再试"
```

### CORS 插件

```yaml
apiVersion: apisix.apache.org/v2
kind: ApisixPluginConfig
metadata:
  name: cors
  namespace: kubeops
spec:
  plugins:
  - name: cors
    config:
      allow_origins: https://kubeops.local
      allow_methods: GET,POST,PUT,DELETE,OPTIONS
      allow_headers: DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Authorization
      expose_headers: Content-Length,Content-Range
      max_age: 3600
```

### Prometheus 监控

```yaml
apiVersion: apisix.apache.org/v2
kind: ApisixGlobalPluginConfig
metadata:
  name: prometheus
  namespace: kubeops
spec:
  plugins:
  - name: prometheus
    config:
      prefer_name: true
```

## 与 Istio 集成

### Higress + Istio

Higress 原生支持 Istio API：

```bash
# Higress 会自动识别 Istio 的 VirtualService、DestinationRule 等
kubectl apply -f deploy/istio/base/
kubectl apply -f deploy/istio/traffic/
```

### APISIX + Istio

APISIX 可以作为 Istio 的入口网关：

```yaml
apiVersion: networking.istio.io/v1beta1
kind: Gateway
metadata:
  name: apisix-gateway
spec:
  selector:
    app.kubernetes.io/name: apisix
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "kubeops.local"
```

## 生产环境配置

### 高可用部署

```yaml
# APISIX 多副本
replicas: 3

# 配置
apisix:
  ssl:
    ssl_trusted_certificate: /etc/ssl/certs/ca-bundle.crt
  nginx_config:
    worker_processes: auto
    worker_rlimit_nofile: 20480
    event:
      worker_connections: 10620
```

### 性能优化

```yaml
# 连接池配置
apisix:
  upstream:
    keepalive: 320
    keepalive_requests: 1000
    keepalive_timeout: 60s
```

### 安全加固

```yaml
# WAF 防护
apiVersion: apisix.apache.org/v2
kind: ApisixRoute
metadata:
  name: secure-api
spec:
  http:
  - plugins:
    - name: ip-restriction
      config:
        whitelist:
        - 10.0.0.0/8
        - 192.168.0.0/16
    - name: jwt-auth
      config:
        key: secret-key
```

## 故障排查

### 查看 APISIX 日志

```bash
# 查看 APISIX 日志
kubectl logs -n ingress-apisix -l app.kubernetes.io/name=apisix

# 查看 APISIX 配置
kubectl exec -n ingress-apisix <apisix-pod> -- cat apisix_conf/apisix.yaml
```

### 查看 Higress 日志

```bash
# 查看 Higress 日志
kubectl logs -n higress-system -l app.kubernetes.io/name=higress

# 查看路由配置
kubectl get virtualservices -n kubeops -o yaml
```

### 测试路由

```bash
# 测试 APISIX 路由
curl -H "Host: kubeops.local" http://<apisix-ip>/api/v1/health

# 测试 Higress 路由
curl -H "Host: kubeops.local" http://<higress-ip>/api/v1/health
```

## 卸载

### APISIX

```bash
helm uninstall apisix -n ingress-apisix
kubectl delete namespace ingress-apisix
```

### Higress

```bash
helm uninstall higress -n higress-system
kubectl delete namespace higress-system
```

## 参考资料

- [APISIX 官方文档](https://apisix.apache.org/docs/)
- [Higress 官方文档](https://higress.io/docs/)
- [APISIX Ingress 控制器](https://github.com/apache/apisix-ingress-controller)
- [Higress GitHub](https://github.com/alibaba/higress)

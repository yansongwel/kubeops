# Istio 网格配置

此目录包含 KubeOps 平台的 Istio 服务网格配置。

## 目录结构

```
istio/
├── base/              # 基础配置
│   ├── namespace.yaml
│   ├── gateway.yaml
│   └── virtual-service.yaml
├── auth/              # 认证策略
│   ├── peer-authentication.yaml
│   └── request-authentication.yaml
├── traffic/           # 流量管理
│   ├── destination-rules.yaml
│   └── traffic-splitting.yaml
└── monitoring/        # 可观测性
    └── service-monitor.yaml
```

## 快速开始

### 1. 安装 Istio

```bash
# 下载 Istio
curl -L https://istio.io/downloadIstio | sh -

# 切换到 Istio 目录
cd istio-*

# 添加 istioctl 到 PATH
export PATH=$PWD/bin:$PATH

# 安装 Istio（默认配置）
istioctl install --set profile=demo -y

# 或生产配置
istioctl install --set profile=default -y
```

### 2. 启用自动 Sidecar 注入

```bash
# 为 KubeOps 命名空间启用注入
kubectl label namespace kubeops istio-injection=enabled
```

### 3. 部署 KubeOps

```bash
# 部署应用（Pod 会自动注入 Sidecar）
helm install kubeops deploy/helm/kubeops \
  --namespace kubeops \
  --create-namespace \
  -f deploy/examples/values-with-istio.yaml
```

### 4. 应用 Istio 配置

```bash
# 应用基础配置
kubectl apply -f deploy/istio/base/

# 应用认证策略（可选）
kubectl apply -f deploy/istio/auth/

# 应用流量管理规则（可选）
kubectl apply -f deploy/istio/traffic/
```

## 配置说明

### 网关 (Gateway)

定义网格的入口点：

```yaml
apiVersion: networking.istio.io/v1beta1
kind: Gateway
metadata:
  name: kubeops-gateway
spec:
  selector:
    istio: ingressgateway
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "kubeops.example.com"
```

### 虚拟服务 (VirtualService)

定义流量路由规则：

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: kubeops
spec:
  hosts:
  - "kubeops.example.com"
  gateways:
  - kubeops-gateway
  http:
  - match:
    - uri:
        prefix: /api
    route:
    - destination:
        host: api-gateway
        port:
          number: 8080
```

### 目标规则 (DestinationRule)

定义流量策略：

```yaml
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: api-gateway
spec:
  host: api-gateway
  trafficPolicy:
    loadBalancer:
      simple: ROUND_ROBIN
    connectionPool:
      tcp:
        maxConnections: 100
```

### 对等认证 (PeerAuthentication)

启用 mTLS：

```yaml
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
spec:
  mtls:
    mode: STRICT
```

## 流量管理

### 灰度发布

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: kubeops-canary
spec:
  http:
  - match:
    - headers:
        x-canary:
          exact: "true"
    route:
    - destination:
        host: api-gateway
        subset: v2  # 新版本
  - route:
    - destination:
        host: api-gateway
        subset: v1  # 旧版本
        weight: 90
    - destination:
        host: api-gateway
        subset: v2
        weight: 10  # 10% 流量到新版本
```

### 故障注入

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: kubeops-fault-injection
spec:
  http:
  - fault:
      delay:
        percentage:
          value: 10
        fixedDelay: 5s
    route:
    - destination:
        host: api-gateway
```

### 超时和重试

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: kubeops-timeout
spec:
  http:
  - timeout: 3s
    retries:
      attempts: 3
      perTryTimeout: 2s
    route:
    - destination:
        host: api-gateway
```

## 可观测性

### 分布式追踪

Istio 自动收集追踪数据并发送到 Jaeger 或 Zipkin：

```bash
# 安装 Kiali（可视化工具）
kubectl apply -f samples/addons

# 访问 Kiali
istioctl dashboard kiali
```

### 指标收集

Istio 自动收集指标：

```bash
# 访问 Prometheus
istioctl dashboard prometheus

# 访问 Grafana
istioctl dashboard grafana
```

## 安全

### mTLS 认证

```yaml
# 全局启用 mTLS
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
  namespace: kubeops
spec:
  mtls:
    mode: STRICT
```

### JWT 认证

```yaml
apiVersion: security.istio.io/v1beta1
kind: RequestAuthentication
metadata:
  name: jwt-example
  namespace: kubeops
spec:
  selector:
    matchLabels:
      app: api-gateway
  jwtRules:
  - issuer: "https://keycloak.example.com/auth/realms/kubeops"
    jwks: "https://keycloak.example.com/auth/realms/kubeops/protocol/openid-connect/certs"
```

### 授权策略

```yaml
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: api-gateway-policy
  namespace: kubeops
spec:
  selector:
    matchLabels:
      app: api-gateway
  action: ALLOW
  rules:
  - from:
    - source:
        requestPrincipals: ["*"]
    to:
    - operation:
        methods: ["GET", "POST"]
```

## 生产环境配置

### 资源限制

```yaml
# 配置 Sidecar 资源限制
apiVersion: v1
kind: Pod
metadata:
  annotations:
    sidecar.istio.io/proxyCPU: "200m"
    sidecar.istio.io/proxyMemory: "256Mi"
    sidecar.istio.io/proxyCPULimit: "500m"
    sidecar.istio.io/proxyMemoryLimit: "512Mi"
```

### 性能优化

```yaml
# 减少 Sidecar 拦获范围
apiVersion: v1
kind: Pod
metadata:
  annotations:
    # 只拦截特定端口
    traffic.sidecar.istio.io/includeOutboundPorts: "8080,8443"
    # 排除特定 IP
    traffic.sidecar.istio.io/excludeOutboundIPRanges: "10.0.0.0/8"
```

## 故障排查

### 查看 Sidecar 状态

```bash
# 查看 Pod 是否注入了 Sidecar
kubectl get pods -n kubeops -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.containers[*].name}{"\n"}{end}'

# 查看 Sidecar 配置
istioctl proxy-config pods <pod-name> -n kubeops

# 查看代理状态
istioctl proxy-status -n kubeops
```

### 查看路由配置

```bash
# 查看虚拟服务
istioctl get virtualservices -n kubeops

# 查看目标规则
istioctl get destinationrules -n kubeops

# 查看网关
istioctl get gateways -n kubeops
```

### 日志查看

```bash
# 查看 Istiod 日志
kubectl logs -n istio-system deployment/istiod

# 查看 Sidecar 代理日志
kubectl logs -n kubeops <pod-name> -c istio-proxy
```

## 卸载

```bash
# 删除 KubeOps
helm uninstall kubeops -n kubeops

# 删除 Istio 配置
kubectl delete -f deploy/istio/

# 卸载 Istio
istioctl uninstall --purge -y

# 删除命名空间
kubectl delete namespace kubeops istio-system
```

## 参考资料

- [Istio 官方文档](https://istio.io/latest/docs/)
- [Istio 流量管理](https://istio.io/latest/docs/concepts/traffic-management/)
- [Istio 安全](https://istio.io/latest/docs/concepts/security/)
- [Istio 可观测性](https://istio.io/latest/docs/concepts/observability/)

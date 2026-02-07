# API 网关和服务网格安装指南

本文档介绍如何在 Kubernetes 集群中安装和配置 APISIX/Higress 和 Istio。
API 网关与 Istio 均为可选外部组件，默认不启用；仅在需要统一流量入口或服务治理时安装。

## 目录

- [前置要求](#前置要求)
- [安装 APISIX](#安装-apisix)
- [安装 Higress](#安装-higress)
- [安装 Istio](#安装-istio)
- [部署 KubeOps](#部署-kubeops)
- [验证安装](#验证安装)

## 前置要求

- Kubernetes 1.19+
- kubectl 已配置
- Helm 3.x

## Helm 默认配置

Helm 默认不启用外部网关与 Istio，可按需在 values 中开启：

```yaml
gateway:
  enabled: false
  type: apisix

istio:
  enabled: false
```

## 安装 APISIX

### 方式一：使用 Helm 安装

```bash
# 添加 APISIX Helm 仓库
helm repo add apisix https://charts.apiseven.com
helm repo update

# 创建命名空间
kubectl create namespace ingress-apisix

# 安装 APISIX
helm install apisix apisix/apisix \
  --namespace ingress-apisix \
  --set gateway.type=LoadBalancer \
  --set admin.enabled=true \
  --set admin.credentials.admin="admin" \
  --set dashboard.enabled=true

# 等待 Pod 就绪
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=apisix -n ingress-apisix --timeout=300s
```

### 方式二：使用 APISIX Ingress 控制器

```bash
# 安装 APISIX 和 Ingress 控制器
helm install apisix apisix/ingress-controller \
  --namespace ingress-apisix \
  --set apisix.enabled=true \
  --set ingress-controller.enabled=true \
  --set config.apisix.serviceNamespace=ingress-apisix

# 应用路由配置
kubectl apply -f deploy/gateway/apisix/routes.yaml
```

### 访问 APISIX Dashboard

```bash
# 端口转发
kubectl port-forward -n ingress-apisix svc/apisix-dashboard 8080:8080

# 访问 http://localhost:8080
# 默认用户名：admin
# 默认密码：admin
```

### 配置 APISIX 路由

```bash
# 应用 KubeOps 路由
kubectl apply -f deploy/gateway/apisix/routes.yaml

# 验证路由
kubectl get apisixroute -n kubeops
```

## 安装 Higress

### 使用 Helm 安装

```bash
# 添加 Higress Helm 仓库
helm repo add higress https://higress.io/helm-charts
helm repo update

# 创建命名空间
kubectl create namespace higress-system

# 安装 Higress
helm install higress higress/higress \
  --namespace higress-system \
  --set global.enableIstioAPI=true \
  --set console.enabled=true

# 等待 Pod 就绪
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=higress -n higress-system --timeout=300s
```

### 访问 Higress 控制台

```bash
# 端口转发
kubectl port-forward -n higress-system svc/higress-console 8080:8080

# 访问 http://localhost:8080
```

### 配置 Higress 路由

```bash
# 应用 KubeOps 路由
kubectl apply -f deploy/gateway/higress/routes.yaml

# 验证路由
kubectl get ingress -n kubeops
kubectl get virtualservice -n kubeops
```

## 安装 Istio

### 下载 Istio

```bash
# 下载 Istio
curl -L https://istio.io/downloadIstio | sh -

# 切换到 Istio 目录
cd istio-*

# 添加 istioctl 到 PATH
export PATH=$PWD/bin:$PATH
```

### 安装 Istio

```bash
# 查看可用的配置 profile
istioctl profile list

# 安装 Istio（demo 配置适合开发）
istioctl install --set profile=demo -y

# 生产环境使用 default profile
istioctl install --set profile=default -y

# 验证安装
istioctl verify-install
```

### 启用 Sidecar 注入

```bash
# 为 KubeOps 命名空间启用自动注入
kubectl label namespace kubeops istio-injection=enabled

# 验证标签
kubectl get namespace kubeops -L istio-injection
```

### 安装 Kiali（可选）

```bash
# 安装 Kiali 和其他插件
kubectl apply -f samples/addons

# 访问 Kiali
istioctl dashboard kiali
```

### 应用 Istio 配置

```bash
# 应用基础配置
kubectl apply -f deploy/istio/base/

# 应用认证策略
kubectl apply -f deploy/istio/auth/

# 应用流量管理规则（可选）
kubectl apply -f deploy/istio/traffic/
```

## 部署 KubeOps

### 使用 APISIX

```bash
# 部署 KubeOps
helm install kubeops deploy/helm/kubeops \
  --namespace kubeops \
  --create-namespace \
  --set gateway.type=apisix \
  -f deploy/examples/values-with-gateway.yaml

# 部署 APISIX 路由
kubectl apply -f deploy/gateway/apisix/routes.yaml
```

### 使用 Higress

```bash
# 部署 KubeOps
helm install kubeops deploy/helm/kubeops \
  --namespace kubeops \
  --create-namespace \
  --set gateway.type=higress \
  -f deploy/examples/values-with-gateway.yaml

# 部署 Higress 路由
kubectl apply -f deploy/gateway/higress/routes.yaml
```

### 使用 Istio

```bash
# 部署 KubeOps
helm install kubeops deploy/helm/kubeops \
  --namespace kubeops \
  --create-namespace \
  --set istio.enabled=true \
  -f deploy/examples/values-with-gateway.yaml

# 应用 Istio 配置
kubectl apply -f deploy/istio/base/
kubectl apply -f deploy/istio/auth/
kubectl apply -f deploy/istio/traffic/
```

## 验证安装

### 检查 APISIX

```bash
# 查看 APISIX Pod
kubectl get pods -n ingress-apisix

# 查看 APISIX 服务
kubectl get svc -n ingress-apisix

# 查看 APISIX 路由
kubectl get apisixroute -n kubeops

# 测试路由
export APISIX_IP=$(kubectl get svc -n ingress-apisix apisix-gateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
curl -H "Host: api.kubeops.local" http://$APISIX_IP/health
```

### 检查 Higress

```bash
# 查看 Higress Pod
kubectl get pods -n higress-system

# 查看 Higress 服务
kubectl get svc -n higress-system

# 查看 Ingress
kubectl get ingress -n kubeops

# 测试路由
export HIGRESS_IP=$(kubectl get svc -n higress-system higress-gateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
curl -H "Host: api.kubeops.local" http://$HIGRESS_IP/api/v1/health
```

### 检查 Istio

```bash
# 查看 Istio Pod
kubectl get pods -n istio-system

# 查看 Sidecar 注入
kubectl get pods -n kubeops -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.containers[*].name}{"\n"}{end}'

# 查看虚拟服务
kubectl get virtualservices -n kubeops

# 查看目标规则
kubectl get destinationrules -n kubeops

# 访问 Kiali
istioctl dashboard kiali
```

### 检查 KubeOps

```bash
# 查看 KubeOps Pod
kubectl get pods -n kubeops

# 查看服务
kubectl get svc -n kubeops

# 查看日志
kubectl logs -n kubeops -l app.kubernetes.io/name=kubeops --all-containers=true

# 测试 API
export GATEWAY_IP=$(kubectl get svc -n ingress-apisix apisix-gateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
curl -H "Host: api.kubeops.local" http://$GATEWAY_IP/api/v1/namespaces
```

## 配置对比

### APISIX vs Higress vs Istio

| 特性 | APISIX | Higress | Istio |
|------|--------|---------|-------|
| **部署复杂度** | 简单 | 简单 | 中等 |
| **性能** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **路由功能** | 丰富 | 丰富 | 非常丰富 |
| **服务网格** | 无 | 基础 | 完整 |
| **可观测性** | 好 | 好 | 优秀 |
| **学习曲线** | 平缓 | 平缓 | 陡峭 |
| **社区支持** | 活跃 | 活跃 | 非常活跃 |
| **适用场景** | API 网关 | API 网关 | 微服务 |

### 推荐组合

1. **小型项目**：APISIX 单独使用
2. **中型项目**：Higress + Istio（基础功能）
3. **大型项目**：APISIX（入口） + Istio（服务网格）

## 卸载

### 卸载 APISIX

```bash
# 删除 KubeOps
helm uninstall kubeops -n kubeops

# 删除路由
kubectl delete -f deploy/gateway/apisix/

# 卸载 APISIX
helm uninstall apisix -n ingress-apisix
kubectl delete namespace ingress-apisix
```

### 卸载 Higress

```bash
# 删除 KubeOps
helm uninstall kubeops -n kubeops

# 删除路由
kubectl delete -f deploy/gateway/higress/

# 卸载 Higress
helm uninstall higress -n higress-system
kubectl delete namespace higress-system
```

### 卸载 Istio

```bash
# 删除 KubeOps
helm uninstall kubeops -n kubeops

# 删除 Istio 配置
kubectl delete -f deploy/istio/

# 卸载 Istio
istioctl uninstall --purge -y
kubectl delete namespace istio-system
```

## 故障排查

### APISIX 问题

```bash
# 查看日志
kubectl logs -n ingress-apisix -l app.kubernetes.io/name=apisix

# 查看配置
kubectl exec -n ingress-apisix <apisix-pod> -- cat apisix_conf/apisix.yaml

# 重新加载配置
kubectl exec -n ingress-apisix <apisix-pod> -- apisix reload
```

### Higress 问题

```bash
# 查看日志
kubectl logs -n higress-system -l app.kubernetes.io/name=higress

# 查看 Ingress 配置
kubectl describe ingress -n kubeops

# 查看 VirtualService
kubectl get virtualservice -n kubeops -o yaml
```

### Istio 问题

```bash
# 查看 Istiod 日志
kubectl logs -n istio-system deployment/istiod

# 查看 Sidecar 代理状态
istioctl proxy-status -n kubeops

# 查看 Sidecar 配置
istioctl proxy-config pods <pod-name> -n kubeops

# 重启 Sidecar
kubectl rollout restart deployment/<deployment-name> -n kubeops
```

## 最佳实践

1. **开发环境**：使用 APISIX，部署简单
2. **测试环境**：使用 Higress，验证 Istio 集成
3. **生产环境**：APISIX（入口） + Istio（服务网格）
4. **监控**：集成 Prometheus + Grafana
5. **日志**：集成 ELK 或 Loki
6. **追踪**：启用 Jaeger 分布式追踪

## 参考资源

- [APISIX 官方文档](https://apisix.apache.org/docs/)
- [Higress 官方文档](https://higress.io/docs/)
- [Istio 官方文档](https://istio.io/latest/docs/)
- [KubeOps 架构文档](../architecture/README.md)

# KubeOps

企业级 Kubernetes 统一管理平台，集成 AI 智能巡检、DevOps/CI/CD、日志平台和监控平台。

## 🎯 项目愿景

KubeOps 致力于成为最全面、可扩展性最强的 Kubernetes 管理平台，提供：

- ✅ **完整的 K8s 资源管理**：管理所有 Kubernetes 资源和 CRD，支持扩展
- 🤖 **AI 智能巡检**：智能集群健康分析和优化建议
- 🚀 **集成 DevOps**：CI/CD 流水线，集成 Jenkins + ArgoCD
- 📊 **灵活的日志方案**：ELK 或 Vector+Loki 方案（可扩展到 Graylog 等）
- 📈 **灵活的监控方案**：Prometheus 或 VictoriaMetrics（可扩展）
- 🌐 **开源免费**：适配所有企业需求

## 🏗️ 系统架构

KubeOps 采用**模块化单体架构**，后端以单个进程运行，内部按领域拆分模块，便于演进和扩展；API 网关与服务网格可作为可选的外部组件接入：

```
┌─────────────────────────────────────────────────────────────┐
│                   可选：APISIX / Higress                     │
│              (API 网关、认证、限流、路由)                      │
└─────────────────────────────────────────────────────────────┘
                             │
┌─────────────────────────────────────────────────────────────┐
│                    可选：Istio 服务网格                        │
│              (流量管理、安全、可观测性)                       │
└─────────────────────────────────────────────────────────────┘
                             │
┌─────────────────────────────────────────────────────────────┐
│                     KubeOps 单体后端                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐       │
│  │ Kube 管理模块 │  │ AI 巡检模块  │  │ DevOps 模块  │       │
│  └──────────────┘  └──────────────┘  └──────────────┘       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐       │
│  │ 日志集成模块 │  │ 监控集成模块 │  │ API 网关模块 │       │
│  └──────────────┘  └──────────────┘  └──────────────┘       │
└─────────────────────────────────────────────────────────────┘
                             │
┌─────────────────────────────────────────────────────────────┐
│                           前端界面                            │
│                          (Vue 3)                             │
└─────────────────────────────────────────────────────────────┘
```

### 核心模块

| 模块 | 描述 | 技术栈 |
|---------|-------------|------------|
| **API 网关模块** | 统一入口、认证、限流、路由 | Go + Gin |
| **Kube 管理模块** | K8s 资源管理 | Go + controller-runtime |
| **AI 巡检模块** | AI 智能分析 | Go + OpenAI/Ollama |
| **DevOps 模块** | Jenkins + ArgoCD 集成 | Go |
| **日志集成模块** | 日志聚合 (ELK/Loki) | Go + Vector |
| **监控集成模块** | 指标和告警 | Go + Prometheus |
| **前端界面** | Web 仪表盘 | Vue 3 + TypeScript |

## 🚀 快速开始

### 前置要求

- Go 1.25
- Node.js 24+
- Docker/Podman
- kubectl
- helm
- kind（本地开发）或可访问的 K8s 集群（1.33+）

### 本地开发

1. **克隆仓库**
```bash
git clone https://github.com/yansongwel/kubeops.git
cd kubeops
```

2. **启动依赖服务**
```bash
docker-compose -f deploy/docker-compose-dev.yaml up -d
```

3. **运行后端服务**
```bash
# 单体后端
cd backend
go run cmd/server/main.go
```

4. **运行前端**
```bash
cd frontend
npm install
npm run dev
```

5. **访问仪表盘**
```
http://localhost:5173
```

### 部署到 Kubernetes

```bash
# 创建命名空间
kubectl create namespace kubeops

# 使用 Helm 安装
helm install kubeops deploy/helm/kubeops \
  --namespace kubeops \
  --set logging.stack=elk \
  --set monitoring.stack=prometheus

# 或使用自定义配置
helm install kubeops deploy/helm/kubeops \
  --namespace kubeops \
  -f deploy/examples/values-production.yaml
```

## 📦 功能特性

### 🔧 Kubernetes 资源管理
- 管理所有 K8s 原生资源（Deployment、Service、ConfigMap 等）
- 完整的 CRD 支持，动态模式处理
- 实时资源监控和事件流
- 资源验证和策略执行
- 多集群管理

### 🤖 AI 智能巡检
- **资源优化**：检测资源配置过高/过低
- **安全分析**：漏洞扫描、RBAC 审计
- **可靠性检查**：健康检查、SLO/SLI 分析
- **成本优化**：识别节省成本的机会
- **预测分析**：容量规划预测
- **日志分析**：基于 NLP 的日志分析

### 🚀 DevOps & CI/CD
- **Jenkins 集成**：创建、触发、监控 Jenkins 任务
- **ArgoCD 集成**：GitOps 部署、应用同步
- **代码质量**：SonarQube、测试覆盖率集成
- **流水线可视化**：实时流水线状态
- **环境管理**：开发、测试、生产环境管理

### 📊 日志平台（可插拔）
**方案 1：ELK Stack**
- Vector + Fluentd → Elasticsearch → Kibana

**方案 2：Loki Stack**
- Vector → VictoriaLogs → Loki → Grafana

**可扩展**：易于添加 Graylog、Watchalert 等

### 📈 监控平台（可插拔）
**方案 1：Prometheus**
- Prometheus + Grafana
- 自定义告警规则
- 仪表盘模板

**方案 2：VictoriaMetrics**
- Prometheus + VictoriaMetrics（集群版）
- 长期指标存储
- 高性能查询

## 🌐 API 网关与服务网格

### API 网关（可选）

**APISIX**（推荐）：
- 云原生 API 网关
- 动态路由、负载均衡
- 认证授权（JWT、Keycloak、LDAP）
- 限流熔断、灰度发布
- 高性能（基于 OpenResty + Lua）

**Higress**：
- 阿里云开源的云原生 API 网关
- 支持 Ingress 和 Gateway API
- 与 Istio 深度集成
- 插件生态丰富

### 服务网格（Istio）

**核心功能**：
- **流量管理**：智能路由、灰度发布、故障注入
- **安全**：mTLS 认证、细粒度授权
- **可观测性**：分布式追踪、指标收集
- **策略执行**：访问控制、速率限制

**集成特性**：
- 自动 Sidecar 注入
- 虚拟服务和目标规则
- 流量分裂和 A/B 测试
- 超时和重试策略

## 🔐 安全特性

- JWT/OAuth2 认证
- 基于角色的访问控制（RBAC）
- 网络策略
- 密钥管理（Kubernetes secrets、Vault 集成）
- 审计日志
- 安全通信（TLS/mTLS）

## 🧪 测试

```bash
# 后端测试
cd backend
make test-unit        # 单元测试
make test-integration # 集成测试
make test-e2e         # 端到端测试

# 前端测试
cd frontend
npm run test:unit     # 单元测试
npm run test:integration # 集成测试
npm run test:e2e      # 端到端测试
```

## 📚 文档

- [架构概述](docs/architecture/README.md)
- [快速开始](docs/QUICKSTART.md)
- [API 文档](docs/api/README.md)
- [开发指南](docs/development/README.md)
- [部署指南](docs/deployment/README.md)
- [贡献指南](CONTRIBUTING.md)

## 🗺️ 开发路线图

### 第一阶段：基础架构（当前）
- [x] 项目架构和初始化
- [ ] API 网关和认证
- [ ] Kube 管理器（基础 K8s 资源）
- [ ] 前端仪表盘框架

### 第二阶段：核心功能
- [ ] 完整的 K8s 资源管理
- [ ] AI 巡检（基础分析）
- [ ] 多集群支持

### 第三阶段：DevOps 集成
- [ ] Jenkins 集成
- [ ] ArgoCD 集成
- [ ] 流水线可视化

### 第四阶段：可观测性
- [ ] 日志服务（ELK + Loki）
- [ ] 监控服务（Prometheus + VictoriaMetrics）
- [ ] 告警管理

### 第五阶段：高级功能
- [ ] 高级 AI 能力
- [ ] 自定义 CRD 构建器
- [ ] 策略引擎（OPA 集成）
- [ ] 成本管理

## 🤝 贡献指南

我们欢迎贡献！请查看 [贡献指南](CONTRIBUTING.md)。

### 开发流程

1. Fork 仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'feat: 添加某个功能'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 📄 开源协议

Apache License 2.0 - 详见 [LICENSE](LICENSE)

## 🙏 致谢

- [Kubernetes](https://kubernetes.io/)
- [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime)
- [Vue.js](https://vuejs.org/)
- [Gin](https://gin-gonic.com/)
- [Prometheus](https://prometheus.io/)
- [Grafana](https://grafana.com/)

## 📞 联系方式

- 网站：https://kubeops.io
- 文档：https://docs.kubeops.io
- GitHub：https://github.com/yansongwel/kubeops
- 邮箱：hello@kubeops.io

---

**注意**：本项目处于早期开发阶段，并非所有功能都已实现。请查看 [开发路线图](#-开发路线图) 了解进展。

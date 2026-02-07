# KubeOps 项目初始化完成 🎉

## 项目定位

**KubeOps** 是一个**全栈学习项目**，面向希望提升前后端开发能力的高级运维工程师。

**核心决策**：自己实现前端，最大化学习价值。

### 核心特性

- ✅ **K8s 资源管理**：管理所有 Kubernetes 资源和 CRD
- ✅ **AI 智能巡检**：集群健康分析和优化建议
- ✅ **DevOps 集成**：Jenkins + ArgoCD
- ✅ **日志平台**：ELK 或 Loki 方案（可插拔）
- ✅ **监控平台**：Prometheus 或 VictoriaMetrics（可插拔）
- ✅ **微服务架构**：Go 后端 + Vue 3 前端

## 技术栈

### 后端
- **语言**：Go 1.25
- **框架**：Gin v1.11.0 (HTTP), gRPC v1.70.0 (内部通信)
- **K8s 集成**：controller-runtime v0.20.0, client-go v0.35.0 (Kubernetes 1.33)
- **数据库**：PostgreSQL + Redis
- **日志**：Zap v1.27.0
- **仓库**：https://github.com/yansongwel/kubeops.git

### 前端
- **框架**：Vue 3.5 + Composition API
- **语言**：TypeScript 5.7
- **UI 库**：Element Plus 2.9
- **构建工具**：Vite 6.0
- **Node 版本**：24.0+

### 基础设施
- **容器**：Docker / Containerd
- **编排**：Kubernetes
- **部署**：Helm Charts
- **本地开发**：kind (Kubernetes in Docker)

## 项目结构

```
KubeOps/
├── backend/                    # Go 后端微服务
│   ├── api-gateway/           # API 网关 ✅
│   ├── kube-manager/          # K8s 资源管理 ✅
│   ├── ai-inspector/          # AI 巡检服务 🚧
│   ├── devops-service/        # DevOps 集成 🚧
│   ├── logging-service/       # 日志平台 🚧
│   ├── monitoring-service/    # 监控平台 🚧
│   ├── common/                # 共享库 🚧
│   ├── proto/                 # gRPC 定义 🚧
│   ├── Dockerfile            # 多阶段构建 ✅
│   └── Makefile              # 构建脚本 ✅
│
├── frontend/                   # Vue 3 前端
│   ├── src/
│   │   ├── views/            # 页面组件 ✅
│   │   ├── layouts/          # 布局组件 ✅
│   │   ├── router/           # 路由配置 ✅
│   │   └── main.ts           # 入口文件 ✅
│   ├── package.json          # 依赖配置 ✅
│   ├── vite.config.ts        # Vite 配置 ✅
│   ├── Dockerfile            # 前端构建 ✅
│   └── nginx.conf            # Nginx 配置 ✅
│
├── deploy/                     # 部署配置
│   ├── docker-compose-dev.yaml  # 本地开发 ✅
│   ├── helm/kubeops/         # Helm Chart ✅
│   │   ├── Chart.yaml        # Chart 定义 ✅
│   │   ├── values.yaml       # 默认配置 ✅
│   │   └── templates/        # K8s 模板 ✅
│   └── examples/             # 示例配置 ✅
│       ├── values-dev.yaml
│       └── values-production.yaml
│
├── scripts/                    # 开发脚本
│   ├── dev.sh                # Linux/Mac 脚本 ✅
│   ├── setup.ps1             # Windows 脚本 ✅
│   └── kind-config.yaml      # Kind 配置 ✅
│
├── docs/                       # 文档
│   ├── architecture/README.md # 架构文档 ✅
│   ├── QUICKSTART.md         # 快速开始 ✅
│   └── DOCUMENTATION.md      # 文档规范 ✅
│
├── README.md                   # 项目说明 ✅
├── CLAUDE.md                   # Claude Code 指南 ✅
├── CONTRIBUTING.md             # 贡献指南 ✅
├── LICENSE                     # Apache 2.0 ✅
└── go.mod                      # Go 模块 ✅

✅ 已完成  🚧 待开发
```

## 已完成工作

### 1. 项目架构 ✅
- [x] 微服务架构设计
- [x] 目录结构规划
- [x] 技术栈选型
- [x] 开发规范制定

### 2. 后端服务 ✅
- [x] API 网关基础代码
- [x] Kube 管理器基础代码
- [x] Docker 多阶段构建
- [x] Makefile 构建脚本
- [x] Go 模块配置

### 3. 前端应用 ✅
- [x] Vue 3 项目搭建
- [x] 路由配置
- [x] 主布局组件
- [x] 页面组件框架
- [x] Element Plus 集成
- [x] TypeScript 配置

### 4. 部署配置 ✅
- [x] Docker Compose 开发环境
- [x] Helm Chart 结构
- [x] K8s 部署模板
- [x] Kind 集群配置
- [x] 开发/生产环境配置

### 5. 文档中文化 ✅
- [x] README.md 中文化
- [x] 架构文档中文化
- [x] 贡献指南中文化
- [x] 代码注释中文化
- [x] 用户界面中文化

### 6. 开发工具 ✅
- [x] Linux/Mac 开发脚本
- [x] Windows 设置脚本
- [x] Gitignore 配置

## 快速开始

### 本地开发

```bash
# 1. 启动依赖服务
docker-compose -f deploy/docker-compose-dev.yaml up -d

# 2. 运行后端服务
cd backend/api-gateway && go run cmd/server/main.go
cd backend/kube-manager && go run cmd/server/main.go

# 3. 运行前端
cd frontend && npm install && npm run dev

# 4. 访问
open http://localhost:5173
```

### Kind 集群部署（基础版）

```bash
# 1. 创建集群
./scripts/dev.sh kind

# 2. 部署应用
helm install kubeops deploy/helm/kubeops \
  --namespace kubeops \
  --create-namespace \
  -f deploy/examples/values-dev.yaml

# 3. 访问
kubectl port-forward -n kubeops svc/kubeops 8080:80
open http://localhost:8080
```

### 使用 APISIX 部署（推荐）

```bash
# 1. 安装 APISIX
helm repo add apisix https://charts.apiseven.com
helm install apisix apisix/apisix -n ingress-apisix --create-namespace

# 2. 部署 KubeOps
helm install kubeops deploy/helm/kubeops \
  --namespace kubeops \
  --create-namespace \
  --set gateway.type=apisix \
  -f deploy/examples/values-with-gateway.yaml

# 3. 配置路由
kubectl apply -f deploy/gateway/apisix/routes.yaml
```

### 使用 Higress + Istio 部署（高级）

```bash
# 1. 安装 Istio
istioctl install --set profile=demo -y

# 2. 启用 Sidecar 注入
kubectl label namespace kubeops istio-injection=enabled

# 3. 安装 Higress
helm repo add higress https://higress.io/helm-charts
helm install higress higress/higress -n higress-system --create-namespace

# 4. 部署 KubeOps
helm install kubeops deploy/helm/kubeops \
  --namespace kubeops \
  --create-namespace \
  --set gateway.type=higress \
  --set istio.enabled=true \
  -f deploy/examples/values-with-gateway.yaml

# 5. 应用 Istio 配置
kubectl apply -f deploy/istio/base/
kubectl apply -f deploy/istio/auth/
```

## 开发路线图

### 第一阶段：基础架构（当前 - 50% 完成）
- [x] 项目架构和初始化
- [ ] API 网关完整实现（JWT 认证、RBAC）
- [ ] Kube 管理器完整实现（所有 K8s 资源）
- [ ] 前端仪表盘完善

### 第二阶段：核心功能（0% 完成）
- [ ] 完整的 K8s 资源管理
- [ ] AI 巡检基础分析
- [ ] 多集群支持

### 第三阶段：DevOps 集成（0% 完成）
- [ ] Jenkins 集成
- [ ] ArgoCD 集成
- [ ] 流水线可视化

### 第四阶段：可观测性（0% 完成）
- [ ] 日志服务（ELK + Loki）
- [ ] 监控服务（Prometheus + VictoriaMetrics）
- [ ] 告警管理

### 第五阶段：高级功能（0% 完成）
- [ ] 高级 AI 能力
- [ ] 自定义 CRD 构建器
- [ ] 策略引擎（OPA 集成）
- [ ] 成本管理

## 核心服务状态

| 服务 | 状态 | 说明 |
|------|------|------|
| **API 网关层** | ✅ 已配置 | APISIX / Higress 配置完成 |
| APISIX | ✅ 已配置 | Helm Chart 和路由配置完成 |
| Higress | ✅ 已配置 | Helm Chart 和路由配置完成 |
| **服务网格** | ✅ 已配置 | Istio 配置完成 |
| Istio | ✅ 已配置 | 基础配置、认证策略、流量管理 |
| **应用服务** | 🚧 开发中 | 后端微服务 |
| API 网关 | 🚧 开发中 | 基础框架完成，需实现认证/授权 |
| Kube 管理器 | 🚧 开发中 | 基础 API 完成，需扩展所有资源 |
| AI 巡检 | 📝 规划中 | 需实现数据采集和分析引擎 |
| DevOps 服务 | 📝 规划中 | 需集成 Jenkins 和 ArgoCD |
| 日志服务 | 📝 规划中 | 需实现日志路由和提供商接口 |
| 监控服务 | 📝 规划中 | 需实现指标查询和提供商接口 |
| 前端界面 | 🚧 开发中 | 框架完成，需实现各功能页面 |

## 下一步工作

### 立即可做
1. **完善 API 网关**
   - 实现 JWT 认证中间件
   - 实现 RBAC 授权
   - 添加请求限流
   - 实现 API 版本控制

2. **扩展 Kube 管理器**
   - 实现所有 K8s 资源的 CRUD
   - 添加 CRD 动态发现
   - 实现 WebSocket 实时推送
   - 添加资源验证

3. **完善前端界面**
   - 实现集群管理页面
   - 实现工作负载页面
   - 实现 Pod 列表和详情
   - 实现实时日志查看

### 需要设计
1. **数据库设计**
   - 用户表
   - 角色表
   - 权限表
   - 集群表
   - 审计日志表

2. **API 设计**
   - RESTful API 规范
   - 错误码定义
   - 响应格式统一

3. **认证流程**
   - JWT Token 生成和验证
   - 刷新 Token 机制
   - 权限检查中间件

## 文档资源

### 核心文档
- **README**：`README.md` - 项目主文档
- **架构文档**：`docs/architecture/README.md` - 系统架构详解
- **快速开始**：`docs/QUICKSTART.md` - 5分钟快速开始
- **贡献指南**：`CONTRIBUTING.md` - 贡献指南
- **文档规范**：`docs/DOCUMENTATION.md` - 文档语言规范

### 前端学习（⭐ 重点）
- **学习指南**：`docs/frontend/learning-guide.md` - 前端技能学习路径
- **学习计划**：`docs/frontend/learning-plan.md` - 6 个月详细计划
- **前端架构**：`docs/frontend/architecture.md` - 前端项目架构
- **代码模板**：`docs/frontend/code-templates.md` - 常用代码模板
- **自己实现**：`docs/frontend/self-implementation.md` - 自研前端说明

### 安装配置
- **API 网关和服务网格**：`docs/installation/gateway-istio.md` - APISIX/Higress/Istio 安装指南
- **APISIX 配置**：`deploy/gateway/apisix/` - APISIX 路由配置
- **Higress 配置**：`deploy/gateway/higress/` - Higress 路由配置
- **Istio 配置**：`deploy/istio/` - Istio 服务网格配置

### 配置示例
- **开发环境**：`deploy/examples/values-dev.yaml`
- **生产环境**：`deploy/examples/values-production.yaml`
- **带网关配置**：`deploy/examples/values-with-gateway.yaml`

## 开发规范

### 代码规范
- **后端**：Go 标准规范 + 中文注释
- **前端**：Vue 3 Composition API + TypeScript
- **提交信息**：约定式提交（中文）

### 文档规范
- **语言**：中文为主要文档语言
- **代码**：标识符英文，注释中文
- **API**：端点英文，文档中文

## 贡献指南

欢迎贡献！请查看 `CONTRIBUTING.md` 了解详细信息。

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'feat: 添加某个功能'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

Apache License 2.0 - 详见 `LICENSE` 文件

## 联系方式

- GitHub：https://github.com/yansongwel/kubeops.git
- 邮箱：hello@kubeops.io

---

**当前版本**：v0.2.0-alpha
**最后更新**：2026-02-07
**项目状态**：🚧 开发中
**最新特性**：✅ APISIX/Higress + Istio 集成

感谢您对 KubeOps 的关注！🎉

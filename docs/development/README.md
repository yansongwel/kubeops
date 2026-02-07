# KubeOps 开发指南

本文档介绍 KubeOps 项目的开发环境设置、代码规范和最佳实践。

---

## 环境准备

### 必需工具

| 工具 | 版本要求 | 用途 |
|------|----------|------|
| **Go** | 1.25+ | 后端开发 |
| **Node.js** | 24.0+ | 前端开发 |
| **Docker** | 最新版 | 容器化 |
| **kubectl** | 1.33+ | K8s 集群管理 |
| **helm** | 3.x+ | K8s 部署 |

### 可选工具

| 工具 | 版本要求 | 用途 |
|------|----------|------|
| **kind** | 最新版 | 本地 K8s 集群 |
| **golangci-lint** | 最新版 | Go 代码检查 |
| **buf** | 最新版 | Protobuf 代码生成 |
| **make** | 任意版本 | 构建自动化 |

---

## 本地开发设置

### 1. 克隆项目

```bash
git clone https://github.com/yansongwel/kubeops.git
cd kubeops
```

### 2. 启动依赖服务

使用 Docker Compose 启动 PostgreSQL 和 Redis：

```bash
docker-compose -f deploy/docker-compose-dev.yaml up -d
```

### 3. 运行后端服务

**方式一：直接运行（单体架构）**

```bash
cd backend
go run cmd/server/main.go
```

**方式二：使用 make**

```bash
cd backend
make dev
```

### 4. 运行前端

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

访问 http://localhost:5173

---

## 后端开发

### 目录结构

```
backend/
├── cmd/
│   └── server/           # 单体入口
├── internal/
│   ├── gateway/          # API 网关模块
│   ├── kube/             # K8s 资源管理模块
│   ├── ai/               # AI 巡检模块
│   ├── devops/           # DevOps 集成模块
│   ├── logging/          # 日志集成模块
│   ├── monitoring/       # 监控集成模块
│   ├── plugin/           # 插件适配层
│   └── common/           # 共享库
└── pkg/                  # 公共包
```

### 分层架构

KubeOps 使用清晰的三层架构：

```
Request → Handler → Service → Repository → K8s API
  ↓         ↓          ↓            ↓          ↓
HTTP     参数验证    业务逻辑    数据访问    K8s集群
```

**各层职责**

- **Handler 层**：处理 HTTP 请求，参数验证，调用 Service
- **Service 层**：业务逻辑，数据处理，调用 Repository
- **Repository 层**：数据访问，与 K8s API 交互

### 添加新的 API 端点

**示例：添加 Deployment 列表 API**

#### 1. 创建 Repository

创建 `backend/internal/kube/repository/deployment.go`：

```go
package repository

import (
	"context"
	"k8s.io/client-go/kubernetes"
)

type DeploymentRepository struct {
	client *kubernetes.Clientset
}

func NewDeploymentRepository(client *kubernetes.Clientset) *DeploymentRepository {
	return &DeploymentRepository{client: client}
}

func (r *DeploymentRepository) ListByNamespace(namespace string) ([]interface{}, error) {
	deployments, err := r.client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// 转换为业务模型
	var result []interface{}
	for _, dep := range deployments.Items {
		result = append(result, dep)
	}
	return result, nil
}
```

#### 2. 创建 Service

创建 `backend/internal/kube/service/deployment.go`：

```go
package service

import "github.com/yansongwel/kubeops/backend/internal/kube/repository"

type DeploymentService struct {
	repo *repository.DeploymentRepository
}

func NewDeploymentService(repo *repository.DeploymentRepository) *DeploymentService {
	return &DeploymentService{repo: repo}
}

func (s *DeploymentService) ListDeployments(namespace string) ([]string, error) {
	deployments, err := s.repo.ListByNamespace(namespace)
	if err != nil {
		return nil, err
	}

	// 业务逻辑：提取 Deployment 名称
	var result []string
	for _, dep := range deployments {
		result = append(result, dep.Name)
	}
	return result, nil
}
```

#### 3. 创建 Handler

创建 `backend/internal/kube/handler/deployment.go`：

```go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yansongwel/kubeops/backend/internal/kube/service"
)

type DeploymentHandler struct {
	service *service.DeploymentService
}

func NewDeploymentHandler(svc *service.DeploymentService) *DeploymentHandler {
	return &DeploymentHandler{service: svc}
}

func (h *DeploymentHandler) ListDeployments(c *gin.Context) {
	namespace := c.Param("namespace")

	deployments, err := h.service.ListDeployments(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list deployments",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": deployments,
	})
}
```

#### 4. 注册路由

在 `backend/cmd/server/main.go` 中添加路由：

```go
// API v1 路由组
v1 := router.Group("/api/v1")
{
	// ... 现有路由

	// Deployment 相关路由
	v1.GET("/namespaces/:namespace/deployments", deploymentHandler.ListDeployments)
}
```

### 代码规范

- **遵循 `gofmt`**：所有代码必须通过 `gofmt` 格式化
- **使用中文注释**：导出的函数和类型必须有中文注释
- **错误处理**：不要忽略错误，使用 `if err != nil` 处理
- **命名规范**：
  - 文件名：小写，如 `user_handler.go`
  - 包名：小写单词，如 `handler`, `service`
  - 导出函数：大驼峰，如 `ListUsers`
  - 私有函数：小驼峰，如 `validateUser`

---

## 前端开发

### 目录结构

```
frontend/src/
├── api/              # API 调用封装
│   └── kube/
│       ├── cluster.ts
│       ├── namespace.ts
│       └── pod.ts
├── components/       # 可复用组件
│   └── common/
├── composables/      # Vue 3 Composables
├── layouts/          # 布局组件
├── router/           # 路由配置
├── stores/           # Pinia 状态管理
├── types/            # TypeScript 类型定义
├── utils/            # 工具函数
├── views/            # 页面组件
├── App.vue
└── main.ts
```

### 添加新页面

**示例：添加 Deployments 页面**

#### 1. 创建 API 调用

创建 `frontend/src/api/kube/deployment.ts`：

```typescript
import request from '@/utils/request'
import type { Deployment } from '@/types/kube'

export function getDeployments(namespace: string) {
  return request.get<Deployment[]>(`/namespaces/${namespace}/deployments`)
}

export function getDeployment(namespace: string, name: string) {
  return request.get<Deployment>(`/namespaces/${namespace}/deployments/${name}`)
}
```

#### 2. 创建类型定义

在 `frontend/src/types/kube.ts` 中添加：

```typescript
export interface Deployment {
  name: string
  namespace: string
  ready: number
  desired: number
  updated: number
  available: number
  age: string
}
```

#### 3. 创建页面组件

创建 `frontend/src/views/Deployments.vue`：

```vue
<template>
  <div class="deployments-page">
    <el-table :data="deployments" v-loading="loading">
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="ready" label="就绪" />
      <el-table-column prop="desired" label="期望" />
      <el-table-column prop="age" label="年龄" />
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import * as deploymentApi from '@/api/kube/deployment'

const route = useRoute()
const deployments = ref([])
const loading = ref(false)

async function fetchDeployments() {
  loading.value = true
  try {
    const namespace = route.params.namespace as string
    deployments.value = await deploymentApi.getDeployments(namespace)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchDeployments()
})
</script>
```

#### 4. 添加路由

在 `frontend/src/router/index.ts` 中添加：

```typescript
{
  path: '/namespaces/:namespace/deployments',
  name: 'Deployments',
  component: () => import('@/views/Deployments.vue')
}
```

### 代码规范

- **使用 TypeScript 严格模式**：所有变量必须有类型
- **使用 Composition API**：优先使用 `<script setup>` 语法
- **组件命名**：使用 PascalCase，如 `UserList.vue`
- **使用中文注释**：复杂的逻辑需要中文注释
- **响应式数据**：
  - 基本类型：使用 `ref`
  - 对象类型：可以使用 `reactive`
- **样式**：使用 Scoped CSS，避免样式污染

---

## 测试

### 后端测试

```bash
# 运行所有测试
cd backend
go test ./...

# 运行特定包的测试
go test ./internal/kube/handler

# 查看测试覆盖率
go test -cover ./...
```

### 前端测试

```bash
cd frontend

# 运行单元测试
npm run test

# 查看覆盖率
npm run test:coverage

# 运行 E2E 测试
npm run test:e2e
```

---

## 调试技巧

### 后端调试

**使用 Delve（Go 调试器）**

```bash
# 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 调试运行
dlv debug cmd/server/main.go
```

**使用 VSCode**

创建 `.vscode/launch.json`：

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Package",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/backend/cmd/server",
      "env": {
        "KUBECONFIG": "${env:HOME}/.kube/config"
      }
    }
  ]
}
```

### 前端调试

1. **使用 Vue DevTools**：安装浏览器扩展
2. **网络请求**：在浏览器控制台的 Network 标签查看
3. **日志**：使用 `console.log` 或 `console.error`

---

## 常见问题

### 1. 无法连接 Kubernetes 集群

**问题**：`failed to build kubeconfig`

**解决方案**：

```bash
# 设置 KUBECONFIG 环境变量
export KUBECONFIG=~/.kube/config

# Windows
set KUBECONFIG=%USERPROFILE%\.kube/config
```

### 2. 前端构建失败

**问题**：`npm install` 失败

**解决方案**：

```bash
# 清除缓存
rm -rf node_modules package-lock.json

# 重新安装
npm install
```

### 3. Go 模块依赖问题

**问题**：依赖版本冲突

**解决方案**：

```bash
# 清理依赖
go mod tidy

# 更新依赖
go get -u ./...
```

### 4. 端口已被占用

**问题**：`bind: address already in use`

**解决方案**：

```bash
# 查找占用端口的进程
lsof -i :8080

# 终止进程
kill -9 <PID>
```

---

## 有用的命令

### Go 相关

```bash
# 格式化代码
go fmt ./...

# 代码检查
golangci-lint run

# 构建
go build -o bin/server ./cmd/server

# 运行
go run cmd/server/main.go
```

### 前端相关

```bash
# 开发模式
npm run dev

# 生产构建
npm run build

# 预览构建
npm run preview

# 类型检查
npm run type-check

# 代码检查
npm run lint
```

### Docker 相关

```bash
# 构建镜像
docker build -t kubeops/kubeops:latest .

# 运行容器
docker run -p 8080:8080 kubeops/kubeops:latest
```

---

## 提交代码

使用约定式提交格式：

```
<type>: <subject>

<body>

<footer>
```

**类型（type）**

- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建或辅助工具变动

**示例**

```
feat: 添加 Deployment 列表 API

- 添加 DeploymentRepository
- 添加 DeploymentService
- 添加 DeploymentHandler
- 注册路由

Closes #123
```

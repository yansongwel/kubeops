# 贡献指南

感谢您对 KubeOps 项目的关注！本文档提供贡献的指南和说明。

## 行为准则

- 互相尊重和包容
- 欢迎新手并帮助他们学习
- 关注建设性反馈
- 开放合作和沟通

## 如何贡献

### 报告 Bug

在创建 Bug 报告之前，请先检查现有 issue 以避免重复。

创建 Bug 报告时，请包含：
- **标题**：清晰且描述性强
- **描述**：问题的详细说明
- **重现步骤**：分步说明
- **预期行为**：您期望发生什么
- **实际行为**：实际发生了什么
- **环境信息**：操作系统、K8s 版本、应用版本
- **日志**：相关的日志片段

### 功能建议

欢迎提出功能建议！请：
- 使用清晰且描述性的标题
- 提供功能的详细描述
- 解释为什么这个功能有用
- 列出可能的实现方案
- 考虑是否可以将其实现为插件

### Pull Request

#### 创建 PR 之前

1. **Fork 仓库**并从 `develop` 分支创建您的分支
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **代码风格**：遵循项目的代码风格
   - Go：运行 `gofmt` 和 `golangci-lint`
   - Vue：运行 `eslint` 和 `prettier`

3. **提交信息**：使用[约定式提交](https://www.conventionalcommits.org/zh-hans/)
   ```
   feat: 添加用户认证
   fix: 修复指标收集器中的内存泄漏
   docs: 更新 API 文档
   test: 为 RBAC 模块添加单元测试
   refactor: 简化资源管理器代码
   ```

4. **测试**：为新功能编写测试
   - 为业务逻辑编写单元测试
   - 为 API 端点编写集成测试
   - 更新文档

#### 创建 PR

1. **标题**：遵循约定式提交格式
   - `feat: 添加对自定义 CRD 的支持`
   - `fix: 修复资源监听器中的竞态条件`

2. **描述**：包含：
   - **摘要**：此 PR 的作用
   - **动机**：为什么需要此更改
   - **更改内容**：主要更改列表
   - **测试方式**：您如何测试此更改
   - **截图**：UI 更改的截图（如适用）
   - **破坏性更改**：任何破坏性更改和迁移指南
   - **相关问题**：相关 issue 的链接

3. **检查清单**：
   - [ ] 代码遵循项目风格指南
   - [ ] 已添加/更新测试并通过
   - [ ] 已更新文档
   - [ ] 提交信息遵循约定式提交
   - [ ] 没有合并冲突
   - [ ] 已填写 PR 描述

#### PR 审核流程

- **自动检查**：CI/CD 流水线自动运行
- **代码审核**：维护者审核您的代码
- **反馈**：处理审核意见
- **批准**：至少需要一名维护者批准
- **合并**：Squash and merge 到 `develop` 分支

## 开发环境搭建

### 前置要求

- Go 1.21+
- Node.js 18+
- Docker/Podman
- kubectl
- helm
- kind（用于本地开发）

### 初始化步骤

1. **克隆您的 Fork**
   ```bash
   git clone https://github.com/YOUR_USERNAME/kubeops.git
   cd kubeops
   ```

2. **安装依赖**
   ```bash
   # 后端
   cd backend
   go mod download

   # 前端
   cd ../frontend
   npm install
   ```

3. **启动开发环境**
   ```bash
   # 启动依赖服务
   docker-compose -f deploy/docker-compose-dev.yaml up -d

   # 运行后端服务（每个终端）
   cd backend/api-gateway && go run cmd/server/main.go
   cd backend/kube-manager && go run cmd/server/main.go

   # 运行前端
   cd frontend && npm run dev
   ```

### 运行测试

```bash
# 后端
cd backend
make test

# 前端
cd frontend
npm run test
```

### 代码风格

#### 后端 (Go)

```bash
# 格式化代码
gofmt -w .

# 代码检查
golangci-lint run

# 自动修复
golangci-lint run --fix
```

#### 前端 (Vue)

```bash
# 代码检查
npm run lint

# 格式化
npm run format
```

## 项目结构

```
KubeOps/
├── backend/
│   ├── api-gateway/          # API 网关服务
│   ├── kube-manager/         # K8s 资源管理器
│   ├── ai-inspector/         # AI 巡检服务
│   ├── devops-service/       # CI/CD 集成
│   ├── logging-service/      # 日志平台集成
│   ├── monitoring-service/   # 监控集成
│   ├── common/               # 共享库
│   └── proto/                # gRPC 定义
├── frontend/                 # Vue 3 仪表盘
│   ├── src/
│   │   ├── components/       # 可复用组件
│   │   ├── views/            # 页面组件
│   │   ├── stores/           # Pinia 状态管理
│   │   ├── api/              # API 客户端
│   │   └── router/           # Vue Router 配置
├── deploy/                   # 部署配置
│   ├── helm/                 # Helm Charts
│   └── examples/             # 示例配置
└── docs/                     # 文档
```

## 编码指南

### Go 后端

1. **包命名**：使用小写单词
2. **错误处理**：始终检查并处理错误
3. **日志**：使用结构化日志 (zap)
4. **配置**：使用环境变量
5. **上下文**：在 API 调用中始终传递上下文

示例：
```go
func (s *Service) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
    pod, err := s.client.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
    if err != nil {
        s.logger.Error("获取 Pod 失败",
            zap.String("命名空间", namespace),
            zap.String("名称", name),
            zap.Error(err))
        return nil, fmt.Errorf("获取 Pod: %w", err)
    }
    return pod, nil
}
```

### Vue 前端

1. **组合式 API**：使用 `<script setup>` 语法
2. **TypeScript**：使用严格 TypeScript
3. **组件**：保持组件小而专注
4. **状态管理**：使用 Pinia 管理全局状态
5. **API 调用**：使用专用的 API 模块

示例：
```vue
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useClusterStore } from '@/stores/cluster'

const clusterStore = useClusterStore()
const clusters = ref<Cluster[]>([])

onMounted(async () => {
  clusters.value = await clusterStore.listClusters()
})
</script>
```

## 添加新功能

### 1. 后端服务

```bash
# 创建新服务目录
mkdir backend/new-service
cd backend/new-service

# 创建服务结构
mkdir -p cmd/server pkg/api pkg/service

# 实现服务
# 添加到 docker-compose 和 helm charts
```

### 2. 前端页面

```bash
# 创建视图
touch frontend/src/views/NewFeature.vue

# 添加路由
# 编辑 frontend/src/router/index.ts

# 添加导航项
# 编辑 frontend/src/layouts/MainLayout.vue
```

### 3. API 端点

```go
// backend/service/pkg/api/handler.go
func (h *Handler) NewEndpoint(c *gin.Context) {
    var req Request
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    result, err := h.service.NewMethod(c.Request.Context(), req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, result)
}
```

## 文档

### 代码注释

- 导出的函数必须有注释
- 复杂逻辑应加以说明
- Go 使用 godoc 格式

```go
// GetCluster 根据 ID 获取集群。
// 如果集群未找到，返回错误。
func (s *Service) GetCluster(ctx context.Context, id string) (*Cluster, error) {
    // ...
}
```

### API 文档

- 为新端点更新 OpenAPI 规范
- 添加请求/响应示例
- 记录错误代码

### README 更新

- 为新功能更新 README
- 添加/更新架构图
- 更新开发指南

## 测试指南

### 单元测试

- 测试公共 API
- 模拟外部依赖
- 测试错误情况
- 目标覆盖率 >80%

示例：
```go
func TestGetPod(t *testing.T) {
    tests := []struct {
        name      string
        namespace string
        podName   string
        wantErr   bool
    }{
        {
            name:      "有效的 Pod",
            namespace: "default",
            podName:   "test-pod",
            wantErr:   false,
        },
        // ... 更多测试用例
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 测试实现
        })
    }
}
```

### 集成测试

- 测试 API 端点
- 使用 testcontainers 管理依赖
- 测试后清理

### 端到端测试

- 测试用户工作流
- 使用 Selenium/Puppeteer 进行 UI 测试

## 发布流程

1. **版本**：遵循语义化版本 (MAJOR.MINOR.PATCH)
2. **更新日志**：更新 CHANGELOG.md
3. **标签**：创建 git 标签
4. **构建**：构建发布产物
5. **部署**：部署到测试环境
6. **测试**：全面测试
7. **公告**：发布说明和公告

## 获取帮助

- **文档**：查看 `docs/` 目录
- **Issues**：搜索或创建 GitHub issue
- **Discussions**：使用 GitHub Discussions
- **聊天**：加入我们的 Discord/Slack

## 致谢

贡献者将获得：
- 在 CONTRIBUTORS.md 中列出
- 在发布说明中提及
- 邀请成为维护者（对于重要贡献）

## 开源协议

通过贡献，您同意您的贡献将根据 Apache License 2.0 许可。

---

感谢您对 KubeOps 的贡献！🎉

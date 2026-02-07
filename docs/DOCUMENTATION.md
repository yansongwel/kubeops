# 项目文档说明

## 文档语言规范

本项目采用 **中文** 作为主要文档和注释语言，具体规范如下：

### 1. Markdown 文档
所有 `.md` 文档均使用中文编写，包括：
- `README.md` - 项目说明
- `CLAUDE.md` - Claude Code 使用指南
- `CONTRIBUTING.md` - 贡献指南
- `docs/architecture/README.md` - 架构文档
- `docs/QUICKSTART.md` - 快速开始指南
- `LICENSE` - 开源协议（保持英文 Apache 2.0）

### 2. 代码注释

#### Go 后端代码
```go
// 初始化日志记录器
logger, _ := zap.NewProduction()

// 初始化 Kubernetes 客户端
k8sClient, err := kubernetes.NewForConfig(k8sConfig)

// 健康检查端点
router.GET("/health", func(c *gin.Context) {
    // ...
})
```

#### Vue 前端代码
```vue
<script setup lang="ts">
import { ref, onMounted } from 'vue'
// 导入集群状态管理
import { useClusterStore } from '@/stores/cluster'

// 初始化集群列表
const clusters = ref<Cluster[]>([])
</script>
```

### 3. 用户界面
前端用户界面采用中文，包括：
- 菜单项：仪表盘、集群管理、工作负载、AI 巡检等
- 页面标题和描述
- 表单标签和提示信息
- 错误消息和成功提示

### 4. 配置文件注释
```yaml
# 日志配置
logging:
  stack: elk  # 可选: elk, loki

# 监控配置
monitoring:
  stack: prometheus  # 可选: prometheus, victoriametrics
```

### 5. API 文档
虽然代码中使用中文注释，但：
- API 端点名称保持英文（RESTful 规范）
- HTTP 状态码使用标准英文
- 错误消息可使用中文便于理解

### 6. 代码规范

#### Go 代码
- 包名使用英文小写
- 变量名、函数名使用英文（Go 命名规范）
- 注释使用中文说明功能和逻辑

#### Vue/TypeScript 代码
- 组件名使用 PascalCase 英文
- 变量名、函数名使用英文驼峰命名
- 注释使用中文

## 示例

### API 端点示例
```go
// GetPod 获取指定命名空间的 Pod 详情
// 如果 Pod 不存在，返回 404 错误
func (h *Handler) GetPod(c *gin.Context) {
    // 获取命名空间和 Pod 名称
    namespace := c.Param("namespace")
    name := c.Param("name")

    // 调用 Kubernetes API
    pod, err := h.client.CoreV1().Pods(namespace).Get(c, name, metav1.GetOptions{})
    // ...
}
```

### Vue 组件示例
```vue
<template>
  <div class="cluster-list">
    <!-- 集群列表容器 -->
    <el-table :data="clusters">
      <el-table-column prop="name" label="集群名称" />
      <el-table-column prop="status" label="状态" />
    </el-table>
  </div>
</template>

<script setup lang="ts">
// 集群列表组件
import { ref, onMounted } from 'vue'

// 定义集群数据类型
interface Cluster {
  name: string
  status: string
}

// 初始化集群列表
const clusters = ref<Cluster[]>([])

// 加载集群数据
onMounted(async () => {
  clusters.value = await loadClusters()
})
</script>
```

## 文档清单

### 已完成中文化的文档

| 文档路径 | 说明 | 状态 |
|---------|------|------|
| `README.md` | 项目主文档 | ✅ 已完成 |
| `CLAUDE.md` | Claude Code 使用指南 | ✅ 已完成 |
| `CONTRIBUTING.md` | 贡献指南 | ✅ 已完成 |
| `docs/architecture/README.md` | 架构文档 | ✅ 已完成 |
| `docs/QUICKSTART.md` | 快速开始 | ✅ 已完成 |
| `docs/deployment/README.md` | 部署指南 | ✅ 已完成 |
| `backend/api-gateway/cmd/server/main.go` | API 网关代码 | ✅ 已完成 |
| `backend/kube-manager/cmd/server/main.go` | Kube 管理器代码 | ✅ 已完成 |
| `frontend/src/main.ts` | 前端入口 | ✅ 已完成 |
| `frontend/src/router/index.ts` | 路由配置 | ✅ 已完成 |
| `frontend/src/layouts/MainLayout.vue` | 主布局 | ✅ 已完成 |
| `frontend/src/views/*.vue` | 页面组件 | ✅ 已完成 |
| `deploy/helm/kubeops/values.yaml` | Helm 配置 | ✅ 已完成 |

### 待创建的文档

| 文档路径 | 说明 | 优先级 |
|---------|------|--------|
| `docs/api/README.md` | API 文档 | 高 |
| `docs/development/README.md` | 开发指南 | 高 |
| `docs/configuration/ai-provider.md` | AI 配置指南 | 中 |
| `docs/configuration/cicd.md` | CI/CD 配置指南 | 中 |
| `docs/deployment/logging.md` | 日志部署指南 | 中 |
| `docs/deployment/monitoring.md` | 监控部署指南 | 中 |
| `CHANGELOG.md` | 变更日志 | 低 |
| `docs/troubleshooting.md` | 故障排查 | 低 |

## 维护规范

### 添加新功能时
1. 代码注释使用中文
2. API 文档使用中文描述，英文示例
3. 用户界面文本使用中文
4. 更新相关 Markdown 文档

### 修复 Bug 时
1. 提交信息使用中文描述
2. 代码注释说明修复原因（中文）
3. 更新 CHANGELOG.md（如有）

### 文档更新时
1. 保持中文为主要语言
2. 技术术语可保留英文（如 Kubernetes、Prometheus）
3. 代码示例保持英文，注释使用中文

## 注意事项

1. **代码标识符**：变量名、函数名、类名等使用英文
2. **技术术语**：如 Kubernetes、Docker、Prometheus 等专有名词保持英文
3. **API 设计**：RESTful API 端点、字段名使用英文
4. **日志输出**：开发环境可用中文，生产环境建议英文便于国际社区协作

## 贡献者指南

所有贡献者应遵循本文档规范：
- 新增代码的注释使用中文
- 更新文档时使用中文
- 用户界面文本使用中文
- 保持代码标识符使用英文

---

**文档版本**：v0.1.0
**最后更新**：2026-02-07

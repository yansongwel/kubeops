# KubeOps 前端开发学习指南

> 本指南面向希望提升全栈开发能力的运维工程师

## 学习目标

通过 KubeOps 项目，您将掌握：

### 前端技能
- ✅ Vue 3 Composition API（现代化开发范式）
- ✅ TypeScript（类型安全）
- ✅ 前端工程化（Vite、构建优化）
- ✅ 状态管理（Pinia）
- ✅ 组件化开发思维
- ✅ 前后端联调经验

### 后端技能
- ✅ Go 语言（云原生开发）
- ✅ 微服务架构设计
- ✅ RESTful API 设计
- ✅ gRPC 通信
- ✅ 数据库设计
- ✅ Kubernetes 编程

### 架构能力
- ✅ 前后端分离架构
- ✅ 微服务拆分原则
- ✅ API 设计规范
- ✅ 数据流向设计
- ✅ 安全架构
- ✅ 性能优化

## 学习路径

### 阶段一：前端基础（2-3 周）

#### Week 1: Vue 3 基础
```bash
# 学习资源
- Vue 3 官方文档: https://cn.vuejs.org/
- TypeScript: https://www.typescriptlang.org/zh/docs/
- Vite: https://cn.vitejs.dev/

# 实践任务
1. 完成登录页面
2. 实现主布局
3. 创建路由配置
```

#### Week 2: 组件开发
```bash
# 学习重点
- 组件通信（props、emit、provide/inject）
- 生命周期
- 组合式函数
- TypeScript 类型定义

# 实践任务
1. 开发表格组件（Pod 列表）
2. 开发表单组件（创建 Deployment）
3. 开发详情组件（YAML 编辑器）
```

#### Week 3: 状态管理
```bash
# 学习重点
- Pinia Store
- API 封装
- 错误处理
- Loading 状态

# 实践任务
1. 集成 Pinia
2. 封装 API 请求
3. 实现全局状态管理
```

### 阶段二：后端开发（3-4 周）

#### Week 1-2: Go 基础
```bash
# 学习资源
- Go 语言圣经: https://gopl-zh.github.io/
- Go by Example: https://gobyexample.com/
- Kubernetes client-go: https://github.com/kubernetes/client-go

# 实践任务
1. 实现 API 网关基础框架
2. 实现 Kube Manager 核心 API
3. 实现中间件（认证、日志、限流）
```

#### Week 3-4: 微服务架构
```bash
# 学习重点
- 微服务拆分原则
- 服务间通信（gRPC）
- 数据库设计
- 缓存策略
- 错误处理和重试

# 实践任务
1. 实现 AI 巡检服务
2. 实现 DevOps 服务
3. 实现日志和监控服务
```

### 阶段三：前后端联调（2 周）

#### Week 1: API 对接
```bash
# 实践任务
1. 定义 API 接口规范
2. 前端封装 API
3. 后端实现接口
4. 联调测试
```

#### Week 2: 功能完善
```bash
# 实践任务
1. 集群管理功能
2. 工作负载管理
3. 实时日志查看
4. 权限管理
```

### 阶段四：高级特性（4 周）

#### Week 1-2: 实时通信
```bash
# 学习重点
- WebSocket
- Server-Sent Events
- 实时日志流
- 实时指标更新
```

#### Week 3-4: 性能优化
```bash
# 学习重点
- 前端性能优化（懒加载、虚拟滚动）
- 后端性能优化（缓存、连接池）
- 数据库优化
- 监控和调试
```

## 核心模块实现顺序

### 1. 用户认证模块 ⭐⭐⭐

**学习重点**：
- JWT Token 机制
- 前端路由守卫
- 后端中间件
- 会话管理

**实现步骤**：
```bash
# 前端
1. 登录页面 UI
2. Token 存储和刷新
3. 路由守卫
4. Axios 请求拦截器

# 后端
1. JWT 生成和验证
2. 用户认证接口
3. 权限中间件
4. 会话管理
```

**代码示例**：
```typescript
// 前端：src/api/auth.ts
export async function login(username: string, password: string) {
  const { data } = await axios.post('/api/v1/auth/login', {
    username,
    password,
  });
  localStorage.setItem('token', data.token);
  return data;
}

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token');
  if (!token && to.path !== '/login') {
    next('/login');
  } else {
    next();
  }
});
```

```go
// 后端：jwt.go
func GenerateToken(userID string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })
    return token.SignedString([]byte("secret"))
}
```

### 2. 集群管理模块 ⭐⭐⭐⭐⭐

**学习重点**：
- CRUD 操作
- 表格组件
- 详情页
- YAML 编辑器

**实现步骤**：
```bash
1. 集群列表页（表格）
2. 创建集群对话框
3. 集群详情页
4. 资源 YAML 编辑器
```

**核心代码结构**：
```typescript
// src/views/kube/clusters/index.vue
<template>
  <div class="clusters-page">
    <!-- 工具栏 -->
    <el-button type="primary" @click="handleCreate">
      创建集群
    </el-button>

    <!-- 集群列表 -->
    <el-table :data="clusters" v-loading="loading">
      <el-table-column prop="name" label="集群名称" />
      <el-table-column prop="version" label="版本" />
      <el-table-column prop="nodeCount" label="节点数" />
      <el-table-column label="操作">
        <template #default="{ row }">
          <el-button @click="handleView(row)">查看</el-button>
          <el-button @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { listClusters, deleteCluster } from '@/api/kube/clusters';

const clusters = ref([]);
const loading = ref(false);

const loadClusters = async () => {
  loading.value = true;
  clusters.value = await listClusters();
  loading.value = false;
};

onMounted(() => {
  loadClusters();
});
</script>
```

### 3. 工作负载模块 ⭐⭐⭐⭐⭐

**学习重点**：
- 复杂表格（展开行、嵌套表格）
- Tab 切换
- 实时状态更新
- 操作确认

**实现步骤**：
```bash
1. Pod 列表（带展开行显示容器信息）
2. Deployment 列表（副本数、更新策略）
3. Service 列表（端点、端口映射）
4. StatefulSet/DaemonSet
```

### 4. YAML 编辑器 ⭐⭐⭐⭐

**学习重点**：
- 代码编辑器集成（Monaco Editor）
- 语法高亮
- 格式化
- 验证

**实现方案**：
```typescript
// 安装 Monaco Editor
npm install monaco-editor @monaco-editor/loader

// src/components/YamlEditor.vue
<template>
  <div class="yaml-editor" ref="editorContainer"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import * as monaco from 'monaco-editor';

const editorContainer = ref<HTMLElement>();
let editor: monaco.editor.IStandaloneCodeEditor;

onMounted(() => {
  editor = monaco.editor.create(editorContainer.value!, {
    value: '# YAML 内容',
    language: 'yaml',
    theme: 'vs-dark',
    minimap: { enabled: false },
  });
});

onUnmounted(() => {
  editor?.dispose();
});
</script>
```

### 5. 实时日志模块 ⭐⭐⭐⭐⭐

**学习重点**：
- WebSocket 连接
- 日志流式显示
- 自动滚动
- 关键字搜索

**实现方案**：
```typescript
// src/composables/useWebSocket.ts
export function useWebSocket(url: string) {
  const ws = ref<WebSocket | null>(null);
  const data = ref<any[]>([]);
  const connected = ref(false);

  const connect = () => {
    ws.value = new WebSocket(url);
    ws.value.onopen = () => {
      connected.value = true;
    };
    ws.value.onmessage = (event) => {
      data.value.push(JSON.parse(event.data));
    };
    ws.value.onclose = () => {
      connected.value = false;
    };
  };

  const disconnect = () => {
    ws.value?.close();
  };

  return { connect, disconnect, data, connected };
}

// 使用
const { connect, data, connected } = useWebSocket('ws://localhost:8080/logs');
connect();
```

## 编程思维提升

### 1. 分层思维

```
前端分层：
├── 视图层
├── 逻辑层
│   ├── API 层
│   ├── Store 层
│   └── 业务层
└── 组件层

后端分层：
├── 接口层
├── 业务层
│   ├── Service
│   └── Domain
├── 数据层
│   ├── Repository
│   └── Model
└── 基础设施层
```

### 2. 模块化思维

**如何拆分模块**：
- 按业务功能拆分（集群、工作负载、AI）
- 按技术层次拆分（API、Store、组件）
- 单一职责原则
- 依赖倒置

### 3. 数据流思维

```
用户操作 → 组件事件 → Store Action → API 请求
  ↓
后端处理 → 数据库/K8s API → 返回结果
  ↓
Store Mutation → 组件响应式更新 → UI 渲染
```

### 4. 错误处理思维

```typescript
// 统一错误处理
class ApiError extends Error {
  constructor(
    public code: number,
    message: string,
    public details?: any
  ) {
    super(message);
  }
}

// 使用
try {
  await createCluster(data);
} catch (error) {
  if (error instanceof ApiError) {
    handleError(error);
  }
}
```

## 实践建议

### 1. 小步快跑
- 不要一次实现所有功能
- 每个功能点独立开发和测试
- 及时提交代码

### 2. 写测试
- 单元测试
- 组件测试
- API 测试

### 3. 代码审查
- 自己 review 代码
- 学习优秀开源项目
- 关注代码质量

### 4. 记录学习
- 写技术博客
- 记录踩坑经验
- 总结最佳实践

## 学习资源推荐

### Vue 3
- [Vue 3 官方文档](https://cn.vuejs.org/)
- [Vue Mastery](https://www.vuemastery.com/)
- [Composition API RFC](https://github.com/vuejs/rfcs)

### Go
- [Go 语言圣经](https://gopl-zh.github.io/)
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://go.dev/doc/effective_go)

### Kubernetes
- [Kubernetes 官方文档](https://kubernetes.io/zh-cn/docs/)
- [client-go 示例](https://github.com/kubernetes/client-go/tree/master/examples)
- [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime)

### 架构设计
- [《架构整洁之道》](https://book.douban.com/subject/27022909/)
- [《微服务架构设计模式》](https://book.douban.com/subject/30356827/)
- [《凤凰项目》](https://book.douban.com/subject/25779198/)

## 项目检查点

### 第 1 个月
- [ ] 前端基础框架搭建
- [ ] 登录认证功能
- [ ] 集群列表页
- [ ] 后端 API 网关基础

### 第 2 个月
- [ ] 工作负载管理（Pod、Deployment）
- [ ] YAML 编辑器
- [ ] 实时日志查看
- [ ] 后端 Kube Manager 完善

### 第 3 个月
- [ ] AI 巡检基础功能
- [ ] DevOps 集成（Jenkins）
- [ ] 监控图表展示
- [ ] 性能优化

### 第 4 个月
- [ ] 完整功能测试
- [ ] 文档完善
- [ ] 部署优化
- [ ] 开源发布

## 总结

作为高级运维工程师，您已经具备：
- ✅ Kubernetes 专业知识
- ✅ 系统架构理解
- ✅ 运维最佳实践

通过 KubeOps 项目，您将获得：
- 🎯 全栈开发能力
- 🎯 产品思维
- 🎯 工程化能力
- 🎯 架构设计能力

这是一个非常有价值的学习项目！加油！💪

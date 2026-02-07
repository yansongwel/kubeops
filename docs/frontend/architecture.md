# 前端项目架构文档

## 技术栈

```
Vue 3.4+          # 渐进式框架
├── TypeScript    # 类型系统
├── Vite 5.x      # 构建工具
├── Element Plus  # UI 组件库
├── Pinia         # 状态管理
├── Vue Router    # 路由管理
├── Axios         # HTTP 客户端
├── Monaco Editor # YAML 编辑器
├── ECharts       # 图表库
└── Day.js        # 时间处理
```

## 项目结构

```
frontend/
├── public/                 # 静态资源
│   └── favicon.ico
├── src/
│   ├── api/               # API 接口层
│   │   ├── kube/         # Kubernetes 相关 API
│   │   │   ├── cluster.ts
│   │   │   ├── workload.ts
│   │   │   └── config.ts
│   │   ├── ai/           # AI 巡检 API
│   │   ├── devops/       # DevOps API
│   │   └── auth.ts       # 认证 API
│   │
│   ├── assets/           # 资源文件
│   │   ├── images/
│   │   ├── styles/
│   │   └── icons/
│   │
│   ├── components/       # 公共组件
│   │   ├── common/       # 通用组件
│   │   │   ├── Table.vue
│   │   │   ├── Form.vue
│   │   │   └── Dialog.vue
│   │   ├── kube/         # K8s 专用组件
│   │   │   ├── PodList.vue
│   │   │   ├── YamlEditor.vue
│   │   │   └── LogViewer.vue
│   │   └── charts/       # 图表组件
│   │       ├── LineChart.vue
│   │       └── PieChart.vue
│   │
│   ├── composables/      # 组合式函数
│   │   ├── useWebSocket.ts
│   │   ├── useTable.ts
│   │   └── useRequest.ts
│   │
│   ├── layouts/          # 布局组件
│   │   ├── MainLayout.vue
│   │   └── BlankLayout.vue
│   │
│   ├── router/           # 路由配置
│   │   ├── index.ts
│   │   └── routes/       # 路由模块
│   │       ├── kube.ts
│   │       ├── ai.ts
│   │       └── devops.ts
│   │
│   ├── stores/           # Pinia 状态管理
│   │   ├── user.ts       # 用户状态
│   │   ├── cluster.ts    # 集群状态
│   │   └── app.ts        # 应用状态
│   │
│   ├── types/            # TypeScript 类型定义
│   │   ├── kube.ts       # K8s 类型
│   │   ├── api.ts        # API 类型
│   │   └── global.d.ts   # 全局类型
│   │
│   ├── utils/            # 工具函数
│   │   ├── request.ts    # Axios 封装
│   │   ├── storage.ts    # 本地存储
│   │   ├── format.ts     # 格式化
│   │   └── validate.ts   # 验证
│   │
│   ├── views/            # 页面视图
│   │   ├── auth/         # 认证
│   │   │   └── Login.vue
│   │   ├── dashboard/    # 仪表盘
│   │   ├── kube/         # K8s 资源
│   │   │   ├── clusters/
│   │   │   ├── pods/
│   │   │   ├── deployments/
│   │   │   └── services/
│   │   ├── ai/           # AI 巡检
│   │   ├── devops/       # DevOps
│   │   ├── logs/         # 日志
│   │   └── monitoring/   # 监控
│   │
│   ├── App.vue           # 根组件
│   └── main.ts           # 入口文件
│
├── index.html
├── vite.config.ts        # Vite 配置
├── tsconfig.json         # TypeScript 配置
├── .eslintrc.cjs         # ESLint 配置
└── package.json
```

## 核心设计模式

### 1. 分层架构

```
┌─────────────────────────────────────┐
│          Views (视图层)              │
│  - 组件交互                          │
│  - 用户操作                          │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│       Composables (逻辑层)           │
│  - 业务逻辑                          │
│  - 状态管理                          │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│         API (接口层)                 │
│  - HTTP 请求                         │
│  - 错误处理                          │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│        Backend (后端服务)            │
└─────────────────────────────────────┘
```

### 2. 组件通信

```typescript
// Props Down, Events Up
<template>
  <ChildComponent
    :data="parentData"
    @update="handleUpdate"
  />
</template>

// Provide/Inject (跨层级)
// 父组件
provide('clusterInfo', clusterInfo)

// 子组件
const clusterInfo = inject('clusterInfo')

// Store (全局状态)
const clusterStore = useClusterStore()
clusterStore.setCluster(cluster)
```

### 3. 代码组织

#### API 层示例

```typescript
// src/api/kube/pods.ts
import request from '@/utils/request';
import type { Pod, PodList, PodQuery } from '@/types/kube';

enum Api {
  List = '/kube/pods',
  Get = '/kube/pods/:namespace/:name',
  Create = '/kube/pods',
  Delete = '/kube/pods/:namespace/:name',
  Logs = '/kube/pods/:namespace/:name/logs',
}

export function listPods(params: PodQuery) {
  return request.get<PodList>(Api.List, { params });
}

export function getPod(namespace: string, name: string) {
  return request.get<Pod>(Api.Get, {
    params: { namespace, name },
  });
}

export function createPod(data: Pod) {
  return request.post(Api.Create, data);
}

export function deletePod(namespace: string, name: string) {
  return request.delete(Api.Delete, {
    params: { namespace, name },
  });
}

export function getPodLogs(namespace: string, name: string) {
  return new EventSource(
    `${import.meta.env.VITE_API_BASE_URL}${Api.Logs}?namespace=${namespace}&name=${name}`
  );
}
```

#### Composable 示例

```typescript
// src/composables/useTable.ts
import { ref, computed } from 'vue';
import { ElMessage } from 'element-plus';

export function useTable<T>(apiFn: () => Promise<T[]>) {
  const data = ref<T[]>([]);
  const loading = ref(false);
  const error = ref<Error | null>(null);

  const fetch = async () => {
    loading.value = true;
    try {
      data.value = await apiFn();
    } catch (e) {
      error.value = e as Error;
      ElMessage.error('加载数据失败');
    } finally {
      loading.value = false;
    }
  };

  const isEmpty = computed(() => data.value.length === 0);

  return {
    data,
    loading,
    error,
    fetch,
    isEmpty,
  };
}

// 使用
const { data: pods, loading, fetch } = useTable(listPods);
```

#### Store 示例

```typescript
// src/stores/cluster.ts
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { Cluster } from '@/types/kube';

export const useClusterStore = defineStore('cluster', () => {
  const clusters = ref<Cluster[]>([]);
  const currentCluster = ref<Cluster | null>(null);

  const clusterCount = computed(() => clusters.value.length);

  function setClusters(data: Cluster[]) {
    clusters.value = data;
  }

  function setCurrentCluster(cluster: Cluster) {
    currentCluster.value = cluster;
    localStorage.setItem('currentCluster', JSON.stringify(cluster));
  }

  function getCurrentCluster() {
    if (!currentCluster.value) {
      const saved = localStorage.getItem('currentCluster');
      if (saved) {
        currentCluster.value = JSON.parse(saved);
      }
    }
    return currentCluster.value;
  }

  return {
    clusters,
    currentCluster,
    clusterCount,
    setClusters,
    setCurrentCluster,
    getCurrentCluster,
  };
});
```

## 页面开发模板

### 列表页模板

```vue
<template>
  <div class="page-container">
    <!-- 工具栏 -->
    <div class="toolbar">
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        创建
      </el-button>
      <el-button @click="handleRefresh">
        <el-icon><Refresh /></el-icon>
        刷新
      </el-button>
    </div>

    <!-- 筛选器 -->
    <div class="filter">
      <el-form :model="query" inline>
        <el-form-item label="命名空间">
          <el-select v-model="query.namespace" placeholder="请选择">
            <el-option label="全部" value="" />
            <el-option
              v-for="ns in namespaces"
              :key="ns.name"
              :label="ns.name"
              :value="ns.name"
            />
          </el-select>
        </el-form-item>
      </el-form>
    </div>

    <!-- 表格 -->
    <el-table
      v-loading="loading"
      :data="tableData"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="namespace" label="命名空间" />
      <el-table-column label="状态">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">
            {{ row.status }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" fixed="right" width="200">
        <template #default="{ row }">
          <el-button link @click="handleView(row)">查看</el-button>
          <el-button link @click="handleEdit(row)">编辑</el-button>
          <el-button link type="danger" @click="handleDelete(row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <el-pagination
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.pageSize"
      :total="pagination.total"
      @current-change="handlePageChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';

// 查询参数
const query = ref({
  namespace: '',
});

// 表格数据
const tableData = ref([]);
const loading = ref(false);

// 分页
const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0,
});

// 加载数据
const loadData = async () => {
  loading.value = true;
  try {
    const { data, total } = await listResources({
      ...query.value,
      ...pagination.value,
    });
    tableData.value = data;
    pagination.value.total = total;
  } catch (error) {
    ElMessage.error('加载数据失败');
  } finally {
    loading.value = false;
  }
};

// 操作
const handleCreate = () => {
  // 打开创建对话框
};

const handleView = (row: any) => {
  // 跳转到详情页
};

const handleEdit = (row: any) => {
  // 打开编辑对话框
};

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定删除吗？', '提示');
    await deleteResource(row.id);
    ElMessage.success('删除成功');
    loadData();
  } catch (error) {
    // 用户取消
  }
};

const handleRefresh = () => {
  loadData();
};

const handlePageChange = (page: number) => {
  pagination.value.page = page;
  loadData();
};

// 生命周期
onMounted(() => {
  loadData();
});
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.toolbar {
  margin-bottom: 20px;
}

.filter {
  margin-bottom: 20px;
}
</style>
```

## 最佳实践

### 1. TypeScript 类型定义

```typescript
// src/types/kube/pod.ts
export interface Pod {
  name: string;
  namespace: string;
  status: string;
  createTime: string;
  labels: Record<string, string>;
  containers: Container[];
}

export interface Container {
  name: string;
  image: string;
  resources: Resources;
}

export interface Resources {
  requests: {
    cpu: string;
    memory: string;
  };
  limits: {
    cpu: string;
    memory: string;
  };
}

export type PodStatus = 'Running' | 'Pending' | 'Failed' | 'Succeeded';
```

### 2. API 请求封装

```typescript
// src/utils/request.ts
import axios, { AxiosError, AxiosRequestConfig } from 'axios';
import { ElMessage } from 'element-plus';

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 30000,
});

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    return response.data;
  },
  (error: AxiosError) => {
    const message = error.response?.data?.message || '请求失败';
    ElMessage.error(message);
    return Promise.reject(error);
  }
);

export default request;
```

### 3. 路由懒加载

```typescript
// src/router/index.ts
import { createRouter, createWebHistory } from 'vue-router';

const routes = [
  {
    path: '/pods',
    component: () => import('@/views/kube/pods/index.vue'),
    meta: { title: 'Pod 管理' },
  },
];
```

### 4. 环境变量

```bash
# .env.development
VITE_API_BASE_URL=http://localhost:8080

# .env.production
VITE_API_BASE_URL=https://api.kubeops.com
```

## 性能优化

### 1. 虚拟滚动（长列表）

```typescript
import { useVirtualList } from '@vueuse/core';

const { list, containerProps, wrapperProps } = useVirtualList(
  largeDataSource,
  {
    itemHeight: 50,
  }
);
```

### 2. 防抖和节流

```typescript
import { useDebounceFn } from '@vueuse/core';

const search = useDebounceFn(async (keyword: string) => {
  // 搜索逻辑
}, 500);
```

### 3. 代码分割

```typescript
// 路由级别分割
const Dashboard = () => import('@/views/Dashboard.vue');

// 组件级别分割
const HeavyComponent = defineAsyncComponent(() =>
  import('./HeavyComponent.vue')
);
```

## 调试技巧

### 1. Vue DevTools

```bash
# 安装
npm add -D @vue/devtools

# 使用
main.ts
import { createApp } from 'vue'
import VueDevTools from '@vue/devtools'

const app = createApp(App)
app.use(VueDevTools)
```

### 2. 网络请求调试

```typescript
// 添加请求日志
if (import.meta.env.DEV) {
  request.interceptors.request.use((config) => {
    console.log('Request:', config);
    return config;
  });
}
```

### 3. 性能分析

```bash
# Vite 可视化构建分析
npm add -D rollup-plugin-visualizer

# vite.config.ts
import { visualizer } from 'rollup-plugin-visualizer';

export default defineConfig({
  plugins: [
    visualizer({ open: true }),
  ],
});
```

这个架构设计既适合学习，又能支持项目长期演进！

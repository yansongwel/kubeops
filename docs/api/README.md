# KubeOps API 文档

## 概述

KubeOps 采用 RESTful API 设计，所有 API 请求和响应均使用 JSON 格式。

## 基础信息

- **Base URL**: `http://localhost:8080/api/v1`
- **认证方式**: JWT Bearer Token
- **响应格式**: JSON
- **字符编码**: UTF-8

## 统一响应格式

### 成功响应

```json
{
  "code": 200,
  "message": "success",
  "data": { ... }
}
```

### 错误响应

```json
{
  "code": 400,
  "message": "请求参数错误",
  "details": "详细错误信息"
}
```

## 认证流程

### 1. 登录获取 Token

**请求**

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin"
}
```

**响应**

```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresAt": "2026-02-08T16:00:00Z"
  }
}
```

### 2. 使用 Token 访问 API

在请求头中添加 Authorization：

```http
Authorization: Bearer {token}
```

---

## 集群管理 API

### 获取集群列表

```http
GET /api/v1/clusters
Authorization: Bearer {token}
```

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": "cluster-1",
      "name": "production-cluster",
      "endpoint": "https://k8s.example.com",
      "status": "Connected",
      "version": "v1.33.0",
      "nodeCount": 5,
      "createdAt": "2026-01-01T00:00:00Z"
    }
  ]
}
```

### 获取集群详情

```http
GET /api/v1/clusters/:id
Authorization: Bearer {token}
```

### 创建集群

```http
POST /api/v1/clusters
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "my-cluster",
  "endpoint": "https://k8s.example.com",
  "kubeconfig": "..."
}
```

### 更新集群

```http
PUT /api/v1/clusters/:id
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "updated-cluster"
}
```

### 删除集群

```http
DELETE /api/v1/clusters/:id
Authorization: Bearer {token}
```

---

## 命名空间 API

### 获取命名空间列表

```http
GET /api/v1/namespaces
Authorization: Bearer {token}
```

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": ["default", "kube-system", "monitoring"]
}
```

### 获取命名空间详情

```http
GET /api/v1/namespaces/:name
Authorization: Bearer {token}
```

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "name": "default",
    "status": "Active",
    "age": "365d",
    "labels": {}
  }
}
```

### 创建命名空间

```http
POST /api/v1/namespaces
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "my-namespace",
  "labels": {
    "environment": "development"
  }
}
```

### 删除命名空间

```http
DELETE /api/v1/namespaces/:name
Authorization: Bearer {token}
```

---

## Pod API

### 获取 Pod 列表

**获取指定命名空间的 Pod**

```http
GET /api/v1/namespaces/{namespace}/pods
Authorization: Bearer {token}
```

**获取所有命名空间的 Pod**

```http
GET /api/v1/pods
Authorization: Bearer {token}
```

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "name": "my-pod",
      "namespace": "default",
      "status": "Running",
      "nodeName": "node-1",
      "ip": "10.244.1.5",
      "createdAt": "2026-02-07T10:00:00Z",
      "labels": {
        "app": "my-app"
      },
      "containers": [
        {
          "name": "container-1",
          "image": "nginx:latest",
          "ready": true,
          "restartCount": 0
        }
      ]
    }
  ]
}
```

### 获取 Pod 详情

```http
GET /api/v1/namespaces/{namespace}/pods/{name}
Authorization: Bearer {token}
```

### 获取 Pod 日志

```http
GET /api/v1/namespaces/{namespace}/pods/{name}/logs?tailLines=100
Authorization: Bearer {token}
```

**查询参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| tailLines | number | 返回最近的行数，默认 100 |
| follow | boolean | 是否持续跟踪日志 |
| previous | boolean | 是否查看上次重启的日志 |

### 删除 Pod

```http
DELETE /api/v1/namespaces/{namespace}/pods/{name}
Authorization: Bearer {token}
```

---

## 错误码

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权（Token 无效或过期） |
| 403 | 禁止访问（权限不足） |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 500 | 服务器内部错误 |
| 503 | 服务不可用 |

---

## 示例代码

### cURL

```bash
# 登录获取 Token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'

# 获取集群列表
curl -X GET http://localhost:8080/api/v1/clusters \
  -H "Authorization: Bearer YOUR_TOKEN"

# 获取 Pod 列表
curl -X GET http://localhost:8080/api/v1/namespaces/default/pods \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### JavaScript/Axios

```javascript
import axios from 'axios'

// 创建 Axios 实例
const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器：添加 Token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 登录
async function login(username, password) {
  const { data } = await api.post('/auth/login', { username, password })
  localStorage.setItem('token', data.data.token)
  return data
}

// 获取集群列表
async function getClusters() {
  const { data } = await api.get('/clusters')
  return data.data
}

// 获取 Pod 列表
async function getPods(namespace) {
  const { data } = await api.get(`/namespaces/${namespace}/pods`)
  return data.data
}
```

### TypeScript/Vue 3

```typescript
import { ref } from 'vue'
import request from '@/utils/request'
import type { Cluster, Pod } from '@/types/kube'

// 获取集群列表
export function useClusters() {
  const clusters = ref<Cluster[]>([])
  const loading = ref(false)

  async function fetchClusters() {
    loading.value = true
    try {
      clusters.value = await request.get<Cluster[]>('/clusters')
    } finally {
      loading.value = false
    }
  }

  return { clusters, loading, fetchClusters }
}

// 获取 Pod 列表
export function usePods(namespace: string) {
  const pods = ref<Pod[]>([])
  const loading = ref(false)

  async function fetchPods() {
    loading.value = true
    try {
      pods.value = await request.get<Pod[]>(`/namespaces/${namespace}/pods`)
    } finally {
      loading.value = false
    }
  }

  return { pods, loading, fetchPods }
}
```

---

## 注意事项

1. **Token 过期**：Token 有效期为 24 小时，过期后需要重新登录
2. **分页**：列表类 API 支持分页，使用 `page` 和 `pageSize` 参数
3. **速率限制**：API 有速率限制，建议合理控制请求频率
4. **WebSocket**：部分实时功能使用 WebSocket，详见相关文档

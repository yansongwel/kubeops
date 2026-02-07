/**
 * Kubernetes 资源类型定义
 */

// ============================================================================
// 集群相关
// ============================================================================

export interface Cluster {
  id: string
  name: string
  endpoint: string
  status: 'Connected' | 'Disconnected' | 'Error'
  version: string
  nodeCount: number
  createdAt: string
}

export interface ClusterCreateParams {
  name: string
  endpoint: string
  kubeconfig?: string
}

// ============================================================================
// 命名空间相关
// ============================================================================

export interface Namespace {
  name: string
  status: string
  age: string
  labels: Record<string, string>
}

// ============================================================================
// Pod 相关
// ============================================================================

export interface Container {
  name: string
  image: string
  ready: boolean
  restartCount: number
}

export interface Pod {
  name: string
  namespace: string
  status: 'Running' | 'Pending' | 'Failed' | 'Succeeded' | 'Unknown'
  nodeName: string
  ip: string
  createdAt: string
  labels: Record<string, string>
  containers: Container[]
}

// ============================================================================
// Deployment 相关
// ============================================================================

export interface Deployment {
  name: string
  namespace: string
  ready: number
  desired: number
  updated: number
  available: number
  age: string
}

// ============================================================================
// Service 相关
// ============================================================================

export interface Service {
  name: string
  namespace: string
  type: 'ClusterIP' | 'NodePort' | 'LoadBalancer' | 'ExternalName'
  clusterIP: string
  externalIP: string
  ports: number[]
  age: string
}

// ============================================================================
// ConfigMap 和 Secret 相关
// ============================================================================

export interface ConfigMap {
  name: string
  namespace: string
  data: Record<string, string>
  age: string
}

export interface Secret {
  name: string
  namespace: string
  type: string
  data: Record<string, string> // base64 编码
  age: string
}

// ============================================================================
// API 响应格式
// ============================================================================

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface ApiError {
  code: number
  message: string
  details?: string
}

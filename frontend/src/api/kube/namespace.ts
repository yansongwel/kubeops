/**
 * 命名空间 API
 */
import request from '@/utils/request'
import type { Namespace } from '@/types/kube'

// 获取命名空间列表
export function getNamespaces() {
  return request.get<Namespace[]>('/namespaces')
}

// 获取命名空间详情
export function getNamespace(name: string) {
  return request.get<Namespace>(`/namespaces/${name}`)
}

// 创建命名空间
export function createNamespace(data: { name: string; labels?: Record<string, string> }) {
  return request.post<Namespace>('/namespaces', data)
}

// 删除命名空间
export function deleteNamespace(name: string) {
  return request.delete(`/namespaces/${name}`)
}

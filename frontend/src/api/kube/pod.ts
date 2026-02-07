/**
 * Pod 管理 API
 */
import request from '@/utils/request'
import type { Pod } from '@/types/kube'

// 获取指定命名空间的 Pod 列表
export function getPods(namespace: string) {
  return request.get<Pod[]>(`/namespaces/${namespace}/pods`)
}

// 获取所有命名空间的 Pod
export function getAllPods() {
  return request.get<Pod[]>('/pods')
}

// 获取 Pod 详情
export function getPod(namespace: string, name: string) {
  return request.get<Pod>(`/namespaces/${namespace}/pods/${name}`)
}

// 获取 Pod 日志
export function getPodLogs(namespace: string, name: string, options?: {
  tailLines?: number
  follow?: boolean
  previous?: boolean
}) {
  return request.get<string>(`/namespaces/${namespace}/pods/${name}/logs`, {
    params: options
  })
}

// 删除 Pod
export function deletePod(namespace: string, name: string) {
  return request.delete(`/namespaces/${namespace}/pods/${name}`)
}

// 重启 Pod (通常通过删除实现)
export function restartPod(namespace: string, name: string) {
  return deletePod(namespace, name)
}

// 获取 Pod 容器日志
export function getContainerLogs(
  namespace: string,
  podName: string,
  containerName: string,
  options?: { tailLines?: number; follow?: boolean }
) {
  return request.get<string>(
    `/namespaces/${namespace}/pods/${podName}/logs/${containerName}`,
    { params: options }
  )
}

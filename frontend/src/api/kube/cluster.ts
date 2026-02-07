/**
 * 集群管理 API
 */
import request from '@/utils/request'
import type { ApiResponse, Cluster, ClusterCreateParams } from '@/types/kube'

// 获取集群列表
export function getClusters() {
  return request.get<Cluster[]>('/clusters')
}

// 获取集群详情
export function getCluster(id: string) {
  return request.get<Cluster>(`/clusters/${id}`)
}

// 创建集群
export function createCluster(data: ClusterCreateParams) {
  return request.post<Cluster>('/clusters', data)
}

// 更新集群
export function updateCluster(id: string, data: Partial<Cluster>) {
  return request.put<Cluster>(`/clusters/${id}`, data)
}

// 删除集群
export function deleteCluster(id: string) {
  return request.delete(`/clusters/${id}`)
}

// 测试集群连接
export function testClusterConnection(id: string) {
  return request.post<{ status: string }>(`/clusters/${id}/test`)
}

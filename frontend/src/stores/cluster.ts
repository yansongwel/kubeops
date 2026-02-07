/**
 * 集群状态管理
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Cluster } from '@/types/kube'
import * as clusterApi from '@/api/kube/cluster'

export const useClusterStore = defineStore('cluster', () => {
  // 状态
  const clusters = ref<Cluster[]>([])
  const currentCluster = ref<Cluster | null>(null)
  const loading = ref(false)

  // 计算属性
  const clusterCount = computed(() => clusters.value.length)

  const connectedClusters = computed(() =>
    clusters.value.filter((c) => c.status === 'Connected')
  )

  const currentClusterName = computed(() => currentCluster.value?.name || '')

  // 方法

  /**
   * 获取所有集群
   */
  async function fetchClusters() {
    loading.value = true
    try {
      clusters.value = await clusterApi.getClusters()
    } catch (error) {
      console.error('获取集群列表失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  /**
   * 设置当前集群
   */
  function setCurrentCluster(cluster: Cluster) {
    currentCluster.value = cluster
    localStorage.setItem('currentCluster', JSON.stringify(cluster))
  }

  /**
   * 获取当前集群
   */
  function getCurrentCluster(): Cluster | null {
    if (!currentCluster.value) {
      const saved = localStorage.getItem('currentCluster')
      if (saved) {
        try {
          currentCluster.value = JSON.parse(saved)
        } catch (e) {
          console.error('解析保存的集群信息失败:', e)
        }
      }
    }
    return currentCluster.value
  }

  /**
   * 清除当前集群
   */
  function clearCurrentCluster() {
    currentCluster.value = null
    localStorage.removeItem('currentCluster')
  }

  /**
   * 创建集群
   */
  async function createCluster(data: Parameters<typeof clusterApi.createCluster>[0]) {
    const newCluster = await clusterApi.createCluster(data)
    clusters.value.push(newCluster)
    return newCluster
  }

  /**
   * 更新集群
   */
  async function updateCluster(id: string, data: Parameters<typeof clusterApi.updateCluster>[1]) {
    const updated = await clusterApi.updateCluster(id, data)
    const index = clusters.value.findIndex((c) => c.id === id)
    if (index !== -1) {
      clusters.value[index] = updated
    }
    return updated
  }

  /**
   * 删除集群
   */
  async function removeCluster(id: string) {
    await clusterApi.deleteCluster(id)
    clusters.value = clusters.value.filter((c) => c.id !== id)
    if (currentCluster.value?.id === id) {
      clearCurrentCluster()
    }
  }

  return {
    // 状态
    clusters,
    currentCluster,
    loading,

    // 计算属性
    clusterCount,
    connectedClusters,
    currentClusterName,

    // 方法
    fetchClusters,
    setCurrentCluster,
    getCurrentCluster,
    clearCurrentCluster,
    createCluster,
    updateCluster,
    removeCluster,
  }
})

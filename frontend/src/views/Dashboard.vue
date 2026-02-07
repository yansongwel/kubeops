<template>
  <div class="dashboard">
    <el-row :gutter="16">
      <el-col :span="6">
        <el-card class="stat-card">
          <el-statistic title="集群数量" :value="stats.clusters">
            <template #suffix>
              <el-icon><Connection /></el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <el-statistic title="Pods" :value="stats.pods">
            <template #suffix>
              <el-icon><Box /></el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <el-statistic title="命名空间" :value="stats.namespaces">
            <template #suffix>
              <el-icon><Folder /></el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <el-statistic title="健康评分" :value="stats.healthScore" suffix="%">
            <template #prefix>
              <el-icon><CircleCheck /></el-icon>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="12">
        <el-card title="资源使用趋势">
          <template #header>
            <span>资源使用趋势</span>
          </template>
          <div style="height: 300px">图表区域</div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card title="AI 巡检报告">
          <template #header>
            <span>AI 巡检报告</span>
          </template>
          <el-table :data="aiReports" style="width: 100%">
            <el-table-column prop="type" label="类型" width="100" />
            <el-table-column prop="severity" label="严重性" width="80">
              <template #default="{ row }">
                <el-tag v-if="row.severity === 'high'" type="danger">高</el-tag>
                <el-tag v-else-if="row.severity === 'medium'" type="warning">中</el-tag>
                <el-tag v-else type="info">低</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="message" label="描述" />
            <el-table-column prop="time" label="时间" width="180" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const stats = ref({
  clusters: 3,
  pods: 156,
  namespaces: 12,
  healthScore: 95,
})

const aiReports = ref([
  {
    type: '资源优化',
    severity: 'medium',
    message: '命名空间 default 中有 3 个 Pod 资源利用率低于 20%',
    time: '2024-01-15 10:30:00',
  },
  {
    type: '安全',
    severity: 'high',
    message: '发现 2 个 Pod 使用特权模式运行',
    time: '2024-01-15 10:25:00',
  },
  {
    type: '可靠性',
    severity: 'low',
    message: 'Deployment nginx 没有设置资源限制',
    time: '2024-01-15 10:20:00',
  },
])
</script>

<style scoped>
.dashboard {
  height: 100%;
}

.stat-card {
  text-align: center;
}
</style>

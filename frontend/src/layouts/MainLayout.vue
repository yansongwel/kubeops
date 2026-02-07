<template>
  <el-container class="layout-container">
    <el-aside width="240px" class="sidebar">
      <div class="logo">
        <h2>KubeOps</h2>
      </div>
      <el-menu
        :default-active="activeMenu"
        class="sidebar-menu"
        router
        background-color="#001529"
        text-color="#fff"
        active-text-color="#1890ff"
      >
        <el-menu-item index="/dashboard">
          <el-icon><Dashboard /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/clusters">
          <el-icon><Connection /></el-icon>
          <span>集群管理</span>
        </el-menu-item>
        <el-sub-menu index="workloads">
          <template #title>
            <el-icon><Box /></el-icon>
            <span>工作负载</span>
          </template>
          <el-menu-item index="/workloads/pods">Pods</el-menu-item>
          <el-menu-item index="/workloads/deployments">Deployments</el-menu-item>
          <el-menu-item index="/workloads/statefulsets">StatefulSets</el-menu-item>
          <el-menu-item index="/workloads/daemonsets">DaemonSets</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="/ai-inspector">
          <el-icon><MagicStick /></el-icon>
          <span>AI 巡检</span>
        </el-menu-item>
        <el-menu-item index="/devops">
          <el-icon><Opportunity /></el-icon>
          <span>DevOps</span>
        </el-menu-item>
        <el-menu-item index="/logs">
          <el-icon><Document /></el-icon>
          <span>日志平台</span>
        </el-menu-item>
        <el-menu-item index="/monitoring">
          <el-icon><TrendCharts /></el-icon>
          <span>监控平台</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <div class="header-left">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item>{{ currentRoute }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-button text @click="handleNotification">
            <el-icon><Bell /></el-icon>
          </el-button>
          <el-dropdown>
            <span class="user-dropdown">
              <el-avatar :size="32" :src="userAvatar" />
              <span class="username">Admin</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>个人设置</el-dropdown-item>
                <el-dropdown-item divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const activeMenu = computed(() => route.path)
const currentRoute = computed(() => route.meta.title as string || 'Home')
const userAvatar = ref('')

const handleNotification = () => {
  console.log('Notifications clicked')
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.sidebar {
  background-color: #001529;
  overflow-x: hidden;
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 20px;
  font-weight: bold;
  border-bottom: 1px solid #002140;
}

.sidebar-menu {
  border-right: none;
  height: calc(100vh - 64px);
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #fff;
  border-bottom: 1px solid #f0f0f0;
  padding: 0 24px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-dropdown {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.username {
  color: #333;
}

.main-content {
  background: #f0f2f5;
  padding: 24px;
}
</style>

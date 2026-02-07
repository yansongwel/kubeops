import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Layout',
    component: () => import('@/layouts/MainLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: '/dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '仪表盘' },
      },
      {
        path: '/clusters',
        name: 'Clusters',
        component: () => import('@/views/Clusters.vue'),
        meta: { title: '集群管理' },
      },
      {
        path: '/workloads',
        name: 'Workloads',
        component: () => import('@/views/Workloads.vue'),
        meta: { title: '工作负载' },
      },
      {
        path: '/ai-inspector',
        name: 'AIInspector',
        component: () => import('@/views/AIInspector.vue'),
        meta: { title: 'AI 巡检' },
      },
      {
        path: '/devops',
        name: 'DevOps',
        component: () => import('@/views/DevOps.vue'),
        meta: { title: 'DevOps' },
      },
      {
        path: '/logs',
        name: 'Logs',
        component: () => import('@/views/Logs.vue'),
        meta: { title: '日志平台' },
      },
      {
        path: '/monitoring',
        name: 'Monitoring',
        component: () => import('@/views/Monitoring.vue'),
        meta: { title: '监控平台' },
      },
    ],
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录' },
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

export default router

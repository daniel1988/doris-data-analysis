import type { RouteRecordRaw } from 'vue-router'
import { createRouter, createWebHistory } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/Login.vue'),
    meta: { title: '登录', public: true }
  },
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/dashboard/index.vue'),
    meta: { title: '看板管理' }
  },
  {
    path: '/analytics',
    redirect: '/analytics/event'
  },
  {
    path: '/analytics/ai',
    name: 'AIAnalysis',
    component: () => import('@/views/analytics/ai/index.vue'),
    meta: { title: '智能分析' }
  },
  {
    path: '/analytics/event',
    name: 'EventAnalysis',
    component: () => import('@/views/analytics/EventAnalysis.vue'),
    meta: { title: '事件分析' }
  },
  {
    path: '/analytics/retention',
    name: 'RetentionAnalysis',
    component: () => import('@/views/analytics/Retention.vue'),
    meta: { title: '留存分析' }
  },
  {
    path: '/analytics/funnel',
    name: 'FunnelAnalysis',
    component: () => import('@/views/analytics/Funnel.vue'),
    meta: { title: '漏斗分析' }
  },
  {
    path: '/analytics/scatter',
    name: 'ScatterAnalysis',
    component: () => import('@/views/analytics/ScatterAnalysis.vue'),
    meta: { title: '分布分析' }
  },
  {
    path: '/analytics/property',
    name: 'PropertyAnalysis',
    component: () => import('@/views/analytics/PropertyAnalysis.vue'),
    meta: { title: '属性分析' }
  },
  {
    path: '/metadata',
    redirect: '/metadata/event'
  },
  {
    path: '/metadata/event',
    name: 'MetadataEvent',
    component: () => import('@/views/metadata/ProjectEvent.vue'),
    meta: { title: '元事件' }
  },
  {
    path: '/metadata/field',
    name: 'MetadataField',
    component: () => import('@/views/metadata/ProjectProperty.vue'),
    meta: { title: '事件属性' }
  },
  {
    path: '/metadata/dimension',
    name: 'MetadataDimension',
    component: () => import('@/views/metadata/Dimension.vue'),
    meta: { title: '维度' }
  },
  {
    path: '/metadata/metrics',
    name: 'MetadataMetrics',
    component: () => import('@/views/metadata/Metrics.vue'),
    meta: { title: '指标中心' }
  },
  {
    path: '/metadata/log',
    name: 'LogDetail',
    component: () => import('@/views/metadata/LogDetail.vue'),
    meta: { title: '日志明细' }
  },
  {
    path: '/user',
    redirect: '/user/list'
  },
  {
    path: '/user/list',
    name: 'UserList',
    component: () => import('@/views/user/list/index.vue'),
    meta: { title: '用户列表' }
  },
  {
    path: '/system',
    redirect: '/system/users'
  },
  {
    path: '/system/projects',
    name: 'ProjectManagement',
    component: () => import('@/views/system/Projects.vue'),
    meta: { title: '项目管理' }
  },
  {
    path: '/system/users',
    name: 'UserManagement',
    component: () => import('@/views/system/Users.vue'),
    meta: { title: '用户管理' }
  },
  {
    path: '/system/roles',
    name: 'RoleManagement',
    component: () => import('@/views/system/Roles.vue'),
    meta: { title: '角色管理' }
  },
  {
    path: '/system/tenants',
    name: 'TenantManagement',
    component: () => import('@/views/system/Tenants.vue'),
    meta: { title: '租户管理' }
  },
  {
    path: '/system/ai-models',
    name: 'AIModelManagement',
    component: () => import('@/views/system/AIModels.vue'),
    meta: { title: 'AI 模型配置' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router

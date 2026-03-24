<template>
  <div class="app-wrapper">
    <!-- 顶部导航栏 -->
    <header class="navbar">
      <div class="navbar-left">
        <div class="logo">DMP Admin V2</div>
      </div>

      <div class="navbar-center">
        <el-menu :default-active="activeModule" mode="horizontal" background-color="#fff" text-color="#303133"
          active-text-color="#409eff" router class="module-menu" :ellipsis="false">
          <el-menu-item index="/dashboard">看板</el-menu-item>

          <!-- 分析模块下拉 -->
          <el-sub-menu index="/analytics">
            <template #title>分析</template>
            <el-menu-item index="/analytics/ai">智能分析</el-menu-item>
            <el-menu-item index="/analytics/event">事件分析</el-menu-item>
            <el-menu-item index="/analytics/retention">留存分析</el-menu-item>
            <el-menu-item index="/analytics/funnel">漏斗分析</el-menu-item>
            <el-menu-item index="/analytics/scatter">分布分析</el-menu-item>
            <el-menu-item index="/analytics/property">属性分析</el-menu-item>
          </el-sub-menu>

          <!-- 元数据模块下拉 -->
          <el-sub-menu index="/metadata">
            <template #title>元数据</template>
            <el-menu-item index="/metadata/field">事件属性</el-menu-item>
            <el-menu-item index="/metadata/event">元事件</el-menu-item>
            <el-menu-item index="/metadata/dimension">维度</el-menu-item>
            <el-menu-item index="/metadata/metrics">指标中心</el-menu-item>
            <el-menu-item index="/metadata/log">日志明细</el-menu-item>
          </el-sub-menu>

          <!-- 用户模块下拉 -->
          <el-sub-menu index="/user">
            <template #title>用户</template>
            <el-menu-item index="/user/list">用户列表</el-menu-item>
          </el-sub-menu>

          <!-- 系统管理下拉 -->
          <el-sub-menu index="/system">
            <template #title>系统管理</template>
            <el-menu-item index="/system/projects">项目管理</el-menu-item>
            <el-menu-item index="/system/users">用户管理</el-menu-item>
            <el-menu-item index="/system/roles">角色管理</el-menu-item>
            <el-menu-item index="/system/tenants">租户管理</el-menu-item>
            <el-menu-item index="/system/ai-models">AI 模型配置</el-menu-item>
          </el-sub-menu>
        </el-menu>
      </div>

      <div class="navbar-right">
        <div class="project-selector" style="margin-right: 30px">
          <el-select v-model="appStore.activeProjectAlias" placeholder="切换项目" filterable size="default"
            style="width: 220px" @change="handleProjectChange">
            <el-option v-for="item in appStore.projectList" :key="item.project_alias" :label="item.project_name"
              :value="item.project_alias">
              <span style="float: left">{{ item.project_name }}</span>
              <span style="float: right; color: var(--el-text-color-secondary); font-size: 12px">
                {{ item.project_alias }}
              </span>
            </el-option>
          </el-select>
        </div>
        <el-dropdown trigger="click">
          <span class="user-info">
            {{ userStore.userInfo?.name || '管理员' }}
            <el-icon>
              <ArrowDown />
            </el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>个人中心</el-dropdown-item>
              <el-dropdown-item divided @click="handleLogout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </header>

    <div class="app-container">
      <!-- 侧边栏 -->
      <aside v-if="showSidebar" :class="['sidebar-container', { collapsed: !appStore.sidebarOpened }]">
        <div class="sidebar-toggle" @click="appStore.toggleSidebar">
          <el-icon>
            <Fold v-if="appStore.sidebarOpened" />
            <Expand v-else />
          </el-icon>
        </div>

        <el-scrollbar v-show="appStore.sidebarOpened">
          <!-- 二级菜单 -->
          <el-menu :default-active="activeSubMenu" :collapse="!appStore.sidebarOpened" background-color="#fff"
            text-color="#303133" active-text-color="#409eff" unique-opened router class="sidebar-menu">
            <el-menu-item v-for="item in subMenuItems" :key="item.path" :index="item.path">
              <el-icon>
                <component :is="item.icon" />
              </el-icon>
              <template #title>{{ item.title }}</template>
            </el-menu-item>
          </el-menu>
        </el-scrollbar>
      </aside>

      <!-- 主体内容 -->
      <div class="main-container">
        <div class="tags-view">
          <el-tabs v-model="appStore.currentTab" type="card" closable @tab-click="handleTabClick"
            @tab-remove="handleTabRemove">
            <el-tab-pane v-for="tab in appStore.activeTabs" :key="tab.path" :label="tab.title" :name="tab.path" />
          </el-tabs>
        </div>

        <main class="app-main">
          <div class="content-wrapper">
            <slot />
          </div>
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAppStore } from '@/store/app'
import { useUserStore } from '@/store/user'
import {
  ArrowDown, Expand, Fold
} from '@element-plus/icons-vue'
import { computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const appStore = useAppStore()

onMounted(() => {
  appStore.fetchProjects()
  // 初始化当前页面的 tab
  if (route.meta?.title) {
    appStore.addTab({
      path: route.path,
      title: route.meta.title
    })
  }
})

// 监听路由变化，自动添加 tab
watch(() => route.path, () => {
  if (route.meta?.title) {
    appStore.addTab({
      path: route.path,
      title: route.meta.title
    })
  }
})

const handleProjectChange = (val: string) => {
  appStore.setActiveProject(val)
  // 刷新当前页面以重新加载元数据
  window.location.reload()
}

const activeModule = computed(() => {
  const pathParts = route.path.split('/')
  return pathParts.length > 1 ? `/${pathParts[1]}` : '/dashboard'
})

const subMenuItems = computed(() => {
  const module = activeModule.value
  if (module === '/metadata') {
    return [
      { path: '/metadata/event', title: '元事件', icon: 'Collection' },
      { path: '/metadata/field', title: '事件属性', icon: 'List' },
      { path: '/metadata/dimension', title: '维度', icon: 'Management' },
      { path: '/metadata/metrics', title: '指标中心', icon: 'DataLine' },
      { path: '/metadata/log', title: '日志明细', icon: 'Memo' }
    ]
  } else if (module === '/system') {
    return [
      { path: '/system/projects', title: '项目管理', icon: 'Folder' },
      { path: '/system/users', title: '用户管理', icon: 'User' },
      { path: '/system/roles', title: '角色管理', icon: 'Lock' },
      { path: '/system/tenants', title: '租户管理', icon: 'OfficeBuilding' },
      { path: '/system/ai-models', title: 'AI 模型配置', icon: 'Setting' }
    ]
  } else if (module === '/analytics') {
    return [
      { path: '/analytics/ai', title: '智能分析', icon: 'Cpu' },
      { path: '/analytics/event', title: '事件分析', icon: 'TrendCharts' },
      { path: '/analytics/retention', title: '留存分析', icon: 'Timer' },
      { path: '/analytics/funnel', title: '漏斗分析', icon: 'Filter' },
      { path: '/analytics/scatter', title: '分布分析', icon: 'TrendCharts' },
      { path: '/analytics/property', title: '属性分析', icon: 'User' }
    ]
  } else if (module === '/user') {
    return [
      { path: '/user/list', title: '用户列表', icon: 'UserFilled' }
    ]
  }
  return []
})

const showSidebar = computed(() => subMenuItems.value.length > 0)
const activeSubMenu = computed(() => route.path)

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

const handleTabClick = (tab: any) => {
  router.push(tab.props.name)
}

const handleTabRemove = (path: string) => {
  appStore.removeTab(path)
}
</script>

<style scoped>
.app-wrapper {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100%;
}

.navbar {
  height: 60px;
  background: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  flex-shrink: 0;
  z-index: 1000;
}

.navbar-left {
  display: flex;
  align-items: center;
}

.navbar-center {
  flex: 1;
  display: flex;
  justify-content: center;
  min-width: 0;
  /* 允许 flex 项目缩小 */
}

.navbar-right {
  width: 200px;
  display: flex;
  justify-content: flex-end;
}

.logo {
  font-size: 20px;
  font-weight: bold;
  color: var(--el-color-primary);
  margin-right: 10px;
}

.project-selector {
  margin-left: 10px;
}

.module-menu {
  border-bottom: none;
}

:deep(.el-sub-menu__title) {
  border-bottom: none !important;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  color: var(--el-text-color-primary);
  white-space: nowrap;
  padding: 0 8px;
  gap: 5px;
}

.user-info:hover {
  color: var(--el-color-primary);
}

.user-info .el-icon {
  margin-left: 4px;
}

.app-container {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.sidebar-container {
  width: 240px;
  background-color: #fff;
  border-right: 1px solid #e6e6e6;
  height: 100%;
  flex-shrink: 0;
  transition: width 0.3s;
  display: flex;
  flex-direction: column;
}

.sidebar-container.collapsed {
  width: 40px;
}

.sidebar-toggle {
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-bottom: 1px solid #f0f0f0;
  color: #909399;
}

.config-sidebar {
  padding: 20px 15px;
}

.sidebar-title {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 20px;
  color: #303133;
}

.sidebar-menu {
  border-right: none;
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: #f0f2f5;
}

.tags-view {
  height: 34px;
  background: #fff;
  border-bottom: 1px solid #d8dce5;
  padding: 2px 15px 0;
}

.tags-view :deep(.el-tabs--card > .el-tabs__header) {
  border-bottom: none;
  margin-bottom: 0;
}

.app-main {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.content-wrapper {
  flex: 1;
  padding: 20px;
  overflow: auto;
  display: flex;
  flex-direction: column;
}

.filter-tag {
  margin-right: 5px;
  margin-bottom: 5px;
}
</style>

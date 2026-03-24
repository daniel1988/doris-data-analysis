<template>
  <div class="dashboard-container">
    <!-- 左侧看板列表 -->
    <aside class="dashboard-sidebar">
      <div class="sidebar-header">
        <span>看板列表</span>
        <el-button type="primary" link icon="Plus" @click="showCreateDialog = true">新建</el-button>
      </div>
      <el-scrollbar>
        <div v-for="item in dashboardList" :key="item.id"
          :class="['dashboard-list-item', { active: currentDashboardId === item.id }]" @click="selectDashboard(item)">
          <el-icon>
            <DataAnalysis />
          </el-icon>
          <span class="item-name">{{ item.display_name || item.name }}</span>
          <el-popconfirm title="确定删除该看板吗？" @confirm="handleDeleteDashboard(item.id!)">
            <template #reference>
              <el-icon class="delete-icon" @click.stop>
                <Delete />
              </el-icon>
            </template>
          </el-popconfirm>
        </div>
        <el-empty v-if="dashboardList.length === 0" description="暂无看板" :image-size="60" />
      </el-scrollbar>
    </aside>

    <!-- 右侧看板内容区 -->
    <main class="dashboard-main" v-loading="loading">
      <template v-if="currentDashboard">
        <div class="dashboard-header">
          <div class="header-left">
            <h2 class="dashboard-title">{{ currentDashboard.display_name || currentDashboard.name }}</h2>
            <p v-if="currentDashboard.description" class="dashboard-desc">{{ currentDashboard.description }}</p>
          </div>
          <div class="dashboard-actions">
            <el-button-group>
              <el-button :icon="Plus" @click="showAddReportDialog = true" v-if="!isEditMode">添加报表</el-button>
              <el-button :icon="Refresh" @click="refreshAll">刷新</el-button>
              <el-button :icon="isEditMode ? Check : Edit" :type="isEditMode ? 'success' : 'primary'"
                @click="toggleEditMode">
                {{ isEditMode ? '保存布局' : '编辑布局' }}
              </el-button>
            </el-button-group>
          </div>
        </div>

        <!-- 全局筛选栏 -->
        <DashboardFilterBar v-model="globalFilters" @change="handleFilterChange" v-if="!isEditMode" />

        <!-- 栅格布局展示图表卡片 -->
        <div class="dashboard-content">
          <grid-layout v-if="layout && layout.length > 0" v-model:layout="layout" :col-num="24" :row-height="30"
            :is-draggable="isEditMode" :is-resizable="isEditMode" :vertical-compact="true" :margin="[10, 10]"
            :use-css-transforms="true" style="min-height: 500px;">
            <grid-item v-for="item in layout" :key="item.i" :x="item.x" :y="item.y" :w="item.w" :h="item.h" :i="item.i"
              class="grid-item-container">
              <DashboardItemCard v-if="findItemById(item.i)" :item="findItemById(item.i)"
                :global-filters="globalFilters" @remove="handleRemoveItem" @update-size="handleUpdateItemSize" />
            </grid-item>
          </grid-layout>
          <el-empty v-else description="看板内暂无报表，请从分析页面保存并添加，或点击右上角“添加报表”按钮" />
        </div>
      </template>
      <template v-else-if="!loading">
        <el-empty description="请选择或创建一个看板" />
      </template>
    </main>

    <!-- 新建看板对话框 -->
    <el-dialog v-model="showCreateDialog" title="新建看板" width="450px">
      <el-form :model="createForm" label-width="80px" ref="createFormRef" :rules="formRules">
        <el-form-item label="看板名称" prop="name">
          <el-input v-model="createForm.name" placeholder="请输入看板名称" />
        </el-form-item>
        <el-form-item label="显示名称" prop="display_name">
          <el-input v-model="createForm.display_name" placeholder="用于菜单展示的名称" />
        </el-form-item>
        <el-form-item label="看板描述">
          <el-input v-model="createForm.description" type="textarea" :rows="3" placeholder="可选描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreateDashboard" :loading="creating">确定</el-button>
      </template>
    </el-dialog>

    <!-- 添加报表对话框 -->
    <AddReportDialog v-if="currentDashboardId !== null" v-model="showAddReportDialog"
      :dashboard-id="currentDashboardId as number" @confirm="handleBatchAddReports" />
  </div>
</template>

<script setup lang="ts">
import type { Dashboard, DashboardItem } from '@/api/dashboard'
import {
  addDashboardItem,
  batchUpdateDashboardItems,
  createDashboard,
  deleteDashboard,
  deleteDashboardItem,
  getDashboardById,
  getDashboardList
} from '@/api/dashboard'
import type { Report } from '@/api/report'
import { useAppStore } from '@/store/app'
import { Check, DataAnalysis, Delete, Edit, Plus, Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { GridItem, GridLayout } from 'vue3-grid-layout'
import AddReportDialog from './components/AddReportDialog.vue'
import DashboardFilterBar from './components/DashboardFilterBar.vue'
import DashboardItemCard from './components/DashboardItemCard.vue'

const appStore = useAppStore()
const loading = ref(false)
const creating = ref(false)
const showCreateDialog = ref(false)
const showAddReportDialog = ref(false)
const isEditMode = ref(false)
const dashboardList = ref<Dashboard[]>([])
const currentDashboardId = ref<number | null>(null)
const currentDashboard = ref<Dashboard | null>(null)
const layout = ref<any[]>([])

const globalFilters = ref({
  timeRange: null as [string, string] | null,
  filters: [] as any[]
})

const createForm = ref({
  name: '',
  display_name: '',
  description: ''
})

const formRules = {
  name: [{ required: true, message: '请输入看板名称', trigger: 'blur' }],
  display_name: [{ required: true, message: '请输入显示名称', trigger: 'blur' }]
}

const fetchDashboards = async () => {
  if (!appStore.activeProjectAlias) {
    dashboardList.value = []
    return
  }
  loading.value = true
  try {
    const res = await getDashboardList({ project_alias: appStore.activeProjectAlias })
    dashboardList.value = res as unknown as Dashboard[]
    if (dashboardList.value.length > 0 && !currentDashboardId.value) {
      await selectDashboard(dashboardList.value[0])
    }
  } catch (error) {
    console.error('Failed to fetch dashboards:', error)
    ElMessage.error('获取看板列表失败')
  } finally {
    loading.value = false
  }
}

const selectDashboard = async (item: Dashboard) => {
  currentDashboardId.value = item.id!
  loading.value = true
  try {
    const res = await getDashboardById(item.id!)
    currentDashboard.value = res as unknown as Dashboard
    syncLayoutFromItems()
  } catch (error) {
    console.error('Failed to fetch dashboard detail:', error)
    ElMessage.error('获取看板详情失败')
  } finally {
    loading.value = false
  }
}

const syncLayoutFromItems = () => {
  if (!currentDashboard.value?.items) {
    layout.value = []
    return
  }
  // 必须确保 layout 的每一项都有完整的属性，且 i 为 string
  layout.value = currentDashboard.value.items.map(item => ({
    x: item.position_x || 0,
    y: item.position_y || 0,
    w: item.width || 12,
    h: item.height || 8,
    i: String(item.id)
  }))
}

const handleFilterChange = () => {
  // 全局筛选器变化时，所有卡片会通过 watch 自动刷新
}

const findItemById = (id: number | string | undefined) => {
  if (id === undefined) return null
  return currentDashboard.value?.items?.find(item => item.id === Number(id)) as DashboardItem
}

const handleCreateDashboard = async () => {
  creating.value = true
  try {
    const payload: Dashboard = {
      project_alias: appStore.activeProjectAlias,
      name: createForm.value.name,
      display_name: createForm.value.display_name,
      description: createForm.value.description,
      category: 'custom',
      layout_type: 'grid',
      status: 'published'
    }
    const res = await createDashboard(payload)
    ElMessage.success('看板创建成功')
    showCreateDialog.value = false
    await fetchDashboards()
    selectDashboard(res as unknown as Dashboard)
  } catch (error) {
    console.error('Failed to create dashboard:', error)
  } finally {
    creating.value = false
  }
}

const handleDeleteDashboard = async (id: number) => {
  try {
    await deleteDashboard(id)
    ElMessage.success('删除成功')
    if (currentDashboardId.value === id) {
      currentDashboardId.value = null
      currentDashboard.value = null
      layout.value = []
    }
    fetchDashboards()
  } catch (error) {
    console.error('Failed to delete dashboard:', error)
  }
}

const handleRemoveItem = async (itemId: number) => {
  try {
    await deleteDashboardItem(itemId)
    ElMessage.success('移除成功')
    if (currentDashboard.value?.items) {
      currentDashboard.value.items = currentDashboard.value.items.filter(i => i.id !== itemId)
      syncLayoutFromItems()
    }
  } catch (error) {
    console.error('Failed to remove item:', error)
  }
}

const handleUpdateItemSize = async (itemId: number, width: number) => {
  const item = layout.value.find(l => Number(l.i) === itemId)
  if (!item) return

  // 更新本地布局状态
  item.w = width
  
  // 实时持久化保存
  try {
    loading.value = true
    const updateData = layout.value.map(l => ({
      id: Number(l.i),
      position_x: l.x,
      position_y: l.y,
      width: l.w,
      height: l.h,
      report_id: 0,
      type: '',
      title: ''
    }))
    await batchUpdateDashboardItems(updateData)
    
    // 更新 currentDashboard 中的原始数据以保持同步
    if (currentDashboard.value?.items) {
      const originalItem = currentDashboard.value.items.find(i => i.id === itemId)
      if (originalItem) originalItem.width = width
    }
    
    ElMessage.success('尺寸更新成功')
  } catch (error) {
    console.error('Failed to update item size:', error)
    ElMessage.error('尺寸更新失败')
    // 如果失败了，建议重新同步布局
    syncLayoutFromItems()
  } finally {
    loading.value = false
  }
}

const handleBatchAddReports = async (reports: Report[]) => {
  if (!currentDashboardId.value) return

  loading.value = true
  try {
    // 计算起始 Y 轴
    let startY = 0
    if (layout.value.length > 0) {
      startY = Math.max(...layout.value.map(l => l.y + l.h))
    }

    for (let i = 0; i < reports.length; i++) {
      const report = reports[i]
      await addDashboardItem({
        dashboard_id: currentDashboardId.value,
        report_id: report.id!,
        type: 'chart',
        title: report.name,
        position_x: (i % 2) * 12, // 简单分两列
        position_y: startY + Math.floor(i / 2) * 8,
        width: 12,
        height: 8,
        is_visible: true
      })
    }
    ElMessage.success(`成功添加 ${reports.length} 个报表`)
    await selectDashboard(currentDashboard.value!)
  } catch (error) {
    console.error('Failed to batch add reports:', error)
  } finally {
    loading.value = false
  }
}

const toggleEditMode = async () => {
  if (isEditMode.value) {
    // 保存布局
    try {
      loading.value = true
      const updateData = layout.value.map(l => ({
        id: Number(l.i),
        position_x: l.x,
        position_y: l.y,
        width: l.w,
        height: l.h,
        report_id: 0, // 占位，后端不更新
        type: '',
        title: ''
      }))
      await batchUpdateDashboardItems(updateData)
      ElMessage.success('布局保存成功')
    } catch (error) {
      console.error('Failed to save layout:', error)
    } finally {
      loading.value = false
    }
  }
  isEditMode.value = !isEditMode.value
}

const refreshAll = () => {
  if (currentDashboardId.value) {
    // 重新获取详情，这会触发 syncLayoutFromItems，并由于 key 的存在可能触发组件重新挂载
    selectDashboard(currentDashboard.value!)
    // 也可以通过修改 globalFilters 触发所有卡片 reload
    globalFilters.value = { ...globalFilters.value }
  }
}

let refreshTimer: any = null

const startAutoRefresh = () => {
  stopAutoRefresh()
  if (currentDashboard.value?.refresh_interval && currentDashboard.value.refresh_interval > 0) {
    refreshTimer = setInterval(() => {
      handleFilterChange()
    }, currentDashboard.value.refresh_interval * 1000)
  }
}

const stopAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

onMounted(() => {
  fetchDashboards()
})

onUnmounted(() => {
  stopAutoRefresh()
})

watch(currentDashboard, () => {
  startAutoRefresh()
})

watch(() => appStore.activeProjectAlias, () => {
  currentDashboardId.value = null
  currentDashboard.value = null
  fetchDashboards()
})
</script>

<style scoped>
.dashboard-container {
  display: flex;
  height: 100%;
  margin: -20px;
}

.dashboard-sidebar {
  width: 240px;
  background: #fff;
  border-right: 1px solid #e6e6e6;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.sidebar-header {
  padding: 15px;
  font-weight: bold;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.dashboard-list-item {
  padding: 12px 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 10px;
  transition: all 0.2s;
  color: #606266;
  font-size: 14px;
  position: relative;
}

.dashboard-list-item:hover {
  background: #f5f7fa;
  color: #409eff;
}

.dashboard-list-item:hover .delete-icon {
  opacity: 1;
}

.dashboard-list-item.active {
  background: #ecf5ff;
  color: #409eff;
  border-right: 2px solid #409eff;
}

.item-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.delete-icon {
  opacity: 0;
  transition: opacity 0.2s;
  font-size: 14px;
}

.delete-icon:hover {
  color: #f56c6c;
}

.dashboard-main {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 20px;
  background: #f0f2f5;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 10px;
}

.dashboard-title {
  margin: 0;
  font-size: 20px;
  color: #303133;
}

.dashboard-desc {
  margin: 5px 0 0;
  font-size: 13px;
  color: #909399;
}

.dashboard-actions {
  display: flex;
  gap: 12px;
}

.dashboard-content {
  flex: 1;
}

.grid-item-container {
  background: transparent;
}

:deep(.vue-grid-layout) {
  background: transparent;
}

:deep(.vue-grid-item.vue-grid-placeholder) {
  background: #409eff !important;
  border-radius: 4px;
  opacity: 0.1 !important;
}

:deep(.vue-grid-item.resizing) {
  opacity: 0.9;
}

:deep(.vue-grid-item.static) {
  background: #cce;
}
</style>

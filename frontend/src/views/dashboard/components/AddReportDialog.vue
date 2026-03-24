<template>
  <el-dialog
    v-model="visible"
    title="添加报表至看板"
    width="600px"
    destroy-on-close
  >
    <div class="add-report-dialog">
      <el-input
        v-model="searchText"
        placeholder="搜索报表名称..."
        prefix-icon="Search"
        clearable
        class="search-input"
      />
      
      <el-table
        v-loading="loading"
        :data="filteredReports"
        height="350px"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column property="name" label="报表名称" show-overflow-tooltip />
        <el-table-column property="category" label="类型" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ categoryMap[row.category] || row.category }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column property="update_time" label="更新时间" width="160">
          <template #default="{ row }">
            {{ formatTime(row.update_time) }}
          </template>
        </el-table-column>
      </el-table>
    </div>
    
    <template #footer>
      <div class="dialog-footer">
        <span class="selected-count">已选 {{ selectedReports.length }} 个报表</span>
        <div>
          <el-button @click="visible = false">取消</el-button>
          <el-button 
            type="primary" 
            :disabled="selectedReports.length === 0" 
            :loading="submitting"
            @click="handleConfirm"
          >
            确定添加
          </el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { getReportList } from '@/api/report'
import type { Report } from '@/api/report'
import { useAppStore } from '@/store/app'
import dayjs from 'dayjs'

const props = defineProps<{
  modelValue: boolean
  dashboardId: number
}>()

const emit = defineEmits(['update:modelValue', 'confirm'])

const appStore = useAppStore()
const visible = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v)
})

const loading = ref(false)
const submitting = ref(false)
const searchText = ref('')
const allReports = ref<Report[]>([])
const selectedReports = ref<Report[]>([])

const categoryMap: Record<string, string> = {
  'events': '事件分析',
  'funnel': '漏斗分析',
  'retention': '留存分析',
  'property': '属性分析',
  'scatter': '分布分析'
}

const filteredReports = computed(() => {
  if (!searchText.value) return allReports.value
  const query = searchText.value.toLowerCase()
  return allReports.value.filter(r => 
    r.name.toLowerCase().includes(query) || 
    (r.description && r.description.toLowerCase().includes(query))
  )
})

const fetchReports = async () => {
  if (!appStore.activeProjectAlias) return
  loading.value = true
  try {
    const res = await getReportList({ project_alias: appStore.activeProjectAlias })
    allReports.value = res as unknown as Report[]
  } catch (error) {
    console.error('Failed to fetch reports:', error)
  } finally {
    loading.value = false
  }
}

const handleSelectionChange = (val: Report[]) => {
  selectedReports.value = val
}

const handleConfirm = () => {
  emit('confirm', selectedReports.value)
  visible.value = false
}

const formatTime = (time?: string) => {
  return time ? dayjs(time).format('YYYY-MM-DD HH:mm') : '-'
}

watch(visible, (val) => {
  if (val) {
    searchText.value = ''
    selectedReports.value = []
    fetchReports()
  }
})
</script>

<style scoped>
.add-report-dialog {
  display: flex;
  flex-direction: column;
  gap: 15px;
}
.search-input {
  width: 100%;
}
.dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}
.selected-count {
  font-size: 13px;
  color: #606266;
}
</style>

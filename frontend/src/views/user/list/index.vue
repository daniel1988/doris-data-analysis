<template>
  <div class="user-list-container">
    <el-card class="filter-card" shadow="never">
      <el-form :model="filterForm" inline class="compact-form">
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="filterForm.timeRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            size="default"
          />
        </el-form-item>

        <el-form-item label="用户ID">
          <DimensionValueSelect
            v-model="filterForm.userId"
            :project-alias="projectAlias"
            :table="TableNames.USER"
            field="u_openid"
            placeholder="请选择或搜索用户ID"
          />
        </el-form-item>

        <el-form-item label="A/B测试">
          <DimensionValueSelect
            v-model="filterForm.abTest"
            :project-alias="projectAlias"
            :table="TableNames.USER"
            field="u_ab_test"
            placeholder="请选择或搜索A/B测试"
            multiple
            clearable
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleQuery" size="default">查询</el-button>
          <el-button @click="handleFilterConfig" size="default">过滤器</el-button>
        </el-form-item>
      </el-form>

      <!-- 已应用的过滤器 - 单独区域，支持换行 -->
      <div v-if="activeFilters.filters.length > 0" class="active-filters-bar">
        <div class="filter-controls">
          <el-icon class="control-icon"><Filter /></el-icon>
          <div class="filters-compact">
            <div 
              v-for="(filter, index) in activeFilters.filters" 
              :key="index" 
              class="filter-chip-compact"
            >
              <span class="filter-field-compact">{{ getPropertyName(filter.column.field) }}</span>
              <span class="filter-operator-compact">{{ operatorLabelMap[filter.operator] || filter.operator }}</span>
              <div v-if="needsValue(filter.operator)" class="filter-value-compact">
                <DimensionValueSelect
                  :model-value="filter.value"
                  :project-alias="projectAlias"
                  :table="(filter.column.table as any)"
                  :field="filter.column.field"
                  :operator="filter.operator"
                  :placeholder="getValuePlaceholder(filter.operator)"
                  class="filter-value-select-compact"
                  @update:modelValue="(value: any) => updateFilterValue(filter, value)"
                />
              </div>
              <span v-else class="filter-no-value-compact">-</span>
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <el-card class="table-card" shadow="never">
      <div class="table-header">
        <div class="header-left">
          <span class="result-count">共找到 {{ total.toLocaleString('zh-CN') }} 个用户</span>
        </div>
        <div class="header-actions">
          <el-button @click="handleColumnSelector" size="small" :icon="Filter">选择列</el-button>
        </div>
      </div>
      <div class="table-wrapper">
        <el-table
          v-loading="loading"
          :data="tableData"
          stripe
          border
          class="fill-table"
          style="width: 100%"
        >
          <!-- 统一按 columns 顺序渲染所有列 -->
          <el-table-column
            v-for="col in columns"
            :key="col.prop"
            :prop="col.prop"
            :label="col.label"
            :min-width="col.prop === 'u_event_time' ? 160 : 120"
            :fixed="col.prop === 'u_event_time' ? 'left' : false"
            :width="col.prop === 'u_openid' ? 180 : undefined"
          >
            <template #default="{ row }">
              <template v-if="col.prop === 'u_event_time'">
                {{ formatDateTime(row.u_event_time) }}
              </template>
              <template v-else-if="col.prop === 'u_openid'">
                <div 
                  class="user-id-cell copyable-cell" 
                  :title="`${row.u_openid}（双击复制）`"
                  @dblclick="handleDoubleClickCopy(row.u_openid, '用户ID')"
                >
                  {{ row.u_openid }}
                </div>
              </template>
              <template v-else>
                <span>{{ row[col.prop] ?? '-' }}</span>
              </template>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <div class="table-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="pageNum"
          :page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          size="default"
          background
          @size-change="onPageSizeChange"
          @current-change="onPageChange"
        />
      </div>
    </el-card>

    <!-- 过滤器配置弹窗 -->
    <el-dialog v-model="filterDialogVisible" title="过滤器配置" width="70%" destroy-on-close>
      <GlobalFilterComponent
        v-model="tempFilters"
        :scope="1"
        :project-alias="projectAlias"
        filter-type="user"
        :user-properties="allUserProperties"
      />
      <template #footer>
        <el-button @click="filterDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleFilterSave">确定</el-button>
      </template>
    </el-dialog>

    <!-- 列选择器弹窗 -->
    <el-dialog v-model="columnSelectorVisible" title="选择显示的列" width="700px">
      <div class="column-selector-body">
        <div class="available-columns">
          <h3>可选列</h3>
          <el-input v-model="searchKeyword" placeholder="搜索字段" clearable />
          <div class="column-list">
            <el-checkbox-group v-model="tempSelectedColumns">
              <el-checkbox
                v-for="prop in filteredProperties"
                :key="prop.propertyId"
                :value="prop.propertyId"
                class="column-checkbox"
              >
                {{ prop.propertyName }} ({{ prop.propertyId }})
              </el-checkbox>
            </el-checkbox-group>
          </div>
        </div>
        <div class="selected-columns">
          <h3>显示列 ({{ tempSelectedColumns.length }})</h3>
          <div class="selected-list">
            <div
              v-for="columnId in tempSelectedColumns"
              :key="columnId"
              class="selected-item"
            >
              <span>{{ getPropertyName(columnId) }}</span>
              <el-button
                link
                type="danger"
                :icon="Close"
                @click="removeSelectedColumn(columnId)"
              />
            </div>
          </div>
          <div v-if="tempSelectedColumns.length === 0" class="selected-empty">
            请选择左侧字段以展示
          </div>
        </div>
      </div>
      <template #footer>
        <el-button @click="columnSelectorVisible = false">取消</el-button>
        <el-button type="primary" @click="handleColumnSave">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { getUserList, getUserProperties } from '@/api/analytics'
import { OperatorTypes as OP, TableNames } from '@/constants/analysis'
import { useAppStore } from '@/store/app'
import type { UserListReq } from '@/types/doris/analysis'
import DimensionValueSelect from '@/components/common/DimensionValueSelect.vue'
import GlobalFilterComponent from '@/components/common/GlobalFilterComponent.vue'
import { Close, Filter } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, onMounted, onUnmounted, reactive, ref, watch } from 'vue'

const operatorLabelMap: Record<number, string> = {
  1: '等于',
  2: '不等于',
  3: '小于',
  4: '小于等于',
  5: '大于',
  6: '大于等于',
  7: '为空',
  8: '不为空',
  9: '在...之间',
  10: '包含',
  11: '不包含',
  12: '相似',
  13: '正则匹配',
  14: '不相似',
  17: '以...开头',
  18: '以...结尾',
}

const filterForm = reactive({
  timeRange: [] as [Date, Date] | [],
  userId: '' as string | string[],
  abTest: [] as string[]
})

const tableData = ref([])
const loading = ref(false)
const columns = ref([
  { prop: 'u_event_time', label: '注册时间' },
  { prop: 'u_openid', label: '用户ID' },
  { prop: 'u_ab_test', label: 'A/B测试' },
  { prop: 'u_gid', label: '游戏ID' },
  { prop: 'u_channel', label: '渠道' },
  { prop: 'u_from_ad_id', label: '来源广告' },
  { prop: 'u_from_advertiser_id', label: '来源广告主' },
  { prop: 'u_reg_cfg_version', label: '注册配置版本' },
  { prop: 'u_reg_game_version', label: '注册游戏版本' },
])
const filterDialogVisible = ref(false)
const columnSelectorVisible = ref(false)

const appStore = useAppStore()
const projectAlias = computed(() => appStore.activeProjectAlias)

const initialDefaultColumns = ref<string[]>(columns.value.map(c => c.prop))
const userDefinedColumns = ref(true)

const lastServerColumns = ref<string[]>([])
const tempSelectedColumns = ref<string[]>([])

function arraysEqual(a: string[] = [], b: string[] = []) {
  if (a === b) return true
  if (!Array.isArray(a) || !Array.isArray(b)) return false
  if (a.length !== b.length) return false
  for (let i = 0; i < a.length; i++) {
    if (a[i] !== b[i]) return false
  }
  return true
}

watch(tempSelectedColumns, (keys) => {
  const currentItemKeys = tempSelectedColumns.value
  if (arraysEqual(keys, currentItemKeys)) return

  const map = new Map(tempSelectedColumns.value.map(k => [k, k]))
  tempSelectedColumns.value = (keys || []).map(k => map.get(k) || k)
})

watch(tempSelectedColumns, (keys) => {
  const nextKeys = (keys || []).map(k => k)
  if (arraysEqual(nextKeys, tempSelectedColumns.value)) return
  tempSelectedColumns.value = nextKeys
})

const pageSize = ref(10)
const pageNum = ref(1)
const total = ref(0)

const activeFilters = ref<any>({ scope: 1, filters: [] })
const tempFilters = ref<any[]>([])

const allUserProperties = ref<any[]>([])
const searchKeyword = ref('')

const STORAGE_KEY = 'userList_filters_v2'

const saveFiltersToStorage = () => {
  try {
    const filtersData = {
      scope: activeFilters.value.scope,
      filters: activeFilters.value.filters,
      timestamp: Date.now()
    }
    localStorage.setItem(STORAGE_KEY, JSON.stringify(filtersData))
  } catch (error) {
    console.warn('保存过滤器失败:', error)
  }
}

const loadFiltersFromStorage = () => {
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) {
      const filtersData = JSON.parse(stored)
      if (filtersData && Array.isArray(filtersData.filters)) {
        activeFilters.value = {
          scope: filtersData.scope || 1,
          filters: filtersData.filters
        }
        return true
      }
    }
  } catch (error) {
    console.warn('加载过滤器失败:', error)
  }
  return false
}

const filteredProperties = computed(() => {
  if (!searchKeyword.value) {
    return allUserProperties.value
  }
  const lowerKeyword = searchKeyword.value.toLowerCase()
  return allUserProperties.value.filter(prop =>
    prop.propertyName.toLowerCase().includes(lowerKeyword) ||
    prop.propertyId.toLowerCase().includes(lowerKeyword)
  )
})

const fieldLabelMap: Record<string, string> = {
  'u_event_time': '注册时间',
  'u_openid': '用户ID',
  'u_ab_test': 'A/B测试',
  'u_gid': '游戏ID',
  'u_channel': '渠道',
  'u_from_ad_id': '来源广告',
  'u_from_advertiser_id': '来源广告主',
  'u_reg_cfg_version': '注册配置版本',
  'u_reg_game_version': '注册游戏版本'
}

const getPropertyName = (propId: string) => {
  const prop = allUserProperties.value.find(p => p.propertyId === propId)
  if (prop && prop.propertyName) return prop.propertyName
  if (fieldLabelMap[propId]) return fieldLabelMap[propId]
  return propId
}

const formatDateTime = (dateTime: string) => {
  if (!dateTime) return '-'
  try {
    const date = new Date(dateTime)
    if (isNaN(date.getTime())) return dateTime
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    const seconds = String(date.getSeconds()).padStart(2, '0')
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
  } catch {
    return dateTime
  }
}

const copyToClipboard = async (text: string, fieldName: string = '') => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success(`${fieldName || '数据'}已复制`)
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

const handleDoubleClickCopy = (text: string, fieldName: string) => {
  if (text && text !== '-') {
    copyToClipboard(text, fieldName)
  }
}

const setDefaultTimeRange = () => {
  const end = new Date()
  const start = new Date()
  start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
  filterForm.timeRange = [start, end]
}
setDefaultTimeRange()

const handleQuery = async () => {
  loading.value = true
  try {
    const final_filters: any[] = []

    if (filterForm.timeRange && filterForm.timeRange.length === 2) {
      final_filters.push({
        column: { table: TableNames.USER, field: 'u_event_time' },
        operator: 9, // Between
        value: [filterForm.timeRange[0].toISOString(), filterForm.timeRange[1].toISOString()],
      })
    }

    if (filterForm.userId && (Array.isArray(filterForm.userId) ? filterForm.userId.length > 0 : true)) {
      final_filters.push({
        column: { table: TableNames.USER, field: 'u_openid' },
        operator: 1, // 等于
        value: filterForm.userId,
      })
    }

    if (filterForm.abTest && filterForm.abTest.length > 0) {
      final_filters.push({
        column: { table: TableNames.USER, field: 'u_ab_test' },
        operator: 10, // 包含
        value: filterForm.abTest,
      })
    }

    if (activeFilters.value && activeFilters.value.filters.length > 0) {
      const validFilters = activeFilters.value.filters.filter((filter: any) => {
        if (!needsValue(filter.operator)) return true
        const value = filter.value
        if (value === null || value === undefined || value === '') return false
        if (Array.isArray(value)) {
          return value.length > 0 && value.some((v: any) => v !== null && v !== undefined && v !== '')
        }
        return true
      })
      if (validFilters.length > 0) final_filters.push(...validFilters)
    }

    const payload: UserListReq = {
      project_alias: projectAlias.value,
      columns: columns.value.map(c => ({ table: TableNames.USER, field: c.prop, alias: c.prop })),
      page_size: pageSize.value,
      page_num: pageNum.value,
      global_filter_groups: {
        query_dates: [],
        dashboard_form_filters: { scope: 1, filters: [], tag_filters: [], user_group_filters: [] },
        global_filters: {
          scope: activeFilters.value.scope || 1,
          filters: final_filters,
          tag_filters: [],
          user_group_filters: []
        }
      }
    }

    const res: any = await getUserList(payload)
    const code = Number(res?.code)
    if (code !== 0) throw new Error(res?.message || '查询失败')

    const data = res?.data || {}
    const rows = Array.isArray(data?.rows) ? data.rows : []
    const columnsRaw = Array.isArray(data?.columns) ? data.columns : []
    const count = Number(data?.count || rows.length || 0)

    tableData.value = rows
    total.value = count

    if (!userDefinedColumns.value) {
      let serverColumns: string[] = []
      if (columnsRaw.length > 0) {
        serverColumns = columnsRaw.map((c: any) => typeof c === 'string' ? c : c?.name).filter(Boolean)
      } else if (rows.length > 0) {
        serverColumns = Object.keys(rows[0] || {})
      }
      if (serverColumns.length) {
        columns.value = serverColumns.slice(0, 12).map(k => ({ prop: k, label: getPropertyName(k) }))
        lastServerColumns.value = [...serverColumns]
      }
    }
  } catch (error) {
    console.error('查询用户列表失败:', error)
    tableData.value = []
  } finally {
    loading.value = false
  }
}

const handleFilterConfig = () => {
  tempFilters.value = JSON.parse(JSON.stringify(activeFilters.value.filters))
  filterDialogVisible.value = true
}

const handleFilterSave = () => {
  activeFilters.value = {
    scope: 1,
    filters: JSON.parse(JSON.stringify(tempFilters.value)),
  }
  saveFiltersToStorage()
  filterDialogVisible.value = false
  handleQuery()
}

function needsValue(operator: number): boolean {
  return operator !== OP.IsNull && operator !== OP.IsNotNull
}

function getValuePlaceholder(operator: number): string {
  switch (operator) {
    case 9: return '如: 10,100'
    case 10: return '选择多个值，满足任一即可'
    case 11: return '选择多个值，都不包含'
    case 12:
    case 14: return '如: %关键词%'
    case 13: return '如: ^[a-z]+$'
    default: return '请输入值'
  }
}

const removeSelectedColumn = (columnId: string) => {
  const index = tempSelectedColumns.value.indexOf(columnId)
  if (index > -1) tempSelectedColumns.value.splice(index, 1)
}

const updateFilterValue = (filter: any, value: any) => {
  const filterIndex = activeFilters.value.filters.findIndex((f: any) => f === filter)
  if (filterIndex !== -1) {
    const newFilters = [...activeFilters.value.filters]
    newFilters[filterIndex] = { ...filter, value }
    activeFilters.value = { scope: activeFilters.value.scope, filters: newFilters }
  }
}

const handleColumnSelector = () => {
  const current = columns.value.map(c => c.prop)
  const fallback = lastServerColumns.value.length > 0 ? lastServerColumns.value : initialDefaultColumns.value
  const base = current.length > 0 ? current : fallback
  tempSelectedColumns.value = Array.from(new Set(base))
  columnSelectorVisible.value = true
}

const handleColumnSave = () => {
  const selected = [...tempSelectedColumns.value]
  const defaultPart = initialDefaultColumns.value.filter(p => selected.includes(p))
  const extraPart = selected.filter(p => !initialDefaultColumns.value.includes(p))
  const finalOrder = [...defaultPart, ...extraPart]

  columns.value = finalOrder.map(prop => ({ prop, label: getPropertyName(prop) }))
  userDefinedColumns.value = true
  columnSelectorVisible.value = false
  handleQuery()
}

const onPageChange = (p: number) => {
  pageNum.value = p
  handleQuery()
}

const onPageSizeChange = (ps: number) => {
  pageSize.value = ps
  pageNum.value = 1
  handleQuery()
}

const fetchUserProperties = async () => {
  try {
    const res: any = await getUserProperties(projectAlias.value)
    const props = Array.isArray(res) ? res : Array.isArray(res.data) ? res.data : Array.isArray(res.data?.data) ? res.data.data : []
    allUserProperties.value = props
  } catch (error) {
    console.error('获取用户属性失败:', error)
    allUserProperties.value = []
  }
}

watch(projectAlias, (newVal, oldVal) => {
  if (newVal && newVal !== oldVal) {
    pageNum.value = 1
    fetchUserProperties()
    handleQuery()
  }
})

onMounted(() => {
  if (projectAlias.value) {
    fetchUserProperties()
    handleQuery()
  }
  
  loadFiltersFromStorage()
})

onUnmounted(() => {
})
</script>

<style lang="scss" scoped>
.user-list-container {
  padding: 12px;
  background: #f8f9fa;
  min-height: calc(100vh - 60px);
  
  .filter-card {
    margin-bottom: 12px;
    border-radius: 4px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    border: 1px solid #e5e7eb;
    :deep(.el-card__body) { padding: 12px; }
  }

  .table-card {
    border-radius: 4px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    border: 1px solid #e5e7eb;
    :deep(.el-card__body) {
      padding: 12px;
      display: flex;
      flex-direction: column;
      height: 100%;
    }
    
    .table-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;
      padding-bottom: 12px;
      border-bottom: 1px solid #e5e7eb;
      flex-shrink: 0;
      .header-left {
        flex: 1;
        .result-count { font-size: 14px; color: #6b7280; font-weight: 500; }
      }
      .header-actions {
        display: flex;
        align-items: center;
        gap: 8px;
      }
    }
    
    .table-wrapper {
      display: flex;
      flex-direction: column;
      flex: 1;
      min-height: 0;
      :deep(.el-table) { flex-shrink: 0; }
      :deep(.el-table__body-wrapper) { flex-shrink: 0; }
    }

    .table-pagination {
      display: flex;
      justify-content: flex-end;
      margin-top: 16px;
      padding: 16px 0;
      border-top: 1px solid #e5e7eb;
      background-color: #f8f9fa;
    }
  }
  
  .compact-form {
    :deep(.el-form-item) {
      margin-right: 16px;
      margin-bottom: 12px;
    }
    :deep(.el-form-item__label) {
      font-size: 14px;
      color: #374151;
      font-weight: 500;
    }
  }
  
  :deep(.el-button) {
    border-radius: 4px;
    font-weight: 400;
  }

  :deep(.el-table) {
    border-radius: 4px;
    overflow: hidden;
    border: 1px solid #e5e7eb;
    .el-table__empty-block { min-height: 60px; }
    .el-table__body-wrapper { min-height: auto !important; }
  }

  :deep(.el-table th) {
    background: #f8f9fa;
    color: #495057;
    font-weight: 500;
    border-bottom: 1px solid #dee2e6;
  }

  :deep(.el-table td) { border-bottom: 1px solid #dee2e6; }
  :deep(.el-table tr:hover td) { background-color: #f8f9fa; }
  :deep(.el-table--striped .el-table__body tr.el-table__row--striped td) {
    background-color: #ffffff;
  }
  
  .fill-table {
    flex: 1 1 auto;
    width: 100%;
    height: 100%;
  }
  .fill-table :deep(.el-table__body-wrapper) {
    height: auto !important;
    overflow-y: auto;
    flex: 1;
  }
  .fill-table :deep(.el-table__inner-wrapper) { height: 100% !important; }
  .fill-table :deep(.el-table__header), .fill-table :deep(.el-table__body) { width: 100% !important; }
  .fill-table :deep(.el-table__header-wrapper) { flex-shrink: 0; }
}

.active-filters-bar {
  margin-top: 12px;
  padding: 8px 12px;
  background: #f8f9fa;
  border-radius: 4px;
  border: 1px solid #e5e7eb;
}

.filter-controls {
  display: flex;
  align-items: center;
  gap: 12px;
}

.control-icon { color: #6c757d; font-size: 16px; }
.filters-compact { display: flex; gap: 10px; flex-wrap: wrap; }
.filter-chip-compact {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  background: #ffffff;
  border-radius: 4px;
  font-size: 13px;
  border: 1px solid #d1d5db;
}
.filter-chip-compact:hover { background: #f3f4f6; }
.filter-field-compact { color: #374151; font-weight: 500; }
.filter-operator-compact { color: #6b7280; }
.filter-value-compact { min-width: 80px; }
.filter-no-value-compact { color: #9ca3af; font-style: italic; }
.logic-tag-compact { font-size: 10px; height: 20px; line-height: 20px; }
.filter-value-select-compact :deep(.el-input) { font-size: 12px; }
.filter-value-select-compact :deep(.el-input__wrapper) { padding: 1px 8px; border-radius: 3px; }

:deep(.el-transfer) {
  display: flex;
  justify-content: space-between;
  align-items: center;
  .el-transfer-panel { width: 250px; }
}

.column-selector-body {
  display: flex;
  gap: 20px;
  .available-columns, .selected-columns {
    flex: 1;
    border: 1px solid #e5e7eb;
    border-radius: 4px;
    padding: 12px;
    height: 400px;
    background: #f8f9fa;
  }
  .available-columns {
    .el-input { margin-bottom: 10px; }
    .column-list { height: calc(400px - 52px); overflow-y: auto; }
  }
  .selected-columns {
    overflow-y: auto;
    .selected-empty { color: #909399; padding: 8px 0; }
  }
  .column-checkbox { display: block; margin-bottom: 5px; }
  .selected-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 5px;
    border-bottom: 1px solid #eee;
  }
  .selected-list { height: calc(400px - 72px); overflow-y: auto; }
}

.user-id-cell, .ad-id-cell {
  font-family: 'JetBrains Mono', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
  color: #495057;
  background: #f8f9fa;
  padding: 4px 8px;
  border-radius: 4px;
  font-weight: 400;
  border: 1px solid #dee2e6;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
  cursor: help;
}

.user-id-cell:hover, .ad-id-cell:hover { background: #e9ecef; border-color: #adb5bd; }
.copyable-cell { position: relative; user-select: text; transition: all 0.2s ease; }
.copyable-cell:hover::after {
  content: '📋';
  position: absolute;
  top: -2px;
  right: -2px;
  font-size: 10px;
  opacity: 0.6;
  pointer-events: none;
}
.copyable-cell:active { transform: scale(0.98); }
</style>

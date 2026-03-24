<template>
  <div class="user-list-content">
    <!-- 过滤条件回显 -->
    <el-card v-if="hasActiveFilters" class="filter-card" shadow="never">
      <div class="active-filters-bar">
        <div class="filter-controls">
          <el-icon class="control-icon"><Filter /></el-icon>
          <span class="filter-title">当前过滤条件：</span>
          <div class="filters-compact">
            <!-- 显示标签过滤器 -->
            <div v-if="initialTagFilter?.tagCode" class="filter-chip-compact tag-filter">
              <span class="filter-field-compact">标签</span>
              <span class="filter-operator-compact">等于</span>
              <span class="filter-value-compact tag-value">{{ initialTagFilter.tagShowName }}</span>
            </div>
            
            <!-- 显示其他过滤器 -->
            <div 
              v-for="(filter, index) in activeFilters.filters" 
              :key="index" 
              class="filter-chip-compact"
            >
              <span class="filter-field-compact">{{ getPropertyName(filter.column.field) }}</span>
              <span class="filter-operator-compact">{{ operatorLabelMap[filter.operator] || filter.operator }}</span>
              <div v-if="needsValue(filter.operator)" class="filter-value-compact">
                <span class="filter-value-text">{{ formatFilterValue(filter.value) }}</span>
              </div>
              <span v-else class="filter-no-value-compact">-</span>
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <el-card class="table-card" shadow="never">
      <div class="table-header">
        <div class="table-info">
          <span class="result-count">
            共找到 {{ total.toLocaleString('zh-CN') }} 个用户
          </span>
        </div>
        <div class="table-actions">
          <el-button @click="handleColumnSelector" size="small">选择列</el-button>
          <el-button @click="handleExport" size="small" :icon="Download">导出</el-button>
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
          :height="tableHeight"
        >
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
                <div 
                  class="cell-content copyable-cell"
                  :title="`${row[col.prop]}（双击复制）`"
                  @dblclick="handleDoubleClickCopy(row[col.prop], col.label)"
                >
                  {{ row[col.prop] }}
                </div>
              </template>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="table-pagination">
          <el-pagination
            v-model:current-page="pageNum"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            layout="sizes, prev, pager, next"
            :total="total"
            @size-change="onPageSizeChange"
            @current-change="onPageChange"
          />
        </div>
      </div>
    </el-card>


    <!-- 列选择对话框 -->
    <el-dialog
      v-model="columnSelectorVisible"
      title="选择显示列"
      width="60%"
      :close-on-click-modal="false"
    >
      <div class="column-selector">
        <div class="available-columns">
          <h3>可选字段</h3>
          <el-input
            v-model="searchKeyword"
            placeholder="搜索字段"
            clearable
            class="search-input"
          />
          <div class="columns-list">
            <el-checkbox-group v-model="tempSelectedColumns" class="checkbox-group">
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
import { getUserList } from '@/api/analytics'
import { OperatorTypes as OP, TableNames } from '@/constants/analysis'
import type { FilterGroup } from '@/types/doris/filter'
import type { UserListReq } from '@/types/doris/analysis'
import { Close, Download, Filter } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'

interface TagFilter {
  id: number
  tagCode: string
  tagShowName: string
  userCount?: number
  projectAlias?: string
  [key: string]: any
}

interface UserGroupFilterPayload {
  groupName?: string
  groupCode: string
  operator?: number
}

interface Props {
  initialTagFilter?: TagFilter | null
  initialUserGroupFilter?: UserGroupFilterPayload | null
  projectAlias?: string
}

interface Emits {
  (e: 'close'): void
  (e: 'totalChange', total: number): void
}

const props = withDefaults(defineProps<Props>(), {
  initialTagFilter: null,
  initialUserGroupFilter: null,
  projectAlias: ''
})

const emit = defineEmits<Emits>()

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

const loading = ref(false)
const tableData = ref([])
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

const pageSize = ref(20)
const pageNum = ref(1)
const total = ref(0)

const activeFilters = ref<any>({ scope: 1, filters: [] })
const columnSelectorVisible = ref(false)
const tempSelectedColumns = ref<string[]>([])
const userDefinedColumns = ref(false)
const lastServerColumns = ref<string[]>([])

const searchKeyword = ref('')

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

const tableHeight = ref<string | number>('auto')

const calculateTableHeight = () => {
  const reservedHeight = 300
  const windowHeight = window.innerHeight
  const calculatedHeight = windowHeight - reservedHeight
  tableHeight.value = Math.max(calculatedHeight, 400)
}

const filteredProperties = computed(() => {
  const defaultProps = Object.keys(fieldLabelMap).map(key => ({
    propertyId: key,
    propertyName: fieldLabelMap[key]
  }))
  
  if (!searchKeyword.value) {
    return defaultProps
  }
  
  const lowerKeyword = searchKeyword.value.toLowerCase()
  return defaultProps.filter(prop =>
    prop.propertyName.toLowerCase().includes(lowerKeyword) ||
    prop.propertyId.toLowerCase().includes(lowerKeyword)
  )
})

const hasActiveFilters = computed(() => {
  return !!(props.initialTagFilter?.tagCode || activeFilters.value.filters.length > 0)
})

const getPropertyName = (propId: string) => {
  if (fieldLabelMap[propId]) {
    return fieldLabelMap[propId]
  }
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

const handleDoubleClickCopy = async (text: string, fieldName: string) => {
  if (!text) return
  
  try {
    await navigator.clipboard.writeText(String(text))
    ElMessage.success(`${fieldName}已复制到剪贴板`)
  } catch (error) {
    console.error('复制失败:', error)
    ElMessage.error('复制失败')
  }
}

const needsValue = (operator: number): boolean => {
  return ![7, 8].includes(operator)
}

const formatFilterValue = (value: any): string => {
  if (value === null || value === undefined) return '-'
  if (Array.isArray(value)) {
    return value.join(', ')
  }
  return String(value)
}

const initializeTagFilter = () => {
}

const handleQuery = async () => {
  if (!props.projectAlias) {
    ElMessage.warning('项目别名不能为空')
    return
  }

  loading.value = true
  try {
    const final_filters = activeFilters.value.filters.filter((filter: any) => 
      filter.column?.table !== TableNames.TAG
    )

    const payload: UserListReq = {
      project_alias: props.projectAlias,
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

    if (props.initialTagFilter?.tagCode) {
      payload.global_filter_groups.global_filters.tag_filters.push({
        tag_code: props.initialTagFilter.tagCode,
        operator: OP.EqualTo as any, 
        tag_value: 1 
      })
    }

    if (props.initialUserGroupFilter?.groupCode) {
      payload.global_filter_groups.global_filters.user_group_filters.push({
        group_name: props.initialUserGroupFilter.groupName || '',
        group_code: props.initialUserGroupFilter.groupCode,
        operator: (props.initialUserGroupFilter.operator ?? OP.In) as any
      })
    }

    const res: any = await getUserList(payload)
    console.log('用户列表API响应:', res)

    const code = Number(res?.code)
    if (code !== 0) {
      throw new Error(res?.message || '查询失败')
    }

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
        serverColumns = Object.keys(rows[0])
      }
      if (serverColumns.length > 0) {
        columns.value = serverColumns.slice(0, 12).map(k => ({ prop: k, label: getPropertyName(k) }))
        lastServerColumns.value = [...serverColumns]
      }
    }
    
    emit('totalChange', total.value)
  } catch (error) {
    console.error('查询用户列表失败:', error)
    tableData.value = []
    total.value = 0
    emit('totalChange', 0)
    ElMessage.error('查询失败')
  } finally {
    loading.value = false
  }
}


const handleExport = () => {
  ElMessage.info('导出功能开发中...')
}

const handleColumnSelector = () => {
  tempSelectedColumns.value = columns.value.map(c => c.prop)
  columnSelectorVisible.value = true
}

const handleColumnSave = () => {
  columns.value = tempSelectedColumns.value.map(prop => ({
    prop,
    label: getPropertyName(prop),
  }))
  columnSelectorVisible.value = false
  handleQuery()
}

const removeSelectedColumn = (columnId: string) => {
  const index = tempSelectedColumns.value.indexOf(columnId)
  if (index > -1) {
    tempSelectedColumns.value.splice(index, 1)
  }
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


watch(() => props.initialTagFilter, () => {
  initializeTagFilter()
  handleQuery()
}, { immediate: true })

onMounted(() => {
  calculateTableHeight()
  window.addEventListener('resize', calculateTableHeight)
  
  if (props.initialTagFilter) {
    initializeTagFilter()
  }
  
  handleQuery()
})

onUnmounted(() => {
  window.removeEventListener('resize', calculateTableHeight)
})
</script>

<style lang="scss" scoped>
.user-list-content {
  height: 100%;
  padding: 16px;
  background: #f8f9fa;
  display: flex;
  flex-direction: column;
  
  .filter-card {
    margin-bottom: 12px;
    border-radius: 4px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    border: 1px solid #e5e7eb;
    
    :deep(.el-card__body) {
      padding: 12px;
    }
  }

  .table-card {
    flex: 1;
    border-radius: 4px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    border: 1px solid #e5e7eb;
    display: flex;
    flex-direction: column;
    
    :deep(.el-card__body) {
      padding: 12px;
      display: flex;
      flex-direction: column;
      height: 100%;
      flex: 1;
    }
    
    .table-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;
      flex-shrink: 0;
      
      .table-info {
        .result-count {
          font-size: 14px;
          color: #606266;
          font-weight: 500;
        }
      }
      
      .table-actions {
        display: flex;
        gap: 8px;
      }
    }
    
    .table-wrapper {
      display: flex;
      flex-direction: column;
      flex: 1;
      min-height: 0;
    }

    .table-pagination {
      display: flex;
      justify-content: center;
      margin-top: 8px;
      padding: 6px 0;
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
  
  .active-filters-bar {
    padding: 0;
    
    .filter-controls {
      display: flex;
      align-items: flex-start;
      gap: 8px;
      
      .control-icon {
        margin-top: 2px;
        color: #606266;
        font-size: 16px;
      }
      
      .filter-title {
        font-size: 14px;
        color: #606266;
        font-weight: 500;
        margin-top: 1px;
      }
      
      .filters-compact {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
        flex: 1;
        
        .filter-chip-compact {
          display: flex;
          align-items: center;
          gap: 6px;
          background: #f0f2f5;
          border: 1px solid #d9d9d9;
          border-radius: 4px;
          padding: 4px 8px;
          font-size: 12px;
          line-height: 1.2;
          
          &.tag-filter {
            background: #e6f7ff;
            border-color: #91d5ff;
            
            .tag-value {
              font-weight: 600;
              color: #1890ff;
            }
          }
          
          .filter-field-compact {
            font-weight: 500;
            color: #262626;
          }
          
          .filter-operator-compact {
            color: #595959;
          }
          
          .filter-value-compact {
            max-width: 120px;
            
            .filter-value-text {
              color: #262626;
              font-weight: 500;
            }
          }
          
          .filter-no-value-compact {
            color: #8c8c8c;
            font-style: italic;
          }
        }
      }
    }
  }
  
  .copyable-cell {
    cursor: pointer;
    transition: all 0.2s;
    
    &:hover {
      background-color: #f0f9ff;
      color: #1890ff;
    }
  }
  
  .user-id-cell {
    font-family: 'Courier New', monospace;
    font-weight: 500;
  }
  
  .column-selector {
    display: flex;
    gap: 20px;
    height: 400px;
    
    .available-columns,
    .selected-columns {
      flex: 1;
      display: flex;
      flex-direction: column;
      
      h3 {
        margin: 0 0 12px 0;
        font-size: 16px;
        color: #303133;
      }
      
      .search-input {
        margin-bottom: 12px;
      }
      
      .columns-list {
        flex: 1;
        overflow-y: auto;
        border: 1px solid #e5e7eb;
        border-radius: 4px;
        padding: 8px;
        
        .checkbox-group {
          display: flex;
          flex-direction: column;
          gap: 8px;
          
          .column-checkbox {
            margin: 0;
            
            :deep(.el-checkbox__label) {
              font-size: 14px;
              line-height: 1.4;
            }
          }
        }
      }
      
      .selected-list {
        flex: 1;
        overflow-y: auto;
        border: 1px solid #e5e7eb;
        border-radius: 4px;
        padding: 8px;
        
        .selected-item {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 4px 8px;
          margin-bottom: 4px;
          background: #f5f7fa;
          border-radius: 3px;
          
          span {
            font-size: 14px;
            color: #303133;
          }
        }
      }
      
      .selected-empty {
        flex: 1;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #909399;
        font-style: italic;
        border: 1px dashed #e5e7eb;
        border-radius: 4px;
      }
    }
  }
}
</style>

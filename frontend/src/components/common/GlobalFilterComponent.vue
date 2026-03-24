<template>
  <div class="global-filter-component">
    <!-- 筛选条件列表 -->
    <div class="gf-container" v-if="modelValue.length > 0">
      <!-- 条件关系头部 -->
      <div class="gf-header-compact" v-if="modelValue.length > 1">
        <div class="relation-info-compact">
          <span class="relation-label-compact">筛选关系：</span>
          <el-tag size="small" type="primary">且 (AND)</el-tag>
          <span class="relation-desc-compact">所有条件必须同时满足</span>
        </div>
      </div>
      
      <!-- 条件列表 -->
      <div class="gf-list">
        <div 
          v-for="(f, idx) in modelValue" 
          :key="idx" 
          class="gf-item-wrapper"
        >
          <!-- 条件卡片 -->
          <div class="gf-item-modern">
            <!-- 条件标题栏 -->
            <div class="gf-header-bar">
              <span class="gf-index">条件 {{ idx + 1 }}</span>
              <el-button 
                link 
                type="danger" 
                size="small" 
                @click="removeFilter(idx)"
                :icon="Close"
              >
                删除
              </el-button>
            </div>
            
            <!-- 条件配置区 -->
            <div class="gf-config-area">
              <div class="gf-config-row">
                <div class="gf-config-item">
                  <label class="gf-label">字段</label>
                  <PropertySelect
                    :model-value="f.column.field"
                    :table="f.column.table"
                    :project-alias="projectAlias"
                    :show-formula-popover="false"
                    :is-global-filter="true"
                    :user-props-only="filterType === 'user'"
                    :event-props-only="filterType === 'event'"
                    :user-properties="userProperties"
                    :event-properties="eventProperties"
                    placeholder="选择属性"
                    @update:model-value="(v:string)=> updateFilterField(idx, v)"
                    @update:table="(t:any)=> updateFilterTable(idx, t)"
                    @change="(payload:any)=> onPropertyChange(idx, payload)"
                  />
                </div>
                
                <div class="gf-config-item">
                  <label class="gf-label">操作符</label>
                  <el-select 
                    v-model="f.operator" 
                    placeholder="选择操作符" 
                    size="small"
                    style="width: 100%;"
                    @change="(op:number)=> { handleOperatorChange(idx, op); ensureOperatorValid(idx) }"
                  >
                    <el-option v-for="op in getOperatorOptions(f)" :key="op.value" :label="op.label" :value="op.value" />
                  </el-select>
                </div>
                
                <div class="gf-config-item" v-if="needsValue(f.operator)">
                  <label class="gf-label">值</label>
                  <DimensionValueSelect
                    v-model="f.value"
                    :project-alias="projectAlias"
                    :table="f.column.table"
                    :field="f.column.field"
                    :operator="f.operator"
                    :placeholder="getValuePlaceholder(f.operator)"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 空状态 -->
    <div v-else class="gf-empty-modern">
      <div class="empty-content">
        <el-icon class="empty-icon"><Filter /></el-icon>
        <p class="empty-text">暂无筛选条件</p>
        <p class="empty-desc">点击下方按钮添加第一个筛选条件</p>
      </div>
    </div>
    
    <!-- 底部操作按钮 -->
    <div class="gf-actions-modern">
      <el-button @click="addFilter" type="primary" plain>
        <el-icon><Plus /></el-icon>
        添加筛选条件
      </el-button>
      <el-button v-if="modelValue.length > 0" @click="clearAllFilters" type="danger" plain>
        <el-icon><Delete /></el-icon>
        清空所有条件
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import DimensionValueSelect from '@/components/common/DimensionValueSelect.vue'
import PropertySelect from '@/components/common/PropertySelect.vue'
import { OperatorTypes as OP, TableNames } from '@/constants/analysis'
import { Close, Delete, Filter, Plus } from '@element-plus/icons-vue'
import { ref, watch } from 'vue'

const operatorLabelMap: Record<number, string> = {
  [OP.EqualTo]: '等于',
  [OP.NotEqualTo]: '不等于',
  [OP.LessThan]: '小于',
  [OP.LessOrEqual]: '小于等于',
  [OP.GreaterThan]: '大于',
  [OP.GreaterOrEqual]: '大于等于',
  [OP.In]: '包含',
  [OP.NotIn]: '不包含',
  [OP.StartWith]: '以...开头',
  [OP.EndWith]: '以...结尾',
}

interface GlobalFilter {
  column: { 
    table: any
    field: string 
  }
  operator: number
  value?: string | string[]  // 支持单选和多选值
  _dataType?: string
}

interface Props {
  modelValue: GlobalFilter[]
  scope?: number
  projectAlias?: string
  filterType?: 'all' | 'user' | 'event'
  userProperties?: any[]
  eventProperties?: any[]
}

interface Emits {
  (e: 'update:modelValue', value: GlobalFilter[]): void
  (e: 'update:scope', value: number): void
}

const props = withDefaults(defineProps<Props>(), {
  scope: 1,
  projectAlias: '',
  filterType: 'all',
  userProperties: () => [],
  eventProperties: () => [],
})

const emit = defineEmits<Emits>()

// 内部状态
const filterScope = ref(props.scope)

// 监听scope变化
watch(() => props.scope, (newScope) => {
  filterScope.value = newScope
})

// 添加筛选条件
const addFilter = () => {
  const newFilters = [...props.modelValue]
  const newTable = props.filterType === 'user' ? TableNames.USER : TableNames.EVENT
  newFilters.push({
    column: { table: newTable, field: '' },
    operator: 1,
    value: ''
  })
  emit('update:modelValue', newFilters)
}

// 删除筛选条件
const removeFilter = (index: number) => {
  const newFilters = [...props.modelValue]
  newFilters.splice(index, 1)
  emit('update:modelValue', newFilters)
}

// 清空所有条件
const clearAllFilters = () => {
  emit('update:modelValue', [])
}

// 更新字段
const updateFilterField = (index: number, field: string) => {
  const newFilters = [...props.modelValue]
  newFilters[index].column.field = field
  emit('update:modelValue', newFilters)
}

// 更新表
const updateFilterTable = (index: number, table: any) => {
  const newFilters = [...props.modelValue]
  newFilters[index].column.table = table
  emit('update:modelValue', newFilters)
}

function getDataTypeForFilter(f: any): string {
  const dt = String(f?._dataType || '').toLowerCase()
  if (dt) return dt
  return 'string'
}

function getOperatorOptions(f: any) {
  const dt = getDataTypeForFilter(f)
  let values: number[]
  if (['int', 'integer', 'bigint', 'decimal', 'double', 'float', 'number', 'numeric'].includes(dt)) {
    values = [OP.EqualTo, OP.NotEqualTo, OP.GreaterThan, OP.GreaterOrEqual, OP.LessThan, OP.LessOrEqual]
  } else if (['datetime', 'timestamp', 'date'].includes(dt)) {
    values = [OP.EqualTo, OP.NotEqualTo, OP.GreaterThan, OP.GreaterOrEqual, OP.LessThan, OP.LessOrEqual]
  } else {
    values = [OP.EqualTo, OP.NotEqualTo, OP.In, OP.NotIn, OP.StartWith, OP.EndWith]
  }
  return values.map(v => ({ value: v, label: operatorLabelMap[v] }))
}

function ensureOperatorValid(index: number) {
  const newFilters = [...props.modelValue]
  const f = newFilters[index]
  const allowed = getOperatorOptions(f).map(o => o.value)
  if (!allowed.includes(Number(f.operator))) {
    f.operator = allowed[0]
    emit('update:modelValue', newFilters)
  }
}

function onPropertyChange(index: number, payload: { field: string; table: string; property?: any }) {
  try {
    const newFilters = [...props.modelValue]
    const f = newFilters[index]
    f.column.field = payload.field
    f.column.table = payload.table
    f._dataType = payload?.property?.data_type || ''
    emit('update:modelValue', newFilters)
    ensureOperatorValid(index)
  } catch {}
}

// 操作符变化
const handleOperatorChange = (index: number, operator: number) => {
  const newFilters = [...props.modelValue]
  const filter = newFilters[index]
  const oldOperator = filter.operator
  
  filter.operator = operator
  
  // 如果操作符不需要值，清空值
  if (!needsValue(operator)) {
    filter.value = undefined
  } else {
    // 处理单选/多选模式切换
    const isNewMultiple = operator === OP.In || operator === OP.NotIn
    const wasOldMultiple = oldOperator === OP.In || oldOperator === OP.NotIn
    
    if (isNewMultiple && !wasOldMultiple) {
      // 从单选切换到多选：将字符串转换为数组
      if (filter.value && typeof filter.value === 'string') {
        filter.value = [filter.value]
      } else if (!filter.value) {
        filter.value = []
      }
    } else if (!isNewMultiple && wasOldMultiple) {
      // 从多选切换到单选：将数组转换为字符串
      if (Array.isArray(filter.value) && filter.value.length > 0) {
        filter.value = filter.value[0]
      } else {
        filter.value = undefined
      }
    }
  }
  
  emit('update:modelValue', newFilters)
}


// 判断操作符是否需要值
const needsValue = (operator: number): boolean => {
  return ![OP.IsNull, OP.IsNotNull].includes(operator) // 为空、不为空不需要值
}

// 获取值的占位符
const getValuePlaceholder = (operator: number): string => {
  switch (operator) {
    case OP.Between: // 区间
      return '如: 10,100'
    case OP.In: // 包含
      return '选择多个值，满足任一即可'
    case OP.NotIn: // 不包含
      return '选择多个值，都不包含'
    case OP.Like: // 相似
    case OP.NotLike: // 不相似
      return '如: %关键词%'
    case 13: // 正则
      return '如: ^[a-z]+$'
    default:
      return '请输入值'
  }
}
</script>

<style scoped>
.global-filter-component {
  padding: 12px;
}

.gf-container {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.gf-header-compact {
  padding: 8px 12px;
  background: #f8f9fa;
  border-radius: 4px;
  border: 1px solid #e9ecef;
}

.relation-info-compact {
  display: flex;
  align-items: center;
  gap: 8px;
}

.relation-label-compact {
  font-size: 13px;
  font-weight: 500;
  color: #495057;
}

.relation-desc-compact {
  font-size: 11px;
  color: #6c757d;
}

.gf-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.gf-item-wrapper {
  position: relative;
  margin: 4px;
}

.gf-item-modern {
  background: #ffffff;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.04);
  overflow: hidden;
  transition: all 0.3s ease;
}

.gf-item-modern:hover {
  border-color: #409eff;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.15);
}

.gf-header-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: transparent;
  border-bottom: 1px solid #e4e7ed;
}

.gf-index {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.gf-config-area {
  padding: 16px;
}

.gf-config-row {
  display: grid;
  grid-template-columns: 1fr 120px 1fr;
  gap: 16px;
  align-items: end;
}

.gf-config-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.gf-label {
  font-size: 12px;
  font-weight: 500;
  color: #606266;
  margin: 0;
}

.gf-empty-modern {
  padding: 60px 20px;
  text-align: center;
}

.empty-content {
  max-width: 300px;
  margin: 0 auto;
}

.empty-icon {
  font-size: 48px;
  color: #d3d4d6;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 16px;
  font-weight: 500;
  color: #606266;
  margin: 0 0 8px 0;
}

.empty-desc {
  font-size: 14px;
  color: #909399;
  margin: 0 0 24px 0;
  line-height: 1.5;
}

.gf-actions-modern {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 20px;
  padding: 20px;
  background: transparent;
  border-top: 1px solid #e4e7ed;
}

@media (max-width: 768px) {
  .gf-config-row {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .gf-actions-modern {
    flex-direction: column;
    gap: 8px;
  }
  
  .gf-actions-modern .el-button {
    width: 100%;
  }
}
</style>

<template>
  <div class="gf-item">
    <div class="gf-remove" @click="$emit('remove')" title="删除条件"><el-icon><Close /></el-icon></div>
    <div class="gf-content">
      <div class="gf-simple-row">
        <div class="field-selector">
          <UnifiedPropertySelector
            v-model="f.column.field"
            :event-id="eventId"
            :placeholder="'选择字段'"
            @change="onFieldChange"
          />
        </div>
        <div class="formula-display" @click="showFormulaPanel = true" :class="{ 'active': showFormulaPanel }">
          <span class="formula-summary">{{ getFormulaSummary() }}</span>
          <el-icon class="expand-icon"><ArrowDown /></el-icon>
        </div>
      </div>

      <div v-if="showFormulaPanel" class="formula-selection-panel" @click.stop>
        <div class="panel-row">
          <div class="panel-field-display" v-if="f.column.field">
            <span class="field-label">{{ getFieldDisplayName() }}</span>
          </div>
          <div class="panel-operator-group">
            <el-select v-model="f.operator" size="small" placeholder="选择操作符" @change="ensureOperatorValid">
              <el-option v-for="op in operatorOptions" :key="op.value" :value="op.value" :label="op.label" />
            </el-select>
          </div>

          <div class="panel-value-group" v-if="!isUserGroup() && f.operator !== OP.IsNull && f.operator !== OP.IsNotNull">
            <DateDiffPanel
              v-if="f.operator === OP.DateDiff"
              v-model="f.value"
            />
            <el-input-number 
              v-else-if="f.operator === OP.NDayRegiste"
              v-model="f.value"
              :min="0"
              :step="1"
              :precision="0"
              :controls="false"
              size="small"
              placeholder="输入天数"
              style="width: 100%;"
            />
            <DimensionValueSelect
              v-else-if="f.column.table !== 'user_tags' && f.column.table !== TableNames.TAG_TABLE"
              v-model="f.value"
              :project-alias="projectAlias"
              :table="f.column.table as any"
              :field="f.column.field"
              :event-id="eventId"
              :operator="(f.operator as any)"
              :placeholder="getValuePlaceholder((f.operator as any))"
            />
            <DimensionValueSelect
              v-else
              v-model="f.value"
              :project-alias="projectAlias"
              :table="TableNames.TAG_TABLE"
              :field="f.column.field"
              :operator="(f.operator as any)"
              placeholder="选择标签值"
            />
          </div>

          <div class="panel-value-group" v-else>
            <span class="hint-text">无需选择值</span>
          </div>
        </div>

        <div class="panel-footer">
          <el-button size="small" @click="showFormulaPanel = false">取消</el-button>
          <el-button type="primary" size="small" @click="confirmFormula">确定</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import DateDiffPanel from '@/components/common/DateDiffPanel.vue'
import DimensionValueSelect from '@/components/common/DimensionValueSelect.vue'
import { OperatorTypes as OP, TableNames } from '@/constants/analysis'
import { OperatorLabels } from '@/types/doris/constants'
import { ArrowDown, Close } from '@element-plus/icons-vue'
import { computed, inject, onBeforeUnmount, onMounted, ref } from 'vue'
import { ANALYSIS_CONTEXT_KEY } from '../../context'
import UnifiedPropertySelector from './UnifiedPropertySelector.vue'

type SimpleFilter = { column: { table: string; field: string }, operator: number, value?: any, _dataType?: string, _tagType?: string }

const props = defineProps<{ 
  modelValue: SimpleFilter; 
  eventId?: string; 
  isGlobalFilter?: boolean;
  /** 控制属性选择范围：'user-only' 仅用户属性；'event-only' 仅事件属性；默认全部 */
  propertySelectMode?: 'user-only' | 'event-only' | 'all'
}>()

const emit = defineEmits<{ (e:'update:modelValue', v: SimpleFilter): void; (e:'remove'): void }>()

const context = inject(ANALYSIS_CONTEXT_KEY)
const projectAlias = computed(() => context?.state.projectAlias || '')

const showFormulaPanel = ref(false)

function handleClickOutside(event: Event) {
  const target = event.target as Element
  if (!(target instanceof Element)) return
  if (!target.closest('.gf-item') && showFormulaPanel.value) {
    showFormulaPanel.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})

const f = props.modelValue
const eventId = props.eventId
const isGlobalFilter = props.isGlobalFilter ?? false
const propertySelectMode = props.propertySelectMode ?? 'all'

const dataType = computed(() => String(f?._dataType || 'string').toLowerCase())

const operatorOptions = computed(() => {
  if (isUserGroup()) return [
    { value: OP.In, label: '属于分群' },
    { value: OP.NotIn, label: '不属于分群' },
  ]
  
  let ops: number[] = []
  const dt = dataType.value
  
  if (['int', 'integer', 'bigint', 'decimal', 'double', 'float', 'number', 'numeric'].includes(dt)) {
    ops = [OP.EqualTo, OP.NotEqualTo, OP.GreaterThan, OP.GreaterOrEqual, OP.LessThan, OP.LessOrEqual, OP.In]
  } else if (['datetime', 'timestamp', 'date'].includes(dt)) {
    ops = [OP.EqualTo, OP.NotEqualTo, OP.GreaterThan, OP.GreaterOrEqual, OP.LessThan, OP.LessOrEqual, OP.DateDiff, OP.NDayRegiste]
  } else if (dt === 'boolean' || dt === 'bool') {
    ops = [OP.IsTrue, OP.IsFalse]
  } else {
    ops = [OP.EqualTo, OP.NotEqualTo, OP.In, OP.NotIn, OP.Like, OP.NotLike, OP.StartWith, OP.EndWith]
  }
  
  return ops.map(v => ({ value: v, label: OperatorLabels[v] || `Op ${v}` }))
})

function ensureOperatorValid() {
  const options = operatorOptions.value
  if (!options.find(o => o.value === Number(f.operator))) {
    f.operator = options[0]?.value || OP.EqualTo
    notifyChange()
  }
}

function isUserGroup(): boolean {
  return String(f?._tagType || '') === 'user_group'
}

function onFieldChange(payload: { id: string; name: string; type: string; table: string; property?: any }) {
  f.column.field = payload.id
  f.column.table = payload.table
  ;(f as any)._dataType = payload?.property?.data_type || ''
  ;(f as any)._tagType = payload?.property?.tag_type || ''
  ensureOperatorValid()
  notifyChange()
}

function getValuePlaceholder(operator: number) {
  switch (operator) {
    case OP.In:
    case OP.NotIn:
      return '如: 关键字'
    case OP.StartWith:
      return '如: 前缀'
    case OP.EndWith:
      return '如: 后缀'
    case OP.NDayRegiste:
      return '输入天数'
    default:
      return '输入筛选值'
  }
}

function getValueDisplayText(): string {
  if ((f.operator === OP.IsNull || f.operator === OP.IsNotNull) || isUserGroup()) {
    return ''
  }
  if (f.operator === OP.DateDiff) {
    const value = f.value
    if (value?.type === 1) return '当天'
    if (value?.type === 2 && value?.values) return `${value.values[0]}-${value.values[1]}天`
    return ''
  }
  if (f.operator === OP.NDayRegiste) {
    const n = Number(f.value)
    return Number.isFinite(n) && n >= 0 ? `${n}天` : ''
  }
  if (Array.isArray(f.value)) return f.value.length > 0 ? f.value.join(', ') : ''
  return f.value || ''
}

function getFormulaSummary(): string {
  const fieldName = f.column.field || '选择字段'
  const operator = OperatorLabels[f.operator] || '选择操作符'
  const value = getValueDisplayText()
  if (!f.column.field) return '点击设置筛选条件'
  return value ? `${fieldName} ${operator} ${value}` : `${fieldName} ${operator}`
}

function confirmFormula() {
  showFormulaPanel.value = false
  notifyChange()
}

function getFieldDisplayName(): string {
  return f.column.field || '等于'
}

function notifyChange() {
  emit('update:modelValue', { ...(f as any) })
}
</script>

<style scoped>
.gf-item { 
  position: relative; 
  background: #ffffff; 
  border: none; 
  border-radius: 4px; 
  padding: 4px 6px;
  transition: all 0.2s ease; 
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
  margin-bottom: 0;
}
.gf-item:hover {
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.gf-remove { 
  position: absolute; 
  top: 4px; 
  right: 6px; 
  width: 18px; 
  height: 18px; 
  border-radius: 50%; 
  background: #ffffff; 
  border: 1px solid var(--el-border-color-light); 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  cursor: pointer; 
  color: #909399; 
  font-size: 11px; 
  transition: all 0.2s ease; 
  z-index: 10; 
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.12);
  opacity: 0;
  visibility: hidden;
}
.gf-item:hover .gf-remove { 
  opacity: 1; 
  visibility: visible;
}
.gf-remove:hover { 
  background: #fde2e2; 
  color: #f56c6c; 
  border-color: #f56c6c; 
  transform: scale(1.1);
}

.gf-simple-row { 
  display: flex; 
  align-items: center; 
  gap: 6px;
  width: 100%; 
  min-height: 26px;
  overflow: hidden; 
  padding: 0;
}

.field-selector { 
  flex: 0 0 auto; 
  min-width: 80px; 
  max-width: 120px;
}
.field-selector :deep(.el-select) {
  width: 100%;
}
.field-selector :deep(.el-input__wrapper) {
  padding: 1px 8px;
}
.field-selector :deep(.el-input__inner) {
  font-size: 13px;
  height: 26px;
  line-height: 26px;
}

.formula-display { 
  flex: 1; 
  min-width: 0;
  max-width: 180px;
  display: flex; 
  align-items: center; 
  justify-content: space-between; 
  padding: 4px 8px; 
  border: 1px solid #e4e7ed; 
  border-radius: 4px; 
  background: #ffffff; 
  cursor: pointer; 
  transition: all 0.2s ease; 
  min-height: 26px; 
  overflow: hidden;
}
.formula-display:hover { 
  border-color: #409eff; 
  background: #f0f9ff;
}
.formula-display.active { 
  border-color: #409eff; 
  background: #f0f9ff; 
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.08);
}
.formula-summary { 
  flex: 1; 
  font-size: 13px; 
  color: #606266; 
  overflow: hidden; 
  text-overflow: ellipsis; 
  white-space: nowrap;
}
.expand-icon { 
  color: #c0c4cc; 
  font-size: 12px; 
  transition: transform 0.3s; 
  flex-shrink: 0;
  margin-left: 4px;
}
.formula-display.active .expand-icon { 
  transform: rotate(180deg);
}

.formula-selection-panel { 
  margin-top: 8px; 
  padding: 12px; 
  border: 1px solid #e4e7ed; 
  border-radius: 6px; 
  background: #fafafa; 
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}
.panel-field-display { 
  flex: 0 0 auto; 
  padding: 4px 10px; 
  background: #f0f9ff; 
  border: 1px solid #b3d8ff; 
  border-radius: 4px; 
  display: flex; 
  align-items: center; 
  min-width: 80px;
  max-width: 120px;
}
.field-label { 
  font-size: 12px; 
  font-weight: 500; 
  color: #409eff; 
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.panel-row { 
  display: flex; 
  gap: 10px; 
  align-items: center; 
  margin-bottom: 12px;
}
.panel-operator-group, 
.panel-value-group { 
  flex: 1; 
  min-width: 0;
}
.panel-operator-group :deep(.el-select),
.panel-value-group :deep(.el-select),
.panel-value-group :deep(.el-input-number) {
  width: 100%;
}
.panel-footer { 
  display: flex; 
  justify-content: flex-end; 
  gap: 8px; 
  padding-top: 10px; 
  margin-top: 6px; 
  border-top: 1px solid #e4e7ed;
}
.hint-text { 
  color: #909399; 
  font-size: 12px;
  font-style: italic;
}
</style>

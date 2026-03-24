<template>
  <div class="metric-property-picker">
    <UnifiedPropertySelector
      v-model="internalField"
      :event-id="eventId"
      :placeholder="'总次数'"
      @change="handlePropertySelect"
    >
      <template #trigger>
        <el-input
          v-model="displayValue"
          readonly
          placeholder="总次数"
          size="small"
          class="property-input"
        >
          <template #suffix>
            <el-icon><ArrowDown /></el-icon>
          </template>
        </el-input>
      </template>
    </UnifiedPropertySelector>
  </div>
</template>

<script setup lang="ts">
import { MetricFormulaLabels, MetricFormulas } from '@/types/doris/constants';
import { ArrowDown } from '@element-plus/icons-vue';
import { computed, inject, ref } from 'vue';
import { ANALYSIS_CONTEXT_KEY } from '../../context';
import UnifiedPropertySelector from '../shared/UnifiedPropertySelector.vue';

const props = withDefaults(defineProps<{ 
  eventId: string;
  field?: string;
  table?: string;
  formula?: number;
  propertyName?: string;
  autoClose?: boolean;
}>(), {
  autoClose: true,
})
const emit = defineEmits<{ 
  (e: 'select', payload: { field: string; table: string; formula: number; propertyName: string; eventName?: string }): void 
}>()

const context = inject(ANALYSIS_CONTEXT_KEY)

const internalField = ref(props.field)

// 计算显示值
const displayValue = computed(() => {
  const field = props.field
  const formula = Number(props.formula === undefined ? MetricFormulas.Count : props.formula)

  if (field === '__TOTAL_TIMES__') return '总次数'
  if (field === '__TOTAL_USERS__') return '总用户数'
  
  const nameToShow = props.propertyName || field

  if (props.propertyName === '总次数' && formula === MetricFormulas.Count) return '总次数'
  if (props.propertyName === '总用户数' && formula === MetricFormulas.CountDistinctUserId) return '总用户数'

  if (nameToShow) {
    const formulaText = getFormulaName(formula)
    return formulaText ? `${nameToShow} · ${formulaText}` : nameToShow
  }
  
  return '总次数'
})

function getFormulaName(formula: number): string {
  return MetricFormulaLabels[Number(formula)] || ''
}

function handlePropertySelect(payload: { id: string; name: string; type: string; table: string }) {
  // Default to Count if formula is not set
  const formula = props.formula !== undefined ? props.formula : MetricFormulas.Count
  
  emit('select', {
    field: payload.id,
    table: payload.table,
    formula: formula,
    propertyName: payload.name
  })
}
</script>

<style scoped>
.metric-property-picker {
  position: relative;
  flex: 1;
  min-width: 0;
  max-width: 160px;
}

.property-input {
  width: 100%;
  min-width: 140px;
  max-width: 160px;
}

.property-input :deep(.el-input__wrapper) {
  padding: 1px 8px;
  border-color: transparent;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.06);
  transition: box-shadow 0.15s ease, background-color 0.15s ease;
  background-color: #f7f8fa;
  border-radius: 6px;
}

.property-input :deep(.el-input__inner) {
  cursor: pointer;
  font-size: 14px;
  height: 28px;
  line-height: 28px;
}

.property-input :deep(.el-input__wrapper:hover) {
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.10);
  background-color: #f5f7fa;
}
.property-input :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px var(--el-color-primary) inset, 0 2px 8px rgba(64, 158, 255, 0.20);
  background-color: #fff;
}
</style>

<template>
  <div class="filter-value-input">
    <!-- Boolean Type -->
    <el-select
      v-if="dataType === DataType.Boolean"
      v-model="innerValue"
      :disabled="disabled"
      size="small"
      style="width: 100%"
      @change="handleChange"
    >
      <el-option label="True" :value="true" />
      <el-option label="False" :value="false" />
    </el-select>

    <!-- DateTime Type: BETWEEN -->
    <el-date-picker
      v-else-if="dataType === DataType.DateTime && isRangeOperator"
      v-model="innerValue"
      type="datetimerange"
      value-format="YYYY-MM-DD HH:mm:ss"
      format="YYYY-MM-DD HH:mm:ss"
      range-separator="至"
      start-placeholder="开始时间"
      end-placeholder="结束时间"
      size="small"
      :disabled="disabled"
      style="width: 100%"
      @change="handleChange"
    />

    <!-- DateTime Type: Single -->
    <el-date-picker
      v-else-if="dataType === DataType.DateTime"
      v-model="innerValue"
      type="datetime"
      value-format="YYYY-MM-DD HH:mm:ss"
      format="YYYY-MM-DD HH:mm:ss"
      :placeholder="placeholder || '选择时间'"
      size="small"
      :disabled="disabled"
      style="width: 100%"
      @change="handleChange"
    />

    <!-- Number Type: BETWEEN -->
    <div v-else-if="dataType === DataType.Number && isRangeOperator" class="number-range">
      <el-input-number 
        v-model="numberRangeStart" 
        size="small" 
        :disabled="disabled"
        :controls="false"
        style="flex: 1"
        @change="handleNumberRangeChange"
      />
      <span class="separator">-</span>
      <el-input-number 
        v-model="numberRangeEnd" 
        size="small" 
        :disabled="disabled"
        :controls="false"
        style="flex: 1"
        @change="handleNumberRangeChange"
      />
    </div>

    <!-- Number Type: IN / NOT_IN -->
    <el-select
      v-else-if="dataType === DataType.Number && isMultipleOperator"
      v-model="innerValue"
      multiple
      filterable
      allow-create
      default-first-option
      :reserve-keyword="false"
      :placeholder="placeholder || '输入数字后回车'"
      size="small"
      :loading="loading"
      :disabled="disabled"
      style="width: 100%"
      @focus="handleFocus"
      @visible-change="handleVisibleChange"
      @change="handleNumberMultipleChange"
    >
      <el-option
        v-for="opt in allValues"
        :key="opt"
        :label="opt"
        :value="opt"
      />
    </el-select>

    <!-- Number Type: Single -->
    <el-input-number
      v-else-if="dataType === DataType.Number"
      v-model="innerValue"
      size="small"
      :disabled="disabled"
      :controls="false"
      :placeholder="placeholder || '输入数字'"
      style="width: 100%"
      @change="handleChange"
    />

    <!-- String Type: All Operators (Single/Multiple depends on Operator) -->
    <el-select
      v-else
      v-model="innerValue"
      :multiple="isMultipleOperator"
      filterable
      allow-create
      default-first-option
      :collapse-tags="isMultipleOperator"
      :collapse-tags-tooltip="isMultipleOperator"
      :max-collapse-tags="2"
      :loading="loading"
      :placeholder="placeholder || (isMultipleOperator ? '搜索或输入多个值' : '搜索或输入值')"
      size="small"
      :disabled="disabled"
      style="width: 100%"
      @focus="handleFocus"
      @visible-change="handleVisibleChange"
      @change="handleChange"
    >
      <el-option
        v-for="opt in allValues"
        :key="opt"
        :label="opt"
        :value="opt"
      />
    </el-select>
  </div>
</template>

<script setup lang="ts">
import { DataType, Operator } from '@/types/doris/common';
import { computed, ref, watch } from 'vue';
import { useDimensionValues } from '../../composables/useDimensionValues';

const props = withDefaults(defineProps<{
  modelValue: any
  dataType: DataType | string
  operator: Operator | number
  projectAlias?: string
  tableName?: string
  fieldName?: string
  eventId?: string
  placeholder?: string
  disabled?: boolean
}>(), {
  disabled: false
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: any): void
}>()

// --- Computed Operators ---
const isMultipleOperator = computed(() => {
  return props.operator === Operator.In || props.operator === Operator.NotIn
})

const isRangeOperator = computed(() => {
  return props.operator === Operator.Between
})

// --- Dimension Values Loading ---
const { allValues, loading, loadValues } = useDimensionValues({
  projectAlias: computed(() => props.projectAlias),
  tableName: computed(() => props.tableName),
  fieldName: computed(() => props.fieldName),
  eventId: computed(() => props.eventId)
})

const handleFocus = () => {
  // 确保输入框聚焦时可以拉取数据 (支持 String 以及 Number 下拉多选等情况)
  loadValues(false)
}

const handleVisibleChange = (visible: boolean) => {
  // 当下拉框展开时，如果没有加载过数据，则拉取数据
  if (visible) {
    loadValues(false)
  }
}

// --- Value Normalization ---
// A pure function to safely convert values when operators or types change
const normalizeValue = (val: any, op: Operator | number, type: DataType | string) => {
  const isMulti = op === Operator.In || op === Operator.NotIn
  const isRange = op === Operator.Between

  if (isRange) {
    if (Array.isArray(val) && val.length === 2) return val
    if (Array.isArray(val) && val.length > 0) return [val[0], val[0]]
    if (val !== null && val !== undefined) return [val, val]
    return type === DataType.Number ? [0, 0] : []
  }

  if (isMulti) {
    if (Array.isArray(val)) return val
    if (val !== null && val !== undefined) return [val]
    return []
  }

  // Single value
  if (Array.isArray(val)) return val.length > 0 ? val[0] : undefined
  return val
}

// --- Internal State ---
const innerValue = ref<any>(normalizeValue(props.modelValue, props.operator, props.dataType))

// Special handling for number range since it uses two inputs
const numberRangeStart = ref<number>(0)
const numberRangeEnd = ref<number>(0)

const syncNumberRangeFromInner = () => {
  if (props.dataType === DataType.Number && isRangeOperator.value) {
    const arr = Array.isArray(innerValue.value) ? innerValue.value : [0, 0]
    numberRangeStart.value = Number(arr[0]) || 0
    numberRangeEnd.value = Number(arr[1]) || 0
  }
}

watch(innerValue, syncNumberRangeFromInner, { immediate: true })

const handleNumberRangeChange = () => {
  innerValue.value = [numberRangeStart.value, numberRangeEnd.value]
  handleChange()
}

// Ensure string numbers are converted to actual numbers in multiple select
const handleNumberMultipleChange = (val: any[]) => {
  if (Array.isArray(val)) {
    innerValue.value = val.map(v => Number(v)).filter(n => !isNaN(n))
  }
  handleChange()
}

const handleChange = () => {
  emit('update:modelValue', innerValue.value)
}

// --- Watchers ---
// Sync from parent
watch(() => props.modelValue, (newVal) => {
  // Only update if fundamentally different to avoid circular updates
  if (JSON.stringify(newVal) !== JSON.stringify(innerValue.value)) {
    innerValue.value = normalizeValue(newVal, props.operator, props.dataType)
  }
}, { deep: true })

// Re-normalize when operator or datatype changes
watch(() => [props.operator, props.dataType], ([newOp, newType]) => {
  innerValue.value = normalizeValue(innerValue.value, newOp as Operator, newType as DataType)
  handleChange()
})

</script>

<style scoped lang="scss">
.filter-value-input {
  width: 100%;
  
  .number-range {
    display: flex;
    align-items: center;
    gap: 8px;
    
    .separator {
      color: var(--el-text-color-secondary);
    }
  }
}
</style>

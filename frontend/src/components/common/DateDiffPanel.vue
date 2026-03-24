<template>
  <div class="date-diff-panel-compact">
    <el-radio-group v-model="selectedType" size="small" @change="onTypeChange">
      <el-radio-button :value="1">当天</el-radio-button>
      <el-radio-button :value="2">区间</el-radio-button>
    </el-radio-group>
    
    <div v-if="selectedType === 2" class="range-group-compact">
      <el-input-number
        v-model="rangeValues[0]"
        size="small"
        :min="-9999"
        :max="9999"
        placeholder="最小"
        controls-position="right"
        @change="onRangeChange"
      />
      <span class="range-separator">至</span>
      <el-input-number
        v-model="rangeValues[1]"
        size="small"
        :min="rangeValues[0] ?? 0"
        :max="9999"
        placeholder="最大"
        controls-position="right"
        @change="onRangeChange"
      />
      <span class="unit">天</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';

interface DateDiffValue {
  type: number
  values?: number[]
}

const props = defineProps<{
  modelValue?: DateDiffValue | null
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: DateDiffValue): void
}>()

// Use undefined for values to handle empty inputs gracefully
const selectedType = ref<number>(1)
const rangeValues = ref<(number | undefined)[]>([0, 7])

watch(() => props.modelValue, (newVal) => {
  if (newVal && typeof newVal === 'object' && 'type' in newVal) {
    selectedType.value = newVal.type
    if (newVal.type === 2 && newVal.values) {
      rangeValues.value = [...newVal.values]
    }
  }
}, { deep: true, immediate: true })

onMounted(() => {
  // Initialize with a default value if the provided v-model is not valid.
  if (!props.modelValue || typeof props.modelValue.type !== 'number') {
    emitValue()
  }
})

function onTypeChange() {
  emitValue()
}

function onRangeChange() {
  if (selectedType.value === 2) {
    const [min, max] = rangeValues.value
    if (typeof min === 'number' && typeof max === 'number' && min > max) {
      rangeValues.value[1] = min
    }
    emitValue()
  }
}

function emitValue() {
  const value: DateDiffValue = {
    type: selectedType.value,
  }
  if (selectedType.value === 2) {
    // Ensure emitted values are numbers, defaulting to 0 if undefined.
    value.values = [rangeValues.value[0] ?? 0, rangeValues.value[1] ?? 0]
  }
  emit('update:modelValue', value)
}
</script>

<style scoped>
.date-diff-panel-compact {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
}

.range-group-compact {
  display: flex;
  align-items: center;
  gap: 6px;
}

.range-separator,
.unit {
  font-size: 12px;
  color: var(--el-text-color-regular);
  white-space: nowrap;
  flex-shrink: 0;
}

:deep(.el-input-number) {
  width: 75px;
}

:deep(.el-input-number .el-input__inner) {
  text-align: center;
}
</style>

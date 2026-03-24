<template>
  <el-select
    v-model="innerValue"
    placeholder="操作符"
    size="small"
    style="width: 100%"
    @change="handleChange"
  >
    <el-option
      v-for="op in availableOperators"
      :key="op.value"
      :label="op.label"
      :value="op.value"
    />
  </el-select>
</template>

<script setup lang="ts">
import { Operator, DataType } from '@/types/doris/common'
import { computed, watch } from 'vue'
import { useFilterOperator } from '../../composables/useFilterOperator'

const props = defineProps<{
  modelValue: Operator | number
  dataType: DataType | string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: Operator): void
}>()

const dataTypeRef = computed(() => props.dataType as DataType)

const { availableOperators, validateOperator } = useFilterOperator(dataTypeRef)

const innerValue = computed({
  get: () => props.modelValue as Operator,
  set: (val) => emit('update:modelValue', val)
})

const handleChange = (val: Operator) => {
  emit('update:modelValue', val)
}

// Watch for datatype changes to auto-correct invalid operators
watch(() => props.dataType, () => {
  const validOp = validateOperator(props.modelValue)
  if (validOp !== props.modelValue) {
    emit('update:modelValue', validOp)
  }
}, { immediate: true })

</script>

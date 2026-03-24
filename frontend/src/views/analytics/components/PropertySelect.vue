<template>
  <el-select
    v-model="internalValue"
    filterable
    clearable
    placeholder="请选择属性"
    size="small"
    style="width: 100%"
  >
    <el-option
      v-for="item in properties"
      :key="item.id"
      :label="item.name"
      :value="item.id"
    >
      <div class="property-option">
        <span>{{ item.name }}</span>
        <span class="property-id">{{ item.id }}</span>
      </div>
    </el-option>
  </el-select>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Option } from '@/api/selector'

const props = defineProps<{
  modelValue: string
  properties: Option[]
}>()

const emit = defineEmits(['update:modelValue', 'change'])

const internalValue = computed({
  get: () => props.modelValue,
  set: (v) => {
    emit('update:modelValue', v)
    const selectedProp = props.properties.find(p => p.id === v)
    emit('change', selectedProp)
  }
})
</script>

<style scoped lang="scss">
.property-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  
  .property-id {
    font-size: 12px;
    color: var(--el-text-color-secondary);
    margin-left: 10px;
  }
}
</style>

<template>
  <div class="funnel-step-item" :class="{ 'is-active': m.filterPopoverVisible }">
    <div class="step-header">
      <div class="step-number">{{ index + 1 }}</div>
      <div class="step-main">
        <el-select 
          v-model="m.e_event_id" 
          placeholder="选择步骤事件" 
          size="default" 
          filterable
          class="event-select"
          @change="$emit('event-change')"
        >
          <el-option
            v-for="ev in projectEvents"
            :key="ev.id"
            :label="ev.name || ev.id"
            :value="ev.id"
          />
        </el-select>
        
        <div class="step-actions">
          <el-button 
            size="small" 
            :type="hasFilters ? 'primary' : 'default'" 
            plain 
            icon="Filter"
            @click="$emit('toggle-filter')"
          >
            {{ hasFilters ? '已设过滤' : '筛选' }}
          </el-button>
          <el-button 
            v-if="canRemove"
            size="small" 
            type="danger" 
            plain 
            icon="Delete"
            @click="$emit('remove')"
          />
        </div>
      </div>
    </div>

    <!-- 过滤器面板 -->
    <div v-if="m.filterPopoverVisible" class="step-filter-panel">
      <slot name="filters"></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Filter, Delete } from '@element-plus/icons-vue'
import { computed } from 'vue'

const props = defineProps<{
  index: number
  m: any
  projectEvents: any[]
  canRemove: boolean
}>()

defineEmits(['remove', 'toggle-filter', 'event-change'])

const hasFilters = computed(() => {
  return props.m.filters?.length > 0 || props.m.filter_group?.filters?.length > 0
})
</script>

<style scoped lang="scss">
.funnel-step-item {
  background: #ffffff;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  padding: 12px;
  margin-bottom: 12px;
  transition: all 0.2s;
  
  &:hover {
    border-color: var(--el-color-primary-light-5);
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
  }
  
  &.is-active {
    border-color: var(--el-color-primary);
  }
}

.step-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.step-number {
  width: 24px;
  height: 24px;
  background: var(--el-color-primary);
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: bold;
  flex-shrink: 0;
}

.step-main {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 12px;
}

.event-select {
  flex: 1;
}

.step-actions {
  display: flex;
  gap: 8px;
}

.step-filter-panel {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px dashed #e4e7ed;
}
</style>

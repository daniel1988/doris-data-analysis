<template>
  <el-drawer
    v-model="visible"
    title="我的报表"
    size="400px"
    destroy-on-close
  >
    <div v-loading="loading" class="report-drawer-content">
      <div v-if="reports.length === 0" class="empty-state">
        <el-empty description="暂无保存的报表" />
      </div>
      
      <div v-for="item in reports" :key="item.id" class="report-item" @click="handleLoad(item)">
        <div class="report-info">
          <div class="report-name">{{ item.name }}</div>
          <div class="report-desc" v-if="item.description">{{ item.description }}</div>
          <div class="report-meta">更新时间：{{ formatDate(item.update_time) }}</div>
        </div>
        <div class="report-actions">
          <el-button 
            type="danger" 
            link 
            icon="Delete" 
            @click.stop="handleDelete(item)"
          />
        </div>
      </div>
    </div>
  </el-drawer>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Report } from '@/api/report'

const props = defineProps<{
  modelValue: boolean
  reports: Report[]
  loading?: boolean
}>()

const emit = defineEmits(['update:modelValue', 'load', 'delete'])

const visible = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v)
})

const handleLoad = (report: Report) => {
  emit('load', report)
  visible.value = false
}

const handleDelete = (report: Report) => {
  emit('delete', report.id)
}

const formatDate = (dateStr?: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString()
}
</script>

<style scoped lang="scss">
.report-drawer-content {
  height: 100%;
  overflow-y: auto;
}

.report-item {
  padding: 15px;
  border-bottom: 1px solid var(--el-border-color-lighter);
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: background 0.2s;

  &:hover {
    background: var(--el-fill-color-light);
  }

  .report-info {
    flex: 1;
    min-width: 0;

    .report-name {
      font-weight: bold;
      font-size: 14px;
      margin-bottom: 4px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .report-desc {
      font-size: 12px;
      color: var(--el-text-color-secondary);
      margin-bottom: 4px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .report-meta {
      font-size: 11px;
      color: var(--el-text-color-placeholder);
    }
  }

  .report-actions {
    margin-left: 10px;
  }
}

.empty-state {
  margin-top: 50px;
}
</style>

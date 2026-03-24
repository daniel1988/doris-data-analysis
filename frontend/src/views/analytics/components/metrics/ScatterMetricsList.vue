<template>
  <el-card class="scatter-metrics-card" shadow="never">
    <template #header>
      <div class="card-header">
        <span class="title">分析指标</span>
        <el-tooltip content="分布分析仅支持单一指标">
          <el-icon class="info-icon"><InfoFilled /></el-icon>
        </el-tooltip>
      </div>
    </template>
    
    <div class="metrics-container">
      <MetricItem
        v-if="metric"
        :index="0"
        :m="metric"
        :hide-copy="true"
        :hide-formula="true"
        :hide-remove="true"
        :hide-drag-handle="true"
        @edit-name="onEditName"
        @finish-edit-name="onFinishEditName"
        @toggle-formula="onToggleFormula"
        @toggle-filter="onToggleFilter"
        @event-change="onEventChange"
      >
        <template #filters>
          <MetricFilterEditor :metric="metric" :project-alias="projectAlias" />
        </template>
        <template #property>
          <MetricPropertyPicker
            :event-id="metric.e_event_id"
            :field="metric.metric.column.field"
            :table="metric.metric.column.table"
            :formula="metric.metric.formula"
            :property-name="getDisplayPropertyName(metric)"
            @select="onApplyFormula"
          />
        </template>
      </MetricItem>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { InfoFilled } from '@element-plus/icons-vue'
import { computed, inject } from 'vue'
import { ANALYSIS_CONTEXT_KEY } from '../../context'
import MetricFilterEditor from './MetricFilterEditor.vue'
import MetricItem from './MetricItem.vue'
import MetricPropertyPicker from './MetricPropertyPicker.vue'

const context = inject(ANALYSIS_CONTEXT_KEY)
const metric = computed(() => (context?.state.form as any)?.metric)
const projectAlias = computed(() => context?.state.projectAlias || '')

defineProps<{
  getDisplayPropertyName: (m: any) => string
}>()

const emit = defineEmits<{
  (e: 'edit-name', m: any): void
  (e: 'finish-edit-name', m: any): void
  (e: 'toggle-formula', m: any): void
  (e: 'toggle-filter', payload: { m: any; index: number }): void
  (e: 'event-change', m: any): void
  (e: 'apply-formula', payload: { m: any; field: string; formula: number; table: string; propertyName?: string; eventName?: string }): void
}>()

function onEditName(m: any) { emit('edit-name', m) }
function onFinishEditName(m: any) { emit('finish-edit-name', m) }
function onToggleFormula(m: any) { emit('toggle-formula', m) }
function onToggleFilter(m: any) { emit('toggle-filter', { m, index: 0 }) }
function onEventChange() { emit('event-change', metric.value) }
function onApplyFormula(payload: { field: string; formula: number; table: string; propertyName?: string; eventName?: string }) {
  emit('apply-formula', { m: metric.value, ...payload })
}
</script>

<style scoped lang="scss">
.scatter-metrics-card {
  border: none;
  :deep(.el-card__header) {
    padding: 12px 16px;
    border-bottom: 1px solid var(--el-border-color-lighter);
  }
  :deep(.el-card__body) {
    padding: 12px;
  }
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  .title {
    font-size: 14px;
    font-weight: 600;
    color: var(--el-text-color-primary);
  }
  .info-icon {
    color: var(--el-text-color-secondary);
    font-size: 14px;
    cursor: help;
  }
}

.metrics-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>

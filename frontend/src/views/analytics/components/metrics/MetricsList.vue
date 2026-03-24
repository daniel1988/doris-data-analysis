<template>
  <MetricsCard @add-metric="onAddMetric()" @add-custom="onAddCustom()">
    <div class="metrics-container">
      <div
        v-for="(m, idx) in metrics"
        :key="m._uid"
        class="metric-draggable-wrapper"
        :class="{ 'drag-over': dragOverIndex === idx }"
        draggable="true"
        @dragstart="handleDragStart(idx)"
        @dragend="handleDragEnd"
        @dragover.prevent="handleDragOver(idx)"
        @drop.prevent="handleDrop(idx)"
      >
        <MetricItem
          :index="idx"
          :m="m"
          @remove="onRemoveMetric(idx)"
          @edit-name="onEditName(m)"
          @finish-edit-name="onFinishEditName(m)"
          @toggle-formula="onToggleFormula(m)"
          @copy="onCopyMetric(m)"
          @toggle-filter="onToggleFilter(m, idx)"
          @event-change="onEventChange(m)"
        >
          <template #filters>
            <MetricFilterEditor :metric="m" :project-alias="projectAlias" />
          </template>
          <template #property>
            <MetricPropertyPicker
              :event-id="m.e_event_id"
              :field="m.metric.column.field"
              :table="m.metric.column.table"
              :formula="m.metric.formula"
              :property-name="getDisplayPropertyName(m)"
              @select="(payload: any) => onApplyFormula(m, payload)"
            />
          </template>
          <template #custom-formula>
            <CustomFormulaPanel
              :m="m"
              @update-elements="(elements:any[]) => onUpdateCustomElements(m, elements)"
              @request-new-event="() => onRequestNewEvent(m)"
              @open-filter="(mm:any)=> onOpenCustomFilter(mm)"
            />
          </template>
        </MetricItem>
      </div>
      <div v-if="!metrics.length" class="empty-state">
        <span style="color: #999; margin-right: 12px;">暂无指标</span>
        <el-button size="small" type="primary" :icon="Plus" @click="onAddMetric()"></el-button>
      </div>
    </div>
  </MetricsCard>
</template>

<script setup lang="ts">
import { Plus } from '@element-plus/icons-vue'
import { ref, inject, computed } from 'vue'
import CustomFormulaPanel from './CustomFormulaPanel.vue'
import MetricFilterEditor from './MetricFilterEditor.vue'
import MetricItem from './MetricItem.vue'
import MetricPropertyPicker from './MetricPropertyPicker.vue'
import MetricsCard from './MetricsCard.vue'
import { ANALYSIS_CONTEXT_KEY } from '../../context'

const context = inject(ANALYSIS_CONTEXT_KEY)
const metrics = computed(() => (context?.state.form as any)?.event_metrics || [])
const projectAlias = computed(() => context?.state.projectAlias || '')

const props = defineProps<{
  getDisplayPropertyName: (m:any)=> string
  getFormulaName: (formula:number)=> string
}>()

const emit = defineEmits<{
  (e:'add-metric'): void
  (e:'add-custom'): void
  (e:'remove-metric', index:number): void
  (e:'edit-name', m:any): void
  (e:'finish-edit-name', m:any): void
  (e:'toggle-formula', m:any): void
  (e:'copy-metric', m:any): void
  (e:'toggle-filter', payload:{ m:any; index:number }): void
  (e:'event-change', m:any): void
  (e:'property-field-change', payload:{ m:any; value:string }): void
  (e:'apply-formula', payload:{ m:any; field:string; formula:number; table:string; propertyName?:string; eventName?:string }): void
  (e:'start-edit-property', m:any): void
  (e:'update-custom-elements', payload:{ m:any; elements:any[] }): void
  (e:'request-new-event', m:any): void
  (e:'open-custom-filter', m:any): void
  (e:'reorder', value:any[]): void
}>()

function onAddMetric() { emit('add-metric') }
function onAddCustom() { emit('add-custom') }
function onRemoveMetric(index:number) { emit('remove-metric', index) }
function onEditName(m:any) { emit('edit-name', m) }
function onFinishEditName(m:any) { emit('finish-edit-name', m) }
function onToggleFormula(m:any) { emit('toggle-formula', m) }
function onCopyMetric(m:any) { emit('copy-metric', m) }
function onToggleFilter(m:any, index:number) { emit('toggle-filter', { m, index }) }
function onEventChange(m:any) { emit('event-change', m) }
function onApplyFormula(m:any, payload:{field:string; formula:number; table:string; propertyName?:string; eventName?:string}) { emit('apply-formula', { m, ...payload }) }
function onUpdateCustomElements(m:any, elements:any[]) { emit('update-custom-elements', { m, elements }) }
function onRequestNewEvent(m:any) { emit('request-new-event', m) }
function onOpenCustomFilter(m:any) { emit('open-custom-filter', m) }

// drag & drop
const dragOverIndex = ref(-1)
let draggedIndex = -1

function handleDragStart(index: number) { draggedIndex = index }
function handleDragEnd() { draggedIndex = -1; dragOverIndex.value = -1 }
function handleDragOver(index: number) { dragOverIndex.value = index }
function handleDrop(index: number) {
  if (draggedIndex === -1 || draggedIndex === index) return
  const current = metrics.value
  if (!Array.isArray(current)) return
  const next = current.slice()
  const [moved] = next.splice(draggedIndex, 1)
  next.splice(index, 0, moved)
  emit('reorder', next)
  draggedIndex = -1
  dragOverIndex.value = -1
}
</script>

<style scoped>
.metrics-container { 
  display: flex; 
  flex-direction: column; 
  gap: 6px; 
  overflow-y: visible;
  overflow-x: hidden;
  padding: 2px;
}

.metrics-container::-webkit-scrollbar {
  width: 6px;
}

.metrics-container::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.metrics-container::-webkit-scrollbar-thumb {
  background: var(--el-color-primary-light-5);
  border-radius: 3px;
  transition: background 0.2s ease;
}

.metrics-container::-webkit-scrollbar-thumb:hover {
  background: var(--el-color-primary-light-3);
}
.empty-state { 
  padding: 16px; 
  text-align: center; 
  color: #999; 
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}
.metric-draggable-wrapper { 
  border-radius: 8px;
  transition: all 0.2s ease;
}
.metric-draggable-wrapper.drag-over { 
  outline: 2px dashed var(--el-color-success); 
  outline-offset: 2px;
}
</style>

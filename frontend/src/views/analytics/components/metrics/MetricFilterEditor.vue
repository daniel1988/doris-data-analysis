<template>
  <div class="metric-filter-popover bordered-filter-list">
    <!-- 有过滤器时显示列表 -->
    <div class="filter-content-popover" v-if="metric && Array.isArray(metric.filters) && metric.filters.length > 0">
      <div class="gf-header" v-if="metric.filters.length > 1">
        <div class="relation-info">
          <el-button size="small" class="relation-switch-btn" @click="toggleScope">{{ scopeText }}</el-button>
        </div>
      </div>
      <div class="gf-list">
        <div v-for="(f, fIdx) in metric.filters" :key="fIdx" class="gf-item-wrapper">
          <BaseFilterItem 
            v-model="metric.filters[fIdx]"
            :event-id="metric.e_event_id"
            :is-global-filter="false"
            :property-select-mode="propertySelectMode"
            @remove="removeFilter(fIdx)"
          />
        </div>
      </div>
      <!-- 底部添加按钮 -->
      <div class="filter-footer-actions mt-10">
        <el-button size="small" type="primary" link icon="Plus" @click="addFilter">添加筛选条件</el-button>
      </div>
    </div>
    
    <!-- 无过滤器时显示空状态 -->
    <div v-else class="filter-empty-state">
      <el-button size="small" type="primary" plain icon="Plus" @click="addFilter">添加第一个筛选条件</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue';
import { ANALYSIS_CONTEXT_KEY } from '../../context';
import BaseFilterItem from '../shared/BaseFilterItem.vue';

const props = defineProps<{ 
  metric: any; 
  /** 控制属性选择范围：'user-only' 仅用户属性；'event-only' 仅事件属性；默认全部 */
  propertySelectMode?: 'user-only' | 'event-only' | 'all'
}>()

const context = inject(ANALYSIS_CONTEXT_KEY)
const metric = computed(() => props.metric)
const projectAlias = computed(() => context?.state.projectAlias || '')
const propertySelectMode = computed(() => props.propertySelectMode || 'all')

const scopeText = computed(() => (metric.value?.scope === 1 ? '且' : '或'))

function toggleScope() {
  if (metric.value) {
    metric.value.scope = metric.value.scope === 1 ? 0 : 1
  }
}

function removeFilter(idx: number) {
  if (metric.value?.filters) {
    metric.value.filters.splice(idx, 1)
  }
}

function addFilter() {
  if (!metric.value.filters) {
    metric.value.filters = []
  }
  metric.value.filters.push({
    column: { table: 'event_data', field: '', alias: '' },
    operator: 1,
    value: ''
  })
}
</script>

<style scoped>
/* 过滤器列表容器 */
.bordered-filter-list {
  position: relative;
  padding-left: 10px; /* 为scope按钮留出空间 */
  margin-top: 8px; /* 减小上边距 */
}
.bordered-filter-list::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0; 
  bottom: 0;
  width: 2px;
  background-color: var(--el-border-color-light);
  border-radius: 1px;
}
.filter-content-popover { 
  padding: 0; 
  max-height: 400px; 
  overflow-y: auto;
}

.filter-empty-state {
  padding: 10px 0;
  display: flex;
  justify-content: center;
}

.filter-footer-actions {
  display: flex;
  justify-content: flex-start;
  padding-top: 8px;
  border-top: 1px dashed var(--el-border-color-lighter);
}

.mt-10 { margin-top: 10px; }

/* 关系按钮（且/或） */
.gf-header {
  position: absolute;
  top: 50%;
  left: -14px;
  transform: translateY(-50%);
  z-index: 10;
  background: transparent;
  padding: 0;
  box-shadow: none;
  width: 26px;
  height: 26px;
}
.relation-switch-btn {
  width: 26px !important;
  height: 26px !important;
  min-width: 26px;
  border-radius: 4px;
  border: 1px solid #d9d9d9;
  box-shadow: 0 1px 3px rgba(0,0,0,0.08);
  padding: 0 !important;
  font-weight: 600;
  font-size: 11px;
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  color: #606266;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
  overflow: visible;
  letter-spacing: 0.5px;
}
.relation-switch-btn:hover {
  border-color: #409eff;
  background: linear-gradient(135deg, #f0f9ff 0%, #e6f4ff 100%);
  color: #409eff;
  box-shadow: 0 2px 6px rgba(64, 158, 255, 0.15);
  transform: translateY(-1px);
}
.relation-info { 
  display: flex; 
  align-items: center; 
  gap: 8px; 
  font-size: 12px; 
  color: var(--el-text-color-regular);
}

/* 过滤项列表 */
.gf-item-wrapper { 
  display: flex; 
  flex-direction: column; 
  gap: 0; 
  margin: 0;
}
.gf-list { 
  margin-top: 0;
  display: flex;
  flex-direction: column;
  gap: 0;
}
</style>

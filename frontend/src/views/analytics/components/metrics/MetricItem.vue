<template>
  <div class="metric-item">
    <!-- 拖拽手柄 -->
    <div v-if="!hideDragHandle" class="drag-handle">
      <el-icon><DCaret /></el-icon>
      <el-icon><DCaret /></el-icon>
    </div>
    
    <!-- 序号标识 -->
    <div class="metric-number">{{ index + 1 }}</div>
    
    <!-- 删除按钮 -->
    <div v-if="!hideRemove" class="metric-close" @click="$emit('remove')">
      <el-icon><Close /></el-icon>
    </div>

    <!-- 第一行：指标名称 + 格式化下拉 + 操作按钮 -->
    <div class="metric-header-row">
      <!-- 指标名称 -->
      <div class="metric-name-section" :class="{ 'is-editing': m.isEditingName }" @click="!m.isEditingName && $emit('edit-name')">
        <div class="metric-name-display" v-if="!m.isEditingName">
          {{ m.name || '点击设置指标名称' }}
        </div>
        <el-input 
          v-else 
          v-model="m.name" 
          size="small" 
          placeholder="请输入指标名称" 
          class="metric-name-input" 
          @blur="$emit('finish-edit-name')" 
          @keyup.enter="$emit('finish-edit-name')" 
        />
      </div>

      <!-- 格式选择器 -->
      <el-select 
        :model-value="m.isCustom ? m.custom_metric?.format : m.metric?.format" 
        @update:model-value="(val: string) => {
          if (m.isCustom) {
            if (m.custom_metric) m.custom_metric.format = val;
          } else {
            if (m.metric) m.metric.format = val;
          }
        }"
        size="small" 
        class="format-select"
        title="数值格式化"
      >
        <el-option label="原始" value="raw" />
        <el-option label="整数" value="int" />
        <el-option label="小数(2位)" value="decimal" />
        <el-option label="%" value="percent" />
      </el-select>

      <!-- 操作按钮组 -->
      <div class="metric-actions">
        <!-- 过滤器按钮 -->
        <el-button 
          v-if="!m.isCustom" 
          size="small" 
          class="icon-button" 
          :class="{ 'is-active': hasFilters }" 
          @click="$emit('toggle-filter')" 
          title="指标过滤"
        >
          <el-icon><Filter /></el-icon>
        </el-button>

        <!-- 自定义公式按钮 -->
        <el-button 
          v-if="!hideFormula"
          size="small" 
          class="icon-button" 
          @click="$emit('toggle-formula')" 
          title="自定义公式"
        >
          Σ
        </el-button>

        <!-- 复制按钮 -->
        <el-button 
          v-if="!hideCopy"
          size="small" 
          class="icon-button" 
          @click="$emit('copy')" 
          title="复制指标配置"
        >
          <el-icon><DocumentCopy /></el-icon>
        </el-button>
      </div>
    </div>

    <!-- 第二行：事件选择 + 聚合方式 -->
    <div class="metric-config-row">
      <template v-if="!m.isCustom">
        <!-- 事件选择器 -->
        <el-select 
          v-model="m.e_event_id" 
          filterable 
          size="small" 
          placeholder="选择事件" 
          class="event-select" 
          popper-class="metric-event-select-popper"
          @visible-change="(v:boolean)=> v && refreshMetadata()" 
          @change="$emit('event-change')"
        >
          <el-option
            v-for="ev in projectEvents"
            :key="ev.id"
            :label="ev.name ? `${ev.name} ` : ev.id"
            :value="ev.id"
          >
            <div class="event-option">
              <span class="event-id">{{ ev.id }}</span>
              <span class="event-name" v-if="ev.name">{{ ev.name }}</span>
            </div>
          </el-option>
        </el-select>
        
        <!-- 连接词 -->
        <span class="connector">的</span>
        
        <!-- 聚合方式选择器 -->
        <slot name="property"></slot>
      </template>
      <template v-else>
        <div class="custom-formula-group">
          <slot name="custom-formula"></slot>
        </div>
      </template>
    </div>

    <!-- 过滤器展开区域 -->
    <div v-show="!m.isCustom && m.filterPopoverVisible" class="inline-filter-container">
      <slot name="filters"></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Close, DCaret, DocumentCopy, Filter } from '@element-plus/icons-vue';
import { computed, inject } from 'vue';
import { ANALYSIS_CONTEXT_KEY } from '../../context';

const props = withDefaults(defineProps<{ 
  index: number; 
  m: any;
  hideFormula?: boolean;
  hideCopy?: boolean;
  hideDragHandle?: boolean;
  hideRemove?: boolean;
}>(), {
  hideFormula: false,
  hideCopy: false,
  hideDragHandle: false,
  hideRemove: false
})
defineEmits<{ (e:'remove'):void; (e:'edit-name'):void; (e:'finish-edit-name'):void; (e:'toggle-formula'):void; (e:'copy'):void; (e:'toggle-filter'):void; (e:'event-change'):void }>()

const context = inject(ANALYSIS_CONTEXT_KEY)
const projectEvents = computed(() => context?.state.metadata.eventOptions || [])
const refreshMetadata = () => context?.actions.refreshMetadata()

const hasFilters = computed(() => !!(props.m?.filters?.length || props.m?.filter_group?.filters?.length || props.m?.filter_group?.filter_groups?.length))
</script>

<style scoped>
.metric-item { 
  position: relative; 
  display: flex;
  flex-direction: column;
  gap: 2px;
  border: none; 
  border-radius: 8px; 
  background: #ffffff; 
  min-width: 380px;
  box-sizing: border-box;
  margin: 2px 0; 
  padding: 10px 30px 10px 44px;
  transition: all 0.2s ease; 
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08); 
}
.metric-item:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12);
  border-color: var(--el-color-primary-light-5);
}

/* 拖拽手柄 */
.drag-handle {
  position: absolute;
  left: 4px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  flex-direction: column;
  gap: 0;
  color: #c0c4cc;
  cursor: grab;
  font-size: 10px;
  line-height: 1;
  opacity: 0;
  transition: all 0.2s ease;
}
.drag-handle :deep(.el-icon) {
  height: 8px;
  margin: -2px 0;
}
.metric-item:hover .drag-handle {
  opacity: 1;
}
.drag-handle:hover {
  color: var(--el-color-primary);
}
.drag-handle:active {
  cursor: grabbing;
}

/* 序号标识 */
.metric-number { 
  position: absolute; 
  left: 4px;
  top: 50%; 
  transform: translateY(-50%);
  width: 22px; 
  height: 22px; 
  border-radius: 4px; 
  background: var(--el-color-primary); 
  color: white; 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  font-size: 12px; 
  font-weight: 600; 
  z-index: 5; 
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.15);
}

/* 删除按钮 */
.metric-close { 
  position: absolute; 
  top: 50%;
  right: 4px; 
  transform: translateY(-50%);
  width: 20px; 
  height: 20px; 
  border-radius: 50%; 
  background: #ffffff; 
  border: 1px solid var(--el-border-color-light); 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  cursor: pointer; 
  color: #909399; 
  font-size: 12px; 
  transition: all 0.2s ease; 
  z-index: 10; 
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.12); 
  opacity: 0; 
  visibility: hidden;
}
.metric-item:hover .metric-close { 
  opacity: 1; 
  visibility: visible;
}
.metric-close:hover {
  background: #fde2e2;
  color: #f56c6c;
  border-color: #f56c6c;
  transform: translateY(-50%) scale(1.1);
}

/* 第一行：指标名称 + 格式选择 + 操作按钮 */
.metric-header-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
  padding-left: 2px;
  padding-right: 2px;
}
.metric-name-section {
  flex: 1;
  min-width: 0;
  max-width: 174px;
}
.metric-name-section.is-editing {
  max-width: 174px;
}
.metric-name-display { 
  font-size: 13px; 
  font-weight: 500; 
  color: var(--el-text-color-secondary); 
  padding: 3px 6px; 
  border: 1px solid transparent; 
  border-radius: 4px; 
  cursor: pointer; 
  transition: all 0.2s ease; 
  background: transparent;
  display: flex;
  width: 100%;
  min-height: 28px;
  align-items: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
}
.metric-name-display:hover {
  background: #f5f7fa;
  border-color: var(--el-border-color-lighter);
}
.metric-name-input { 
  width: 100%;
  max-width: 300px;
}

/* 指标名编辑态：下划线样式 */
.metric-name-input :deep(.el-input__wrapper) {
  background: transparent;
  border: none;
  border-radius: 0;
  box-shadow: none !important;
  padding: 0 0 2px 0;
  border-bottom: 1px solid var(--el-border-color);
  transition: border-color 0.15s ease;
}
.metric-name-input :deep(.el-input__wrapper:hover) {
  border-bottom-color: var(--el-text-color-placeholder);
}
.metric-name-input :deep(.el-input__wrapper.is-focus) {
  border-bottom-color: var(--el-color-primary);
  box-shadow: none;
}
.metric-name-input :deep(.el-input__inner) {
  padding: 0;
  height: 22px;
  line-height: 22px;
  font-size: 13px;
}

/* 第二行：事件选择 + 聚合方式 */
.metric-config-row { 
  display: flex; 
  align-items: center; 
  gap: 6px; 
  flex-wrap: wrap;
  min-width: 0; 
  width: 100%; 
  box-sizing: border-box;
  padding-right: 2px;
}

/* 连接词 */
.connector { 
  color: #606266; 
  font-size: 14px; 
  padding: 0 2px; 
  flex-shrink: 0;
}

/* 事件选择器 */
.event-select {
  min-width: 120px;
  max-width: 200px;
  flex: 0 1 180px;
}
.event-select :deep(.el-input__wrapper) {
  padding: 1px 8px;
  border-color: transparent;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.06);
  transition: box-shadow 0.15s ease, background-color 0.15s ease;
  background-color: #f7f8fa;
  border-radius: 6px;
}
.event-select :deep(.el-input__inner) {
  font-size: 14px;
  height: 28px;
  line-height: 28px;
}
.event-select :deep(.el-input__wrapper:hover) {
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.10);
  background-color: #f5f7fa;
}
.event-select :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px var(--el-color-primary) inset, 0 2px 8px rgba(64, 158, 255, 0.20);
  background-color: #fff;
}

/* 事件选项样式 */
.event-option {
  display: flex;
  align-items: center;
  gap: 8px;
}

.event-id {
  color: #303133;
  font-weight: 500;
}

.event-name {
  color: #909399;
  font-size: 12px;
}

/* 自定义公式组 */
.custom-formula-group {
  flex: 1;
  min-width: 0;
  overflow-x: auto;
}

/* 格式选择器 */
.format-select {
  width: 85px;
  flex-shrink: 0;
  margin-left: 8px;
}

/* 操作按钮组 */
.metric-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-left: auto;
  flex-shrink: 0;
  min-width: fit-content;
}

/* 图标按钮 */
.icon-button { 
  min-width: 24px;
  width: 24px;
  height: 24px; 
  padding: 0; 
  background: transparent;
  border: none;
  color: var(--el-text-color-secondary); 
  border-radius: 4px; 
  transition: all 0.15s ease; 
  font-size: 13px;
  flex-shrink: 0;
}
.icon-button:hover { 
  background: var(--el-color-primary-light-9); 
  color: var(--el-color-primary); 
}
.icon-button.is-active { 
  color: var(--el-color-primary);
  background: var(--el-color-primary-light-9);
}

/* 过滤器展开区域 */
.inline-filter-container { 
  margin-top: 6px; 
  padding-top: 6px;
  padding-left: 2px;
  width: 100%;
  border-top: 1px solid rgba(0, 0, 0, 0.06);
}
</style>

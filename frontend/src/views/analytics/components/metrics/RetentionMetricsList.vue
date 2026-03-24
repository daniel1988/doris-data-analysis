<template>
  <el-card class="metrics-card" shadow="never">
    <template #header>
      <div class="card-header">
        <span class="title">留存指标配置</span>
        <span class="subtitle">分析用户从初始事件到回访事件的留存情况</span>
      </div>
    </template>

    <div class="metrics-flow">
      <!-- 初始事件 -->
      <div class="metric-section">
        <div class="section-header">
          <div class="step-badge start">A</div>
          <div class="label-group">
            <span class="section-label">初始事件</span>
            <span class="section-desc">用户分析的目标起点</span>
          </div>
        </div>
        <MetricItem :index="0" :m="initMetric" class="retention-metric-item" 
          :hide-formula="true" :hide-copy="true"
          @remove="() => { }"
          @edit-name="initMetric.isEditingName = true" @finish-edit-name="initMetric.isEditingName = false"
          @toggle-formula="() => { }" @copy="() => { }"
          @toggle-filter="initMetric.filterPopoverVisible = !initMetric.filterPopoverVisible"
          @event-change="$emit('event-change', initMetric)">
          <template #filters>
            <MetricFilterEditor :metric="initMetric" :project-alias="projectAlias" />
          </template>
          <template #property>
            <MetricPropertyPicker :event-id="initMetric.e_event_id" :field="initMetric.metric.column.field"
              :table="initMetric.metric.column.table" :formula="initMetric.metric.formula"
              :property-name="getDisplayPropertyName(initMetric)"
              @select="(payload: any) => $emit('apply-formula', { m: initMetric, ...payload })" />
          </template>
        </MetricItem>
      </div>

      <!-- 连接箭头 -->
      <div class="flow-connector">
        <div class="connector-line"></div>
        <el-icon class="connector-arrow">
          <Bottom />
        </el-icon>
      </div>

      <!-- 结束事件 -->
      <div class="metric-section">
        <div class="section-header">
          <div class="step-badge end">B</div>
          <div class="label-group">
            <span class="section-label">回访事件</span>
            <span class="section-desc">后续发生的转化行为</span>
          </div>
        </div>
        <MetricItem :index="1" :m="endMetric" class="retention-metric-item" 
          :hide-formula="true" :hide-copy="true"
          @remove="() => { }"
          @edit-name="endMetric.isEditingName = true" @finish-edit-name="endMetric.isEditingName = false"
          @toggle-formula="() => { }" @copy="() => { }"
          @toggle-filter="endMetric.filterPopoverVisible = !endMetric.filterPopoverVisible"
          @event-change="$emit('event-change', endMetric)">
          <template #filters>
            <MetricFilterEditor :metric="endMetric" :project-alias="projectAlias" />
          </template>
          <template #property>
            <MetricPropertyPicker :event-id="endMetric.e_event_id" :field="endMetric.metric.column.field"
              :table="endMetric.metric.column.table" :formula="endMetric.metric.formula"
              :property-name="getDisplayPropertyName(endMetric)"
              @select="(payload: any) => $emit('apply-formula', { m: endMetric, ...payload })" />
          </template>
        </MetricItem>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { Bottom } from '@element-plus/icons-vue';
import MetricFilterEditor from './MetricFilterEditor.vue';
import MetricItem from './MetricItem.vue';
import MetricPropertyPicker from './MetricPropertyPicker.vue';

defineProps<{
  initMetric: any
  endMetric: any
  projectAlias: string
  getDisplayPropertyName: (m: any) => string
}>()

defineEmits<{
  (e: 'event-change', m: any): void
  (e: 'apply-formula', payload: any): void
}>()
</script>

<style scoped lang="scss">
.metrics-card {
  border-radius: 12px;
  height: auto !important; // 确保高度自适应
  flex-shrink: 0; // 不允许被压缩

  :deep(.el-card__body) {
    padding: 12px; // 压缩内边距
    height: auto !important;
    overflow: visible !important; // 严禁出现滚动条
    display: flex;
    flex-direction: column;
  }

  .card-header {
    .title {
      font-weight: bold;
      font-size: 14px;
    }

    .subtitle {
      font-size: 12px;
      color: #909399;
      margin-left: 8px;
    }
  }
}

.metrics-flow {
  display: flex;
  flex-direction: column;
  padding: 4px 0;
  gap: 8px; // 减小间距
}

.metric-section {
  position: relative;
  min-height: 80px; // 进一步适度减小最小高度

  .section-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 6px; // 进一步减小间距
    padding-left: 4px;

    .step-badge {
      width: 24px;
      height: 24px;
      border-radius: 6px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: bold;
      font-size: 12px;

      &.start {
        background-color: var(--el-color-primary-light-9);
        color: var(--el-color-primary);
        border: 1px solid var(--el-color-primary-light-5);
      }

      &.end {
        background-color: var(--el-color-success-light-9);
        color: var(--el-color-success);
        border: 1px solid var(--el-color-success-light-5);
      }
    }

    .label-group {
      display: flex;
      flex-direction: column;

      .section-label {
        font-size: 13px;
        font-weight: 600;
        color: var(--el-text-color-primary);
      }

      .section-desc {
        font-size: 11px;
        color: var(--el-text-color-secondary);
      }
    }
  }
}

.flow-connector {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin: 2px 0 2px 12px; // 进一步压缩外边距
  width: 24px;

  .connector-line {
    width: 2px;
    height: 12px; // 缩短连接线长度
    background: repeating-linear-gradient(to bottom,
        var(--el-border-color-lighter),
        var(--el-border-color-lighter) 4px,
        transparent 4px,
        transparent 8px);
  }

  .connector-arrow {
    font-size: 14px;
    color: var(--el-border-color);
    margin-top: -2px;
  }
}

/* 隐藏 Retention 中不需要的 MetricItem 元素 */
:deep(.retention-metric-item) {

  .metric-number,
  .drag-handle,
  .metric-close {
    display: none !important;
  }

  padding-left: 16px;
  padding-right: 16px;
  background: var(--el-fill-color-blank);
  border: 1px solid var(--el-border-color-lighter);
  box-shadow: none;
  min-width: 0 !important; // 覆盖默认的 380px，防止溢出
  width: 100%;

  &:hover {
    border-color: var(--el-color-primary-light-5);
    background: var(--el-fill-color-light);
  }

  // 确保内部选择器也能自适应
  .metric-header-row {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    align-items: center;
    width: 100%;

    .metric-name-section {
      flex: 1;
      min-width: 100px;
    }

    .format-select {
      width: 80px !important;
    }

    .metric-actions {
      margin-left: auto;
    }
  }

  .metric-config-row {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    align-items: center;
    width: 100%;

    .event-select {
      flex: 1;
      min-width: 150px;
    }

    .connector {
      margin: 0 4px;
    }
  }
}
</style>

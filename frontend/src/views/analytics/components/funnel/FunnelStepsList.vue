<template>
  <el-card class="funnel-steps-card" shadow="never">
    <template #header>
      <div class="card-header">
        <span class="title">漏斗步骤</span>
        <el-button 
          size="small" 
          type="primary" 
          icon="Plus" 
          @click="$emit('add-step')"
          :disabled="steps.length >= 8"
        >
          添加步骤
        </el-button>
      </div>
    </template>
    
    <div class="steps-list">
      <FunnelStepItem
        v-for="(step, index) in steps"
        :key="step._uid"
        :index="index"
        :m="step"
        :project-events="projectEvents"
        :can-remove="steps.length > 2"
        @remove="$emit('remove-step', index)"
        @toggle-filter="$emit('toggle-filter', step)"
        @event-change="$emit('event-change', step)"
      >
        <template #filters>
          <MetricFilterEditor 
            :metric="step" 
            :project-alias="projectAlias" 
            property-select-mode="event-only"
          />
        </template>
      </FunnelStepItem>
    </div>
  </el-card>

  <!-- 转化窗口期配置 -->
  <el-card class="window-config-card mt-12" shadow="never">
    <template #header>
      <div class="card-header">
        <span class="title">转化窗口期</span>
        <el-tooltip content="用户从第一步开始，在多长时间内完成整个漏斗流程" placement="top">
          <el-icon class="help-icon"><QuestionFilled /></el-icon>
        </el-tooltip>
      </div>
    </template>
    <div class="window-config-content">
      <el-input-number
        v-model="windowNumValue"
        :min="1"
        :max="90"
        size="default"
      />
      <span class="unit ml-10">天</span>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { Plus, QuestionFilled } from '@element-plus/icons-vue'
import { computed } from 'vue'
import FunnelStepItem from './FunnelStepItem.vue'
import MetricFilterEditor from '../metrics/MetricFilterEditor.vue'

const props = defineProps<{
  steps: any[]
  projectEvents: any[]
  projectAlias: string
  windowNum: number
}>()

const emit = defineEmits(['add-step', 'remove-step', 'toggle-filter', 'event-change', 'update:window-num'])

const windowNumValue = computed({
  get: () => props.windowNum,
  set: (v) => emit('update:window-num', v)
})
</script>

<style scoped lang="scss">
.funnel-steps-card {
  border-radius: 12px;
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    .title {
      font-weight: bold;
      font-size: 14px;
    }
  }
}

.window-config-card {
  border-radius: 12px;
  .card-header {
    display: flex;
    align-items: center;
    gap: 8px;
    .title {
      font-weight: bold;
      font-size: 14px;
    }
    .help-icon {
      font-size: 14px;
      color: #909399;
      cursor: help;
    }
  }
}

.window-config-content {
  display: flex;
  align-items: center;
  .unit {
    font-size: 14px;
    color: #606266;
  }
}

.mt-12 {
  margin-top: 12px;
}

.ml-10 {
  margin-left: 10px;
}
</style>

<template>
  <el-card class="time-controls-card" shadow="never">
    <div class="time-controls-container">
      <!-- 时间粒度 -->
      <div v-if="!hideInterval" class="control-item">
        <span class="label">粒度</span>
        <el-select v-model="innerInterval" size="small" style="width: 80px">
          <el-option label="合计" :value="1" />
          <el-option label="按天" :value="2" />
          <el-option label="按周" :value="3" />
          <el-option label="按月" :value="4" />
          <el-option label="小时" :value="7" />
        </el-select>
      </div>

      <!-- 时间范围 -->
      <div class="control-item">
        <span class="label">时间</span>
        <el-date-picker
          v-model="innerStaticRange"
          type="daterange"
          size="small"
          range-separator="-"
          start-placeholder="开始"
          end-placeholder="结束"
          value-format="YYYY-MM-DD"
          :shortcuts="dateShortcuts"
          style="width: 240px"
        />
      </div>

      <!-- 动态时间快捷选项 -->
      <div class="control-item">
        <el-dropdown trigger="click" @command="handleDynamicCommand">
          <el-button size="small" type="primary" plain>
            动态时间<el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #footer>
            <el-dropdown-menu>
              <el-dropdown-item command="past-7">过去 7 天</el-dropdown-item>
              <el-dropdown-item command="past-30">过去 30 天</el-dropdown-item>
              <el-dropdown-item command="recent-7">最近 7 天 (含今天)</el-dropdown-item>
              <el-dropdown-item command="recent-30">最近 30 天 (含今天)</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>

      <!-- 对比开关 -->
      <div v-if="!hideComparison" class="control-item ml-auto">
        <el-checkbox v-model="showComparison" size="small">开启对比</el-checkbox>
        <el-date-picker
          v-if="showComparison"
          v-model="innerComparisonRange"
          type="daterange"
          size="small"
          range-separator="-"
          start-placeholder="对比开始"
          end-placeholder="对比结束"
          value-format="YYYY-MM-DD"
          style="width: 240px; margin-left: 10px"
        />
      </div>

      <!-- 额外的插槽内容（用于插入留存天数选择器等） -->
      <div class="control-item extra-slot ml-10">
        <slot name="extra"></slot>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ArrowDown } from '@element-plus/icons-vue';
import { computed, ref } from 'vue';

const props = defineProps<{
  interval?: number
  staticRange: [string, string] | null
  comparisonRange: [string, string] | null
  hideInterval?: boolean
  hideComparison?: boolean
}>()

const emit = defineEmits([
  'update:interval', 
  'update:static-range', 
  'update:comparison-range',
  'dynamic-change'
])

const showComparison = ref(false)

const innerInterval = computed({
  get: () => props.interval ?? 1,
  set: (v) => emit('update:interval', v)
})

const innerStaticRange = computed({
  get: () => props.staticRange,
  set: (v) => emit('update:static-range', v)
})

const innerComparisonRange = computed({
  get: () => props.comparisonRange,
  set: (v) => emit('update:comparison-range', v)
})

const handleDynamicCommand = (command: string) => {
  const [type, days] = command.split('-')
  emit('dynamic-change', { 
    type: type === 'past' ? 'past' : 'recent', 
    amountOrDate: days 
  })
}

const dateShortcuts = [
  {
    text: '最近7天',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
      return [start, end]
    },
  },
  {
    text: '最近30天',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
      return [start, end]
    },
  }
]
</script>

<style scoped lang="scss">
.time-controls-card {
  margin-bottom: 15px;
  :deep(.el-card__body) {
    padding: 10px 15px;
  }
}

.time-controls-container {
  display: flex;
  align-items: center;
  gap: 20px;

  .control-item {
    display: flex;
    align-items: center;
    gap: 8px;

    .label {
      font-size: 13px;
      color: var(--el-text-color-secondary);
    }
  }
}

.ml-auto { margin-left: auto; }
</style>


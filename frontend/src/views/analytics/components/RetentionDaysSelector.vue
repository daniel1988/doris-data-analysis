<template>
  <div class="days-selector-popover">
    <el-popover
      placement="bottom"
      :width="320"
      trigger="click"
      popper-class="retention-days-popover"
    >
      <template #reference>
        <el-button size="default" class="days-trigger-btn">
          <el-icon class="mr-4"><Calendar /></el-icon>
          留存时间点
          <el-tag size="small" type="info" class="ml-4">{{ modelValue.length }}</el-tag>
        </el-button>
      </template>

      <div class="popover-content">
        <div class="popover-header">
          <span class="title">配置留存时间点</span>
          <span class="subtitle">(最多支持 10 个)</span>
        </div>
        
        <div class="days-list">
          <div v-for="(day, index) in modelValue" :key="index" class="day-tag">
            <el-input-number
              v-model="modelValue[index]"
              :min="1"
              :max="365"
              size="small"
              class="day-input"
              :controls="false"
              @change="handleDayChange"
            />
            <span class="day-label">{{ unitLabel }}后</span>
            <el-icon class="remove-btn" @click="removeDay(index)"><Close /></el-icon>
          </div>
          
          <el-button 
            v-if="modelValue.length < 10" 
            size="small" 
            icon="Plus" 
            @click="addDay"
            class="add-btn"
          >
            添加时间点
          </el-button>
        </div>

        <div class="popover-footer mt-15">
          <div class="preset-groups">
            <span class="label">常用：</span>
            <el-button-group size="small">
              <el-button @click="setPreset([1, 3, 7])">1/3/7{{ unitLabel }}</el-button>
              <el-button @click="setPreset([1, 7, 30])">1/7/30{{ unitLabel }}</el-button>
            </el-button-group>
          </div>
        </div>
      </div>
    </el-popover>
  </div>
</template>

<script setup lang="ts">
import { Calendar, Close, Plus } from '@element-plus/icons-vue'
import { computed } from 'vue'

const props = defineProps<{
  modelValue: number[]
  interval: string | number // 可能是字符串或数字
}>()

const emit = defineEmits(['update:modelValue'])

const unitLabel = computed(() => {
  // 统一处理 interval，可能是 TimeGrainInterval 枚举值或数字
  const val = String(props.interval)
  if (val === '2' || val === 'day') return '天'
  if (val === '3' || val === 'week') return '周'
  if (val === '4' || val === 'month') return '月'
  return '天'
})

function handleDayChange() {
  const sorted = [...props.modelValue].sort((a, b) => a - b)
  emit('update:modelValue', sorted)
}

function addDay() {
  if (props.modelValue.length >= 10) return
  const nextDay = props.modelValue.length > 0 ? Math.max(...props.modelValue) + 1 : 1
  emit('update:modelValue', [...props.modelValue, nextDay].sort((a, b) => a - b))
}

function removeDay(index: number) {
  const newDays = [...props.modelValue]
  newDays.splice(index, 1)
  emit('update:modelValue', newDays)
}

function setPreset(days: number[]) {
  emit('update:modelValue', [...days])
}
</script>

<style scoped lang="scss">
.days-trigger-btn {
  display: flex;
  align-items: center;
  font-weight: 500;
  
  .mr-4 {
    margin-right: 4px;
  }
  .ml-4 {
    margin-left: 4px;
  }
}

.popover-content {
  .popover-header {
    margin-bottom: 12px;
    display: flex;
    align-items: center;
    gap: 8px;
    
    .title {
      font-weight: 600;
      font-size: 14px;
      color: var(--el-text-color-primary);
    }
    
    .subtitle {
      font-size: 12px;
      color: var(--el-text-color-secondary);
    }
  }
}

.days-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  max-height: 200px;
  overflow-y: auto;
  padding: 4px;
}

.day-tag {
  display: flex;
  align-items: center;
  gap: 4px;
  background: var(--el-fill-color-light);
  padding: 4px 8px;
  border-radius: 4px;
  border: 1px solid var(--el-border-color-lighter);
  
  .day-input {
    width: 44px !important;
    :deep(.el-input__wrapper) {
      padding: 0 4px;
      box-shadow: none;
      background: transparent;
    }
    :deep(.el-input__inner) {
      text-align: center;
    }
  }
  
  .day-label {
    font-size: 12px;
    color: var(--el-text-color-regular);
  }
  
  .remove-btn {
    cursor: pointer;
    font-size: 12px;
    color: var(--el-text-color-placeholder);
    &:hover {
      color: var(--el-color-danger);
    }
  }
}

.add-btn {
  border-style: dashed;
}

.popover-footer {
  border-top: 1px solid var(--el-border-color-lighter);
  padding-top: 12px;
  
  .preset-groups {
    display: flex;
    align-items: center;
    
    .label {
      font-size: 12px;
      color: var(--el-text-color-secondary);
    }
  }
}

.mt-15 {
  margin-top: 15px;
}
</style>

<template>
  <div class="stat-card">
    <div class="stat-value-container">
      <span class="stat-value">{{ formattedValue }}</span>
      <span class="stat-unit" v-if="unit">{{ unit }}</span>
    </div>
    
    <div v-if="hasComparison" class="stat-comparison">
      <div class="comparison-item">
        <span class="comparison-label">环比</span>
        <span :class="['comparison-value', trendClass]">
          <el-icon v-if="trend === 'up'"><CaretTop /></el-icon>
          <el-icon v-if="trend === 'down'"><CaretBottom /></el-icon>
          {{ percentage }}%
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { CaretTop, CaretBottom } from '@element-plus/icons-vue'

const props = defineProps<{
  rows: any[]
  columns: string[]
}>()

const formattedValue = computed(() => {
  if (!props.rows || props.rows.length === 0) return '0'
  const lastRow = props.rows[props.rows.length - 1]
  const metricCol = props.columns.find(c => !c.includes('日期') && !c.includes('时间') && !c.includes('(对比)'))
  if (!metricCol) return '0'
  
  const val = Number(lastRow[metricCol])
  if (isNaN(val)) return lastRow[metricCol]
  
  return val.toLocaleString()
})

const unit = computed(() => {
  // 可以根据指标名称判断单位，暂时留空
  return ''
})

const comparisonData = computed(() => {
  if (!props.rows || props.rows.length === 0) return null
  const lastRow = props.rows[props.rows.length - 1]
  const metricCol = props.columns.find(c => !c.includes('日期') && !c.includes('时间') && !c.includes('(对比)'))
  const compCol = props.columns.find(c => c.includes('(对比)'))
  
  if (!metricCol || !compCol) return null
  
  const current = Number(lastRow[metricCol])
  const previous = Number(lastRow[compCol])
  
  if (isNaN(current) || isNaN(previous) || previous === 0) return null
  
  const change = current - previous
  const percentage = ((change / previous) * 100).toFixed(2)
  
  return {
    percentage: Math.abs(Number(percentage)),
    trend: change >= 0 ? 'up' : 'down'
  }
})

const hasComparison = computed(() => !!comparisonData.value)
const percentage = computed(() => comparisonData.value?.percentage || '0')
const trend = computed(() => comparisonData.value?.trend || 'up')
const trendClass = computed(() => trend.value === 'up' ? 'trend-up' : 'trend-down')
</script>

<style scoped lang="scss">
.stat-card {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 10px;

  .stat-value-container {
    display: flex;
    align-items: baseline;
    gap: 4px;
    
    .stat-value {
      font-size: 28px;
      font-weight: bold;
      color: #303133;
    }
    
    .stat-unit {
      font-size: 14px;
      color: #909399;
    }
  }

  .stat-comparison {
    margin-top: 8px;
    font-size: 13px;
    color: #606266;

    .comparison-item {
      display: flex;
      align-items: center;
      gap: 8px;
    }

    .comparison-label {
      color: #909399;
    }

    .comparison-value {
      display: flex;
      align-items: center;
      font-weight: 500;
      
      &.trend-up {
        color: #f56c6c;
      }
      
      &.trend-down {
        color: #67c23a;
      }
    }
  }
}
</style>

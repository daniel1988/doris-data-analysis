<template>
  <el-card class="user-property-metric-card" shadow="never">
    <template #header>
      <div class="card-header">
        <span class="title">分析指标</span>
      </div>
    </template>
    
    <div class="metric-config">
      <div class="metric-item">
        <span class="label">指标名称</span>
        <el-input v-model="metric.name" placeholder="请输入指标名称" size="default" />
      </div>

      <div class="metric-item mt-15">
        <span class="label">分析属性</span>
        <div class="property-formula-group">
          <UnifiedPropertySelector
            v-model="metric.metric.column.field"
            property-select-mode="user-only"
            placeholder="选择用户属性"
            @select="handlePropertySelect"
          >
            <template #trigger>
              <el-button class="property-btn">
                {{ displayPropertyName || '选择用户属性' }}
                <el-icon class="ml-5"><ArrowDown /></el-icon>
              </el-button>
            </template>
          </UnifiedPropertySelector>

          <el-select
            v-model="metric.metric.formula"
            placeholder="计算方式"
            size="default"
            class="formula-select"
            :disabled="isTotalUsers"
          >
            <el-option
              v-for="f in availableFormulas"
              :key="f.value"
              :label="f.label"
              :value="f.value"
            />
          </el-select>
        </div>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { FormulaDefaults, TableNames } from '@/constants/analysis'
import { Formula } from '@/types/doris/common'
import { MetricFormulaLabels } from '@/types/doris/constants'
import { ArrowDown } from '@element-plus/icons-vue'
import { computed, ref, watch } from 'vue'
import UnifiedPropertySelector from './shared/UnifiedPropertySelector.vue'

const props = defineProps<{
  metric: any
}>()

const emit = defineEmits(['update:metric'])

const selectedProperty = ref<any>(null)

const isTotalUsers = computed(() => {
  return props.metric.metric.column.field === FormulaDefaults.TOTAL_USERS_FIELD
})

const displayPropertyName = computed(() => {
  if (isTotalUsers.value) return '总用户数'
  return selectedProperty.value?.name || props.metric.metric.column.field
})

const availableFormulas = computed(() => {
  if (isTotalUsers.value) {
    return [{ label: '去重计数', value: Formula.CountDistinct }]
  }

  const type = selectedProperty.value?.type?.toLowerCase()
  if (type === 'number' || type === 'int' || type === 'float' || type === 'double' || type === 'decimal') {
    return [
      { label: '总和', value: Formula.Sum },
      { label: '平均值', value: Formula.Avg },
      { label: '最大值', value: Formula.Max },
      { label: '最小值', value: Formula.Min },
      { label: '去重计数', value: Formula.CountDistinct }
    ]
  }

  // 默认仅支持去重计数
  return [{ label: '去重计数', value: Formula.CountDistinct }]
})

function handlePropertySelect(payload: any) {
  selectedProperty.value = payload
  
  const m = props.metric
  m.metric.column.field = payload.id
  m.metric.column.table = TableNames.USER

  // 如果切换到总用户数，强制设置公式
  if (payload.id === FormulaDefaults.TOTAL_USERS_FIELD) {
    m.metric.formula = Formula.CountDistinct
    m.name = '用户数'
  } else {
    // 默认设置为去重计数，如果属性是数值类型，用户可以再改
    m.metric.formula = Formula.CountDistinct
    m.name = `${payload.name} · ${MetricFormulaLabels[Formula.CountDistinct]}`
  }
}

// 监听公式变化，更新指标名称
watch(() => props.metric.metric.formula, (newFormula) => {
  if (isTotalUsers.value) return
  if (selectedProperty.value) {
    props.metric.name = `${selectedProperty.value.name} · ${MetricFormulaLabels[newFormula]}`
  }
})
</script>

<style scoped lang="scss">
.user-property-metric-card {
  border-radius: 12px;
  .card-header {
    .title {
      font-weight: bold;
      font-size: 14px;
    }
  }
}

.metric-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  .label {
    font-size: 12px;
    color: #909399;
    font-weight: 500;
  }
}

.property-formula-group {
  display: flex;
  gap: 12px;
  
  .property-btn {
    flex: 1;
    justify-content: space-between;
    text-align: left;
  }
  
  .formula-select {
    width: 120px;
  }
}

.mt-15 {
  margin-top: 15px;
}

.ml-5 {
  margin-left: 5px;
}
</style>

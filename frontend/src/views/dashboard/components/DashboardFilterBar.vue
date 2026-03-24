<template>
  <div class="dashboard-filter-bar">
    <div class="filter-group">
      <span class="label">时间范围</span>
      <el-date-picker v-model="timeRange" type="daterange" size="default" range-separator="至" start-placeholder="开始日期"
        end-placeholder="结束日期" value-format="YYYY-MM-DD" :shortcuts="dateShortcuts" @change="handleFilterChange" />
    </div>

    <div class="filter-group">
      <span class="label">全局筛选</span>
      <div v-for="(filter, index) in filters" :key="index" class="filter-item">
        <el-select v-model="filter.column.field" placeholder="选择属性" size="default" filterable clearable
          style="width: 160px" @change="handleFilterChange">
          <el-option v-for="prop in properties" :key="prop.id" :label="prop.name" :value="prop.id" />
        </el-select>

        <el-select v-model="filter.operator" placeholder="操作符" size="default" style="width: 100px"
          @change="handleFilterChange">
          <el-option label="等于" :value="1" />
          <el-option label="不等于" :value="2" />
          <el-option label="包含" :value="10" />
          <el-option label="有值" :value="13" />
        </el-select>

        <el-input v-if="![12, 13].includes(filter.operator)" v-model="filter.value" placeholder="值" size="default"
          style="width: 140px" @blur="handleFilterChange" @keyup.enter="handleFilterChange" />

        <el-button type="danger" link icon="Delete" @click="removeFilter(index)" />
      </div>

      <el-button type="primary" link icon="Plus" @click="addFilter">
        添加筛选
      </el-button>
    </div>

    <div class="filter-actions">
      <el-button type="primary" link icon="Refresh" @click="handleFilterChange">全部刷新</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Option } from '@/api/selector'
import { getProperties } from '@/api/selector'
import { useAppStore } from '@/store/app'
import dayjs from 'dayjs'
import { onMounted, ref, watch } from 'vue'

const props = defineProps<{
  modelValue: {
    timeRange: [string, string] | null
    filters: any[]
  }
}>()

const emit = defineEmits(['update:modelValue', 'change'])

const appStore = useAppStore()
const properties = ref<Option[]>([])
const timeRange = ref<[string, string] | null>(props.modelValue?.timeRange || null)
const filters = ref<any[]>(props.modelValue?.filters || [])

watch(() => props.modelValue, (newVal) => {
  if (newVal) {
    timeRange.value = newVal.timeRange
    filters.value = newVal.filters || []
  }
}, { deep: true })

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

const fetchProperties = async () => {
  if (!appStore.activeProjectAlias) return
  try {
    const res = await getProperties({ project_alias: appStore.activeProjectAlias })
    properties.value = res as unknown as Option[]
  } catch (error) {
    console.error('Failed to fetch properties:', error)
  }
}

const addFilter = () => {
  filters.value.push({
    column: { table: 'event_data', field: '', alias: '' },
    operator: 1,
    value: ''
  })
}

const removeFilter = (index: number) => {
  filters.value.splice(index, 1)
  handleFilterChange()
}

const handleFilterChange = () => {
  emit('update:modelValue', {
    timeRange: timeRange.value,
    filters: filters.value
  })
  emit('change')
}

onMounted(() => {
  fetchProperties()
  if (!timeRange.value) {
    const end = dayjs().format('YYYY-MM-DD')
    const start = dayjs().subtract(7, 'day').format('YYYY-MM-DD')
    timeRange.value = [start, end]
    handleFilterChange()
  }
})

watch(() => appStore.activeProjectAlias, () => {
  fetchProperties()
})
</script>

<style scoped lang="scss">
.dashboard-filter-bar {
  background: #fff;
  padding: 12px 20px;
  border-radius: 4px;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 24px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  margin-bottom: 10px;

  .filter-group {
    display: flex;
    align-items: center;
    gap: 12px;

    .label {
      font-size: 14px;
      font-weight: 500;
      color: #606266;
      white-space: nowrap;
    }

    .filter-item {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 4px 8px;
      background: #f5f7fa;
      border-radius: 4px;
    }
  }

  .filter-actions {
    margin-left: auto;
  }
}
</style>

<template>
  <el-card class="global-filters-card" shadow="never">
    <template #header v-if="!hideHeader">
      <div class="card-header">
        <span>{{ propertySelectMode === 'user-only' ? '用户筛选' : '全局筛选' }}</span>
        <el-button type="primary" link icon="Plus" @click="addFilter">添加筛选</el-button>
      </div>
    </template>
    
    <div v-if="hideHeader" class="header-actions-only">
      <el-button type="primary" link icon="Plus" @click="addFilter" size="small">添加筛选</el-button>
    </div>

    <div class="filter-list">
      <div class="filter-header-actions" v-if="filters.length > 1">
        <el-button size="small" type="primary" plain @click="toggleRelation">
          {{ relation === 'and' ? '且 (AND)' : '或 (OR)' }}
        </el-button>
      </div>

      <div v-for="(filter, index) in filters" :key="index" class="filter-row">
        <el-row :gutter="10" align="middle">
          <el-col :span="8">
            <PropertySelect v-model="filter.column.field" :properties="properties" @change="(prop) => handlePropertyChange(filter, prop)" />
          </el-col>
          <el-col :span="6">
            <OperatorSelect v-model="filter.operator" :data-type="filter.column.dataType || DataType.String" />
          </el-col>
          <el-col :span="8">
            <FilterValueInput 
              v-if="![Operator.IsNull, Operator.IsNotNull].includes(filter.operator)"
              v-model="filter.value"
              :data-type="filter.column.dataType || DataType.String"
              :operator="filter.operator"
              :project-alias="context?.state?.projectAlias || appStore.activeProjectAlias || ''"
              :table-name="filter.column.table"
              :field-name="filter.column.field"
            />
          </el-col>
          <el-col :span="2">
            <el-button link type="danger" icon="Delete" @click="removeFilter(index)"></el-button>
          </el-col>
        </el-row>
      </div>
      <div v-if="filters.length === 0" class="empty-filters">
        暂无全局筛选
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { useAppStore } from '@/store/app';
import { DataType, Operator } from '@/types/doris/common';
import { computed, inject } from 'vue';
import { ANALYSIS_CONTEXT_KEY } from '../context';
import PropertySelect from './PropertySelect.vue';
import FilterValueInput from './shared/FilterValueInput.vue';
import OperatorSelect from './shared/OperatorSelect.vue';

const props = defineProps<{
  filterGroup?: any
  hideHeader?: boolean
  propertySelectMode?: 'user-only' | 'event-only' | 'all'
}>()

const context = inject(ANALYSIS_CONTEXT_KEY) as any
const appStore = useAppStore()

const activeFilterGroup = computed(() => {
  if (props.filterGroup) return props.filterGroup
  const form = context?.state?.form as any
  return form?.global_filter_groups || form?.filter_groups
})

const filters = computed(() => {
  return activeFilterGroup.value?.global_filters?.filters || []
})

const relation = computed(() => {
  return activeFilterGroup.value?.global_filters?.relation || 'and'
})

const properties = computed(() => {
  const allProps = context?.state?.metadata?.propertyOptions || []
  if (props.propertySelectMode === 'user-only') {
    return allProps.filter(p => p.table === 'user_data' || p.table === 'user_tag_data')
  }
  return allProps
})

const addFilter = () => {
  const filterGroups = activeFilterGroup.value
  if (!filterGroups) return
  
  if (!filterGroups.global_filters) {
    filterGroups.global_filters = {
      scope: 1,
      relation: 'and',
      filters: [],
      tag_filters: [],
      user_group_filters: []
    }
  }
  
  if (!filterGroups.global_filters.filters) {
    filterGroups.global_filters.filters = []
  }

  const defaultTable = props.propertySelectMode === 'user-only' ? 'user_data' : 'event_data'

  filterGroups.global_filters.filters.push({
    column: { table: defaultTable, field: '', alias: '' },
    operator: Operator.EqualTo,
    value: ''
  })
}

const removeFilter = (index: number) => {
  const filterGroups = activeFilterGroup.value
  if (filterGroups?.global_filters?.filters) {
    filterGroups.global_filters.filters.splice(index, 1)
  }
}

const handlePropertyChange = (filter: any, property: any) => {
  if (property) {
    const fallbackTable = props.propertySelectMode === 'user-only' ? 'user_data' : 'event_data'
    filter.column.table = property.table || filter.column.table || fallbackTable
    filter.column.dataType = property.data_type
  }
}

const toggleRelation = () => {
  const filterGroups = activeFilterGroup.value
  if (filterGroups?.global_filters) {
    filterGroups.global_filters.relation = relation.value === 'and' ? 'or' : 'and'
  }
}
</script>

<style scoped lang="scss">
.global-filters-card {
  margin-bottom: 12px; // 减小间距

  :deep(.el-card__body) {
    padding: 12px; // 压缩内边距
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 14px;
    font-weight: 600;
  }

  .filter-list {
    .filter-header-actions {
      margin-bottom: 8px;
      display: flex;
      justify-content: flex-end;
    }

    .filter-row {
      margin-bottom: 8px; // 减小行间距
      padding-bottom: 8px;
      border-bottom: 1px dashed var(--el-border-color-extra-light);

      &:last-child {
        margin-bottom: 0;
        padding-bottom: 0;
        border-bottom: none;
      }
    }

    .empty-filters {
      font-size: 13px;
      color: var(--el-text-color-placeholder);
      text-align: center;
      padding: 10px 0;
    }
  }
}
</style>

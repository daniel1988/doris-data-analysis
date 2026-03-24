<template>
  <el-card class="global-groups-card" shadow="never">
    <template #header>
      <div class="card-header">
        <span>分组维度</span>
        <el-button type="primary" link icon="Plus" @click="addGroup">添加分组</el-button>
      </div>
    </template>
    
    <div class="group-list">
      <div v-for="(group, index) in groups" :key="index" class="group-row">
        <el-row :gutter="10" align="middle">
          <el-col :span="6">
            <el-select v-model="group.group_type" placeholder="类型" size="small">
              <el-option label="按值" :value="GroupType.Value" />
              <el-option label="按区间" :value="GroupType.Ranges" />
              <el-option label="按日期" :value="GroupType.Date" />
            </el-select>
          </el-col>
          <el-col :span="8">
            <PropertySelect v-model="group.column.field" :properties="properties" />
          </el-col>
          <el-col :span="8">
            <el-input v-model="group.column.alias" placeholder="显示名" size="small" />
          </el-col>
          <el-col :span="2">
            <el-button link type="danger" icon="Delete" @click="removeGroup(index)"></el-button>
          </el-col>
        </el-row>
        
        <!-- 区间分组配置 -->
        <div v-if="group.group_type === GroupType.Ranges" class="range-config mt-10">
          <el-button size="small" type="primary" link icon="Plus" @click="addRange(group)">添加区间</el-button>
          <div v-for="(range, rIdx) in group.value_ranges" :key="rIdx" class="range-row mt-5">
            <el-input-number v-model="range.min" size="small" placeholder="最小" />
            <span class="mx-5">-</span>
            <el-input-number v-model="range.max" size="small" placeholder="最大" />
            <el-button link type="danger" icon="Delete" @click="group.value_ranges.splice(rIdx, 1)" />
          </div>
        </div>
      </div>
      <div v-if="groups.length === 0" class="empty-groups">
        暂无分组维度
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import type { Group } from '@/types/doris/analysis';
import { GroupType } from '@/types/doris/analysis';
import { computed, inject } from 'vue';
import { ANALYSIS_CONTEXT_KEY } from '../context';
import PropertySelect from './PropertySelect.vue';

const context = inject(ANALYSIS_CONTEXT_KEY)
const groups = computed(() => (context?.state.form as any)?.groups || [])
const properties = computed(() => context?.state.metadata.propertyOptions || [])

const addGroup = () => {
  const form = context?.state.form as any
  if (!form) return
  if (!form.groups) form.groups = []
  
  form.groups.push({
    group_type: GroupType.Value,
    column: { table: 'event_data', field: '', alias: '' },
    value_ranges: [],
    time_grain: {
      column: { table: 'event_data', field: 'e_event_time', alias: '' },
      interval: 2,
      window_num: 0
    },
    tag_group: { tag_code: '', operator: 1, tag_value: '' },
    user_group: { group_name: '', group_code: '', operator: 1 }
  })
}

const removeGroup = (index: number) => {
  const form = context?.state.form as any
  if (form?.groups) {
    form.groups.splice(index, 1)
  }
}

const addRange = (group: Group) => {
  if (!group.value_ranges) group.value_ranges = []
  group.value_ranges.push({ min: 0, max: 100 })
}
</script>

<style scoped lang="scss">
.global-groups-card {
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

  .group-list {
    .group-row {
      margin-bottom: 8px; // 减小行间距
      padding-bottom: 8px;
      border-bottom: 1px dashed var(--el-border-color-lighter);

      &:last-child {
        margin-bottom: 0;
        padding-bottom: 0;
        border-bottom: none;
      }
    }

    .range-config {
      padding-left: 20px;
      border-left: 2px solid var(--el-color-primary-light-8);
    }
  }

  .empty-groups {
    padding: 10px 0;
    text-align: center;
    color: var(--el-text-color-placeholder);
    font-size: 13px;
  }
}

.mt-10 { margin-top: 10px; }
.mt-5 { margin-top: 5px; }
.mx-5 { margin: 0 5px; }
</style>


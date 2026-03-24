<template>
  <div class="property-user-groups-panel">
    <div v-for="(ug, index) in userGroups" :key="index" class="user-group-item">
      <div class="group-header">
        <span class="group-index">人群 {{ index + 1 }}</span>
        <el-input 
          v-model="ug.alias" 
          placeholder="人群名称（如：高消费用户）" 
          size="small" 
          class="group-name-input"
        />
        <el-button 
          v-if="userGroups.length > 1"
          size="small" 
          type="danger" 
          plain 
          icon="Delete"
          @click="removeUserGroup(index)"
        />
      </div>
      <div class="group-filter">
        <GlobalFiltersCard 
          :filter-group="ug.filter_group" 
          :hide-header="true" 
          :property-select-mode="'user-only'"
        />
      </div>
    </div>
    
    <el-button 
      v-if="userGroups.length < 5" 
      size="small" 
      type="primary" 
      plain 
      icon="Plus" 
      @click="addUserGroup"
      class="add-group-btn mt-12"
    >
      添加人群对比
    </el-button>
  </div>
</template>

<script setup lang="ts">
import { Plus, Delete } from '@element-plus/icons-vue'
import GlobalFiltersCard from './GlobalFiltersCard.vue'

const props = defineProps<{
  userGroups: any[]
}>()

const emit = defineEmits(['update:userGroups'])

function addUserGroup() {
  const newGroups = [...props.userGroups, {
    alias: `人群 ${props.userGroups.length + 1}`,
    filter_group: {
      global_filters: { scope: 1, filters: [], tag_filters: [], user_group_filters: [] },
      query_dates: []
    }
  }]
  emit('update:userGroups', newGroups)
}

function removeUserGroup(index: number) {
  const newGroups = [...props.userGroups]
  newGroups.splice(index, 1)
  emit('update:userGroups', newGroups)
}
</script>

<style scoped lang="scss">
.property-user-groups-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.user-group-item {
  border: 1px dashed #dcdfe6;
  border-radius: 8px;
  padding: 12px;
  background: #f8f9fa;
  
  .group-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 12px;
    .group-index {
      font-size: 12px;
      color: #909399;
      font-weight: bold;
    }
    .group-name-input {
      width: 200px;
    }
  }
}

.add-group-btn {
  width: 100%;
  border-style: dashed;
}

.mt-12 {
  margin-top: 12px;
}
</style>

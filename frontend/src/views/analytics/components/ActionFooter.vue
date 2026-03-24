<template>
  <div class="action-footer">
    <div class="left-actions">
      <el-button type="info" plain @click="$emit('reset')" class="action-btn">
        <el-icon class="el-icon--left"><Refresh /></el-icon>重置配置
      </el-button>
    </div>
    <div class="right-actions">
      <el-button type="success" plain @click="$emit('save')" class="action-btn">
        <el-icon class="el-icon--left"><DocumentChecked /></el-icon>保存报表
      </el-button>
      <el-button 
        type="primary" 
        :loading="loading" 
        @click="$emit('analyze')" 
        class="action-btn analyze-btn"
      >
        <el-icon v-if="!loading" class="el-icon--left"><TrendCharts /></el-icon>
        {{ loading ? '计算中...' : '开始计算' }}
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { DocumentChecked, Refresh, TrendCharts } from '@element-plus/icons-vue';

defineProps<{
  loading: boolean
}>()

defineEmits(['analyze', 'save', 'reset'])
</script>

<style scoped lang="scss">
.action-footer {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  box-sizing: border-box;
  background-color: var(--el-bg-color);
  border-top: 1px solid var(--el-border-color-lighter);
  
  /* 确保在父容器中置底 */
  margin-top: auto;

  .right-actions {
    display: flex;
    gap: 12px;
  }

  .action-btn {
    font-weight: 500;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    
    &:hover {
      transform: translateY(-1px);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    }

    &:active {
      transform: translateY(0);
    }
  }

  .analyze-btn {
    min-width: 120px;
    background: linear-gradient(135deg, var(--el-color-primary) 0%, var(--el-color-primary-light-3) 100%);
    border: none;
    color: white;

    &:hover {
      opacity: 0.9;
      box-shadow: 0 6px 16px rgba(var(--el-color-primary-rgb), 0.3);
    }
    
    &.is-loading {
      opacity: 0.8;
      background: var(--el-color-primary-light-5);
    }
  }
}
</style>

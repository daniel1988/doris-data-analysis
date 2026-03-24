<template>
  <el-card class="metrics-card" shadow="never">
    <template #header>
      <div class="section-header">
        <span>事件指标</span>
        <div class="header-actions">
          <el-button size="small" type="primary" :icon="Plus" @click="$emit('add-metric')" title="添加事件指标" />
          <el-button size="small" type="primary" plain @click="$emit('add-custom')" title="添加自定义公式指标">Σ</el-button>
        </div>
      </div>
    </template>
    <slot />
  </el-card>
</template>

<script setup lang="ts">
import { Plus } from '@element-plus/icons-vue';
defineEmits<{ (e:'add-metric'): void; (e:'add-custom'): void }>()
</script>

<style scoped>
.metrics-card { 
  border-radius: 12px; 
  border: 1px solid rgba(0, 0, 0, 0.08);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}
.metrics-card :deep(.el-card__header) {
  padding: 14px 18px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}
.metrics-card :deep(.el-card__body) {
  padding: 12px 18px;
}
.section-header { 
  display: flex; 
  justify-content: space-between; 
  align-items: center; 
  font-weight: 600; 
  font-size: 15px; 
  color: #1f2937;
}
.header-actions { 
  display: flex; 
  gap: 8px; 
  align-items: center; 
}
</style>

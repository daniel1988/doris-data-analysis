<template>
  <el-card class="scatter-settings-card" shadow="never">
    <template #header>
      <div class="card-header">
        <span class="title">分布设置</span>
      </div>
    </template>

    <div class="settings-content">
      <el-form :model="form" label-position="top">
        <el-form-item label="分布类型">
          <el-radio-group v-model="form.scatter_type" class="type-radio-group">
            <el-radio :label="1">离散值</el-radio>
            <el-radio :label="2">自动分桶</el-radio>
            <el-radio :label="3">自定义区间</el-radio>
          </el-radio-group>
        </el-form-item>

        <div v-if="form.scatter_type === 2" class="auto-settings">
          <el-form-item label="分桶数量">
            <el-input-number v-model="form.bin_count" :min="1" :max="100" />
          </el-form-item>
        </div>

        <div v-if="form.scatter_type === 3" class="custom-settings">
          <div class="range-header">
            <span>区间配置</span>
            <el-button type="primary" link :icon="Plus" @click="addRange">添加区间</el-button>
          </div>
          
          <div v-for="(range, index) in form.scatter_ranges" :key="index" class="range-item">
            <el-input-number v-model="range.min" placeholder="最小值" :controls="false" />
            <span class="separator">-</span>
            <el-input-number v-model="range.max" placeholder="最大值" :controls="false" />
            <el-button type="danger" link :icon="Delete" @click="removeRange(index)" />
          </div>

          <div v-if="!form.scatter_ranges?.length" class="empty-ranges">
            <el-empty description="暂无自定义区间" :image-size="40" />
          </div>
        </div>
      </el-form>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { Delete, Plus } from '@element-plus/icons-vue'
import { computed, inject } from 'vue'
import { ANALYSIS_CONTEXT_KEY } from '../context'

const context = inject(ANALYSIS_CONTEXT_KEY)
const form = computed(() => (context?.state.form as any))

function addRange() {
  if (!form.value.scatter_ranges) form.value.scatter_ranges = []
  form.value.scatter_ranges.push({ min: 0, max: 0 })
}

function removeRange(index: number) {
  form.value.scatter_ranges.splice(index, 1)
}
</script>

<style scoped lang="scss">
.scatter-settings-card {
  border: none;
  :deep(.el-card__header) {
    padding: 12px 16px;
    border-bottom: 1px solid var(--el-border-color-lighter);
  }
  :deep(.el-card__body) {
    padding: 12px;
  }
}

.card-header {
  .title {
    font-size: 14px;
    font-weight: 600;
    color: var(--el-text-color-primary);
  }
}

.type-radio-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
  :deep(.el-radio) {
    margin-right: 0;
  }
}

.range-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  span {
    font-size: 12px;
    color: var(--el-text-color-secondary);
  }
}

.range-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  .separator {
    color: var(--el-text-color-placeholder);
  }
  :deep(.el-input-number) {
    flex: 1;
    .el-input__inner {
      text-align: left;
    }
  }
}

.empty-ranges {
  padding: 12px;
  border: 1px dashed var(--el-border-color);
  border-radius: 4px;
  :deep(.el-empty) {
    padding: 0;
  }
}
</style>

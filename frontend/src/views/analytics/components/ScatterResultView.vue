<template>
  <div class="scatter-result-view">
    <el-card shadow="never" class="result-card">
      <template #header>
        <div class="result-header">
          <div class="header-left">
            <span class="title">分布分析结果</span>
            <el-button v-if="sql" link type="primary" size="small" icon="View" @click="$emit('show-sql')"
              class="ml-10">查看SQL</el-button>
          </div>
          <div class="header-right">
            <el-radio-group :model-value="vizMode" size="small" @update:model-value="$emit('update:vizMode', $event)">
              <el-radio-button value="bar">柱状图</el-radio-button>
              <el-radio-button value="pie">饼图</el-radio-button>
              <el-radio-button value="table">表格</el-radio-button>
            </el-radio-group>
            <div class="ops ml-10">
              <el-button link type="primary" size="small" icon="Picture" @click="handleSaveImage">另存为图片</el-button>
            </div>
          </div>
        </div>
      </template>

      <!-- 图表区域 -->
      <div v-if="vizMode !== 'table'" v-loading="loading" class="chart-section">
        <ChartView ref="chartViewRef" :rows="rows" :columns="columns" :viz-mode="vizMode" />
      </div>

      <!-- 表格区域 -->
      <div v-if="vizMode === 'table' || true" v-loading="loading" :class="['table-section', vizMode === 'table' ? '' : 'mt-20']">
        <TableView :rows="rows" :columns="columns" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import ChartView from './ChartView.vue';
import TableView from './TableView.vue';

const props = defineProps<{
  loading: boolean
  rows: any[]
  columns: string[]
  vizMode: string
  sql?: string
}>()

const emit = defineEmits(['update:vizMode', 'show-sql'])

const chartViewRef = ref<InstanceType<typeof ChartView>>()

const handleSaveImage = () => {
  chartViewRef.value?.saveAsImage()
}
</script>

<style scoped lang="scss">
.scatter-result-view {
  .result-card {
    border: none;
    :deep(.el-card__header) {
      padding: 12px 16px;
      border-bottom: 1px solid var(--el-border-color-lighter);
    }
    .result-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .header-left {
        display: flex;
        align-items: center;

        .title {
          font-weight: bold;
          font-size: 14px;
        }
      }

      .header-right {
        display: flex;
        align-items: center;
      }
    }

    .chart-section {
      height: 400px;
      width: 100%;
      position: relative;
      margin-top: 10px;
    }

    .table-section {
      margin-top: 20px;
    }
  }
}

.ml-10 {
  margin-left: 10px;
}

.mt-20 {
  margin-top: 20px;
}
</style>

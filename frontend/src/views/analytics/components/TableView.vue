<template>
  <div class="table-view">
    <div v-if="!hideHeader" class="table-header-bar mb-10">
      <span class="section-title">数据明细</span>
      <div class="table-actions">
        <el-button size="small" icon="Download" @click="exportToExcel">导出Excel</el-button>
      </div>
    </div>

    <el-table :data="displayRows" style="width: 100%" border stripe size="small" :height="hideHeader ? '100%' : 320">
      <el-table-column v-for="col in columns" :key="col" :prop="col" :label="col" min-width="120" sortable>
        <template #default="{ row }">
          {{ formatValue(row[col], col) }}
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-wrap mt-10">
      <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[10, 20, 50, 100]"
        small background layout="total, sizes, prev, pager, next" :total="rows?.length || 0" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import * as XLSX from 'xlsx';

const props = withDefaults(defineProps<{
  rows: any[]
  columns: string[]
  hideHeader?: boolean
}>(), {
  hideHeader: false
})

const currentPage = ref(1)
const pageSize = ref(20)

const displayRows = computed(() => {
  if (!props.rows) return []
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return props.rows.slice(start, end)
})

const formatValue = (val: any, colName: string) => {
  if (typeof val !== 'number') return val

  // Percent formatting for columns containing '%' or '率'
  if (colName.includes('%') || colName.includes('率')) {
    return (val * 100).toFixed(2) + '%'
  }

  // Thousand separator for large numbers
  if (val > 1000) {
    return val.toLocaleString()
  }

  return val
}

const exportToExcel = () => {
  const worksheet = XLSX.utils.json_to_sheet(props.rows || [])
  const workbook = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(workbook, worksheet, 'Data')
  XLSX.writeFile(workbook, `analysis-data-${Date.now()}.xlsx`)
}
</script>

<style scoped lang="scss">
.table-view {
  height: 100%;
  display: flex;
  flex-direction: column;

  .table-header-bar {
    flex-shrink: 0;
    display: flex;
    justify-content: space-between;
    align-items: center;

    .section-title {
      font-weight: bold;
      font-size: 14px;
      color: #333;
      border-left: 4px solid var(--el-color-primary);
      padding-left: 10px;
    }
  }

  .el-table {
    flex: 1;
    min-height: 0;
  }

  .pagination-wrap {
    flex-shrink: 0;
    display: flex;
    justify-content: flex-end;
  }
}

.mb-10 {
  margin-bottom: 10px;
}

.mt-10 {
  margin-top: 10px;
}
</style>

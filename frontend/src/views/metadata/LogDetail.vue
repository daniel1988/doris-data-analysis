<template>
  <div class="log-detail-container">
    <div class="page-header mb-15">
      <div class="header-left">
        <span class="page-title">日志明细</span>
        <p class="page-desc">查询并查看原始事件日志数据</p>
      </div>
    </div>

    <!-- 搜索栏 -->
    <el-card class="search-card mb-15" shadow="never">
      <el-form :model="searchForm" inline @submit.prevent="handleSearch">
        <el-form-item label="用户ID">
          <el-input v-model="searchForm.openId" placeholder="请输入用户ID" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="事件ID">
          <el-select v-model="searchForm.eventId" placeholder="请选择事件" clearable filterable style="width: 200px">
            <el-option v-for="item in projectEvents" :key="item.id" :label="item.name || item.id" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="事件时间">
          <el-date-picker
            v-model="searchForm.eventTimeRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            :shortcuts="dateShortcuts"
            style="width: 360px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch" :loading="loading">查询</el-button>
          <el-button :icon="Refresh" @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-card class="table-card" shadow="never" v-loading="loading">
      <el-table :data="tableData" border stripe style="width: 100%" height="calc(100vh - 350px)">
        <el-table-column prop="e_event_time" label="事件时间" width="180" sortable fixed="left">
          <template #default="{ row }">
            {{ formatEventTime(row.e_event_time) }}
          </template>
        </el-table-column>
        <el-table-column prop="e_openid" label="用户ID" width="200">
          <template #default="{ row }">
            <span class="copyable-text" @dblclick="copyToClipboard(row.e_openid, '用户ID')">{{ row.e_openid }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="e_event_id" label="事件ID" width="150">
          <template #default="{ row }">
            <el-tag size="small">{{ row.e_event_id }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="e_event_name" label="事件名称" width="150" show-overflow-tooltip />
        <el-table-column prop="e_package_name" label="应用包名" width="180" show-overflow-tooltip />
        <el-table-column prop="e_platform" label="平台" width="100">
          <template #default="{ row }">
            <el-tag :type="getPlatformType(row.e_platform)" size="small">{{ row.e_platform }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="e_ip" label="IP地址" width="140" />
        <el-table-column prop="e_request_id" label="请求ID" width="200" show-overflow-tooltip />
        <el-table-column label="操作" width="100" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link :icon="View" @click="viewProperties(row.e_properties)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container mt-15">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="pagination.total"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 属性详情抽屉 -->
    <el-drawer v-model="propertyDrawerVisible" title="事件属性详情" size="500px">
      <div class="property-json-wrapper">
        <pre v-if="formattedProperties">{{ formattedProperties }}</pre>
        <el-empty v-else description="暂无属性数据" />
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { useAppStore } from '@/store/app'
import { Refresh, Search, View } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { onMounted, watch } from 'vue'
import { useEventDetail } from './composables/useEventDetail'

const appStore = useAppStore()
const {
  loading,
  tableData,
  searchForm,
  pagination,
  dateShortcuts,
  propertyDrawerVisible,
  formattedProperties,
  projectEvents,
  fetchMeta,
  loadData,
  handleSearch,
  handlePageChange,
  handleSizeChange,
  viewProperties,
  copyToClipboard
} = useEventDetail()

const formatEventTime = (time: string) => {
  if (!time) return '-'
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

const getPlatformType = (platform: string) => {
  const p = platform?.toLowerCase()
  if (p === 'android') return 'success'
  if (p === 'ios') return 'warning'
  if (p === 'web') return ''
  return 'info'
}

const resetSearch = () => {
  searchForm.openId = ''
  searchForm.eventId = ''
  searchForm.eventTimeRange = []
  handleSearch()
}

onMounted(async () => {
  await fetchMeta()
  if (appStore.activeProjectAlias) {
    loadData()
  }
})

watch(() => appStore.activeProjectAlias, async (newVal) => {
  if (newVal) {
    await fetchMeta()
    handleSearch()
  }
})
</script>

<style scoped lang="scss">
.log-detail-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.page-header {
  .page-title {
    font-size: 18px;
    font-weight: bold;
    color: #303133;
  }
  .page-desc {
    font-size: 13px;
    color: #909399;
    margin-top: 5px;
  }
}

.search-card {
  :deep(.el-card__body) {
    padding-bottom: 2px;
  }
}

.table-card {
  flex: 1;
  min-height: 0;
  
  :deep(.el-card__body) {
    height: 100%;
    display: flex;
    flex-direction: column;
  }
}

.copyable-text {
  cursor: pointer;
  &:hover {
    color: var(--el-color-primary);
    text-decoration: underline;
  }
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
}

.property-json-wrapper {
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 4px;
  height: 100%;
  overflow: auto;
  
  pre {
    margin: 0;
    font-family: monospace;
    font-size: 13px;
    line-height: 1.5;
    color: #303133;
    white-space: pre-wrap;
    word-break: break-all;
  }
}

.mb-15 {
  margin-bottom: 15px;
}

.mt-15 {
  margin-top: 15px;
}
</style>

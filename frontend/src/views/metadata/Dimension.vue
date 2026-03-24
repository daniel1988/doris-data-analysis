<template>
  <div class="metadata-dimension-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>维度管理</span>
          <el-button type="primary" icon="Plus">新增维度</el-button>
        </div>
      </template>

      <el-row :gutter="20">
        <el-col :span="8" v-for="dim in pagedDimensionData" :key="dim.name">
          <el-card shadow="hover" class="dim-card">
            <div class="dim-header">
              <el-icon :size="24" color="#409eff">
                <Management />
              </el-icon>
              <div class="dim-info">
                <div class="dim-name">{{ dim.displayName }}</div>
                <div class="dim-code"><code>{{ dim.name }}</code></div>
              </div>
            </div>
            <div class="dim-body">
              <p>{{ dim.description }}</p>
              <div class="dim-meta">
                <span>包含值: <strong>{{ dim.valueCount }}</strong></span>
                <el-divider direction="vertical" />
                <span>关联字段: <strong>{{ dim.fieldCount }}</strong></span>
              </div>
            </div>
            <div class="dim-footer">
              <el-button type="primary" link>查看详情</el-button>
              <el-button type="primary" link>同步数据</el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <div class="pagination-container">
        <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[6, 12, 24, 48]"
          layout="total, sizes, prev, pager, next, jumper" :total="dimensionData.length" @size-change="handleSizeChange"
          @current-change="handleCurrentChange" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { Management } from '@element-plus/icons-vue'
import { computed, ref } from 'vue'

const dimensionData = ref([
  { name: 'geo_city', displayName: '地理位置-城市', description: '基于 IP 解析或上报的城市维度', valueCount: 342, fieldCount: 1 },
  { name: 'device_brand', displayName: '设备品牌', description: '手机或电脑的品牌维度', valueCount: 56, fieldCount: 2 },
  { name: 'user_level', displayName: '用户等级', description: '用户的会员等级或活跃等级', valueCount: 5, fieldCount: 1 },
  { name: 'campaign_source', displayName: '广告来源', description: '流量渠道的来源维度', valueCount: 120, fieldCount: 3 },
])

const currentPage = ref(1)
const pageSize = ref(6)

const pagedDimensionData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return dimensionData.value.slice(start, end)
})

const handleSizeChange = (val: number) => {
  pageSize.value = val
  currentPage.value = 1
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
}
</script>

<style scoped>
.metadata-dimension-page {
  padding: 0;
}

:deep(.el-card) {
  margin-bottom: 0;
}

.dim-card {
  margin-bottom: 20px;
}

.dim-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
}

.dim-name {
  font-weight: bold;
  font-size: 16px;
}

.dim-code code {
  font-size: 12px;
  color: #909399;
}

.dim-body p {
  font-size: 13px;
  color: #606266;
  height: 40px;
  overflow: hidden;
  text-overflow: ellipsis;
}

.dim-meta {
  font-size: 12px;
  color: #909399;
  margin-top: 10px;
}

.dim-footer {
  margin-top: 15px;
  border-top: 1px solid #f0f0f0;
  padding-top: 10px;
  display: flex;
  justify-content: flex-end;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>

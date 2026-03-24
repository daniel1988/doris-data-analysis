<template>
  <div class="user-management-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-input
              v-model="searchQuery"
              placeholder="搜索用户名/邮箱"
              style="width: 250px"
              prefix-icon="Search"
              clearable
            />
          </div>
          <el-button type="primary" icon="Plus">新增用户</el-button>
        </div>
      </template>

      <el-table :data="pagedData" style="width: 100%" v-loading="loading">
        <el-table-column label="用户" width="250">
          <template #default="scope">
            <div class="user-info-cell">
              <el-avatar :size="32" :src="scope.row.avatar">{{ scope.row.username.charAt(0).toUpperCase() }}</el-avatar>
              <div class="user-details">
                <div class="username">{{ scope.row.username }}</div>
                <div class="email">{{ scope.row.email }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="role" label="角色" width="150">
          <template #default="scope">
            <el-tag :type="getRoleTag(scope.row.role)">{{ scope.row.role }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="tenant" label="所属租户" width="180" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'active' ? 'success' : 'info'" effect="dark">
              {{ scope.row.status === 'active' ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="lastLogin" label="最后登录时间" width="180" />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="scope">
            <el-button type="primary" link>编辑</el-button>
            <el-button type="warning" link>重置密码</el-button>
            <el-dropdown trigger="click">
              <el-button type="info" link style="margin-left: 10px">
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item>权限设置</el-dropdown-item>
                  <el-dropdown-item divided type="danger">禁用账户</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="filteredData.length"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ArrowDown } from '@element-plus/icons-vue'
import { computed, ref } from 'vue'

const searchQuery = ref('')
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)

const userData = ref([
  { id: 1, username: 'admin', email: 'admin@dmp.com', role: '超级管理员', tenant: '默认租户', status: 'active', lastLogin: '2024-03-13 10:30:00', avatar: '' },
  { id: 2, username: 'zhangsan', email: 'zhangsan@corp.com', role: '分析师', tenant: '电商业务部', status: 'active', lastLogin: '2024-03-12 15:45:00', avatar: '' },
  { id: 3, username: 'lisi', email: 'lisi@corp.com', role: '运营', tenant: '市场部', status: 'active', lastLogin: '2024-03-13 09:00:00', avatar: '' },
  { id: 4, username: 'wangwu', email: 'wangwu@test.com', role: '访客', tenant: '外部合作伙伴', status: 'inactive', lastLogin: '2024-03-01 11:20:00', avatar: '' },
  { id: 5, username: 'dev_user', email: 'dev@dmp.com', role: '管理员', tenant: '技术部', status: 'active', lastLogin: '2024-03-13 17:30:00', avatar: '' },
])

const getRoleTag = (role: string) => {
  switch (role) {
    case '超级管理员': return 'danger'
    case '管理员': return 'warning'
    case '分析师': return 'success'
    default: return 'info'
  }
}

const filteredData = computed(() => {
  return userData.value.filter(item => {
    return item.username.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
           item.email.toLowerCase().includes(searchQuery.value.toLowerCase())
  })
})

const pagedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredData.value.slice(start, end)
})
</script>

<style scoped>
.user-management-page {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
}

.user-info-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-details {
  display: flex;
  flex-direction: column;
}

.username {
  font-weight: bold;
  color: #303133;
  font-size: 14px;
}

.email {
  font-size: 12px;
  color: #909399;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>

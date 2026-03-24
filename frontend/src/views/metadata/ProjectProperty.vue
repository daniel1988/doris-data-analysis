<template>
  <div class="event-property-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span class="title">事件属性管理</span>
            <el-tag size="small" type="info" style="margin-left: 10px">{{ projectAlias }}</el-tag>
          </div>
          <div class="header-right" style="display: flex; align-items: center;">
            <el-input v-model="searchQuery" placeholder="搜索属性标识/名称" clearable style="width: 250px; margin-right: 16px">
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-button type="primary" icon="Plus" @click="handleAdd">添加属性</el-button>
          </div>
        </div>
      </template>

      <el-table v-loading="loading" :data="pagedPropertyList" border stripe style="width: 100%">
        <el-table-column prop="property_id" label="属性标识" min-width="150" />
        <el-table-column prop="property_name" label="显示名称" min-width="150" />
        <el-table-column prop="data_type" label="数据类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ row.data_type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="property_type" label="属性类型" width="120">
          <template #default="{ row }">
            {{ row.property_type === 0 ? '普通' : '其他' }}
          </template>
        </el-table-column>
        <el-table-column prop="create_time" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.create_time) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确定删除该属性吗？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button link type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper" :total="filteredPropertyList.length" @size-change="handleSizeChange"
          @current-change="handleCurrentChange" />
      </div>
    </el-card>

    <!-- 编辑/添加弹窗 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑属性' : '添加属性'" width="500px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="属性标识" prop="property_id">
          <el-input v-model="form.property_id" :disabled="isEdit" placeholder="例如：e_os" />
        </el-form-item>
        <el-form-item label="显示名称" prop="property_name">
          <el-input v-model="form.property_name" placeholder="例如：操作系统" />
        </el-form-item>
        <el-form-item label="数据类型" prop="data_type">
          <el-select v-model="form.data_type" placeholder="选择数据类型" style="width: 100%">
            <el-option label="字符串 (string)" value="string" />
            <el-option label="数值 (number)" value="number" />
            <el-option label="日期时间 (datetime)" value="datetime" />
            <el-option label="布尔值 (boolean)" value="boolean" />
          </el-select>
        </el-form-item>
        <el-form-item label="属性类型" prop="property_type">
          <el-radio-group v-model="form.property_type">
            <el-radio :value="0">普通</el-radio>
            <el-radio :value="1">其他</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import type { ProjectProperty } from '@/api/meta'
import { createProjectProperty, deleteProjectProperty, getProjectProperties, updateProjectProperty } from '@/api/meta'
import { useAppStore } from '@/store/app'
import { Search } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'

const appStore = useAppStore()
const projectAlias = ref(appStore.activeProjectAlias)
const loading = ref(false)
const propertyList = ref<ProjectProperty[]>([])

const currentPage = ref(1)
const pageSize = ref(20)

const searchQuery = ref('')

const filteredPropertyList = computed(() => {
  if (!searchQuery.value) {
    return propertyList.value
  }
  const query = searchQuery.value.toLowerCase()
  return propertyList.value.filter(
    (property) =>
      property.property_id.toLowerCase().includes(query) ||
      (property.property_name && property.property_name.toLowerCase().includes(query))
  )
})

const pagedPropertyList = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredPropertyList.value.slice(start, end)
})

const handleSizeChange = (val: number) => {
  pageSize.value = val
  currentPage.value = 1
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
}

const dialogVisible = ref(false)
const isEdit = ref(false)
const submitLoading = ref(false)
const formRef = ref<FormInstance>()

const form = reactive<ProjectProperty>({
  project_alias: projectAlias.value,
  property_id: '',
  property_name: '',
  data_type: 'string',
  property_type: 0
})

const rules = reactive<FormRules>({
  property_id: [{ required: true, message: '请输入属性标识', trigger: 'blur' }],
  property_name: [{ required: true, message: '请输入显示名称', trigger: 'blur' }],
  data_type: [{ required: true, message: '请选择数据类型', trigger: 'change' }]
})

const fetchList = async () => {
  if (!projectAlias.value) return
  loading.value = true
  try {
    const data = await getProjectProperties(projectAlias.value)
    propertyList.value = data as unknown as ProjectProperty[]
  } catch (error: any) {
    ElMessage.error('获取列表失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  isEdit.value = false
  form.project_alias = projectAlias.value
  form.property_id = ''
  form.property_name = ''
  form.data_type = 'string'
  form.property_type = 0
  dialogVisible.value = true
}

const handleEdit = (row: ProjectProperty) => {
  isEdit.value = true
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleDelete = async (row: ProjectProperty) => {
  try {
    await deleteProjectProperty(projectAlias.value, row.property_id)
    ElMessage.success('删除成功')
    fetchList()
  } catch (error: any) {
    ElMessage.error('删除失败: ' + error.message)
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        if (isEdit.value) {
          await updateProjectProperty(form)
          ElMessage.success('更新成功')
        } else {
          await createProjectProperty(form)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        fetchList()
      } catch (error: any) {
        ElMessage.error('操作失败: ' + error.message)
      } finally {
        submitLoading.value = false
      }
    }
  })
}

const formatDate = (dateStr?: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString()
}

onMounted(() => {
  fetchList()
})

</script>

<style scoped lang="scss">
.event-property-container {
  padding: 0;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .header-left {
      display: flex;
      align-items: center;

      .title {
        font-size: 16px;
        font-weight: bold;
      }
    }
  }

  /* 移除外层 card 的底部间距，因为容器已经有 padding */
  :deep(.el-card) {
    margin-bottom: 0;
  }

  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>

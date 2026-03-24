<template>
  <div class="project-event-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span class="title">元事件管理</span>
            <el-tag size="small" type="info" style="margin-left: 10px">{{ projectAlias }}</el-tag>
          </div>
          <div class="header-right">
            <el-input
              v-model="searchQuery"
              placeholder="搜索事件标识或名称"
              prefix-icon="Search"
              clearable
              style="width: 250px; margin-right: 15px"
            />
            <el-button type="primary" icon="Plus" @click="handleAdd">添加事件</el-button>
          </div>
        </div>
      </template>

      <el-table v-loading="loading" :data="pagedEventList" border stripe style="width: 100%">
        <el-table-column prop="event_id" label="事件标识" min-width="150" />
        <el-table-column prop="event_name" label="显示名称" min-width="150" />
        <el-table-column prop="event_type" label="事件类型" width="120">
          <template #default="{ row }">
            <el-tag size="small" :type="row.event_type === 0 ? '' : 'warning'">
              {{ row.event_type === 0 ? '普通事件' : '虚拟事件' }}
            </el-tag>
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
            <el-popconfirm title="确定删除该事件吗？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button link type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper" :total="filteredEventList.length" @size-change="handleSizeChange"
          @current-change="handleCurrentChange" />
      </div>
    </el-card>

    <!-- 编辑/添加弹窗 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑事件' : '添加事件'" width="500px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="事件标识" prop="event_id">
          <el-input v-model="form.event_id" :disabled="isEdit" placeholder="例如：login" />
        </el-form-item>
        <el-form-item label="显示名称" prop="event_name">
          <el-input v-model="form.event_name" placeholder="例如：登录" />
        </el-form-item>
        <el-form-item label="事件类型" prop="event_type">
          <el-radio-group v-model="form.event_type">
            <el-radio :value="0">普通事件</el-radio>
            <el-radio :value="1">虚拟事件</el-radio>
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
import type { ProjectEvent } from '@/api/meta'
import { createProjectEvent, deleteProjectEvent, getProjectEvents, updateProjectEvent } from '@/api/meta'
import { useAppStore } from '@/store/app'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'

const appStore = useAppStore()
const projectAlias = ref(appStore.activeProjectAlias)
const loading = ref(false)
const eventList = ref<ProjectEvent[]>([])
const searchQuery = ref('')

const currentPage = ref(1)
const pageSize = ref(20)

const filteredEventList = computed(() => {
  if (!searchQuery.value) {
    return eventList.value
  }
  const query = searchQuery.value.toLowerCase()
  return eventList.value.filter(
    (event) =>
      event.event_id.toLowerCase().includes(query) ||
      (event.event_name && event.event_name.toLowerCase().includes(query))
  )
})

const pagedEventList = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredEventList.value.slice(start, end)
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

const form = reactive<ProjectEvent>({
  project_alias: projectAlias.value,
  event_id: '',
  event_name: '',
  event_type: 0
})

const rules = reactive<FormRules>({
  event_id: [{ required: true, message: '请输入事件标识', trigger: 'blur' }],
  event_name: [{ required: true, message: '请输入显示名称', trigger: 'blur' }],
  event_type: [{ required: true, message: '请选择事件类型', trigger: 'change' }]
})

const fetchList = async () => {
  if (!projectAlias.value) return
  loading.value = true
  try {
    const data = await getProjectEvents(projectAlias.value)
    eventList.value = data as unknown as ProjectEvent[]
  } catch (error: any) {
    ElMessage.error('获取列表失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  isEdit.value = false
  form.project_alias = projectAlias.value
  form.event_id = ''
  form.event_name = ''
  form.event_type = 0
  dialogVisible.value = true
}

const handleEdit = (row: ProjectEvent) => {
  isEdit.value = true
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleDelete = async (row: ProjectEvent) => {
  try {
    await deleteProjectEvent(projectAlias.value, row.event_id)
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
          await updateProjectEvent(form)
          ElMessage.success('更新成功')
        } else {
          await createProjectEvent(form)
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
.project-event-container {
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

    .header-right {
      display: flex;
      align-items: center;
    }
  }

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

<template>
  <div class="project-management-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span class="title">项目管理</span>
          <el-button type="primary" icon="Plus" @click="handleAdd">新建项目</el-button>
        </div>
      </template>

      <el-table v-loading="loading" :data="projectList" border stripe style="width: 100%">
        <el-table-column prop="project_alias" label="项目标识" min-width="120" />
        <el-table-column prop="project_name" label="项目名称" min-width="150" />
        <el-table-column prop="region" label="地域/设备号" min-width="120" />
        <el-table-column prop="secret" label="签名密钥" min-width="180">
          <template #default="{ row }">
            <code>{{ row.secret }}</code>
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
            <el-popconfirm title="确定删除该项目吗？这可能导致关联数据失效！" @confirm="handleDelete(row)">
              <template #reference>
                <el-button link type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 编辑/添加弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑项目' : '新建项目'"
      width="550px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="项目标识" prop="project_alias">
          <el-input v-model="form.project_alias" :disabled="isEdit" placeholder="唯一英文标识，例如：zgmgr4" />
        </el-form-item>
        <el-form-item label="项目名称" prop="project_name">
          <el-input v-model="form.project_name" placeholder="项目的显示名称" />
        </el-form-item>
        <el-form-item label="地域标识" prop="region">
          <el-input v-model="form.region" placeholder="例如：cn-shanghai" />
        </el-form-item>
        <el-form-item label="签名密钥" prop="secret">
          <el-input v-model="form.secret" placeholder="用于 API 鉴权的密钥">
            <template #append>
              <el-button @click="generateSecret">随机生成</el-button>
            </template>
          </el-input>
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
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { getProjects, createProject, updateProject, deleteProject } from '@/api/project'
import type { ProjectData } from '@/api/project'

const loading = ref(false)
const projectList = ref<ProjectData[]>([])

const dialogVisible = ref(false)
const isEdit = ref(false)
const submitLoading = ref(false)
const formRef = ref<FormInstance>()

const form = reactive<ProjectData>({
  project_alias: '',
  project_name: '',
  region: '',
  secret: ''
})

const rules = reactive<FormRules>({
  project_alias: [
    { required: true, message: '请输入项目标识', trigger: 'blur' },
    { pattern: /^[a-z0-9_]+$/, message: '仅支持小写字母、数字和下划线', trigger: 'blur' }
  ],
  project_name: [{ required: true, message: '请输入项目名称', trigger: 'blur' }],
  secret: [{ required: true, message: '请输入签名密钥', trigger: 'blur' }]
})

const fetchList = async () => {
  loading.value = true
  try {
    const data = await getProjects()
    projectList.value = data as unknown as ProjectData[]
  } catch (error: any) {
    ElMessage.error('获取列表失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  isEdit.value = false
  form.project_alias = ''
  form.project_name = ''
  form.region = ''
  form.secret = ''
  dialogVisible.value = true
}

const handleEdit = (row: ProjectData) => {
  isEdit.value = true
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleDelete = async (row: ProjectData) => {
  try {
    await deleteProject(row.project_alias)
    ElMessage.success('项目已删除')
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
          await updateProject(form)
          ElMessage.success('更新成功')
        } else {
          await createProject(form)
          ElMessage.success('项目创建成功')
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

const generateSecret = () => {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let result = ''
  for (let i = 0; i < 32; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  form.secret = result
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
.project-management-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    .title {
      font-size: 16px;
      font-weight: bold;
    }
  }
  
  code {
    background-color: var(--el-fill-color-light);
    padding: 2px 4px;
    border-radius: 4px;
    font-family: monospace;
  }
}
</style>

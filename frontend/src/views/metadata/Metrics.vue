<template>
  <div class="metrics-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span class="title">指标中心</span>
          <div class="header-right" style="display: flex; align-items: center;">
            <el-input v-model="searchQuery" placeholder="搜索指标名称/标识" clearable style="width: 250px; margin-right: 16px">
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-button type="primary" icon="Plus" @click="handleAdd">新建指标</el-button>
          </div>
        </div>
      </template>

      <el-table v-loading="loading" :data="filteredMetricList" border stripe style="width: 100%">
        <el-table-column prop="metric_name" label="指标名称" width="180" />
        <el-table-column prop="metric_code" label="指标标识" width="180">
          <template #default="{ row }">
            <code>{{ row.metric_code }}</code>
          </template>
        </el-table-column>
        <el-table-column prop="base_table" label="基础表" width="120" />
        <el-table-column prop="expression" label="计算表达式" show-overflow-tooltip>
          <template #default="{ row }">
            <code>{{ row.expression }}</code>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确定删除该指标吗？" @confirm="handleDelete(row)">
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
      :title="isEdit ? '编辑指标' : '新建指标'"
      width="680px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="110px">
        <el-form-item label="指标名称" prop="metric_name">
          <el-input v-model="form.metric_name" placeholder="如：DAU" />
        </el-form-item>
        <el-form-item label="指标标识" prop="metric_code">
          <el-input v-model="form.metric_code" placeholder="如：dau_count (只能包含字母、数字和下划线)" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="基础表" prop="base_table">
          <el-radio-group v-model="form.base_table">
            <el-radio-button label="event_data">事件表 (event_data)</el-radio-button>
            <el-radio-button label="user_data">用户表 (user_data)</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="计算表达式" prop="expression">
          <el-input 
            v-model="form.expression" 
            type="textarea" 
            :rows="4"
            placeholder="如：COUNT(DISTINCT CASE WHEN e_event_id='sys.login' THEN e_openid END)" 
          />
          <div class="form-tip">完整的聚合 SQL 片段，系统将在查询时直接替换。</div>
        </el-form-item>
        <el-form-item label="业务说明" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="2" placeholder="可选，用于说明该指标的业务含义" />
        </el-form-item>
        <el-form-item label="是否启用">
          <el-switch v-model="form.status" :active-value="1" :inactive-value="0" />
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
import {
  createMetric,
  deleteMetric,
  getMetrics,
  updateMetric,
  type ProjectMetric
} from '@/api/metric'
import { Search } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'

const loading = ref(false)
const metricList = ref<ProjectMetric[]>([])
const searchQuery = ref('')

const filteredMetricList = computed(() => {
  if (!searchQuery.value) {
    return metricList.value
  }
  const query = searchQuery.value.toLowerCase()
  return metricList.value.filter(
    (metric) =>
      (metric.metric_name && metric.metric_name.toLowerCase().includes(query)) ||
      (metric.metric_code && metric.metric_code.toLowerCase().includes(query))
  )
})

const dialogVisible = ref(false)
const isEdit = ref(false)
const submitLoading = ref(false)
const formRef = ref<FormInstance>()

const defaultForm = (): ProjectMetric => ({
  metric_name: '',
  metric_code: '',
  expression: '',
  base_table: 'event_data',
  description: '',
  status: 1
})

const form = reactive<ProjectMetric>(defaultForm())

const rules = reactive<FormRules>({
  metric_name: [{ required: true, message: '请输入指标名称', trigger: 'blur' }],
  metric_code: [
    { required: true, message: '请输入指标标识', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '只能包含字母、数字和下划线', trigger: 'blur' }
  ],
  base_table: [{ required: true, message: '请选择基础表', trigger: 'change' }],
  expression: [{ required: true, message: '请输入计算表达式', trigger: 'blur' }]
})

const fetchList = async () => {
  loading.value = true
  try {
    const data = await getMetrics()
    metricList.value = data as unknown as ProjectMetric[]
  } catch (error: any) {
    // 错误已由拦截器处理
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  isEdit.value = false
  Object.assign(form, defaultForm())
  dialogVisible.value = true
}

const handleEdit = (row: ProjectMetric) => {
  isEdit.value = true
  Object.assign(form, { ...row })
  dialogVisible.value = true
}

const handleDelete = async (row: ProjectMetric) => {
  try {
    await deleteMetric(row.metric_code)
    ElMessage.success('删除成功')
    fetchList()
  } catch (error: any) {
    // 错误已由拦截器处理
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        if (isEdit.value) {
          await updateMetric(form.metric_code, form)
          ElMessage.success('更新成功')
        } else {
          await createMetric(form)
          ElMessage.success('新建成功')
        }
        dialogVisible.value = false
        fetchList()
      } catch (error: any) {
        // 错误已由拦截器处理
      } finally {
        submitLoading.value = false
      }
    }
  })
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.metrics-container {
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
    padding: 2px 6px;
    border-radius: 4px;
    font-family: monospace;
    font-size: 12px;
  }
  .form-tip {
    font-size: 12px;
    color: var(--el-text-color-secondary);
    margin-top: 4px;
    line-height: 1.4;
  }
}
</style>

<template>
  <div class="ai-model-management-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span class="title">AI 模型配置</span>
          <el-button type="primary" icon="Plus" @click="handleAdd">添加模型</el-button>
        </div>
      </template>

      <el-table v-loading="loading" :data="modelList" border stripe style="width: 100%">
        <el-table-column prop="provider" label="提供商" width="120">
          <template #default="{ row }">
            <el-tag :type="providerTagType(row.provider)" size="small">{{ providerLabel(row.provider) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="display_name" label="显示名称" min-width="150" />
        <el-table-column prop="model_name" label="模型标识" min-width="150">
          <template #default="{ row }">
            <code>{{ row.model_name }}</code>
          </template>
        </el-table-column>
        <el-table-column prop="base_url" label="API 地址" min-width="250" show-overflow-tooltip />
        <el-table-column prop="api_key" label="API Key" width="150">
          <template #default="{ row }">
            <code>{{ row.api_key }}</code>
          </template>
        </el-table-column>
        <el-table-column label="默认" width="80" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.is_default" type="success" size="small">默认</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.is_enabled ? 'success' : 'info'" size="small">
              {{ row.is_enabled ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleTest(row)">测试</el-button>
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button v-if="!row.is_default" link type="success" @click="handleSetDefault(row)">设为默认</el-button>
            <el-popconfirm title="确定删除该模型配置吗？" @confirm="handleDelete(row)">
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
      :title="isEdit ? '编辑模型配置' : '添加模型配置'"
      width="620px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="110px">
        <el-form-item label="提供商" prop="provider">
          <el-select v-model="form.provider" placeholder="选择提供商" @change="onProviderChange">
            <el-option v-for="p in providers" :key="p.value" :label="p.label" :value="p.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="显示名称" prop="display_name">
          <el-input v-model="form.display_name" placeholder="如：DeepSeek-V3" />
        </el-form-item>
        <el-form-item label="API 地址" prop="base_url">
          <el-input v-model="form.base_url" placeholder="https://api.deepseek.com/v1/chat/completions" />
        </el-form-item>
        <el-form-item label="API Key" prop="api_key">
          <el-input v-model="form.api_key" placeholder="sk-xxxxxxxx" show-password />
        </el-form-item>
        <el-form-item label="模型标识" prop="model_name">
          <el-input v-model="form.model_name" placeholder="如：deepseek-chat" />
        </el-form-item>
        <el-form-item label="最大 Token">
          <el-input-number v-model="form.max_tokens" :min="256" :max="128000" :step="1024" />
        </el-form-item>
        <el-form-item label="温度">
          <el-slider v-model="form.temperature" :min="0" :max="2" :step="0.1" show-input style="width: 100%" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.is_enabled" />
        </el-form-item>
        <el-form-item label="设为默认">
          <el-switch v-model="form.is_default" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort_order" :min="0" :max="999" />
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
import {
  getAIModels,
  createAIModel,
  updateAIModel,
  deleteAIModel,
  testAIModel
} from '@/api/aiModel'
import type { AIModelConfig } from '@/api/aiModel'

const providers = [
  { value: 'deepseek', label: 'DeepSeek', baseUrl: 'https://api.deepseek.com/v1/chat/completions' },
  { value: 'openai', label: 'OpenAI', baseUrl: 'https://api.openai.com/v1/chat/completions' },
  { value: 'qwen', label: '通义千问', baseUrl: 'https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions' },
  { value: 'glm', label: '智谱GLM', baseUrl: 'https://open.bigmodel.cn/api/paas/v4/chat/completions' },
  { value: 'kimi', label: '月之暗面', baseUrl: 'https://api.moonshot.cn/v1/chat/completions' },
  { value: 'ollama', label: 'Ollama (本地)', baseUrl: 'http://localhost:11434/v1/chat/completions' }
]

const providerLabel = (val: string) => providers.find(p => p.value === val)?.label || val
const providerTagType = (val: string): string => {
  const map: Record<string, string> = {
    deepseek: '', openai: 'success', qwen: 'warning', glm: 'info', kimi: 'danger', ollama: ''
  }
  return map[val] || 'info'
}

const loading = ref(false)
const modelList = ref<AIModelConfig[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitLoading = ref(false)
const formRef = ref<FormInstance>()

const defaultForm = (): AIModelConfig => ({
  provider: 'deepseek',
  display_name: '',
  base_url: 'https://api.deepseek.com/v1/chat/completions',
  api_key: '',
  model_name: '',
  max_tokens: 4096,
  temperature: 0.1,
  is_default: false,
  is_enabled: true,
  sort_order: 0
})

const form = reactive<AIModelConfig>(defaultForm())

const rules = reactive<FormRules>({
  provider: [{ required: true, message: '请选择提供商', trigger: 'change' }],
  display_name: [{ required: true, message: '请输入显示名称', trigger: 'blur' }],
  base_url: [{ required: true, message: '请输入 API 地址', trigger: 'blur' }],
  model_name: [{ required: true, message: '请输入模型标识', trigger: 'blur' }]
})

const fetchList = async () => {
  loading.value = true
  try {
    const data = await getAIModels()
    modelList.value = data as unknown as AIModelConfig[]
  } catch (error: any) {
    ElMessage.error('获取列表失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const onProviderChange = (val: string) => {
  const p = providers.find(p => p.value === val)
  if (p) {
    form.base_url = p.baseUrl
  }
}

const handleAdd = () => {
  isEdit.value = false
  Object.assign(form, defaultForm())
  dialogVisible.value = true
}

const handleEdit = (row: AIModelConfig) => {
  isEdit.value = true
  Object.assign(form, { ...row, api_key: '' })
  dialogVisible.value = true
}

const handleDelete = async (row: AIModelConfig) => {
  try {
    await deleteAIModel(row.id!)
    ElMessage.success('已删除')
    fetchList()
  } catch (error: any) {
    ElMessage.error('删除失败: ' + error.message)
  }
}

const handleSetDefault = async (row: AIModelConfig) => {
  try {
    await updateAIModel(row.id!, { ...row, is_default: true })
    ElMessage.success('已设为默认')
    fetchList()
  } catch (error: any) {
    ElMessage.error('设置失败: ' + error.message)
  }
}

const handleTest = async (row: AIModelConfig) => {
  try {
    ElMessage.info('正在测试连接...')
    const res = await testAIModel(row.id!) as any
    ElMessage.success(res?.message || '连接成功')
  } catch (error: any) {
    ElMessage.error('测试失败: ' + error.message)
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        if (isEdit.value) {
          await updateAIModel(form.id!, form)
          ElMessage.success('更新成功')
        } else {
          await createAIModel(form)
          ElMessage.success('添加成功')
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

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.ai-model-management-container {
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
    font-size: 12px;
  }
}
</style>

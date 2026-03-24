<template>
  <el-dialog
    v-model="visible"
    title="保存报表"
    width="480px"
    destroy-on-close
  >
    <el-form :model="form" label-width="80px" ref="formRef" :rules="rules">
      <el-form-item label="报表名称" prop="name">
        <el-input v-model="form.name" placeholder="请输入报表名称" />
      </el-form-item>
      <el-form-item label="报表描述">
        <el-input
          v-model="form.description"
          type="textarea"
          placeholder="请输入报表描述"
          :rows="3"
        />
      </el-form-item>
      
      <el-divider content-position="left">添加到看板 (可选)</el-divider>
      
      <el-form-item label="选择看板">
        <el-select
          v-model="form.dashboardId"
          placeholder="请选择看板"
          clearable
          style="width: 100%"
          :loading="loadingDashboards"
        >
          <el-option
            v-for="db in dashboards"
            :key="db.id"
            :label="db.display_name || db.name"
            :value="db.id!"
          />
        </el-select>
        <div class="tip">保存后将自动把报表添加至选中的看板中</div>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" @click="handleConfirm" :loading="loading">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch, onMounted } from 'vue'
import { getDashboardList } from '@/api/dashboard'
import type { Dashboard } from '@/api/dashboard'
import { useAppStore } from '@/store/app'

const props = defineProps<{
  modelValue: boolean
  loading?: boolean
  initialData?: { name: string; description: string; dashboardId?: number }
}>()

const emit = defineEmits(['update:modelValue', 'confirm'])

const appStore = useAppStore()
const formRef = ref()
const loadingDashboards = ref(false)
const dashboards = ref<Dashboard[]>([])

const visible = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v)
})

const form = reactive({
  name: '',
  description: '',
  dashboardId: undefined as number | undefined
})

const rules = {
  name: [{ required: true, message: '请输入报表名称', trigger: 'blur' }]
}

const fetchDashboards = async () => {
  if (!appStore.activeProjectAlias) return
  loadingDashboards.value = true
  try {
    const res = await getDashboardList({ project_alias: appStore.activeProjectAlias })
    dashboards.value = res as unknown as Dashboard[]
  } catch (error) {
    console.error('Failed to fetch dashboards:', error)
  } finally {
    loadingDashboards.value = false
  }
}

watch(() => props.initialData, (newVal) => {
  if (newVal) {
    form.name = newVal.name || ''
    form.description = newVal.description || ''
    form.dashboardId = newVal.dashboardId
  }
}, { immediate: true })

watch(visible, (val) => {
  if (val) {
    fetchDashboards()
  }
})

const handleConfirm = async () => {
  if (!formRef.value) return
  await formRef.value.validate((valid: boolean) => {
    if (valid) {
      emit('confirm', { ...form })
    }
  })
}
</script>

<style scoped>
.tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
  line-height: 1.4;
}
:deep(.el-divider__text) {
  font-size: 13px;
  color: #606266;
  font-weight: normal;
}
</style>

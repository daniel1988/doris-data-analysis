<template>
  <div class="property-analysis-container">
    <!-- 页面头部 -->
    <div class="page-header mb-15">
      <div class="header-left">
        <span class="page-title">用户属性分析</span>
        <el-tag v-if="currentReport" type="success" closable @close="handleReset" class="ml-10">
          {{ currentReport.name }}
        </el-tag>
      </div>
      <div class="header-right">
        <el-button icon="Files" @click="showReportDrawer = true">我的报表</el-button>
      </div>
    </div>

    <!-- 主布局 -->
    <div class="analysis-content">
      <PropertyAnalysisLayout>
        <!-- 左侧配置区 -->
        <template #config>
          <!-- 指标配置 -->
          <UserPropertyMetricCard :metric="form.metric" />

          <!-- 全局过滤 -->
          <GlobalFiltersCard :filter-group="form.filter_groups" />

          <!-- 分组/人群分群 -->
          <el-card class="group-config-card" shadow="never">
            <template #header>
              <div class="card-header">
                <el-radio-group v-model="form.group_type" size="small">
                  <el-radio-button :label="1">维度分组</el-radio-button>
                  <el-radio-button :label="2">人群分群对比</el-radio-button>
                </el-radio-group>
              </div>
            </template>

            <div class="group-config-content">
              <!-- 维度分组模式 -->
              <div v-if="form.group_type === 1">
                <GlobalGroupsCard :hide-header="true" />
              </div>

              <!-- 人群分群模式 -->
              <div v-else>
                <PropertyUserGroupsPanel v-model:user-groups="form.user_groups" />
              </div>
            </div>
          </el-card>
        </template>

        <!-- 底部操作栏 -->
        <template #footer>
          <ActionFooter :loading="loading" @analyze="handleAnalyze" @save="showSaveDialog = true"
            @reset="handleReset" />
        </template>

        <!-- 右侧结果展示区 -->
        <template #header>
          <!-- 时间与粒度控制 -->
          <TimeControls hide-interval v-model:static-range="staticTimeRange"
            v-model:comparison-range="comparisonTimeRange" @dynamic-change="handleDynamicTimeChange" />
        </template>

        <template #results>
          <div class="result-area">
            <!-- 结果图表与表格 -->
            <ResultView :loading="loading" :rows="results" :columns="columns" :viz-mode="vizMode" :sql="sql"
              @update:viz-mode="vizMode = $event" @show-sql="showSqlDrawer = true" />
          </div>
        </template>
      </PropertyAnalysisLayout>
    </div>

    <!-- SQL 预览抽屉 -->
    <SqlDrawer v-model="showSqlDrawer" :sql="sql" />

    <!-- 保存报表弹窗 -->
    <SaveReportDialog v-model="showSaveDialog" :initial-data="currentReport || { name: '', description: '' }"
      @confirm="onSaveConfirm" />

    <!-- 报表列表抽屉 -->
    <ReportDrawer v-model="showReportDrawer" :reports="reportList" :loading="reportLoading" @load="onLoadReport"
      @delete="removeReport" />
  </div>
</template>

<script setup lang="ts">
import { getReportById } from '@/api/report'
import { VizModes } from '@/constants/analysis'
import { useAppStore } from '@/store/app'
import { ElMessage } from 'element-plus'
import { onMounted, provide, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ANALYSIS_CONTEXT_KEY, type AnalysisContext } from './context'

// 组件导入
import ActionFooter from './components/ActionFooter.vue'
import GlobalFiltersCard from './components/GlobalFiltersCard.vue'
import GlobalGroupsCard from './components/GlobalGroupsCard.vue'
import PropertyAnalysisLayout from './components/PropertyAnalysisLayout.vue'
import PropertyUserGroupsPanel from './components/PropertyUserGroupsPanel.vue'
import ReportDrawer from './components/ReportDrawer.vue'
import ResultView from './components/ResultView.vue'
import SaveReportDialog from './components/SaveReportDialog.vue'
import SqlDrawer from './components/SqlDrawer.vue'
import TimeControls from './components/TimeControls.vue'
import UserPropertyMetricCard from './components/UserPropertyMetricCard.vue'

// Composables
import { useEventMeta } from './composables/useEventMeta'
import { usePropertyAnalysisForm } from './composables/usePropertyAnalysisForm'
import { useReportPersistence } from './composables/useReportPersistence'
import { useTimeRangeSync } from './composables/useTimeRangeSync'

import { usePayloadBuilders } from './composables/usePayloadBuilders'

const appStore = useAppStore()
const route = useRoute()
const {
  form,
  loading,
  results,
  columns,
  total,
  sql,
  runPropertyAnalysis,
  resetForm,
  loadForm
} = usePropertyAnalysisForm()

const { eventOptions, propertyOptions, fetchMeta } = useEventMeta()
const { buildPropertyPayload } = usePayloadBuilders()
const {
  staticTimeRange,
  comparisonTimeRange,
  dynamicDateFilter,
  handleDynamicTimeChange,
  syncFromForm
} = useTimeRangeSync()

const {
  currentReport,
  reportList,
  loading: reportLoading,
  saveReport,
  fetchReportList,
  removeReport
} = useReportPersistence()

// UI 状态
const vizMode = ref<VizModes>(VizModes.Bar)
const showSqlDrawer = ref(false)
const showSaveDialog = ref(false)
const showReportDrawer = ref(false)

// --- Context Setup ---
const contextState = reactive({
  projectAlias: appStore.activeProjectAlias,
  metadata: {
    eventOptions: eventOptions,
    propertyOptions: propertyOptions
  },
  loading,
  form,
  results,
  columns,
  sql,
  total,
  vizMode
})

const handleAnalyze = async () => {
  if (!form.metric.metric.column.field) {
    ElMessage.warning('请选择分析属性')
    return
  }

  form.project_alias = appStore.activeProjectAlias
  form.filter_groups.query_dates = staticTimeRange.value || []

  const payload = buildPropertyPayload(
    form.project_alias,
    form.metric,
    form.filter_groups,
    form.groups,
    form.user_groups,
    form.group_type,
    dynamicDateFilter.value
  )

  await runPropertyAnalysis(payload)
}

const contextActions = {
  refreshMetadata: fetchMeta,
  triggerAnalyze: handleAnalyze,
  syncTimeRange: (staticRange: string[] | null, comparisonRange: string[] | null) => {
    staticTimeRange.value = staticRange as [string, string] | null
    comparisonTimeRange.value = comparisonRange as [string, string] | null
  }
}

provide(ANALYSIS_CONTEXT_KEY, {
  state: contextState,
  actions: contextActions
} as AnalysisContext)

const handleReset = () => {
  resetForm()
  currentReport.value = null
}

const onSaveConfirm = async (data: { name: string; description: string; dashboardId?: number }) => {
  const success = await saveReport({
    ...data,
    category: 'property',
    queryParams: {
      form,
      vizMode: vizMode.value
    },
    projectAlias: appStore.activeProjectAlias
  })
  if (success) showSaveDialog.value = false
}

const onLoadReport = (report: any) => {
  const params = typeof report.query_params === 'string'
    ? JSON.parse(report.query_params)
    : report.query_params
  const formData = params.form || params
  loadForm(formData)
  const filterData = formData.filter_groups || formData.global_filter_groups
  if (filterData) {
    syncFromForm(filterData)
  }
  vizMode.value = params.vizMode || VizModes.Bar
  currentReport.value = report
  showReportDrawer.value = false
  handleAnalyze()
}

const loadReportFromRoute = async () => {
  const reportId = Number(route.query.reportId)
  if (!reportId) return
  try {
    const report = await getReportById(reportId)
    onLoadReport(report)
  } catch (e) {
    ElMessage.error('报表加载失败')
  }
}

// --- Lifecycle ---
onMounted(async () => {
  await fetchMeta()
  await fetchReportList(appStore.activeProjectAlias, 'property')
  await loadReportFromRoute()
})

watch(() => appStore.activeProjectAlias, (newVal) => {
  if (newVal) {
    fetchMeta()
    fetchReportList(newVal, 'property')
    handleReset()
  }
})

watch(() => route.query.reportId, () => {
  loadReportFromRoute()
})
</script>

<style scoped lang="scss">
.property-analysis-container {
  padding: 0;
  height: 100%;
  display: flex;
  flex-direction: column;

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .page-title {
      font-size: 18px;
      font-weight: bold;
    }
  }

  .analysis-content {
    flex: 1;
    min-height: 0;
  }

  .result-area {
    display: flex;
    flex-direction: column;
    gap: 16px;
    height: 100%;
  }
}

.group-config-card {
  border-radius: 12px;

  .card-header {
    display: flex;
    justify-content: center;
  }
}

.ml-10 {
  margin-left: 10px;
}
</style>

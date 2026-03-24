<template>
  <div class="funnel-analysis-container">
    <!-- 页面头部 -->
    <div class="page-header mb-15">
      <div class="header-left">
        <span class="page-title">漏斗分析</span>
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
      <FunnelAnalysisLayout>
        <!-- 左侧配置区 -->
        <template #config>
          <!-- 漏斗步骤配置 -->
          <FunnelStepsList v-model:window-num="form.time_grain.window_num" :steps="form.event_metrics"
            :project-events="eventOptions" :project-alias="appStore.activeProjectAlias" @add-step="addStep"
            @remove-step="removeStep" @toggle-filter="(step) => step.filterPopoverVisible = !step.filterPopoverVisible"
            @event-change="onEventChange" />

          <!-- 全局过滤 -->
          <GlobalFiltersCard />

          <!-- 分组配置 -->
          <GlobalGroupsCard />
        </template>

        <!-- 底部操作栏 -->
        <template #footer>
          <ActionFooter :loading="loading" @analyze="handleAnalyze" @save="showSaveDialog = true"
            @reset="handleReset" />
        </template>

        <!-- 右侧结果展示区 -->
        <template #header>
          <!-- 时间与粒度控制 -->
          <TimeControls v-model:interval="form.time_grain.interval" v-model:static-range="staticTimeRange"
            v-model:comparison-range="comparisonTimeRange" @dynamic-change="handleDynamicTimeChange" />
        </template>

        <template #results>
          <div class="result-area">
            <!-- 结果图表与表格 -->
            <ResultView :loading="loading" :rows="results" :columns="columns" :viz-mode="vizMode" :sql="sql"
              @update:viz-mode="vizMode = $event" @show-sql="showSqlDrawer = true" />
          </div>
        </template>
      </FunnelAnalysisLayout>
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
import FunnelStepsList from './components/funnel/FunnelStepsList.vue'
import FunnelAnalysisLayout from './components/FunnelAnalysisLayout.vue'
import GlobalFiltersCard from './components/GlobalFiltersCard.vue'
import GlobalGroupsCard from './components/GlobalGroupsCard.vue'
import ReportDrawer from './components/ReportDrawer.vue'
import ResultView from './components/ResultView.vue'
import SaveReportDialog from './components/SaveReportDialog.vue'
import SqlDrawer from './components/SqlDrawer.vue'
import TimeControls from './components/TimeControls.vue'

// Composables
import { useEventMeta } from './composables/useEventMeta'
import { useFunnelForm } from './composables/useFunnelForm'
import { usePayloadBuilders } from './composables/usePayloadBuilders'
import { useReportPersistence } from './composables/useReportPersistence'
import { useTimeRangeSync } from './composables/useTimeRangeSync'

const appStore = useAppStore()
const route = useRoute()
const {
  form,
  loading,
  results,
  columns,
  total,
  sql,
  addStep,
  removeStep,
  runFunnelAnalysis,
  resetForm,
  loadForm
} = useFunnelForm()

const { eventOptions, propertyOptions, fetchMeta } = useEventMeta()
const { buildFunnelPayload } = usePayloadBuilders()
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
const vizMode = ref<VizModes>(VizModes.Bar) // 漏斗分析默认用柱状图/漏斗图
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
  // 校验步骤是否完整
  const incompleteStep = form.event_metrics.find(s => !s.e_event_id)
  if (incompleteStep) {
    ElMessage.warning('请为所有步骤选择事件')
    return
  }

  form.project_alias = appStore.activeProjectAlias
  form.filter_groups.query_dates = staticTimeRange.value || []

  const payload = buildFunnelPayload(
    form.project_alias,
    form.event_metrics,
    form.filter_groups,
    form.time_grain,
    form.groups,
    dynamicDateFilter.value
  )

  await runFunnelAnalysis(payload)
}

const contextActions = {
  refreshMetadata: fetchMeta,
  triggerAnalyze: handleAnalyze,
  syncTimeRange: (staticRange: string[] | null, comparisonRange: string[] | null) => {
    staticTimeRange.value = staticRange
    comparisonTimeRange.value = comparisonRange
  }
}

provide(ANALYSIS_CONTEXT_KEY, {
  state: contextState,
  actions: contextActions
} as AnalysisContext)

const onEventChange = (step: any) => {
  const event = eventOptions.value.find(e => e.id === step.e_event_id)
  if (event) {
    step.name = event.name || event.id
  }
}

const handleReset = () => {
  resetForm()
  currentReport.value = null
  syncFromForm(form.filter_groups)
}

const onSaveConfirm = async (data: { name: string; description: string; dashboardId?: number }) => {
  const success = await saveReport({
    ...data,
    category: 'funnel',
    queryParams: {
      form,
      staticTimeRange: staticTimeRange.value,
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
  staticTimeRange.value = params.staticTimeRange || filterData?.query_dates || null
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
  await fetchReportList(appStore.activeProjectAlias, 'funnel')
  await loadReportFromRoute()
})

watch(() => appStore.activeProjectAlias, (newVal) => {
  if (newVal) {
    fetchMeta()
    fetchReportList(newVal, 'funnel')
    handleReset()
  }
})

watch(() => route.query.reportId, () => {
  loadReportFromRoute()
})
</script>

<style scoped lang="scss">
.funnel-analysis-container {
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
</style>

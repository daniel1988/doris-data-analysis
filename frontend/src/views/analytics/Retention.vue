<template>
  <div class="retention-analysis-container">
    <!-- 页面头部 -->
    <div class="page-header mb-15">
      <div class="header-left">
        <span class="page-title">留存分析</span>
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
      <RetentionAnalysisLayout>
        <!-- 左侧配置区 -->
        <template #config>
          <!-- 留存指标配置 -->
          <RetentionMetricsList :init-metric="form.init_event_metric" :end-metric="form.end_event_metric"
            :project-alias="appStore.activeProjectAlias" :get-display-property-name="getDisplayPropertyName"
            @event-change="onEventChange" @apply-formula="onApplyFormula" />

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
            :hide-comparison="true" @dynamic-change="handleDynamicTimeChange">
            <template #extra>
              <RetentionDaysSelector v-model="form.day_n_array" :interval="form.time_grain.interval" />
            </template>
          </TimeControls>
        </template>

        <template #results>
          <div class="result-area">
            <!-- 结果图表与表格 -->
            <ResultView :loading="loading" :rows="results" :columns="columns" :viz-mode="vizMode" :sql="sql"
              :hide-chart="true" @update:viz-mode="vizMode = $event" @show-sql="showSqlDrawer = true" />
          </div>
        </template>
      </RetentionAnalysisLayout>
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
import { FormulaDefaults, TableNames, VizModes } from '@/constants/analysis'
import { useAppStore } from '@/store/app'
import { MetricFormulaLabels } from '@/types/doris/constants'
import { ElMessage } from 'element-plus'
import { onMounted, provide, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ANALYSIS_CONTEXT_KEY, type AnalysisContext } from './context'

// 组件导入
import ActionFooter from './components/ActionFooter.vue'
import GlobalFiltersCard from './components/GlobalFiltersCard.vue'
import GlobalGroupsCard from './components/GlobalGroupsCard.vue'
import RetentionMetricsList from './components/metrics/RetentionMetricsList.vue'
import ReportDrawer from './components/ReportDrawer.vue'
import ResultView from './components/ResultView.vue'
import RetentionAnalysisLayout from './components/RetentionAnalysisLayout.vue'
import RetentionDaysSelector from './components/RetentionDaysSelector.vue'
import SaveReportDialog from './components/SaveReportDialog.vue'
import SqlDrawer from './components/SqlDrawer.vue'
import TimeControls from './components/TimeControls.vue'

// Composables
import { useEventMeta } from './composables/useEventMeta'
import { usePayloadBuilders } from './composables/usePayloadBuilders'
import { useReportPersistence } from './composables/useReportPersistence'
import { useRetentionForm } from './composables/useRetentionForm'
import { useTimeRangeSync } from './composables/useTimeRangeSync'
import { normalizeMetricField } from './utils/payloadNormalizer'

const appStore = useAppStore()
const route = useRoute()
const {
  form,
  loading,
  results,
  columns,
  total,
  sql,
  runRetentionAnalysis,
  resetForm,
  loadForm
} = useRetentionForm()

const { eventOptions, propertyOptions, fetchMeta } = useEventMeta()
const { buildPayload, buildRetentionPayload } = usePayloadBuilders()
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
const vizMode = ref<VizModes>(VizModes.Line)
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
  if (!form.init_event_metric.e_event_id || !form.end_event_metric.e_event_id) {
    ElMessage.warning('请选择初始事件和结束事件')
    return
  }

  form.project_alias = appStore.activeProjectAlias
  form.filter_groups.query_dates = staticTimeRange.value || []

  const payload = buildRetentionPayload(
    form.project_alias,
    form.init_event_metric,
    form.end_event_metric,
    form.filter_groups,
    form.time_grain,
    form.groups,
    form.day_n_array,
    dynamicDateFilter.value
  )

  await runRetentionAnalysis(payload)
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

// --- 辅助函数 ---
const getDisplayPropertyName = (m: any) => {
  if (m.metric?.column?.field === FormulaDefaults.TOTAL_TIMES_FIELD) return '总次数'
  if (m.metric?.column?.field === FormulaDefaults.TOTAL_USERS_FIELD) return '总用户数'
  const prop = propertyOptions.value.find(p => p.id === m.metric?.column?.field)
  return prop?.name || m.metric?.column?.field || ''
}

const getFormulaName = (formula: number) => {
  return MetricFormulaLabels[formula] || ''
}

const onApplyFormula = (payload: any) => {
  const { m, field, table, formula, propertyName } = payload
  if (!m.metric) m.metric = { column: { table: TableNames.EVENT, field: '' }, formula: 2, format: 'int' }
  
  if (['__TOTAL_TIMES__', '__TOTAL_USERS__'].includes(field)) {
    normalizeMetricField(m, field)
  } else {
    m.metric.column.field = field
    m.metric.column.table = table
    m.metric.formula = formula
  }

  const eventName = eventOptions.value.find(e => e.id === m.e_event_id)?.name || m.e_event_id
  const propName = propertyName || field
  const formulaName = getFormulaName(m.metric.formula)
  m.name = `${eventName}.${propName} · ${formulaName}`
}

const onEventChange = (m: any) => {
  if (m.metric && m.metric.column) {
    normalizeMetricField(m, '__TOTAL_USERS__')
  }
  m.name = eventOptions.value.find(e => e.id === m.e_event_id)?.name || m.e_event_id
}

const handleReset = () => {
  resetForm()
  currentReport.value = null
  syncFromForm(form.filter_groups)
}

const onSaveConfirm = async (data: { name: string; description: string; dashboardId?: number }) => {
  const success = await saveReport({
    ...data,
    category: 'retention',
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
  vizMode.value = params.vizMode || VizModes.Line
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
  await fetchReportList(appStore.activeProjectAlias, 'retention')
  await loadReportFromRoute()
})

watch(() => appStore.activeProjectAlias, (newVal) => {
  if (newVal) {
    fetchMeta()
    fetchReportList(newVal, 'retention')
    handleReset()
  }
})

watch(() => route.query.reportId, () => {
  loadReportFromRoute()
})
</script>

<style scoped lang="scss">
.retention-analysis-container {
  padding: 0;
  height: calc(100vh - 150px); // 减去导航栏、面包屑/标签页以及 content-wrapper 的 padding
  display: flex;
  flex-direction: column;
  overflow: hidden;

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

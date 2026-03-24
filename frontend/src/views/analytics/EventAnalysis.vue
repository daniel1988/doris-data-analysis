<template>
  <div class="event-analysis-container">
    <div class="page-header mb-15">
      <div class="header-left">
        <span class="page-title">事件分析</span>
        <el-tag v-if="currentReport" type="success" closable @close="handleReset" class="ml-10">
          {{ currentReport.name }}
        </el-tag>
      </div>
      <div class="header-right">
        <el-button icon="Files" @click="openReportDrawer">我的报表</el-button>
      </div>
    </div>

    <div class="analysis-content">
      <EventAnalysisLayout>
        <!-- 左侧配置区 -->
        <template #config>
          <MetricsList :get-display-property-name="getDisplayPropertyName" :get-formula-name="getFormulaName"
            @add-metric="addMetric" @add-custom="addCustomMetric" @remove-metric="removeMetric"
            @edit-name="(m: any) => m.isEditingName = true" @finish-edit-name="(m: any) => m.isEditingName = false"
            @toggle-formula="onToggleFormula" @copy-metric="onCopyMetric" @toggle-filter="onToggleFilter"
            @event-change="onEventChange" @apply-formula="(payload: any) => onApplyFormula(payload.m, payload)"
            @update-custom-elements="(payload: any) => onUpdateCustomElements(payload.m, payload.elements)"
            @request-new-event="onRequestNewEvent" @open-custom-filter="onOpenCustomFilter" @reorder="onReorder" />

          <GlobalFiltersCard />

          <GlobalGroupsCard />
        </template>

        <template #footer>
          <ActionFooter :loading="loading" @analyze="handleAnalyze" @save="showSaveDialog = true"
            @reset="handleReset" />
        </template>

        <!-- 右侧结果区 -->
        <template #header>
          <TimeControls v-model:interval="form.time_grain.interval" v-model:static-range="staticTimeRange"
            v-model:comparison-range="comparisonTimeRange" @dynamic-change="handleDynamicTimeChange" />
        </template>

        <template #results>
          <div class="result-area">
            <ResultView :loading="loading" :rows="results" :columns="columns" :viz-mode="vizMode" :sql="sql"
              @update:viz-mode="vizMode = $event" @show-sql="showSqlDrawer = true" />
          </div>
        </template>
      </EventAnalysisLayout>
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
import { FormulaDefaults, TableNames } from '@/constants/analysis'
import { useAppStore } from '@/store/app'
import { MetricFormulaLabels } from '@/types/doris/constants'
import { FilterScope } from '@/types/doris/filter'
import { ElMessage } from 'element-plus'
import { onMounted, provide, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import ActionFooter from './components/ActionFooter.vue'
import EventAnalysisLayout from './components/EventAnalysisLayout.vue'
import GlobalFiltersCard from './components/GlobalFiltersCard.vue'
import GlobalGroupsCard from './components/GlobalGroupsCard.vue'
import MetricsList from './components/metrics/MetricsList.vue'
import ReportDrawer from './components/ReportDrawer.vue'
import ResultView from './components/ResultView.vue'
import SaveReportDialog from './components/SaveReportDialog.vue'
import SqlDrawer from './components/SqlDrawer.vue'
import TimeControls from './components/TimeControls.vue'
import { useAnalysisRunner } from './composables/useAnalysisRunner'
import { useCustomFormula } from './composables/useCustomFormula'
import { useEventAnalysisForm } from './composables/useEventAnalysisForm'
import { useEventMeta } from './composables/useEventMeta'
import { usePayloadBuilders } from './composables/usePayloadBuilders'
import { useReportPersistence } from './composables/useReportPersistence'
import { useTimeRangeSync } from './composables/useTimeRangeSync'
import { ANALYSIS_CONTEXT_KEY, type AnalysisContext } from './context'
import { normalizeMetricField } from './utils/payloadNormalizer'

const appStore = useAppStore()
const route = useRoute()
const {
  form, vizMode, loading, results, columns, sql, total,
  addMetric, addCustomMetric, removeMetric, resetForm, loadForm
} = useEventAnalysisForm()
const { eventOptions, propertyOptions, fetchMeta } = useEventMeta()
const { convertToCustomMetric, convertToNormalMetric, syncWebMetricToApiMetric } = useCustomFormula()

const { buildPayload } = usePayloadBuilders()
const { runAnalysis } = useAnalysisRunner()
const {
  staticTimeRange,
  comparisonTimeRange,
  dynamicDateFilter,
  handleStaticTimeChange,
  handleDynamicTimeChange
} = useTimeRangeSync()

const {
  currentReport,
  reportList,
  loading: reportLoading,
  saveReport,
  fetchReportList,
  removeReport
} = useReportPersistence()

const showSqlDrawer = ref(false)
const showSaveDialog = ref(false)
const showReportDrawer = ref(false)

// --- Analysis Context Setup ---
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
  if (form.event_metrics.length === 0) {
    ElMessage.warning('请至少添加一个分析指标')
    return
  }

  form.project_alias = appStore.activeProjectAlias
  form.filter_groups.query_dates = staticTimeRange.value || []
  form.filter_groups.comparison_query_dates = comparisonTimeRange.value || []

  const payload = buildPayload(
    form.project_alias,
    form.event_metrics,
    form.filter_groups,
    form.time_grain,
    form.groups,
    form.orders,
    dynamicDateFilter.value
  )

  // 执行分析并更新结果
  const res = await runAnalysis(payload)
  if (res) {
    results.value = res.rows
    columns.value = res.columns
    total.value = res.count
    sql.value = res.sql
  }
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

// --- 指标操作增强 ---
const getDisplayPropertyName = (m: any) => {
  if (m.metric?.column?.field === '__TOTAL_TIMES__') return '总次数'
  if (m.metric?.column?.field === '__TOTAL_USERS__') return '总用户数'
  // 这里可以进一步从 propertyOptions 中查找 property_name
  const prop = propertyOptions.value.find(p => p.id === m.metric?.column?.field)
  return prop?.name || m.metric?.column?.field || ''
}

const getFormulaName = (formula: number) => {
  return MetricFormulaLabels[formula] || ''
}

const onApplyFormula = (m: any, payload: any) => {
  if (!m.metric) m.metric = { column: { table: TableNames.EVENT, field: '' }, formula: 2, format: 'raw' }
  if (!m.metric.column) m.metric.column = { table: TableNames.EVENT, field: '' }

  const { field, table, formula, propertyName } = payload

  if (['__TOTAL_TIMES__', '__TOTAL_USERS__'].includes(field)) {
    normalizeMetricField(m, field)
  } else {
    m.metric.column.field = field
    m.metric.column.table = table
    m.metric.formula = formula
  }

  // 自动生成指标名称
  const eventName = eventOptions.value.find(e => e.id === m.e_event_id)?.name || m.e_event_id
  const propName = propertyName || field
  const formulaName = getFormulaName(m.metric.formula)
  m.name = `${eventName}.${propName} · ${formulaName}`
}

const onToggleFormula = (m: any) => {
  if (m.isCustom) {
    convertToNormalMetric(m)
  } else {
    convertToCustomMetric(m)
  }
}

const onCopyMetric = (m: any) => {
  const copy = JSON.parse(JSON.stringify(m))
  copy._uid = Date.now()
  form.event_metrics.push(copy)
}

const onToggleFilter = (payload: any) => {
  const m = payload.m || payload
  m.filterPopoverVisible = true
  if (!m.filter_group) {
    m.filter_group = { scope: FilterScope.And, filters: [] }
  }
  if (!m.filter_group.filters) {
    m.filter_group.filters = []
  }
  m.filter_group.filters.push({
    column: { table: TableNames.EVENT, field: '', alias: '' },
    operator: 1,
    value: ''
  })
}

const onEventChange = (m: any) => {
  // 当事件变化时，规范化为总次数
  if (m.metric && m.metric.column) {
    normalizeMetricField(m, '__TOTAL_TIMES__')
  }
  m.name = eventOptions.value.find(e => e.id === m.e_event_id)?.name || m.e_event_id
}

const onUpdateCustomElements = (m: any, elements: any[]) => {
  if (!m.custom_web_metric) m.custom_web_metric = { eventMetrics: [], format: FormulaDefaults.DEFAULT_FORMAT }
  m.custom_web_metric.eventMetrics = [...elements]
  syncWebMetricToApiMetric(m)
}

const onRequestNewEvent = (m: any) => {
  // 此处原用于自定义公式内插入事件，现已改为由 CustomFormulaInput 组件内部处理。
  // 如果未来有其他入口触发插入事件，可在此处扩展。
}

const onOpenCustomFilter = (metricInFormula: any) => {
  // 处理公式内部事件的过滤器打开（已经在组件内处理）
}

const onReorder = (newMetrics: any[]) => {
  form.event_metrics.splice(0, form.event_metrics.length, ...newMetrics)
}

const onSaveConfirm = async (data: { name: string; description: string; dashboardId?: number }) => {
  const success = await saveReport({
    ...data,
    category: 'events',
    queryParams: {
      form,
      staticTimeRange: staticTimeRange.value,
      comparisonTimeRange: comparisonTimeRange.value,
      dynamicDateFilter: dynamicDateFilter.value
    },
    projectAlias: appStore.activeProjectAlias
  })
  if (success) {
    showSaveDialog.value = false
  }
}

const openReportDrawer = () => {
  fetchReportList(appStore.activeProjectAlias, 'events')
  showReportDrawer.value = true
}

const onLoadReport = (report: any) => {
  try {
    const params = typeof report.query_params === 'string'
      ? JSON.parse(report.query_params)
      : report.query_params
    const formData = params.form || params

    // 使用 loadForm 保持响应性并更新状态
    loadForm(formData)

    // 确保加载后的自定义指标也是同步的
    form.event_metrics.forEach((m: any) => {
      if (m.isCustom) syncWebMetricToApiMetric(m)
    })

    const filterData = formData.filter_groups || formData.global_filter_groups
    staticTimeRange.value = params.staticTimeRange || filterData?.query_dates || null
    comparisonTimeRange.value = params.comparisonTimeRange || filterData?.comparison_query_dates || null
    dynamicDateFilter.value = params.dynamicDateFilter
    currentReport.value = report

    // 加载完成后自动分析
    setTimeout(() => {
      handleAnalyze()
    }, 0)
  } catch (e) {
    console.error('Report parse error:', e)
    ElMessage.error('报表解析失败')
  }
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

const handleReset = () => {
  // 重置表单到初始状态
  resetForm()
  form.project_alias = appStore.activeProjectAlias

  // 重置时间范围
  const end = new Date()
  const start = new Date()
  start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
  const format = (d: Date) => d.toISOString().split('T')[0]
  staticTimeRange.value = [format(start), format(end)]
  comparisonTimeRange.value = null
  dynamicDateFilter.value = null

  currentReport.value = null
}

// 监听项目切换
watch(() => appStore.activeProjectAlias, (newAlias) => {
  if (newAlias) {
    fetchMeta()
    handleReset()
  }
})

watch(() => route.query.reportId, () => {
  loadReportFromRoute()
})

onMounted(async () => {
  // 初始化默认日期范围：过去7天
  const end = new Date()
  const start = new Date()
  start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
  const format = (d: Date) => d.toISOString().split('T')[0]
  staticTimeRange.value = [format(start), format(end)]
  handleStaticTimeChange(staticTimeRange.value)

  // 默认添加一个指标
  if (form.event_metrics.length === 0) {
    addMetric()
  }

  // 获取元数据
  await fetchMeta()
  await loadReportFromRoute()
})
</script>

<style scoped lang="scss">
.event-analysis-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.analysis-content {
  flex: 1;
  min-height: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;

  .page-title {
    font-size: 18px;
    font-weight: bold;
  }
}

.result-area {
  display: flex;
  flex-direction: column;
  gap: 16px;
  height: 100%;
}

.mb-15 {
  margin-bottom: 15px;
}

.ml-10 {
  margin-left: 10px;
}
</style>

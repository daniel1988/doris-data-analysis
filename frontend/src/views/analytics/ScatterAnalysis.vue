<template>
  <div class="scatter-analysis-container">
    <div class="page-header mb-15">
      <div class="header-left">
        <span class="page-title">分布分析</span>
        <el-tag v-if="currentReport" type="success" closable @close="handleReset" class="ml-10">
          {{ currentReport.name }}
        </el-tag>
      </div>
      <div class="header-right">
        <el-button icon="Files" @click="showReportDrawer = true">我的报表</el-button>
      </div>
    </div>

    <div class="analysis-content">
      <ScatterAnalysisLayout>
        <!-- 左侧配置区 -->
        <template #config>
          <ScatterMetricsList 
            :get-display-property-name="getDisplayPropertyName"
            @edit-name="(m: any) => m.isEditingName = true"
            @finish-edit-name="(m: any) => m.isEditingName = false"
            @toggle-formula="onToggleFormula"
            @toggle-filter="onToggleFilter"
            @event-change="onEventChange"
            @apply-formula="onApplyFormula"
          />

          <ScatterSettingsCard />

          <GlobalFiltersCard />

          <GlobalGroupsCard />
        </template>

        <template #footer>
          <ActionFooter 
            :loading="loading" 
            @analyze="handleAnalyze" 
            @save="showSaveDialog = true"
            @reset="handleReset" 
          />
        </template>

        <!-- 右侧结果区 -->
        <template #header>
          <TimeControls 
            v-model:interval="form.time_grain.interval" 
            :static-range="staticTimeRange"
            :hide-comparison="true"
            @update:static-range="handleStaticTimeChange"
            @dynamic-change="handleDynamicTimeChange" 
          />
        </template>

        <template #results>
          <div class="result-area">
            <ScatterResultView 
              :loading="loading" 
              :rows="results" 
              :columns="columns" 
              :viz-mode="vizMode" 
              :sql="sql"
              @update:viz-mode="vizMode = $event" 
              @show-sql="showSqlDrawer = true" 
            />
          </div>
        </template>
      </ScatterAnalysisLayout>
    </div>

    <!-- SQL 预览抽屉 -->
    <SqlDrawer v-model="showSqlDrawer" :sql="sql" />

    <!-- 保存报表弹窗 -->
    <SaveReportDialog 
      v-model="showSaveDialog" 
      :initial-data="currentReport || { name: '', description: '' }"
      @confirm="onSaveConfirm" 
    />

    <!-- 报表列表抽屉 -->
    <ReportDrawer 
      v-model="showReportDrawer" 
      :reports="reportList" 
      :loading="reportLoading" 
      @load="onLoadReport"
      @delete="removeReport" 
    />
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'ScatterAnalysis'
})
import { getReportById } from '@/api/report'
import { useAppStore } from '@/store/app'
import { ElMessage } from 'element-plus'
import { onMounted, provide, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

import ActionFooter from './components/ActionFooter.vue'
import GlobalFiltersCard from './components/GlobalFiltersCard.vue'
import GlobalGroupsCard from './components/GlobalGroupsCard.vue'
import ScatterMetricsList from './components/metrics/ScatterMetricsList.vue'
import ReportDrawer from './components/ReportDrawer.vue'
import SaveReportDialog from './components/SaveReportDialog.vue'
import ScatterAnalysisLayout from './components/ScatterAnalysisLayout.vue'
import ScatterResultView from './components/ScatterResultView.vue'
import ScatterSettingsCard from './components/ScatterSettingsCard.vue'
import SqlDrawer from './components/SqlDrawer.vue'
import TimeControls from './components/TimeControls.vue'

import { useEventMeta } from './composables/useEventMeta'
import { useReportPersistence } from './composables/useReportPersistence'
import { useScatterForm } from './composables/useScatterForm'
import { useTimeRangeSync } from './composables/useTimeRangeSync'
import { ANALYSIS_CONTEXT_KEY, type AnalysisContext } from './context'
import { normalizeMetricField } from './utils/payloadNormalizer'

const appStore = useAppStore()
const route = useRoute()

const {
  form, vizMode, loading, results, columns, sql, total, duration,
  runAnalysis, resetForm, loadForm
} = useScatterForm()

const { eventOptions, propertyOptions, fetchMeta } = useEventMeta()

const {
  staticTimeRange,
  handleStaticTimeChange,
  handleDynamicTimeChange,
  dynamicDateFilter,
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
  if (!form.metric?.e_event_id) {
    ElMessage.warning('请选择分析事件')
    return
  }

  form.project_alias = appStore.activeProjectAlias
  form.global_filter_groups.query_dates = staticTimeRange.value || []

  // 处理动态日期：如果有动态日期，则添加到全局过滤器中
  // 首先移除已有的动态日期过滤器 (operator 19)
  form.global_filter_groups.global_filters.filters = 
    form.global_filter_groups.global_filters.filters.filter(f => f.operator !== 19)
  
  if (dynamicDateFilter.value) {
    form.global_filter_groups.global_filters.filters.push({
      column: { table: 'event_data', field: 'e_event_time', alias: '日期' },
      operator: 19,
      value: dynamicDateFilter.value.value
    })
  }

  // 执行分析
  await runAnalysis(form)
}

const handleReset = () => {
  resetForm()
  currentReport.value = null
}

const contextActions = {
  refreshMetadata: fetchMeta,
  triggerAnalyze: handleAnalyze,
  syncTimeRange: (staticRange: string[] | null, _comparisonRange: string[] | null) => {
    handleStaticTimeChange(staticRange as [string, string] | null)
  }
}

provide(ANALYSIS_CONTEXT_KEY, {
  state: contextState,
  actions: contextActions
} as AnalysisContext)

// --- Event Handlers ---
const getDisplayPropertyName = (m: any) => {
  if (m.metric?.column?.field === '__TOTAL_TIMES__') return '总次数'
  if (m.metric?.column?.field === '__TOTAL_USERS__') return '总用户数'
  const prop = propertyOptions.value.find(p => p.id === m.metric?.column?.field)
  return prop?.name || m.metric?.column?.field || ''
}

const onToggleFormula = (m: any) => {
  // Logic to toggle between count/sum/avg/etc. 
  // For scatter, we usually let MetricPropertyPicker handle this via onApplyFormula
}

const onToggleFilter = (payload: { m: any; index: number }) => {
  if (payload.m) {
    payload.m.filterPopoverVisible = !payload.m.filterPopoverVisible
  }
}

const onEventChange = (m: any) => {
  if (m && m.metric) {
    normalizeMetricField(m, '__TOTAL_TIMES__')
  }
}

const onApplyFormula = (payload: any) => {
  const { m, field, formula, table } = payload
  if (m && m.metric) {
    if (['__TOTAL_TIMES__', '__TOTAL_USERS__'].includes(field)) {
      normalizeMetricField(m, field)
    } else {
      m.metric.column.field = field
      m.metric.column.table = table
      m.metric.formula = formula
    }
  }
}

const onSaveConfirm = async (data: { name: string; description: string; dashboardId?: number }) => {
  const success = await saveReport({
    ...data,
    category: 'scatter',
    queryParams: form,
    projectAlias: appStore.activeProjectAlias
  })
  if (success) {
    showSaveDialog.value = false
    fetchReportList(appStore.activeProjectAlias, 'scatter')
  }
}

const onLoadReport = (report: any) => {
  currentReport.value = report
  const params = typeof report.query_params === 'string' 
    ? JSON.parse(report.query_params) 
    : report.query_params
  loadForm(params)
  showReportDrawer.value = false

  // 同步时间状态到 useTimeRangeSync
  syncFromForm(params.global_filter_groups || params.filter_groups)
  
  // 自动触发分析
  handleAnalyze()
}

const loadReportFromRoute = async () => {
  const reportId = Number(route.query.reportId)
  if (!reportId) return
  try {
    const report = await getReportById(reportId)
    onLoadReport(report)
  } catch (e) {
    console.error('Failed to load report:', e)
  }
}

// --- Lifecycle ---
onMounted(async () => {
  await fetchMeta()
  await fetchReportList(appStore.activeProjectAlias, 'scatter')
  await loadReportFromRoute()
})

watch(() => route.query.reportId, () => {
  loadReportFromRoute()
})

watch(() => appStore.activeProjectAlias, (newAlias) => {
  if (newAlias) {
    contextState.projectAlias = newAlias
    fetchMeta()
    fetchReportList(newAlias, 'scatter')
  }
})
</script>

<style scoped lang="scss">
.scatter-analysis-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 20px;
  background-color: var(--el-bg-color-page);

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-shrink: 0;

    .header-left {
      display: flex;
      align-items: center;
      .page-title {
        font-size: 18px;
        font-weight: 600;
      }
    }
  }

  .analysis-content {
    flex: 1;
    min-height: 0;
  }
}

.mb-15 {
  margin-bottom: 15px;
}

.ml-10 {
  margin-left: 10px;
}
</style>

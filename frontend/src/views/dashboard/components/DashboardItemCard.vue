<template>
  <el-card shadow="hover" class="dashboard-item-card"
    :body-style="{ padding: '10px', height: '100%', boxSizing: 'border-box' }">
    <template #header>
      <div class="card-header">
        <span class="title clickable-title" @click="handleOpenReport">
          {{ item.title || report?.name || '未命名报表' }}
        </span>
        <div class="actions">
          <el-radio-group v-if="canToggleView" v-model="showTable" size="small" class="view-toggle">
            <el-radio-button :value="true">
              <el-icon><List /></el-icon>
            </el-radio-button>
            <el-radio-button :value="false">
              <el-icon><DataLine /></el-icon>
            </el-radio-button>
          </el-radio-group>
          <el-icon class="action-icon" @click="refreshData" title="刷新数据">
            <Refresh />
          </el-icon>
          <el-dropdown trigger="click" @command="handleCommand">
            <el-icon class="action-icon">
              <More />
            </el-icon>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="edit">编辑报表</el-dropdown-item>
                <el-dropdown-item command="size-medium" v-if="item.width !== 12">设为中图</el-dropdown-item>
                <el-dropdown-item command="size-large" v-if="item.width !== 24">设为大图</el-dropdown-item>
                <el-dropdown-item command="remove" divided style="color: #f56c6c">从看板移除</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </template>

    <div v-loading="loading" class="card-content">
      <template v-if="error">
        <el-result icon="error" title="加载失败" :sub-title="error"></el-result>
      </template>
      <template v-else-if="chartData">
        <StatCardView v-if="isStatMode" :rows="chartData.rows" :columns="chartData.columns" />
        <template v-else>
          <TableView v-if="showTable" :rows="chartData.rows" :columns="chartData.columns" hide-header />
          <ChartView v-else :rows="chartData.rows" :columns="chartData.columns" :viz-mode="vizMode" />
        </template>
      </template>
      <template v-else-if="!loading">
        <el-empty description="暂无数据" :image-size="60" />
      </template>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import {
  eventAnalysis,
  funnelAnalysis,
  retentionAnalysis,
  scatterAnalysis,
  userPropertyAnalysis
} from '@/api/analytics'
import type { DashboardItem } from '@/api/dashboard'
import type { Report } from '@/api/report'
import { getReportById } from '@/api/report'
import ChartView from '@/views/analytics/components/ChartView.vue'
import TableView from '@/views/analytics/components/TableView.vue'
import { usePayloadBuilders } from '@/views/analytics/composables/usePayloadBuilders'
import { DataLine, List, More, Refresh } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import StatCardView from './StatCardView.vue'

const props = defineProps<{
  item: DashboardItem
  globalFilters: {
    timeRange: [string, string] | null
    filters: any[]
  }
}>()

const emit = defineEmits(['remove', 'update-size'])

const router = useRouter()
const { buildPayload, buildFunnelPayload, buildRetentionPayload, buildPropertyPayload } = usePayloadBuilders()
const report = ref<Report | null>(null)
const loading = ref(false)
const error = ref('')
const chartData = ref<{ rows: any[], columns: string[] } | null>(null)
const showTable = ref(true)

const canToggleView = computed(() => {
  if (isStatMode.value) return false
  if (report.value?.category === 'retention') return false
  return true
})

const normalizeAnalysisResponse = (res: any) => {
  const payload = res?.rows ? res : (res?.data || {})
  const rawColumns = Array.isArray(payload.columns) ? payload.columns : []
  let columns = rawColumns.map((c: any) => {
    if (typeof c === 'string') return c
    return c?.name || c?.label || String(c?.field || '')
  }).filter(Boolean)

  const rawRows = Array.isArray(payload.rows) ? payload.rows : []
  const rows = rawRows.map((row: any) => {
    if (Array.isArray(row) && columns.length > 0) {
      return columns.reduce((acc: Record<string, any>, col, index) => {
        acc[col] = row[index]
        return acc
      }, {})
    }
    return row
  })

  if (rows.length > 0 && !Array.isArray(rows[0]) && typeof rows[0] === 'object') {
    const rowKeys = Object.keys(rows[0])
    if (columns.length === 0) {
      columns = rowKeys
    } else {
      // 补全缺失的列名（如果 rows 中的 key 在 columns 中没有）
      const missingKeys = rowKeys.filter(key => !columns.includes(key))
      if (missingKeys.length > 0) {
        columns = [...columns, ...missingKeys]
      }
    }
  }

  // 特殊处理：如果列名中包含 e_openid/u_openid，尝试映射为更有意义的名称
  columns = columns.map(c => {
    if (c === 'e_openid' || c === 'u_openid') return '用户数'
    if (c === 'e_count') return '总次数'
    return c
  })

  // 同时更新 rows 中的 key，保持一致
  const finalRows = rows.map(row => {
    const newRow: any = {}
    Object.keys(row).forEach(key => {
      let newKey = key
      if (key === 'e_openid' || key === 'u_openid') newKey = '用户数'
      if (key === 'e_count') newKey = '总次数'
      newRow[newKey] = row[key]
    })
    return newRow
  })

  return { rows: finalRows, columns }
}

const vizMode = computed(() => {
  // 1. 优先使用看板项覆盖的配置
  if (props.item.config_override) {
    try {
      const config = typeof props.item.config_override === 'string'
        ? JSON.parse(props.item.config_override)
        : props.item.config_override
      if (config.vizMode) return config.vizMode
    } catch (e) { }
  }

  // 2. 其次使用报表保存时的配置
  if (report.value?.query_params) {
    try {
      const params = typeof report.value.query_params === 'string'
        ? JSON.parse(report.value.query_params)
        : report.value.query_params

      // 检查外层或 form 内层
      return params.vizMode || params.form?.vizMode || 'line'
    } catch (e) { }
  }
  return 'line'
})

const isStatMode = computed(() => {
  return props.item.type === 'stat' || vizMode.value === 'stat' || vizMode.value === 'metric'
})

const fetchReportAndData = async () => {
  loading.value = true
  error.value = ''
  try {
    const res = await getReportById(props.item.report_id)
    report.value = res as unknown as Report
    if (import.meta.env.DEV) {
      console.log('[DashboardItemCard] fetchReportAndData', {
        reportId: props.item.report_id,
        category: report.value?.category,
        queryParams: report.value?.query_params
      })
    }
    await loadData()
  } catch (err: any) {
    console.error('[DashboardItemCard] fetchReportAndData error:', err)
    error.value = err.message || '获取报表详情失败'
  } finally {
    loading.value = false
  }
}

const loadData = async () => {
  if (!report.value) {
    if (import.meta.env.DEV) console.warn('[DashboardItemCard] No report available for loadData')
    return
  }

  try {
    const params = typeof report.value.query_params === 'string'
      ? JSON.parse(report.value.query_params)
      : report.value.query_params
    const form = params.form || params
    const projectAlias = report.value.project_alias || appStore.activeProjectAlias

    if (!projectAlias) {
      if (import.meta.env.DEV) console.warn('[DashboardItemCard] No project_alias for loadData')
      return
    }

    // 合并全局过滤器逻辑
    let queryDates = props.globalFilters.timeRange || params.staticTimeRange || form.filter_groups?.query_dates || []
    let comparisonDates = params.comparisonTimeRange || form.filter_groups?.comparison_query_dates || []

    // 如果是指标卡，需要特殊处理时间以支持环比
    if (isStatMode.value && queryDates.length === 2) {
      const start = dayjs(queryDates[0])
      const end = dayjs(queryDates[1])
      const diff = end.diff(start, 'day') + 1
      const prevStart = start.subtract(diff, 'day').format('YYYY-MM-DD')
      const prevEnd = start.subtract(1, 'day').format('YYYY-MM-DD')
      comparisonDates = [prevStart, prevEnd]
    }

    const mergedGlobalFilters = [
      ...(form?.filter_groups?.global_filters?.filters || []),
      ...props.globalFilters.filters
    ]
    const filterGroups = {
      ...form?.filter_groups,
      query_dates: queryDates,
      comparison_query_dates: comparisonDates,
      global_filters: {
        relation: form?.filter_groups?.global_filters?.relation || 'and',
        filters: mergedGlobalFilters
      }
    }

    let analysisFn: any
    let payload: any

    const category = report.value.category || 'events'

    if (import.meta.env.DEV) {
      console.log('[DashboardItemCard] Preparing payload', {
        reportId: props.item.report_id,
        category,
        form
      })
    }

    switch (category) {
      case 'events':
        payload = buildPayload(
          projectAlias,
          form.event_metrics || [],
          filterGroups,
          form.time_grain || 'day',
          form.groups || [],
          form.orders || [],
          params.dynamicDateFilter
        )
        analysisFn = eventAnalysis
        break
      case 'funnel':
        payload = buildFunnelPayload(
          projectAlias,
          form.event_metrics || [],
          filterGroups,
          form.time_grain || 'day',
          form.groups || [],
          params.dynamicDateFilter
        )
        analysisFn = funnelAnalysis
        break
      case 'retention':
        payload = buildRetentionPayload(
          projectAlias,
          form.init_event_metric,
          form.end_event_metric,
          filterGroups,
          form.time_grain || 'day',
          form.groups || [],
          form.day_n_array || [1, 3, 7, 30],
          params.dynamicDateFilter
        )
        analysisFn = retentionAnalysis
        break
      case 'property':
        payload = buildPropertyPayload(
          projectAlias,
          form.metric,
          filterGroups,
          form.groups || [],
          form.user_groups || [],
          form.group_type,
          params.dynamicDateFilter
        )
        analysisFn = userPropertyAnalysis
        break
      case 'scatter':
        payload = { ...form, project_alias: projectAlias, filter_groups: filterGroups }
        analysisFn = scatterAnalysis
        break
      default:
        console.warn(`[DashboardItemCard] Unknown category: ${category}`)
    }

    if (analysisFn && payload) {
      if (import.meta.env.DEV) {
        console.log('[DashboardItemCard] Executing analysis', {
          reportId: props.item.report_id,
          payload
        })
      }
      const res = await analysisFn(payload)
      const normalized = normalizeAnalysisResponse(res)
      if (import.meta.env.DEV) {
        console.log('[DashboardItemCard] normalized data', {
          reportId: props.item.report_id,
          category: report.value?.category,
          rows: normalized.rows.length,
          columns: normalized.columns
        })
      }
      chartData.value = {
        rows: normalized.rows,
        columns: normalized.columns
      }
    } else {
      console.warn('[DashboardItemCard] analysisFn or payload missing', { analysisFn: !!analysisFn, payload: !!payload })
    }
  } catch (err: any) {
    console.error('[DashboardItemCard] loadData error:', err)
    error.value = '数据加载解析失败: ' + (err.message || '')
  }
}

const refreshData = () => {
  loadData()
}

const handleOpenReport = () => {
  const pathMap: Record<string, string> = {
    'events': '/analytics/event',
    'funnel': '/analytics/funnel',
    'retention': '/analytics/retention',
    'property': '/analytics/property',
    'scatter': '/analytics/scatter'
  }
  const path = pathMap[report.value?.category || '']
  if (path) {
    router.push({ path, query: { reportId: props.item.report_id } })
  }
}

const handleCommand = (command: string) => {
  if (command === 'edit') {
    handleOpenReport()
  } else if (command === 'size-medium') {
    emit('update-size', props.item.id, 12)
  } else if (command === 'size-large') {
    emit('update-size', props.item.id, 24)
  } else if (command === 'remove') {
    emit('remove', props.item.id)
  }
}

onMounted(() => {
  fetchReportAndData()
})

watch(() => props.item.report_id, () => {
  fetchReportAndData()
})

watch(() => props.globalFilters, () => {
  loadData()
}, { deep: true })
</script>

<style scoped>
.dashboard-item-card {
  height: 100%;
}

:deep(.el-card) {
  display: flex;
  flex-direction: column;
  height: 100%;
  border: none;
}

:deep(.el-card__header) {
  padding: 8px 12px;
  border-bottom: 1px solid #f0f2f5;
  flex-shrink: 0;
}

:deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  padding: 10px !important;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-weight: bold;
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 80%;
}

.clickable-title {
  cursor: pointer;
}

.clickable-title:hover {
  color: var(--el-color-primary);
}

.actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.view-toggle :deep(.el-radio-button__inner) {
  padding: 5px 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-icon {
  cursor: pointer;
  color: #909399;
  transition: color 0.2s;
  font-size: 16px;
}

.action-icon:hover {
  color: #409eff;
}

.card-content {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}
</style>

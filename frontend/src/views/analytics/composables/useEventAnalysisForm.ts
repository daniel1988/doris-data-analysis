import { VizModes } from '@/constants/analysis'
import type { AnalysisPayload } from '@/types/doris/analysis'
import { reactive, ref } from 'vue'
import { useAnalysisRunner } from './useAnalysisRunner'
import { useFilterState } from './useFilterState'
import { useGroupState } from './useGroupState'
import { useMetricState } from './useMetricState'

export function useEventAnalysisForm() {
  const { metrics, addMetric, addCustomMetric, removeMetric, clearMetrics } = useMetricState()
  const { filterGroups, addGlobalFilter, removeGlobalFilter, clearFilters } = useFilterState()
  const { groups, addGroup, removeGroup, clearGroups } = useGroupState()
  const { loading, results, columns, total, sql, runAnalysis, clearResults } = useAnalysisRunner()

  const form = reactive<AnalysisPayload>({
    project_alias: '',
    event_metrics: metrics,
    filter_groups: filterGroups,
    time_grain: {
      column: { table: 'event_data', field: 'e_event_time', alias: '日期' },
      interval: 2, // Tg_Interval_Day
      window_num: 0
    },
    groups: groups,
    orders: []
  } as any)

  const vizMode = ref<VizModes>(VizModes.Line)

  const resetForm = () => {
    form.project_alias = ''
    clearMetrics()
    clearFilters()
    clearGroups()
    clearResults()
    addMetric() // Add one default metric
  }

  const loadForm = (data: any) => {
    if (!data) return

    // Update project alias
    if (data.project_alias) form.project_alias = data.project_alias

    // Update metrics without breaking reference if possible, 
    // but here we just ensure the form property is updated.
    // If the caller uses the functions from this composable, 
    // we need to make sure those functions point to the current form.event_metrics.

    if (data.event_metrics) {
      metrics.splice(0, metrics.length, ...data.event_metrics)
    }

    const filterData = data.filter_groups || data.global_filter_groups
    if (filterData) {
      Object.assign(filterGroups, filterData)
    }

    if (data.groups) {
      groups.splice(0, groups.length, ...data.groups)
    }

    if (data.time_grain) {
      Object.assign(form.time_grain, data.time_grain)
    }
  }

  return {
    form,
    loading,
    results,
    columns,
    total,
    sql,
    vizMode,
    addMetric,
    addCustomMetric,
    removeMetric,
    resetForm,
    loadForm,
    runAnalysis
  }
}

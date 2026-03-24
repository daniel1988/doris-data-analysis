import { FormulaDefaults, TableNames } from '@/constants/analysis'
import { TimeGrainInterval } from '@/types/doris/analysis'
import { FormatDefault, Formula } from '@/types/doris/common'
import { FilterScope } from '@/types/doris/filter'
import { EventType, type EventMetric } from '@/types/doris/metric'
import { reactive, ref } from 'vue'
import { useFilterState } from './useFilterState'
import { useGroupState } from './useGroupState'
import { useFunnelAnalysis } from './useFunnelAnalysis'

export function useFunnelForm() {
  const { filterGroups, clearFilters } = useFilterState()
  const { groups, clearGroups } = useGroupState()
  const { loading, results, columns, total, sql, runFunnelAnalysis, clearResults } = useFunnelAnalysis()

  const createDefaultStep = (index: number): EventMetric => ({
    _uid: Date.now() + Math.random(),
    e_event_id: '',
    name: `步骤 ${index + 1}`,
    type: EventType.Normal,
    metric: {
      column: { table: TableNames.EVENT, field: FormulaDefaults.TOTAL_USERS_FIELD, alias: '' },
      formula: Formula.BitmapUnion,
      format: FormatDefault
    },
    scope: FilterScope.And,
    filters: [],
    tag_filters: [],
    user_group_filters: []
  })

  const form = reactive({
    project_alias: '',
    event_metrics: [createDefaultStep(0), createDefaultStep(1)], // 默认两步
    filter_groups: filterGroups,
    groups: groups,
    time_grain: {
      column: { table: TableNames.EVENT, field: 'e_event_time', alias: '日期' },
      interval: TimeGrainInterval.Day,
      window_num: 1 // 默认 1 天转化窗口
    },
    page_size: 1000,
    page_num: 1
  })

  const addStep = () => {
    if (form.event_metrics.length >= 8) return
    form.event_metrics.push(createDefaultStep(form.event_metrics.length))
  }

  const removeStep = (index: number) => {
    if (form.event_metrics.length <= 2) return
    form.event_metrics.splice(index, 1)
    // 更新剩余步骤名称
    form.event_metrics.forEach((m, idx) => {
      if (m.name.startsWith('步骤 ')) {
        m.name = `步骤 ${idx + 1}`
      }
    })
  }

  const resetForm = () => {
    form.project_alias = ''
    form.event_metrics = [createDefaultStep(0), createDefaultStep(1)]
    form.time_grain.interval = TimeGrainInterval.Day
    form.time_grain.window_num = 1
    clearFilters()
    clearGroups()
    clearResults()
  }

  const loadForm = (data: any) => {
    if (!data) return
    if (data.project_alias) form.project_alias = data.project_alias
    if (data.event_metrics) {
      form.event_metrics = data.event_metrics.map((m: any, idx: number) => ({
        ...m,
        _uid: Date.now() + Math.random() + idx
      }))
    }
    const filterData = data.filter_groups || data.global_filter_groups
    if (filterData) Object.assign(filterGroups, filterData)
    if (data.groups) {
      groups.splice(0, groups.length, ...data.groups)
    }
    if (data.time_grain) Object.assign(form.time_grain, data.time_grain)
  }

  return {
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
  }
}

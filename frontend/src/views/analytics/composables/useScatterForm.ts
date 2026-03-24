import { VizModes } from '@/constants/analysis'
import type { ScatterAnalysisReq } from '@/types/doris/analysis'
import { FormatDefault, Formula } from '@/types/doris/common'
import { FilterScope } from '@/types/doris/filter'
import { EventType } from '@/types/doris/metric'
import { reactive, ref } from 'vue'
import { useFilterState } from './useFilterState'
import { useGroupState } from './useGroupState'
import { useScatterRunner } from './useScatterRunner'

export function useScatterForm() {
  const { filterGroups, addGlobalFilter, removeGlobalFilter, clearFilters } = useFilterState()
  const { groups, addGroup, removeGroup, clearGroups } = useGroupState()
  const { loading, results, columns, sql, total, duration, runAnalysis, clearResults } = useScatterRunner()

  const defaultMetric = () => ({
    _uid: Date.now(),
    e_event_id: '',
    name: '分析指标',
    type: EventType.Normal,
    metric: {
      column: { table: 'event_data', field: 'e_openid', alias: '' },
      formula: Formula.Count,
      format: FormatDefault
    },
    scope: FilterScope.And,
    filters: [],
    tag_filters: [],
    user_group_filters: [],
    filterPopoverVisible: false
  })

  const form = reactive<ScatterAnalysisReq>({
    project_alias: '',
    metric: defaultMetric() as any,
    global_filter_groups: filterGroups,
    groups: groups,
    time_grain: {
      column: { table: 'event_data', field: 'e_event_time', alias: '日期' },
      interval: 2, // Tg_Interval_Day
      window_num: 0
    },
    scatter_type: 1, // 1: 离散值, 2: 自动分桶, 3: 自定义区间
    scatter_ranges: [],
    bin_count: 10
  } as any)

  const vizMode = ref<VizModes>(VizModes.Bar)

  const resetForm = () => {
    form.project_alias = ''
    form.metric = defaultMetric() as any
    form.scatter_type = 1
    form.scatter_ranges = []
    clearFilters()
    clearGroups()
    clearResults()
  }

  const loadForm = (data: any) => {
    if (!data) return
    if (data.project_alias) form.project_alias = data.project_alias
    if (data.metric) Object.assign(form.metric, data.metric)
    const filterData = data.filter_groups || data.global_filter_groups
    if (filterData) Object.assign(filterGroups, filterData)
    if (data.groups) groups.splice(0, groups.length, ...data.groups)
    if (data.time_grain) Object.assign(form.time_grain, data.time_grain)
    if (data.scatter_type) form.scatter_type = data.scatter_type
    if (data.scatter_ranges) form.scatter_ranges = [...data.scatter_ranges]
    if (data.bin_count) form.bin_count = data.bin_count
  }

  const addRange = () => {
    if (!form.scatter_ranges) form.scatter_ranges = []
    form.scatter_ranges.push({ min: 0, max: 0 })
  }

  const removeRange = (index: number) => {
    if (form.scatter_ranges) {
      form.scatter_ranges.splice(index, 1)
    }
  }

  return {
    form,
    loading,
    results,
    columns,
    sql,
    total,
    duration,
    vizMode,
    resetForm,
    loadForm,
    runAnalysis,
    addRange,
    removeRange
  }
}

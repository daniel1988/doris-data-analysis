import { FormulaDefaults, TableNames } from '@/constants/analysis'
import { TimeGrainInterval } from '@/types/doris/analysis'
import { FormatDefault, Formula } from '@/types/doris/common'
import { FilterScope } from '@/types/doris/filter'
import { EventType, type EventMetric } from '@/types/doris/metric'
import { reactive } from 'vue'
import { useFilterState } from './useFilterState'
import { useGroupState } from './useGroupState'
import useRetentionAnalysis from './useRetentionAnalysis'

export function useRetentionForm() {
  const { filterGroups, clearFilters } = useFilterState()
  const { groups, clearGroups } = useGroupState()
  const { loading, results, columns, total, sql, runRetentionAnalysis, clearResults } = useRetentionAnalysis()

  const createDefaultMetric = (name: string): EventMetric => ({
    _uid: Date.now() + Math.random(),
    e_event_id: '',
    name,
    type: EventType.Normal,
    metric: {
      column: { table: TableNames.EVENT, field: FormulaDefaults.TOTAL_USERS_FIELD, alias: '' },
      formula: Formula.Count,
      format: FormatDefault
    },
    scope: FilterScope.And,
    filters: [],
    tag_filters: [],
    user_group_filters: []
  })

  const form = reactive({
    project_alias: '',
    init_event_metric: createDefaultMetric('初始事件'),
    end_event_metric: createDefaultMetric('结束事件'),
    filter_groups: filterGroups,
    groups: groups,
    time_grain: {
      column: { table: 'event_data', field: 'e_event_time', alias: '日期' },
      interval: TimeGrainInterval.Day,
      window_num: 0
    },
    day_n_array: [1, 3, 7, 14, 30]
  })

  const resetForm = () => {
    form.project_alias = ''
    form.init_event_metric = createDefaultMetric('初始事件')
    form.end_event_metric = createDefaultMetric('结束事件')
    form.day_n_array = [1, 3, 7, 14, 30]
    form.time_grain.interval = TimeGrainInterval.Day
    clearFilters()
    clearGroups()
    clearResults()
  }

  const loadForm = (data: any) => {
    if (!data) return
    if (data.project_alias) form.project_alias = data.project_alias

    if (data.init_event_metric) {
      Object.assign(form.init_event_metric, data.init_event_metric)
      form.init_event_metric._uid = Date.now() + Math.random()
    }

    if (data.end_event_metric) {
      Object.assign(form.end_event_metric, data.end_event_metric)
      form.end_event_metric._uid = Date.now() + Math.random()
    }

    if (data.day_n_array) form.day_n_array = [...data.day_n_array]
    if (data.filter_groups) {
      Object.assign(filterGroups, data.filter_groups)
    } else if (data.global_filter_groups) {
      // 兼容可能存在的旧版数据字段名
      Object.assign(filterGroups, data.global_filter_groups)
    }

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
    runRetentionAnalysis,
    resetForm,
    loadForm
  }
}

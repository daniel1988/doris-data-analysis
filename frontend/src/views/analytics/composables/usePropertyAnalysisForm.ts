import { FormulaDefaults, TableNames } from '@/constants/analysis'
import { FormatDefault, Formula } from '@/types/doris/common'
import { FilterScope } from '@/types/doris/filter'
import { reactive, ref } from 'vue'
import { useFilterState } from './useFilterState'
import { useGroupState } from './useGroupState'
import { usePropertyAnalysisRunner } from './usePropertyAnalysisRunner'

export function usePropertyAnalysisForm() {
  const { filterGroups, clearFilters } = useFilterState()
  const { groups, clearGroups } = useGroupState()
  const { loading, results, columns, total, sql, runPropertyAnalysis, clearResults } = usePropertyAnalysisRunner()

  const userGroups = ref<any[]>([])

  const createDefaultMetric = () => ({
    name: '用户数',
    metric: {
      column: { table: TableNames.USER, field: FormulaDefaults.TOTAL_USERS_FIELD },
      formula: Formula.CountDistinct,
      format: FormatDefault
    }
  })

  const form = reactive({
    project_alias: '',
    metric: createDefaultMetric(),
    filter_groups: filterGroups,
    groups: groups,
    user_groups: userGroups,
    group_type: 1, // 1: 维度分组, 2: 人群分群
    page_size: 1000,
    page_num: 1
  })

  const resetForm = () => {
    form.project_alias = ''
    form.metric = createDefaultMetric()
    form.group_type = 1
    userGroups.value = []
    clearFilters()
    clearGroups()
    clearResults()
  }

  const loadForm = (data: any) => {
    if (!data) return
    if (data.project_alias) form.project_alias = data.project_alias
    if (data.metric) Object.assign(form.metric, data.metric)
    if (data.group_type) form.group_type = data.group_type
    const filterData = data.filter_groups || data.global_filter_groups
    if (filterData) Object.assign(filterGroups, filterData)
    if (data.groups) {
      groups.splice(0, groups.length, ...data.groups)
    }
    if (data.user_groups) {
      userGroups.value = [...data.user_groups]
    }
  }

  return {
    form,
    loading,
    results,
    columns,
    total,
    sql,
    runPropertyAnalysis,
    resetForm,
    loadForm
  }
}

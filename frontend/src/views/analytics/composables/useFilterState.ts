import { reactive } from 'vue'
import { FilterScope, type GlobalFilterGroups } from '@/types/doris/filter'

export function useFilterState(initialFilterGroups?: GlobalFilterGroups) {
  const filterGroups = reactive<GlobalFilterGroups>(initialFilterGroups || {
    global_filters: {
      scope: FilterScope.And,
      filters: [],
      tag_filters: [],
      user_group_filters: []
    },
    dashboard_form_filters: {
      scope: FilterScope.And,
      filters: [],
      tag_filters: [],
      user_group_filters: []
    },
    query_dates: []
  })

  const addGlobalFilter = () => {
    filterGroups.global_filters.filters.push({
      column: { table: 'event_data', field: '', alias: '' },
      operator: 1,
      value: ''
    })
  }

  const removeGlobalFilter = (index: number) => {
    filterGroups.global_filters.filters.splice(index, 1)
  }

  const clearFilters = () => {
    filterGroups.global_filters.filters = []
    filterGroups.global_filters.tag_filters = []
    filterGroups.global_filters.user_group_filters = []
  }

  return {
    filterGroups,
    addGlobalFilter,
    removeGlobalFilter,
    clearFilters
  }
}

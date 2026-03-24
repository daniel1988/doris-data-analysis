import { ref, reactive } from 'vue'
import { FormatDefault, Formula } from '@/types/doris/common'
import { FilterScope } from '@/types/doris/filter'
import type { EventMetric } from '@/types/doris/metric'
import { EventType } from '@/types/doris/metric'

export function useMetricState(initialMetrics: EventMetric[] = []) {
  const metrics = reactive<EventMetric[]>(initialMetrics)

  const addMetric = () => {
    const newMetric: EventMetric = {
      _uid: Date.now(),
      e_event_id: '',
      name: '新指标',
      type: EventType.Normal,
      metric: {
        column: { table: 'event_data', field: 'e_openid', alias: '' },
        formula: Formula.Count,
        format: FormatDefault
      },
      scope: FilterScope.And,
      filters: [],
      tag_filters: [],
      user_group_filters: []
    }
    metrics.push(newMetric)
  }

  const removeMetric = (index: number) => {
    metrics.splice(index, 1)
  }

  const updateMetric = (index: number, updated: Partial<EventMetric>) => {
    if (metrics[index]) {
      Object.assign(metrics[index], updated)
    }
  }

  const addCustomMetric = () => {
    const newMetric: EventMetric = {
      _uid: Date.now(),
      e_event_id: '',
      name: '自定义公式',
      type: EventType.Custom,
      metric: {
        column: { table: 'event_data', field: 'e_openid', alias: '' },
        formula: Formula.Count,
        format: FormatDefault
      },
      is_custom: true,
      isCustom: true,
      custom_web_metric: {
        eventMetrics: []
      },
      web_metric: {
        tokens: [],
        format: FormatDefault
      },
      custom_metric: {
        event_metrics: [],
        custom_formula: '',
        format: FormatDefault
      },
      scope: FilterScope.And,
      filters: [],
      tag_filters: [],
      user_group_filters: []
    }
    metrics.push(newMetric)
  }

  const clearMetrics = () => {
    metrics.splice(0, metrics.length)
  }

  return {
    metrics,
    addMetric,
    addCustomMetric,
    removeMetric,
    updateMetric,
    clearMetrics
  }
}

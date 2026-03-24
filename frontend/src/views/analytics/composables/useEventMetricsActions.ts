import { TableNames } from '@/constants/analysis'
import { FormatDefault, Formula } from '@/types/doris/common'
import { FilterScope } from '@/types/doris/filter'
import type { EventMetric } from '@/types/doris/metric'
import { EventType } from '@/types/doris/metric'

export function useEventMetricsActions(eventMetrics: EventMetric[]) {
  const addMetric = () => {
    const newMetric: EventMetric = {
      _uid: Date.now(),
      e_event_id: '',
      name: '新指标',
      type: EventType.Normal,
      metric: {
        column: { table: TableNames.EVENT, field: 'e_openid', alias: '' },
        formula: Formula.Count,
        format: FormatDefault
      },
      scope: FilterScope.And,
      filters: [],
      tag_filters: [],
      user_group_filters: []
    }
    eventMetrics.push(newMetric)
  }

  const removeMetric = (index: number) => {
    eventMetrics.splice(index, 1)
  }

  const updateMetric = (index: number, updated: Partial<EventMetric>) => {
    Object.assign(eventMetrics[index], updated)
  }

  const addCustomMetric = () => {
    const newMetric: any = {
      _uid: Date.now(),
      name: `自定义指标${eventMetrics.length + 1}`,
      type: EventType.Custom,
      isCustom: true,
      is_custom: true,
      custom_web_metric: {
        eventMetrics: [],
        format: FormatDefault
      },
      custom_metric: {
        event_metrics: [],
        custom_formula: '',
        format: FormatDefault
      }
    }
    eventMetrics.push(newMetric)
  }

  return {
    addMetric,
    addCustomMetric,
    removeMetric,
    updateMetric
  }
}

import { Group, TimeGrain } from './analysis'
import { Column, Formula } from './common'
import { FilterGroup } from './filter'

/**
 * 基础指标定义
 */
export interface Metric {
  column: Column
  formula: Formula
  format: string
}

/**
 * 自定义指标配置
 */
export interface CustomMetric {
  event_metrics: EventMetric[]
  custom_formula: string
  format: string
}

/**
 * 前端可视化编辑使用的公式 Token
 */
export type FormulaTokenType = 'metric' | 'operator' | 'number'

export interface FormulaToken {
  type: FormulaTokenType
  value: string // 如果是 metric，则是别名如 m0；如果是 operator/number，则是实际值
  metric?: EventMetric // 仅当 type 为 metric 时存在
}

/**
 * 前端 UI 状态管理模型
 */
export interface WebMetric {
  tokens: FormulaToken[]
  format: string
}

/**
 * 事件类型枚举
 */
export enum EventType {
  Normal = 1, // 普通事件
  Custom = 2, // 自定义事件
}

/**
 * 事件指标定义
 */
export interface EventMetric extends FilterGroup {
  _uid?: number | string // 前端唯一标识
  e_event_id: string
  name: string
  type: EventType
  metric: Metric
  custom_metric?: CustomMetric
  web_metric?: WebMetric
  is_custom?: boolean // 前端标记，是否为公式指标

  // UI 状态
  isEditingName?: boolean
  filterPopoverVisible?: boolean
  isCustom?: boolean // 兼容 old_frontend
  custom_web_metric?: {
    eventMetrics: any[] // 兼容 old_frontend 的 CustomElement[]
  }
}

/**
 * 窗口函数指标枚举
 */
export enum WindowFunc {
  FirstValue = 1,
  LastValue = 2,
  Count = 3,
  Sum = 4,
  Avg = 5,
}

/**
 * 窗口函数指标定义
 */
export interface WindowMetric {
  column: Column
  window_func: WindowFunc
  time_grain: TimeGrain
  groups: Group[]
}

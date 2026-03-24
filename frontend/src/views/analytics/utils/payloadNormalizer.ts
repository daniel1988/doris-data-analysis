import { FormulaDefaults, FormulaTypes, TableNames } from '@/constants/analysis'

/**
 * 规范化指标字段，将前端占位符（如 __TOTAL_USERS__）转换为后端实际列名
 * 
 * @param metric 指标对象，支持普通指标和自定义公式中的基础指标
 * @param selectedField 可选：强制指定的字段名（如 __TOTAL_TIMES__）
 */
export function normalizeMetricField(metric: any, selectedField?: string) {
  if (!metric || !metric.metric) return metric

  const col = metric.metric.column || (metric.metric.column = {})
  const field = selectedField ?? col.field
  if (selectedField) col.field = selectedField

  // 1. 处理“总次数”占位符
  if (field === FormulaDefaults.TOTAL_TIMES_FIELD) {
    if (col.table === TableNames.USER || col.table === TableNames.USER_PROPERTY) {
      col.field = 'u_openid'
    } else {
      col.field = 'e_openid'
    }
    if (!col.table) col.table = TableNames.EVENT
    metric.metric.formula = FormulaTypes.Count
  }

  // 2. 处理“总用户数”占位符
  if (field === FormulaDefaults.TOTAL_USERS_FIELD) {
    if (col.table === TableNames.USER || col.table === TableNames.USER_PROPERTY) {
      col.field = 'u_openid'
    } else {
      col.field = 'e_openid'
    }
    if (!col.table) col.table = TableNames.EVENT
    metric.metric.formula = FormulaTypes.CountDistinctUserId
  }

  const formula = Number(metric.metric.formula)

  // 3. 兜底逻辑：如果是计数/去重计数且字段为空，默认使用 ID 字段
  if (!col.field && (formula === FormulaTypes.Count || formula === FormulaTypes.CountDistinctUserId)) {
    if (col.table === TableNames.USER || col.table === TableNames.USER_PROPERTY) {
      col.field = 'u_openid'
    } else {
      col.field = 'e_openid'
    }
  }

  // 4. 确保 table 属性存在
  if (!col.table) {
    col.table = TableNames.EVENT
  }

  return metric
}

/**
 * 深度规范化指标对象（处理自定义公式中的所有子指标）
 */
export function deepNormalizeMetric(metric: any) {
  if (!metric) return metric

  // 如果是自定义公式指标
  if (metric.isCustom && metric.custom_metric?.event_metrics) {
    metric.custom_metric.event_metrics.forEach((m: any) => normalizeMetricField(m))
  } else {
    // 普通指标
    normalizeMetricField(metric)
  }

  return metric
}

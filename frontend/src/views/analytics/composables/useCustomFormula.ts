import { FormulaDefaults, FormulaTokenTypes, FormulaTypes, TableNames } from '@/constants/analysis'
import { CustomElement, CustomMetric } from '@/types/doris/customFormula'
import { normalizeMetricField } from '../utils/payloadNormalizer'

export function useCustomFormula() {
  /**
   * 核心同步逻辑：将前端 UI 的 Token 序列 (CustomElement[]) 转换为 API 所需的 custom_metric 结构
   * 遵循“自动索引”原则，为每个基础指标分配唯一的别名（如 m0, m1...）
   * 
   * @param m 包含 custom_web_metric 和 custom_metric 的指标对象
   */
  const syncWebMetricToApiMetric = (m: any) => {
    if (!m.custom_web_metric || !m.custom_metric) return

    const elements = m.custom_web_metric.eventMetrics || []
    const eventMetricsForApi: any[] = []
    let formulaStr = ''
    let metricCounter = 0

    elements.forEach((element: any) => {
      const t = String(element?.type).toLowerCase()

      if (t === FormulaTokenTypes.Metric || t === '1' || t === 'event') {
        // 生成别名，如 m0, m1...
        const alias = `${FormulaDefaults.METRIC_PREFIX}${metricCounter}`

        // 深度复制基础指标并注入别名作为其名称
        const metricCopy = JSON.parse(JSON.stringify(element.value))
        metricCopy.name = alias

        // 规范化字段（如总次数映射）
        normalizeMetricField(metricCopy)

        eventMetricsForApi.push(metricCopy)
        formulaStr += alias
        metricCounter++
      } else {
        // 运算符或数字直接追加到公式字符串中
        formulaStr += String(element?.value ?? '')
      }
    })

    // 更新 API Payload 结构
    m.custom_metric.event_metrics = eventMetricsForApi
    m.custom_metric.custom_formula = formulaStr
    m.custom_metric.format = m.custom_web_metric.format || FormulaDefaults.DEFAULT_FORMAT
  }

  /**
   * 将普通指标转换为自定义公式指标
   * @param m 当前指标对象
   */
  const convertToCustomMetric = (m: any) => {
    if (m.isCustom) return

    const snapshot = {
      type: 1,
      e_event_id: m.e_event_id || '',
      name: `${m.name || '指标'}_0`,
      metric: JSON.parse(JSON.stringify(m.metric || { column: { table: TableNames.EVENT, field: 'e_openid' }, formula: FormulaTypes.Count })),
      filter_group: m.filter_group ? JSON.parse(JSON.stringify(m.filter_group)) : { scope: 1, filters: [] }
    }

    m.isCustom = true
    m.type = 2 // EventType.Custom
    m.custom_web_metric = {
      eventMetrics: [{ type: FormulaTokenTypes.Metric, value: snapshot, _uid: Date.now() }],
      format: FormulaDefaults.DEFAULT_FORMAT
    }
    m.custom_metric = {
      event_metrics: [snapshot],
      custom_formula: snapshot.name,
      format: FormulaDefaults.DEFAULT_FORMAT
    }

    // 清理普通指标字段
    delete m.e_event_id
    delete m.metric
    delete m.filter_group

    syncWebMetricToApiMetric(m)
  }

  /**
   * 将自定义公式指标回退为普通指标（取第一个事件指标）
   * @param m 当前指标对象
   */
  const convertToNormalMetric = (m: any) => {
    if (!m.isCustom) return

    const firstEm = m.custom_metric?.event_metrics?.[0]
    m.isCustom = false
    m.type = 1 // EventType.Normal
    m.e_event_id = firstEm?.e_event_id || ''
    m.metric = firstEm?.metric ? JSON.parse(JSON.stringify(firstEm.metric)) : { column: { table: TableNames.EVENT, field: 'e_openid' }, formula: FormulaTypes.Count }
    m.filter_group = firstEm?.filter_group ? JSON.parse(JSON.stringify(firstEm.filter_group)) : { scope: 1, filters: [] }

    // 清理自定义指标字段
    delete m.custom_web_metric
    delete m.custom_metric
  }

  /**
   * 将前端的 WebMetric 转换为后端所需的 CustomMetric (兼容旧调用)
   */
  const buildCustomMetric = (webMetric: { tokens: CustomElement[], format: string }): CustomMetric => {
    const mockMetric = {
      custom_web_metric: {
        eventMetrics: webMetric.tokens,
        format: webMetric.format
      },
      custom_metric: {
        event_metrics: [],
        custom_formula: '',
        format: ''
      }
    }
    syncWebMetricToApiMetric(mockMetric)
    return mockMetric.custom_metric
  }

  /**
   * 自动生成自定义指标的展示名称
   * 将 Token 序列转换为可读性较高的字符串，如 "支付成功(求和) / 100"
   */
  const generateFormulaName = (tokens: CustomElement[]): string => {
    if (!tokens || tokens.length === 0) return FormulaDefaults.UNKNOWN_FORMULA

    return tokens
      .map((t) => {
        if ((t.type === FormulaTokenTypes.Metric || t.type === '1' || t.type === 'event') && t.value) {
          // 如果有具体的指标名，则使用指标名，否则回退到基础名称
          return t.value.name || '指标'
        }
        return t.value
      })
      .join(' ')
  }

  return {
    syncWebMetricToApiMetric,
    buildCustomMetric,
    generateFormulaName,
    normalizeMetricField,
    convertToCustomMetric,
    convertToNormalMetric
  }
}

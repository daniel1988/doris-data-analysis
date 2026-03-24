/**
 * 自定义公式指标相关类型定义
 */

// 公式元素类型
export type CustomElementType = '1' | '2' | '3' | 'event' | 'number' | 'operator' | 'metric' // 1-事件指标, 2-数字, 3-运算符, metric-指标(新)

// 公式元素结构
export interface CustomElement {
    type: CustomElementType
    value: any // 根据类型存储不同值
    _uid?: number // 前端唯一标识
}

// 事件指标结构
export interface EventMetric {
    type: number
    eventId?: string // 兼容字段
    e_event_id: string // 事件ID
    name: string // 指标名称
    metric: {
        column: {
            table: string
            field: string
            alias?: string
        }
        formula: number // 指标类型：1-求和, 2-总次数, 3-去重数, 4-用户数等
        format: string // 数值格式：raw, int, decimal, percent
    }
    filter_group?: {
        scope: number // 1-AND, 2-OR
        filters: FilterCondition[]
    }
    // 前端显示相关属性
    showFilter?: boolean
    isEditing?: boolean
    _uid?: number
}

// 过滤条件
export interface FilterCondition {
    field: string
    operator: string
    value: any
    type?: string
}

// 前端自定义公式指标整体结构
export interface CustomWebMetric {
    eventMetrics: CustomElement[] // 公式元素数组（包含事件指标、数字、运算符）
    formula?: string // 完整的公式表达式（可选，由元素自动生成）
    format?: string // 数值格式化方式
}

// API请求中的自定义指标结构
export interface CustomMetric {
    event_metrics: EventMetric[] // 事件指标列表
    custom_formula: string // 自定义公式表达式
    format?: string // 数值格式化方式
}

// 完整的指标结构（包含普通指标和自定义指标）
export interface MetricConfig {
    _uid: number
    type: number // 1-普通指标, 2-自定义指标
    name: string

    // 普通指标字段
    e_event_id?: string
    metric?: {
        column: {
            table: string
            field: string
            alias?: string
        }
        formula: number
        format: string
    }
    filter_group?: {
        scope: number
        filters: FilterCondition[]
    }

    // 自定义指标标识
    isCustom?: boolean

    // 自定义指标字段
    custom_web_metric?: CustomWebMetric // 前端展示结构
    custom_metric?: CustomMetric // API请求结构

    // 前端状态字段
    isEditingName?: boolean
    isEditingProperty?: boolean
    isEditingFormula?: boolean
    filterPopoverVisible?: boolean
    showFilter?: boolean
}

// 事件选项
export interface EventOption {
    event_id: string
    event_name?: string
    event_type?: string
}

// 公式运算符定义
export interface FormulaOperator {
    label: string
    value: string
    title: string
    category: 'arithmetic' | 'parentheses' | 'function'
}

export const FORMULA_OPERATORS: FormulaOperator[] = [
    { label: '+', value: '+', title: '加', category: 'arithmetic' },
    { label: '-', value: '-', title: '减', category: 'arithmetic' },
    { label: '×', value: '*', title: '乘', category: 'arithmetic' },
    { label: '÷', value: '/', title: '除', category: 'arithmetic' },
    { label: '(', value: '(', title: '左括号', category: 'parentheses' },
    { label: ')', value: ')', title: '右括号', category: 'parentheses' }
]

// 数值格式选项
export interface FormatOption {
    label: string
    value: string
    description: string
}

export const FORMAT_OPTIONS: FormatOption[] = [
    { label: '原始', value: 'raw', description: '显示原始数值' },
    { label: '整数', value: 'int', description: '显示为整数' },
    { label: '小数(2位)', value: 'decimal', description: '显示2位小数' },
    { label: '百分比', value: 'percent', description: '显示为百分比' }
]

// 指标公式类型映射
export const METRIC_FORMULA_MAP: Record<number, string> = {
    1: '求和',
    2: '总次数',
    3: '去重数',
    4: '用户数',
    5: '任意值',
    6: '最大值',
    7: '最小值',
    8: '平均值'
}

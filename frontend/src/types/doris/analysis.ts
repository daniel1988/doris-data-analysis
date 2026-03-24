import { Column, Order } from './common'
import { GlobalFilterGroups, TagFilter, UserGroupDataFilter } from './filter'
import { EventMetric } from './metric'
import { OperatorTypes } from '@/constants/analysis'

/**
 * 分析指标 (别名，确保接口规范)
 */
export type AnalysisMetric = EventMetric

/**
 * 分析过滤器
 */
export interface AnalysisFilter {
    table: string
    field: string
    operator: OperatorTypes
    value: any
    type: string // 字段类型，如 string, number, date
}

/**
 * 时间粒度枚举
 */
export enum TimeGrainInterval {
    Empty = 1,
    Day = 2,
    Week = 3,
    Month = 4,
    Quarter = 5,
    Year = 6,
    Hour = 7,
    Minute = 8,
}

/**
 * 时间粒度配置
 */
export interface TimeGrain {
    column: Column
    interval: TimeGrainInterval
    window_num: number
}

/**
 * 分组类型枚举
 */
export enum GroupType {
    Value = 0, // 按值
    Ranges = 1, // 按照区间
    Date = 2, // 按照日期
    TagGroup = 3, //按照标签分组
    UserGroupData = 4, //按照用户分群
    UserGroups = 5, // 用户属性分群
}

/**
 * 值范围定义
 */
export interface ValueRange {
    min: number
    max: number
}

/**
 * 分组配置
 */
export interface Group {
    group_type: GroupType
    column: Column
    value_ranges: ValueRange[]
    time_grain: TimeGrain
    tag_group: TagFilter
    user_group: UserGroupDataFilter
}

/**
 * 分析分组 (别名，确保接口规范)
 */
export type AnalysisGroup = Group

/**
 * 基础请求模型
 */
export interface QueryRequest {
    project_alias: string
    sql?: string
    time_zone?: string
    page_size?: number
    page_num?: number
    query_id?: string
    user_id?: string
}

/**
 * 基础响应模型
 */
export interface QueryResponse<T = any> {
    rows: T[]
    columns: string[]
    count: number
    sql: string
    duration: number
}

/**
 * 事件分析请求
 */
export interface EventAnalysisReq extends QueryRequest {
    event_metrics: EventMetric[]
    filter_groups: GlobalFilterGroups
    time_grain: TimeGrain
    groups: Group[]
    orders: Order[]
}

/**
 * 留存分析请求
 */
export interface RetentionAnalysisReq extends QueryRequest {
    init_event_metric: EventMetric
    end_event_metric: EventMetric
    global_filter_groups: GlobalFilterGroups
    groups: Group[]
    time_grain: TimeGrain
    day_n_array: number[]
}

/**
 * 漏斗分析请求
 */
export interface FunnelAnalysisReq extends QueryRequest {
    steps: EventMetric[]
    global_filter_groups: GlobalFilterGroups
    groups: Group[]
    window: number
    time_grain: TimeGrain
}

/**
 * 分布分析请求
 */
export interface ScatterAnalysisReq extends QueryRequest {
    metric: EventMetric
    global_filter_groups: GlobalFilterGroups
    groups: Group[]
    time_grain: TimeGrain
    scatter_type?: number
    scatter_ranges?: Array<{ min: number; max: number }>
    bin_count?: number
}

/**
 * 用户属性分析请求
 */
export interface UserPropertyAnalysisReq extends QueryRequest {
    metric: Metric
    filter_groups: GlobalFilterGroups
    groups: Group[]
    user_groups: Array<{ alias: string; filter_group: GlobalFilterGroups }>
    group_type: number // 1: 维度分组, 2: 人群分群
}

/**
 * 分析请求负载统一结构
 */
export type AnalysisPayload = EventAnalysisReq | FunnelAnalysisReq | RetentionAnalysisReq | ScatterAnalysisReq | UserPropertyAnalysisReq

/**
 * 用户列表请求
 */
export interface UserListReq extends QueryRequest {
    global_filter_groups: GlobalFilterGroups
    columns: Column[]
}

/**
 * 维度查询请求
 */
export interface DimensionReq {
    project_alias: string
    table: string
    field: string
    e_event_id?: string
}

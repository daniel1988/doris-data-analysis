import { Column, Operator } from './common'

/**
 * 过滤器定义
 */
export interface Filter {
  column: Column
  value: any
  operator: Operator
}

/**
 * 过滤器范围枚举
 */
export enum FilterScope {
  Or = 0,
  And = 1,
}

/**
 * 标签过滤器定义
 */
export interface TagFilter {
  tag_code: string
  operator: Operator
  tag_value: any
}

/**
 * 用户分群数据过滤器
 */
export interface UserGroupDataFilter {
  group_name: string
  group_code: string
  operator: Operator
}

/**
 * 过滤组定义
 */
export interface FilterGroup {
  scope: FilterScope
  filters: Filter[]
  tag_filters: TagFilter[]
  user_group_filters: UserGroupDataFilter[]
}

/**
 * 全局过滤组
 */
export interface GlobalFilterGroups {
  global_filters: FilterGroup
  dashboard_form_filters: FilterGroup
  query_dates: string[]
  comparison_query_dates?: string[]
}

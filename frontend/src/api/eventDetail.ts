import request from '@/utils/request'

// 事件明细查询参数接口
export interface EventDetailRequest {
  project_alias: string
  page_size?: number
  page_num?: number
  select_fields?: string[]
  event_filter_group?: FilterGroup | null
}

// 过滤组接口
export interface FilterGroup {
  scope: number // 0: OR, 1: AND
  filters: Filter[]
}

// 过滤条件接口
export interface Filter {
  column: {
    table?: string
    field: string
    alias?: string
  }
  operator: number // 操作符枚举值
  value: any
}

// 事件明细数据接口
export interface EventDetailItem {
  e_openid: string         // 事件用户ID
  e_event_id: string       // 事件ID
  e_event_name: string     // 事件名称
  e_event_time: string     // 事件时间
  e_package_name: string   // 事件应用包名
  e_platform: string       // 事件平台
  e_ip: string             // IP 地址
  e_request_id: string     // 请求 ID
  e_properties?: string    // 自定义属性 JSON 字符串
}

// 操作符标签映射
export const OperatorLabels: Record<number, string> = {
  1: '等于',
  2: '不等于',
  3: '小于',
  4: '小于等于',
  5: '大于',
  6: '大于等于',
  7: '为空',
  8: '不为空',
  9: '区间',
  10: '包含',
  11: '不包含',
  12: '模糊匹配',
  13: '日期差异',
  14: '不模糊匹配',
  15: '属于',
  16: '不属于',
  17: '以...开始',
  18: '以...结束',
  19: '动态日期'
}

/**
 * 查询事件明细
 */
export function getEventDetail(data: EventDetailRequest) {
  return request<{
    rows: EventDetailItem[]
    count: number
    columns?: string[]
    sql?: string
    duration?: number
  }>({
    url: '/v1/event-detail',
    method: 'post',
    data
  })
}

/**
 * 获取事件属性选项（用于过滤器下拉框）
 */
export function getEventProperties(projectAlias: string) {
  return request<Array<{
    id: string
    name: string
    data_type: string
  }>>({
    url: '/v1/selector/properties',
    method: 'get',
    params: { project_alias: projectAlias }
  })
}

/**
 * 根据项目别名获取事件列表（用于事件ID下拉框）
 */
export function getProjectEvents(projectAlias: string) {
  return request<Array<{
    id: string
    name: string
  }>>({
    url: '/v1/selector/events',
    method: 'get',
    params: { project_alias: projectAlias }
  })
}

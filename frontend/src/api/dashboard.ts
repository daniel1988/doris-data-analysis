import request from '@/utils/request'

export interface DashboardItem {
  id?: number
  dashboard_id?: number
  report_id: number
  type: string
  title: string
  position_x: number
  position_y: number
  width: number
  height: number
  z_index?: number
  config_override?: string
  is_visible?: boolean
  create_time?: string
  update_time?: string
}

export interface Dashboard {
  id?: number
  project_alias: string
  name: string
  display_name: string
  description: string
  category: string
  layout_type: string
  grid_config?: string
  theme?: string
  refresh_interval?: number
  filters?: string
  variables?: string
  status: string
  owner_id?: number
  create_time?: string
  update_time?: string
  items?: DashboardItem[]
}

export function getDashboardList(params: { project_alias: string }) {
  return request<Dashboard[]>({
    url: '/v1/dashboard/list',
    method: 'get',
    params
  })
}

export function getDashboardById(id: number) {
  return request<Dashboard>({
    url: `/v1/dashboard/${id}`,
    method: 'get'
  })
}

export function createDashboard(data: Dashboard) {
  return request<Dashboard>({
    url: '/v1/dashboard',
    method: 'post',
    data
  })
}

export function updateDashboard(data: Dashboard) {
  return request({
    url: '/v1/dashboard',
    method: 'put',
    data
  })
}

export function deleteDashboard(id: number) {
  return request({
    url: `/v1/dashboard/${id}`,
    method: 'delete'
  })
}

export function addDashboardItem(data: DashboardItem) {
  return request<DashboardItem>({
    url: '/v1/dashboard/item',
    method: 'post',
    data
  })
}

export function updateDashboardItem(data: DashboardItem) {
  return request({
    url: '/v1/dashboard/item',
    method: 'put',
    data
  })
}

export function deleteDashboardItem(id: number) {
  return request({
    url: `/v1/dashboard/item/${id}`,
    method: 'delete'
  })
}

export function batchUpdateDashboardItems(data: DashboardItem[]) {
  return request({
    url: '/v1/dashboard/items/batch',
    method: 'put',
    data
  })
}

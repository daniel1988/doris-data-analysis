import request from '@/utils/request'

export interface Report {
  id?: number
  project_alias: string
  name: string
  category: string
  description: string
  query_params: string
  create_user?: number
  create_time?: string
  update_time?: string
}

export function getReportList(params: { project_alias: string; category?: string }) {
  return request<Report[]>({
    url: '/v1/report/list',
    method: 'get',
    params
  })
}

export function getReportById(id: number) {
  return request<Report>({
    url: `/v1/report/${id}`,
    method: 'get'
  })
}

export function createReport(data: Report) {
  return request<Report>({
    url: '/v1/report',
    method: 'post',
    data
  })
}

export function updateReport(data: Report) {
  return request({
    url: '/v1/report',
    method: 'put',
    data
  })
}

export function deleteReport(id: number) {
  return request({
    url: `/v1/report/${id}`,
    method: 'delete'
  })
}

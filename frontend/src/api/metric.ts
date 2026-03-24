import { useAppStore } from '@/store/app'
import request from '@/utils/request'

export interface ProjectMetric {
  id?: number
  project_alias?: string
  metric_name: string
  metric_code: string
  expression: string
  base_table: string
  description: string
  status: number
  create_time?: string
  update_time?: string
}

const getHeaders = () => {
  const appStore = useAppStore()
  return {
    'X-Project-Alias': appStore.activeProjectAlias
  }
}

export function getMetrics() {
  return request({
    url: '/v1/meta/metrics',
    method: 'get',
    headers: getHeaders()
  })
}

export function createMetric(data: ProjectMetric) {
  return request({
    url: '/v1/meta/metrics',
    method: 'post',
    data,
    headers: getHeaders()
  })
}

export function updateMetric(metricCode: string, data: ProjectMetric) {
  return request({
    url: `/v1/meta/metrics/${metricCode}`,
    method: 'put',
    data,
    headers: getHeaders()
  })
}

export function deleteMetric(metricCode: string) {
  return request({
    url: `/v1/meta/metrics/${metricCode}`,
    method: 'delete',
    headers: getHeaders()
  })
}

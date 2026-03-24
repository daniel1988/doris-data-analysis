import request from '@/utils/request'

export interface DimensionReq {
  project_alias: string
  table: string
  field: string
  e_event_id?: string
}

export const dimensionsApi = {
  listValues(params: DimensionReq) {
    return request({
      url: '/v1/analytics/dimension',
      method: 'post',
      data: params
    })
  }
}

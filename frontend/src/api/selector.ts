import request from '@/utils/request'

export interface Option {
  id: string
  name: string
  data_type?: string
}

export function getEvents(params: { project_alias: string }) {
  return request<any[]>({
    url: '/v1/selector/events',
    method: 'get',
    params
  })
}

export function getProperties(params: { project_alias: string; event_id?: string }) {
  return request<any[]>({
    url: '/v1/selector/properties',
    method: 'get',
    params
  })
}

export function getEventOptions(projectAlias: string) {
  return getEvents({ project_alias: projectAlias }) as Promise<Option[]>
}

export function getPropertyOptions(projectAlias: string, eventId?: string) {
  return getProperties({
    project_alias: projectAlias,
    event_id: eventId
  }) as Promise<Option[]>
}

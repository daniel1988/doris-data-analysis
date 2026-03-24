import request from '@/utils/request'

export interface ProjectProperty {
    project_alias: string
    property_id: string
    property_name: string
    data_type: string
    property_type: number
    create_time?: string
    update_time?: string
}

export interface ProjectEvent {
    project_alias: string
    event_id: string
    event_name: string
    event_type: number
    create_time?: string
    update_time?: string
}

export function getProjectEvents(projectAlias: string) {
    return request<ProjectEvent[]>({
        url: '/v1/meta/project_events',
        method: 'get',
        params: { project_alias: projectAlias }
    })
}

export function createProjectEvent(data: ProjectEvent) {
    return request({
        url: '/v1/meta/project_events',
        method: 'post',
        data
    })
}

export function updateProjectEvent(data: ProjectEvent) {
    return request({
        url: '/v1/meta/project_events',
        method: 'put',
        data
    })
}

export function deleteProjectEvent(projectAlias: string, eventId: string) {
    return request({
        url: `/v1/meta/project_events/${eventId}`,
        method: 'delete',
        params: { project_alias: projectAlias }
    })
}

export function getProjectProperties(projectAlias: string) {
    return request<ProjectProperty[]>({
        url: '/v1/meta/project_properties',
        method: 'get',
        params: { project_alias: projectAlias }
    })
}

export function createProjectProperty(data: ProjectProperty) {
    return request({
        url: '/v1/meta/project_properties',
        method: 'post',
        data
    })
}

export function updateProjectProperty(data: ProjectProperty) {
    return request({
        url: '/v1/meta/project_properties',
        method: 'put',
        data
    })
}

export function deleteProjectProperty(projectAlias: string, propertyId: string) {
    return request({
        url: `/v1/meta/project_properties/${propertyId}`,
        method: 'delete',
        params: { project_alias: projectAlias }
    })
}

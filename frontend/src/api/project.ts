import request from '@/utils/request'

export interface ProjectData {
  project_alias: string
  project_name: string
  region: string
  secret: string
  create_time?: string
  update_time?: string
}

export function getProjects() {
  return request<ProjectData[]>({
    url: '/v1/system/projects',
    method: 'get'
  })
}

export function createProject(data: ProjectData) {
  return request({
    url: '/v1/system/projects',
    method: 'post',
    data
  })
}

export function updateProject(data: ProjectData) {
  return request({
    url: '/v1/system/projects',
    method: 'put',
    data
  })
}

export function deleteProject(projectAlias: string) {
  return request({
    url: `/v1/system/projects/${projectAlias}`,
    method: 'delete'
  })
}

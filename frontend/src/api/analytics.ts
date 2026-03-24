import type {
  DimensionReq,
  EventAnalysisReq,
  FunnelAnalysisReq,
  QueryResponse,
  RetentionAnalysisReq,
  ScatterAnalysisReq,
  UserListReq,
  UserPropertyAnalysisReq
} from '@/types/doris/analysis'
import request from '@/utils/request'

export function eventAnalysis(data: EventAnalysisReq) {
  return request<QueryResponse>({
    url: '/v1/analytics/event',
    method: 'post',
    data
  })
}

export function funnelAnalysis(data: FunnelAnalysisReq) {
  return request<QueryResponse>({
    url: '/v1/analytics/funnel',
    method: 'post',
    data
  })
}

export function retentionAnalysis(data: RetentionAnalysisReq) {
  return request<QueryResponse>({
    url: '/v1/analytics/retention',
    method: 'post',
    data
  })
}

export function scatterAnalysis(data: ScatterAnalysisReq) {
  return request<QueryResponse>({
    url: '/v1/analytics/scatter',
    method: 'post',
    data
  })
}

export function userPropertyAnalysis(data: UserPropertyAnalysisReq) {
  return request<QueryResponse>({
    url: '/v1/analytics/user_property',
    method: 'post',
    data
  })
}

export function getUserList(data: UserListReq) {
  return request<QueryResponse>({
    url: '/v1/analytics/user_list',
    method: 'post',
    data
  })
}

export function getDimensions(data: DimensionReq) {
  return request<any[]>({
    url: '/v1/analytics/dimension',
    method: 'post',
    data
  })
}

/**
 * 获取用户属性列表
 */
export function getUserProperties(projectAlias: string) {
  return request({
    url: '/v1/analytics/user_properties',
    method: 'get',
    params: { project_alias: projectAlias }
  })
}

/**
 * 获取用户标签列表
 */
export function getUserTags(projectAlias: string) {
  return request({
    url: '/v1/analytics/user_tags',
    method: 'get',
    params: { project_alias: projectAlias }
  })
}

/**
 * 获取标签候选值
 */
export function getTagValues(projectAlias: string, tagCode: string) {
  return request({
    url: '/v1/analytics/tag_values',
    method: 'get',
    params: { project_alias: projectAlias, tag_code: tagCode }
  })
}

/**
 * AI Chat
 */
export function aiChat(data: { 
  project_alias: string; 
  query: string; 
  model_id?: number;
  mode?: 'sql' | 'chat';
  context_data?: any;
}) {
  return request({
    url: '/v1/analytics/ai/chat',
    method: 'post',
    data
  })
}

// ================= AI Session 相关接口 =================

export interface AISession {
  id?: number
  project_alias?: string
  user_query: string
  llm_sql: string
  viz_type: string
  x_axis: string
  y_axis: string
  narrative: string
  create_time?: string
}

export function saveAISession(data: AISession) {
  return request({
    url: '/v1/analytics/ai/sessions',
    method: 'post',
    data
  })
}

export function getAISessions(projectAlias: string) {
  return request({
    url: '/v1/analytics/ai/sessions',
    method: 'get',
    params: { project_alias: projectAlias }
  })
}

export function deleteAISession(id: number) {
  return request({
    url: `/v1/analytics/ai/sessions/${id}`,
    method: 'delete'
  })
}

export function executeAISession(id: number, projectAlias: string) {
  return request({
    url: `/v1/analytics/ai/sessions/${id}/execute`,
    method: 'post',
    params: { project_alias: projectAlias }
  })
}

/**
 * 获取用户分群列表
 */
export function getUserGroups(projectAlias: string) {
  return request({
    url: '/v1/analytics/user_groups',
    method: 'get',
    params: { project_alias: projectAlias }
  })
}

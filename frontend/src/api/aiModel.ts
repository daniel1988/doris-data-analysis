import request from '@/utils/request'

export interface AIModelConfig {
  id?: number
  provider: string
  display_name: string
  base_url: string
  api_key: string
  model_name: string
  max_tokens: number
  temperature: number
  is_default: boolean
  is_enabled: boolean
  sort_order: number
  create_time?: string
  update_time?: string
}

export interface AIModelBrief {
  id: number
  provider: string
  display_name: string
  is_default: boolean
}

// ===== 管理端 API =====

export function getAIModels() {
  return request<AIModelConfig[]>({
    url: '/v1/system/ai-models',
    method: 'get'
  })
}

export function createAIModel(data: AIModelConfig) {
  return request({
    url: '/v1/system/ai-models',
    method: 'post',
    data
  })
}

export function updateAIModel(id: number, data: AIModelConfig) {
  return request({
    url: `/v1/system/ai-models/${id}`,
    method: 'put',
    data
  })
}

export function deleteAIModel(id: number) {
  return request({
    url: `/v1/system/ai-models/${id}`,
    method: 'delete'
  })
}

export function testAIModel(id: number) {
  return request<{ message: string }>({
    url: `/v1/system/ai-models/${id}/test`,
    method: 'post'
  })
}

// ===== 前端可用模型列表 =====

export function getEnabledAIModels() {
  return request<AIModelBrief[]>({
    url: '/v1/analytics/ai/models',
    method: 'get'
  })
}

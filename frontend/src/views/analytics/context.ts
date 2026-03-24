import { InjectionKey, Ref } from 'vue'
import { Option } from '@/api/selector'
import { AnalysisPayload } from '@/types/doris/analysis'
import { VizModes } from '@/constants/analysis'

/**
 * 分析元数据结构
 */
export interface AnalysisMetadata {
  eventOptions: Option[]
  propertyOptions: Option[]
}

/**
 * 分析上下文状态
 */
export interface AnalysisContextState {
  projectAlias: string
  metadata: AnalysisMetadata
  loading: boolean
  form: AnalysisPayload
  results: any[]
  columns: string[]
  sql: string
  total: number
  vizMode: VizModes
}

/**
 * 分析上下文操作
 */
export interface AnalysisContextActions {
  refreshMetadata: () => Promise<void>
  triggerAnalyze: () => Promise<void>
  syncTimeRange: (staticRange: string[] | null, comparisonRange: string[] | null) => void
}

/**
 * 完整的分析上下文
 */
export interface AnalysisContext {
  state: AnalysisContextState
  actions: AnalysisContextActions
}

/**
 * Context Injection Key
 */
export const ANALYSIS_CONTEXT_KEY: InjectionKey<AnalysisContext> = Symbol('AnalysisContext')

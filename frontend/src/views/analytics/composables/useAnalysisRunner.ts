import { ref } from 'vue'
import { eventAnalysis } from '@/api/analytics'
import type { AnalysisPayload, QueryResponse } from '@/types/doris/analysis'
import { ElMessage } from 'element-plus'

export function useAnalysisRunner() {
  const loading = ref(false)
  const results = ref<any[]>([])
  const columns = ref<string[]>([])
  const total = ref(0)
  const sql = ref('')
  const duration = ref(0)

  const runAnalysis = async (payload: AnalysisPayload) => {
    loading.value = true
    try {
      // Basic implementation, can be extended for other analysis types
      const resp = await eventAnalysis(payload as any)
      const data = resp as unknown as QueryResponse
      
      // Independent unpacking logic
      results.value = data.rows || []
      columns.value = data.columns || []
      total.value = data.count || 0
      sql.value = data.sql || ''
      duration.value = data.duration || 0
      
      return data
    } catch (error: any) {
      console.error('Analysis failed:', error)
      ElMessage.error(error.message || '分析执行失败')
      return null
    } finally {
      loading.value = false
    }
  }

  const clearResults = () => {
    results.value = []
    columns.value = []
    total.value = 0
    sql.value = ''
    duration.value = 0
  }

  return {
    loading,
    results,
    columns,
    total,
    sql,
    duration,
    runAnalysis,
    clearResults
  }
}

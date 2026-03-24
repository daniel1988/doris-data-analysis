import { funnelAnalysis } from '@/api/analytics'
import type { QueryResponse, FunnelAnalysisReq } from '@/types/doris/analysis'
import { ElMessage } from 'element-plus'
import { ref } from 'vue'

export function useFunnelAnalysis() {
  const loading = ref(false)
  const results = ref<any[]>([])
  const columns = ref<string[]>([])
  const total = ref(0)
  const sql = ref('')

  const runFunnelAnalysis = async (payload: FunnelAnalysisReq) => {
    loading.value = true
    try {
      const resp = await funnelAnalysis(payload)
      const data = resp as unknown as QueryResponse
      
      results.value = data.rows || []
      columns.value = data.columns || []
      total.value = data.count || 0
      sql.value = data.sql || ''
      
      return data
    } catch (error: any) {
      console.error('Funnel analysis failed:', error)
      ElMessage.error(error.message || '漏斗分析执行失败')
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
  }

  return {
    loading,
    results,
    columns,
    total,
    sql,
    runFunnelAnalysis,
    clearResults
  }
}

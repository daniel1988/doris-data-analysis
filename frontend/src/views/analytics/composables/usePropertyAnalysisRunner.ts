import { userPropertyAnalysis } from '@/api/analytics'
import type { QueryResponse, UserPropertyAnalysisReq } from '@/types/doris/analysis'
import { ElMessage } from 'element-plus'
import { ref } from 'vue'

export function usePropertyAnalysisRunner() {
  const loading = ref(false)
  const results = ref<any[]>([])
  const columns = ref<string[]>([])
  const total = ref(0)
  const sql = ref('')

  const runPropertyAnalysis = async (payload: UserPropertyAnalysisReq) => {
    loading.value = true
    try {
      const resp = await userPropertyAnalysis(payload)
      const data = resp as unknown as QueryResponse
      
      results.value = data.rows || []
      columns.value = data.columns || []
      total.value = data.count || 0
      sql.value = data.sql || ''
      
      return data
    } catch (error: any) {
      console.error('Property analysis failed:', error)
      ElMessage.error(error.message || '属性分析执行失败')
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
    runPropertyAnalysis,
    clearResults
  }
}

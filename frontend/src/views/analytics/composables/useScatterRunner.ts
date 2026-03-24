import { scatterAnalysis } from '@/api/analytics'
import type { QueryResponse, ScatterAnalysisReq } from '@/types/doris/analysis'
import { ElMessage } from 'element-plus'
import { ref } from 'vue'

export function useScatterRunner() {
  const loading = ref(false)
  const results = ref<any[]>([])
  const columns = ref<string[]>([])
  const sql = ref('')
  const total = ref(0)
  const duration = ref(0)

  const runAnalysis = async (payload: ScatterAnalysisReq) => {
    loading.value = true
    try {
      const res = await scatterAnalysis(payload)
      results.value = res.rows || []
      columns.value = res.columns || []
      sql.value = res.sql || ''
      duration.value = res.duration || 0
      total.value = results.value.length
      return res
    } catch (error: any) {
      console.error('Scatter Analysis Error:', error)
      ElMessage.error(error.message || '分析请求失败')
      return null
    } finally {
      loading.value = false
    }
  }

  const clearResults = () => {
    results.value = []
    columns.value = []
    sql.value = ''
    total.value = 0
    duration.value = 0
  }

  return {
    loading,
    results,
    columns,
    sql,
    total,
    duration,
    runAnalysis,
    clearResults
  }
}

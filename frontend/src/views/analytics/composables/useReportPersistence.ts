import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { createReport, updateReport, getReportList, deleteReport } from '@/api/report'
import type { Report } from '@/api/report'
import { addDashboardItem } from '@/api/dashboard'

export function useReportPersistence() {
  const currentReport = ref<Report | null>(null)
  const reportList = ref<Report[]>([])
  const loading = ref(false)

  const saveReport = async (data: { 
    name: string; 
    description: string; 
    category: string; 
    queryParams: any; 
    projectAlias: string;
    dashboardId?: number
  }) => {
    const payload: Report = {
      project_alias: data.projectAlias,
      name: data.name,
      category: data.category,
      description: data.description,
      query_params: JSON.stringify(data.queryParams)
    }

    try {
      let savedReport: Report
      if (currentReport.value?.id) {
        payload.id = currentReport.value.id
        await updateReport(payload)
        savedReport = { ...payload }
        ElMessage.success('报表更新成功')
      } else {
        const res = await createReport(payload)
        savedReport = res as unknown as Report
        currentReport.value = savedReport
        ElMessage.success('报表保存成功')
      }

      // 如果选择了看板，则添加到看板
      if (data.dashboardId && savedReport.id) {
        await addDashboardItem({
          dashboard_id: data.dashboardId,
          report_id: savedReport.id,
          type: 'chart',
          title: savedReport.name,
          position_x: 0,
          position_y: 0,
          width: 12,
          height: 8,
          is_visible: true
        })
        ElMessage.success('已添加至看板')
      }
      return true
    } catch (error) {
      console.error('Failed to save report:', error)
      ElMessage.error('保存失败')
      return false
    }
  }

  const fetchReportList = async (projectAlias: string, category: string) => {
    loading.value = true
    try {
      const res = await getReportList({ project_alias: projectAlias, category })
      reportList.value = res as unknown as Report[]
    } catch (error) {
      console.error('Failed to fetch reports:', error)
    } finally {
      loading.value = false
    }
  }

  const removeReport = async (id: number) => {
    try {
      await ElMessageBox.confirm('确定删除该报表吗？', '提示', { type: 'warning' })
      await deleteReport(id)
      ElMessage.success('删除成功')
      return true
    } catch (error) {
      return false
    }
  }

  return {
    currentReport,
    reportList,
    loading,
    saveReport,
    fetchReportList,
    removeReport
  }
}

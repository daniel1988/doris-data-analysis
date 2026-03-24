import {
  getEventDetail,
  getEventProperties,
  getProjectEvents,
  type EventDetailItem,
  type EventDetailRequest,
  type Filter
} from '@/api/eventDetail'
import { useAppStore } from '@/store/app'
import { ElMessage } from 'element-plus'
import { computed, reactive, ref } from 'vue'

export function useEventDetail() {
  const appStore = useAppStore()
  const loading = ref(false)
  const tableData = ref<EventDetailItem[]>([])
  const availableColumns = ref<string[]>([])

  const filterDialogVisible = ref(false)
  const propertyDrawerVisible = ref(false)
  const currentProperties = ref('')

  const allEventProperties = ref<any[]>([])
  const projectEvents = ref<any[]>([])
  const activeFilters = ref<Filter[]>([])

  const searchForm = reactive({
    openId: '',
    eventId: '',
    eventTimeRange: [] as string[]
  })

  const pagination = reactive({
    page: 1,
    pageSize: 20,
    total: 0
  })

  const dateShortcuts = [
    {
      text: '今天',
      value: () => {
        const now = new Date()
        const start = new Date(now.getFullYear(), now.getMonth(), now.getDate())
        return [start, now]
      }
    },
    {
      text: '昨天',
      value: () => {
        const now = new Date()
        const start = new Date(now.getFullYear(), now.getMonth(), now.getDate() - 1)
        const end = new Date(now.getFullYear(), now.getMonth(), now.getDate() - 1, 23, 59, 59)
        return [start, end]
      }
    },
    {
      text: '近7天',
      value: () => {
        const end = new Date()
        const start = new Date()
        start.setDate(end.getDate() - 6)
        start.setHours(0, 0, 0, 0)
        return [start, end]
      }
    },
    {
      text: '近30天',
      value: () => {
        const end = new Date()
        const start = new Date()
        start.setDate(end.getDate() - 29)
        start.setHours(0, 0, 0, 0)
        return [start, end]
      }
    }
  ]

  const formattedProperties = computed(() => {
    if (!currentProperties.value) return ''
    try {
      const parsed = typeof currentProperties.value === 'string'
        ? JSON.parse(currentProperties.value)
        : currentProperties.value
      return JSON.stringify(parsed, null, 2)
    } catch {
      return currentProperties.value
    }
  })

  const fetchMeta = async () => {
    if (!appStore.activeProjectAlias) return
    try {
      const [props, events] = await Promise.all([
        getEventProperties(appStore.activeProjectAlias),
        getProjectEvents(appStore.activeProjectAlias)
      ])
      allEventProperties.value = props as any
      projectEvents.value = events as any
    } catch (error) {
      console.error('Failed to fetch metadata:', error)
    }
  }

  const loadData = async () => {
    if (!appStore.activeProjectAlias) return
    loading.value = true
    try {
      const filters: Filter[] = [...activeFilters.value]

      if (searchForm.openId) {
        filters.push({ column: { field: 'e_openid' }, operator: 1, value: searchForm.openId })
      }
      if (searchForm.eventId) {
        filters.push({ column: { field: 'e_event_id' }, operator: 1, value: searchForm.eventId })
      }
      if (searchForm.eventTimeRange && searchForm.eventTimeRange.length === 2) {
        filters.push({ column: { field: 'e_event_time' }, operator: 9, value: searchForm.eventTimeRange })
      }

      const params: EventDetailRequest = {
        project_alias: appStore.activeProjectAlias,
        page_num: pagination.page,
        page_size: pagination.pageSize,
        select_fields: ['e_event_time', 'e_openid', 'e_event_id', 'e_event_name', 'e_package_name', 'e_platform', 'e_ip', 'e_request_id', 'e_properties'],
        event_filter_group: filters.length > 0 ? { scope: 1, filters } : null
      }

      const res = await getEventDetail(params)
      tableData.value = res.rows || []
      pagination.total = res.count || 0
      if (tableData.value.length > 0) {
        availableColumns.value = res.columns || Object.keys(tableData.value[0])
      }
    } catch (error: any) {
      ElMessage.error(error.message || '查询失败')
    } finally {
      loading.value = false
    }
  }

  const handleSearch = () => {
    pagination.page = 1
    loadData()
  }

  const handlePageChange = (page: number) => {
    pagination.page = page
    loadData()
  }

  const handleSizeChange = (size: number) => {
    pagination.pageSize = size
    pagination.page = 1
    loadData()
  }

  const viewProperties = (props: string) => {
    currentProperties.value = props
    propertyDrawerVisible.value = true
  }

  const copyToClipboard = (text: string, label: string) => {
    if (!text) return
    navigator.clipboard.writeText(text).then(() => {
      ElMessage.success(`${label}已复制`)
    }).catch(() => {
      ElMessage.error('复制失败')
    })
  }

  return {
    loading,
    tableData,
    availableColumns,
    searchForm,
    pagination,
    dateShortcuts,
    filterDialogVisible,
    propertyDrawerVisible,
    formattedProperties,
    activeFilters,
    allEventProperties,
    projectEvents,
    fetchMeta,
    loadData,
    handleSearch,
    handlePageChange,
    handleSizeChange,
    viewProperties,
    copyToClipboard
  }
}

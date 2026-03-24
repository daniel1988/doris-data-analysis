import { dimensionsApi } from '@/api/dimensions'
import { ref, unref, watch, type Ref } from 'vue'

// 简单的内存缓存，避免频繁聚焦导致重复请求
const valuesCache = new Map<string, string[]>()

export function useDimensionValues(props: {
  projectAlias?: string | Ref<string | undefined>
  tableName?: string | Ref<string | undefined>
  fieldName?: string | Ref<string | undefined>
  eventId?: string | Ref<string | undefined>
}) {
  const allValues = ref<string[]>([])
  const loading = ref(false)
  let loadedOnce = false

  const getCacheKey = () => {
    return `${unref(props.projectAlias)}_${unref(props.tableName)}_${unref(props.fieldName)}_${unref(props.eventId) || 'global'}`
  }

  const loadValues = async (force = false) => {
    // 如果已经加载过且不需要强制加载，或者正在加载中，则直接返回
    if ((loadedOnce && !force) || loading.value) {
      if (loadedOnce) {
        allValues.value = valuesCache.get(getCacheKey()) || []
      }
      return
    }

    const projectAlias = unref(props.projectAlias)
    const tableName = unref(props.tableName)
    const fieldName = unref(props.fieldName)
    const eventId = unref(props.eventId)

    if (!projectAlias || !tableName || !fieldName) return

    const cacheKey = getCacheKey()
    if (!force && valuesCache.has(cacheKey)) {
      allValues.value = valuesCache.get(cacheKey) || []
      loadedOnce = true
      return
    }

    // 如果是强制刷新，且缓存中已经有数据，先赋值缓存数据以供展示
    if (valuesCache.has(cacheKey)) {
      allValues.value = valuesCache.get(cacheKey) || []
    }

    loading.value = true
    try {
      const resp = await dimensionsApi.listValues({
        project_alias: projectAlias,
        table: tableName,
        field: fieldName,
        e_event_id: eventId || undefined,
      })

      const raw = resp as any
      const list = Array.isArray(raw) ? raw : (Array.isArray(raw?.data) ? raw.data : [])
      const values: string[] = []

      for (const row of list) {
        // 兼容后端返回 'value' 或原始字段名
        const v = row?.value ?? row?.[fieldName]
        if (v !== null && v !== undefined) {
          values.push(String(v))
        }
      }

      // 去重并排序
      const uniqueValues = Array.from(new Set(values)).sort()
      allValues.value = uniqueValues
      valuesCache.set(cacheKey, uniqueValues)
      loadedOnce = true
    } catch (error) {
      console.error('Failed to load dimension values:', error)
    } finally {
      loading.value = false
    }
  }

  // 监听依赖参数变化，重置状态
  watch(
    () => [unref(props.projectAlias), unref(props.tableName), unref(props.fieldName), unref(props.eventId)],
    () => {
      loadedOnce = false
      allValues.value = []
    }
  )

  return {
    allValues,
    loading,
    loadValues
  }
}

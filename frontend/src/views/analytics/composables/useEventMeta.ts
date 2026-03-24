import { ref, computed } from 'vue'
import { getEventOptions, getPropertyOptions, type Option } from '@/api/selector'
import { useAppStore } from '@/store/app'

export function useEventMeta() {
  const appStore = useAppStore()
  const eventOptions = ref<Option[]>([])
  const propertyOptions = ref<Option[]>([])
  const loading = ref(false)

  const fetchMeta = async () => {
    if (!appStore.activeProjectAlias) return
    loading.value = true
    try {
      const [events, properties] = await Promise.all([
        getEventOptions(appStore.activeProjectAlias),
        getPropertyOptions(appStore.activeProjectAlias)
      ])
      eventOptions.value = events
      propertyOptions.value = properties
    } catch (error) {
      console.error('Failed to fetch metadata:', error)
    } finally {
      loading.value = false
    }
  }

  return {
    eventOptions,
    propertyOptions,
    loading,
    fetchMeta
  }
}

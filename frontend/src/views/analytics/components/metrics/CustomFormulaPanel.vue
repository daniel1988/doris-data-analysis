<template>
  <div class="custom-formula-group">
    <CustomFormulaInput
      :model-value="m.custom_web_metric.eventMetrics"
      @update:model-value="(elements: any[]) => emit('update-elements', elements)"
      @request-new-event="() => emit('request-new-event')"
      @open-filter="emit('open-filter', m.value)"
    />
  </div>
</template>

<script setup lang="ts">
import CustomFormulaInput from '../CustomFormulaInput.vue';
import { computed } from 'vue';

const props = defineProps<{ m: any }>()
const emit = defineEmits<{ (e:'update-elements', elements:any[]):void; (e:'request-new-event'):void; (e:'open-filter', metric:any):void }>()

const m = computed(()=> props.m)

// 当自定义公式中任一事件 token 的 e_event_id 变化时，强制重建子组件，
// 确保子组件内部（如 PropertySelect/DimensionValueSelect）的请求参数使用最新事件ID
const eventsKey = computed(() => {
  try {
    const elements = (m.value?.custom_web_metric?.eventMetrics) || []
    return elements.map((el:any) => String(el?.value?.e_event_id || '')).join('|')
  } catch { return '' }
})
</script>

<style scoped>
.custom-formula-group { width: 100%; }
</style>

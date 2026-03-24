import { ref } from 'vue'

export function useTimeRangeSync() {
  const queryDateRange = ref<[string, string] | null>(null)
  const staticTimeRange = ref<[string, string] | null>(null)
  const comparisonTimeRange = ref<[string, string] | null>(null)
  const dynamicDateFilter = ref<{ operator: number; value: [string, string] } | null>(null)

  function handleStaticTimeChange(value: [string, string] | null) {
    staticTimeRange.value = value
    if (value && value.length === 2) {
      // 保持与后端格式一致，补全毫秒
      queryDateRange.value = [`${value[0]} 00:00:00`, `${value[1]} 23:59:59`]
      dynamicDateFilter.value = null
    } else {
      queryDateRange.value = null
    }
  }

  function handleComparisonChange(value: [string, string] | null) {
    comparisonTimeRange.value = value
  }

  /**
   * 动态时间选择
   * operator=19，value为长度为2的字符串数组
   * value[0]: '0'-过去N天, '1'-最近N天, '2'-自某日至今天
   * value[1]: N或日期字符串
   */
  function handleDynamicTimeChange(meta: { type: 'past' | 'recent' | 'since'; amountOrDate: string }) {
    if (!meta) {
      dynamicDateFilter.value = null
      return
    }

    if (meta.type === 'past') {
      dynamicDateFilter.value = { operator: 19, value: ['0', String(meta.amountOrDate)] as [string, string] }
    } else if (meta.type === 'recent') {
      dynamicDateFilter.value = { operator: 19, value: ['1', String(meta.amountOrDate)] as [string, string] }
    } else if (meta.type === 'since') {
      dynamicDateFilter.value = { operator: 19, value: ['2', String(meta.amountOrDate)] as [string, string] }
    }

    // 动态时间与静态时间互斥
    queryDateRange.value = null
    staticTimeRange.value = null
  }

  /**
   * 从表单数据同步时间状态
   */
  function syncFromForm(filterGroups: any) {
    if (!filterGroups) return

    // 同步静态日期
    if (filterGroups.query_dates && filterGroups.query_dates.length === 2) {
      handleStaticTimeChange([filterGroups.query_dates[0], filterGroups.query_dates[1]])
    } else {
      staticTimeRange.value = null
      queryDateRange.value = null
    }

    // 同步对比日期
    if (filterGroups.comparison_query_dates && filterGroups.comparison_query_dates.length === 2) {
      comparisonTimeRange.value = [filterGroups.comparison_query_dates[0], filterGroups.comparison_query_dates[1]]
    } else {
      comparisonTimeRange.value = null
    }

    // 同步动态日期过滤器 (operator 19)
    const dynamicFilter = (filterGroups.global_filters?.filters || []).find((f: any) => f.operator === 19)
    if (dynamicFilter) {
      dynamicDateFilter.value = { operator: 19, value: dynamicFilter.value }
    } else {
      dynamicDateFilter.value = null
    }
  }

  return {
    queryDateRange,
    staticTimeRange,
    comparisonTimeRange,
    dynamicDateFilter,
    handleStaticTimeChange,
    handleComparisonChange,
    handleDynamicTimeChange,
    syncFromForm
  }
}

<template>
  <div class="chart-view">
    <div ref="chartRef" class="echarts-container"></div>
  </div>
</template>

<script setup lang="ts">
import { useResizeObserver } from '@vueuse/core';
import * as echarts from 'echarts';
import { nextTick, onMounted, onUnmounted, ref, shallowRef, watch } from 'vue';

const props = defineProps<{
  rows: any[]
  columns: string[]
  vizMode: string
}>()

const chartRef = ref<HTMLElement>()
const chartInstance = shallowRef<echarts.ECharts | null>(null)

const initChart = () => {
  if (!chartRef.value) return

  if (!chartInstance.value) {
    chartInstance.value = echarts.init(chartRef.value)
  }

  const option = buildChartOption()
  if (Object.keys(option).length > 0) {
    // Use notMerge: true for smooth transitions
    chartInstance.value.setOption(option, { notMerge: true })
  }
}

const buildChartOption = () => {
  if (!props.rows || props.rows.length === 0 || !props.columns || props.columns.length === 0) return {}

  // 识别维度列和指标列
  // 1. 优先寻找包含“日期”或“时间”的列作为主要维度
  // 2. 否则，如果列数 > 1，通常第一列是维度（如“城市”、“渠道”），后续是指标
  // 3. 如果只有 1 列，或者所有列都是数值型，则可能没有维度，只有指标

  let dimensionCol = props.columns.find(c => c.includes('日期') || c.includes('时间'))

  // 如果没找到日期列，但有多个列，取第一列作为维度
  if (!dimensionCol && props.columns.length > 1) {
    dimensionCol = props.columns[0]
  }

  // 如果还是没找到（只有一列），或者该列看起来像指标（数值型），则视作无维度
  const isNumeric = (col: string) => {
    const val = props.rows[0]?.[col]
    return typeof val === 'number' || (!isNaN(Number(val)) && val !== '')
  }

  if (dimensionCol && isNumeric(dimensionCol) && props.columns.length === 1) {
    dimensionCol = undefined
  }

  const metricCols = props.columns.filter(c => c !== dimensionCol)

  // 处理饼图
  if (props.vizMode === 'pie') {
    const nameCol = dimensionCol || props.columns[0]
    const valueCol = metricCols[0] || props.columns[0]

    const pieData = props.rows.map(row => ({
      name: String(row[nameCol] || ''),
      value: Number(row[valueCol]) || 0
    })).filter(item => item.value > 0)

    return {
      tooltip: { trigger: 'item', formatter: '{a} <br/>{b}: {c} ({d}%)' },
      legend: { orient: 'vertical', left: 'left', type: 'scroll' },
      series: [{
        name: valueCol,
        type: 'pie',
        radius: '50%',
        data: pieData,
        emphasis: {
          itemStyle: { shadowBlur: 10, shadowOffsetX: 0, shadowColor: 'rgba(0, 0, 0, 0.5)' }
        }
      }]
    }
  }

  // 处理无维度的情况（单点展示，使用柱状图）
  if (!dimensionCol) {
    return {
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: ['合计'] },
      yAxis: { type: 'value' },
      series: metricCols.map(name => ({
        name,
        type: 'bar',
        data: [Number(props.rows[0]?.[name]) || 0],
        barMaxWidth: 50
      }))
    }
  }

  // 处理有时序或维度的折线/柱状图
  const sortedRows = [...props.rows]
  if (dimensionCol.includes('日期') || dimensionCol.includes('时间')) {
    sortedRows.sort((a, b) => {
      const valA = String(a[dimensionCol!] || '')
      const valB = String(b[dimensionCol!] || '')
      return valA.localeCompare(valB)
    })
  }

  const xAxisData = sortedRows.map(row => row[dimensionCol!])

  const series = metricCols.map(name => {
    const isComparison = name.includes('(对比)')
    const type = (props.vizMode === 'bar') ? 'bar' : 'line'
    return {
      name,
      type: type,
      data: sortedRows.map(row => row[name]),
      smooth: type === 'line',
      symbol: 'circle',
      symbolSize: 8,
      lineStyle: isComparison ? { type: 'dashed', width: 2 } : { width: 3 },
      emphasis: { focus: 'series' }
    }
  })

  return {
    tooltip: {
      trigger: 'axis',
      confine: true,
      backgroundColor: 'rgba(255, 255, 255, 0.95)',
      borderColor: '#eee',
      borderWidth: 1,
      textStyle: { color: '#333' },
      axisPointer: { type: 'shadow' }
    },
    legend: { data: metricCols, bottom: 10, type: 'scroll', icon: 'roundRect' },
    grid: { left: '3%', right: '4%', bottom: '15%', top: '5%', containLabel: true },
    xAxis: {
      type: 'category',
      data: xAxisData,
      axisLabel: { rotate: 30, color: '#666' },
      axisLine: { lineStyle: { color: '#eee' } }
    },
    yAxis: {
      type: 'value',
      axisLabel: { color: '#666' },
      splitLine: { lineStyle: { type: 'dashed', color: '#f5f5f5' } }
    },
    series
  }
}

watch(() => [props.rows, props.vizMode], () => {
  initChart()
}, { deep: true })

onMounted(() => {
  nextTick(() => {
    initChart()
  })
})

onUnmounted(() => {
  chartInstance.value?.dispose()
})

useResizeObserver(chartRef, () => {
  chartInstance.value?.resize()
})

// Expose save as image action
defineExpose({
  saveAsImage: () => {
    if (chartInstance.value) {
      const url = chartInstance.value.getDataURL({
        type: 'png',
        pixelRatio: 2,
        backgroundColor: '#fff'
      })
      const a = document.createElement('a')
      a.href = url
      a.download = `chart-${Date.now()}.png`
      a.click()
    }
  }
})
</script>

<style scoped>
.chart-view {
  width: 100%;
  height: 100%;
  min-height: 320px;
}

.echarts-container {
  height: 100%;
  width: 100%;
}
</style>

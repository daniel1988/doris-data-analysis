<template>
  <div class="chart-container" ref="chartRef"></div>
</template>

<script setup lang="ts">
import * as echarts from 'echarts'
import { onMounted, ref, watch, onUnmounted, nextTick } from 'vue'

const props = defineProps<{
  vizType: string
  xAxis: string
  yAxis: string
  data: any[]
}>()

const chartRef = ref<HTMLElement | null>(null)
let chartInstance: echarts.ECharts | null = null

const initChart = () => {
  if (!chartRef.value) return
  
  if (chartInstance) {
    chartInstance.dispose()
  }

  chartInstance = echarts.init(chartRef.value)
  renderChart()
}

const renderChart = () => {
  if (!chartInstance) return

  let option = {}

  if (props.vizType === 'line' || props.vizType === 'bar') {
    const xAxisData = props.data.map((item: any) => item[props.xAxis])
    // y_axis can be comma separated or single
    const yKeys = props.yAxis ? props.yAxis.split(',').map((s: string) => s.trim()) : []

    // If y_axis not provided, try to guess numeric columns excluding x_axis
    if (yKeys.length === 0 && props.data.length > 0) {
      const cols = Object.keys(props.data[0])
      cols.forEach(c => {
        if (c !== props.xAxis && typeof props.data[0][c] === 'number') {
          yKeys.push(c)
        }
      })
    }

    const series = yKeys.map((key: string) => ({
      name: key,
      type: props.vizType,
      data: props.data.map((item: any) => item[key])
    }))

    option = {
      tooltip: { trigger: 'axis' },
      legend: { data: yKeys },
      xAxis: { type: 'category', data: xAxisData },
      yAxis: { type: 'value' },
      series: series
    }
  } else if (props.vizType === 'pie') {
    // Expect x_axis as name, y_axis as value
    const nameKey = props.xAxis
    const valueKey = props.yAxis

    const data = props.data.map((item: any) => ({
      name: item[nameKey],
      value: item[valueKey]
    }))

    option = {
      tooltip: { trigger: 'item' },
      legend: { orient: 'vertical', left: 'left' },
      series: [
        {
          name: valueKey || 'Value',
          type: 'pie',
          radius: '50%',
          data: data,
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)'
            }
          }
        }
      ]
    }
  }

  chartInstance.setOption(option)
}

onMounted(() => {
  nextTick(() => {
    initChart()
  })
})

watch(() => props.data, () => {
  if (chartInstance) {
    renderChart()
  }
}, { deep: true })

onUnmounted(() => {
  if (chartInstance) {
    chartInstance.dispose()
    chartInstance = null
  }
})
</script>

<style scoped>
.chart-container {
  width: 100%;
  height: 300px;
}
</style>

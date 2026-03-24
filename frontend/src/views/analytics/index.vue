<template>
  <div class="analytics-page">
    <el-card class="analytics-card">
      <template #header>
        <div class="card-header">
          <span>事件分析趋势</span>
          <div class="header-ops">
            <el-radio-group v-model="timeRange" size="small">
              <el-radio-button value="7d">最近7天</el-radio-button>
              <el-radio-button value="30d">最近30天</el-radio-button>
            </el-radio-group>
            <el-button type="primary" size="small" style="margin-left: 10px">导出数据</el-button>
          </div>
        </div>
      </template>
      
      <div ref="chartRef" class="chart-container"></div>
    </el-card>
    
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card>
          <template #header>事件分布</template>
          <div ref="pieChartRef" class="sub-chart-container"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>明细数据</template>
          <el-table :data="tableData" style="width: 100%" height="250">
            <el-table-column prop="date" label="日期" width="120" />
            <el-table-column prop="event" label="事件名称" />
            <el-table-column prop="count" label="次数" width="100" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import * as echarts from 'echarts'
import { onMounted, onUnmounted, ref, watch } from 'vue'

const chartRef = ref<HTMLElement>()
const pieChartRef = ref<HTMLElement>()
const timeRange = ref('7d')
let lineChart: echarts.ECharts | null = null
let pieChart: echarts.ECharts | null = null

const tableData = [
  { date: '2024-03-01', event: 'page_view', count: 1200 },
  { date: '2024-03-01', event: 'button_click', count: 450 },
  { date: '2024-03-02', event: 'page_view', count: 1500 },
  { date: '2024-03-02', event: 'button_click', count: 600 },
]

const initCharts = () => {
  if (chartRef.value) {
    lineChart = echarts.init(chartRef.value)
    lineChart.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: {
        type: 'category',
        data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
      },
      yAxis: { type: 'value' },
      series: [{
        data: [820, 932, 901, 934, 1290, 1330, 1320],
        type: 'line',
        smooth: true,
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(64, 158, 255, 0.5)' },
            { offset: 1, color: 'rgba(64, 158, 255, 0)' }
          ])
        }
      }]
    })
  }
  
  if (pieChartRef.value) {
    pieChart = echarts.init(pieChartRef.value)
    pieChart.setOption({
      tooltip: { trigger: 'item' },
      series: [
        {
          name: '事件分布',
          type: 'pie',
          radius: ['40%', '70%'],
          data: [
            { value: 1048, name: '页面访问' },
            { value: 735, name: '按钮点击' },
            { value: 580, name: '表单提交' },
            { value: 484, name: '搜索' },
            { value: 300, name: '其他' }
          ]
        }
      ]
    })
  }
}

const handleResize = () => {
  lineChart?.resize()
  pieChart?.resize()
}

onMounted(() => {
  initCharts()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  lineChart?.dispose()
  pieChart?.dispose()
})

watch(timeRange, () => {
  // 模拟数据更新
  lineChart?.setOption({
    series: [{
      data: Array.from({ length: 7 }, () => Math.floor(Math.random() * 1000 + 500))
    }]
  })
})
</script>

<style scoped>
.analytics-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-container {
  height: 400px;
}

.sub-chart-container {
  height: 250px;
}
</style>

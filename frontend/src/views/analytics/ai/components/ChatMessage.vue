<template>
  <div class="message-row" :class="msg.role">
    <div class="avatar">
      <el-avatar 
        :icon="msg.role === 'user' ? UserFilled : Cpu" 
        :size="36" 
        :class="msg.role === 'user' ? 'user-avatar' : 'ai-avatar'" 
      />
    </div>
    
    <div class="message-content">
      <!-- User Query -->
      <div v-if="msg.role === 'user'" class="bubble user-bubble">
        {{ msg.content }}
      </div>

      <!-- AI Response -->
      <div v-else class="bubble ai-bubble">
        <template v-if="msg.loading">
          <div class="loading-state">
            <el-icon class="is-loading"><Loading /></el-icon> Thinking...
          </div>
        </template>
        <template v-else-if="msg.mode === 'chat'">
          <div class="markdown-content" v-html="renderMarkdown(msg.content || '')"></div>
        </template>
        <template v-else>
          <div v-if="msg.narrative" class="narrative">{{ msg.narrative }}</div>
          
          <div v-if="msg.sql" class="sql-box">
            <div class="sql-header">Generated SQL</div>
            <pre>{{ msg.sql }}</pre>
          </div>

          <div v-if="msg.data && msg.data.length > 0" class="viz-box">
            <el-tabs 
              v-if="canVisualize" 
              v-model="activeTab" 
              type="card" 
              class="viz-tabs"
            >
              <el-tab-pane label="Chart" name="chart">
                <ChatChart 
                  v-if="activeTab === 'chart'"
                  :viz-type="msg.viz_type"
                  :x-axis="msg.x_axis"
                  :y-axis="msg.y_axis"
                  :data="msg.data"
                />
              </el-tab-pane>
              <el-tab-pane label="Data Table" name="table">
                <el-table 
                  :data="msg.data" 
                  border 
                  size="small" 
                  style="width: 100%;"
                  max-height="400"
                >
                  <el-table-column 
                    v-for="col in columns" 
                    :key="col" 
                    :prop="col" 
                    :label="col" 
                    min-width="120"
                    show-overflow-tooltip
                  />
                </el-table>
              </el-tab-pane>
            </el-tabs>
            
            <el-table 
              v-else 
              :data="msg.data" 
              border 
              size="small" 
              style="width: 100%; margin-top: 10px;"
              max-height="400"
            >
              <el-table-column 
                v-for="col in columns" 
                :key="col" 
                :prop="col" 
                :label="col" 
                min-width="120"
                show-overflow-tooltip
              />
            </el-table>
          </div>
          
          <div v-else-if="msg.data && msg.data.length === 0" class="no-data">
            No data found for this query.
          </div>
          
          <div v-if="msg.error" class="error-msg">
            Error: {{ msg.error }}
          </div>
          
          <!-- Save Action Area -->
          <div v-if="!msg.error && msg.sql" class="action-bar">
            <el-button 
              type="primary" 
              link 
              size="small" 
              :icon="msg.saved ? StarFilled : Star" 
              :disabled="msg.saved"
              @click="$emit('save', msg)"
            >
              {{ msg.saved ? '已保存为会话记录' : '保存为有效分析' }}
            </el-button>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Cpu, Loading, Star, StarFilled, UserFilled } from '@element-plus/icons-vue'
import MarkdownIt from 'markdown-it'
import { computed, ref } from 'vue'
import ChatChart from './ChatChart.vue'

const md = new MarkdownIt({
  html: true,
  breaks: true,
  linkify: true
})

interface Message {
  role: 'user' | 'ai'
  mode?: 'sql' | 'chat'
  content?: string
  narrative?: string
  sql?: string
  viz_type?: string
  x_axis?: string
  y_axis?: string
  data?: any[]
  error?: string
  loading?: boolean
  saved?: boolean
}

const props = defineProps<{
  msg: Message
}>()

const emit = defineEmits<{
  (e: 'save', msg: Message): void
}>()

const activeTab = ref('chart')

const renderMarkdown = (text: string) => {
  return md.render(text)
}

const canVisualize = computed(() => {
  const { msg } = props
  return msg.viz_type && 
         msg.viz_type !== 'table' && 
         msg.data && 
         msg.data.length > 0 &&
         (msg.x_axis || msg.viz_type === 'pie')
})

const columns = computed(() => {
  if (!props.msg.data || props.msg.data.length === 0) return []
  return Object.keys(props.msg.data[0])
})
</script>

<style scoped>
.message-row {
  display: flex;
  margin-bottom: 20px;
}

.message-row.user {
  flex-direction: row-reverse;
}

.avatar {
  margin: 0 10px;
}

.user-avatar {
  background-color: #409eff;
}

.ai-avatar {
  background-color: #67c23a;
}

.message-content {
  max-width: 80%;
  min-width: 300px;
}

.bubble {
  padding: 12px 16px;
  border-radius: 8px;
  position: relative;
}

.user-bubble {
  background-color: #409eff;
  color: white;
  border-top-right-radius: 0;
}

.ai-bubble {
  background-color: white;
  color: #303133;
  border-top-left-radius: 0;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.narrative {
  margin-bottom: 10px;
  white-space: pre-wrap;
}

.sql-box {
  background-color: #282c34;
  color: #abb2bf;
  padding: 10px;
  border-radius: 4px;
  font-family: monospace;
  margin-bottom: 10px;
  font-size: 12px;
  overflow-x: auto;
}

.sql-header {
  color: #61afef;
  font-weight: bold;
  margin-bottom: 5px;
}

.viz-box {
  background-color: #fff;
  padding: 10px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
}

.viz-tabs :deep(.el-tabs__header) {
  margin-bottom: 10px;
}

.error-msg {
  color: #f56c6c;
  margin-top: 5px;
  font-size: 12px;
}

.loading-state {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #909399;
}

.action-bar {
  margin-top: 15px;
  padding-top: 10px;
  border-top: 1px solid #ebeef5;
  text-align: right;
}

.markdown-content :deep(p) {
  margin-bottom: 10px;
  line-height: 1.6;
}

.markdown-content :deep(pre) {
  background-color: #282c34;
  color: #abb2bf;
  padding: 10px;
  border-radius: 4px;
  overflow-x: auto;
}

.markdown-content :deep(code) {
  font-family: monospace;
  background-color: #f0f2f5;
  padding: 2px 4px;
  border-radius: 2px;
  color: #d63200;
}
</style>

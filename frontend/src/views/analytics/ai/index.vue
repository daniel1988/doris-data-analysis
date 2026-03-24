<template>
  <div class="ai-chat-container">
    <!-- Model Selector -->
    <div class="model-selector-bar">
      <div class="left">
        <span class="label">AI 模型：</span>
        <el-select 
          v-model="selectedModelId" 
          placeholder="选择模型" 
          size="small"
          :loading="modelsLoading"
          style="width: 200px"
        >
          <el-option 
            v-for="m in availableModels" 
            :key="m.id" 
            :label="m.display_name" 
            :value="m.id"
          >
            <span>{{ m.display_name }}</span>
            <el-tag v-if="m.is_default" size="small" type="success" style="margin-left: 8px">默认</el-tag>
          </el-option>
        </el-select>
      </div>
      <div class="right">
        <el-button size="small" :icon="Clock" @click="showHistoryDrawer = true">历史分析记录</el-button>
      </div>
    </div>

    <div class="messages-area" ref="messagesArea">
      <!-- Welcome Screen -->
      <div v-if="messages.length === 0" class="welcome-screen">
        <h2>AI Data Analyst</h2>
        <p>Ask questions about your data in natural language.</p>
        <div class="examples">
          <el-tag 
            v-for="ex in examples" 
            :key="ex" 
            class="example-tag" 
            @click="setInput(ex)"
          >
            {{ ex }}
          </el-tag>
        </div>
      </div>

      <!-- Message List -->
      <ChatMessage 
        v-for="(msg, index) in messages" 
        :key="index" 
        :msg="msg" 
        @save="handleSaveSession(msg)"
      />

      <!-- Loading State -->
      <ChatMessage 
        v-if="loading" 
        :msg="{ role: 'ai', loading: true }" 
      />
    </div>

    <!-- Input Area -->
    <ChatInput 
      ref="chatInputRef"
      :loading="loading"
      :hasContextData="hasContextData"
      @send="sendMessage" 
    />

    <!-- History Drawer -->
    <el-drawer
      v-model="showHistoryDrawer"
      title="历史分析记录"
      direction="rtl"
      size="400px"
      @open="fetchHistorySessions"
    >
      <div v-loading="historyLoading" class="history-list">
        <el-empty v-if="historySessions.length === 0" description="暂无历史保存记录" />
        <div 
          v-for="session in historySessions" 
          :key="session.id" 
          class="history-item"
        >
          <div class="header">
            <span class="time">{{ session.create_time }}</span>
            <el-button link type="danger" size="small" @click="handleDeleteSession(session.id!)">删除</el-button>
          </div>
          <div class="query">{{ session.user_query }}</div>
          <div class="actions">
            <el-tag size="small" type="info">{{ session.viz_type }}</el-tag>
            <el-button type="primary" size="small" @click="handleLoadSession(session)">点击加载</el-button>
          </div>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { getEnabledAIModels } from '@/api/aiModel'
import type { AISession } from '@/api/analytics'
import { aiChat, deleteAISession, executeAISession, getAISessions, saveAISession } from '@/api/analytics'
import { useAppStore } from '@/store/app'
import { Clock } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, nextTick, onMounted, ref } from 'vue'
import ChatInput from './components/ChatInput.vue'
import ChatMessage from './components/ChatMessage.vue'

const appStore = useAppStore()
const loading = ref(false)
const messages = ref<any[]>([])
const messagesArea = ref<HTMLElement | null>(null)
const chatInputRef = ref<InstanceType<typeof ChatInput> | null>(null)

// Check if there's any previous successful result data
const hasContextData = computed(() => {
  for (let i = messages.value.length - 1; i >= 0; i--) {
    const msg = messages.value[i]
    if (msg.role === 'ai' && msg.data && msg.data.length > 0) {
      return true
    }
  }
  return false
})

// Model selector
const availableModels = ref<any[]>([])
const selectedModelId = ref<number>(0)
const modelsLoading = ref(false)

// History drawer
const showHistoryDrawer = ref(false)
const historyLoading = ref(false)
const historySessions = ref<AISession[]>([])

const examples = [
  "展示上周每个渠道的总用户数",
  "展示前5个付费用户",
  "展示每日活跃用户趋势"
]

const fetchModels = async () => {
  modelsLoading.value = true
  try {
    const data = await getEnabledAIModels() as any[]
    availableModels.value = data || []
    // 自动选择默认模型
    const defaultModel = data?.find((m: any) => m.is_default)
    if (defaultModel) {
      selectedModelId.value = defaultModel.id
    } else if (data?.length > 0) {
      selectedModelId.value = data[0].id
    }
  } catch {
    // 模型列表获取失败时不阻塞页面
  } finally {
    modelsLoading.value = false
  }
}

const setInput = (text: string) => {
  sendMessage({ query: text, mode: 'sql', useContext: false })
}

const sendMessage = async (payload: { query: string; mode: 'sql' | 'chat'; useContext: boolean }) => {
  const { query, mode, useContext } = payload
  if (!query.trim()) return
  if (loading.value) return

  messages.value.push({ role: 'user', content: query })
  loading.value = true
  scrollToBottom()

  let contextData = undefined
  if (mode === 'chat' && useContext) {
    // Find the latest result data
    for (let i = messages.value.length - 1; i >= 0; i--) {
      const msg = messages.value[i]
      if (msg.role === 'ai' && msg.data && msg.data.length > 0) {
        contextData = msg.data
        break
      }
    }
  }

  try {
    const res = await aiChat({
      project_alias: appStore.activeProjectAlias,
      query: query,
      model_id: selectedModelId.value || undefined,
      mode: mode,
      context_data: contextData
    })
    
    const result = res as any
    
    if (mode === 'chat' && result.type === 'text') {
      messages.value.push({
        role: 'ai',
        mode: 'chat',
        content: result.content
      })
    } else {
      const aiMsg = {
        role: 'ai',
        mode: 'sql',
        userQuery: query, // Save the query for session saving
        narrative: result.narrative || "Here is the result:",
        sql: result.sql,
        viz_type: result.viz_type,
        x_axis: result.x_axis,
        y_axis: result.y_axis,
        data: result.data,
        saved: false
      }
      messages.value.push(aiMsg)
    }
    
    nextTick(() => {
      scrollToBottom()
    })

  } catch (err: any) {
    console.error(err)
    messages.value.push({
      role: 'ai',
      narrative: "Sorry, I encountered an error processing your request.",
      error: err.message || "Unknown error"
    })
    scrollToBottom()
  } finally {
    loading.value = false
  }
}

const scrollToBottom = () => {
  nextTick(() => {
    if (messagesArea.value) {
      messagesArea.value.scrollTop = messagesArea.value.scrollHeight
    }
  })
}

// Session Management
const handleSaveSession = async (msg: any) => {
  try {
    const sessionData: AISession = {
      project_alias: appStore.activeProjectAlias,
      user_query: msg.userQuery || '未记录的问题',
      llm_sql: msg.sql || '',
      viz_type: msg.viz_type || 'table',
      x_axis: msg.x_axis || '',
      y_axis: msg.y_axis || '',
      narrative: msg.narrative || ''
    }
    await saveAISession(sessionData)
    ElMessage.success('已保存为有效分析')
    msg.saved = true
  } catch (error) {
    console.error('Failed to save session', error)
  }
}

const fetchHistorySessions = async () => {
  historyLoading.value = true
  try {
    const res = await getAISessions(appStore.activeProjectAlias)
    historySessions.value = res as unknown as AISession[]
  } catch (error) {
    console.error('Failed to load history', error)
  } finally {
    historyLoading.value = false
  }
}

const handleDeleteSession = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除该历史记录吗？', '提示', { type: 'warning' })
    await deleteAISession(id)
    ElMessage.success('删除成功')
    fetchHistorySessions()
  } catch (e) {
    // cancelled or error
  }
}

const handleLoadSession = async (session: AISession) => {
  showHistoryDrawer.value = false
  
  messages.value.push({ 
    role: 'user', 
    content: `[加载历史分析]: ${session.user_query}` 
  })
  
  loading.value = true
  scrollToBottom()

  try {
    const res = await executeAISession(session.id!, appStore.activeProjectAlias)
    const result = res as any
    
    messages.value.push({
      role: 'ai',
      userQuery: session.user_query,
      narrative: result.narrative || session.narrative || "Here is the historical result:",
      sql: result.sql,
      viz_type: result.viz_type,
      x_axis: result.x_axis,
      y_axis: result.y_axis,
      data: result.data,
      saved: true
    })
    
    nextTick(() => scrollToBottom())
  } catch (err: any) {
    messages.value.push({
      role: 'ai',
      narrative: "加载历史数据失败。",
      error: err.message || "Unknown error"
    })
    scrollToBottom()
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchModels()
})
</script>

<style scoped>
.ai-chat-container {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 120px);
  padding: 20px;
  background-color: #f5f7fa;
}

.model-selector-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding: 8px 12px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
}

.model-selector-bar .left {
  display: flex;
  align-items: center;
}

.model-selector-bar .label {
  font-size: 13px;
  color: #606266;
  margin-right: 8px;
}

.messages-area {
  flex: 1;
  overflow-y: auto;
  margin-bottom: 20px;
  padding-right: 10px;
}

.welcome-screen {
  text-align: center;
  margin-top: 100px;
  color: #606266;
}

.examples {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  gap: 10px;
}

.example-tag {
  cursor: pointer;
}

.history-list {
  padding: 10px;
}

.history-item {
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 12px;
  margin-bottom: 12px;
  background: #fafafa;
  transition: all 0.3s;
}

.history-item:hover {
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
}

.history-item .header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.history-item .time {
  font-size: 12px;
  color: #909399;
}

.history-item .query {
  font-size: 14px;
  color: #303133;
  margin-bottom: 12px;
  font-weight: 500;
}

.history-item .actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>

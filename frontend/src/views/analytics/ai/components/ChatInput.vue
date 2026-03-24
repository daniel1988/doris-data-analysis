<template>
  <div class="input-area">
    <div class="mode-switch">
      <el-radio-group v-model="mode" size="small">
        <el-radio-button value="sql">智能取数</el-radio-button>
        <el-radio-button value="chat">AI 对话</el-radio-button>
      </el-radio-group>
      <el-checkbox 
        v-if="mode === 'chat' && hasContextData" 
        v-model="useContext" 
        style="margin-left: 16px"
        size="small"
      >
        引用当前结果数据
      </el-checkbox>
    </div>
    <el-input
      v-model="inputValue"
      :placeholder="mode === 'sql' ? 'Ask something about your data... (Enter to send)' : 'Chat with AI... (Enter to send)'"
      @keyup.enter="handleSend"
      :disabled="loading"
    >
      <template #append>
        <el-button @click="handleSend" :loading="loading">
          <el-icon><Position /></el-icon>
        </el-button>
      </template>
    </el-input>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Position } from '@element-plus/icons-vue'

const props = defineProps<{
  loading: boolean;
  hasContextData?: boolean;
}>()

const emit = defineEmits<{
  (e: 'send', payload: { query: string; mode: 'sql' | 'chat'; useContext: boolean }): void
}>()

const inputValue = ref('')
const mode = ref<'sql' | 'chat'>('sql')
const useContext = ref(false)

const handleSend = () => {
  if (!inputValue.value.trim() || props.loading) return
  emit('send', { 
    query: inputValue.value, 
    mode: mode.value, 
    useContext: useContext.value 
  })
  inputValue.value = ''
}
</script>

<style scoped>
.input-area {
  padding: 10px 0;
}
.mode-switch {
  margin-bottom: 8px;
  display: flex;
  align-items: center;
}
</style>

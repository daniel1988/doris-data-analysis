<template>
  <el-drawer
    v-model="visible"
    title="生成的 SQL 查询语句"
    direction="rtl"
    size="50%"
    destroy-on-close
  >
    <div class="sql-drawer-content">
      <div class="sql-actions">
        <el-button size="small" type="primary" icon="DocumentCopy" @click="copySql">复制 SQL</el-button>
        <el-button size="small" icon="MagicStick" @click="formatSql">格式化</el-button>
      </div>

      <div class="sql-container">
        <pre class="sql-code">{{ displaySql }}</pre>
      </div>

      <div class="sql-tip">
        <el-icon><InfoFilled /></el-icon>
        此 SQL 由系统根据分析配置自动生成，仅供数据审计参考。
      </div>
    </div>
  </el-drawer>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'

const props = defineProps<{
  modelValue: boolean
  sql: string
}>()

const emit = defineEmits(['update:modelValue'])

const visible = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v)
})

const formattedSql = ref('')
const displaySql = computed(() => formattedSql.value || props.sql || '-- 暂无 SQL 生成')

const formatSql = () => {
  if (!props.sql) return
  let formatted = props.sql
    .replace(/\bSELECT\b/gi, 'SELECT')
    .replace(/\bFROM\b/gi, '\nFROM')
    .replace(/\bWHERE\b/gi, '\nWHERE')
    .replace(/\bGROUP BY\b/gi, '\nGROUP BY')
    .replace(/\bORDER BY\b/gi, '\nORDER BY')
    .replace(/\bLIMIT\b/gi, '\nLIMIT')
    .replace(/\bAND\b/gi, '\n  AND')
    .replace(/\bOR\b/gi, '\n  OR')
    .replace(/,/g, ',\n  ')
  
  formattedSql.value = formatted.replace(/\n\s*\n/g, '\n').trim()
  ElMessage.success('SQL 已格式化')
}

const copySql = async () => {
  try {
    await navigator.clipboard.writeText(displaySql.value)
    ElMessage.success('已复制到剪贴板')
  } catch (err) {
    ElMessage.error('复制失败')
  }
}

watch(() => props.sql, () => {
  formattedSql.value = ''
})
</script>

<style scoped lang="scss">
.sql-drawer-content {
  display: flex;
  flex-direction: column;
  height: 100%;

  .sql-actions {
    margin-bottom: 16px;
    display: flex;
    gap: 10px;
  }

  .sql-container {
    flex: 1;
    background: #282c34;
    color: #abb2bf;
    padding: 15px;
    border-radius: 4px;
    overflow: auto;
    
    .sql-code {
      margin: 0;
      font-family: 'Fira Code', monospace;
      font-size: 13px;
      line-height: 1.6;
      white-space: pre-wrap;
      word-break: break-all;
    }
  }

  .sql-tip {
    margin-top: 16px;
    font-size: 12px;
    color: var(--el-text-color-secondary);
    display: flex;
    align-items: center;
    gap: 5px;
  }
}
</style>

<template>
  <div class="unified-property-selector">
    <el-popover
      ref="popoverRef"
      placement="bottom-start"
      :width="360"
      trigger="click"
      popper-class="prop-select-popover"
      @show="onPopoverShow"
    >
      <template #reference>
        <div class="selector-trigger" :class="{ 'is-active': popoverVisible }">
          <slot name="trigger">
            <el-input
              v-model="displayValue"
              readonly
              :placeholder="placeholder"
              size="small"
              class="trigger-input"
            >
              <template #suffix>
                <el-icon><ArrowDown /></el-icon>
              </template>
            </el-input>
          </slot>
        </div>
      </template>

      <div class="selector-content">
        <!-- 搜索框 -->
        <div class="search-bar">
          <el-input
            v-model="searchQuery"
            placeholder="搜索属性名称或ID..."
            size="small"
            clearable
            prefix-icon="Search"
          />
        </div>

        <!-- Tab 切换 -->
        <el-tabs v-model="activeTab" class="selector-tabs" v-loading="loadingSpecificProps">
          <el-tab-pane label="项目属性" name="event">
            <div class="option-list">
              <div
                v-for="item in filteredEventProps"
                :key="item.id"
                class="option-item"
                @click="handleSelect(item, 'event')"
              >
                <span class="option-icon"><el-icon><Memo /></el-icon></span>
                <span class="option-name">{{ item.name }}</span>
                <span class="option-id">{{ item.id }}</span>
              </div>
              <div v-if="filteredEventProps.length === 0" class="empty-placeholder">
                暂无项目属性
              </div>
            </div>
          </el-tab-pane>
          <el-tab-pane label="用户属性" name="user">
            <div class="option-list">
              <div
                v-for="item in filteredUserProps"
                :key="item.id"
                class="option-item"
                @click="handleSelect(item, 'user')"
              >
                <span class="option-icon"><el-icon><User /></el-icon></span>
                <span class="option-name">{{ item.name }}</span>
                <span class="option-id">{{ item.id }}</span>
              </div>
            </div>
          </el-tab-pane>
          <el-tab-pane label="用户标签" name="tag">
            <div class="option-list">
              <div
                v-for="item in filteredTags"
                :key="item.id"
                class="option-item"
                @click="handleSelect(item, 'tag')"
              >
                <span class="option-icon"><el-icon><PriceTag /></el-icon></span>
                <span class="option-name">{{ item.name }}</span>
                <span class="option-id">{{ item.id }}</span>
              </div>
            </div>
          </el-tab-pane>
          <el-tab-pane label="用户分群" name="group">
            <div class="option-list">
              <div
                v-for="item in filteredGroups"
                :key="item.id"
                class="option-item"
                @click="handleSelect(item, 'group')"
              >
                <span class="option-icon"><el-icon><Files /></el-icon></span>
                <span class="option-name">{{ item.name }}</span>
                <span class="option-id">{{ item.id }}</span>
              </div>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-popover>
  </div>
</template>

<script setup lang="ts">
import { Option, getPropertyOptions } from '@/api/selector';
import { ArrowDown, Files, Memo, PriceTag, User } from '@element-plus/icons-vue';
import { computed, inject, ref, watch } from 'vue';
import { ANALYSIS_CONTEXT_KEY } from '../../context';

const props = defineProps<{
  modelValue?: string
  placeholder?: string
  type?: 'event' | 'user' | 'tag' | 'group'
  eventId?: string // 新增：特定事件ID，用于过滤属性
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'change', payload: { id: string; name: string; type: string; table: string; property?: any }): void
}>()

const context = inject(ANALYSIS_CONTEXT_KEY)
const popoverVisible = ref(false)
const searchQuery = ref('')
const activeTab = ref(props.type || 'event')
const popoverRef = ref()

// 特定事件的属性列表
const specificEventProps = ref<Option[]>([])
const loadingSpecificProps = ref(false)

const onPopoverShow = () => {
  popoverVisible.value = true
  if (activeTab.value === 'event' && props.eventId) {
    fetchSpecificEventProps()
  }
}

// 获取特定事件的属性
const fetchSpecificEventProps = async () => {
  if (!props.eventId || !context?.state.projectAlias) return
  
  loadingSpecificProps.value = true
  try {
    const res = await getPropertyOptions(context.state.projectAlias, props.eventId)
    specificEventProps.value = res || []
  } catch (err) {
    console.error('Failed to fetch event properties:', err)
  } finally {
    loadingSpecificProps.value = false
  }
}

// 监听 eventId 变化，如果已经打开弹窗则刷新
watch(() => props.eventId, () => {
  if (popoverVisible.value && activeTab.value === 'event') {
    fetchSpecificEventProps()
  }
})

// 监听标签切换
watch(activeTab, (newTab) => {
  if (newTab === 'event' && props.eventId && specificEventProps.value.length === 0) {
    fetchSpecificEventProps()
  }
})

const eventProps = computed(() => {
  // 如果有特定事件ID，展示该事件关联的属性；否则展示所有项目属性
  if (props.eventId && specificEventProps.value.length > 0) {
    return specificEventProps.value
  }
  // 兜底：展示所有项目属性（从 metadata 中过滤出 event_data 表的属性）
  return context?.state.metadata.propertyOptions.filter(p => p.table === 'event_data') || []
})

const userProps = computed(() => context?.state.metadata.propertyOptions.filter(p => p.table === 'user_data') || [])
const tags = computed(() => context?.state.metadata.propertyOptions.filter(p => p.table === 'user_tag_data') || [])
const groups = computed(() => context?.state.metadata.propertyOptions.filter(p => p.table === 'user_group') || [])

const filterOptions = (options: Option[]) => {
  if (!searchQuery.value) return options
  const q = searchQuery.value.toLowerCase()
  return options.filter(o => o.name.toLowerCase().includes(q) || o.id.toLowerCase().includes(q))
}

const filteredEventProps = computed(() => filterOptions(eventProps.value))
const filteredUserProps = computed(() => filterOptions(userProps.value))
const filteredTags = computed(() => filterOptions(tags.value))
const filteredGroups = computed(() => filterOptions(groups.value))

const displayValue = computed(() => {
  if (!props.modelValue) return ''
  // 查找顺序：特定事件属性 -> 全局属性 -> 原始值
  const all = [...specificEventProps.value, ...context?.state.metadata.propertyOptions || []]
  return all.find(o => o.id === props.modelValue)?.name || props.modelValue
})

const handleSelect = (item: Option, type: string) => {
  emit('update:modelValue', item.id)
  emit('change', { 
    id: item.id, 
    name: item.name, 
    type, 
    table: item.table || (type === 'event' ? 'event_data' : (type === 'user' ? 'user_data' : 'user_tag_data')),
    property: item // 传递完整的属性对象，包含 data_type 等信息
  })
  popoverVisible.value = false
  popoverRef.value?.hide()
}
</script>

<style scoped lang="scss">
.unified-property-selector {
  display: inline-block;
  width: 100%;
}

.selector-trigger {
  cursor: pointer;
  &.is-active {
    .el-input__wrapper {
      box-shadow: 0 0 0 1px var(--el-color-primary) inset;
    }
  }
}

.selector-content {
  padding: 8px;
}

.search-bar {
  margin-bottom: 12px;
}

.selector-tabs {
  :deep(.el-tabs__header) {
    margin-bottom: 8px;
  }
}

.option-list {
  max-height: 250px;
  overflow-y: auto;
}

.option-item {
  display: flex;
  align-items: center;
  padding: 8px;
  cursor: pointer;
  border-radius: 4px;
  transition: background 0.2s;

  &:hover {
    background: var(--el-fill-color-light);
  }

  .option-icon {
    margin-right: 8px;
    color: var(--el-text-color-secondary);
  }

  .option-name {
    flex: 1;
    font-size: 13px;
    color: var(--el-text-color-primary);
  }

  .option-id {
    font-size: 11px;
    color: var(--el-text-color-placeholder);
    margin-left: 8px;
  }
}

.empty-placeholder {
  padding: 20px;
  text-align: center;
  color: var(--el-text-color-placeholder);
  font-size: 13px;
}
</style>

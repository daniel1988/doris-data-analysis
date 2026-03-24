<template>
  <el-drawer
    v-model="visible"
    :title="drawerTitle"
    direction="rtl"
    size="80%"
    :before-close="handleClose"
    class="user-drilldown-drawer"
  >
    <div class="drawer-content">
      <UserListContent
        v-if="mode === 'tag'"
        :initial-tag-filter="tagFilterPayload"
        :project-alias="projectAlias"
        @total-change="handleTotalChange"
      />
      <UserListContent
        v-else
        :initial-user-group-filter="groupFilterPayload"
        :project-alias="projectAlias"
        @total-change="handleTotalChange"
      />
    </div>
  </el-drawer>
</template>

<script setup lang="ts">
import UserListContent from '@/components/common/UserListContent.vue'
import { computed, ref, watch } from 'vue'

type Mode = 'tag' | 'group'

interface TagFilterPayload {
  id?: number
  tagCode: string
  tagShowName: string
  userCount?: number
}

interface GroupFilterPayload {
  groupName?: string
  groupCode: string
  userCount?: number
  operator?: number
}

interface Props {
  modelValue: boolean
  mode: Mode
  tagFilter?: TagFilterPayload | null
  groupFilter?: GroupFilterPayload | null
  projectAlias?: string
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: false,
  projectAlias: '',
  tagFilter: null,
  groupFilter: null,
})

const emit = defineEmits<{
  (e: 'update:modelValue', v: boolean): void
}>()

const visible = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v)
})

const projectAlias = computed(() => props.projectAlias)

const actualTotal = ref<number | undefined>(undefined)

const tagFilterPayload = computed(() => {
  if (props.mode !== 'tag' || !props.tagFilter) return null
  return {
    id: props.tagFilter.id ?? 0,
    tagCode: props.tagFilter.tagCode,
    tagShowName: props.tagFilter.tagShowName,
    userCount: props.tagFilter.userCount,
  }
})

const groupFilterPayload = computed(() => {
  if (props.mode !== 'group' || !props.groupFilter) return null
  return {
    groupName: props.groupFilter.groupName,
    groupCode: props.groupFilter.groupCode,
    operator: props.groupFilter.operator ?? 10,
  }
})

const drawerTitle = computed(() => {
  if (props.mode === 'tag') {
    if (!props.tagFilter) return '用户列表'
    const count = actualTotal.value ?? props.tagFilter.userCount
    const suffix = count !== undefined ? ` - ${count.toLocaleString('zh-CN')}人` : ''
    return `${props.tagFilter.tagShowName}${suffix}`
  }
  // group mode
  if (!props.groupFilter) return '用户列表'
  const count = actualTotal.value ?? props.groupFilter.userCount
  const suffix = count !== undefined ? ` - ${count.toLocaleString('zh-CN')}人` : ''
  return `${props.groupFilter.groupName || props.groupFilter.groupCode}${suffix}`
})

const handleClose = () => {
  visible.value = false
}

const handleTotalChange = (total: number) => {
  actualTotal.value = total
}

watch(visible, (nv) => {
  if (!nv) actualTotal.value = undefined
})

watch(() => props.tagFilter?.tagCode, () => { if (props.mode === 'tag') actualTotal.value = undefined })
watch(() => props.groupFilter?.groupCode, () => { if (props.mode === 'group') actualTotal.value = undefined })
</script>

<style scoped lang="scss">
.user-drilldown-drawer {
  :deep(.el-drawer__header) {
    padding: 16px 20px;
    border-bottom: 1px solid #e5e7eb;
    margin-bottom: 0;
  }
  :deep(.el-drawer__body) {
    padding: 0;
    height: calc(100% - 64px);
    overflow: hidden;
  }
}
.drawer-content {
  height: 100%;
  display: flex;
  flex-direction: column;
}
</style>

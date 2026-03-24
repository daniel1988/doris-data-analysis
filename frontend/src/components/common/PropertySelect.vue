<template>
	<div class="prop-select-wrap">
		<el-popover 
			ref="popoverRef"
			placement="bottom-start" 
			:width="360" 
			trigger="click"
			popper-class="prop-select-popover"
		>
					<div class="prop-select-content">
			<!-- 全局搜索框 -->
			<div class="global-search">
				<el-input 
					v-model="searchKeyword" 
					placeholder="搜索属性名称、ID或数据类型..."
					size="small" 
					clearable
					class="search-input"
				>
					<template #prefix>
						<el-icon><Search /></el-icon>
					</template>
				</el-input>
			</div>
			
            <el-tabs v-model="activeTab" size="small" class="prop-tabs">
				<el-tab-pane v-if="!userPropsOnly" label="事件属性" name="event_data">
					<div class="prop-list">
							<!-- 特殊选项 -->
							<div 
								v-for="opt in extraOptions"
								:key="`extra-${opt.value}`"
								class="prop-item special-item"
								@click="selectProperty(opt.value, 'event_data')"
							>
								<span class="prop-name">{{ opt.label }}</span>
								<span class="prop-type">特殊</span>
							</div>
							
							<!-- 事件属性列表 -->
							<el-popover 
								v-for="p in eventPropsFiltered"
								:key="p.property_id || p.property_name"
								v-if="showFormulaPopover"
								placement="right" 
								trigger="hover" 
								:width="200" 
								:show-after="200"
								:hide-after="100"
							>
								<div class="formula-list">
									<div class="formula-title">选择公式:</div>
									<el-tag
										v-for="opt in getFormulaOptionsForProperty(p)"
										:key="opt.value"
										size="small"
										class="formula-tag"
										@click.stop="applyFormula(p, opt.value)"
									>{{ opt.label }}</el-tag>
								</div>
								<template #reference>
									<div 
										class="prop-item hoverable-item"
										@click="selectProperty(p.property_id || p.property_name || '', 'event_data')"
									>
										<div class="prop-info">
											<span class="prop-name">{{ p.property_name || p.property_id }}</span>
											<span class="prop-type">{{ getDataTypeText(p.data_type) }}</span>
										</div>
									</div>
								</template>
							</el-popover>
							
							<!-- 无公式选择时的普通项 -->
							<div 
								v-for="p in eventPropsFiltered"
								:key="p.property_id || p.property_name"
								v-if="!showFormulaPopover"
								class="prop-item"
								@click="selectProperty(p.property_id || p.property_name || '', 'event_data')"
							>
								<div class="prop-info">
									<span class="prop-name">{{ p.property_name || p.property_id }}</span>
									<span class="prop-type">{{ getDataTypeText(p.data_type) }}</span>
								</div>
							</div>
						</div>
					</el-tab-pane>
					
									<el-tab-pane v-if="!eventPropsOnly" label="用户属性" name="user_data">
					<div class="prop-list">
							<!-- 用户属性列表 -->
							<el-popover 
								v-for="p in userPropsFiltered"
								:key="p.property_id || p.property_name"
								v-if="showFormulaPopover"
								placement="right" 
								trigger="hover" 
								:width="200" 
								:show-after="200"
								:hide-after="100"
							>
								<div class="formula-list">
									<div class="formula-title">选择公式:</div>
									<el-tag
										v-for="opt in getFormulaOptionsForProperty(p)"
										:key="opt.value"
										size="small"
										class="formula-tag"
										@click.stop="applyFormula(p, opt.value)"
									>{{ opt.label }}</el-tag>
								</div>
								<template #reference>
									<div 
										class="prop-item hoverable-item"
										@click="selectProperty(p.property_id || p.property_name || '', 'user_data')"
									>
										<div class="prop-info">
											<span class="prop-name">{{ p.property_name || p.property_id }}</span>
											<span class="prop-type">{{ getDataTypeText(p.data_type) }}</span>
										</div>
									</div>
								</template>
							</el-popover>
							
							<!-- 无公式选择时的普通项 -->
							<div 
								v-for="p in userPropsFiltered"
								:key="p.property_id || p.property_name"
								v-if="!showFormulaPopover"
								class="prop-item"
								@click="selectProperty(p.property_id || p.property_name || '', 'user_data')"
							>
								<div class="prop-info">
									<span class="prop-name">{{ p.property_name || p.property_id }}</span>
									<span class="prop-type">{{ getDataTypeText(p.data_type) }}</span>
								</div>
							</div>
						</div>
					</el-tab-pane>
					
                    <!-- 用户标签选项卡 -->
					<el-tab-pane v-if="!eventPropsOnly && !userPropsOnly" label="用户标签" name="user_tag_data">
						<div class="prop-list">
							<!-- 用户标签列表 -->
							<div 
								v-for="tag in userTagsFiltered"
								:key="tag.tag_code"
								class="prop-item"
								@click="selectProperty(tag.tag_code, 'user_tag_data')"
							>
								<div class="prop-info">
									<span class="prop-name">{{ tag.tag_show_name }}</span>
									<span class="prop-type">{{ getTagTypeText(tag.tag_type) }}</span>
								</div>
							</div>
						</div>
					</el-tab-pane>

                    <!-- 用户分群选项卡 -->
                    <el-tab-pane v-if="!eventPropsOnly && !userPropsOnly" label="用户分群" name="user_group">
                        <div class="prop-list">
                            <div 
                                v-for="g in userGroupsFiltered"
                                :key="g.group_code"
                                class="prop-item"
                                @click="selectProperty(g.group_code, 'user_group')"
                            >
                                <div class="prop-info">
                                    <span class="prop-name">{{ g.group_name }}</span>
                                    <span class="prop-type">分群</span>
                                </div>
                            </div>
                        </div>
                    </el-tab-pane>
				</el-tabs>
			</div>
			
			<template #reference>
				<el-input
					ref="inputRef"
					v-model="displayValue"
					size="small"
					readonly
					:placeholder="placeholder || '选择属性字段'"
					:style="{ width: width || '100%' }"
					@click="onInputClick"
				>
					<template #suffix>
						<el-icon class="select-arrow">
							<ArrowDown />
						</el-icon>
					</template>
				</el-input>
			</template>
		</el-popover>
	</div>
</template>

<script setup lang="ts">
import { getUserGroups, getUserProperties, getUserTags } from '@/api/analytics';
import { getPropertyOptions } from '@/api/selector';
import { ArrowDown, Search } from '@element-plus/icons-vue';
import { computed, ref, watch } from 'vue';

type ColumnSource = 'event_data'|'user_data'|'user_tag_data'
interface PropertyItem {
	property_id?: string
	property_name?: string
	data_type?: string
	source?: ColumnSource
}

interface UserTag {
	id: number
	project_alias: string
	tag_code: string
	tag_show_name: string
	tag_type: string
	description?: string
	status: number
	user_count?: number
}

interface UserGroupItem {
  id: number
  group_code: string
  group_name: string
}

const props = withDefaults(defineProps<{
	modelValue?: string
	table?: ColumnSource
	projectAlias?: string
	eventId?: string
	placeholder?: string
	extraOptions?: Array<{ label: string; value: string }>
	showFormulaPopover?: boolean
	width?: string
	/** 是否为全局筛选模式，true: 从project_property表获取所有属性，false: 从project_event_property获取特定事件的属性 */
	isGlobalFilter?: boolean;
	/** 是否只显示属性，不显示用户属性标签页 */
	eventPropsOnly?: boolean;
	/** 是否只显示用户属性，不显示事件属性标签页 */
	userPropsOnly?: boolean
	/** 外部传入的用户属性列表 */
	userProperties?: any[]
	/** 外部传入的事件属性列表 */
	eventProperties?: any[]
	/** 是否在选择后自动关闭弹窗 */
	autoClose?: boolean
}>(), {
	placeholder: '选择属性字段',
	extraOptions: () => [],
	showFormulaPopover: false,
	width: '100%',
	isGlobalFilter: false,
	eventPropsOnly: false,
	userPropsOnly: false,
	userProperties: () => [],
	eventProperties: () => [],
	autoClose: true,
})
const emit = defineEmits<{
	(e: 'update:modelValue', v: string): void
	(e: 'update:table', v: ColumnSource): void
	(e: 'change', v: { field: string; table: ColumnSource; property?: PropertyItem }): void
	(e: 'apply-formula', v: { field: string; formula: number; table: ColumnSource }): void
}>()

const popoverRef = ref()
const inputRef = ref()
const innerValue = ref<string | undefined>(props.modelValue)
watch(() => props.modelValue, v => { innerValue.value = v })

const allProperties = ref<PropertyItem[]>([])
const eventPropertiesRel = ref<Array<{ event_id: string; property_id: string }>>([])
const userTagsList = ref<UserTag[]>([])
const userGroupsList = ref<UserGroupItem[]>([])
let loaded = false
const activeTab = ref<any>('event_data')
const searchKeyword = ref('')

// 如果只显示用户属性，则默认打开用户属性tab
watch(() => props.userPropsOnly, (isUserOnly) => {
  if (isUserOnly) {
    activeTab.value = 'user_data'
  }
}, { immediate: true })

async function ensureMeta() {
	if (loaded) return
	const alias = props.projectAlias || ''
	if (!alias) return
	
	// 如果外部传入了用户属性，则不再单独请求
	const userPropsPromise = Array.isArray(props.userProperties) && props.userProperties.length > 0
		? Promise.resolve({ data: props.userProperties })
		: getUserProperties(alias)

	// 如果外部传入了事件属性，则不再单独请求
	const eventPropsPromise = Array.isArray(props.eventProperties) && props.eventProperties.length > 0
		? Promise.resolve({ data: props.eventProperties })
		: getPropertyOptions(alias)

  const [propsResp, userPropsResp, relResp, userTagsResp, userGroupsResp] = await Promise.all([
		eventPropsPromise,
		userPropsPromise,
		getProjectProperties(alias),
      getUserTags(alias),
      getUserGroups(alias).then((r:any)=>({ data: (r?.data?.list)||r?.data||[] })).catch(()=>({ data: [] })),
	])
	
	// 合并属性和用户属性
	const eventPropsList = ((propsResp as any)?.data || propsResp || []).map((p: any) => ({
		property_id: p.propertyId || p.property_id || p.id,
		property_name: p.propertyName || p.property_name || p.name,
		data_type: p.dataType || p.data_type,
		source: 'event_data'
	}))
	
	const userPropsList = ((userPropsResp as any)?.data || userPropsResp || []).map((p: any) => ({
		property_id: p.propertyId || p.property_id || p.id,
		property_name: p.propertyName || p.property_name || p.name,
		data_type: p.dataType || p.data_type,
		source: 'user_data'
	}))
	
	allProperties.value = [...eventPropsList, ...userPropsList]
	eventPropertiesRel.value = ((relResp as any)?.data || relResp || []).map((rel: any) => ({
		event_id: rel.eventId || rel.event_id,
		property_id: rel.propertyId || rel.property_id
	}))
	
	// 设置用户标签数据
	userTagsList.value = ((userTagsResp as any)?.data || userTagsResp || []).map((tag: any) => ({
		id: tag.id,
		project_alias: tag.project_alias || tag.projectAlias,
		tag_code: tag.tag_code || tag.tagCode,
		tag_show_name: tag.tag_show_name || tag.tagShowName,
		tag_type: tag.tag_type || tag.tagType,
		description: tag.description,
		status: tag.status,
		user_count: tag.user_count || tag.userCount
	}))

  // 设置用户分群数据
  userGroupsList.value = (((userGroupsResp as any)?.data) || userGroupsResp || []).map((g: any) => ({
    id: g.id,
    group_code: g.group_code || g.groupCode,
    group_name: g.group_name || g.groupName,
  }))
	
	loaded = true
}

watch(() => props.projectAlias, () => { loaded = false })

const eventProps = computed<PropertyItem[]>(() => {
	const eIdRaw = props.eventId
	let list = allProperties.value.filter(p => p.source === 'event_data')
	
	// 如果是全局筛选模式，返回所有属性
	if (props.isGlobalFilter) {
		return list
	}
	
	// 如果是事件指标模式，只返回与特定事件关联的属性
	if (eIdRaw) {
		const eId = String(eIdRaw)
		const set = new Set(eventPropertiesRel.value.filter(x => String(x.event_id) === eId).map(x => x.property_id))
		const related = list.filter(p => p.property_id && set.has(p.property_id as string))
		if (related.length) list = related
	}
	return list
})
const eventPropsFiltered = computed<PropertyItem[]>(() => {
  const kw = (searchKeyword.value || '').trim().toLowerCase()
  if (!kw) return eventProps.value
  return eventProps.value.filter(p =>
    String(p.property_name || '').toLowerCase().includes(kw) ||
    String(p.property_id || '').toLowerCase().includes(kw) ||
    String(getDataTypeText(p.data_type) || '').toLowerCase().includes(kw)
  )
})

const userProps = computed<PropertyItem[]>(() => {
	// 优先使用外部传入的用户属性
	if (Array.isArray(props.userProperties) && props.userProperties.length > 0) {
		return props.userProperties.map((p: any) => ({
			property_id: p.propertyId || p.property_id,
			property_name: p.propertyName || p.property_name,
			data_type: p.dataType || p.data_type,
			source: 'user_data'
		}))
	}

	// 否则，使用内部加载的数据
	const list = allProperties.value.filter(p => p.source === 'user_data')
	if (list.length) return list
	return [{ property_id: 'u_openid', property_name: '用户ID', data_type: 'string', source: 'user_data' }]
})
const userPropsFiltered = computed<PropertyItem[]>(() => {
  const kw = (searchKeyword.value || '').trim().toLowerCase()
  if (!kw) return userProps.value
  return userProps.value.filter(p =>
    String(p.property_name || '').toLowerCase().includes(kw) ||
    String(p.property_id || '').toLowerCase().includes(kw) ||
    String(getDataTypeText(p.data_type) || '').toLowerCase().includes(kw)
  )
})

// 用户标签过滤计算属性
const userTagsFiltered = computed<UserTag[]>(() => {
  const kw = (searchKeyword.value || '').trim().toLowerCase()
  if (!kw) return userTagsList.value.filter(tag => tag.status === 1) // 只显示启用的标签
  return userTagsList.value.filter(tag => 
    tag.status === 1 && (
      String(tag.tag_show_name || '').toLowerCase().includes(kw) ||
      String(tag.tag_code || '').toLowerCase().includes(kw) ||
      String(tag.description || '').toLowerCase().includes(kw)
    )
  )
})

// 用户分群过滤计算
const userGroupsFiltered = computed<UserGroupItem[]>(() => {
  const kw = (searchKeyword.value || '').trim().toLowerCase()
  const list = userGroupsList.value || []
  if (!kw) return list
  return list.filter(g =>
    String(g.group_name || '').toLowerCase().includes(kw) ||
    String(g.group_code || '').toLowerCase().includes(kw)
  )
})

function getTagTypeText(tagType?: string): string {
  const t = String(tagType || '').toLowerCase()
  if (t === 'condition') return '条件标签'
  if (t === 'first_last') return '首末次标签'
  if (t === 'metric_value') return '指标值标签'
  if (t === 'id') return 'ID 标签'
  return '标签'
}

// 显示值计算
const displayValue = computed(() => {
	if (!innerValue.value) return ''
	
	// 检查是否是特殊选项
	const extraOpt = props.extraOptions.find(opt => opt.value === innerValue.value)
	if (extraOpt) return extraOpt.label
	
	// 检查事件属性
	const eventProp = eventProps.value.find(p => (p.property_id || p.property_name) === innerValue.value)
	if (eventProp) return eventProp.property_name || eventProp.property_id
	
	// 检查用户属性
	const userProp = userProps.value.find(p => (p.property_id || p.property_name) === innerValue.value)
	if (userProp) return userProp.property_name || userProp.property_id
	
	// 检查用户标签
	const userTag = userTagsList.value.find(tag => tag.tag_code === innerValue.value)
	if (userTag) return `${userTag.tag_show_name}`
	
	return innerValue.value
})

function onInputClick() {
	ensureMeta()
}

function selectProperty(value: string, table: ColumnSource | 'user_group') {
	innerValue.value = value
	emit('update:modelValue', value)
	// 将用户分群映射到 user_tag_data 表以便统一值选择逻辑
	const tableToEmit = table === 'user_group' ? 'user_tag_data' : table
	emit('update:table', tableToEmit as ColumnSource)
	
	// 查找属性信息
	let found: any = null
	
	if (tableToEmit === 'user_tag_data') {
		// 用户标签
		const tag = userTagsList.value.find(tag => tag.tag_code === value)
		if (tag) {
			found = {
				property_id: tag.tag_code,
				property_name: `${tag.tag_show_name}`,
				data_type: 'string', // 标签值一般为字符串类型
				source: 'user_tag_data',
				tag_type: tag.tag_type
			}
		}
		// 用户分群：映射为标签样式字段，property 显示 group_name
		if (!found) {
			const g = userGroupsList.value.find(x => x.group_code === value)
			if (g) {
				found = {
					property_id: g.group_code,
					property_name: `${g.group_name}`,
					data_type: 'string',
					source: 'user_tag_data',
					tag_type: 'user_group'
				}
			}
		}
	} else {
		// 事件属性或用户属性
		const findIn = (arr: PropertyItem[]) => arr.find(p => (p.property_id || p.property_name) === value)
		found = findIn(table === 'event_data' ? eventProps.value : userProps.value)
	}
	
	emit('change', { field: value, table: tableToEmit as ColumnSource, property: found })
	
	// 根据autoClose属性决定是否关闭弹窗
	if (props.autoClose) {
		try { popoverRef.value?.hide?.() } catch {}
	}
}

const formulaOptions = [
	{ value: 1, label: '求和' },
	{ value: 2, label: '总次数' },
	{ value: 3, label: '去重数' },
	{ value: 4, label: '用户数' },
	{ value: 5, label: '任意值' },
	{ value: 6, label: '最大值' },
	{ value: 7, label: '最小值' },
	{ value: 8, label: '平均值' },
]
function getFormulaOptionsForProperty(p: PropertyItem) {
	const dt = String(p?.data_type || '').toLowerCase()
	const numeric = [1,2,3,6,7,8]
	const stringy = [2,3,5]
	const fallback = [2,3]
	let allow: number[]
	if ([ 'int', 'integer', 'bigint', 'decimal', 'double', 'float', 'number', 'numeric' ].includes(dt)) allow = numeric
	else if ([ 'varchar', 'char', 'text', 'string' ].includes(dt)) allow = stringy
	else allow = fallback
	return formulaOptions.filter(o => allow.includes(o.value))
}
function getDataTypeText(dataType?: string): string {
	if (!dataType) return '未知类型'
	const dt = String(dataType).toLowerCase()
	
	// 数值类型
	if (['int', 'integer', 'bigint'].includes(dt)) return '整数'
	if (['decimal', 'double', 'float', 'number', 'numeric'].includes(dt)) return '数值'
	
	// 字符串类型
	if (['varchar', 'char', 'text', 'string'].includes(dt)) return '文本'
	
	// 其他类型
	if (dt === 'datetime' || dt === 'timestamp') return '时间'
	if (dt === 'date') return '日期'
	if (dt === 'boolean' || dt === 'bool') return '布尔'
	
	// 默认显示原始类型
	return dt.toUpperCase()
}

function applyFormula(p: PropertyItem, formula: number) {
	const field = p.property_id || p.property_name || ''
	const table: ColumnSource = p.source || 'event_data'
	if (!field) return
	emit('apply-formula', { field, formula, table })
	// 根据autoClose属性决定是否关闭弹窗
	if (props.autoClose) {
		try { popoverRef.value?.hide?.() } catch {}
	}
}
</script>

<style scoped>
.prop-select-wrap {
  position: relative;
}

.select-arrow {
  color: #c0c4cc;
  font-size: 14px;
  transition: transform 0.3s;
}

.prop-select-content {
  max-height: 400px;
  overflow: hidden;
}

.prop-tabs {
  margin: -8px -8px 0;
}

.prop-tabs :deep(.el-tabs__header) {
  margin-bottom: 8px;
  padding: 0 8px;
  background: #f5f7fa;
}

.prop-tabs :deep(.el-tabs__content) {
  padding: 8px 12px;
}

.prop-list {
  max-height: 300px;
  overflow-y: auto;
}

.prop-item {
  padding: 8px 12px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  border-bottom: 1px solid #f0f0f0;
  min-height: 36px;
  display: flex;
  align-items: center;
}

.prop-item:hover {
  background: #f5f7fa;
}

.prop-item.hoverable-item {
  border: 1px solid transparent;
  transition: all 0.2s;
}

.prop-item.hoverable-item:hover {
  background: #f0f9ff;
  border-color: #b3d8ff;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.1);
}

.prop-item:last-child {
  border-bottom: none;
}

.special-item {
  background: #f0f9ff;
  border: 1px solid #b3d8ff;
  margin-bottom: 8px;
}

.special-item:hover {
  background: #e1f5fe;
}

.prop-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.prop-name {
  font-size: 13px;
  color: #303133;
  font-weight: 500;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 220px;
}

.prop-type {
  font-size: 10px;
  color: #909399;
  background: #f0f0f0;
  padding: 2px 6px;
  border-radius: 3px;
  line-height: 1.3;
  flex-shrink: 0;
}



.special-item .prop-type {
  background: #409eff;
  color: white;
}

.formula-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  padding: 8px;
}

.formula-title {
  width: 100%;
  font-size: 12px;
  color: #606266;
  margin-bottom: 8px;
  font-weight: 500;
}

.formula-tag {
  cursor: pointer;
  font-size: 11px;
  height: 22px;
  line-height: 20px;
  padding: 0 8px;
  border-radius: 4px;
  transition: all 0.2s;
}

.formula-tag:hover {
  transform: scale(1.05);
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}
/* 全局搜索框样式 */
.global-search {
  padding: 12px;
  border-bottom: 1px solid #f0f0f0;
  background: #fafbfc;
}

.search-input {
  width: 100%;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 6px;
  border: 1px solid #dcdfe6;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  transition: all 0.2s ease;
}

.search-input :deep(.el-input__wrapper:hover) {
  border-color: var(--el-color-primary-light-5);
}

.search-input :deep(.el-input__wrapper.is-focus) {
  border-color: var(--el-color-primary);
  box-shadow: 0 0 0 2px var(--el-color-primary-light-9);
}

.search-input :deep(.el-input__prefix) {
  color: #909399;
}
</style>

<style>
.prop-select-popover {
  padding: 6px !important;
  z-index: 3000 !important;
}

.prop-select-popover .el-popover__title {
  margin-bottom: 12px;
}
</style>

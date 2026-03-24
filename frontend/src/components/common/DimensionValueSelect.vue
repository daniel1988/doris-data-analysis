<template>
	<!-- 时间字段：区间（BETWEEN）使用范围时间选择器 -->
	<el-date-picker
		v-if="isTimeRange"
		v-model="dateRangeValue"
		type="datetimerange"
		value-format="YYYY-MM-DD HH:mm:ss"
		format="YYYY-MM-DD HH:mm:ss"
		range-separator="至"
		start-placeholder="开始时间"
		end-placeholder="结束时间"
		size="small"
		style="width: 100%; min-width: 0;"
	/>

	<!-- 时间字段：单值时间选择器 -->
	<el-date-picker
		v-else-if="isTimeSingle"
		v-model="dateSingleValue"
		type="datetime"
		value-format="YYYY-MM-DD HH:mm:ss"
		format="YYYY-MM-DD HH:mm:ss"
		:placeholder="getPlaceholder()"
		size="small"
		style="width: 100%; min-width: 0;"
	/>

	<!-- 非时间字段：下拉选择（原逻辑） -->
	<el-select
		v-else
		v-model="innerValue"
		filterable
		reserve-keyword
		clearable
		allow-create
		default-first-option
		:multiple="isMultiple"
		:collapse-tags="isMultiple"
		:collapse-tags-tooltip="isMultiple"
		:max-collapse-tags="2"
		:loading="loading"
		size="small"
		:placeholder="getPlaceholder()"
		style="width: 100%; min-width: 0;"
		@focus="handleFocus"
	>
		<el-option
			v-for="opt in allValues"
			:label="opt"
			:value="opt"
			:key="opt"
		/>
	</el-select>
</template>

<script setup lang="ts">
import { dimensionsApi } from '@/api/dimensions';
import { FilterOperators as OP } from '@/types/doris/constants';
import { computed, ref, watch } from 'vue';

const props = withDefaults(defineProps<{
	modelValue?: string | string[]
	projectAlias?: string
	table?: 'event_data'|'user_data'|'user_tag_data'
	field?: string
	eventId?: string
	placeholder?: string
	operator?: number  // 操作符：10=包含，11=不包含
	multiple?: boolean
	dataType?: string   // 可选：字段数据类型（如 datetime/timestamp/date）
}>(), {
	placeholder: '输入以搜索或直接回车确认',
	multiple: false,
})

const emit = defineEmits<{ (e: 'update:modelValue', v?: string | string[]): void }>()

// 判断是否为多选模式（包含或不包含操作符，或显式指定）
const isMultiple = computed(() => props.multiple || props.operator === OP.In || props.operator === OP.NotIn)

// 判断是否为时间类型字段（优先使用 dataType，其次用字段名启发式判断）
const isTimeField = computed(() => {
	const dt = String(props.dataType || '').toLowerCase()
	if (['datetime', 'timestamp', 'date'].includes(dt)) return true
	const f = String(props.field || '').toLowerCase()
	return f === 'e_event_time' || f.endsWith('time') || f.endsWith('_time') || f.includes('date')
})

// 是否时间区间（仅在 BETWEEN 时使用范围选择）
const isTimeRange = computed(() => isTimeField.value && props.operator === OP.Between)
// 是否单值时间
const isTimeSingle = computed(() => isTimeField.value && !isTimeRange.value)

const innerValue = ref<string | string[] | undefined>(props.modelValue)

// 监听props变化
watch(() => props.modelValue, v => { 
	innerValue.value = v 
})

// 监听innerValue变化并emit
watch(innerValue, v => {
	if (isTimeRange.value) {
		const range = Array.isArray(v) ? (v as string[]) : []
		const normalized = range.length === 2 ? [range[0], range[1]] : []
		emit('update:modelValue', normalized)
		return
	}
	if (isMultiple.value) {
		// 多选模式：确保返回数组
		const arrayValue = Array.isArray(v) ? v : (v ? [v] : [])
		emit('update:modelValue', arrayValue)
	} else {
		// 单选模式：确保返回字符串
		const stringValue = Array.isArray(v) ? v[0] : v
		emit('update:modelValue', stringValue ?? '')
	}
})

// 当操作符变化时，按需规范化值形态
watch(() => props.operator, () => {
	if (isTimeRange.value) {
		// 确保为 [start, end] 结构
		if (!Array.isArray(innerValue.value) || (innerValue.value as string[]).length !== 2) {
			innerValue.value = []
		}
	} else if (isTimeSingle.value) {
		// 确保为单个字符串
		if (Array.isArray(innerValue.value)) {
			innerValue.value = innerValue.value[0] || ''
		}
	}
})

const allValues = ref<string[]>([])
const loading = ref(false)
let loadedOnce = false

// 根据模式获取占位符
const getPlaceholder = () => {
	if (isTimeRange.value) {
		return '选择时间范围'
	}
	if (isTimeSingle.value) {
		return '选择时间'
	}
	if (isMultiple.value) {
		return props.placeholder || '输入以搜索或选择多个值'
	}
	return props.placeholder || '输入以搜索或直接回车确认'
}

async function ensureLoaded() {
	if (loadedOnce) return
	const alias = props.projectAlias || ''
	if (!alias || !props.field || !props.table) return
	try {
		loading.value = true
		const resp = await dimensionsApi.listValues({
			project_alias: alias,
			table: props.table,
			field: props.field,
			e_event_id: props.eventId || undefined,
		})
		
		// 根据 interceptor 的逻辑，resp 应该是 res.data，即直接的数组或包含 rows 的对象
		let list: any[] = []
		if (Array.isArray(resp)) {
			list = resp
		} else if ((resp as any)?.rows) {
			list = (resp as any).rows
		} else if ((resp as any)?.data) {
			// 万一拦截器没生效或者返回结构多了一层
			list = Array.isArray((resp as any).data) ? (resp as any).data : ((resp as any).data.rows || [])
		}

		// 增强兼容性：后端可能返回 [{ value: "xxx" }] 或 [{ fieldName: "xxx" }]
		const values: string[] = []
		for (const row of list) {
			// 1. 优先尝试 "value" 键 (后端统一后的 key)
			let v = row?.["value"]
			// 2. 其次尝试请求的 field 键 (旧逻辑兼容)
			if (v === undefined || v === null) {
				v = row?.[props.field!]
			}
			// 3. 如果还是没有且 row 是简单类型，则直接使用 row
			if (v === undefined || v === null) {
				if (typeof row !== 'object') v = row
			}
			
			if (v !== null && v !== undefined) values.push(String(v))
		}
		// 去重并排序
		allValues.value = Array.from(new Set(values)).sort()
		loadedOnce = true
	} finally {
		loading.value = false
	}
}

function handleFocus() {
	// 仅在下拉选择模式下加载可选值
	if (!isTimeField.value && !loadedOnce) ensureLoaded()
}

// 当字段、表、项目变化时，重置
watch(() => [props.projectAlias, props.table, props.field, props.eventId], () => {
	loadedOnce = false
	allValues.value = []
})

// 当必要参数就绪后，自动尝试加载一次（避免用户必须再次聚焦才能发起请求）
watch(() => [props.projectAlias, props.table, props.field], async ([alias, table, field]) => {
    const hasReady = Boolean(alias && table && field)
    if (!hasReady) return
    if (isTimeField.value) return
    if (!loadedOnce) {
        try { await ensureLoaded() } catch {}
    }
}, { immediate: false })

// 便捷的 v-model 适配：时间选择绑定
const dateSingleValue = computed<string | ''>({
	get() {
		return (Array.isArray(innerValue.value) ? (innerValue.value[0] || '') : (innerValue.value as string)) || ''
	},
	set(v: string | '') {
		innerValue.value = v || ''
	}
})

const dateRangeValue = computed<[string, string] | null>({
	get() {
		if (Array.isArray(innerValue.value) && innerValue.value.length === 2) {
			return [String(innerValue.value[0] || ''), String(innerValue.value[1] || '')]
		}
		return null
	},
	set(v: [string, string] | null) {
		innerValue.value = v ? [v[0], v[1]] : []
	}
})
</script>

<style scoped>
</style>

<template>
  <div class="custom-formula-input" ref="formulaInputRef">
    <el-popover
      :visible="isPanelVisible"
      placement="bottom"
      :width="280"
      trigger="manual"
      popper-class="formula-popover"
      :virtual-ref="formulaInputRef"
      virtual-triggering
    >
      <FormulaPanel
        @add-token="addToken"
        @delete-token="deleteLastToken"
        @clear="clearAllTokens"
        @insert-event="onInsertEvent"
        @close="isPanelVisible = false"
      />
    </el-popover>

    <div 
      class="formula-display-area" 
      :class="{ 'is-focused': isPanelVisible }"
      tabindex="0"
      ref="displayAreaRef"
      @click="onDisplayAreaClick"
      @keydown="handleKeyDown"
    >
      <template v-if="localTokens.length === 0">
        <span class="placeholder">点击编辑公式</span>
          </template>
          <template v-else>
        <div v-for="(element, index) in localTokens" :key="element._uid" class="token-wrapper">
          <span v-if="String(element.type) === FormulaTokenTypes.Operator || String(element.type) === FormulaTokenTypes.Number" :class="`token-${String(element.type) === FormulaTokenTypes.Number ? 'number' : 'operator'}`">
            {{ element.value }}
          </span>
          <!-- Event Metric Token (Full Config UI) -->
          <div v-if="String(element.type) === FormulaTokenTypes.Metric" class="inline-event-metric" @click.stop>
            <el-select v-model="element.value.e_event_id" placeholder="选择事件" size="small" filterable clearable @change="() => handleEventChangeInFormula(element)">
              <el-option
                v-for="ev in projectEvents"
                :key="ev.id"
                :label="ev.name ? `${ev.name} (${ev.id})` : ev.id"
                :value="ev.id"
              />
                  </el-select>
            <el-popover :visible="activeFilterTokenUid === element._uid" placement="bottom-start" :width="520" trigger="manual">
               <template #reference>
                 <el-button size="small" :icon="Filter" :type="hasFilters(element.value) ? 'primary' : 'default'" @click.stop="handleFilterIconClick(element)">
                   过滤 {{ getFilterCount(element.value) }}
                  </el-button>
               </template>
               <!-- Full Filter UI inside Popover -->
               <div class="metric-filter-popover">
                 <div class="filter-header">
                   <span class="filter-title">事件过滤条件</span>
                   <div style="display:flex; gap:6px; align-items:center;">
                     <el-button size="small" type="primary" @click="addFilter(element.value)"><el-icon><Plus /></el-icon></el-button>
                   </div>
                 </div>
                 <div v-if="element.value.filter_group?.filters?.length > 0" class="filter-content-popover">
                   <div class="gf-header" v-if="element.value.filter_group.filters.length > 1">
                     <div class="relation-info">
                       <el-button size="small" class="relation-switch-btn" @click="toggleFilterScope(element.value)">{{ getFilterScopeText(element.value) }}</el-button>
                     </div>
                   </div>
                   <div class="gf-list">
                     <div v-for="(f, fIdx) in element.value.filter_group.filters" :key="fIdx" class="gf-item-wrapper">
                       <BaseFilterItem 
                         v-model="element.value.filter_group.filters[fIdx]"
                         :event-id="element.value.e_event_id"
                         :is-global-filter="false"
                         @remove="removeFilter(element.value, fIdx)"
                       />
                     </div>
                </div>
              </div>
                 <div v-else class="filter-content-empty">点击右上角 "+" 添加过滤条件</div>
                 
                 <!-- 确认按钮区域 -->
                 <div class="filter-footer">
                   <el-button size="small" @click="toggleFilter(element)">取消</el-button>
                   <el-button size="small" type="primary" @click="confirmFilter(element)">确认</el-button>
                 </div>
        </div>
             </el-popover>
            <div class="property-formula-group">
              <MetricPropertyPicker
                :event-id="element.value.e_event_id"
                :field="element.value.metric.column.field"
                :table="element.value.metric.column.table"
                :formula="element.value.metric.formula"
                :auto-close="false"
                @select="(payload: any) => handleApplyFormula(element, payload)"
              />
            </div>
      </div>
      </div>
      </template>
      <!-- Blinking cursor simulation -->
      <span v-if="isPanelVisible" class="blinking-cursor"></span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { FormulaDefaults, FormulaTokenTypes, TableNames } from '@/constants/analysis';
import { MetricFormulas as MF } from '@/types/doris/constants';
import type { CustomElement, EventMetric } from '@/types/doris/customFormula';
import { Filter, Plus } from '@element-plus/icons-vue';
import { computed, inject, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { useCustomFormula } from '../composables/useCustomFormula';
import { ANALYSIS_CONTEXT_KEY } from '../context';
import FormulaPanel from './metrics/FormulaPanel.vue';
import MetricPropertyPicker from './metrics/MetricPropertyPicker.vue';
import BaseFilterItem from './shared/BaseFilterItem.vue';

const props = defineProps<{ 
  modelValue: CustomElement[]
}>();
const emit = defineEmits(['update:modelValue', 'request-new-event', 'open-filter']);

const { normalizeMetricField } = useCustomFormula()

const context = inject(ANALYSIS_CONTEXT_KEY)
const projectEvents = computed(() => context?.state.metadata.eventOptions || [])

const isPanelVisible = ref(false);
const localTokens = ref<CustomElement[]>([]);
const formulaInputRef = ref<HTMLDivElement | null>(null);
const displayAreaRef = ref<HTMLDivElement | null>(null);
const activeFilterTokenUid = ref<number | string | null>(null);

// --- Click Outside Logic ---
const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as HTMLElement | null;
  if (!target) return;

  if (activeFilterTokenUid.value !== null) {
    if (!target.closest('.metric-filter-popover') && !target.closest('.el-popper')) {
      return;
    }
  }

  if (target.closest && (
    target.closest('.el-popper') ||
    target.closest('.metric-filter-popover') ||
    target.closest('.prop-select-popover')
  )) {
    return;
  }

  if (formulaInputRef.value && !formulaInputRef.value.contains(target as Node)) {
    isPanelVisible.value = false;
  }
};

onMounted(() => {
  document.addEventListener('click', handleClickOutside, true);
  localTokens.value = JSON.parse(JSON.stringify(props.modelValue || []));
  ensureEventMetricDefaults();
});

watch(() => props.modelValue, (newVal) => {
  localTokens.value = JSON.parse(JSON.stringify(newVal || []));
  ensureEventMetricDefaults();
}, { deep: true, immediate: true });

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside, true);
});

function ensureEventMetricDefaults() {
  localTokens.value.forEach((t, i) => {
    if (!t._uid) {
      t._uid = Date.now() + i;
    }
    const tt = String(t?.type)
    if (tt === FormulaTokenTypes.Metric || tt === '1') {
      const em: any = t.value || (t.value = {})
      if (em.e_event_id === undefined) em.e_event_id = ''
      if (em.type === undefined) em.type = 1 // Normal event type
      if (!em.metric) em.metric = {}
      if (!em.metric.column) em.metric.column = { table: TableNames.EVENT, field: '' }
      if (em.metric.formula === undefined || em.metric.formula === null) em.metric.formula = MF.Count
      if (!em.filter_group) em.filter_group = { scope: 1, filters: [] }
    }
  })
}

function handleFilterIconClick(element: any) {
  if (!element?._uid) return;
  activeFilterTokenUid.value = element._uid;
}

function confirmFilter(element: any) {
  if (!element?.value) return
  activeFilterTokenUid.value = null;
  updateModel()
}

function toggleFilterScope(metric: any) {
  if (metric && metric.filter_group) {
    metric.filter_group.scope = metric.filter_group.scope === 1 ? 2 : 1;
    updateModel()
  }
}

function getFilterScopeText(metric: any): string {
  return metric?.filter_group?.scope === 1 ? '且' : '或'
}

const updateModel = () => {
  emit('update:modelValue', localTokens.value);
};

function handleEventChangeInFormula(element: CustomElement) {
  const eventMetric = element.value as EventMetric;
  if (!eventMetric.e_event_id) {
    updateModel();
    return;
  }
  
  const selectedEvent = projectEvents.value.find(e => e.id === eventMetric.e_event_id);
  const eventName = selectedEvent?.name || eventMetric.e_event_id;
  
  const formulaMap: Record<number, string> = {
    1: '求和', 2: '总次数', 3: '去重数', 4: '用户数', 5: '任意值',
    6: '最大值', 7: '最小值', 8: '平均值'
  };
  
  const field = eventMetric.metric.column.field;
  const formula = eventMetric.metric.formula || 2;
  
  let metricDisplayName: string;
  if (field === FormulaDefaults.TOTAL_TIMES_FIELD) {
    metricDisplayName = '总次数';
  } else if (field === FormulaDefaults.TOTAL_USERS_FIELD) {
    metricDisplayName = '总用户数';
  } else if (field) {
    metricDisplayName = formulaMap[formula] || field;
  } else {
    metricDisplayName = '总次数';
  }
  
  eventMetric.name = `${eventName}.${metricDisplayName}`;
  updateModel();
}

const addToken = (token: Pick<CustomElement, 'type' | 'value'>) => {
  const t = String(token.type);
  if (t === FormulaTokenTypes.Number) {
    const v = String(token.value ?? '');
    if (v === '.') {
      appendDot();
      return;
    }
    if (/^\d$/.test(v)) {
      appendDigit(v);
      return;
    }
  }
  if (t === FormulaTokenTypes.Operator) {
    addOperatorToken(String(token.value ?? ''));
    return;
  }
  localTokens.value.push({ ...token, _uid: Date.now() } as CustomElement);
  ensureEventMetricDefaults();
  updateModel();
};

const deleteLastToken = () => {
  handleBackspaceOrDelete();
};

const clearAllTokens = () => {
  localTokens.value = [];
  updateModel();
};

const onInsertEvent = () => {
  const newEventToken = {
    type: FormulaTokenTypes.Metric,
    _uid: Date.now(),
    value: {
      e_event_id: '',
      name: '新事件',
      type: 1, // EventType.Normal
      metric: { column: { table: TableNames.EVENT, field: '' }, formula: MF.Count },
      filter_group: { scope: 1, filters: [] }
    }
  };
  localTokens.value.push(newEventToken as any);
  ensureEventMetricDefaults();
  updateModel();
};

const hasFilters = (metric: EventMetric) => (metric.filter_group?.filters?.length ?? 0) > 0;
const getFilterCount = (metric: EventMetric) => {
  const len = metric?.filter_group?.filters?.length ?? 0
  return len > 0 ? `(${len})` : ''
};

const addFilter = (metric: EventMetric) => {
  if (!metric.filter_group) {
    metric.filter_group = { scope: 1, filters: [] };
  }
  metric.filter_group.filters.push({ field: '', operator: '1', value: '' });
  updateModel();
};

const removeFilter = (metric: EventMetric, index: number) => {
  metric.filter_group?.filters.splice(index, 1);
  updateModel();
};

const toggleFilter = (element: any) => {
  activeFilterTokenUid.value = null;
};

const onDisplayAreaClick = () => {
  isPanelVisible.value = true;
  nextTick(() => displayAreaRef.value?.focus());
};

function handleKeyDown(e: KeyboardEvent) {
  if (!isPanelVisible.value) return;
  if (isTypingInInput(e.target)) return;

  const key = e.key;
  if (/^[0-9]$/.test(key)) {
    appendDigit(key);
    e.preventDefault();
    e.stopPropagation();
    return;
  }
  if (key === '.' || key === 'Decimal') {
    appendDot();
    e.preventDefault();
    e.stopPropagation();
    return;
  }
  if (key === '+' || key === '-' || key === '*' || key === '/' || key === '(' || key === ')') {
    addOperatorToken(key);
    e.preventDefault();
    e.stopPropagation();
    return;
  }
  if (key === 'Backspace' || key === 'Delete') {
    handleBackspaceOrDelete();
    e.preventDefault();
    e.stopPropagation();
    return;
  }
  if (key === 'Enter' || key === 'Escape') {
    isPanelVisible.value = false;
    e.preventDefault();
    e.stopPropagation();
    return;
  }
}

function handleApplyFormula(element: CustomElement, payload: any) {
  if (payload.event) payload.event.stopPropagation();

  const em = element.value as EventMetric;
  em.metric.column.field = payload.field;
  em.metric.column.table = payload.table;
  em.metric.formula = payload.formula;
  
  if (payload.field === FormulaDefaults.TOTAL_TIMES_FIELD || payload.field === FormulaDefaults.TOTAL_USERS_FIELD) {
    normalizeMetricField(element.value as any);
  }
  
  handleEventChangeInFormula(element);
}

function isTypingInInput(el: EventTarget | null): boolean {
  const node = el as HTMLElement | null;
  if (!node) return false;
  const tag = node.tagName;
  if (tag === 'INPUT' || tag === 'TEXTAREA') return true;
  if ((node as any).isContentEditable) return true;
  if (node.closest && node.closest('.el-input, .el-select, .el-cascader, .el-autocomplete')) return true;
  return false;
}

function appendDigit(d: string) {
  const tokens = localTokens.value;
  const last = tokens[tokens.length - 1];
  if (last && String(last.type) === FormulaTokenTypes.Number) {
    last.value = String(last.value ?? '') + d;
  } else {
    tokens.push({ type: FormulaTokenTypes.Number, value: d, _uid: Date.now() } as CustomElement);
  }
  updateModel();
}

function appendDot() {
  const tokens = localTokens.value;
  const last = tokens[tokens.length - 1];
  if (last && String(last.type) === FormulaTokenTypes.Number) {
    const val = String(last.value ?? '');
    if (!val.includes('.')) {
      last.value = val.length > 0 ? val + '.' : '0.';
      updateModel();
    }
  } else {
    tokens.push({ type: FormulaTokenTypes.Number, value: '0.', _uid: Date.now() } as CustomElement);
    updateModel();
  }
}

function addOperatorToken(op: string) {
  localTokens.value.push({ type: FormulaTokenTypes.Operator, value: op, _uid: Date.now() } as CustomElement);
  updateModel();
}

function handleBackspaceOrDelete() {
  const tokens = localTokens.value;
  if (tokens.length === 0) return;
  const last = tokens[tokens.length - 1];
  if (String(last.type) === FormulaTokenTypes.Number) {
    const val = String(last.value ?? '');
    if (val.length > 1) {
      last.value = val.slice(0, -1);
    } else {
      tokens.pop();
    }
  } else {
    tokens.pop();
  }
  updateModel();
}
</script>

<style scoped>
.custom-formula-input {
  width: 100%;
  min-width: 350px;
  box-sizing: border-box;
}

.formula-display-area {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
  border: none;
  border-bottom: 1px solid #dcdfe6;
  border-radius: 4px 4px 0 0;
  padding: 8px 12px;
  min-height: 48px;
  min-width: 300px;
  width: 100%;
  background-color: #f7f8fa;
  transition: all 0.2s;
  cursor: text;
  position: relative;
  box-sizing: border-box;
}
.formula-display-area:hover {
  background-color: #f1f3f5;
}

.formula-display-area.is-focused {
  border-bottom-color: var(--el-color-primary);
  border-bottom-width: 2px;
  box-shadow: none;
}

.blinking-cursor {
  width: 1px;
  height: 1.1em;
  background-color: #333;
  animation: blink 1s infinite;
  margin-left: 2px;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}

.token-wrapper { 
  display: inline-flex;
  align-items: center;
  gap: 4px; 
}

.token-number, .token-operator {
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
  font-size: 14px;
  height: 32px;
  line-height: 28px;
}
.token-number {
  background-color: #ecf5ff;
  color: #409eff;
}
.token-operator {
  background-color: #fdf6ec;
  color: #e6a23c;
}

.inline-event-metric { 
  display: inline-flex;
  gap: 6px; 
  align-items: center;
  background-color: #ffffff;
  padding: 4px 8px;
  border-radius: 4px;
}
.property-formula-group {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
.formula-popover { padding: 0 !important; }
.metric-filter-popover {
  padding: 12px;
}
.filter-content-popover {
  position: relative;
  padding-left: 20px;
  margin-top: 12px;
}
.filter-content-popover::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0; 
  bottom: 0;
  width: 2px;
  background-color: var(--el-border-color-lighter);
  border-radius: 1px;
}
.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}
.filter-content-empty {
  text-align: center;
  color: #909399;
  padding: 20px 0;
}
.gf-header {
  position: absolute; 
  top: 50%; 
  left: -16px; 
  transform: translateY(-50%);
  z-index: 10; 
  background: transparent;
  padding: 0;
  box-shadow: none;
  margin-bottom: 0;
}
.gf-list {
  margin-top: 0;
}
.gf-item-wrapper { 
  display: flex; 
  flex-direction: column; 
  gap: 4px; 
  margin: 0; 
}
.filter-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 12px 0 0 0;
  border-top: 1px solid #f0f0f0;
  margin-top: 12px;
}
.relation-switch-btn {
  min-width: 32px;
  height: 32px;
  border-radius: 50%;
  border: 1px solid var(--el-border-color-lighter);
  box-shadow: 0 1px 3px rgba(0,0,0,0.08);
  padding: 0;
  font-weight: bold;
}
.relation-info { 
  display: flex; 
  align-items: center; 
  gap: 8px; 
  font-size: 12px; 
  color: var(--el-text-color-regular); 
}
</style>

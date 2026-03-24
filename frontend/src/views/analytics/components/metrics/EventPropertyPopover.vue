<template>
  <el-popover
    :visible="visible"
    placement="bottom-start"
    :width="580"
    trigger="manual"
    popper-class="event-property-popover"
    :hide-after="0"
    :show-after="0"
    @update:visible="handlePopoverVisibleChange"
  >
    <div class="popover-content" @click.stop>
      <div class="selector-sentence">
        <span class="text">将</span>
        <span class="event-name">{{ currentEventName }}</span>
        <span class="text">的</span>
        
        <!-- 事件属性下拉 -->
        <el-select 
          v-model="selectedPropertyId" 
          placeholder="选择属性" 
          size="small"
          class="property-selector"
          popper-append-to-body
          filterable
          @change="onPropertyChange"
          @click.stop
          @visible-change="onPropertyDropdownVisibleChange"
        >
          <!-- 特殊选项优先展示 -->
          <el-option label="总次数" value="__TOTAL_TIMES__" @click.stop />
          <el-option label="总用户数" value="__TOTAL_USERS__" @click.stop />
          <!-- 事件属性列表 -->
          <el-option 
            v-for="property in eventProperties" 
            :key="property.property_id" 
            :label="formatPropertyLabel(property)" 
            :value="property.property_id"
            @click.stop
          />
        </el-select>
        
        <span class="text">的</span>
        
        <!-- 公式列表 -->
        <el-select 
          v-model="selectedFormulaValue" 
          placeholder="选择公式" 
          size="small"
          class="formula-selector"
          popper-append-to-body
          @change="onFormulaChange"
          @click.stop
          @visible-change="onFormulaDropdownVisibleChange"
        >
          <el-option 
            v-for="formula in availableFormulas" 
            :key="formula.value" 
            :label="formula.label" 
            :value="formula.value"
            @click.stop
          />
        </el-select>
        
        <span class="text">作为指标</span>
      </div>
      
      <!-- 确认按钮 -->
      <div class="action-buttons">
        <el-button size="small" @click="emit('update:visible', false)">取消</el-button>
        <el-button type="primary" size="small" @click="confirmSelection" :disabled="!canConfirm">确认</el-button>
      </div>
    </div>
    
    <template #reference>
      <div></div>
    </template>
  </el-popover>
</template>

<script setup lang="ts">
import { getProjectProperties } from '@/api/meta';
import { getEventOptions, getPropertyOptions } from '@/api/selector';
import { MetricFormulas as MF, TableNames } from '@/types/doris/constants';
import { computed, onMounted, ref, watch } from 'vue';

interface PropertyItem {
  property_id: string;
  property_name?: string;
  data_type?: string;
}

interface SelectionItem {
  field: string;
  formula: number;
  table: string;
  propertyName: string;
  label: string;
  eventName: string;
}

const props = withDefaults(defineProps<{
  visible: boolean;
  projectAlias: string;
  eventId: string;
  /** 确认后是否自动关闭弹层（由上层控制） */
  closeOnConfirm?: boolean;
  /** 父组件传入的初始字段（用于回显） */
  initialField?: string;
  /** 父组件传入的初始公式（用于回显） */
  initialFormula?: number | null;
  /** 父组件传入的属性名称（可选，仅作显示，不影响回显逻辑） */
  initialPropertyName?: string;
}>(), {
  closeOnConfirm: true
});

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void;
  (e: 'select', selection: SelectionItem): void;
}>();

const eventProperties = ref<PropertyItem[]>([]);
const projectEvents = ref<any[]>([]);
const selectedPropertyId = ref<string>('');
const selectedFormulaValue = ref<number | null>(null);
const isSelectingFromDropdown = ref<boolean>(false);

// 计算可用的公式选项
const availableFormulas = computed(() => {
  if (!selectedPropertyId.value) return [];
  
  // 特殊属性的公式
  if (selectedPropertyId.value === '__TOTAL_TIMES__') {
    return [{ value: MF.Count, label: '总次数' }];
  }
  if (selectedPropertyId.value === '__TOTAL_USERS__') {
    return [{ value: MF.CountDistinctUserId, label: '总用户数' }];
  }
  
  // 普通属性的公式
  const property = eventProperties.value.find(p => p.property_id === selectedPropertyId.value);
  if (property) {
    return getFormulaOptions(property);
  }
  
  return [];
});

// 检查是否可以确认
const canConfirm = computed(() => {
  return selectedPropertyId.value && selectedFormulaValue.value !== null;
});

// 计算当前事件名称
const currentEventName = computed(() => {
  if (!props.eventId) return '未选择事件';
  const event = projectEvents.value.find(e => (e.eventId || e.event_id) === props.eventId);
  return event?.eventName || event?.event_name || event?.eventId || event?.event_id || props.eventId;
});

function formatPropertyLabel(property: PropertyItem): string {
  const name = property.property_name || property.property_id;
  const type = getDataTypeText(property.data_type);
  // 使用不换行空格 (U+00A0) 来确保空格被正确渲染
  return `${name}\u00A0\u00A0\u00A0\u00A0${type}`;
}


// 公式选项
const formulaOptions = [
  { value: MF.Sum, label: '求和' },
  { value: MF.Count, label: '总次数' },
  { value: MF.CountDistinct, label: '去重数' },
  { value: MF.CountDistinctUserId, label: '用户数' },
  { value: MF.CountDistinctDailyUserId, label: '按天去重数' },
  { value: MF.Any, label: '任意值' },
  { value: MF.Max, label: '最大值' },
  { value: MF.Min, label: '最小值' },
  { value: MF.Avg, label: '平均值' },
];

// 根据数据类型获取可用公式
function getFormulaOptions(property: PropertyItem) {
  const dt = String(property?.data_type || '').toLowerCase();
  
  // 数值类
  const numericTypes = ['int', 'integer', 'bigint', 'decimal', 'double', 'float', 'number', 'numeric'];
  const numericFormulas = [MF.Count, MF.CountDistinct, MF.Sum, MF.Avg, MF.Max, MF.Min];

  // 字符串类
  const stringTypes = ['varchar', 'char', 'text', 'string'];
  const stringFormulas = [MF.Count, MF.CountDistinct, MF.CountDistinctDailyUserId, MF.Any];

  // 时间类
  const timeTypes = ['datetime', 'timestamp', 'date'];
  const timeFormulas = [MF.Count, MF.CountDistinct, MF.Any];

  let allow: number[];
  
  if (numericTypes.includes(dt)) {
    allow = numericFormulas;
  } else if (stringTypes.includes(dt)) {
    allow = stringFormulas;
  } else if (timeTypes.includes(dt)) {
    allow = timeFormulas;
  } else {
    // 默认回退（例如布尔型等）
    allow = [MF.Count, MF.CountDistinct, MF.Any];
  }
  
  return formulaOptions.filter(o => allow.includes(o.value));
}

// 获取数据类型显示文本
function getDataTypeText(dataType?: string): string {
  if (!dataType) return '未知';
  const dt = String(dataType).toLowerCase();
  
  if (['int', 'integer', 'bigint'].includes(dt)) return '整数';
  if (['decimal', 'double', 'float', 'number', 'numeric'].includes(dt)) return '数值';
  if (['varchar', 'char', 'text', 'string'].includes(dt)) return '文本';
  if (dt === 'datetime' || dt === 'timestamp') return '时间';
  if (dt === 'date') return '日期';
  if (dt === 'boolean' || dt === 'bool') return '布尔';
  
  return dt.toUpperCase();
}

// 处理属性变化
function onPropertyChange(silent = false) {
  // 标记正在进行选择操作
  if (!silent) isSelectingFromDropdown.value = true;
  
  // 在变更时尽量保留当前已选公式（若仍然可用）
  const previousFormula = selectedFormulaValue.value;
  selectedFormulaValue.value = null;
  
  // 如果有可用公式，自动选择第一个
  if (availableFormulas.value.length > 0) {
    if (previousFormula !== null && availableFormulas.value.some(f => f.value === previousFormula)) {
      // 保留原有公式
      selectedFormulaValue.value = previousFormula;
    } else {
      // 回退为第一个可用公式
      selectedFormulaValue.value = availableFormulas.value[0].value;
    }
  }
  
  // 延迟重置选择状态
  setTimeout(() => {
    if (!silent) isSelectingFromDropdown.value = false;
  }, 200);
}

// 处理公式变化
function onFormulaChange() {
  // 标记正在进行选择操作
  isSelectingFromDropdown.value = true;
  
  // 延迟重置选择状态
  setTimeout(() => {
    isSelectingFromDropdown.value = false;
  }, 200);
}

// 处理弹层可见性变化
function handlePopoverVisibleChange(visible: boolean) {
  // 如果正在进行下拉选择，阻止弹层关闭
  if (!visible && isSelectingFromDropdown.value) {
    isSelectingFromDropdown.value = false;
    return;
  }
  emit('update:visible', visible);
}

// 处理属性下拉框可见性变化
function onPropertyDropdownVisibleChange(visible: boolean) {
  if (visible) {
    isSelectingFromDropdown.value = true;
  } else {
    setTimeout(() => {
      isSelectingFromDropdown.value = false;
    }, 100);
  }
}

// 处理公式下拉框可见性变化  
function onFormulaDropdownVisibleChange(visible: boolean) {
  if (visible) {
    isSelectingFromDropdown.value = true;
  } else {
    setTimeout(() => {
      isSelectingFromDropdown.value = false;
    }, 100);
  }
}

// 确认选择
function confirmSelection() {
  if (!canConfirm.value) return;
  
  const propertyId = selectedPropertyId.value;
  const propertyName = getPropertyName(propertyId);
  const formulaLabel = availableFormulas.value.find(f => f.value === selectedFormulaValue.value)?.label || '';
  
  let field = propertyId;
  let formula = selectedFormulaValue.value!;

  if (propertyId === '__TOTAL_TIMES__') {
    field = 'e_openid'; // 特殊处理
    formula = MF.Count;
  } else if (propertyId === '__TOTAL_USERS__') {
    field = 'e_openid'; // 特殊处理
    formula = MF.CountDistinctUserId;
  }

  const selection: SelectionItem = {
    field: field,
    formula: formula,
    table: TableNames.EVENT_TABLE,
    propertyName,
    label: selectedPropertyId.value.startsWith('__') ? propertyName : `${propertyName} · ${formulaLabel}`,
    eventName: currentEventName.value // 添加事件名称
  };
  
  emit('select', selection);
  // 根据配置在确认后关闭弹层
  emit('update:visible', false);
}

// 获取属性名称
function getPropertyName(propertyId: string): string {
  if (propertyId === '__TOTAL_TIMES__') return '总次数';
  if (propertyId === '__TOTAL_USERS__') return '总用户数';
  
  const property = eventProperties.value.find(p => p.property_id === propertyId);
  return property?.property_name || property?.property_id || propertyId;
}

// 加载事件属性
async function loadEventProperties() {
  if (!props.projectAlias || !props.eventId) return;
  
  try {
    const [propertiesResp, eventPropsResp, eventsResp] = await Promise.all([
      getPropertyOptions(props.projectAlias),
      getProjectProperties(props.projectAlias),
      getEventOptions(props.projectAlias)
    ]);
    
    // 处理项目事件数据
    projectEvents.value = ((eventsResp as any)?.data || eventsResp || []).map((e: any) => ({
      event_id: e.eventId || e.event_id || e.id,
      event_name: e.eventName || e.event_name || e.name,
    }));
    
    const allProperties = ((propertiesResp as any)?.data || propertiesResp || []).map((p: any) => ({
      property_id: p.propertyId || p.property_id || p.id,
      property_name: p.propertyName || p.property_name || p.name,
      data_type: p.dataType || p.data_type,
    }));
    
    const eventPropertyRelations = ((eventPropsResp as any)?.data || eventPropsResp || []).map((rel: any) => ({
      event_id: rel.eventId || rel.event_id,
      property_id: rel.propertyId || rel.property_id
    }));
    
    // 筛选与当前事件相关的属性
    const currentEventId = String(props.eventId || '');
    const relatedPropertyIds = new Set(
      eventPropertyRelations
        .filter((rel: any) => String(rel.event_id) === currentEventId)
        .map((rel: any) => String(rel.property_id))
    );

    const filteredByRelation = allProperties.filter((p: any) => {
      const pid = String(p.property_id || '');
      return pid && relatedPropertyIds.has(pid);
    });

    // 如果当前事件没有维护关联关系，则回退为显示全部项目属性
    eventProperties.value = filteredByRelation.length > 0 ? filteredByRelation : allProperties;

    // 优先根据父组件传入的初始值进行回显（如果提供）
    const appliedFromInitial = applyInitialSelectionFromParent();

    // 自动选择默认项（"总次数"），除非已有有效选择
    const currentSelectionIsValid = selectedPropertyId.value && 
      (
        ['__TOTAL_TIMES__', '__TOTAL_USERS__'].includes(selectedPropertyId.value) || 
        eventProperties.value.some(p => p.property_id === selectedPropertyId.value)
      );

    if (!currentSelectionIsValid && !appliedFromInitial) {
      selectedPropertyId.value = '__TOTAL_TIMES__';
      selectedFormulaValue.value = MF.Count;
    } else {
      // 触发一次属性变更检查，以确保公式是正确的
      onPropertyChange(true); // silent = true
    }
  } catch (error) {
    console.error('加载事件属性失败:', error);
  }
}

// 根据父组件的已选值进行一次回显设置（在属性加载完成后调用）
function applyInitialSelectionFromParent(): boolean {
  const initialField = props.initialField;
  const initialFormula = props.initialFormula ?? null;
  if (!initialField) return false;

  // 处理特殊项映射
  if (initialField === '__TOTAL_TIMES__' || (initialField === 'e_openid' && initialFormula === MF.Count)) {
    selectedPropertyId.value = '__TOTAL_TIMES__';
    selectedFormulaValue.value = MF.Count;
    return true;
  }
  if (initialField === '__TOTAL_USERS__' || (initialField === 'e_openid' && initialFormula === MF.CountDistinctUserId)) {
    selectedPropertyId.value = '__TOTAL_USERS__';
    selectedFormulaValue.value = MF.CountDistinctUserId;
    return true;
  }

  // 普通属性：需要在事件属性列表中存在
  const target = eventProperties.value.find(p => p.property_id === initialField);
  if (target) {
    selectedPropertyId.value = target.property_id;
    // 确保公式可用，否则回退为第一个允许的公式
    const allowed = getFormulaOptions(target).map(o => o.value);
    if (initialFormula !== null && allowed.some(v => v === initialFormula)) {
      selectedFormulaValue.value = initialFormula;
    } else if (allowed.length > 0) {
      selectedFormulaValue.value = allowed[0];
    } else {
      selectedFormulaValue.value = null;
    }
    return true;
  }
  return false;
}

// 监听可见性变化
watch(() => props.visible, (visible) => {
  if (visible) {
    loadEventProperties();
  }
});

// 监听项目和事件变化
watch([() => props.projectAlias, () => props.eventId], () => {
  if (props.visible) {
    loadEventProperties();
  }
});

onMounted(() => {
  if (props.visible) {
    loadEventProperties();
  }
});
</script>

<style scoped>
.popover-content {
  max-height: 450px;
  overflow-y: auto;
}

.selector-sentence {
  display: flex;
  align-items: center;
  flex-wrap: nowrap;
  gap: 8px;
  margin-bottom: 20px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
  overflow-x: auto;
}

.text {
  font-size: 14px;
  color: #303133;
  font-weight: 500;
  white-space: nowrap;
}

.event-name {
  font-size: 14px;
  color: #409eff;
  font-weight: 600;
  background: #e1f5fe;
  padding: 4px 8px;
  border-radius: 4px;
  border: 1px solid #b3d8ff;
  white-space: nowrap;
}

.property-selector,
.formula-selector {
  min-width: 100px;
  width: 140px;
  flex-shrink: 0;
}

.property-selector :deep(.el-input__inner),
.formula-selector :deep(.el-input__inner) {
  font-size: 13px;
}

.action-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding-top: 16px;
  border-top: 1px solid #ebeef5;
}
</style>

<style>
.event-property-popover {
  padding: 16px !important;
  z-index: 2000 !important;
}

.property-selector .el-select-dropdown,
.formula-selector .el-select-dropdown {
  z-index: 3000 !important;
}

.property-selector,
.formula-selector {
  position: relative;
  z-index: 1;
}

.property-selector .el-input,
.formula-selector .el-input,
.property-selector .el-select,
.formula-selector .el-select {
  pointer-events: auto;
}

.el-select-dropdown .el-select-dropdown__item {
  pointer-events: auto;
}

.event-property-popover .el-select-dropdown {
  z-index: 3000 !important;
}

.event-property-popover .el-select-dropdown .el-select-dropdown__item {
  position: relative;
  z-index: 3001 !important;
}
</style>

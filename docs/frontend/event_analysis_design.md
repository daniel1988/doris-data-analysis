# 事件分析功能与页面布局详细设计文档

本文档详细总结了 [EventAnalysis.vue](file:///d:/golang/game_bi_platform/frontend/src/views/analysis/EventAnalysis.vue) 的核心功能、页面布局设计、状态管理及数据交互流程。

---

## 1. 页面布局设计 (Page Layout)

页面采用典型的 **左右分栏 + 响应式** 布局，通过 [EventAnalysisLayout.vue](file:///d:/golang/game_bi_platform/frontend/src/views/analysis/components/EventAnalysisLayout.vue) 实现。

### 1.1 左侧配置面板 (Left Pane)
- **宽度**：固定 520px (响应式下可缩减至 480px 或 100%)。
- **组成部分**：
    - **指标配置列表 (MetricsList)**：管理分析的核心指标。支持普通指标和自定义公式指标。
    - **全局筛选卡片 (GlobalFiltersCard)**：定义对所有指标生效的过滤条件，支持事件属性、用户属性、标签及用户分群。
    - **全局分组卡片 (GlobalGroupsCard)**：定义数据的聚合维度，支持按值、按区间、按日期、按标签及按分群分组。
    - **操作底栏 (ActionFooter)**：固定在底部，包含“新建”、“保存报表”和“开始分析”按钮。
- **交互特性**：独立滚动，确保配置项较多时仍能方便操作。

### 1.2 右侧结果面板 (Right Pane)
- **宽度**：弹性自适应。
- **组成部分**：
    - **时间控制栏 (TimeControlsWithDynamicPicker)**：固定在顶部，支持时间粒度选择（天、小时等）、静态/动态时间范围选择及对比时间段设置。
    - **结果展示区 (ResultView)**：独立滚动的区域，包含：
        - **可视化图表**：支持表格、折线图、柱状图、饼图切换。
        - **数据详情表格**：支持分页、多列排序、单元格格式化。
        - **查询状态**：显示 SQL 生成预览按钮和请求计时器。

---

## 2. 状态管理与核心逻辑 (State & Logic)

页面采用 Vue 3 Composition API，业务逻辑高度解耦到多个 Composables 中。

### 2.1 状态中心 ([useEventAnalysisForm.ts](file:///d:/golang/game_bi_platform/frontend/src/views/analysis/composables/useEventAnalysisForm.ts))
- **`form` (Reactive)**：核心表单对象，符合 `EventAnalysisV1Payload` 接口。
    - `event_metrics`: 数组，存储所有已配置的指标。
    - `time_grain`: 时间粒度配置。
    - `filter_groups`: 结构化的过滤条件。
    - `groups`: 分组维度。
- **视图状态**：`loading` (加载中)、`rows/columns` (结果数据)、`vizMode` (图表类型)、`total` (总行数)。

### 2.2 指标管理逻辑 ([useEventMetricsActions.ts](file:///d:/golang/game_bi_platform/frontend/src/views/analysis/composables/useEventMetricsActions.ts))
- **普通指标**：包含 `e_event_id` (事件ID)、`metric` (包含字段、公式、格式)。默认公式为“总次数”。
- **自定义公式指标**：类型为 `isCustom: true`，包含 `custom_metric` (API 结构) 和 `custom_web_metric` (前端 UI 编辑结构)。
- **交互**：支持指标的添加、删除、重命名联动（根据事件和公式自动生成名称）。

### 2.3 自定义公式深度解析 ([useCustomFormula.ts](file:///d:/golang/game_bi_platform/frontend/src/views/analysis/composables/useCustomFormula.ts))
- **可视化编辑**：通过 [CustomFormulaInput.vue](file:///d:/golang/game_bi_platform/frontend/src/components/common/CustomFormulaInput.vue) 实现类 IDE 的 Token 化编辑体验。
- **Token 类型**：
    - **事件指标 Token**：每个 Token 都是一个独立的配置单元，包含：
        - 事件选择（如“支付成功”）。
        - 指标计算（如“求和(金额)”）。
        - **指标级过滤**：支持为公式中的每个基础指标单独设置过滤条件（如“支付成功 且 渠道=AppStore”）。
    - **数值 Token**：支持直接输入常量数字。
    - **运算符 Token**：支持 `+`, `-`, `*`, `/` 及括号 `(`, `)`。
- **元素映射与同步**：
    - 将 UI 上的 Token 序列实时同步到 API 所需的 `custom_formula` 字符串（如 `(m0 + m1) / 100`）。
    - **自动索引**：为公式中的每个事件指标分配唯一的内部别名（如 `m0`, `m1`），确保后端能准确解析复杂表达式。
- **输出格式化**：支持对计算结果进行格式化处理，包括原始值、整数、两位小数及百分比展示。

---

## 3. 数据交互流程 (Data Flow)

### 3.1 请求构建阶段 ([usePayloadBuilders.ts](file:///d:/golang/game_bi_platform/frontend/src/views/analysis/composables/usePayloadBuilders.ts))
在点击“开始分析”后，系统会执行以下构建步骤：
1. **构建过滤组 (`buildFilterGroups`)**：
    - 合并全局筛选器。
    - 特殊处理“用户分群”和“标签”类型的过滤，分别放入 `user_group_filters` 和 `tag_filters`。
    - 处理动态时间 (如“过去7天”) 和静态时间范围。
2. **构建分组项 (`buildGroups`)**：
    - 根据分组类型（按值、区间、日期等）生成对应的 API 结构。
    - 确保 `tag_group` 和 `user_group` 的编码正确传递。

### 3.2 校验与发送 ([useAnalysisRunner.ts](file:///d:/golang/game_bi_platform/frontend/src/views/analysis/composables/useAnalysisRunner.ts))
- **规范化 (Normalization)**：将前端便捷表示（如 `__TOTAL_TIMES__`）映射为后端标准的 `e_openid` + `COUNT` 公式。
- **自动填充**：补全缺失的默认值（如时间粒度默认按天）。
- **异步执行**：调用 `eventAnalysisV1Api.execute`，并处理多层 Data Envelope 的解包，提取 `rows`、`columns`、`total` 和 `sql`。

---

## 4. 关键功能特性

- **双层过滤体系**：
    - **指标级过滤**：仅影响单个指标的计算逻辑。
    - **全局级过滤**：作为 WHERE 条件作用于整个查询。
- **属性联动机制**：
    - 切换事件时，系统通过 `useEventMeta.ts` 自动刷新可用的属性列表。
- **报表持久化**：
    - 通过 [useReportPersistence.ts](file:///d:/golang/game_bi_platform/frontend/src/views/analysis/composables/useReportPersistence.ts) 将整个 `form` 状态序列化并保存。
    - 支持通过 URL 中的 `report_id` 在页面挂载时自动还原复杂的分析场景。
- **开发者友好**：
    - 集成 SQL 预览功能，实时展示后端生成的 Doris/ClickHouse 查询语句，方便数据审计和排障。

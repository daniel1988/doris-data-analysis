# 漏斗分析 - 业务逻辑与 Composables (old_frontend)

本文档整理了 `old_frontend` 中漏斗分析模块的业务逻辑封装。

---

## 1. 组合式逻辑 (Composables)

漏斗分析的业务逻辑主要拆分为三个 Composables：

### 1.1 [useFunnelForm.ts](file:///d:/gitee/dmp_admin_v2/old_frontend/views/funnel/composables/useFunnelForm.ts)
-   **核心职能**：管理漏斗分析的所有表单状态及其变更逻辑。
-   **职责**：
    -   定义 `FunnelForm` 的初始状态，包括默认的两个步骤。
    -   提供 `addStep()` 和 `removeStep()` 函数，管理步骤序列的增删。
    -   提供 `validate()` 函数，在分析执行前校验：
        -   步骤数量（2-8 步）。
        -   各步骤是否已选择事件。
        -   转化窗口期是否合法（> 0）。
        -   是否选择了查询时间范围。

### 1.2 [useFunnelAnalysis.ts](file:///d:/gitee/dmp_admin_v2/old_frontend/views/funnel/composables/useFunnelAnalysis.ts)
-   **核心职能**：负责与后端漏斗分析接口进行通信。
-   **职责**：
    -   调用 `funnelAnalysisApi.query(form)` 执行查询。
    -   管理 `loading` 状态。
    -   **数据解析 (Unwrap)**：通过 `unwrapEnvelope` 函数处理后端多层嵌套的响应结构，提取 `rows` 和 `columns`。
    -   **结果标准化**：将后端返回的原始数据转换为前端组件可识别的 `result` 对象，包含数据行、列定义、总数及生成的 SQL。
    -   提供统一的错误处理机制（400 参数错误、500 服务器错误、504 超时等）。

### 1.3 [useFunnelReportPersistence.ts](file:///d:/gitee/dmp_admin_v2/old_frontend/views/funnel/composables/useFunnelReportPersistence.ts)
-   **核心职能**：管理漏斗报表的保存、加载与回填。
-   **职责**：
    -   管理 `currentReport` 状态。
    -   **报表保存**：将当前表单状态（步骤、过滤、分组、时间配置）打包为 `query_params` 并调用 `reportApi`。
    -   **报表加载**：从保存的 `query_params` 中解析并恢复 `FunnelForm` 的各项状态。
    -   管理保存弹窗的可见性。

---

## 2. 数据处理流

1.  **初始化**：`useFunnelForm` 初始化表单，默认包含两个空白步骤。
2.  **用户配置**：用户选择事件、添加过滤条件、设置窗口期和分组。
3.  **校验与提交**：点击“开始分析”触发 `validate()`，通过后调用 `runAnalysis(form)`。
4.  **响应处理**：`useFunnelAnalysis` 接收响应，进行解包映射，更新 `result` 状态。
5.  **渲染展示**：`FunnelResultView` 监听 `result` 变化，重新渲染漏斗图和数据表。
6.  **报表操作**：用户点击“保存报表”，由 `useFunnelReportPersistence` 处理配置持久化。

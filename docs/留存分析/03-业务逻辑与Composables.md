# 留存分析 - 业务逻辑与 Composables (old_frontend)

本文档整理了 `old_frontend` 中留存分析模块的业务逻辑封装。

---

## 1. 组合式逻辑 (Composables)

留存分析的主要业务逻辑被拆分为四个 Composables：

### 1.1 [useRetentionForm.ts](file:///d:/gitee/dmp_admin_v2/old_frontend/views/retention/composables/useRetentionForm.ts)
-   **核心职能**：管理留存分析的所有表单状态。
-   **职责**：
    -   定义 `RetentionForm` 的初始状态。
    -   提供 `validate()` 函数，在分析执行前校验必填项。
    -   管理 `day_n_array` 的默认值和约束条件。
    -   管理 `init_event_metric` 和 `end_event_metric`。

### 1.2 [useRetentionAnalysis.ts](file:///d:/gitee/dmp_admin_v2/old_frontend/views/retention/composables/useRetentionAnalysis.ts)
-   **核心职能**：负责与后端 API 进行通信。
-   **职责**：
    -   调用 `retentionAnalysisApi.query()`。
    -   处理请求过程中的 `loading` 状态。
    -   统一处理 API 错误信息。
    -   存储并提供分析结果 `result`。

### 1.3 [useRetentionReportPersistence.ts](file:///d:/gitee/dmp_admin_v2/old_frontend/views/retention/composables/useRetentionReportPersistence.ts)
-   **核心职能**：管理留存分析报表的保存和加载。
-   **职责**：
    -   处理报表保存弹窗的开启/关闭。
    -   调用 `reportApi` 保存分析配置。
    -   处理报表加载后的表单回填逻辑。

### 1.4 [useRetentionResult.ts](file:///d:/gitee/dmp_admin_v2/old_frontend/views/retention/composables/useRetentionResult.ts)
-   **核心职能**：对分析结果进行转换和后处理，用于 UI 渲染。
-   **职责**：
    -   将后端的二维数组或原始数据转换为图表所需的 `series`。
    -   将结果映射为留存矩阵表格。
    -   处理“留存率”与“留存人数”的单位转换。

---

## 2. 数据处理流

1.  **初始化**：从 `useRetentionForm` 获取默认表单配置。
2.  **配置**：用户在界面中配置初始事件、结束事件、分析天数、时间范围等。
3.  **校验**：用户点击“开始分析”，触发 `validate()` 校验。
4.  **请求**：校验通过，`useRetentionAnalysis` 携带表单数据调用后端。
5.  **展示**：API 返回数据，存储到 `result`，由 `useRetentionResult` 处理并渲染到 `RetentionResultView`。
6.  **持久化**：用户可随时开启保存报表流程。

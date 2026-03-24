# 用户属性分析 - 业务逻辑与 Composables (old_frontend)

本文档整理了 `old_frontend` 中用户属性分析模块的业务逻辑封装。

---

## 1. 组合式逻辑 (Composables)

属性分析的主要业务逻辑被拆分为三个 Composables：

### 1.1 [usePropertyAnalysisForm.ts](file:///d:/gitee/dmp_admin_v2/old_frontend/views/property-analysis/composables/usePropertyAnalysisForm.ts)
-   **核心职能**：管理属性分析的所有表单状态。
-   **职责**：
    -   定义 `UserPropertyForm` 的初始状态。
    -   管理 `user_property_metric` 的默认值（默认为“用户数”）。
    -   管理 `user_property_group_type` (1: 普通, 2: 人群分群)。
    -   提供 `validate()` 和 `resetForm()`。

### 1.2 [usePropertyAnalysisRunner.ts](file:///d:/gitee/dmp_admin_v2/old_frontend/views/property-analysis/composables/usePropertyAnalysisRunner.ts)
-   **核心职能**：负责负载构建与 API 执行。
-   **职责**：
    -   **负载转换**：将表单中的分组 (groups) 和人群分群 (user_groups) 按照当前模式转换成后端识别的统一结构。
    -   调用 `userPropertyAnalysisApi.query()`。
    -   管理结果集 (`rows`, `columns`, `total`)。
    -   提供生成的 SQL 字符串。

### 1.3 [usePropertyReportPersistence.ts](file:///d:/gitee/dmp_admin_v2/old_frontend/views/property-analysis/composables/usePropertyReportPersistence.ts)
-   **核心职能**：处理报表的保存与加载。
-   **职责**：
    -   打包当前分析配置为 JSON 字符串。
    -   处理报表加载后的表单回填（需处理分组模式的兼容性）。

---

## 2. 数据处理流

1.  **初始化**：`usePropertyAnalysisForm` 生成默认表单。
2.  **配置**：用户选择分析属性、设置聚合公式、添加过滤器。
3.  **模式选择**：用户选择按“维度”分组或按“人群分群”对比。
4.  **执行分析**：`usePropertyAnalysisRunner` 构建 Payload，如果是“总用户数”，转换占位符。
5.  **响应解析**：API 返回聚合结果，`ResultView` 负责渲染图表与表格。
6.  **持久化**：用户可将当前配置保存至报表系统。

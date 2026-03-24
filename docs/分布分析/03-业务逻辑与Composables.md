# 03-业务逻辑与Composables

分布分析的业务逻辑将通过多个 Composable 进行拆分管理，保持逻辑清晰且易于维护。

## 1. 逻辑拆分

### 1.1 useScatterForm.ts
- **职能**：管理分布分析的表单状态。
- **状态项**：
  - `project_alias`: 项目别名。
  - `metric`: 单个分析指标。
  - `scatter_type`: 分布类型（1-离散, 2-自动, 3-自定义）。
  - `scatter_ranges`: 自定义区间数组。
  - `time_grain`: 时间粒度。
  - `groups`: 分组维度。
- **方法**：`resetForm()`, `loadForm(data)`。

### 1.2 useScatterRunner.ts
- **职能**：执行分析请求并管理结果。
- **状态项**：`loading`, `results`, `columns`, `sql`, `total`。
- **方法**：`runAnalysis(payload)`。

### 1.3 useScatterReportPersistence.ts
- **职能**：处理分布分析报表的保存、加载、删除。
- **方法**：`saveReport(name)`, `fetchReportList()`, `deleteReport(id)`。

## 2. 数据转换逻辑

分布分析的结果集通常包含：
- 分组字段（可选）。
- `区间` (或 `metric_val`): 统计桶。
- `用户数`: 落入该桶的人数。

### 2.1 图表数据转换
- 需要将多行数据转换为 ECharts 柱状图所需的 `xAxis.data` 和 `series.data`。
- 如果有分组，需要生成多个 Series。

### 2.2 排序逻辑
- 默认按区间值从小到大排序。
- 自定义区间需要保持用户定义的顺序。

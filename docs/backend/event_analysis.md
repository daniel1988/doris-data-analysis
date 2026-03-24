# 事件分析 (Event Analysis)

事件分析用于分析特定事件在一段时间内的变化趋势，支持按维度分组和过滤。

## 1. 接口定义

- **函数**: `EventAnalysis(req *EventAnalysisReq) (interface{}, error)`
- **请求结构体**: `EventAnalysisReq`

### 1.1 请求参数 (EventAnalysisReq)

| 参数名 | 类型 | 描述 |
| :--- | :--- | :--- |
| `ProjectAlias` | `string` | 项目别名 |
| `EventMetrics` | `[]doris.EventMetric` | 待分析的事件及指标（如：次数、人数、人均次数等） |
| `FilterGroups` | `doris.GlobalFilterGroups` | 全局过滤条件 |
| `TimeGrain` | `doris.TimeGrain` | 时间粒度（天、小时等） |
| `Groups` | `[]doris.Group` | 分组维度 |
| `Orders` | `[]doris.Order` | 排序规则 |

## 2. 业务逻辑

1. **SQL 构建**: 通过 `BuildEventAnalysisSql` 函数动态构建 Doris SQL。
2. **多指标处理**: 支持同时查询多个事件的不同度量指标。
3. **维度分组**: 支持按多个属性进行分组。
4. **时间聚合**: 根据选定的时间粒度（如 `day`, `hour`）对数据进行聚合。
5. **执行查询**: 调用 `SqlService` 在 Doris 集群中执行生成的 SQL 并返回结果。

## 3. SQL 构建流程

- 使用 `WITH` 语句构建中间表（包含标签关联等）。
- 构建主事件表 `EVENT_TABLE_ALIAS`。
- 应用时间粒度公式和维度分组。
- 应用全局过滤条件。
- 处理排序和结果返回。

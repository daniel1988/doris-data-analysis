# 留存分析 (Retention Analysis)

留存分析用于分析执行了初始事件的用户，在后续一段时间内执行回访事件的情况。

## 1. 接口定义

- **函数**: `RetentionAnalysis(req *RetentionAnalysisReq) (*model.QueryResponse, error)`
- **请求结构体**: `RetentionAnalysisReq`

### 1.1 请求参数 (RetentionAnalysisReq)

| 参数名 | 类型 | 描述 |
| :--- | :--- | :--- |
| `ProjectAlias` | `string` | 项目别名 |
| `InitEventMetric` | `doris.EventMetric` | 初始事件（起始动作） |
| `EndEventMetric` | `doris.EventMetric` | 回访事件（后续动作） |
| `GlobalFilterGroups`| `doris.GlobalFilterGroups` | 全局过滤条件 |
| `Groups` | `[]doris.Group` | 分组维度 |
| `TimeGrain` | `doris.TimeGrain` | 时间粒度 |
| `DayNArray` | `[]int` | 分析的留存天数序列（如 [1, 3, 7]） |

## 2. 业务逻辑

1. **初始用户群**: 筛选在指定时间内执行了 `InitEventMetric` 的用户。
2. **回访校验**: 检查上述用户在后续的第 N 天是否执行了 `EndEventMetric`。
3. **矩阵构建**:
   - `buildInitRetentionMetric`: 计算每日/每单位时间的新增用户基数。
   - `buildEndRetentionMetric`: 计算后续回访的用户。
   - 通过 `LEFT JOIN` 关联初始表和回访表。
4. **留存率计算**: `ROUND(COUNT(DISTINCT 回访用户) / COUNT(DISTINCT 初始用户), 4)`。

## 3. 展示形式

- **留存矩阵**: 横轴为后续天数，纵轴为初始日期。
- **留存趋势**: 特定留存天数（如次留）随日期的变化曲线。
 bitumen

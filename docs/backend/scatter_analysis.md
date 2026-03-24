# 分布分析 (Scatter/Distribution Analysis)

分布分析用于研究特定指标在用户群中的分布情况（如：用户购买次数的分布、访问时长的分布）。

## 1. 接口定义

- **函数**: `ScatterAnalysis(req *ScatterAnalysisReq) (*model.QueryResponse, error)`
- **请求结构体**: `ScatterAnalysisReq`

### 1.1 请求参数 (ScatterAnalysisReq)

| 参数名 | 类型 | 描述 |
| :--- | :--- | :--- |
| `ProjectAlias` | `string` | 项目别名 |
| `ScatterMetric` | `doris.EventMetric` | 待分析的分布指标 |
| `FilterGroups` | `doris.GlobalFilterGroups` | 全局过滤条件 |
| `Groups` | `[]doris.Group` | 分组维度 |
| `TimeGrain` | `doris.TimeGrain` | 时间粒度 |
| `ScatterType` | `int` | 分布类型（1: 离散, 2: 默认区间, 3: 自定义区间, 4: 百分比） |
| `ScatterRanges` | `[]ScatterRange` | 自定义区间定义 |

## 2. 业务逻辑

1. **确定范围**: 如果是自动区间，首先查询指标的最大值。
2. **数据分桶**:
   - `Scatter_Type_Value`: 每一个值作为一个桶。
   - `Scatter_Type_Custom_Range`: 根据用户定义的 `Min/Max` 进行分桶。
   - 自动分桶: 将数据范围近似均分为 12 个区间。
3. **SQL 构建**:
   - 使用 `CASE WHEN` 语句在 SQL 层面对数据进行打标归类。
   - 统计每个桶内的去重用户数（`COUNT(DISTINCT e_openid)`）。
4. **结果聚合**: 按时间粒度和维度分组展示分布变化。

## 3. 分布类型

- **离散数字**: 适用于小范围整数（如活跃天数）。
- **自定义区间**: 业务定义的价值分层（如金额区间）。
- **自适应区间**: 系统自动根据数据现状生成的等分区间。
 bitumen

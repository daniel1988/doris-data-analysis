# 漏斗分析 (Funnel Analysis)

漏斗分析用于分析用户在预设步骤中的流失与转化情况。

## 1. 接口定义

- **函数**: `FunnelAnalysis(req *FunnelAnalysisReq) (*model.QueryResponse, error)`
- **请求结构体**: `FunnelAnalysisReq`

### 1.1 请求参数 (FunnelAnalysisReq)

| 参数名 | 类型 | 描述 |
| :--- | :--- | :--- |
| `ProjectAlias` | `string` | 项目别名 |
| `EventMetrics` | `[]doris.EventMetric` | 漏斗的步骤序列（有序） |
| `FilterGroups` | `doris.GlobalFilterGroups` | 全局过滤条件 |
| `Groups` | `[]doris.Group` | 分组维度 |
| `TimeGrain` | `doris.TimeGrain` | 时间粒度 |

## 2. 业务逻辑

1. **步骤定义**: 用户按顺序定义一系列事件作为漏斗的步骤。
2. **窗口计算**: 在指定的时间窗口内，计算用户从步骤 N 转化到步骤 N+1 的比例。
3. **SQL 构建**:
   - `buildFunnelBaseTable`: 构建初始用户表。
   - `buildFunnelEventTable`: 构建符合步骤序列的事件路径表。
   - 利用 Doris 的窗口函数或序列匹配能力计算转化。
4. **结果返回**: 返回各步骤的用户数及转化率。

## 3. 特色功能

- **分组分析**: 支持按维度（如城市、渠道）对比不同人群的漏斗转化。
- **灵活窗口**: 支持自定义漏斗完成的有效时间窗口。

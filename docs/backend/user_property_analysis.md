# 用户属性分析 (User Property Analysis)

用户属性分析用于根据用户的静态属性或画像标签进行人群画像分析和指标统计。

## 1. 接口定义

- **函数**: `UserPropertyAnalysis(req *UserPropertyAnalysisReq) (*model.QueryResponse, error)`
- **请求结构体**: `UserPropertyAnalysisReq`

### 1.1 请求参数 (UserPropertyAnalysisReq)

| 参数名 | 类型 | 描述 |
| :--- | :--- | :--- |
| `ProjectAlias` | `string` | 项目别名 |
| `UserPropertyMetric`| `UserPropertyMetric` | 待统计的属性指标 |
| `FilterGroups` | `doris.GlobalFilterGroups` | 全局过滤条件 |
| `UserPropertyGroupType` | `int` | 分组类型（1: 普通分组, 2: 属性分群） |
| `Groups` | `[]doris.Group` | 分组维度 |
| `UserGroups` | `[]doris.UserGroup` | 属性分群定义 |

## 2. 业务逻辑

1. **指标处理**: 支持统计用户总数（`__TOTAL_USERS__` 映射为 `u_openid`）或属性的聚合值（Sum, Avg 等）。
2. **分组逻辑**:
   - 支持常规的属性分组（如：按省份）。
   - 支持 **用户属性分群**: 类似于“活跃用户”、“流失用户”的逻辑分组。
3. **SQL 构建**:
   - `BuildUserPropertySql`: 构建针对用户表（`users`）的查询。
   - 关联标签表（如果涉及标签过滤）。
   - 应用分组聚合。
4. **查询执行**: 返回按维度划分的属性统计结果。

## 3. 应用场景

- 查看不同等级用户的平均客单价。
- 统计各渠道注册用户的性别分布。
- 分析特定画像人群的人数占比。

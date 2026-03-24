# 用户列表 (User List)

用户列表接口用于根据复杂的筛选条件（行为、属性、标签）查询并分页展示明细用户数据。

## 1. 接口定义

- **函数**: `UserList(req *UserListReq) (interface{}, error)`
- **请求结构体**: `UserListReq`

### 1.1 请求参数 (UserListReq)

| 参数名 | 类型 | 描述 |
| :--- | :--- | :--- |
| `ProjectAlias` | `string` | 项目别名 |
| `PageSize` | `int` | 每页数量 |
| `PageNum` | `int` | 当前页码 |
| `SelectFields` | `[]string` | 需查询展示的字段列表 |
| `UserFilterGroup` | `doris.FilterGroup` | 用户属性过滤条件 |
| `EventFilterGroup`| `doris.FilterGroup` | 事件行为过滤条件 |
| `TagFilters` | `[]doris.TagFilter` | 标签过滤条件 |
| `UserGroupFilters`| `[]UserGroupDataFilter` | 分群过滤条件 |

## 2. 业务逻辑

1. **默认字段**: 如果未指定 `SelectFields`，默认返回 `u_openid` 和 `u_event_time`。
2. **多维筛选**:
   - **属性过滤**: 直接作用于 `users` 表。
   - **标签/分群过滤**: 通过子查询（`IN` 语句）关联 `user_tag_data` 或 `user_group_data` 表。
3. **SQL 构建**:
   - 构建 `SELECT` 和 `FROM` 子句。
   - 组合所有 `WHERE` 条件（AND 连接）。
   - 应用 `ORDER BY u_event_time DESC`（默认按最后活跃时间排序）。
4. **分页处理**: 调用 `SqlService` 的分页执行逻辑，返回总数和当前页数据。

## 3. 核心功能

- **人群透放**: 配合分析模型，下钻查看符合特定行为特征的用户清单。
- **用户画像展示**: 导出或展示特定人群的详细属性信息。
 bitumen

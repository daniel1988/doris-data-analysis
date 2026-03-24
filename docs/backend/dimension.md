# 维度值获取 (Dimension Values)

该接口用于获取特定字段（事件属性、用户属性或标签）的所有可选维度值，通常用于前端筛选器的下拉列表。

## 1. 接口定义

- **函数**: `GetDimensions(req *DimensionReq) ([]map[string]interface{}, error)`
- **请求结构体**: `DimensionReq`

### 1.1 请求参数 (DimensionReq)

| 参数名 | 类型 | 描述 |
| :--- | :--- | :--- |
| `ProjectAlias` | `string` | 项目别名 |
| `Table` | `string` | 数据表名（events, users, tags） |
| `Field` | `string` | 目标字段名 |
| `EventId` | `string` | 事件 ID（仅当查询事件属性时需要） |

## 2. 业务逻辑

1. **类型分发**:
   - 如果 `Table` 是 `tags`，则调用 `TagValues` 处理。
   - 否则根据 `events` 或 `users` 表进行查询。
2. **连接获取**: 根据项目别名动态获取对应的 Doris 数据库连接。
3. **数据查询**:
   - 执行 `SELECT Field FROM Table GROUP BY Field` 获取去重后的值。
   - **自动限时**: 仅查询最近 7 天内出现过的值（性能优化）。
   - **限额返回**: 默认返回前 100 条数据。
4. **排序**: 按维度值升序排列。

## 3. 标签处理 (TagValues)

- 特殊处理标签表逻辑。
- 从标签明细表中根据 `tag_code` 提取 `tag_value`。
- 返回友好的别名映射。
 bitumen

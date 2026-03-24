---
name: "doris-utils"
description: "Doris 分析工具包，包含指标计算公式常量、过滤操作符、以及指标/过滤器/分组结构体定义。在处理后端 doris SQL 拼接、分析模块开发时调用。"
---

# doris-utils (Doris SQL 工具包)

该技能封装了后端 `pkg/doris` 目录下核心的数据结构与常量定义。在进行与 Doris 数据库交互、拼接分析 SQL 时（如事件分析、漏斗分析、留存分析等），请遵循以下定义的结构和常量。

## 1. 表与别名常量 (`const.go`)

```go
const (
	EVENT_TABLE = "event_data"
	USER_TABLE  = "user_data"
	TAG_TABLE   = "user_tag_data"
	USER_GROUP_TABLE = "user_group_data"

	EVENT_TABLE_ALIAS = "events"
	USER_TABLE_ALIAS  = "users"
)
```

## 2. 项目元数据表结构

除了上述常量，`dmp_center` 数据库中还包含以下核心的项目元数据管理表，在进行系统设置或配置读取时可直接引用：

- **项目数据 (`project_data`)**：记录项目的基础信息（别名 `project_alias`、名称、region等）。
- **元事件定义 (`project_event`)**：记录每个项目的事件字典（`event_id`，`event_name`）。
- **元属性定义 (`project_property`)**：记录每个项目的事件属性字典（`property_id`，`property_name`，数据类型等）。
- **事件-属性关联 (`project_event_property`)**：记录特定事件包含哪些属性。
- **用户属性定义 (`user_properties`)**：记录 `user_data` 表的字段字典。
- **用户标签/分群 (`user_tags`)**：记录人群圈选的标签定义与对应 SQL。
- **指标配置表 (`project_metrics`)**：记录预定义的指标公式表达式。

在涉及到后端读取元数据字典时，通常通过 `core.GetProjectCenter()` 结合上述表进行 GORM 查询。

## 3. 过滤操作符 (Operators)

```go
const (
	OperEmpty          = iota // 0
	OperEqualTo               // 1: =
	OperNotEqualTo            // 2: !=
	OperLessThan              // 3: <
	OperLessOrEqual           // 4: <=
	OperGreaterThan           // 5: >
	OperGreaterOrEqual        // 6: >=
	OperIsNull                // 7: IS NULL
	OperIsNotNull             // 8: IS NOT NULL
	OperBetween               // 9: BETWEEN
	OperIn                    // 10: IN
	OperNotIn                 // 11: NOT IN
	OperLike                  // 12: LIKE
	OperDateDiff              // 13: 日期差
	OperNotLike               // 14: NOT LIKE
	OperBelongTo              // 15: 属于
	OperNotBelongTo           // 16: 不属于
	OperStartWith             // 17: LIKE 'xxx%'
	OperEndWith               // 18: LIKE '%xxx'
	OperDynamicDates          // 19: 动态日期
	OperNDayRegiste           // 20: N 天前注册
)
```

## 3. 指标计算公式常量 (Formulas)

```go
const (
	Formula_Empty                       = iota // 0
	Formula_Sum                                // 1: SUM
	Formula_Count                              // 2: COUNT
	Formula_Count_Distinct                     // 3: COUNT(DISTINCT)
	Formula_Count_Distinct_UserId              // 4: COUNT(DISTINCT user_id)
	Formula_Any                                // 5: ANY_VALUE
	Formula_Max                                // 6: MAX
	Formula_Min                                // 7: MIN
	Formula_Avg                                // 8: AVG
	// ... 更多参考 pkg/doris/const.go
	Formula_Count_Days                  = 15   // 15: COUNT(DISTINCT DATE())
	Formula_Bitmap_Union                = 18   // 18: BITMAP_UNION(BITMAP_HASH())
	Formula_Count_Distinct_Daily_UserId = 19   // 19: 按天去重计数
)
```

## 4. 核心结构体定义

### 4.1 Filter (过滤器)
```go
type Filter struct {
	Column   Column      `json:"column"`
	Value    interface{} `json:"value"`
	Operator int         `json:"operator"` // 对应上方 OperXXX 常量
}
```

### 4.2 Metric (指标)
```go
type Metric struct {
	Column  Column `json:"column"`
	Formula int    `json:"formula"` // 对应上方 Formula_XXX 常量
	Format  string `json:"format"`  // "raw", "int", "decimal", "percent"
}
```

### 4.3 Group (分组)
```go
type GroupType int
const (
	GroupBy_Value         = 0 // 按值
	GroupBy_Ranges        = 1 // 按照区间
	GroupBy_Date          = 2 // 按照日期
	GroupBy_TagGroup      = 3 // 按照标签分组
	GroupBy_UserGroupData = 4 // 按照用户分群
	GroupBy_UserGroups    = 5 // 用户属性分群
)

type Group struct {
	GroupType   GroupType           `json:"group_type"`
	Column      Column              `json:"column"`
	ValueRanges []ValueRange        `json:"value_ranges"`
	TimeGrain   TimeGrain           `json:"time_grain"`
	TagGroup    TagFilter           `json:"tag_group"`
	UserGroup   UserGroupDataFilter `json:"user_group"`
}
```

## 使用说明
当需要开发后端分析模块构建复杂的 Doris SQL 时，你可以直接参考并组合使用该技能中定义的这些常量和结构，来构建 `BuildFilter`、`BuildMetric` 和 `BuildGroup` 的调用。
package doris

const (
	EVENT_TABLE = "event_data"
	USER_TABLE  = "user_data"
	TAG_TABLE   = "user_tag_data"

	USER_GROUP_TABLE = "user_group_data"

	EVENT_TABLE_ALIAS = "events"
	USER_TABLE_ALIAS  = "users"

	USER_PROPERTY_TABLE       = "ups"
	USER_PROPERTY_TABLE_ALIAS = "ups"
)

// 过滤操作符
const (
	OperEmpty = iota
	OperEqualTo
	OperNotEqualTo
	OperLessThan
	OperLessOrEqual
	OperGreaterThan
	OperGreaterOrEqual
	OperIsNull
	OperIsNotNull
	OperBetween
	OperIn
	OperNotIn
	OperLike
	OperDateDiff
	OperNotLike
	OperBelongTo
	OperNotBelongTo
	OperStartWith
	OperEndWith
	OperDynamicDates // 动态日期
	OperNDayRegiste  // N 天前注册
)

const (
	Formula_Empty = iota
	Formula_Sum
	Formula_Count
	Formula_Count_Distinct
	Formula_Count_Distinct_UserId
	Formula_Any
	Formula_Max
	Formula_Min
	Formula_Avg
	Formula_IsNull
	Formula_IsNotNull
	Formula_IsTrue
	Formula_IsFalse
	Formula_First
	Formula_Last
	Formula_Count_Days
	Formula_Last_Value  // 指定列最大值： 即按照事件时间倒序排序第一条
	Formula_First_Value // 指定列最小值： 即按照事件时间正序排序第一条
	Formula_Bitmap_Union
	Formula_Count_Distinct_Daily_UserId // 按天去重计数
)

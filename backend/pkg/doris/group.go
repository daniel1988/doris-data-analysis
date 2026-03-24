package doris

import (
	"fmt"
	"strings"
)

const (
	Display_Value = iota
	Display_Row
	Display_Date
	Display_Column
)

type GroupType int

const (
	GroupBy_Value         = 0 // 按值
	GroupBy_Ranges        = 1 // 按照区间
	GroupBy_Date          = 2 // 按照日期
	GroupBy_TagGroup      = 3 //按照标签分组
	GroupBy_UserGroupData = 4 //按照用户分群
	GroupBy_UserGroups    = 5 // 用户属性分群
)

type Group struct {
	GroupType GroupType `json:"group_type"`
	Column    Column    `json:"column"`

	ValueRanges []ValueRange        `json:"value_ranges"`
	TimeGrain   TimeGrain           `json:"time_grain"`
	TagGroup    TagFilter           `json:"tag_group"`
	UserGroup   UserGroupDataFilter `json:"user_group"`
}
type UserGroup struct {
	Alias       string      `json:"alias"`
	FilterGroup FilterGroup `json:"filter_group"`
}

type ValueRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

func BuildGroups(groups []Group) ([]string, []string) {
	groupFields := make([]string, 0)
	selectFields := make([]string, 0)
	for _, group := range groups {
		switch group.GroupType {
		case GroupBy_Value:
			groupFields = append(groupFields, group.Column.Field)
			selectFields = append(selectFields, group.Column.Field)
		case GroupBy_Ranges:
			caseSql, alias := GroupByRanges(group)
			groupFields = append(groupFields, alias)
			selectFields = append(selectFields, caseSql)
		case GroupBy_TagGroup:
			selectField, alias := BuildTagGroup(group)
			groupFields = append(groupFields, alias)
			selectFields = append(selectFields, selectField)
		case GroupBy_UserGroupData:
			selectField, alias := BuildUserGroupData(group)
			groupFields = append(groupFields, alias)
			selectFields = append(selectFields, selectField)
		}
	}
	return groupFields, selectFields
}

// 构建分组字段
func BuildGroup(group Group) (string, string) {
	selectField := group.Column.Field

	alias := selectField
	switch group.GroupType {
	case GroupBy_Value:
		alias = selectField
	case GroupBy_Date:
		if group.Column.Alias == "" {
			group.Column.Alias = group.Column.Field
		}
		tg := TimeGrain{
			Column: Column{
				Field: group.Column.Field,
				Alias: group.Column.Alias,
				Table: group.Column.Table,
			},
			Interval: group.TimeGrain.Interval,
		}
		selectField, alias = BuildTimeGrainV2(tg)
	case GroupBy_Ranges:
		selectField, alias = GroupByRanges(group)
	case GroupBy_TagGroup:
		selectField, alias = BuildTagGroup(group)
	case GroupBy_UserGroupData:
		selectField, alias = BuildUserGroupData(group)
	}

	return selectField, alias
}

func GroupByRanges(group Group) (string, string) {
	whenConds := make([]string, 0)
	for _, vr := range group.ValueRanges {
		vrGroupalias := getValueRangeAlias(vr)
		if vr.Min <= 0 {
			whenConds = append(whenConds, fmt.Sprintf(" WHEN %s >= %v AND %s < %v THEN '%s'", group.Column.Field, vr.Min, group.Column.Field, vr.Max, vrGroupalias))
		} else {
			whenConds = append(whenConds, fmt.Sprintf(" WHEN %s >= %v AND %s < %v THEN '%s'", group.Column.Field, vr.Min, group.Column.Field, vr.Max, vrGroupalias))
		}
	}
	alias := fmt.Sprintf("%s区间", group.Column.Field)
	caseSql := fmt.Sprintf(`CASE %s	END AS "%s"`, strings.Join(whenConds, "\n"), alias)
	return caseSql, alias
}

const (
	Infinity       = "Infinity"
	Infinity_Value = 1000000
)

func getValueRangeAlias(vr ValueRange) string {
	min := fmt.Sprintf("%v", vr.Min)
	if vr.Min == -Infinity_Value {
		min = fmt.Sprintf("-%v", Infinity)
	}

	max := fmt.Sprintf("%v", vr.Max)
	if vr.Max == Infinity_Value {
		max = fmt.Sprintf("+%v", Infinity)
	}
	return fmt.Sprintf("[%v,%v)", min, max)
}

func BuildTagGroup(group Group) (string, string) {
	tagCode := group.TagGroup.TagCode
	selectField := fmt.Sprintf("`%s`.tag_value ", tagCode)
	return selectField, tagCode
}

func BuildUserGroupData(group Group) (string, string) {
	alias := group.UserGroup.GroupName
	return fmt.Sprintf("`%s`.user_id AS `%s`", group.UserGroup.GroupCode, alias), alias
}
func BuildUserGroups(ugs []UserGroup) (string, string) {
	whenSqls := make([]string, 0)

	for _, ug := range ugs {
		filterGroupSql := BuildFilterGroup(ug.FilterGroup)
		whenSqls = append(whenSqls, fmt.Sprintf(" WHEN %s THEN '%s'", filterGroupSql, ug.Alias))
	}

	caseSql := fmt.Sprintf(`CASE %s END `, strings.Join(whenSqls, "\n"))
	return caseSql, "指标"
}

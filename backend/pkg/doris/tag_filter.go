package doris

import (
	"fmt"
	"strings"
	"time"
)

type SetOperate int

const (
	SetOperateEmpty     = iota // 空
	SetOperateIntersect = 1    // 交集
	SetOperateUnion     = 2    // 并集
)

type TagFilter struct {
	TagCode  string      `json:"tag_code"`
	Operator int         `json:"operator"`
	TagValue interface{} `json:"tag_value"`
}

func BuildTagFilter(tf TagFilter) string {
	return fmt.Sprintf("`%s`.user_id IS NOT NULL", tf.TagCode)
}

func ExtractTagFilters(ems []EventMetric, fgs GlobalFilterGroups, groups []Group) []TagFilter {
	tagFilterMap := make(map[string]TagFilter)
	for _, em := range ems {
		if em.Type == EventType_Custom {
			for _, cem := range em.CustomMetric.EventMetrics {
				for _, tf := range cem.FilterGroup.TagFilters {
					tagFilterMap[tf.TagCode] = tf
				}
			}
		} else {
			for _, tf := range em.FilterGroup.TagFilters {
				tagFilterMap[tf.TagCode] = tf
			}
		}
	}
	for _, fg := range fgs.GlobalFilters.TagFilters {
		tagFilterMap[fg.TagCode] = fg
	}
	for _, g := range groups {
		if g.GroupType == GroupBy_TagGroup {
			if _, ok := tagFilterMap[g.TagGroup.TagCode]; ok {
				continue
			}
			tagFilterMap[g.TagGroup.TagCode] = g.TagGroup
		}
	}

	tagFilters := make([]TagFilter, 0)
	for _, tf := range tagFilterMap {
		tagFilters = append(tagFilters, tf)
	}
	return tagFilters
}

func BuildTagTables(tagFilters []TagFilter) []string {
	tagTables := make([]string, 0)
	for _, tf := range tagFilters {

		tableSqls := make([]string, 0)
		tableSqls = append(tableSqls, "SELECT user_id, tag_code, tag_value FROM user_tag_data")
		whereSqls := make([]string, 0)
		if tf.TagCode != "" {
			whereSqls = append(whereSqls, buildTagWhereSql(tf))
		}
		if len(whereSqls) > 0 {
			tableSqls = append(tableSqls, fmt.Sprintf("WHERE %s", strings.Join(whereSqls, " AND ")))
		}
		tagTables = append(tagTables, fmt.Sprintf("`%s` AS (%s)", tf.TagCode, strings.Join(tableSqls, " ")))
	}
	return tagTables
}

func buildTagWhereSql(tf TagFilter) string {
	switch tf.Operator {
	case OperEqualTo:
		return fmt.Sprintf("tag_code='%s' AND tag_value='%v'", tf.TagCode, tf.TagValue)
	case OperIn:
		return fmt.Sprintf("tag_code='%s' AND tag_value IN (%v)", tf.TagCode, formatTagValue(tf.TagValue))
	case OperNotEqualTo, OperNotIn:
		return fmt.Sprintf("tag_code='%s' AND tag_value NOT IN (%v)", tf.TagCode, formatTagValue(tf.TagValue))
	case OperGreaterOrEqual:
		return fmt.Sprintf("tag_code='%s' AND tag_value >= '%v'", tf.TagCode, tf.TagValue)
	case OperGreaterThan:
		return fmt.Sprintf("tag_code='%s' AND tag_value > '%v'", tf.TagCode, tf.TagValue)
	case OperLessOrEqual:
		return fmt.Sprintf("tag_code='%s' AND tag_value <= '%v'", tf.TagCode, tf.TagValue)
	case OperLessThan:
		return fmt.Sprintf("tag_code='%s' AND tag_value < '%v'", tf.TagCode, tf.TagValue)
	case OperIsNotNull:
		return fmt.Sprintf("tag_code='%s' AND (tag_value !='' OR tag_value IS NOT NULL)", tf.TagCode)
	case OperIsNull:
		return fmt.Sprintf("tag_code='%s' AND (tag_value ='' OR tag_value IS NULL)", tf.TagCode)
	}
	return fmt.Sprintf("tag_code='%s'", tf.TagCode)
}

func formatTagValue(value interface{}) string {
	switch value.(type) {
	case []string:
		arrayValues := make([]string, 0)
		for _, v := range value.([]string) {
			arrayValues = append(arrayValues, fmt.Sprintf("\"%v\"", v))
		}
		return strings.Join(arrayValues, ",")
	case []interface{}:
		arrayValues := make([]string, 0)
		for _, v := range value.([]interface{}) {
			arrayValues = append(arrayValues, fmt.Sprintf("\"%v\"", v))
		}
		return strings.Join(arrayValues, ",")
	case string:
		return fmt.Sprintf("\"%v\"", value)
	case int, int64, int32, int16, int8:
		return fmt.Sprintf("%v", value)
	case float64, float32:
		return fmt.Sprintf("%v", value)
	case bool:
		return fmt.Sprintf("%v", value)
	case time.Time:
		return fmt.Sprintf("\"%v\"", value)
	case time.Duration:
		return fmt.Sprintf("%v", value)
	case nil:
		return "''"
	default:
		return fmt.Sprintf("\"%v\"", value)
	}
}

func BuildJoinTagTables(baseTable string, tagFilters []TagFilter) string {
	joinSqls := make([]string, 0)
	joinField := "e_openid"
	if baseTable == USER_TABLE || baseTable == USER_TABLE_ALIAS {
		joinField = "u_openid"
	}
	for _, tf := range tagFilters {
		joinSqls = append(joinSqls, fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.user_id",
			tf.TagCode, baseTable, joinField, tf.TagCode))
	}
	return strings.Join(joinSqls, " \n")
}

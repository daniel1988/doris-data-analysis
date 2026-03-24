package doris

import (
	"fmt"
	"strconv"
	"strings"
)

type CustomMetric struct {
	EventMetrics  []EventMetric `json:"event_metrics"`
	CustomFormula string        `json:"custom_formula"`
	Format        string        `json:"format"`
}

const (
	EventType_Normal = 1 // 普通事件
	EventType_Custom = 2 // 自定义事件
)

type EventMetric struct {
	EventId      string       `json:"e_event_id"`
	Name         string       `json:"name"`
	Type         int          `json:"type"`
	Metric       Metric       `json:"metric"`
	CustomMetric CustomMetric `json:"custom_metric"`
	FilterGroup
}

func BuildEventMetric(em EventMetric, args ...any) string {
	if em.Metric.Column.Field == "__TOTAL_TIMES__" || em.Metric.Column.Field == "__TOTAL_USERS__" || em.Metric.Column.Field == "" {
		em.Metric.Column.Field = "e_openid"
	}
	if em.Type == 0 {
		em.Type = EventType_Normal
	}

	// 自定义事件处理
	if em.Type == EventType_Custom || len(em.CustomMetric.EventMetrics) > 0 {
		customFormula := em.CustomMetric.CustomFormula
		for cemIdx, cem := range em.CustomMetric.EventMetrics {
			if cem.Name == "" {
				continue
			}
			if cem.Metric.Format == "" {
				cem.Metric.Format = FormatDecimal
			}
			metricFormula := BuildEventMetric(cem, args...)

			if cem.Name == "" {
				cem.Name = fmt.Sprintf("%s_%v", em.Name, cemIdx)
			}

			customFormula = strings.ReplaceAll(customFormula, cem.Name, metricFormula)
		}

		return RoundFormat(customFormula, FormatDecimal)
	}
	if em.Metric.Formula == Formula_Last_Value || em.Metric.Formula == Formula_First_Value {
		em.Metric.Formula = Formula_Any
		metricName := em.Name
		if strings.Contains(metricName, "`") {
			metricName = strings.ReplaceAll(metricName, "`", "")
		}
		em.Metric.Column.Field = fmt.Sprintf("`%s`", metricName)
	}

	fg := FilterGroup{
		Scope:      FilterScope(em.Scope),
		Filters:    em.Filters,
		TagFilters: em.TagFilters,
	}

	caseFilterSqls := make([]string, 0)
	eventIdFilterSql := BuildVirtualEventFilter(em)
	caseFilterSqls = append(caseFilterSqls, eventIdFilterSql)
	if fgSql := BuildFilterGroup(fg); fgSql != "true" {
		caseFilterSqls = append(caseFilterSqls, fgSql)
	}

	tableAlias := EVENT_TABLE_ALIAS
	switch em.Metric.Column.Table {
	case USER_TABLE:
		tableAlias = USER_TABLE_ALIAS
	case EVENT_TABLE:
		tableAlias = EVENT_TABLE_ALIAS
	}
	caseSql := fmt.Sprintf("CASE WHEN %s THEN %s.%s END",
		strings.Join(caseFilterSqls, " AND "),
		tableAlias,
		em.Metric.Column.Field)

	return BuildMetric(Metric{
		Column: Column{
			Field: caseSql,
		},
		Formula: em.Metric.Formula,
		Format:  em.Metric.Format,
	})
}

func HasUserDataField(ems []EventMetric, fgs GlobalFilterGroups, groups []Group) bool {
	for _, column := range ExtractColumns(ems, fgs, groups) {
		if column.Table == USER_TABLE {
			return true
		}
	}
	return false
}

// 处理虚拟事件ID
func BuildVirtualEventFilter(em EventMetric) string {
	eventIdFilterSqls := make([]string, 0)
	if em.Type == EventType_Custom {
		for _, cem := range em.CustomMetric.EventMetrics {
			eventIdFilterSqls = append(eventIdFilterSqls, fmt.Sprintf("(%s)", BuildVirtualEventFilter(cem)))
		}
	}

	switch em.Type {
	case EventType_Normal:
		eventIdFilterSqls = append(eventIdFilterSqls, fmt.Sprintf("e_event_id = '%s'", em.EventId))
	}
	return fmt.Sprintf("(%s)", strings.Join(eventIdFilterSqls, " OR "))
}

func BuildWindowMetric(em EventMetric, tg TimeGrain, groups []Group) (string, string) {
	tgField := tg.Column.Field

	if !strings.Contains(tg.Column.Field, ".") {
		tgField = fmt.Sprintf("%s.%s", EVENT_TABLE_ALIAS, tg.Column.Field)
	}
	tgFormula := fmt.Sprintf(`DATE(%s)`, tgField)
	switch tg.Interval {
	case Tg_Interval_Day:
		tgFormula = fmt.Sprintf(`DATE(%s)`, tgField)
	case Tg_Interval_Week:
		tgFormula = fmt.Sprintf(`WEEK(%s)`, tgField)
	case Tg_Interval_Month:
		tgFormula = fmt.Sprintf(`MONTH(%s)`, tgField)
	case Tg_Interval_Quarter:
		tgFormula = fmt.Sprintf(`QUARTER(%s)`, tgField)
	case Tg_Interval_Year:
		tgFormula = fmt.Sprintf(`YEAR(%s)`, tgField)
	case Tg_Interval_Hour:
		tgFormula = fmt.Sprintf(`HOUR(%s)`, tgField)
	case Tg_Interval_Minute:
		tgFormula = fmt.Sprintf(`MINUTE(%s)`, tgField)
	}

	whenSqls := make([]string, 0)
	whenSqls = append(whenSqls, BuildVirtualEventFilter(em))

	if fgSql := BuildFilterGroup(FilterGroup{
		Scope:   FilterScope(em.Scope),
		Filters: em.Filters,
	}); fgSql != "true" {
		whenSqls = append(whenSqls, fgSql)
	}

	eventTimeField := "events.event_time"
	caseSql := fmt.Sprintf("CASE WHEN %s THEN %s.%s END",
		strings.Join(whenSqls, " AND "),
		EVENT_TABLE_ALIAS,
		em.Metric.Column.Field)
	if em.Metric.Column.Table == USER_TABLE {
		caseSql = em.Metric.Column.Field
		eventTimeField = "u_event_time"
	}
	partitionBys := make([]string, 0)
	partitionBys = append(partitionBys, tgFormula)
	partitionBys = append(partitionBys, "events.event_id")
	for _, group := range groups {
		switch group.Column.Table {
		case USER_TABLE, USER_PROPERTY_TABLE:
			partitionBys = append(partitionBys, USER_TABLE_ALIAS+"."+group.Column.Field)
		case EVENT_TABLE:
			partitionBys = append(partitionBys, EVENT_TABLE_ALIAS+"."+group.Column.Field)
		}
	}

	selectFormula := fmt.Sprintf(`FIRST_VALUE(%s) OVER (PARTITION BY %s ORDER BY %s DESC)`, caseSql, strings.Join(partitionBys, ","), eventTimeField)
	switch em.Metric.Formula {
	case Formula_First_Value:
		selectFormula = fmt.Sprintf(`FIRST_VALUE(%s) OVER (PARTITION BY %s ORDER BY %s ASC)`, caseSql, strings.Join(partitionBys, ","), eventTimeField)
	}
	return selectFormula, em.Name
}

func HasWindowMetric(em EventMetric) bool {
	if em.Metric.Formula == Formula_First_Value || em.Metric.Formula == Formula_Last_Value {
		return true
	}
	if em.Type == EventType_Custom {
		for _, cem := range em.CustomMetric.EventMetrics {
			if cem.Metric.Formula == Formula_First_Value || cem.Metric.Formula == Formula_Last_Value {
				return true
			}
		}
	}
	return false
}

func BuildEventMetricFilterSql(ems []EventMetric) string {
	whereSqls := make([]string, 0)
	for _, em := range ems {
		metricFilterSqls := make([]string, 0)
		if em.Type == EventType_Custom {
			whereSqls = append(whereSqls, BuildEventMetricFilterSql(em.CustomMetric.EventMetrics))
		} else {
			if em.EventId == "" {
				continue
			}
			metricFilterSqls = append(metricFilterSqls, fmt.Sprintf("event_id = '%s'", em.EventId))
			if fgSql := BuildFilterGroup(FilterGroup{
				Scope:   FilterScope(em.Scope),
				Filters: em.Filters,
			}); fgSql != "true" {
				metricFilterSqls = append(metricFilterSqls, fmt.Sprintf("(%s)", fgSql))
			}
		}
		if len(metricFilterSqls) > 0 {
			whereSqls = append(whereSqls, fmt.Sprintf("(%s)", strings.Join(metricFilterSqls, " AND ")))
		}
	}
	return fmt.Sprintf("(%s)", strings.Join(whereSqls, " OR "))
}

func ConvertMetricName(name string) string {
	specialChars := []string{",", ".", " ", "(", ")", "-", "_", ":", ";", "'", "\"", "!", "@", "#", "$", "%", "^", "&", "*", "+", "=", "?", "<", ">", "/", "\\", "|", "[", "]", "{", "}", "`", "~"}
	for _, char := range specialChars {
		if strings.Contains(name, char) {
			return fmt.Sprintf("`%s`", name)
		}
	}
	_, err := strconv.Atoi(name)
	if err == nil {
		return fmt.Sprintf("`%s`", name)
	}
	return fmt.Sprintf("`%s`", name)
}

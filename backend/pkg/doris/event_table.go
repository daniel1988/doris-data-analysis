package doris

import (
	"fmt"
	"strings"
	"time"
)

// 构建日志明细表 - 处理虚拟属性
func BuildEventDataTable(ems []EventMetric, tg TimeGrain, fgs GlobalFilterGroups, groups []Group, tz string) string {
	tableFields := make([]string, 0)
	uniqueFields := make(map[string]bool)
	if _, ok := uniqueFields["e_event_time"]; !ok {
		tableFields = append(tableFields, "e_event_time")
		uniqueFields["e_event_time"] = true
	}
	if _, ok := uniqueFields["e_event_id"]; !ok {
		tableFields = append(tableFields, "e_event_id")
		uniqueFields["e_event_id"] = true
	}
	if _, ok := uniqueFields["e_openid"]; !ok {
		tableFields = append(tableFields, "e_openid")
		uniqueFields["e_openid"] = true
	}

	for _, em := range ems {
		if em.Metric.Formula == Formula_First_Value || em.Metric.Formula == Formula_Last_Value {
			selectFormula, metricName := BuildWindowMetric(em, tg, groups)
			if strings.Contains(metricName, "`") {
				metricName = strings.ReplaceAll(metricName, "`", "")
			}
			tableFields = append(tableFields, fmt.Sprintf("%s AS `%s`", selectFormula, metricName))
		}
	}

	for _, column := range ExtractColumns(ems, fgs, groups) {
		if column.Table == EVENT_TABLE || column.Table == "" {
			if column.Field == "" {
				continue
			}
			if _, ok := uniqueFields[column.Field]; ok {
				continue
			}
			uniqueFields[column.Field] = true
			tableFields = append(tableFields, fmt.Sprintf("%s.%s", EVENT_TABLE_ALIAS, column.Field))
		}
	}

	eventTableSqls := make([]string, 0)
	eventTableSqls = append(eventTableSqls, fmt.Sprintf("SELECT %s", strings.Join(tableFields, ", ")))
	eventTableSqls = append(eventTableSqls, fmt.Sprintf("FROM %s %s", EVENT_TABLE, EVENT_TABLE_ALIAS))
	whereSqls := make([]string, 0)
	eventIdFilterSqls := make([]string, 0)
	normalEventIdValues := make([]string, 0)
	eventIdMap := make(map[string]bool)
	for _, em := range ems {
		if em.Type == EventType_Custom {
			for _, cem := range em.CustomMetric.EventMetrics {
				if cem.EventId == "" {
					continue
				}
				if _, ok := eventIdMap[cem.EventId]; ok {
					continue
				}
				eventIdMap[cem.EventId] = true
				normalEventIdValues = append(normalEventIdValues, fmt.Sprintf("'%s'", cem.EventId))
			}
		} else if em.EventId != "" {
			if _, ok := eventIdMap[em.EventId]; ok {
				continue
			}
			eventIdMap[em.EventId] = true
			normalEventIdValues = append(normalEventIdValues, fmt.Sprintf("'%s'", em.EventId))
		}
	}
	if len(normalEventIdValues) > 0 {
		eventIdFilterSqls = append(eventIdFilterSqls, fmt.Sprintf("events.e_event_id IN(%s)", strings.Join(normalEventIdValues, ",")))
	}
	if len(eventIdFilterSqls) > 0 {
		whereSqls = append(whereSqls, fmt.Sprintf("(%s)", strings.Join(eventIdFilterSqls, " OR ")))
	}

	if len(fgs.QueryDates) == 2 {
		tzField := fmt.Sprintf("%s.e_event_time", EVENT_TABLE_ALIAS)
		queryDates := ConvertTzQueryDates(fgs.QueryDates, tz)
		if strings.Contains(queryDates[0], "convert_tz") {
			whereSqls = append(whereSqls, fmt.Sprintf("(%s >= %s AND %s < %s)",
				tzField, queryDates[0], tzField, queryDates[1]))
		} else {
			whereSqls = append(whereSqls, fmt.Sprintf("(%s >= '%s' AND %s < '%s')",
				tzField, queryDates[0], tzField, queryDates[1]))
		}
	}
	eventTableSqls = append(eventTableSqls, fmt.Sprintf("WHERE %s", strings.Join(whereSqls, " AND ")))
	return strings.Join(eventTableSqls, " \n")
}

func ConvertTzQueryDates(queryDates []string, tz string) []string {
	return queryDates
	layout := "2006-01-02 15:04:05.999"
	convertQueryDates := make([]string, len(queryDates))
	for i, date := range queryDates {
		t, err := time.Parse(layout, date)
		if err == nil {
			t = AddDateForTzTime(t, tz)
			date = t.Format(layout)
			convertQueryDates[i] = fmt.Sprintf("convert_tz('%s', '+08:00', '%s')", date, tz)
		} else {
			convertQueryDates[i] = fmt.Sprintf("convert_tz('%s', '+08:00', '%s')", date, tz)
		}

	}
	return convertQueryDates
}

func BuildExtraWithTables(ems []EventMetric, fgs GlobalFilterGroups, groups []Group) []string {
	withTables := make([]string, 0)

	if userGroupFilters := ExtractUserGroupFilters(ems, fgs, groups); len(userGroupFilters) > 0 {
		withTables = append(withTables, BuildUserGroupFilterTables(userGroupFilters)...)
	}

	if tagFilters := ExtractTagFilters(ems, fgs, groups); len(tagFilters) > 0 {
		withTables = append(withTables, BuildTagTables(tagFilters)...)
	}
	return withTables
}

func BuildExtraJoinTables(baseTable string, ems []EventMetric, fgs GlobalFilterGroups, groups []Group) []string {
	joinTables := make([]string, 0)

	if HasUserDataField(ems, fgs, groups) {
		joinTables = append(joinTables, fmt.Sprintf("LEFT JOIN %s %s ON %s.e_openid = %s.u_openid",
			USER_TABLE, USER_TABLE_ALIAS, baseTable, USER_TABLE_ALIAS))
	}

	if tagFilters := ExtractTagFilters(ems, fgs, groups); len(tagFilters) > 0 {
		joinTables = append(joinTables, BuildJoinTagTables(baseTable, tagFilters))
	}

	if userGroupFilters := ExtractUserGroupFilters(ems, fgs, groups); len(userGroupFilters) > 0 {
		joinTables = append(joinTables, BuildJoinUserGroupFilterTables(baseTable, userGroupFilters))
	}

	return joinTables
}

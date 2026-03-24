package doris

import (
	"fmt"
)

type Column struct {
	Table string `json:"table"`
	Field string `json:"field"`
	Alias string `json:"alias"`
}

// 提取Column
func ExtractColumns(ems []EventMetric, fgs GlobalFilterGroups, groups []Group) []Column {
	columns := make([]Column, 0)
	for _, em := range ems {
		if em.Type == EventType_Custom || len(em.CustomMetric.EventMetrics) > 0 {
			for _, cem := range em.CustomMetric.EventMetrics {
				if cem.Metric.Column.Field != "" {
					columns = append(columns, cem.Metric.Column)
				}
				for _, filter := range cem.Filters {
					columns = append(columns, filter.Column)
				}
			}
		} else {

			if em.Metric.Column.Field != "" {
				columns = append(columns, em.Metric.Column)
			}
			for _, filter := range em.Filters {
				columns = append(columns, filter.Column)
			}
		}
	}

	if len(fgs.GlobalFilters.Filters) > 0 {
		for _, filter := range fgs.GlobalFilters.Filters {
			columns = append(columns, filter.Column)
		}
	}

	if len(fgs.DashBoardFormFilters.Filters) > 0 {
		for _, filter := range fgs.DashBoardFormFilters.Filters {
			columns = append(columns, filter.Column)
		}
	}

	for _, group := range groups {
		if group.GroupType == GroupBy_TagGroup {
			continue
		}
		columns = append(columns, group.Column)
	}
	return columns
}

func ConvertTzWithAlias(col Column, timeZone string) string {
	return ConvertTzField(col, timeZone) + " AS " + col.Field
}

func ConvertTzField(col Column, timeZone string) string {

	tableAlias := EVENT_TABLE_ALIAS
	switch col.Table {
	case EVENT_TABLE:
		tableAlias = EVENT_TABLE_ALIAS
	case USER_TABLE:
		tableAlias = USER_TABLE_ALIAS
	}
	tz := "+08:00"
	switch timeZone {
	case "+00:00":
		tz = "+00:00"
	case "+03:00":
		tz = "+03:00"
	case "-05:00":
		tz = "-05:00"
	case "-08:00":
		tz = "-08:00"
	case "+08:00", "":
		tz = "+08:00"
		return fmt.Sprintf("%s.%s", tableAlias, col.Field)
	}
	return fmt.Sprintf("convert_tz(%s.%s, '+08:00', '%s')", tableAlias, col.Field, tz)
}

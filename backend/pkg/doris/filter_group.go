package doris

import (
	"fmt"
	"strings"
)

type FilterScope int

const (
	ScopeOr  = 0
	ScopeAnd = 1
)

type FilterGroup struct {
	Scope            FilterScope           `json:"scope"`
	Filters          []Filter              `json:"filters"`
	TagFilters       []TagFilter           `json:"tag_filters"`
	UserGroupFilters []UserGroupDataFilter `json:"user_group_filters"`
}

func BuildFilterGroup(fg FilterGroup) string {
	filterSqls := make([]string, 0)

	if len(fg.TagFilters) > 0 {
		for _, tf := range fg.TagFilters {
			filterSqls = append(filterSqls, fmt.Sprintf("(%s)", BuildTagFilter(tf)))
		}
	}

	if len(fg.UserGroupFilters) > 0 {
		for _, ugf := range fg.UserGroupFilters {
			filterSqls = append(filterSqls, fmt.Sprintf("(%s)", BuildUserGroupFilterSql(ugf)))
		}
	}

	for _, filter := range fg.Filters {
		filterSql := BuildFilter(filter)
		if filterSql == "true" {
			continue
		}
		filterSqls = append(filterSqls, fmt.Sprintf("(%s)", filterSql))
	}

	splitScope := " OR "
	if fg.Scope == ScopeAnd {
		splitScope = " AND "
	}
	if len(filterSqls) > 0 {
		return fmt.Sprintf("(%s)", strings.Join(filterSqls, splitScope))
	}
	return "true"
}

type GlobalFilterGroups struct {
	GlobalFilters        FilterGroup `json:"global_filters"`
	DashBoardFormFilters FilterGroup `json:"dashboard_form_filters"`
	QueryDates           []string    `json:"query_dates"`
}

func BuildGlobalFilterGroups(fgs GlobalFilterGroups, tz string) []string {
	globalFilterSqls := make([]string, 0)
	if len(fgs.GlobalFilters.Filters) > 0 || len(fgs.GlobalFilters.TagFilters) > 0 {
		filterSql := BuildFilterGroup(fgs.GlobalFilters)
		if filterSql != "true" {
			globalFilterSqls = append(globalFilterSqls, filterSql)
		}
	}

	if len(fgs.DashBoardFormFilters.Filters) > 0 || len(fgs.DashBoardFormFilters.TagFilters) > 0 {
		filterSql := BuildFilterGroup(fgs.DashBoardFormFilters)
		if filterSql != "true" {
			globalFilterSqls = append(globalFilterSqls, filterSql)
		}
	}

	if len(fgs.QueryDates) > 0 && tz != "" {
		tzField := fmt.Sprintf("%s.e_event_time", EVENT_TABLE_ALIAS)
		queryDates := ConvertTzQueryDates(fgs.QueryDates, tz)

		if strings.Contains(queryDates[0], "convert_tz") {
			globalFilterSqls = append(globalFilterSqls, fmt.Sprintf("(%s >= %s AND %s < %s)",
				tzField, queryDates[0], tzField, queryDates[1]))
		} else {
			globalFilterSqls = append(globalFilterSqls, fmt.Sprintf("(%s >= '%s' AND %s < '%s')",
				tzField, queryDates[0], tzField, queryDates[1]))
		}
	}
	return globalFilterSqls
}

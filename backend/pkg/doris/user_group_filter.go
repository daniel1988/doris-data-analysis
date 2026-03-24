package doris

import (
	"fmt"
	"strings"
)

type UserGroupDataFilter struct {
	GroupName string `json:"group_name"`
	GroupCode string `json:"group_code"`
	Operator  int    `json:"operator"`
}

func BuildUserGroupFilterSql(ugf UserGroupDataFilter) string {
	switch ugf.Operator {
	case OperIn:
		return fmt.Sprintf("`%s`.user_id IS NOT NULL", ugf.GroupCode)
	case OperNotIn:
		return fmt.Sprintf("`%s`.user_id IS NULL", ugf.GroupCode)
	default:
		return fmt.Sprintf("`%s`.user_id IS NOT NULL", ugf.GroupCode)
	}
}

func ExtractUserGroupFilters(ems []EventMetric, fgs GlobalFilterGroups, groups []Group) []UserGroupDataFilter {
	userGroupFilters := make([]UserGroupDataFilter, 0)

	for _, em := range ems {
		if len(em.UserGroupFilters) > 0 {
			userGroupFilters = append(userGroupFilters, em.UserGroupFilters...)
		}
	}

	if len(fgs.GlobalFilters.UserGroupFilters) > 0 {
		userGroupFilters = append(userGroupFilters, fgs.GlobalFilters.UserGroupFilters...)
	}

	if len(fgs.DashBoardFormFilters.UserGroupFilters) > 0 {
		userGroupFilters = append(userGroupFilters, fgs.DashBoardFormFilters.UserGroupFilters...)
	}

	for _, group := range groups {
		if group.GroupType == GroupBy_UserGroupData {
			userGroupFilters = append(userGroupFilters, group.UserGroup)
		}
	}
	return userGroupFilters
}

func BuildUserGroupFilterTables(userGroupFilters []UserGroupDataFilter) []string {
	tables := make([]string, 0)
	for _, ugf := range userGroupFilters {
		tables = append(tables, fmt.Sprintf("`%s` AS (SELECT user_id,group_code,group_name FROM user_group_data WHERE group_code = '%s')",
			ugf.GroupCode, ugf.GroupCode))
	}
	return tables
}

func BuildJoinUserGroupFilterTables(baseTable string, userGroupFilters []UserGroupDataFilter) string {
	joinSqls := make([]string, 0)
	joinSqls = append(joinSqls, baseTable)
	joinField := "e_openid"
	if baseTable == USER_TABLE || baseTable == USER_TABLE_ALIAS {
		joinField = "u_openid"
	}
	for _, ugf := range userGroupFilters {
		joinSqls = append(joinSqls, fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.user_id",
			ugf.GroupCode, baseTable, joinField, ugf.GroupCode))
	}
	return strings.Join(joinSqls, " \n")
}

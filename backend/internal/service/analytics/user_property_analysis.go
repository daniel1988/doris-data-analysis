package analytics

import (
	"encoding/json"
	"fmt"
	"strings"

	"gitee.com/dmp_admin_v2/backend/internal/common"
	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/pkg/doris"
)

type UserPropertyAnalysisReq struct {
	model.QueryRequest
	Metric       doris.Metric             `json:"metric"`
	FilterGroups doris.GlobalFilterGroups `json:"filter_groups"`
	Groups       []doris.Group            `json:"groups"`
	UserGroups   []UserGroup              `json:"user_groups"`
	GroupType    int                      `json:"group_type"` // 1: 维度分组, 2: 人群分群
}

type UserGroup struct {
	Alias       string            `json:"alias"`
	FilterGroup doris.FilterGroup `json:"filter_group"`
}

func UserPropertyAnalysis(req *UserPropertyAnalysisReq) (*model.QueryResponse, error) {
	binJs, _ := json.Marshal(req)
	common.Logger.Infof("用户属性分析请求参数: %s", string(binJs))

	return NewSqlService().ExecuteQuery(&model.QueryRequest{
		ProjectAlias: req.ProjectAlias,
		SQL:          BuildUserPropertySql(req),
	})
}

func BuildUserPropertySql(req *UserPropertyAnalysisReq) string {
	selectFields := make([]string, 0)
	groupFields := make([]string, 0)
	whereSqls := make([]string, 0)

	// 1. 处理指标 (Metric)
	// 确保指标字段带有表别名
	if req.Metric.Column.Table == doris.USER_TABLE || req.Metric.Column.Table == "" {
		if req.Metric.Column.Field == "" {
			req.Metric.Column.Field = "u_openid"
		}
		req.Metric.Column.Field = fmt.Sprintf("%s.%s", doris.USER_TABLE_ALIAS, req.Metric.Column.Field)
	}

	metricSql := doris.BuildMetric(req.Metric)
	selectFields = append(selectFields, fmt.Sprintf("%s AS `指标值`", metricSql))

	// 2. 处理分组 (Grouping)
	if req.GroupType == 2 && len(req.UserGroups) > 0 {
		// 人群分群模式
		caseSqls := make([]string, 0)
		for _, ug := range req.UserGroups {
			fgSql := doris.BuildFilterGroup(ug.FilterGroup)
			caseSqls = append(caseSqls, fmt.Sprintf("WHEN %s THEN '%s'", fgSql, ug.Alias))
		}
		groupField := fmt.Sprintf("CASE %s ELSE '其他' END", strings.Join(caseSqls, " "))
		selectFields = append(selectFields, fmt.Sprintf("%s AS `分组`", groupField))
		groupFields = append(groupFields, "`分组`")
	} else {
		// 维度分组模式
		if len(req.Groups) > 0 {
			for _, group := range req.Groups {
				// 确保分组字段带有表别名
				if group.Column.Table == doris.USER_TABLE || group.Column.Table == "" {
					if group.Column.Field != "" {
						group.Column.Field = fmt.Sprintf("%s.%s", doris.USER_TABLE_ALIAS, group.Column.Field)
					}
				}

				selectField, groupField := doris.BuildGroup(group)
				if selectField != "" {
					selectFields = append(selectFields, selectField)
					groupFields = append(groupFields, groupField)
				}
			}
		}
	}

	// 3. 构建基础查询
	querySqls := make([]string, 0)
	querySqls = append(querySqls, fmt.Sprintf("SELECT %s", strings.Join(selectFields, ", ")))
	querySqls = append(querySqls, fmt.Sprintf("FROM %s AS %s", doris.USER_TABLE, doris.USER_TABLE_ALIAS))

	// 4. 处理全局过滤
	if fgSql := doris.BuildFilterGroup(req.FilterGroups.GlobalFilters); fgSql != "true" {
		whereSqls = append(whereSqls, fgSql)
	}

	if len(whereSqls) > 0 {
		querySqls = append(querySqls, fmt.Sprintf("WHERE %s", strings.Join(whereSqls, " AND ")))
	}

	// 5. 分组与排序
	if len(groupFields) > 0 {
		querySqls = append(querySqls, fmt.Sprintf("GROUP BY %s", strings.Join(groupFields, ", ")))
	}
	querySqls = append(querySqls, "ORDER BY `指标值` DESC")

	return strings.Join(querySqls, " \n")
}

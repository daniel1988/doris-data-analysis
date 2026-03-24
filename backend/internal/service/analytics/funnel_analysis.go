package analytics

import (
	"encoding/json"
	"fmt"
	"strings"

	"gitee.com/dmp_admin_v2/backend/internal/common"
	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/pkg/doris"
)

type FunnelAnalysisReq struct {
	model.QueryRequest
	Steps              []doris.EventMetric      `json:"steps"`
	GlobalFilterGroups doris.GlobalFilterGroups `json:"global_filter_groups"`
	Groups             []doris.Group            `json:"groups"`
	Window             int64                    `json:"window"` // 漏斗窗口期，秒
	TimeGrain          doris.TimeGrain          `json:"time_grain"`
}

func FunnelAnalysis(req *FunnelAnalysisReq) (*model.QueryResponse, error) {
	binJs, _ := json.Marshal(req)
	common.Logger.Infof("漏斗分析请求参数: %s", string(binJs))

	return NewSqlService().ExecuteQuery(&model.QueryRequest{
		ProjectAlias: req.ProjectAlias,
		SQL:          BuildFunnelSql(req),
	})
}

func BuildFunnelSql(req *FunnelAnalysisReq) string {
	// 构建 WITH 表
	withTables := doris.BuildExtraWithTables(req.Steps, req.GlobalFilterGroups, req.Groups)
	withTables = append(withTables, fmt.Sprintf("%s AS (%s)",
		doris.EVENT_TABLE_ALIAS,
		doris.BuildEventDataTable(req.Steps, req.TimeGrain, req.GlobalFilterGroups, req.Groups, req.TimeZone)))

	// 构建 window_funnel SQL
	stepConditions := make([]string, 0)
	for _, step := range req.Steps {
		// 每个步骤的条件包含：事件ID + 步骤内过滤器
		conds := []string{fmt.Sprintf("e_event_id = '%s'", step.EventId)}

		// 步骤内过滤器
		fg := doris.FilterGroup{
			Scope:      doris.FilterScope(step.Scope),
			Filters:    step.Filters,
			TagFilters: step.TagFilters,
		}
		if fgSql := doris.BuildFilterGroup(fg); fgSql != "true" {
			conds = append(conds, fgSql)
		}

		stepConditions = append(stepConditions, fmt.Sprintf("(%s)", strings.Join(conds, " AND ")))
	}

	selectFields := make([]string, 0)
	groupFields := make([]string, 0)

	if len(req.Groups) > 0 {
		for _, group := range req.Groups {
			selectField, groupField := doris.BuildGroup(group)
			selectFields = append(selectFields, selectField)
			groupFields = append(groupFields, groupField)
		}
	}

	// 使用 Doris window_funnel 函数
	funnelFormula := fmt.Sprintf("window_funnel(%d, 'fixed', e_event_time, %s)",
		req.Window, strings.Join(stepConditions, ", "))

	selectFields = append(selectFields, fmt.Sprintf("%s AS step", funnelFormula))

	// 最终查询
	querySqls := make([]string, 0)
	querySqls = append(querySqls, fmt.Sprintf("WITH %s", strings.Join(withTables, ", \n")))

	// 子查询计算每个用户的最大步骤
	subQuerySelect := append([]string{"e_openid"}, selectFields...)
	subQueryGroupBy := append([]string{"e_openid"}, groupFields...)

	subQuery := fmt.Sprintf("SELECT %s FROM %s GROUP BY %s",
		strings.Join(subQuerySelect, ", "),
		doris.EVENT_TABLE_ALIAS,
		strings.Join(subQueryGroupBy, ", "))

	// 外层查询统计各步骤人数
	finalSelect := make([]string, 0)
	if len(groupFields) > 0 {
		finalSelect = append(finalSelect, groupFields...)
	}

	for i := 1; i <= len(req.Steps); i++ {
		colName := fmt.Sprintf("步骤%d", i)
		finalSelect = append(finalSelect, fmt.Sprintf("SUM(CASE WHEN step >= %d THEN 1 ELSE 0 END) AS `%s`", i, colName))

		// 计算转化率
		if i > 1 {
			// 步骤转化率：当前步 / 上一步
			finalSelect = append(finalSelect, fmt.Sprintf("ROUND(SUM(CASE WHEN step >= %d THEN 1 ELSE 0 END) / NULLIF(SUM(CASE WHEN step >= %d THEN 1 ELSE 0 END), 0), 4) AS `%s转化率`", i, i-1, colName))

			// 总转化率：当前步 / 第一步
			finalSelect = append(finalSelect, fmt.Sprintf("ROUND(SUM(CASE WHEN step >= %d THEN 1 ELSE 0 END) / NULLIF(SUM(CASE WHEN step >= 1 THEN 1 ELSE 0 END), 0), 4) AS `%s总转化率`", i, colName))
		}
	}

	querySqls = append(querySqls, fmt.Sprintf("SELECT %s FROM (%s) t", strings.Join(finalSelect, ", "), subQuery))
	if len(groupFields) > 0 {
		querySqls = append(querySqls, fmt.Sprintf("GROUP BY %s", strings.Join(groupFields, ", ")))
	}

	return strings.Join(querySqls, " \n")
}

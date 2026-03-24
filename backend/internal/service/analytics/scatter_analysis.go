package analytics

import (
	"encoding/json"
	"fmt"
	"strings"

	"gitee.com/dmp_admin_v2/backend/internal/common"
	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/pkg/doris"
)

type ScatterAnalysisReq struct {
	model.QueryRequest
	Metric             doris.EventMetric      `json:"metric"`
	GlobalFilterGroups doris.GlobalFilterGroups `json:"global_filter_groups"`
	Groups             []doris.Group            `json:"groups"`
	TimeGrain          doris.TimeGrain          `json:"time_grain"`
}

func ScatterAnalysis(req *ScatterAnalysisReq) (*model.QueryResponse, error) {
	binJs, _ := json.Marshal(req)
	common.Logger.Infof("分布分析请求参数: %s", string(binJs))

	return NewSqlService().ExecuteQuery(&model.QueryRequest{
		ProjectAlias: req.ProjectAlias,
		SQL:          BuildScatterSql(req),
	})
}

func BuildScatterSql(req *ScatterAnalysisReq) string {
	ems := []doris.EventMetric{req.Metric}
	withTables := doris.BuildExtraWithTables(ems, req.GlobalFilterGroups, req.Groups)
	withTables = append(withTables, fmt.Sprintf("%s AS (%s)",
		doris.EVENT_TABLE_ALIAS,
		doris.BuildEventDataTable(ems, req.TimeGrain, req.GlobalFilterGroups, req.Groups, req.TimeZone)))

	selectFields := make([]string, 0)
	groupFields := make([]string, 0)

	if len(req.Groups) > 0 {
		for _, group := range req.Groups {
			selectField, groupField := doris.BuildGroup(group)
			selectFields = append(selectFields, selectField)
			groupFields = append(groupFields, groupField)
		}
	}

	metricFormula := doris.BuildEventMetric(req.Metric, req.TimeGrain)
	selectFields = append(selectFields, fmt.Sprintf("%s AS metric_val", metricFormula))

	querySqls := make([]string, 0)
	querySqls = append(querySqls, fmt.Sprintf("WITH %s", strings.Join(withTables, ", \n")))

	// 子查询计算每个用户的指标值
	subQuerySelect := append([]string{"e_openid"}, selectFields...)
	subQueryGroupBy := append([]string{"e_openid"}, groupFields...)

	subQuery := fmt.Sprintf("SELECT %s FROM %s GROUP BY %s",
		strings.Join(subQuerySelect, ", "),
		doris.EVENT_TABLE_ALIAS,
		strings.Join(subQueryGroupBy, ", "))

	// 外层查询按指标值分区间统计人数
	finalSelect := make([]string, 0)
	if len(groupFields) > 0 {
		finalSelect = append(finalSelect, groupFields...)
	}

	finalSelect = append(finalSelect, "metric_val AS `区间`", "COUNT(DISTINCT e_openid) AS `用户数`")

	querySqls = append(querySqls, fmt.Sprintf("SELECT %s FROM (%s) t", strings.Join(finalSelect, ", "), subQuery))

	finalGroupBy := append(groupFields, "metric_val")
	querySqls = append(querySqls, fmt.Sprintf("GROUP BY %s", strings.Join(finalGroupBy, ", ")))
	querySqls = append(querySqls, "ORDER BY metric_val ASC")

	return strings.Join(querySqls, " \n")
}

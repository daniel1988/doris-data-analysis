package analytics

import (
	"encoding/json"
	"fmt"
	"strings"

	"gitee.com/dmp_admin_v2/backend/internal/common"
	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/pkg/doris"
)

type EventAnalysisReq struct {
	model.QueryRequest
	EventMetrics []doris.EventMetric      `json:"event_metrics"`
	FilterGroups doris.GlobalFilterGroups `json:"filter_groups"`
	TimeGrain    doris.TimeGrain          `json:"time_grain"`
	Groups       []doris.Group            `json:"groups"`
	Orders       []doris.Order            `json:"orders"`
}

func EventAnalysis(req *EventAnalysisReq) (interface{}, error) {
	binJs, _ := json.Marshal(req)
	common.Logger.Infof("事件分析请求参数: %s", string(binJs))
	return NewSqlService().ExecuteQuery(&model.QueryRequest{
		ProjectAlias: req.ProjectAlias,
		SQL:          BuildEventAnalysisSql(req),
	})
}

func BuildEventAnalysisSql(req *EventAnalysisReq) string {
	withTables := doris.BuildExtraWithTables(req.EventMetrics, req.FilterGroups, req.Groups)

	withTables = append(withTables, fmt.Sprintf("%s AS (%s)",
		doris.EVENT_TABLE_ALIAS,
		doris.BuildEventDataTable(req.EventMetrics, req.TimeGrain, req.FilterGroups, req.Groups, req.TimeZone)))
	selectMetrics := make([]string, 0)
	groupFields := make([]string, 0)

	if req.TimeGrain.Interval != doris.Tg_Interval_Empty {
		tgFormula := doris.BuildTimeGrainFormula(req.TimeGrain)
		selectMetrics = append(selectMetrics, fmt.Sprintf("%s AS `%s`", tgFormula, req.TimeGrain.Column.Alias))

		groupFields = append(groupFields, req.TimeGrain.Column.Alias)
	}

	if len(req.Groups) > 0 {
		for _, group := range req.Groups {
			selectField, groupField := doris.BuildGroup(group)
			selectMetrics = append(selectMetrics, selectField)
			groupFields = append(groupFields, groupField)
		}
	}

	for _, em := range req.EventMetrics {
		metricFormula := doris.BuildEventMetric(em, req.TimeGrain)
		selectMetrics = append(selectMetrics, fmt.Sprintf("%s AS `%s`", metricFormula, em.Name))
	}

	eventQuerySqls := make([]string, 0)
	eventQuerySqls = append(eventQuerySqls, fmt.Sprintf("WITH %s", strings.Join(withTables, ", \n")))
	eventQuerySqls = append(eventQuerySqls, fmt.Sprintf("SELECT %s", strings.Join(selectMetrics, ", \n")))
	eventQuerySqls = append(eventQuerySqls, "FROM events")

	if joinSqls := doris.BuildExtraJoinTables(doris.EVENT_TABLE_ALIAS, req.EventMetrics, req.FilterGroups, req.Groups); len(joinSqls) > 0 {
		eventQuerySqls = append(eventQuerySqls, strings.Join(joinSqls, " \n"))
	}

	whereSqls := make([]string, 0)
	if fgSqls := doris.BuildGlobalFilterGroups(req.FilterGroups, ""); len(fgSqls) > 0 {
		whereSqls = append(whereSqls, fmt.Sprintf("(%s)", strings.Join(fgSqls, " AND ")))
	}
	if len(whereSqls) > 0 {
		eventQuerySqls = append(eventQuerySqls, fmt.Sprintf("WHERE %s", strings.Join(whereSqls, " AND ")))
	}
	if len(groupFields) > 0 {
		eventQuerySqls = append(eventQuerySqls, fmt.Sprintf("GROUP BY %s", strings.Join(groupFields, ", \n")))
	}

	if len(req.Orders) > 0 {
		orderSqls := doris.BuildOrders(req.Orders)
		if len(orderSqls) > 0 {
			eventQuerySqls = append(eventQuerySqls, fmt.Sprintf("ORDER BY %s", strings.Join(orderSqls, ", \n")))
		}
	} else {
		if req.TimeGrain.Interval != doris.Tg_Interval_Empty {
			eventQuerySqls = append(eventQuerySqls, "ORDER BY 日期 DESC")
		}
	}
	return strings.Join(eventQuerySqls, " \n")
}

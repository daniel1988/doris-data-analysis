package analytics

import (
	"encoding/json"
	"fmt"
	"strings"

	"gitee.com/dmp_admin_v2/backend/internal/common"
	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/pkg/doris"
)

const (
	Init_Table = "itbl"
	End_Table  = "etbl"
)

type RetentionAnalysisReq struct {
	model.QueryRequest
	InitEventMetric          doris.EventMetric `json:"init_event_metric"`
	EndEventMetric           doris.EventMetric `json:"end_event_metric"`
	doris.GlobalFilterGroups `json:"global_filter_groups"`
	Groups                   []doris.Group   `json:"groups"`
	TimeGrain                doris.TimeGrain `json:"time_grain"`
	DayNArray                []int           `json:"day_n_array"`
}

func RetentionAnalysis(req *RetentionAnalysisReq) (*model.QueryResponse, error) {
	binJs, _ := json.Marshal(req)
	common.Logger.Infof("留存分析请求参数: %s", string(binJs))

	return NewSqlService().ExecuteQuery(&model.QueryRequest{
		ProjectAlias: req.ProjectAlias,
		SQL:          BuildRetentionSql(req),
	})
}

func BuildRetentionSql(req *RetentionAnalysisReq) string {
	ems := []doris.EventMetric{req.InitEventMetric, req.EndEventMetric}
	tabFilters := doris.ExtractTagFilters(ems, req.GlobalFilterGroups, req.Groups)
	withTables := doris.BuildTagTables(tabFilters)
	withTables = append(withTables, fmt.Sprintf("%s AS (%s)",
		doris.EVENT_TABLE_ALIAS,
		doris.BuildEventDataTable(ems, req.TimeGrain, req.GlobalFilterGroups, req.Groups, req.TimeZone)))

	withTables = append(withTables, fmt.Sprintf("%s AS (%s)", Init_Table, buildInitRetentionMetric(req)))
	withTables = append(withTables, fmt.Sprintf("%s AS (%s)", End_Table, buildEndRetentionMetric(req)))

	querySqls := make([]string, 0)
	querySqls = append(querySqls, fmt.Sprintf("WITH %s", strings.Join(withTables, ",\n ")))
	selectFields := make([]string, 0)
	groupFields := make([]string, 0)
	selectFields = append(selectFields, fmt.Sprintf("%s.日期", Init_Table))
	groupFields = append(groupFields, fmt.Sprintf("%s.日期", Init_Table))
	for _, group := range req.Groups {
		groupAlias, _ := doris.BuildGroup(group)
		selectFields = append(selectFields, fmt.Sprintf("%s.%s", Init_Table, groupAlias))
		groupFields = append(groupFields, fmt.Sprintf("%s.%s", Init_Table, groupAlias))
	}

	selectFields = append(selectFields, fmt.Sprintf("COUNT(DISTINCT %s.e_openid) AS `用户数`", Init_Table))
	for _, dayN := range req.DayNArray {
		selectFields = append(selectFields, calculateRetentionDayN(dayN))
	}
	querySqls = append(querySqls, fmt.Sprintf("SELECT %s", strings.Join(selectFields, ",\n ")))
	querySqls = append(querySqls, fmt.Sprintf("FROM %s", "itbl"))
	querySqls = append(querySqls, fmt.Sprintf("LEFT JOIN %s ON itbl.e_openid = etbl.e_openid", "etbl"))

	querySqls = append(querySqls, fmt.Sprintf("GROUP BY %s", strings.Join(groupFields, ", ")))
	querySqls = append(querySqls, fmt.Sprintf("ORDER BY %s.日期 DESC", Init_Table))

	return strings.Join(querySqls, " \n ")
}

func calculateRetentionDayN(dayN int) string {
	formula := fmt.Sprintf("ROUND(COUNT(DISTINCT CASE WHEN DATEDIFF(%s.日期, %s.日期) = %d THEN %s.e_openid END) / COUNT(DISTINCT %s.e_openid), 4) AS `%v日`",
		End_Table, Init_Table, dayN, End_Table, Init_Table, dayN)
	return formula
}

func buildInitRetentionMetric(req *RetentionAnalysisReq) string {
	selectFields := make([]string, 0)
	groupFields := make([]string, 0)
	tgAlias := "日期"
	tg := doris.BuildTimeGrainFormula(req.TimeGrain)
	selectFields = append(selectFields, fmt.Sprintf("%s AS %s", tg, tgAlias))
	groupFields = append(groupFields, tgAlias)

	selectFields = append(selectFields, fmt.Sprintf("%s.e_openid", doris.EVENT_TABLE_ALIAS))
	groupFields = append(groupFields, fmt.Sprintf("%s.e_openid", doris.EVENT_TABLE_ALIAS))
	if len(req.Groups) > 0 {
		for _, group := range req.Groups {
			selectField, groupField := doris.BuildGroup(group)
			selectFields = append(selectFields, selectField)
			groupFields = append(groupFields, groupField)
		}
	}
	initSqls := make([]string, 0)

	initSqls = append(initSqls, fmt.Sprintf("SELECT %s", strings.Join(selectFields, ", ")))
	initSqls = append(initSqls, fmt.Sprintf("FROM %s", doris.EVENT_TABLE_ALIAS))
	initEms := []doris.EventMetric{req.InitEventMetric}
	if joinTables := doris.BuildExtraJoinTables(doris.EVENT_TABLE_ALIAS, initEms, req.GlobalFilterGroups, req.Groups); len(joinTables) > 0 {
		initSqls = append(initSqls, strings.Join(joinTables, " "))
	}
	whereSqls := make([]string, 0)
	whereSqls = append(whereSqls, fmt.Sprintf("e_event_id = '%s'", req.InitEventMetric.EventId))

	if fgSql := doris.BuildFilterGroup(req.InitEventMetric.FilterGroup); fgSql != "true" {
		whereSqls = append(whereSqls, fgSql)
	}

	if fgSqls := doris.BuildGlobalFilterGroups(req.GlobalFilterGroups, ""); len(fgSqls) > 0 {
		whereSqls = append(whereSqls, fmt.Sprintf("(%s)", strings.Join(fgSqls, " AND ")))
	}

	if len(whereSqls) > 0 {
		initSqls = append(initSqls, fmt.Sprintf("WHERE %s", strings.Join(whereSqls, " AND ")))
	}
	initSqls = append(initSqls, fmt.Sprintf("GROUP BY %s", strings.Join(groupFields, ", ")))
	return strings.Join(initSqls, " \n ")
}

func buildEndRetentionMetric(req *RetentionAnalysisReq) string {
	selectFields := make([]string, 0)
	groupFields := make([]string, 0)
	tg := doris.BuildTimeGrainFormula(req.TimeGrain)
	tgAlias := "日期"
	selectFields = append(selectFields, fmt.Sprintf("%s AS %s", tg, tgAlias))
	groupFields = append(groupFields, tgAlias)

	selectFields = append(selectFields, fmt.Sprintf("%s.e_openid", doris.EVENT_TABLE_ALIAS))
	groupFields = append(groupFields, fmt.Sprintf("%s.e_openid", doris.EVENT_TABLE_ALIAS))
	if len(req.Groups) > 0 {
		for _, group := range req.Groups {
			selectField, groupField := doris.BuildGroup(group)
			selectFields = append(selectFields, selectField)
			groupFields = append(groupFields, groupField)
		}
	}

	endSqls := make([]string, 0)
	endSqls = append(endSqls, fmt.Sprintf("SELECT %s", strings.Join(selectFields, ",")))
	endSqls = append(endSqls, fmt.Sprintf("FROM %s", doris.EVENT_TABLE_ALIAS))
	endEms := []doris.EventMetric{req.EndEventMetric}
	if joinTables := doris.BuildExtraJoinTables(doris.EVENT_TABLE_ALIAS, endEms, req.GlobalFilterGroups, req.Groups); len(joinTables) > 0 {
		endSqls = append(endSqls, strings.Join(joinTables, " "))
	}
	whereSqls := make([]string, 0)
	whereSqls = append(whereSqls, fmt.Sprintf("e_event_id = '%s'", req.EndEventMetric.EventId))
	if fgSql := doris.BuildFilterGroup(req.InitEventMetric.FilterGroup); fgSql != "true" {
		whereSqls = append(whereSqls, fgSql)
	}

	if fgSqls := doris.BuildGlobalFilterGroups(req.GlobalFilterGroups, ""); len(fgSqls) > 0 {
		whereSqls = append(whereSqls, fmt.Sprintf("(%s)", strings.Join(fgSqls, " AND ")))
	}
	if len(whereSqls) > 0 {
		endSqls = append(endSqls, fmt.Sprintf("WHERE %s", strings.Join(whereSqls, " AND ")))
	}
	endSqls = append(endSqls, fmt.Sprintf("GROUP BY %s", strings.Join(groupFields, ", ")))
	return strings.Join(endSqls, " \n ")
}

package analytics

import (
	"fmt"
	"strings"

	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/internal/service/meta"
	"gitee.com/dmp_admin_v2/backend/pkg/doris"
)

var (
	defaultEvSelectFields = []string{
		"e_openid",
		"e_event_id",
		"e_event_time",
		"e_event_name",
		"e_properties",
	}

	// 基础必选字段，不受前端控制
	baseEvFields = []string{
		"e_openid",
		"e_event_id",
		"e_event_time",
	}
)

type EventDetailReq struct {
	ProjectAlias     string            `json:"project_alias"`
	PageSize         int               `json:"page_size"`
	PageNum          int               `json:"page_num"`
	SelectFields     []string          `json:"select_fields"`
	EventFilterGroup doris.FilterGroup `json:"event_filter_group"`
}

func EventDetail(req *EventDetailReq) (interface{}, error) {
	// 1. 获取该项目在元数据中定义的属性列表
	propSrv := meta.NewProjectPropertyService()
	props, err := propSrv.GetList(req.ProjectAlias)
	if err != nil {
		return nil, fmt.Errorf("failed to get project properties: %v", err)
	}

	// 2. 建立有效字段集合 (Doris 表中的实际列)
	validFields := make(map[string]bool)
	// 添加基础字段
	for _, f := range defaultEvSelectFields {
		validFields[f] = true
	}
	// 添加扩展字段 (e_platform, e_package_name, e_ip, e_request_id 等通常也在基础字段或元数据中)
	validFields["e_platform"] = true
	validFields["e_ip"] = true
	validFields["e_request_id"] = true

	// 添加元数据中定义的属性字段
	for _, p := range props {
		validFields[p.PropertyId] = true
	}

	// 3. 校验并剔除无效字段
	finalSelectFields := make([]string, 0)
	uniqueFields := make(map[string]bool)

	// 强制包含基础字段
	for _, f := range baseEvFields {
		finalSelectFields = append(finalSelectFields, f)
		uniqueFields[f] = true
	}

	// 处理前端传来的字段
	for _, f := range req.SelectFields {
		if validFields[f] && !uniqueFields[f] {
			finalSelectFields = append(finalSelectFields, f)
			uniqueFields[f] = true
		}
	}

	// 如果没有传有效字段，则使用默认字段
	if len(finalSelectFields) <= len(baseEvFields) {
		for _, f := range defaultEvSelectFields {
			if !uniqueFields[f] {
				finalSelectFields = append(finalSelectFields, f)
				uniqueFields[f] = true
			}
		}
	}

	req.SelectFields = finalSelectFields

	return NewSqlService().ExecuteQuery(&model.QueryRequest{
		ProjectAlias: req.ProjectAlias,
		PageSize:     req.PageSize,
		PageNum:      req.PageNum,
		SQL:          buildEventDetailSql(req),
	})
}

func buildEventDetailSql(req *EventDetailReq) string {
	selectFields := req.SelectFields
	if len(selectFields) == 0 {
		selectFields = append(selectFields, defaultEvSelectFields...)
	}

	eventDetailSqls := make([]string, 0)
	eventDetailSqls = append(eventDetailSqls, fmt.Sprintf("SELECT %s", strings.Join(selectFields, ",")))
	eventDetailSqls = append(eventDetailSqls, fmt.Sprintf("FROM %s %s", doris.EVENT_TABLE, doris.EVENT_TABLE_ALIAS))

	if fgSql := doris.BuildFilterGroup(req.EventFilterGroup); fgSql != "true" {
		eventDetailSqls = append(eventDetailSqls, fmt.Sprintf("WHERE %s", fgSql))
	}

	eventDetailSqls = append(eventDetailSqls, "ORDER BY e_event_time DESC")

	return strings.Join(eventDetailSqls, " ")
}

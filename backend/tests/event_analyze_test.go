package tests

import (
	"fmt"
	"testing"

	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/internal/service/analytics"
	"gitee.com/dmp_admin_v2/backend/pkg/doris"
)

func TestEventAnalyze(t *testing.T) {

	req := &analytics.EventAnalysisReq{
		QueryRequest: model.QueryRequest{
			ProjectAlias: "zgmgr4",
		},
		EventMetrics: []doris.EventMetric{
			{
				Type:    doris.EventType_Normal,
				EventId: "login",
				Name:    "登录",
				Metric: doris.Metric{
					Formula: doris.Formula_Count,
					Column: doris.Column{
						Field: "e_openid",
						Table: doris.EVENT_TABLE,
					},
				},
				FilterGroup: doris.FilterGroup{
					TagFilters: []doris.TagFilter{
						{
							TagCode:  "cjxcx3_20250908_1018",
							Operator: doris.OperGreaterOrEqual,
							TagValue: "1",
						},
					},
				},
			},
		},
		TimeGrain: doris.TimeGrain{
			Column: doris.Column{
				Field: "e_event_time",
				Table: doris.EVENT_TABLE,
				Alias: "日期",
			},
			Interval: doris.Tg_Interval_Day,
		},
	}

	sql := analytics.BuildEventAnalysisSql(req)
	fmt.Println(sql)
}

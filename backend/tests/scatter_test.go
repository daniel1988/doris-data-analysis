package tests

import (
	"fmt"
	"testing"

	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/internal/service/analytics"
	"gitee.com/dmp_admin_v2/backend/pkg/doris"
)

func TestScatterQuery(t *testing.T) {

	req := analytics.ScatterAnalysisReq{
		QueryRequest: model.QueryRequest{
			ProjectAlias: "test_project",
		},
		Metric: doris.EventMetric{
			Type:    1,
			EventId: "sys.login",
			Name:    "login_count",
			Metric: doris.Metric{
				Formula: doris.Formula_Count,
				Column: doris.Column{
					Field: "e_openid",
					Table: doris.EVENT_TABLE,
				},
			},
		},
		GlobalFilterGroups: doris.GlobalFilterGroups{
			QueryDates: []string{"2025-11-01", "2025-11-13"},
		},
		TimeGrain: doris.TimeGrain{
			Interval: doris.Tg_Interval_Day,
			Column: doris.Column{
				Field: "e_event_time",
				Table: doris.EVENT_TABLE,
				Alias: "日期",
			},
		},
	}

	sql := analytics.BuildScatterSql(&req)
	fmt.Println(sql)
}

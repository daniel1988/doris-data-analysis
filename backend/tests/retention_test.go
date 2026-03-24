package tests

import (
	"fmt"
	"testing"

	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/internal/service/analytics"
	"gitee.com/dmp_admin_v2/backend/pkg/doris"
)

func TestRetention(t *testing.T) {
	req := &analytics.RetentionAnalysisReq{
		QueryRequest: model.QueryRequest{
			ProjectAlias: "test",
		},
		InitEventMetric: doris.EventMetric{
			EventId: "sys.register",
		},
		EndEventMetric: doris.EventMetric{
			EventId: "sys.login",
		},
		DayNArray: []int{1, 2, 3, 4, 5, 6, 7},
		TimeGrain: doris.TimeGrain{
			Column: doris.Column{
				Table: "event_data",
				Field: "e_event_time",
			},
			Interval: doris.Tg_Interval_Day,
		},
		GlobalFilterGroups: doris.GlobalFilterGroups{
			QueryDates: []string{"2025-10-01", "2025-10-07"},
		},
	}

	sql := analytics.BuildRetentionSql(req)
	fmt.Println(sql)
}

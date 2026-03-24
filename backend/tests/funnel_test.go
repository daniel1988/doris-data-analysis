package tests

import (
	"fmt"
	"testing"

	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/internal/service/analytics"
	"gitee.com/dmp_admin_v2/backend/pkg/doris"
)

func TestFunnel(t *testing.T) {

	req := analytics.FunnelAnalysisReq{
		QueryRequest: model.QueryRequest{
			ProjectAlias: "zgmgr4",
		},
		Steps: []doris.EventMetric{
			{
				Type:    doris.EventType_Normal,
				EventId: "sys.register",
				Name:    "注册",
				Metric: doris.Metric{
					Formula: doris.Formula_Bitmap_Union,
					Column: doris.Column{
						Field: "e_openid",
						Table: doris.EVENT_TABLE,
					},
				},
			},
			{
				Type:    doris.EventType_Normal,
				EventId: "sys.login",
				Name:    "登录",
				Metric: doris.Metric{
					Formula: doris.Formula_Bitmap_Union,
					Column: doris.Column{
						Field: "e_openid",
						Table: doris.EVENT_TABLE,
					},
				},
			},
			{
				Type:    doris.EventType_Normal,
				EventId: "sys.ad",
				Name:    "广告",
				Metric: doris.Metric{
					Formula: doris.Formula_Bitmap_Union,
					Column: doris.Column{
						Field: "e_openid",
						Table: doris.EVENT_TABLE,
					},
				},
			},
		},
		TimeGrain: doris.TimeGrain{
			Column: doris.Column{
				Table: doris.EVENT_TABLE,
				Field: "e_event_time",
				Alias: "日期",
			},
			Interval: doris.Tg_Interval_Day,
		},
		Window: 3600, // 1 hour
	}

	sql := analytics.BuildFunnelSql(&req)
	fmt.Println(sql)
}

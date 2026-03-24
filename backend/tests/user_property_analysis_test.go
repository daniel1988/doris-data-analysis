package tests

import (
	"fmt"
	"testing"

	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/internal/service/analytics"
	"gitee.com/dmp_admin_v2/backend/pkg/doris"
)

func TestUserPropertyAnalysis(t *testing.T) {

	req := analytics.UserPropertyAnalysisReq{
		QueryRequest: model.QueryRequest{
			ProjectAlias: "test_project",
		},
		Property: doris.Column{
			Table: doris.USER_TABLE,
			Field: "province",
		},
		Groups: []doris.Group{
			{
				Column: doris.Column{
					Field: "city",
					Table: doris.USER_TABLE,
				},
			},
		},
	}

	sql := analytics.BuildUserPropertySql(&req)
	fmt.Println(sql)
}

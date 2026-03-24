package tests

import (
	"fmt"
	"testing"

	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/internal/service/analytics"
	"gitee.com/dmp_admin_v2/backend/pkg/doris"
)

func TestUserList(t *testing.T) {
	req := &analytics.UserListReq{
		QueryRequest: model.QueryRequest{
			ProjectAlias: "cjxcx3",
			PageSize:     10,
			PageNum:      1,
		},
		Columns: []doris.Column{
			{Field: "e_openid", Alias: "用户ID", Table: doris.USER_TABLE},
			{Field: "u_event_time", Alias: "激活时间", Table: doris.USER_TABLE},
		},
	}

	sql := analytics.BuildUserListSql(req)
	fmt.Println(sql)
}

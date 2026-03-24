package analytics

import (
	"encoding/json"
	"fmt"
	"strings"

	"gitee.com/dmp_admin_v2/backend/internal/common"
	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/pkg/doris"
)

type UserListReq struct {
	model.QueryRequest
	GlobalFilterGroups doris.GlobalFilterGroups `json:"global_filter_groups"`
	Columns            []doris.Column           `json:"columns"`
}

func GetUserList(req *UserListReq) (*model.QueryResponse, error) {
	binJs, _ := json.Marshal(req)
	common.Logger.Infof("用户列表请求参数: %s", string(binJs))

	return NewSqlService().ExecuteQuery(&model.QueryRequest{
		ProjectAlias: req.ProjectAlias,
		SQL:          BuildUserListSql(req),
		PageSize:     req.PageSize,
		PageNum:      req.PageNum,
	})
}

func BuildUserListSql(req *UserListReq) string {
	selectFields := make([]string, 0)
	for _, col := range req.Columns {
		tableAlias := doris.USER_TABLE_ALIAS
		if col.Table == doris.EVENT_TABLE {
			tableAlias = doris.EVENT_TABLE_ALIAS
		}
		selectFields = append(selectFields, fmt.Sprintf("%s.%s AS `%s`", tableAlias, col.Field, col.Alias))
	}

	querySqls := make([]string, 0)
	querySqls = append(querySqls, fmt.Sprintf("SELECT %s", strings.Join(selectFields, ", ")))
	querySqls = append(querySqls, fmt.Sprintf("FROM %s AS %s", doris.USER_TABLE, doris.USER_TABLE_ALIAS))

	whereSqls := make([]string, 0)
	if fgSqls := doris.BuildGlobalFilterGroups(req.GlobalFilterGroups, doris.USER_TABLE_ALIAS); len(fgSqls) > 0 {
		whereSqls = append(whereSqls, fmt.Sprintf("(%s)", strings.Join(fgSqls, " AND ")))
	}

	if len(whereSqls) > 0 {
		querySqls = append(querySqls, fmt.Sprintf("WHERE %s", strings.Join(whereSqls, " AND ")))
	}

	return strings.Join(querySqls, " \n")
}

package analytics

import (
	"encoding/json"
	"fmt"
	"time"

	"gitee.com/dmp_admin_v2/backend/internal/common"
	"gitee.com/dmp_admin_v2/backend/internal/core"
	"gitee.com/dmp_admin_v2/backend/pkg/doris"

	"go.uber.org/zap"
)

type DimensionReq struct {
	ProjectAlias string `json:"project_alias"`
	Table        string `json:"table"`
	Field        string `json:"field"`
	EventId      string `json:"e_event_id"`
}

func GetDimensions(req *DimensionReq) ([]map[string]interface{}, error) {
	binJs, _ := json.Marshal(req)
	common.Logger.Infof("GetDimensions \n%s", string(binJs))

	if req.Table == doris.TAG_TABLE {
		return TagValues(req)
	}

	dorisDB, err := core.GetProjectDoris(req.ProjectAlias)
	if err != nil {
		common.Logger.Error("Failed to get Doris connection for project",
			zap.String("project", req.ProjectAlias),
			zap.Error(err))
		return nil, fmt.Errorf("project %s not found or not available", req.ProjectAlias)
	}

	if req.Table != doris.EVENT_TABLE && req.Table != doris.USER_TABLE {
		return nil, fmt.Errorf("table %s not found or not available", req.Table)
	}

	// 采用原生 SQL 以避免 GORM 生成的 LIMIT 语句在 Doris 中报错
	query := fmt.Sprintf("SELECT DISTINCT `%s` AS value FROM `%s` WHERE 1=1 ", req.Field, req.Table)
	args := []interface{}{}

	switch req.Table {
	case doris.EVENT_TABLE:
		if req.EventId != "" {
			query += " AND e_event_id = ? "
			args = append(args, req.EventId)
		}
		query += " AND e_event_time >= ? "
		args = append(args, time.Now().AddDate(0, 0, -7))
	case doris.USER_TABLE:
		query += " AND u_event_time >= ? "
		args = append(args, time.Now().AddDate(0, 0, -7))
	}

	query += " ORDER BY value ASC LIMIT 100"

	results := make([]map[string]interface{}, 0)
	err = dorisDB.Raw(query, args...).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func TagValues(req *DimensionReq) ([]map[string]interface{}, error) {
	dorisDB, err := core.GetProjectDoris(req.ProjectAlias)
	if err != nil {
		common.Logger.Error("Failed to get Doris connection for project",
			zap.String("project", req.ProjectAlias),
			zap.Error(err))
		return nil, fmt.Errorf("project %s not found or not available", req.ProjectAlias)
	}

	// 采用原生 SQL
	query := fmt.Sprintf("SELECT DISTINCT tag_value AS value FROM `%s` WHERE tag_code = ? ORDER BY value ASC LIMIT 100", doris.TAG_TABLE)
	results := make([]map[string]interface{}, 0)
	err = dorisDB.Raw(query, req.Field).Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

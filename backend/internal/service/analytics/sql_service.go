package analytics

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitee.com/dmp_admin_v2/backend/internal/common"
	"gitee.com/dmp_admin_v2/backend/internal/core"
	"gitee.com/dmp_admin_v2/backend/internal/model"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SqlService SQL服务
type SqlService struct{}

// NewSqlService 创建SQL服务实例
func NewSqlService() *SqlService {
	return &SqlService{}
}

// ExecuteQuery 执行SQL查询
func (x *SqlService) ExecuteQuery(req *model.QueryRequest) (*model.QueryResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	startTime := time.Now()

	common.Logger.Infof("项目：%s\nSQL:\n%s", req.ProjectAlias, req.SQL)

	userID, _ := strconv.Atoi(req.UserID)
	// 创建查询日志记录
	queryLog := &model.QueryLog{
		ProjectAlias: req.ProjectAlias,
		SQL:          req.SQL,
		Status:       "executing",
		CreateTime:   startTime,
		CreateUser:   userID,
	}

	countSql := fmt.Sprintf("SELECT COUNT(1) AS count FROM (%s) AS subquery", req.SQL)
	x.addDefaultLimit(req)

	// 获取项目的 Doris 连接
	dorisDB, err := core.GetProjectDoris(req.ProjectAlias)
	if err != nil {
		// 记录错误日志
		queryLog.Status = "error"
		queryLog.ErrorMsg = fmt.Sprintf("Failed to get Doris connection: %v", err)
		queryLog.Duration = time.Since(startTime).Seconds()

		if logErr := x.saveQueryLog(queryLog); logErr != nil {
			common.Logger.Error("Failed to save query log", zap.Error(logErr))
		}

		common.Logger.Error("Failed to get Doris connection for project",
			zap.String("project", req.ProjectAlias),
			zap.Error(err))
		return nil, fmt.Errorf("project %s not found or not available", req.ProjectAlias)
	}

	// 执行SQL查询
	tx := dorisDB.WithContext(ctx).Raw(req.SQL)
	if req.QueryId == "" {
		req.QueryId = uuid.New().String()
	}
	x.setQueryTraceId(req.QueryId, dorisDB)

	sqlRows, err := tx.Rows()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			x.cancelQuery(req.QueryId, dorisDB)
		}
		// 记录错误日志
		queryLog.Status = "error"
		queryLog.ErrorMsg = err.Error()
		queryLog.Duration = time.Since(startTime).Seconds()
		if logErr := x.saveQueryLog(queryLog); logErr != nil {
			common.Logger.Error("Failed to save query log", zap.Error(logErr))
		}

		common.Logger.Errorf("项目%s 执行失败, %v", req.ProjectAlias, err)
		return nil, fmt.Errorf("query execution failed: %v", err)
	}
	defer sqlRows.Close()
	// 解析结果
	columns, _ := sqlRows.Columns()
	common.Logger.Infof("Columns: %v", columns)
	var results []map[string]interface{}
	for sqlRows.Next() {
		tx.ScanRows(sqlRows, &results)
	}
	results = common.FormatDataForDisplay(results)
	// 计算执行时间
	duration := time.Since(startTime).Seconds()

	// 记录成功日志
	queryLog.Status = "success"
	queryLog.RowCount = len(results)
	queryLog.Duration = float64(duration) / 1000

	if err := x.saveQueryLog(queryLog); err != nil {
		common.Logger.Error("Failed to save query log", zap.Error(err))
		// 不影响查询结果返回
	}
	// 执行 count 语句
	var count int
	if err := dorisDB.Raw(countSql).Scan(&count).Error; err != nil {
		common.Logger.Error("Failed to execute count query", zap.Error(err))
		// 不影响查询结果返回
	}

	common.Logger.Infof("项目%s 返回: %d 行数据, 执行时间: %v 秒",
		req.ProjectAlias,
		len(results),
		duration)

	return &model.QueryResponse{
		Rows:     results,
		Columns:  columns,
		Count:    count,
		Sql:      req.SQL,
		Duration: duration,
	}, nil
}

func (s *SqlService) addDefaultLimit(req *model.QueryRequest) {
	// 如果没有设置 limit，则默认设置为 100
	if req.PageSize <= 0 {
		req.PageSize = 1000
	}
	// 如果没有设置 offset，则默认设置为 0
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if !strings.Contains(strings.ToLower(req.SQL), "limit") {
		// 添加 limit 和 offset
		offset := (req.PageNum - 1) * req.PageSize
		req.SQL = fmt.Sprintf("%s LIMIT %d OFFSET %d", req.SQL, req.PageSize, offset)
	}
}

// saveQueryLog 保存查询日志到MySQL
func (s *SqlService) saveQueryLog(log *model.QueryLog) error {
	centerDb := core.GetProjectCenter()
	if centerDb == nil {
		return fmt.Errorf("centerDb connection not available")
	}

	if err := centerDb.Create(log).Error; err != nil {
		return fmt.Errorf("failed to save query log: %w", err)
	}

	return nil
}

// GetQueryLogs 获取查询日志
func (s *SqlService) GetQueryLogs(projectAlias string, limit int, offset int) ([]model.QueryLog, error) {
	centerDb := core.GetProjectCenter()
	if centerDb == nil {
		return nil, fmt.Errorf("MySQL connection not available")
	}

	var logs []model.QueryLog
	query := centerDb.Model(&model.QueryLog{}).Order("create_time DESC")

	if projectAlias != "" {
		query = query.Where("project_alias = ?", projectAlias)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to get query logs: %w", err)
	}

	return logs, nil
}

func (x *SqlService) cancelQuery(queryID string, db *gorm.DB) error {
	cancelSql := fmt.Sprintf(`KILL QUERY "%s"`, queryID)
	if err := db.Exec(cancelSql).Error; err != nil {
		common.Logger.Errorf("%s 执行失败", queryID)
		return err
	}
	return nil
}

func (x *SqlService) setQueryTraceId(queryId string, db *gorm.DB) error {
	sql := fmt.Sprintf(`set session_context="trace_id:%s"`, queryId)
	err := db.Exec(sql).Error
	if err != nil {
		common.Logger.Errorf("setQueryTraceId err: %v", err.Error())
		return err
	}
	return nil
}

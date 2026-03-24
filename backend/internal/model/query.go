package model

import "time"

type QueryRequest struct {
	ProjectAlias string `json:"project_alias"`
	SQL          string `json:"sql"`
	TimeZone     string `json:"time_zone"`
	PageSize     int    `json:"page_size"`
	PageNum      int    `json:"page_num"`
	QueryId      string `json:"query_id"`
	UserID       string `json:"user_id"`
}

type QueryResponse struct {
	Rows     []map[string]interface{} `json:"rows"`
	Columns  []string                 `json:"columns"`
	Count    int                      `json:"count"`
	Sql      string                   `json:"sql"`
	Duration float64                  `json:"duration"`
}

type QueryLog struct {
	ProjectAlias string    `gorm:"column:project_alias"              json:"project_alias"`
	SQL          string    `gorm:"column:sql;type:text"              json:"sql"`
	Status       string    `gorm:"column:status"                     json:"status"`
	ErrorMsg     string    `gorm:"column:error_msg;type:text"        json:"error_msg"`
	RowCount     int       `gorm:"column:row_count"                  json:"row_count"`
	Duration     float64   `gorm:"column:duration"                   json:"duration"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	CreateUser   int       `gorm:"column:create_user"                json:"create_user"`
}

func (QueryLog) TableName() string {
	return "query_logs"
}

package model

import "time"

// ProjectMetric 指标配置表模型
type ProjectMetric struct {
	ProjectAlias string    `gorm:"column:project_alias;primaryKey"   json:"project_alias"`
	MetricCode   string    `gorm:"column:metric_code;primaryKey"     json:"metric_code"`
	MetricName   string    `gorm:"column:metric_name"                json:"metric_name"`
	Expression   string    `gorm:"column:expression"                 json:"expression"`
	BaseTable    string    `gorm:"column:base_table"                 json:"base_table"`
	Description  string    `gorm:"column:description"                json:"description"`
	Status       int       `gorm:"column:status"                     json:"status"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (ProjectMetric) TableName() string {
	return "project_metrics"
}

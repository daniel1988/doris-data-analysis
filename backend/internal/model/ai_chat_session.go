package model

import "time"

// AIChatSession AI会话记录表模型
type AIChatSession struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProjectAlias string    `gorm:"column:project_alias" json:"project_alias"`
	UserID       int64     `gorm:"column:user_id" json:"user_id"`
	UserQuery    string    `gorm:"column:user_query" json:"user_query"`
	LLMSQL       string    `gorm:"column:llm_sql" json:"llm_sql"`
	VizType      string    `gorm:"column:viz_type" json:"viz_type"`
	XAxis        string    `gorm:"column:x_axis" json:"x_axis"`
	YAxis        string    `gorm:"column:y_axis" json:"y_axis"`
	Narrative    string    `gorm:"column:narrative" json:"narrative"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (AIChatSession) TableName() string {
	return "ai_chat_sessions"
}

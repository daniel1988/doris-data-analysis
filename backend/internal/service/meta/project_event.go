package meta

import (
	"fmt"
	"time"

	"gitee.com/dmp_admin_v2/backend/internal/core"
)

// ProjectEvent 对应 dmp_center.project_event 表
type ProjectEvent struct {
	ProjectAlias string    `gorm:"column:project_alias;primaryKey"   json:"project_alias"`
	EventId      string    `gorm:"column:event_id;primaryKey"        json:"event_id"`
	EventName    string    `gorm:"column:event_name"                 json:"event_name"`
	EventType    int       `gorm:"column:event_type"                 json:"event_type"` // 0-普通, 1-其他
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

func (ProjectEvent) TableName() string {
	return "project_event"
}

type ProjectEventService struct{}

func NewProjectEventService() *ProjectEventService {
	return &ProjectEventService{}
}

// GetList 获取元事件列表
func (s *ProjectEventService) GetList(projectAlias string) ([]ProjectEvent, error) {
	var list []ProjectEvent
	db := core.GetProjectCenter()
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	err := db.Where("project_alias = ?", projectAlias).Find(&list).Error
	return list, err
}

// Create 创建元事件
func (s *ProjectEventService) Create(e *ProjectEvent) error {
	db := core.GetProjectCenter()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	return db.Create(e).Error
}

// Update 更新元事件
func (s *ProjectEventService) Update(e *ProjectEvent) error {
	db := core.GetProjectCenter()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	// 使用 Select("*") 确保所有字段（包括零值）都被更新
	return db.Model(e).Select("*").Updates(e).Error
}

// Delete 删除元事件
func (s *ProjectEventService) Delete(projectAlias, eventId string) error {
	db := core.GetProjectCenter()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	return db.Where("project_alias = ? AND event_id = ?", projectAlias, eventId).Delete(&ProjectEvent{}).Error
}

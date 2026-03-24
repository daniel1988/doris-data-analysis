package meta

import (
	"fmt"
	"time"

	"gitee.com/dmp_admin_v2/backend/internal/core"
	"gorm.io/gorm"
)

// ProjectEventProperty 对应 dmp_center.project_event_property 表
type ProjectEventProperty struct {
	ProjectAlias string    `gorm:"column:project_alias;primaryKey"   json:"project_alias"`
	EventId      string    `gorm:"column:event_id;primaryKey"        json:"event_id"`
	PropertyId   string    `gorm:"column:property_id;primaryKey"     json:"property_id"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

func (ProjectEventProperty) TableName() string {
	return "project_event_property"
}

type ProjectEventPropertyService struct{}

func NewProjectEventPropertyService() *ProjectEventPropertyService {
	return &ProjectEventPropertyService{}
}

// GetEventProperties 获取指定事件的关联属性
func (s *ProjectEventPropertyService) GetEventProperties(projectAlias, eventId string) ([]ProjectProperty, error) {
	var list []ProjectProperty
	db := core.GetProjectCenter()
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	// 联表查询 project_property 和 project_event_property
	err := db.Table("project_property").
		Select("project_property.*").
		Joins("JOIN project_event_property ON project_property.project_alias = project_event_property.project_alias AND project_property.property_id = project_event_property.property_id").
		Where("project_event_property.project_alias = ? AND project_event_property.event_id = ?", projectAlias, eventId).
		Find(&list).Error

	return list, err
}

// SaveRelations 保存事件与属性的关联关系
func (s *ProjectEventPropertyService) SaveRelations(projectAlias, eventId string, propertyIds []string) error {
	db := core.GetProjectCenter()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	return db.Transaction(func(tx *gorm.DB) error {
		// 先删除旧关系
		if err := tx.Where("project_alias = ? AND event_id = ?", projectAlias, eventId).Delete(&ProjectEventProperty{}).Error; err != nil {
			return err
		}

		// 插入新关系
		if len(propertyIds) > 0 {
			relations := make([]ProjectEventProperty, 0, len(propertyIds))
			for _, pid := range propertyIds {
				relations = append(relations, ProjectEventProperty{
					ProjectAlias: projectAlias,
					EventId:      eventId,
					PropertyId:   pid,
				})
			}
			if err := tx.Create(&relations).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

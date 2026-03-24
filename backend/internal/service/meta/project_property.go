package meta

import (
	"fmt"
	"time"

	"gitee.com/dmp_admin_v2/backend/internal/core"
)

// ProjectProperty 对应 dmp_center.project_property 表
type ProjectProperty struct {
	ProjectAlias string    `gorm:"column:project_alias;primaryKey"   json:"project_alias"`
	PropertyId   string    `gorm:"column:property_id;primaryKey"     json:"property_id"`
	PropertyName string    `gorm:"column:property_name"              json:"property_name"`
	DataType     string    `gorm:"column:data_type"                  json:"data_type"`
	PropertyType int       `gorm:"column:property_type"              json:"property_type"` // 0-普通, 1-其他
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

func (ProjectProperty) TableName() string {
	return "project_property"
}

type ProjectPropertyService struct{}

func NewProjectPropertyService() *ProjectPropertyService {
	return &ProjectPropertyService{}
}

// GetList 获取项目属性列表
func (s *ProjectPropertyService) GetList(projectAlias string) ([]ProjectProperty, error) {
	var list []ProjectProperty
	db := core.GetProjectCenter()
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	err := db.Where("project_alias = ?", projectAlias).Find(&list).Error
	return list, err
}

// Create 创建项目属性
func (s *ProjectPropertyService) Create(p *ProjectProperty) error {
	db := core.GetProjectCenter()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	return db.Create(p).Error
}

// Update 更新项目属性
func (s *ProjectPropertyService) Update(p *ProjectProperty) error {
	db := core.GetProjectCenter()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	// 使用 Select("*") 确保所有字段（包括零值）都被更新
	return db.Model(p).Select("*").Updates(p).Error
}

// Delete 删除项目属性
func (s *ProjectPropertyService) Delete(projectAlias, propertyId string) error {
	db := core.GetProjectCenter()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	return db.Where("project_alias = ? AND property_id = ?", projectAlias, propertyId).Delete(&ProjectProperty{}).Error
}

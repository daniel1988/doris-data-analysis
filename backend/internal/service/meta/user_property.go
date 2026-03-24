package meta

import (
	"fmt"
	"time"

	"gitee.com/dmp_admin_v2/backend/internal/core"
)

// UserProperty 对应 dmp_center.user_properties 表
type UserProperty struct {
	ProjectAlias string    `gorm:"column:project_alias;primaryKey"   json:"project_alias"`
	PropertyId   string    `gorm:"column:property_id;primaryKey"     json:"property_id"`
	PropertyName string    `gorm:"column:property_name"              json:"property_name"`
	EventType    string    `gorm:"column:event_type"                 json:"event_type"`
	DataType     string    `gorm:"column:data_type"                  json:"data_type"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

func (UserProperty) TableName() string {
	return "user_properties"
}

type UserPropertyService struct{}

func NewUserPropertyService() *UserPropertyService {
	return &UserPropertyService{}
}

// GetList 获取用户属性列表
func (s *UserPropertyService) GetList(projectAlias string) ([]UserProperty, error) {
	var list []UserProperty
	db := core.GetProjectCenter()
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	err := db.Where("project_alias = ?", projectAlias).Find(&list).Error
	return list, err
}

// Create 创建用户属性
func (s *UserPropertyService) Create(p *UserProperty) error {
	db := core.GetProjectCenter()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	return db.Create(p).Error
}

// Update 更新用户属性
func (s *UserPropertyService) Update(p *UserProperty) error {
	db := core.GetProjectCenter()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	return db.Model(p).Select("*").Updates(p).Error
}

// Delete 删除用户属性
func (s *UserPropertyService) Delete(projectAlias, propertyId string) error {
	db := core.GetProjectCenter()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	return db.Where("project_alias = ? AND property_id = ?", projectAlias, propertyId).Delete(&UserProperty{}).Error
}

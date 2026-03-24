package system

import (
	"time"

	"gitee.com/dmp_admin_v2/backend/internal/core"
)

// ProjectData 对应 dmp_center.project_data 表
type ProjectData struct {
	ProjectAlias string    `gorm:"column:project_alias;primaryKey" json:"project_alias"`
	ProjectName  string    `gorm:"column:project_name" json:"project_name"`
	Region       string    `gorm:"column:region" json:"region"`
	Secret       string    `gorm:"column:secret" json:"secret"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

func (ProjectData) TableName() string {
	return "project_data"
}

type ProjectDataService struct{}

func NewProjectDataService() *ProjectDataService {
	return &ProjectDataService{}
}

// GetList 获取项目列表
func (s *ProjectDataService) GetList() ([]ProjectData, error) {
	var list []ProjectData
	db := core.GetProjectCenter()
	err := db.Find(&list).Error
	return list, err
}

// Create 创建项目
func (s *ProjectDataService) Create(p *ProjectData) error {
	return core.GetProjectCenter().Create(p).Error
}

// Update 更新项目
func (s *ProjectDataService) Update(p *ProjectData) error {
	return core.GetProjectCenter().Model(p).Updates(p).Error
}

// Delete 删除项目
func (s *ProjectDataService) Delete(projectAlias string) error {
	return core.GetProjectCenter().Where("project_alias = ?", projectAlias).Delete(&ProjectData{}).Error
}

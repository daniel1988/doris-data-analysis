package analytics

import (
	"time"

	"gitee.com/dmp_admin_v2/backend/internal/core"
)

// Report 报表配置模型
type Report struct {
	ID           uint      `gorm:"primarykey"                        json:"id"`
	ProjectAlias string    `gorm:"column:project_alias"              json:"project_alias"`
	Name         string    `gorm:"column:name"                       json:"name"`
	Category     string    `gorm:"column:category"                   json:"category"` // event, funnel, retention
	Description  string    `gorm:"column:description"                json:"description"`
	QueryParams  string    `gorm:"column:query_params;type:text"     json:"query_params"` // JSON 存储
	CreateUser   int       `gorm:"column:create_user"                json:"create_user"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

func (Report) TableName() string {
	return "reports"
}

type ReportService struct{}

func NewReportService() *ReportService {
	return &ReportService{}
}

func (s *ReportService) Create(r *Report) error {
	return core.GetProjectCenter().Create(r).Error
}

func (s *ReportService) Update(r *Report) error {
	return core.GetProjectCenter().Model(r).Updates(r).Error
}

func (s *ReportService) GetList(projectAlias, category string) ([]Report, error) {
	var list []Report
	db := core.GetProjectCenter().Where("project_alias = ?", projectAlias)
	if category != "" {
		db = db.Where("category = ?", category)
	}
	err := db.Order("update_time DESC").Find(&list).Error
	return list, err
}

func (s *ReportService) Delete(id uint) error {
	return core.GetProjectCenter().Delete(&Report{}, id).Error
}

func (s *ReportService) GetByID(id uint) (*Report, error) {
	var list []Report
	err := core.GetProjectCenter().Where("id = ?", id).Find(&list).Error
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], err
}

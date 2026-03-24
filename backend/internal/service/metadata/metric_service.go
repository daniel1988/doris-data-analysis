package metadata

import (
	"gitee.com/dmp_admin_v2/backend/internal/core"
	"gitee.com/dmp_admin_v2/backend/internal/model"
)

type MetricService struct{}

func NewMetricService() *MetricService {
	return &MetricService{}
}

// GetList 获取项目下的指标列表
func (s *MetricService) GetList(projectAlias string) ([]model.ProjectMetric, error) {
	var list []model.ProjectMetric
	db := core.GetProjectCenter()
	err := db.Where("project_alias = ?", projectAlias).Order("metric_code ASC").Find(&list).Error
	return list, err
}

// GetEnabledList 获取启用的指标列表 (供 AI Schema 注入使用)
func (s *MetricService) GetEnabledList(projectAlias string) ([]model.ProjectMetric, error) {
	var list []model.ProjectMetric
	db := core.GetProjectCenter()
	err := db.Where("project_alias = ? AND status = 1", projectAlias).Order("metric_code ASC").Find(&list).Error
	return list, err
}

// Create 创建指标
func (s *MetricService) Create(m *model.ProjectMetric) error {
	db := core.GetProjectCenter()
	return db.Create(m).Error
}

// Update 更新指标
func (s *MetricService) Update(m *model.ProjectMetric) error {
	db := core.GetProjectCenter()
	// GORM Updates 默认忽略零值，如果需要更新 status 为 0，可能需要使用 Select 或者 map
	return db.Model(m).Where("project_alias = ? AND metric_code = ?", m.ProjectAlias, m.MetricCode).Updates(map[string]interface{}{
		"metric_name": m.MetricName,
		"expression":  m.Expression,
		"base_table":  m.BaseTable,
		"description": m.Description,
		"status":      m.Status,
	}).Error
}

// Delete 删除指标
func (s *MetricService) Delete(metricCode string, projectAlias string) error {
	db := core.GetProjectCenter()
	return db.Where("metric_code = ? AND project_alias = ?", metricCode, projectAlias).Delete(&model.ProjectMetric{}).Error
}

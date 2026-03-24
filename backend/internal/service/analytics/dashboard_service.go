package analytics

import (
	"time"

	"gitee.com/dmp_admin_v2/backend/internal/core"
)

// Dashboard 看板配置模型
type Dashboard struct {
	ID              uint            `gorm:"primarykey"                        json:"id"`
	ProjectAlias    string          `gorm:"column:project_alias"              json:"project_alias"`
	Name            string          `gorm:"column:name"                       json:"name"`
	DisplayName     string          `gorm:"column:display_name"               json:"display_name"`
	Description     string          `gorm:"column:description"                json:"description"`
	Category        string          `gorm:"column:category"                   json:"category"`
	LayoutType      string          `gorm:"column:layout_type"                json:"layout_type"`
	GridConfig      string          `gorm:"column:grid_config;type:text"      json:"grid_config"`
	Theme           string          `gorm:"column:theme"                      json:"theme"`
	RefreshInterval int             `gorm:"column:refresh_interval"           json:"refresh_interval"`
	Filters         string          `gorm:"column:filters;type:text"          json:"filters"`
	Variables       string          `gorm:"column:variables;type:text"        json:"variables"`
	Status          string          `gorm:"column:status"                     json:"status"`
	OwnerID         int             `gorm:"column:owner_id"                   json:"owner_id"`
	CreateTime      time.Time       `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime      time.Time       `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
	Items           []DashboardItem `gorm:"foreignKey:DashboardID"            json:"items"`
}

func (Dashboard) TableName() string {
	return "dashboards"
}

// DashboardItem 看板组件模型
type DashboardItem struct {
	ID             uint      `gorm:"primarykey"                        json:"id"`
	DashboardID    uint      `gorm:"column:dashboard_id"               json:"dashboard_id"`
	ReportID       uint      `gorm:"column:report_id"                  json:"report_id"`
	Type           string    `gorm:"column:type"                       json:"type"` // chart, table, mixed, stat, map
	Title          string    `gorm:"column:title"                      json:"title"`
	PositionX      int       `gorm:"column:position_x"                 json:"position_x"`
	PositionY      int       `gorm:"column:position_y"                 json:"position_y"`
	Width          int       `gorm:"column:width"                      json:"width"`
	Height         int       `gorm:"column:height"                     json:"height"`
	ZIndex         int       `gorm:"column:z_index"                    json:"z_index"`
	ConfigOverride string    `gorm:"column:config_override;type:text"  json:"config_override"`
	IsVisible      bool      `gorm:"column:is_visible;default:true"    json:"is_visible"`
	CreateTime     time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime     time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

func (DashboardItem) TableName() string {
	return "dashboard_items"
}

type DashboardService struct{}

func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

func (s *DashboardService) Create(d *Dashboard) error {
	return core.GetProjectCenter().Create(d).Error
}

func (s *DashboardService) Update(d *Dashboard) error {
	return core.GetProjectCenter().Model(d).Updates(d).Error
}

func (s *DashboardService) GetList(projectAlias string) ([]Dashboard, error) {
	var list []Dashboard
	err := core.GetProjectCenter().Where("project_alias = ?", projectAlias).
		Order("update_time DESC").Find(&list).Error
	return list, err
}

func (s *DashboardService) GetByID(id uint) (*Dashboard, error) {
	var list []Dashboard
	err := core.GetProjectCenter().Preload("Items").Where("id = ?", id).Find(&list).Error
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil // 或者返回自定义错误
	}
	return &list[0], nil
}

func (s *DashboardService) Delete(id uint) error {
	tx := core.GetProjectCenter().Begin()
	if err := tx.Where("dashboard_id = ?", id).Delete(&DashboardItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(&Dashboard{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// AddItem 为看板添加组件
func (s *DashboardService) AddItem(item *DashboardItem) error {
	return core.GetProjectCenter().Create(item).Error
}

// UpdateItem 更新组件配置或布局
func (s *DashboardService) UpdateItem(item *DashboardItem) error {
	return core.GetProjectCenter().Model(item).Updates(item).Error
}

// DeleteItem 删除组件
func (s *DashboardService) DeleteItem(id uint) error {
	return core.GetProjectCenter().Delete(&DashboardItem{}, id).Error
}

// BatchUpdateItems 批量更新看板组件（主要用于保存布局）
func (s *DashboardService) BatchUpdateItems(items []DashboardItem) error {
	tx := core.GetProjectCenter().Begin()
	for _, item := range items {
		// 只更新布局相关的字段，避免误操作
		if err := tx.Model(&DashboardItem{}).Where("id = ?", item.ID).Updates(map[string]interface{}{
			"position_x": item.PositionX,
			"position_y": item.PositionY,
			"width":      item.Width,
			"height":     item.Height,
		}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

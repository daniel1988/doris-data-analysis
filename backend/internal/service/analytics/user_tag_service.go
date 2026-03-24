package analytics

import (
	"time"

	"gitee.com/dmp_admin_v2/backend/internal/core"
)

// UserTag 用户标签/分群定义模型
type UserTag struct {
	ProjectAlias string    `gorm:"column:project_alias;primaryKey"   json:"project_alias"`
	TagCode      string    `gorm:"column:tag_code;primaryKey"        json:"tag_code"`
	TagName      string    `gorm:"column:tag_name"                   json:"tag_name"`
	TagSql       string    `gorm:"column:tag_sql;type:text"          json:"tag_sql"`
	UserCount    float64   `gorm:"column:user_count"                 json:"user_count"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

func (UserTag) TableName() string {
	return "user_tags"
}

type UserTagService struct{}

func NewUserTagService() *UserTagService {
	return &UserTagService{}
}

func (s *UserTagService) Create(t *UserTag) error {
	return core.GetProjectCenter().Create(t).Error
}

func (s *UserTagService) Update(t *UserTag) error {
	return core.GetProjectCenter().Model(t).Updates(t).Error
}

func (s *UserTagService) GetList(projectAlias string) ([]UserTag, error) {
	var list []UserTag
	err := core.GetProjectCenter().Where("project_alias = ?", projectAlias).
		Order("create_time DESC").Find(&list).Error
	return list, err
}

func (s *UserTagService) GetByCode(projectAlias, tagCode string) (*UserTag, error) {
	var list []UserTag
	err := core.GetProjectCenter().Where("project_alias = ? AND tag_code = ?", projectAlias, tagCode).Find(&list).Error
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], err
}

func (s *UserTagService) Delete(projectAlias, tagCode string) error {
	return core.GetProjectCenter().Where("project_alias = ? AND tag_code = ?", projectAlias, tagCode).Delete(&UserTag{}).Error
}

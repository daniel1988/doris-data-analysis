package analytics

import (
	"errors"

	"gitee.com/dmp_admin_v2/backend/internal/core"
	"gitee.com/dmp_admin_v2/backend/internal/model"
)

type AISessionService struct{}

func NewAISessionService() *AISessionService {
	return &AISessionService{}
}

// Create 保存一条新的 AI 会话记录
func (s *AISessionService) Create(session *model.AIChatSession) error {
	db := core.GetProjectCenter()
	return db.Create(session).Error
}

// GetList 获取用户在当前项目下的会话记录
func (s *AISessionService) GetList(projectAlias string, userID int64) ([]model.AIChatSession, error) {
	var list []model.AIChatSession
	db := core.GetProjectCenter()
	err := db.Where("project_alias = ? AND user_id = ?", projectAlias, userID).Order("id DESC").Find(&list).Error
	return list, err
}

// Delete 删除会话记录
func (s *AISessionService) Delete(id int64, projectAlias string, userID int64) error {
	db := core.GetProjectCenter()
	// 鉴权：只能删除属于自己且属于当前项目的记录
	return db.Where("id = ? AND project_alias = ? AND user_id = ?", id, projectAlias, userID).Delete(&model.AIChatSession{}).Error
}

// GetByID 根据ID获取会话记录
func (s *AISessionService) GetByID(id int64, projectAlias string) (*model.AIChatSession, error) {
	db := core.GetProjectCenter()
	var list []model.AIChatSession
	err := db.Where("id = ? AND project_alias = ?", id, projectAlias).Find(&list).Error
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil // not found
	}
	return &list[0], nil
}

// Execute 重新执行某个历史会话的 SQL
func (s *AISessionService) Execute(id int64, projectAlias string) (interface{}, error) {
	session, err := s.GetByID(id, projectAlias)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New("会话记录不存在")
	}

	// 重新执行 SQL
	chatService := NewAIChatService()
	results, err := chatService.ExecuteQuery(projectAlias, session.LLMSQL)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"sql":       session.LLMSQL,
		"viz_type":  session.VizType,
		"x_axis":    session.XAxis,
		"y_axis":    session.YAxis,
		"narrative": session.Narrative,
		"data":      results,
	}, nil
}

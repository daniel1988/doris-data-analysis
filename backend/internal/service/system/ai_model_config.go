package system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gitee.com/dmp_admin_v2/backend/internal/core"
)

// AIModelConfig 对应 dmp_center.ai_model_config 表
type AIModelConfig struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Provider    string    `gorm:"column:provider"                    json:"provider"`
	DisplayName string    `gorm:"column:display_name"                json:"display_name"`
	BaseURL     string    `gorm:"column:base_url"                    json:"base_url"`
	APIKey      string    `gorm:"column:api_key"                     json:"api_key"`
	ModelName   string    `gorm:"column:model_name"                  json:"model_name"`
	MaxTokens   int       `gorm:"column:max_tokens"                  json:"max_tokens"`
	Temperature float64   `gorm:"column:temperature"                 json:"temperature"`
	IsDefault   bool      `gorm:"column:is_default"                  json:"is_default"`
	IsEnabled   bool      `gorm:"column:is_enabled"                  json:"is_enabled"`
	SortOrder   int       `gorm:"column:sort_order"                  json:"sort_order"`
	CreateTime  time.Time `gorm:"column:create_time;autoCreateTime"  json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time;autoUpdateTime"  json:"update_time"`
}

func (AIModelConfig) TableName() string {
	return "ai_model_config"
}

// AIModelConfigBrief 前端可用模型列表的精简结构（不含 API Key）
type AIModelConfigBrief struct {
	ID          int64  `json:"id"`
	Provider    string `json:"provider"`
	DisplayName string `json:"display_name"`
	IsDefault   bool   `json:"is_default"`
}

type AIModelConfigService struct{}

func NewAIModelConfigService() *AIModelConfigService {
	return &AIModelConfigService{}
}

// GetList 获取所有模型配置（管理端，API Key 脱敏）
func (s *AIModelConfigService) GetList() ([]AIModelConfig, error) {
	var list []AIModelConfig
	db := core.GetProjectCenter()
	err := db.Order("sort_order ASC, id ASC").Find(&list).Error
	if err != nil {
		return nil, err
	}
	// API Key 脱敏
	for i := range list {
		list[i].APIKey = maskAPIKey(list[i].APIKey)
	}
	return list, nil
}

// GetEnabledModels 获取已启用的模型列表（前端选择器用）
func (s *AIModelConfigService) GetEnabledModels() ([]AIModelConfigBrief, error) {
	var list []AIModelConfigBrief
	db := core.GetProjectCenter()
	err := db.Model(&AIModelConfig{}).
		Select("id, provider, display_name, is_default").
		Where("is_enabled = ?", true).
		Order("sort_order ASC, id ASC").
		Find(&list).Error
	return list, err
}

// GetByID 根据 ID 获取模型配置（内部使用，包含完整 APIKey）
func (s *AIModelConfigService) GetByID(id int64) (*AIModelConfig, error) {
	var m AIModelConfig
	db := core.GetProjectCenter()
	// Doris 可能不支持 GORM 默认附加的 LIMIT 1，改用 Find 然后取第一条
	var list []AIModelConfig
	err := db.Where("id = ?", id).Find(&list).Error
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("record not found")
	}
	m = list[0]
	return &m, nil
}

// GetDefault 获取默认模型配置
func (s *AIModelConfigService) GetDefault() (*AIModelConfig, error) {
	var m AIModelConfig
	db := core.GetProjectCenter()
	// Doris 可能不支持 GORM 默认附加的 LIMIT 1，改用 Find 然后取第一条
	var list []AIModelConfig
	err := db.Where("is_default = ? AND is_enabled = ?", true, true).Find(&list).Error
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("record not found")
	}
	m = list[0]
	return &m, nil
}

// Create 创建模型配置
func (s *AIModelConfigService) Create(m *AIModelConfig) error {
	db := core.GetProjectCenter()
	// 如果设为默认，先取消其他默认
	if m.IsDefault {
		db.Model(&AIModelConfig{}).Where("1=1").Update("is_default", false)
	}
	return db.Create(m).Error
}

// Update 更新模型配置
func (s *AIModelConfigService) Update(m *AIModelConfig) error {
	db := core.GetProjectCenter()
	// 如果设为默认，先取消其他默认
	if m.IsDefault {
		db.Model(&AIModelConfig{}).Where("id != ?", m.ID).Update("is_default", false)
	}
	return db.Model(m).Updates(m).Error
}

// Delete 删除模型配置
func (s *AIModelConfigService) Delete(id int64) error {
	return core.GetProjectCenter().Where("id = ?", id).Delete(&AIModelConfig{}).Error
}

// TestConnection 测试模型连通性（向 LLM 发送简单请求验证 Key 和网络）
func (s *AIModelConfigService) TestConnection(id int64) (string, error) {
	m, err := s.GetByID(id)
	if err != nil {
		return "", fmt.Errorf("模型配置不存在: %v", err)
	}
	// 当提供商为 ollama 时，跳过鉴权相关的强制要求，或者允许使用任意字符串（例如 "ollama"）作为 API Key
	if m.Provider != "ollama" && m.APIKey == "" {
		return "", fmt.Errorf("API Key 未配置")
	}

	start := time.Now()

	reqBody := map[string]interface{}{
		"model": m.ModelName,
		"messages": []map[string]string{
			{"role": "user", "content": "Hi"},
		},
		"max_tokens": 5,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", m.BaseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("构建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if m.Provider != "ollama" {
		req.Header.Set("Authorization", "Bearer "+m.APIKey)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("网络连接失败: %v", err)
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return "", fmt.Errorf("API Key 无效 (HTTP %d)", resp.StatusCode)
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return "", fmt.Errorf("频率限制 (HTTP 429)")
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API 返回错误 (HTTP %d): %s", resp.StatusCode, string(bodyBytes))
	}

	return fmt.Sprintf("连接成功，模型: %s，延迟: %dms", m.ModelName, elapsed.Milliseconds()), nil
}

// maskAPIKey 脱敏 API Key，保留前4位和后4位
func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}

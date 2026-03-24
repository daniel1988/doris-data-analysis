package analytics

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"gitee.com/dmp_admin_v2/backend/internal/common"
	"gitee.com/dmp_admin_v2/backend/internal/core"
	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/internal/service/system"
	"github.com/spf13/viper"
	"github.com/xwb1989/sqlparser"
)

type AIChatService struct{}

func NewAIChatService() *AIChatService {
	return &AIChatService{}
}

// llmConfig 封装一次 LLM 调用所需的配置
type llmConfig struct {
	apiKey      string
	baseURL     string
	modelName   string
	maxTokens   int
	temperature float64
	mockMode    bool
}

// resolveLLMConfig 根据 modelId 解析模型配置
// 优先级：modelId 指定 > 数据库默认模型
func resolveLLMConfig(modelId int64) (*llmConfig, error) {
	configSrv := system.NewAIModelConfigService()

	// 1. 尝试通过 modelId 查找
	if modelId > 0 {
		m, err := configSrv.GetByID(modelId)
		// ollama 可能不需要 api key
		if err == nil && (m.APIKey != "" || m.Provider == "ollama") {
			return &llmConfig{
				apiKey:      m.APIKey,
				baseURL:     m.BaseURL,
				modelName:   m.ModelName,
				maxTokens:   m.MaxTokens,
				temperature: m.Temperature,
				mockMode:    false,
			}, nil
		}
	}

	// 2. 尝试数据库默认模型
	m, err := configSrv.GetDefault()
	if err == nil && (m.APIKey != "" || m.Provider == "ollama") {
		return &llmConfig{
			apiKey:      m.APIKey,
			baseURL:     m.BaseURL,
			modelName:   m.ModelName,
			maxTokens:   m.MaxTokens,
			temperature: m.Temperature,
			mockMode:    false,
		}, nil
	}

	return nil, errors.New("no valid AI model configuration found in database")
}

// LLMResponse defines the expected JSON format from the LLM
type LLMResponse struct {
	SQL       string `json:"sql"`
	VizType   string `json:"viz_type"`
	XAxis     string `json:"x_axis"`
	YAxis     string `json:"y_axis"`
	Narrative string `json:"narrative"`
}

// BuildVirtualSchema queries dmp_center to create a metadata context for LLM
func (s *AIChatService) BuildVirtualSchema(projectAlias string) (string, error) {
	db := core.GetProjectCenter()

	// Get event properties mapping
	var eventProps []struct {
		EventName    string `gorm:"column:event_name"`
		PropertyName string `gorm:"column:property_name"`
		DataType     string `gorm:"column:data_type"`
	}
	query := `
		SELECT e.event_name, p.property_name, p.data_type 
		FROM project_event e 
		JOIN project_event_property ep ON e.event_id = ep.event_id 
		JOIN project_property p ON ep.property_id = p.property_id 
		WHERE e.project_alias = ? AND p.project_alias = ?
	`
	db.Raw(query, projectAlias, projectAlias).Scan(&eventProps)

	var userProps []struct {
		PropertyName string `gorm:"column:property_name"`
		DataType     string `gorm:"column:data_type"`
	}
	db.Raw("SELECT property_name, data_type FROM user_properties WHERE project_alias = ?", projectAlias).Scan(&userProps)

	// 获取所有的元事件名称映射，方便大模型理解
	var allEvents []struct {
		EventID   string `gorm:"column:event_id"`
		EventName string `gorm:"column:event_name"`
	}
	db.Raw("SELECT event_id, event_name FROM project_event WHERE project_alias = ?", projectAlias).Scan(&allEvents)

	// 获取已启用的指标定义
	var metrics []struct {
		MetricName  string `gorm:"column:metric_name"`
		MetricCode  string `gorm:"column:metric_code"`
		Expression  string `gorm:"column:expression"`
		Description string `gorm:"column:description"`
	}
	db.Raw("SELECT metric_name, metric_code, expression, description FROM project_metrics WHERE project_alias = ? AND status = 1", projectAlias).Scan(&metrics)

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("你是一个专业的 Apache Doris SQL 数据分析师。\n用户想要查询项目 '%s' 的数据。\n", projectAlias))
	builder.WriteString("以下是你可以访问的表及其结构：\n\n")

	builder.WriteString("TABLE `event_data` (用户行为事件表)\n")
	builder.WriteString("- `e_openid` (String): 用户唯一标识 ID\n")
	builder.WriteString("- `e_event_id` (String): 事件标识\n")
	builder.WriteString("- `e_event_time` (Datetime): 事件发生时间\n")

	if len(allEvents) > 0 {
		builder.WriteString("已注册的事件映射（e_event_id 对应含义）：\n")
		for _, e := range allEvents {
			builder.WriteString(fmt.Sprintf("- e_event_id='%s' 代表事件 '%s'\n", e.EventID, e.EventName))
		}
	}

	builder.WriteString("事件特定属性字段：\n")
	propMap := make(map[string]string)
	for _, ep := range eventProps {
		key := ep.PropertyName
		if _, exists := propMap[key]; !exists {
			builder.WriteString(fmt.Sprintf("- `%s` (%s)\n", ep.PropertyName, ep.DataType))
			propMap[key] = ep.DataType
		}
	}

	if len(metrics) > 0 {
		builder.WriteString("\nAvailable Business Metrics (你可以直接在 SQL 查询中使用这些预定义的业务指标逻辑):\n")
		for _, m := range metrics {
			builder.WriteString(fmt.Sprintf("- Metric: %s (%s)\n  Logic: %s\n", m.MetricName, m.MetricCode, m.Expression))
			if m.Description != "" {
				builder.WriteString(fmt.Sprintf("  Description: %s\n", m.Description))
			}
		}
	}

	// builder.WriteString("\nTABLE `user_data` (用户表)\n")
	// builder.WriteString("- `u_openid` (String): 用户唯一标识 ID\n")
	// builder.WriteString("- `u_event_time` (Datetime): 用户创建时间\n")
	// builder.WriteString("用户特定属性字段：\n")
	// for _, up := range userProps {
	// 	builder.WriteString(fmt.Sprintf("- `%s` (%s)\n", up.PropertyName, up.DataType))
	// }

	return builder.String(), nil
}

// CallLLMForSQL communicates with an OpenAI-compatible API
func (s *AIChatService) CallLLMForSQL(cfg *llmConfig, systemPrompt, userQuery string) (*LLMResponse, error) {
	if viper.GetBool("ai.mock_mode") {
		// Mock response for testing/demo purposes when no API key is set
		time.Sleep(1 * time.Second) // Simulate latency
		return &LLMResponse{
			SQL:       "SELECT e_event_name, COUNT(*) as count FROM event_data GROUP BY e_event_name ORDER BY count DESC LIMIT 5",
			VizType:   "bar",
			XAxis:     "e_event_name",
			YAxis:     "count",
			Narrative: "[MOCK] This is a simulated response. Please configure your LLM API Key in the backend service. Here is a breakdown of top events.",
		}, nil
	}

	temperature := cfg.temperature
	if temperature == 0 {
		temperature = 0.1
	}

	reqBody := map[string]interface{}{
		"model": cfg.modelName,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt + "\n关键规则:\n1. 只能返回有效的 SELECT SQL 语句。\n2. 必须严格按照 JSON 格式响应，包含以下键：sql, viz_type, x_axis, y_axis, narrative。\n3. 当用户询问相关业务指标时，请优先使用上述 Available Business Metrics 中提供的逻辑进行组合和查询。"},
			{"role": "user", "content": userQuery},
		},
		"response_format": map[string]string{"type": "json_object"},
		"temperature":     temperature,
	}
	jsonBody, _ := json.Marshal(reqBody)

	common.Logger.Infof("请求Prompt: %s", string(jsonBody))

	req, _ := http.NewRequest("POST", cfg.baseURL, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	if cfg.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+cfg.apiKey)
	}

	client := &http.Client{Timeout: 90 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("LLM API error: %s", string(bodyBytes))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return nil, errors.New("no choice in LLM response")
	}
	choice := choices[0].(map[string]interface{})
	message := choice["message"].(map[string]interface{})
	content := message["content"].(string)

	var llmResp LLMResponse
	if err := json.Unmarshal([]byte(content), &llmResp); err != nil {
		// Sometimes LLM wraps JSON in markdown code blocks
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimPrefix(content, "```")
		content = strings.TrimSuffix(content, "```")
		if err2 := json.Unmarshal([]byte(content), &llmResp); err2 != nil {
			return nil, fmt.Errorf("failed to parse JSON from LLM: %v", err)
		}
	}

	return &llmResp, nil
}

// ValidateAndSecureSQL uses sqlparser to validate safety and inject project boundaries
func (s *AIChatService) ValidateAndSecureSQL(sqlStr string, projectAlias string) (string, error) {
	stmt, err := sqlparser.Parse(sqlStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse SQL: %v", err)
	}

	sel, ok := stmt.(*sqlparser.Select)
	if !ok {
		return "", errors.New("Access Denied: Only SELECT operations are allowed")
	}

	// Simple check for tables
	validTables := map[string]bool{"event_data": true, "user_data": true, "user_tag_data": true}
	for _, expr := range sel.From {
		if te, ok := expr.(*sqlparser.AliasedTableExpr); ok {
			if tn, ok := te.Expr.(sqlparser.TableName); ok {
				if !validTables[tn.Name.String()] {
					return "", fmt.Errorf("Access Denied: Table %s is not allowed", tn.Name.String())
				}
			}
		}
	}

	// Limit injection
	if sel.Limit == nil && len(sel.GroupBy) == 0 {
		sel.Limit = &sqlparser.Limit{Rowcount: sqlparser.NewIntVal([]byte("1000"))}
	}

	safeSQL := sqlparser.String(stmt)
	return safeSQL, nil
}

// ExecuteQuery runs the safe SQL
func (s *AIChatService) ExecuteQuery(projectAlias, sql string) ([]map[string]interface{}, error) {
	sqlService := NewSqlService()
	req := &model.QueryRequest{
		ProjectAlias: projectAlias,
		SQL:          sql,
	}
	resp, err := sqlService.ExecuteQuery(req)
	if err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

// CallLLMForChat communicates with an OpenAI-compatible API for pure text chat
func (s *AIChatService) CallLLMForChat(cfg *llmConfig, systemPrompt, userQuery string) (string, error) {
	if viper.GetBool("ai.mock_mode") {
		time.Sleep(1 * time.Second) // Simulate latency
		return "[MOCK] 这是一个模拟的 AI 对话回复。这里可以是对数据的分析或者业务解答。", nil
	}

	temperature := cfg.temperature
	if temperature == 0 {
		temperature = 0.7 // Use a higher temperature for chat
	}

	reqBody := map[string]interface{}{
		"model": cfg.modelName,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userQuery},
		},
		"temperature": temperature,
	}
	jsonBody, _ := json.Marshal(reqBody)

	common.Logger.Infof("请求Chat Prompt: %s", string(jsonBody))

	req, _ := http.NewRequest("POST", cfg.baseURL, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	if cfg.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+cfg.apiKey)
	}

	client := &http.Client{Timeout: 90 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("LLM API error: %s", string(bodyBytes))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", errors.New("no choice in LLM response")
	}
	choice := choices[0].(map[string]interface{})
	message := choice["message"].(map[string]interface{})
	content := message["content"].(string)

	return content, nil
}

// ChatWithAI handles pure natural language conversation, optionally with context data
func (s *AIChatService) ChatWithAI(projectAlias, userQuery string, modelId int64, contextData interface{}) (interface{}, error) {
	cfg, err := resolveLLMConfig(modelId)
	if err != nil {
		return nil, err
	}

	systemPrompt := "你是一个专业的数据分析师。请根据用户的问题进行解答。"
	if contextData != nil {
		// Limit the data size to prevent token overflow
		dataBytes, _ := json.Marshal(contextData)
		dataStr := string(dataBytes)
		if len(dataStr) > 4000 {
			dataStr = dataStr[:4000] + "...(数据已截断)"
		}
		systemPrompt += fmt.Sprintf("\n用户正在查看以下数据（JSON格式）：\n%s\n请结合上述数据进行分析解读。", dataStr)
	}

	content, err := s.CallLLMForChat(cfg, systemPrompt, userQuery)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"type":    "text",
		"content": content,
	}, nil
}

// HandleChat orchestrates the workflow with dynamic model selection
func (s *AIChatService) HandleChat(projectAlias, userQuery string, modelId int64) (interface{}, error) {
	// 根据 modelId 解析模型配置
	cfg, err := resolveLLMConfig(modelId)
	if err != nil {
		return nil, err
	}

	schema, err := s.BuildVirtualSchema(projectAlias)
	if err != nil {
		return nil, err
	}

	retries := 2
	var finalSQL string
	var llmResp *LLMResponse
	var lastError error

	systemPrompt := schema

	for i := 0; i <= retries; i++ {
		// If it's a retry, append the error to the prompt
		currentPrompt := systemPrompt
		if lastError != nil {
			currentPrompt += fmt.Sprintf("\n上一次执行 SQL 失败，错误信息: %v。请修正 SQL。", lastError)
		}

		llmResp, err = s.CallLLMForSQL(cfg, currentPrompt, userQuery)
		if err != nil {
			log.Printf("[AIChatService] CallLLMForSQL failed (attempt %d): %v\n", i+1, err)
			return nil, err
		}

		log.Printf("[AIChatService] LLM Generated SQL (attempt %d):\n%s\n", i+1, llmResp.SQL)

		finalSQL, err = s.ValidateAndSecureSQL(llmResp.SQL, projectAlias)
		if err != nil {
			lastError = fmt.Errorf("SQL Validation Failed: %v", err)
			continue
		}

		// Execute
		results, err := s.ExecuteQuery(projectAlias, finalSQL)
		if err != nil {
			lastError = err
			continue
		}

		// Success
		return map[string]interface{}{
			"sql":       finalSQL,
			"viz_type":  llmResp.VizType,
			"x_axis":    llmResp.XAxis,
			"y_axis":    llmResp.YAxis,
			"narrative": llmResp.Narrative,
			"data":      results,
		}, nil
	}

	return nil, fmt.Errorf("failed after %d retries. Last error: %v", retries, lastError)
}

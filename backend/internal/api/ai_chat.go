package api

import (
	"net/http"
	"strconv"

	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/internal/service/analytics"
	"github.com/gin-gonic/gin"
)

type AIChatHandler struct {
	aiChatService    *analytics.AIChatService
	aiSessionService *analytics.AISessionService
}

func NewAIChatHandler() *AIChatHandler {
	return &AIChatHandler{
		aiChatService:    analytics.NewAIChatService(),
		aiSessionService: analytics.NewAISessionService(),
	}
}

// HandleChat processes natural language to return SQL and analytical results
func (h *AIChatHandler) HandleChat(c *gin.Context) {
	var req struct {
		ProjectAlias string      `json:"project_alias" binding:"required"`
		Query        string      `json:"query"         binding:"required"`
		ModelID      int64       `json:"model_id"`
		Mode         string      `json:"mode"`
		ContextData  interface{} `json:"context_data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Mode == "chat" {
		result, err := h.aiChatService.ChatWithAI(req.ProjectAlias, req.Query, req.ModelID, req.ContextData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": result})
		return
	}

	result, err := h.aiChatService.HandleChat(req.ProjectAlias, req.Query, req.ModelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": result})
}

// SaveSession 保存 AI 会话记录
func (h *AIChatHandler) SaveSession(c *gin.Context) {
	var req model.AIChatSession
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "msg": "参数绑定失败"})
		return
	}

	projectAlias := c.GetHeader("X-Project-Alias")
	if projectAlias == "" {
		projectAlias = c.Query("project_alias")
	}
	if projectAlias == "" {
		// 如果 Header 和 Query 都没有，看 JSON 里有没有传
		if req.ProjectAlias != "" {
			projectAlias = req.ProjectAlias
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing project_alias"})
			return
		}
	}
	req.ProjectAlias = projectAlias

	// TODO: 从上下文获取当前登录用户 ID，暂用默认 1 代替
	req.UserID = 1

	if err := h.aiSessionService.Create(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": req})
}

// GetSessions 获取当前项目的 AI 会话记录列表
func (h *AIChatHandler) GetSessions(c *gin.Context) {
	projectAlias := c.GetHeader("X-Project-Alias")
	if projectAlias == "" {
		projectAlias = c.Query("project_alias")
	}
	if projectAlias == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing project_alias"})
		return
	}

	// TODO: 从上下文获取当前登录用户 ID，暂用默认 1 代替
	userID := int64(1)

	list, err := h.aiSessionService.GetList(projectAlias, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": list})
}

// DeleteSession 删除 AI 会话记录
func (h *AIChatHandler) DeleteSession(c *gin.Context) {
	projectAlias := c.GetHeader("X-Project-Alias")
	if projectAlias == "" {
		projectAlias = c.Query("project_alias")
	}
	if projectAlias == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing project_alias"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// TODO: 从上下文获取当前登录用户 ID，暂用默认 1 代替
	userID := int64(1)

	if err := h.aiSessionService.Delete(id, projectAlias, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": nil})
}

// ExecuteSession 重新执行 AI 会话记录
func (h *AIChatHandler) ExecuteSession(c *gin.Context) {
	projectAlias := c.GetHeader("X-Project-Alias")
	if projectAlias == "" {
		projectAlias = c.Query("project_alias")
	}
	if projectAlias == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing project_alias"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	result, err := h.aiSessionService.Execute(id, projectAlias)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": result})
}

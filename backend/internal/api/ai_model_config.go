package api

import (
	"strconv"

	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/internal/service/system"

	"github.com/gin-gonic/gin"
)

type AIModelConfigHandler struct {
	srv *system.AIModelConfigService
}

func NewAIModelConfigHandler() *AIModelConfigHandler {
	return &AIModelConfigHandler{
		srv: system.NewAIModelConfigService(),
	}
}

// GetList 获取所有模型配置（管理端）
func (h *AIModelConfigHandler) GetList(c *gin.Context) {
	list, err := h.srv.GetList()
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}
	model.Success(c, list)
}

// GetEnabledModels 获取已启用的模型列表（前端选择器）
func (h *AIModelConfigHandler) GetEnabledModels(c *gin.Context) {
	list, err := h.srv.GetEnabledModels()
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}
	model.Success(c, list)
}

// Create 创建模型配置
func (h *AIModelConfigHandler) Create(c *gin.Context) {
	var m system.AIModelConfig
	if err := c.ShouldBindJSON(&m); err != nil {
		model.BadRequest(c, err.Error())
		return
	}
	if err := h.srv.Create(&m); err != nil {
		model.InternalError(c, err.Error())
		return
	}
	model.Success(c, nil)
}

// Update 更新模型配置
func (h *AIModelConfigHandler) Update(c *gin.Context) {
	var m system.AIModelConfig
	if err := c.ShouldBindJSON(&m); err != nil {
		model.BadRequest(c, err.Error())
		return
	}
	if m.ID == 0 {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			model.BadRequest(c, "invalid id")
			return
		}
		m.ID = id
	}
	if err := h.srv.Update(&m); err != nil {
		model.InternalError(c, err.Error())
		return
	}
	model.Success(c, nil)
}

// Delete 删除模型配置
func (h *AIModelConfigHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		model.BadRequest(c, "invalid id")
		return
	}
	if err := h.srv.Delete(id); err != nil {
		model.InternalError(c, err.Error())
		return
	}
	model.Success(c, nil)
}

// TestConnection 测试模型连通性
func (h *AIModelConfigHandler) TestConnection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		model.BadRequest(c, "invalid id")
		return
	}
	msg, err := h.srv.TestConnection(id)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}
	model.Success(c, gin.H{"message": msg})
}

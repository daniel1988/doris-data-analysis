package api

import (
	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/internal/service/meta"

	"github.com/gin-gonic/gin"
)

type ProjectEventHandler struct {
	srv *meta.ProjectEventService
}

func NewProjectEventHandler() *ProjectEventHandler {
	return &ProjectEventHandler{
		srv: meta.NewProjectEventService(),
	}
}

// GetList 获取元事件列表
func (h *ProjectEventHandler) GetList(c *gin.Context) {
	projectAlias := c.Query("project_alias")
	if projectAlias == "" {
		model.BadRequest(c, "project_alias is required")
		return
	}

	list, err := h.srv.GetList(projectAlias)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, list)
}

// Create 创建元事件
func (h *ProjectEventHandler) Create(c *gin.Context) {
	var e meta.ProjectEvent
	if err := c.ShouldBindJSON(&e); err != nil {
		model.BadRequest(c, err.Error())
		return
	}

	if err := h.srv.Create(&e); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, nil)
}

// Update 更新元事件
func (h *ProjectEventHandler) Update(c *gin.Context) {
	var e meta.ProjectEvent
	if err := c.ShouldBindJSON(&e); err != nil {
		model.BadRequest(c, err.Error())
		return
	}

	if err := h.srv.Update(&e); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, nil)
}

// Delete 删除元事件
func (h *ProjectEventHandler) Delete(c *gin.Context) {
	projectAlias := c.Query("project_alias")
	eventId := c.Param("id")

	if projectAlias == "" || eventId == "" {
		model.BadRequest(c, "project_alias and event_id are required")
		return
	}

	if err := h.srv.Delete(projectAlias, eventId); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, nil)
}

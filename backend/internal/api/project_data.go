package api

import (
	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/internal/service/system"

	"github.com/gin-gonic/gin"
)

type ProjectDataHandler struct {
	srv *system.ProjectDataService
}

func NewProjectDataHandler() *ProjectDataHandler {
	return &ProjectDataHandler{
		srv: system.NewProjectDataService(),
	}
}

// GetList 获取项目列表
func (h *ProjectDataHandler) GetList(c *gin.Context) {
	list, err := h.srv.GetList()
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, list)
}

// Create 创建项目
func (h *ProjectDataHandler) Create(c *gin.Context) {
	var p system.ProjectData
	if err := c.ShouldBindJSON(&p); err != nil {
		model.BadRequest(c, err.Error())
		return
	}

	if err := h.srv.Create(&p); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, nil)
}

// Update 更新项目
func (h *ProjectDataHandler) Update(c *gin.Context) {
	var p system.ProjectData
	if err := c.ShouldBindJSON(&p); err != nil {
		model.BadRequest(c, err.Error())
		return
	}

	if err := h.srv.Update(&p); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, nil)
}

// Delete 删除项目
func (h *ProjectDataHandler) Delete(c *gin.Context) {
	projectAlias := c.Param("alias")
	if projectAlias == "" {
		model.BadRequest(c, "project_alias is required")
		return
	}

	if err := h.srv.Delete(projectAlias); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, nil)
}

package api

import (
	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/internal/service/meta"

	"github.com/gin-gonic/gin"
)

type ProjectPropertyHandler struct {
	srv *meta.ProjectPropertyService
}

func NewProjectPropertyHandler() *ProjectPropertyHandler {
	return &ProjectPropertyHandler{
		srv: meta.NewProjectPropertyService(),
	}
}

// GetList 获取项目属性列表
func (h *ProjectPropertyHandler) GetList(c *gin.Context) {
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

// Create 创建项目属性
func (h *ProjectPropertyHandler) Create(c *gin.Context) {
	var p meta.ProjectProperty
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

// Update 更新项目属性
func (h *ProjectPropertyHandler) Update(c *gin.Context) {
	var p meta.ProjectProperty
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

// Delete 删除项目属性
func (h *ProjectPropertyHandler) Delete(c *gin.Context) {
	projectAlias := c.Query("project_alias")
	propertyId := c.Param("id")

	if projectAlias == "" || propertyId == "" {
		model.BadRequest(c, "project_alias and property_id are required")
		return
	}

	if err := h.srv.Delete(projectAlias, propertyId); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, nil)
}

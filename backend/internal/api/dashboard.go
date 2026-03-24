package api

import (
	"strconv"

	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/internal/service/analytics"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	srv *analytics.DashboardService
}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{
		srv: analytics.NewDashboardService(),
	}
}

func (h *DashboardHandler) Create(c *gin.Context) {
	var d analytics.Dashboard
	if err := c.ShouldBindJSON(&d); err != nil {
		model.BadRequest(c, err.Error())
		return
	}

	if err := h.srv.Create(&d); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, d)
}

func (h *DashboardHandler) Update(c *gin.Context) {
	var d analytics.Dashboard
	if err := c.ShouldBindJSON(&d); err != nil {
		model.BadRequest(c, err.Error())
		return
	}

	if err := h.srv.Update(&d); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, nil)
}

func (h *DashboardHandler) GetList(c *gin.Context) {
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

func (h *DashboardHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	d, err := h.srv.GetByID(uint(id))
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, d)
}

func (h *DashboardHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := h.srv.Delete(uint(id)); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, nil)
}

func (h *DashboardHandler) AddItem(c *gin.Context) {
	var item analytics.DashboardItem
	if err := c.ShouldBindJSON(&item); err != nil {
		model.BadRequest(c, err.Error())
		return
	}

	if err := h.srv.AddItem(&item); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, item)
}

func (h *DashboardHandler) UpdateItem(c *gin.Context) {
	var item analytics.DashboardItem
	if err := c.ShouldBindJSON(&item); err != nil {
		model.BadRequest(c, err.Error())
		return
	}

	if err := h.srv.UpdateItem(&item); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, nil)
}

func (h *DashboardHandler) DeleteItem(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := h.srv.DeleteItem(uint(id)); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, nil)
}

func (h *DashboardHandler) BatchUpdateItems(c *gin.Context) {
	var items []analytics.DashboardItem
	if err := c.ShouldBindJSON(&items); err != nil {
		model.BadRequest(c, err.Error())
		return
	}

	if err := h.srv.BatchUpdateItems(items); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, nil)
}

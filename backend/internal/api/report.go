package api

import (
	"strconv"

	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/internal/service/analytics"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	srv *analytics.ReportService
}

func NewReportHandler() *ReportHandler {
	return &ReportHandler{
		srv: analytics.NewReportService(),
	}
}

func (h *ReportHandler) Create(c *gin.Context) {
	var r analytics.Report
	if err := c.ShouldBindJSON(&r); err != nil {
		model.BadRequest(c, err.Error())
		return
	}

	if err := h.srv.Create(&r); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, r)
}

func (h *ReportHandler) Update(c *gin.Context) {
	var r analytics.Report
	if err := c.ShouldBindJSON(&r); err != nil {
		model.BadRequest(c, err.Error())
		return
	}

	if err := h.srv.Update(&r); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, nil)
}

func (h *ReportHandler) GetList(c *gin.Context) {
	projectAlias := c.Query("project_alias")
	category := c.Query("category")

	if projectAlias == "" {
		model.BadRequest(c, "project_alias is required")
		return
	}

	list, err := h.srv.GetList(projectAlias, category)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, list)
}

func (h *ReportHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := h.srv.Delete(uint(id)); err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, nil)
}

func (h *ReportHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	r, err := h.srv.GetByID(uint(id))
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, r)
}

package api

import (
	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/internal/service/metadata"
	"github.com/gin-gonic/gin"
)

type ProjectMetricHandler struct {
	metricService *metadata.MetricService
}

func NewProjectMetricHandler() *ProjectMetricHandler {
	return &ProjectMetricHandler{
		metricService: metadata.NewMetricService(),
	}
}

func (h *ProjectMetricHandler) GetList(c *gin.Context) {
	projectAlias := c.GetHeader("X-Project-Alias")
	if projectAlias == "" {
		projectAlias = c.Query("project_alias")
	}
	if projectAlias == "" {
		model.BadRequest(c, "missing project_alias")
		return
	}

	list, err := h.metricService.GetList(projectAlias)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}
	model.Success(c, list)
}

func (h *ProjectMetricHandler) Create(c *gin.Context) {
	projectAlias := c.GetHeader("X-Project-Alias")
	if projectAlias == "" {
		projectAlias = c.Query("project_alias")
	}
	if projectAlias == "" {
		model.BadRequest(c, "missing project_alias")
		return
	}

	var m model.ProjectMetric
	if err := c.ShouldBindJSON(&m); err != nil {
		model.BadRequest(c, err.Error())
		return
	}
	m.ProjectAlias = projectAlias

	if err := h.metricService.Create(&m); err != nil {
		model.InternalError(c, err.Error())
		return
	}
	model.Success(c, m)
}

func (h *ProjectMetricHandler) Update(c *gin.Context) {
	projectAlias := c.GetHeader("X-Project-Alias")
	if projectAlias == "" {
		projectAlias = c.Query("project_alias")
	}
	if projectAlias == "" {
		model.BadRequest(c, "missing project_alias")
		return
	}

	metricCode := c.Param("id")
	if metricCode == "" {
		model.BadRequest(c, "invalid metric_code")
		return
	}

	var m model.ProjectMetric
	if err := c.ShouldBindJSON(&m); err != nil {
		model.BadRequest(c, err.Error())
		return
	}
	m.MetricCode = metricCode
	m.ProjectAlias = projectAlias

	if err := h.metricService.Update(&m); err != nil {
		model.InternalError(c, err.Error())
		return
	}
	model.Success(c, m)
}

func (h *ProjectMetricHandler) Delete(c *gin.Context) {
	projectAlias := c.GetHeader("X-Project-Alias")
	if projectAlias == "" {
		projectAlias = c.Query("project_alias")
	}
	if projectAlias == "" {
		model.BadRequest(c, "missing project_alias")
		return
	}

	metricCode := c.Param("id")
	if metricCode == "" {
		model.BadRequest(c, "invalid metric_code")
		return
	}

	if err := h.metricService.Delete(metricCode, projectAlias); err != nil {
		model.InternalError(c, err.Error())
		return
	}
	model.Success(c, nil)
}

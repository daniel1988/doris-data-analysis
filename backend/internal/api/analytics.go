package api

import (
	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/internal/service/analytics"
	"gitee.com/dmp_admin_v2/backend/internal/service/meta"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	userPropSrv *meta.UserPropertyService
	userTagSrv  *analytics.UserTagService
}

func NewAnalyticsHandler() *AnalyticsHandler {
	return &AnalyticsHandler{
		userPropSrv: meta.NewUserPropertyService(),
		userTagSrv:  analytics.NewUserTagService(),
	}
}

func (h *AnalyticsHandler) EventAnalysis(c *gin.Context) {
	var req analytics.EventAnalysisReq
	if err := c.ShouldBindJSON(&req); err != nil {
		model.BadRequest(c, err.Error())
		return
	}

	resp, err := analytics.EventAnalysis(&req)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, resp)
}

func (h *AnalyticsHandler) GetUserProperties(c *gin.Context) {
	projectAlias := c.Query("project_alias")
	if projectAlias == "" {
		model.BadRequest(c, "project_alias is required")
		return
	}

	list, err := h.userPropSrv.GetList(projectAlias)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, list)
}

func (h *AnalyticsHandler) GetUserTags(c *gin.Context) {
	projectAlias := c.Query("project_alias")
	if projectAlias == "" {
		model.BadRequest(c, "project_alias is required")
		return
	}

	list, err := h.userTagSrv.GetList(projectAlias)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, list)
}

func (h *AnalyticsHandler) GetUserGroups(c *gin.Context) {
	// 目前分群数据存放在 user_tags 表中，此处暂时复用 GetUserTags 的逻辑
	h.GetUserTags(c)
}

func (h *AnalyticsHandler) GetTagValues(c *gin.Context) {
	projectAlias := c.Query("project_alias")
	tagCode := c.Query("tag_code")

	if projectAlias == "" || tagCode == "" {
		model.BadRequest(c, "project_alias and tag_code are required")
		return
	}

	// 模拟实现：从维度值查询逻辑中获取标签值
	// 实际上 user_tag_data 表存储了 user_id 和 tag_value
	// 这里我们可以复用 GetDimensions 的部分逻辑，或者直接查询 user_tag_data
	resp, err := analytics.GetDimensions(&analytics.DimensionReq{
		ProjectAlias: projectAlias,
		Table:        "user_tag_data",
		Field:        "tag_value",
		// 实际上需要过滤 tag_code，但目前的 DimensionReq 不支持多条件过滤
		// 暂且返回该表的所有值，或者后续优化 DimensionReq
	})

	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, resp)
}

func (h *AnalyticsHandler) FunnelAnalysis(c *gin.Context) {
	var req analytics.FunnelAnalysisReq
	if err := c.ShouldBindJSON(&req); err != nil {
		model.BadRequest(c, err.Error())
		return
	}

	resp, err := analytics.FunnelAnalysis(&req)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, resp)
}

func (h *AnalyticsHandler) RetentionAnalysis(c *gin.Context) {
	var req analytics.RetentionAnalysisReq
	if err := c.ShouldBindJSON(&req); err != nil {
		model.BadRequest(c, err.Error())
		return
	}

	resp, err := analytics.RetentionAnalysis(&req)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, resp)
}

func (h *AnalyticsHandler) ScatterAnalysis(c *gin.Context) {
	var req analytics.ScatterAnalysisReq
	if err := c.ShouldBindJSON(&req); err != nil {
		model.BadRequest(c, err.Error())
		return
	}

	resp, err := analytics.ScatterAnalysis(&req)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, resp)
}

func (h *AnalyticsHandler) UserPropertyAnalysis(c *gin.Context) {
	var req analytics.UserPropertyAnalysisReq
	if err := c.ShouldBindJSON(&req); err != nil {
		model.BadRequest(c, err.Error())
		return
	}

	resp, err := analytics.UserPropertyAnalysis(&req)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, resp)
}

func (h *AnalyticsHandler) GetUserList(c *gin.Context) {
	var req analytics.UserListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		model.BadRequest(c, err.Error())
		return
	}

	resp, err := analytics.GetUserList(&req)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, resp)
}

func (h *AnalyticsHandler) GetDimensions(c *gin.Context) {
	var req analytics.DimensionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		model.BadRequest(c, err.Error())
		return
	}

	resp, err := analytics.GetDimensions(&req)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, resp)
}

func (h *AnalyticsHandler) EventDetail(c *gin.Context) {
	var req analytics.EventDetailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		model.BadRequest(c, err.Error())
		return
	}

	resp, err := analytics.EventDetail(&req)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	model.Success(c, resp)
}

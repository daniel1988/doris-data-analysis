package api

import (
	"gitee.com/dmp_admin_v2/backend/internal/model"
	"gitee.com/dmp_admin_v2/backend/internal/service/meta"

	"github.com/gin-gonic/gin"
)

type SelectorHandler struct {
	evSrv   *meta.ProjectEventService
	propSrv *meta.ProjectPropertyService
	relSrv  *meta.ProjectEventPropertyService
}

func NewSelectorHandler() *SelectorHandler {
	return &SelectorHandler{
		evSrv:   meta.NewProjectEventService(),
		propSrv: meta.NewProjectPropertyService(),
		relSrv:  meta.NewProjectEventPropertyService(),
	}
}

// GetEvents 获取当前项目的事件下拉列表
func (h *SelectorHandler) GetEvents(c *gin.Context) {
	projectAlias := c.Query("project_alias")
	if projectAlias == "" {
		model.BadRequest(c, "project_alias is required")
		return
	}

	list, err := h.evSrv.GetList(projectAlias)
	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	// 转换为通用的下拉结构 (id, name)
	type Option struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	options := make([]Option, 0, len(list))
	for _, e := range list {
		options = append(options, Option{
			Id:   e.EventId,
			Name: e.EventName,
		})
	}

	model.Success(c, options)
}

// GetProperties 获取当前项目的事件属性下拉列表
func (h *SelectorHandler) GetProperties(c *gin.Context) {
	projectAlias := c.Query("project_alias")
	eventId := c.Query("e_event_id")

	if projectAlias == "" {
		model.BadRequest(c, "project_alias is required")
		return
	}

	var list []meta.ProjectProperty
	var err error

	if eventId != "" {
		// 获取特定事件的属性
		list, err = h.relSrv.GetEventProperties(projectAlias, eventId)
	} else {
		// 获取项目所有属性
		list, err = h.propSrv.GetList(projectAlias)
	}

	if err != nil {
		model.InternalError(c, err.Error())
		return
	}

	// 转换为通用的下拉结构 (id, name)
	type Option struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		DataType string `json:"data_type"`
	}

	options := make([]Option, 0, len(list))
	for _, e := range list {
		options = append(options, Option{
			Id:       e.PropertyId,
			Name:     e.PropertyName,
			DataType: e.DataType,
		})
	}

	model.Success(c, options)
}

package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	analyticsHandler := NewAnalyticsHandler()
	projectEventHandler := NewProjectEventHandler()
	projectPropertyHandler := NewProjectPropertyHandler()
	projectDataHandler := NewProjectDataHandler()
	selectorHandler := NewSelectorHandler()
	reportHandler := NewReportHandler()
	dashboardHandler := NewDashboardHandler()
	aiChatHandler := NewAIChatHandler()
	aiModelConfigHandler := NewAIModelConfigHandler()
	projectMetricHandler := NewProjectMetricHandler()

	v1 := r.Group("/api/v1")
	{
		selectorGroup := v1.Group("/selector")
		{
			selectorGroup.GET("/events", selectorHandler.GetEvents)
			selectorGroup.GET("/properties", selectorHandler.GetProperties)
		}

		analyticsGroup := v1.Group("/analytics")
		{
			analyticsGroup.POST("/event", analyticsHandler.EventAnalysis)
			analyticsGroup.POST("/funnel", analyticsHandler.FunnelAnalysis)
			analyticsGroup.POST("/retention", analyticsHandler.RetentionAnalysis)
			analyticsGroup.POST("/scatter", analyticsHandler.ScatterAnalysis)
			analyticsGroup.POST("/user_property", analyticsHandler.UserPropertyAnalysis)
			analyticsGroup.POST("/user_list", analyticsHandler.GetUserList)
			analyticsGroup.POST("/dimension", analyticsHandler.GetDimensions)

			// 元数据相关 (GET)
			analyticsGroup.GET("/user_properties", analyticsHandler.GetUserProperties)
			analyticsGroup.GET("/user_tags", analyticsHandler.GetUserTags)
			analyticsGroup.GET("/user_groups", analyticsHandler.GetUserGroups)
			analyticsGroup.GET("/tag_values", analyticsHandler.GetTagValues)

			// AI Chat
			analyticsGroup.POST("/ai/chat", aiChatHandler.HandleChat)
			// AI 可用模型列表
			analyticsGroup.GET("/ai/models", aiModelConfigHandler.GetEnabledModels)

			analyticsGroup.POST("/ai/sessions", aiChatHandler.SaveSession)
			analyticsGroup.GET("/ai/sessions", aiChatHandler.GetSessions)
			analyticsGroup.DELETE("/ai/sessions/:id", aiChatHandler.DeleteSession)
			analyticsGroup.POST("/ai/sessions/:id/execute", aiChatHandler.ExecuteSession)
		}

		v1.POST("/event-detail", analyticsHandler.EventDetail)

		reportGroup := v1.Group("/report")
		{
			reportGroup.GET("/list", reportHandler.GetList)
			reportGroup.GET("/:id", reportHandler.GetByID)
			reportGroup.POST("", reportHandler.Create)
			reportGroup.PUT("", reportHandler.Update)
			reportGroup.DELETE("/:id", reportHandler.Delete)
		}

		dashboardGroup := v1.Group("/dashboard")
		{
			dashboardGroup.GET("/list", dashboardHandler.GetList)
			dashboardGroup.GET("/:id", dashboardHandler.GetByID)
			dashboardGroup.POST("", dashboardHandler.Create)
			dashboardGroup.PUT("", dashboardHandler.Update)
			dashboardGroup.DELETE("/:id", dashboardHandler.Delete)

			dashboardGroup.POST("/item", dashboardHandler.AddItem)
			dashboardGroup.PUT("/item", dashboardHandler.UpdateItem)
			dashboardGroup.DELETE("/item/:id", dashboardHandler.DeleteItem)
			dashboardGroup.PUT("/items/batch", dashboardHandler.BatchUpdateItems)
		}

		metaGroup := v1.Group("/meta")
		{
			metaGroup.GET("/project_events", projectEventHandler.GetList)
			metaGroup.POST("/project_events", projectEventHandler.Create)
			metaGroup.PUT("/project_events", projectEventHandler.Update)
			metaGroup.DELETE("/project_events/:id", projectEventHandler.Delete)

			metaGroup.GET("/project_properties", projectPropertyHandler.GetList)
			metaGroup.POST("/project_properties", projectPropertyHandler.Create)
			metaGroup.PUT("/project_properties", projectPropertyHandler.Update)
			metaGroup.DELETE("/project_properties/:id", projectPropertyHandler.Delete)

			// 指标管理
			metaGroup.GET("/metrics", projectMetricHandler.GetList)
			metaGroup.POST("/metrics", projectMetricHandler.Create)
			metaGroup.PUT("/metrics/:id", projectMetricHandler.Update)
			metaGroup.DELETE("/metrics/:id", projectMetricHandler.Delete)
		}

		systemGroup := v1.Group("/system")
		{
			systemGroup.GET("/projects", projectDataHandler.GetList)
			systemGroup.POST("/projects", projectDataHandler.Create)
			systemGroup.PUT("/projects", projectDataHandler.Update)
			systemGroup.DELETE("/projects/:alias", projectDataHandler.Delete)

			// AI 模型配置管理
			systemGroup.GET("/ai-models", aiModelConfigHandler.GetList)
			systemGroup.POST("/ai-models", aiModelConfigHandler.Create)
			systemGroup.PUT("/ai-models/:id", aiModelConfigHandler.Update)
			systemGroup.DELETE("/ai-models/:id", aiModelConfigHandler.Delete)
			systemGroup.POST("/ai-models/:id/test", aiModelConfigHandler.TestConnection)
		}
	}
}

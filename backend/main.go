package main

import (
	"log"

	"gitee.com/dmp_admin_v2/backend/internal/api"
	"gitee.com/dmp_admin_v2/backend/internal/common"
	"gitee.com/dmp_admin_v2/backend/internal/core"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置和数据库
	if err := core.Init(""); err != nil {
		log.Fatalf("Failed to initialize core: %v", err)
	}

	r := gin.Default()

	// 注册路由
	api.RegisterRoutes(r)

	// 启动服务器
	common.Logger.Infof("Starting server on :8080...")
	if err := r.Run(":8080"); err != nil {
		common.Logger.Fatalf("Failed to run server: %v", err)
	}
}

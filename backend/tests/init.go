package tests

import (
	"gitee.com/dmp_admin_v2/backend/internal/core"
)

func init() {
	// 初始化配置和数据库，默认从 ../config/config.yaml 加载
	if err := core.Init("../config/config.yaml"); err != nil {
		panic(err)
	}
}

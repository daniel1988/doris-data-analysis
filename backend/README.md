# DMP Admin Backend

## 数据库初始化

项目提供了自动初始化 `dmp_center` 数据库及相关表结构的脚本。

### 运行环境
- 已正确配置 `backend/config/config.yaml` 中的 Doris 连接信息。

### 运行指令
在 `backend` 目录下执行：

```bash
go run cmd/db_init/main.go
```

### 脚本功能
1. 自动连接 Doris 数据库。
2. 检查并创建 `dmp_center` 数据库。
3. 自动检测各表是否存在，若不存在则根据最新定义进行创建。
4. 包含表：`project_data`, `project_event`, `project_property`, `project_event_property`, `query_logs`, `user_properties`, `user_tags`, `reports`, `dashboards`, `dashboard_items`。

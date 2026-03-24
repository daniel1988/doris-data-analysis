# 数据采集 API 开发文档索引

本目录用于沉淀基于 Golang 的数据采集 API 设计与实施步骤。  
数据结构基线见 [数据格式.md](file:///d:/gitee/dmp_admin_v2/docs/数据采集api/数据格式.md)。

## 技术栈

- 语言：Golang 1.22+
- Web 框架：Gin
- 配置管理：Viper
- 日志：Zap
- 数据存储：
  - 本地分钟日志文件（原始事件落盘，JSONL）
  - MySQL（配置与异常事件存储）
  - Redis（限流、幂等、防重）
  - Doris（通过 Stream Load 导入分析明细）
- 通信协议：HTTP + JSON（`snake_case`）
- 安全校验：HMAC-SHA256 签名 + nonce 防重放

## 目录结构（按阶段）

```text
dc-api/
├── cmd/
│   └── dc-api/
│       └── main.go
├── internal/
│   ├── api/
│   │   └── router.go
│   ├── controller/
│   │   ├── collect_controller.go
│   │   └── health_controller.go
│   ├── service/
│   │   ├── collect_service.go
│   │   ├── normalize_service.go
│   │   └── replay_service.go
│   ├── repo/
│   │   ├── minute_file_repo.go
│   │   ├── schema_repo.go
│   │   ├── stream_load_repo.go
│   │   ├── redis_repo.go
│   │   └── mysql_repo.go
│   ├── middleware/
│   │   ├── signature.go
│   │   ├── rate_limit.go
│   │   ├── trace.go
│   │   └── body_limit.go
│   ├── model/
│   │   ├── request.go
│   │   ├── response.go
│   │   ├── event.go
│   │   └── error_code.go
│   └── config/
│       └── config.go
├── pkg/
│   ├── logger/
│   │   └── logger.go
│   ├── validator/
│   │   └── validator.go
│   ├── signer/
│   │   └── hmac_sha256.go
│   └── metrics/
│       └── metrics.go
├── configs/
│   ├── config.dev.yaml
│   └── config.prod.yaml
├── scripts/
│   └── run-local.ps1
└── go.mod
```

## 阅读顺序

1. [01-开发步骤总览.md](file:///d:/gitee/dmp_admin_v2/docs/数据采集api/01-开发步骤总览.md)
2. [02-接口与数据模型开发步骤.md](file:///d:/gitee/dmp_admin_v2/docs/数据采集api/02-接口与数据模型开发步骤.md)
3. [03-API层开发步骤.md](file:///d:/gitee/dmp_admin_v2/docs/数据采集api/03-API层开发步骤.md)
4. [04-Service与数据处理开发步骤.md](file:///d:/gitee/dmp_admin_v2/docs/数据采集api/04-Service与数据处理开发步骤.md)
5. [05-存储链路与可靠性开发步骤.md](file:///d:/gitee/dmp_admin_v2/docs/数据采集api/05-存储链路与可靠性开发步骤.md)
6. [06-测试发布与运维步骤.md](file:///d:/gitee/dmp_admin_v2/docs/数据采集api/06-测试发布与运维步骤.md)

## 说明

- 本阶段仅输出开发方案与步骤拆解，不包含代码实现。
- 文档按“接口契约 -> API 层 -> Service 层 -> 存储链路 -> 测试发布”顺序推进，可直接作为开发任务拆分依据。

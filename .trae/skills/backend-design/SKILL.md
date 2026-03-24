---
name: "backend-design"
description: "提供后端 (Golang) API 的架构设计、分层规范与错误处理指南。在进行后端接口开发、重构或代码审查时调用。"
---

# backend-design (后端 API 设计与开发规范)

该技能用于指导基于 Golang 的后端项目架构设计、API 开发与编码规范。当涉及后端开发或代码审查时，请务必遵循以下原则。

## 1. 架构分层设计 (MVC 模式)
后端项目严格遵循职责单一的目录分层：
- **API 层 (`internal/api/`)**：即路由与中间件装配层。
  - 负责路由注册、路由分组、中间件挂载（鉴权、限流、日志、trace）。
  - 不承载业务逻辑，不直接处理复杂参数校验与业务编排。
- **Controller 层 (`internal/controller/`)**：即请求处理层。
  - 仅负责接收 HTTP 请求、参数绑定与校验。
  - 调用 Service 层执行业务逻辑。
  - 组装并返回标准响应。不应包含复杂的业务逻辑代码。
- **Service 层 (`internal/service/`)**：
  - 负责核心业务逻辑处理与计算。
  - 与数据库、外部 API（如 AI 模型）进行交互。
- **Repo 层 (`internal/repo/`)**：
  - 负责数据访问抽象（DB、Redis、MQ、外部网关）。
  - 不承载业务编排逻辑，仅提供可复用的数据读写能力。
- **Middleware 层 (`internal/middleware/`)**：
  - 放置通用中间件（鉴权、签名、限流、trace、请求日志）。
  - 由 API 路由层统一装配与复用。
- **Model 层 (`internal/model/`)**：
  - 负责定义数据库表映射结构体（GORM 实体）和公用的数据传输对象 (DTO)。

推荐目录关系：`pkg` 与 `internal` 同级，用于放置可复用基础能力（如 logger、validator、utils），避免与业务分层耦合。

## 2. API 响应与 JSON 规范
- **请求/响应格式**：前后端交互必须统一使用 **JSON** 格式。JSON 字段命名必须使用 **`snake_case`** (下划线)。
- **统一响应结构**：
  所有 API 响应必须通过调用 `internal/model/response.go`（或 `api/response`）中的统一方法返回，标准结构如下：
  ```json
  {
    "code": 200,
    "data": { ... },
    "message": "success"
  }
  ```
- **参数校验**：
  在 Controller 层，必须使用 Gin 框架的 `ShouldBindJSON` 或 `ShouldBindQuery` 方法进行参数绑定与校验。

## 3. 错误处理与日志规范
- **错误处理**：
  - 业务代码中**禁止忽略错误**（不能使用 `_` 丢弃 `err`）。
  - Controller 层发生错误时，必须通过 `model.Error`（或类似方法）返回明确的状态码和错误提示。
- **日志记录**：
  - 统一使用 `internal/common.Logger` 记录日志。
  - **禁止**使用原生的 `fmt.Println` 进行调试输出，必须通过 `Logger.Infof` / `Logger.Errorf` 等方法处理。

## 4. 命名与代码规模限制
- **文件长度**：单个 `.go` 源文件原则上不应超过 **500 行**，若过长应当按职责拆分。
- **导出标识符**：导出的变量、常量、结构体和函数使用 `PascalCase` (大驼峰)。
- **私有标识符**：包内私有的变量和函数使用 `camelCase` (小驼峰)。
- **标签**：数据库字段 (gorm) 和 JSON 标签必须使用 `snake_case` (下划线命名)。

## 5. 基础组件引用
- **配置读取**：所有环境变量和全局参数必须通过 `core.GlobalConfig` 或 `viper` 读取，禁止在业务代码中硬编码。
- **数据库连接**：使用 `core.GetProjectCenter()` 获取 GORM 实例进行数据库操作。

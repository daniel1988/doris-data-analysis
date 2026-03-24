# 02-数据模型与 API 结构

## 1. 核心数据模型 (基于 old_frontend)

### 1.1 事件明细项 (`EventDetailItem`)
Doris 中记录的原始事件数据结构：
- `e_openid`: 事件用户唯一标识 (String)
- `e_event_id`: 事件标识 (String)
- `e_event_name`: 事件显示名称 (String)
- `e_event_time`: 事件发生时间 (String, 格式：`YYYY-MM-DD HH:mm:ss`)
- `e_package_name`: 应用包名 (String)
- `e_platform`: 设备平台 (String: `android`, `ios`, `web`, `pc`)
- `e_ip`: 用户 IP 地址 (String)
- `e_request_id`: 请求 ID (String, 用于追踪完整请求链路)
- `e_properties`: 自定义事件属性 (JSON String)

### 1.2 过滤组接口 (`FilterGroup`)
- `scope`: 逻辑关系 (0: OR, 1: AND)
- `filters`: `Filter[]` 过滤条件数组

### 1.3 过滤条件接口 (`Filter`)
- `column`: 字段定义
  - `table`: 对应表名（通常为 `event_data`）
  - `field`: 对应字段名
  - `alias`: 显示别名
- `operator`: 操作符 (Number, 详见下文)
- `value`: 过滤值 (Any, 视操作符而定)

### 1.4 操作符枚举 (`FilterOperators`)
支持的操作符包括：
- `1` (EQUAL_TO): 等于
- `2` (NOT_EQUAL_TO): 不等于
- `3` (LESS_THAN): 小于
- `9` (BETWEEN): 范围 (通常用于时间段)
- `10` (IN): 包含于列表
- `7` (IS_NULL): 为空
- `12` (LIKE): 模糊匹配
- `19` (DYNAMIC_DATES): 动态日期

## 2. API 接口说明

### 2.1 查询事件明细
- **URL**: `POST /api/v1/event-detail`
- **Payload**: `EventDetailRequest`
  - `project_alias`: 项目别名 (必填)
  - `page_size`: 每页数量
  - `page_num`: 当前页码
  - `select_fields`: 需要查询的字段列表
  - `event_filter_group`: 高级过滤器对象

### 2.2 获取事件属性列表
- **URL**: `GET /api/v1/event-properties/project/{alias}`
- **功能**: 用于在“过滤器”弹窗中选择可用的事件属性字段。

### 2.3 获取事件列表
- **URL**: `GET /api/v1/project-events/project/{alias}`
- **功能**: 用于在快速搜索栏中提供事件 ID 的下拉选择和模糊搜索。

# 04-Service与数据处理开发步骤

## 1. Service 核心职责

`CollectService` 负责：
- 多种输入格式统一拆解
- 数据标准化（NormalizedEvent）
- 业务校验（限额、黑名单、字段白名单）
- 调用下游存储/队列

## 2. 统一处理流程

1. 接收 DTO
2. 展平为 `[]NormalizedEvent`
3. 执行校验与标准化
4. 执行去重（可选）
5. 批量推送到 Repo（DB/MQ）
6. 返回 accepted / dropped 统计

## 3. 标准化规则

- `event_time` 转为统一时区（UTC 或系统配置时区）
- `event_name` 可为空时，使用 `event_id` 兜底
- `properties` 转 JSON 字符串，处理不可序列化值
- 补充系统字段：
  - `ingest_time`
  - `ip`
  - `ua`
  - `trace_id`

## 4. 数据质量规则

建议第一版开启：
- key 合法字符校验（字母、数字、下划线）
- value 类型白名单（string/number/bool）
- 超大字段截断或丢弃
- 单请求最大事件数限制

## 5. 幂等与防重（可选）

- Redis Key：`dedup:{project}:{hash}`
- TTL：1~10 分钟（按业务）
- 命中防重则计入 `dropped`

## 6. 本步骤验收标准

- 三类输入最终均可输出统一标准事件
- 返回统计中包含 `accepted` 与 `dropped`
- 异常数据不会导致整批失败（可配置 fail-fast）

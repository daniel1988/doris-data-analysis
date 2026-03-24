# 03-API层开发步骤

## 1. 路由设计

在 `router.go` 注册：
- `POST /api/v1/collect/event`
- `POST /api/v1/collect/user-batch`
- `POST /api/v1/collect/batch`
- `GET /api/v1/healthz`

建议将采集路由单独分组，便于应用签名校验与限流中间件。

## 2. Handler 职责

每个 Handler 仅做四件事：
1. `ShouldBindJSON` 参数绑定
2. 基础参数校验
3. 调用 `CollectService`
4. 返回统一响应结构

禁止在 Handler 中写入：
- 批量拆解逻辑
- 底层存储逻辑
- 队列重试逻辑

## 3. 错误码建议

- `40001` 参数非法
- `40002` 时间戳非法
- `40003` 批量条数超限
- `40101` 签名缺失（Header 不完整）
- `40102` 签名校验失败
- `40103` nonce 重放
- `40104` 时间戳超出允许窗口
- `42901` 请求过载
- `50001` 内部处理失败

## 4. 接口签名与防重放校验

签名算法采用 `HMAC-SHA256`，并引入 `nonce` 防重放机制。

### 4.1 必需请求头

- `X-Project-Alias`：项目别名
- `X-Timestamp`：时间戳（秒或毫秒）
- `X-Nonce`：随机字符串（建议 16~32 字符）
- `X-Sign`：签名字符串（小写 hex）

### 4.2 服务端校验规则

- 时间戳允许偏差：`±5 分钟`
- nonce 有效期：`10 分钟`
- 同一 `(project_alias, timestamp, nonce)` 在有效期内不可重复
- 时间戳支持秒或毫秒，毫秒值在服务端统一转换为秒

### 4.3 签名串规范

1. 参与签名字段数组：`[project_alias, timestamp, nonce]`
2. 数组元素按字典序升序排序
3. 使用换行符 `\n` 拼接为 `message`
4. 使用项目 `secret` 作为 key 计算 `HMAC-SHA256(message)`
5. 结果编码为小写十六进制，作为 `X-Sign`

### 4.4 服务端实现建议

- 在 API 中间件中完成签名校验，不进入 Handler
- `secret` 通过 `project_alias` 从配置中心或项目配置表读取
- nonce 去重建议使用 Redis：
  - key：`sign_nonce:{project_alias}:{timestamp}:{nonce}`
  - ttl：`10 分钟`
- 若签名不通过，统一返回 `40102`

## 5. 中间件顺序

建议顺序：
1. `trace_id` 注入
2. 签名 Header 完整性校验
3. 时间戳窗口校验
4. nonce 防重放校验
5. HMAC 签名校验
6. 鉴权（可选）
7. 限流
8. 请求体大小限制
9. Handler

## 6. 接口回包建议

### 成功（同步接收成功）
```json
{
  "code": 200,
  "data": {
    "accepted": 100
  },
  "message": "success"
}
```

### 失败（参数问题）
```json
{
  "code": 40001,
  "data": null,
  "message": "event_time is invalid"
}
```

### 失败（签名问题）
```json
{
  "code": 40102,
  "data": null,
  "message": "invalid signature"
}
```

## 7. 本步骤验收标准

- 四条路由可访问
- 参数异常可稳定返回定义错误码
- 签名、时间戳、nonce 防重放校验生效
- 采集成功时可返回 accepted 数量

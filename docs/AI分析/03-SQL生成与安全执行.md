# SQL 生成与安全执行沙箱

对话生成 SQL 给后端直接执行，就像给大门留了钥匙，存在极大的合规和破坏风险（如 `DROP DATABASE`，或者通过 `CROSS JOIN` 直接搞挂集群）。

由于系统采用了 Doris 这样的大数据底座，我们需要非常严密的拦截与审查层：**SQL Validator**。

## 1. SQL 生成的限定规则
在传递给 LLM 的系统级 Prompt 中，需要用最高权重告知：
1. **只能生成 `SELECT` 语句**。严禁生成 `UPDATE`, `DELETE`, `INSERT`, `DROP`, `ALTER`。
2. 必须包含项目过滤：由于平台存在多项目混表隔离方案，所有的查询必须隐式或显式地带有 `project_alias = 'xx'` 的条件（视底层数据流而定）。

## 2. 后端 AST 解析与安全性拦截 (SQL Parser)
当 LLM 返回形如 `SELECT * FROM event_data ...` 的 SQL 时，**绝对不能**直接拼接到 `db.Query()`。

### 步骤一：语法树(AST)解析与验证
在 Go 后端引入轻量级的 SQL Parser（如 `github.com/xwb1989/sqlparser` / PingCAP parser）：
1. **语句类型校验**：检查 AST 的根节点必须是 `SelectStmt`。查到非 SELECT 类型直接抛出异常 `Access Denied: Only SELECT operations are allowed.`
2. **表名访问限制**：遍历 AST 中的所有 `TableName` 节点。校验其访问的表如果在 `[event_data, user_data, user_tag_data]` 允许的白名单白名单之外，直接拒绝。

### 步骤二：项目越权保护 (RLS 注入)
在物理表中，如果是基于 `project_alias` 维度的多租户架构。我们可以在 AST 层面自动并且强制添加一条 `WHERE project_alias = ?` （从当前的会话 Token 解析出用户所查询的项目别名），保证即使用户诱导 LLM 生成了跨界 SQL，也只能在自己当前的项目里进行查找。

### 步骤三：运算量熔断 (Resource Limit)
为防止恶意的深度全量扫描：
1. **强制分页（LIMIT 拦截）**：如果查询未包含 LIMIT 且不含 GROUP BY 聚合函数，强制追加 `LIMIT 100` 或 `LIMIT 1000`。
2. **强制时间收敛**：如果未包含 `e_event_time` 或 `u_event_time` 作为分区维度的限制条件，提示引擎将其查询窗口默认限制为最近 `30天`（或通过 AST 写入 `WHERE event_time >= NOW() - INTERVAL 30 DAY`）。

## 3. Auto-Healing (自动纠错机制)
当 Doris 引擎返回了 Syntax Error（例如：尝试使用不存在的字段）：
1. 捕获 Doris 抛出的错误栈。
2. 后端组合：`"The SQL you generated caused an error: [ERROR CODE: xxx]. Please fix it based on our schema. The wrong SQL was: [xxx]"`。
3. 循环调用 LLM 进行修复，最多重试 2 次。若 2 次依然报错，则抛出让用户手动改动的提示信息。

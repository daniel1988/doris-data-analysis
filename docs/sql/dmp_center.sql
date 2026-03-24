-- Doris dmp_center 实际表结构 (由 SHOW CREATE TABLE 获取)
CREATE DATABASE IF NOT EXISTS dmp_center;
USE dmp_center;

-- 1. 项目数据表
CREATE TABLE `project_data` (
  `project_alias` varchar(58) NULL DEFAULT "" COMMENT "项目别名",
  `project_name` varchar(58) NULL DEFAULT "" COMMENT "项目名称",
  `region` varchar(128) NULL DEFAULT "" COMMENT "设备号",
  `secret` varchar(58) NULL DEFAULT "" COMMENT "签名密钥",
  `create_time` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=OLAP
UNIQUE KEY(`project_alias`)
DISTRIBUTED BY HASH(`project_alias`) BUCKETS 1
PROPERTIES (
  "enable_unique_key_merge_on_write" = "true"
);

-- 2. 元事件定义表
CREATE TABLE `project_event` (
  `project_alias` varchar(58) NULL DEFAULT "",
  `event_id` varchar(58) NULL DEFAULT "",
  `event_name` varchar(65533) NULL DEFAULT "",
  `event_type` int NULL DEFAULT "0",
  `create_time` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=OLAP
UNIQUE KEY(`project_alias`, `event_id`)
DISTRIBUTED BY HASH(`project_alias`, `event_id`) BUCKETS 1
PROPERTIES (
  "enable_unique_key_merge_on_write" = "true"
);

-- 3. 元属性定义表 (各项目 event_data 字段)
CREATE TABLE `project_property` (
  `project_alias` varchar(58) NULL DEFAULT "",
  `property_id` varchar(58) NULL DEFAULT "",
  `property_name` varchar(65533) NULL DEFAULT "",
  `data_type` varchar(58) NULL,
  `property_type` int NULL DEFAULT "0",
  `create_time` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=OLAP
UNIQUE KEY(`project_alias`, `property_id`)
DISTRIBUTED BY HASH(`project_alias`, `property_id`) BUCKETS 1
PROPERTIES (
  "enable_unique_key_merge_on_write" = "true"
);

-- 4. 事件-属性关联表（各项目事件与字段映射表）
CREATE TABLE `project_event_property` (
  `project_alias` varchar(58) NULL DEFAULT "",
  `event_id` varchar(58) NULL DEFAULT "",
  `property_id` varchar(58) NULL DEFAULT "",
  `create_time` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=OLAP
UNIQUE KEY(`project_alias`, `event_id`, `property_id`)
DISTRIBUTED BY HASH(`project_alias`, `event_id`) BUCKETS 1
PROPERTIES (
  "enable_unique_key_merge_on_write" = "true"
);

-- 5. 查询日志表
CREATE TABLE `query_logs` (
  `project_alias` varchar(64) NULL COMMENT "项目标识",
  `create_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT "创建时间",     
  `create_user` int NULL DEFAULT "0" COMMENT "创建人",
  `sql` text NULL COMMENT "SQL语句",
  `status` varchar(20) NULL COMMENT "执行状态",
  `error_msg` text NULL COMMENT "错误信息",
  `row_count` bigint NULL DEFAULT "0" COMMENT "影响行数",
  `duration` double NULL DEFAULT "0"
) ENGINE=OLAP
DUPLICATE KEY(`project_alias`, `create_time`)
PARTITION BY RANGE(`create_time`) (
  PARTITION p202603 VALUES [('2026-03-01 00:00:00'), ('2026-04-01 00:00:00'))
)
DISTRIBUTED BY HASH(`create_time`) BUCKETS AUTO;

-- 6. 用户属性定义表（各项目 user_data 字段）
CREATE TABLE `user_properties` (
  `project_alias` varchar(58) NULL DEFAULT "",
  `property_id` varchar(58) NULL DEFAULT "",
  `property_name` varchar(65533) NULL DEFAULT "",
  `event_type` varchar(28) NULL,
  `data_type` varchar(58) NULL,
  `create_time` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=OLAP
UNIQUE KEY(`project_alias`, `property_id`)
DISTRIBUTED BY HASH(`project_alias`, `property_id`) BUCKETS 1
PROPERTIES (
  "enable_unique_key_merge_on_write" = "true"
);

-- 7. 用户标签/分群表
CREATE TABLE `user_tags` (
  `project_alias` varchar(58) NULL DEFAULT "",
  `tag_code` varchar(58) NULL DEFAULT "",
  `tag_name` varchar(128) NULL DEFAULT "",
  `tag_sql` text NULL DEFAULT "",
  `user_count` double NULL DEFAULT "0",
  `create_time` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=OLAP
UNIQUE KEY(`project_alias`, `tag_code`)
DISTRIBUTED BY HASH(`project_alias`, `tag_code`) BUCKETS 1
PROPERTIES (
  "enable_unique_key_merge_on_write" = "true"
);

-- 8. 报表配置表
CREATE TABLE `reports` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT "报表ID",
  `project_alias` varchar(58) NULL DEFAULT "" COMMENT "项目别名",
  `name` varchar(128) NULL DEFAULT "" COMMENT "报表名称",
  `category` varchar(58) NULL DEFAULT "" COMMENT "分类 (event, funnel, retention)",
  `description` varchar(255) NULL DEFAULT "" COMMENT "报表描述",
  `query_params` text NULL COMMENT "查询参数JSON",
  `create_user` int NULL DEFAULT "0" COMMENT "创建人ID",
  `create_time` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=OLAP
UNIQUE KEY(`id`)
DISTRIBUTED BY HASH(`id`) BUCKETS 1
PROPERTIES (
  "enable_unique_key_merge_on_write" = "true"
);

-- 9. 看板配置表
DROP TABLE IF EXISTS `dashboards`;
CREATE TABLE `dashboards` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT "看板ID",
  `project_alias` varchar(58) NULL DEFAULT "" COMMENT "项目别名",
  `name` varchar(128) NULL DEFAULT "" COMMENT "看板内部名称",
  `display_name` varchar(128) NULL DEFAULT "" COMMENT "显示名称",
  `description` varchar(255) NULL DEFAULT "" COMMENT "描述",
  `category` varchar(58) NULL DEFAULT "custom" COMMENT "分类 (overview, analytical, custom)",
  `layout_type` varchar(58) NULL DEFAULT "grid" COMMENT "布局模式 (grid, flex, custom)",
  `grid_config` text NULL COMMENT "全局网格配置",
  `theme` varchar(20) NULL DEFAULT "light" COMMENT "主题 (light, dark, auto)",
  `refresh_interval` int NULL DEFAULT "0" COMMENT "刷新频率(秒)",
  `filters` text NULL COMMENT "全局筛选器配置",
  `variables` text NULL COMMENT "看板变量",
  `status` varchar(20) NULL DEFAULT "published" COMMENT "状态 (draft, published, archived)",
  `owner_id` int NULL DEFAULT "0" COMMENT "拥有者ID",
  `create_time` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=OLAP
UNIQUE KEY(`id`)
DISTRIBUTED BY HASH(`id`) BUCKETS 1
PROPERTIES (
  "enable_unique_key_merge_on_write" = "true"
);

-- 10. 看板组件表
DROP TABLE IF EXISTS `dashboard_items`;
CREATE TABLE `dashboard_items` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT "组件ID",
  `dashboard_id` bigint NOT NULL COMMENT "所属看板ID",
  `report_id` bigint NULL DEFAULT "0" COMMENT "关联报表ID",
  `type` varchar(58) NULL DEFAULT "" COMMENT "组件类型 (chart, table, mixed, stat, map)",
  `title` varchar(128) NULL DEFAULT "" COMMENT "组件覆盖标题",
  `position_x` int NULL DEFAULT "0" COMMENT "网格 X 坐标",
  `position_y` int NULL DEFAULT "0" COMMENT "网格 Y 坐标",
  `width` int NULL DEFAULT "12" COMMENT "网格宽度",
  `height` int NULL DEFAULT "8" COMMENT "网格高度",
  `z_index` int NULL DEFAULT "1" COMMENT "层级深度",
  `config_override` text NULL COMMENT "报表配置覆盖JSON",
  `is_visible` boolean NULL DEFAULT "1" COMMENT "是否可见",
  `create_time` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=OLAP
UNIQUE KEY(`id`)
DISTRIBUTED BY HASH(`id`) BUCKETS 1
PROPERTIES (
  "enable_unique_key_merge_on_write" = "true"
);

-- 11. AI 模型配置表
CREATE TABLE `ai_model_config` (
  `id`            bigint      NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `provider`      varchar(58) NULL DEFAULT ''  COMMENT '提供商标识 (openai/deepseek/qwen/glm/kimi/ollama)',
  `display_name`  varchar(128) NULL DEFAULT '' COMMENT '显示名称',
  `base_url`      varchar(512) NULL DEFAULT '' COMMENT 'API基地址',
  `api_key`       varchar(512) NULL DEFAULT '' COMMENT 'API密钥',
  `model_name`    varchar(128) NULL DEFAULT '' COMMENT '模型标识',
  `max_tokens`    int          NULL DEFAULT '4096' COMMENT '最大生成Token数',
  `temperature`   double       NULL DEFAULT '0.1'  COMMENT '温度参数',
  `is_default`    boolean      NULL DEFAULT '0'    COMMENT '是否为默认模型',
  `is_enabled`    boolean      NULL DEFAULT '1'    COMMENT '是否启用',
  `sort_order`    int          NULL DEFAULT '0'    COMMENT '排序序号',
  `create_time`   datetime     NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time`   datetime     NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=OLAP
UNIQUE KEY(`id`)
DISTRIBUTED BY HASH(`id`) BUCKETS 1
PROPERTIES (
  "replication_num" = "1",
  "enable_unique_key_merge_on_write" = "true"
);

-- 12. 指标配置表
CREATE TABLE `project_metrics` (
  `project_alias` varchar(58) NOT NULL COMMENT '项目别名',
  `metric_code` varchar(128) NOT NULL COMMENT '指标英文标识 (如: dau_count)',
  `metric_name` varchar(128) NOT NULL COMMENT '指标名称 (如: DAU)',
  `expression` text NOT NULL COMMENT '指标计算表达式 (如: COUNT(DISTINCT CASE WHEN e_event_id=''sys.login'' THEN e_openid END))',
  `base_table` varchar(64) NOT NULL DEFAULT 'event_data' COMMENT '基于哪张表 (event_data/user_data)',
  `description` text NULL COMMENT '业务口径说明',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态: 1启用 0停用',
  `create_time` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=OLAP
UNIQUE KEY(`project_alias`, `metric_code`)
DISTRIBUTED BY HASH(`project_alias`, `metric_code`) BUCKETS 1
PROPERTIES (
    "replication_num" = "1",
  "enable_unique_key_merge_on_write" = "true"
);

-- 12. 指标配置表
CREATE TABLE `project_metrics` (
  `project_alias` varchar(58) NOT NULL COMMENT '项目别名',
  `metric_code` varchar(128) NOT NULL COMMENT '指标英文标识 (如: dau_count)',
  `metric_name` varchar(128) NOT NULL COMMENT '指标名称 (如: DAU)',
  `expression` text NOT NULL COMMENT '指标计算表达式 (如: COUNT(DISTINCT CASE WHEN e_event_id=''sys.login'' THEN e_openid END))',
  `base_table` varchar(64) NOT NULL DEFAULT 'event_data' COMMENT '基于哪张表 (event_data/user_data)',
  `description` text NULL COMMENT '业务口径说明',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态: 1启用 0停用',
  `create_time` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=OLAP
UNIQUE KEY(`project_alias`, `metric_code`)
DISTRIBUTED BY HASH(`project_alias`, `metric_code`) BUCKETS 1
PROPERTIES (
  "enable_unique_key_merge_on_write" = "true"
);

-- 13. AI 会话记录表
CREATE TABLE `ai_chat_sessions` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `project_alias` varchar(58) NOT NULL COMMENT '项目别名',
  `user_id` bigint NOT NULL COMMENT '保存该记录的用户ID',
  `user_query` text NOT NULL COMMENT '用户的自然语言提问',
  `llm_sql` text NOT NULL COMMENT 'AI 生成并验证通过的 SQL',
  `viz_type` varchar(32) NULL COMMENT '图表类型 (bar, line, pie, table 等)',
  `x_axis` varchar(128) NULL COMMENT 'X轴字段',
  `y_axis` varchar(128) NULL COMMENT 'Y轴字段',
  `narrative` text NULL COMMENT 'AI 给出的结论描述',
  `create_time` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=OLAP
UNIQUE KEY(`id`)
DISTRIBUTED BY HASH(`id`) BUCKETS 1
PROPERTIES (
  "enable_unique_key_merge_on_write" = "true"
);

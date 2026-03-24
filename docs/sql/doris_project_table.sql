-- Doris 存储层建表语句模版
-- 数据库名建议使用项目别名，例如：zgmgr4
-- CREATE DATABASE IF NOT EXISTS zgmgr4;
-- USE zgmgr4;

-- 事件数据表 (Duplicate Key 模型)
-- 用于存储所有用户原始行为流水
CREATE TABLE `event_data` (
  `e_event_id` varchar(58) NULL COMMENT "事件标识",
  `e_openid` varchar(128) NULL COMMENT "用户ID",
  `e_event_time` datetime NULL COMMENT "事件触发时间",
  `e_create_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT "入库时间",
  `e_scene` varchar(28) NULL COMMENT "场景",
  `e_device_id` text NULL COMMENT "设备ID",
  `e_ip` text NULL COMMENT "IP地址",
  `e_properties` json NULL COMMENT "事件属性",
  `e_event_type` varchar(28) NULL COMMENT "事件类型",
  `e_event_name` varchar(28) NULL COMMENT "事件名称",
  `e_os` varchar(28) NULL COMMENT "操作系统",
  `e_os_version` varchar(28) NULL COMMENT "操作系统版本",
  `e_request_id` varchar(58) NULL COMMENT "请求ID",
  `e_app_version` varchar(28) NULL COMMENT "应用版本",
  `e_from_scene` varchar(28) NULL COMMENT "来源场景",
  `e_platform` varchar(128) NULL,
  `e_sim_property` double NULL,
  `e_version` varchar(128) NULL,
  `e_project_alias` varchar(128) NULL,
  `e_log_id` varchar(128) NULL,
  `e_id` double NULL,
  `e_is_first` double NULL,
  `e_login_game_version` varchar(128) NULL,
  `e_country` varchar(128) NULL,
  `e_country_code` varchar(128) NULL,
  `e_login_cfg_version` varchar(128) NULL,
  `e_region` varchar(128) NULL,
  `e_register_time` double NULL,
  `e_type` varchar(128) NULL,
  `e_ab_test` double NULL,
  `e_channel` double NULL,
  `e_duration` double NULL,
  `e_positionTag` varchar(128) NULL,
  `e_from_ad_id` varchar(128) NULL,
  `e_from_creative_id` varchar(128) NULL,
  `e_last_login_time` double NULL,
  `e_longitude` double NULL,
  `e_account_id` varchar(128) NULL,
  `e_channel_id` varchar(128) NULL,
  `e_gid` varchar(128) NULL,
  `e_reg_cfg_version` varchar(128) NULL,
  `e_reg_game_version` varchar(128) NULL,
  `e_revive_progress` double NULL,
  `e_first_login_time` double NULL,
  `e_latitude` double NULL,
  `e_level` double NULL,
  `e_lose_progress` double NULL,
  `e_city` varchar(128) NULL,
  `e_from_advertiser_id` varchar(128) NULL,
  `e_mode` varchar(128) NULL,
  `e_tag_id` varchar(128) NULL,
  `e_time` datetime NULL,
  `e_sim_version` varchar(128) NULL
) ENGINE=OLAP
DUPLICATE KEY(`e_event_id`, `e_openid`, `e_event_time`)
COMMENT 'OLAP'
PARTITION BY RANGE(`e_event_time`)
(PARTITION p202602 VALUES [('2026-02-01 00:00:00'), ('2026-03-01 00:00:00')),
PARTITION p202603 VALUES [('2026-03-01 00:00:00'), ('2026-04-01 00:00:00')),
PARTITION p202604 VALUES [('2026-04-01 00:00:00'), ('2026-05-01 00:00:00')),
PARTITION p202605 VALUES [('2026-05-01 00:00:00'), ('2026-06-01 00:00:00')),
PARTITION p202606 VALUES [('2026-06-01 00:00:00'), ('2026-07-01 00:00:00')))
DISTRIBUTED BY HASH(`e_event_time`) BUCKETS AUTO
PROPERTIES (
"replication_allocation" = "tag.location.default: 3",
"in_memory" = "false",
"storage_format" = "V2"
);

-- 用户属性表 (Unique Key 模型)
-- 存储用户最新的画像/属性信息
CREATE TABLE `user_data` (
  `u_openid` varchar(58) NULL DEFAULT "",
  `u_event_time` datetime NULL,
  `u_create_time` datetime NULL,
  `u_update_time` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  `u_ab_test` int NULL,
  `u_gid` varchar(58) NULL,
  `u_channel` varchar(58) NULL,
  `u_from_creative_id` text NULL,
  `u_from_ad_id` text NULL,
  `u_from_advertiser_id` text NULL,
  `u_deviceid` text NULL,
  `u_reg_cfg_version` varchar(58) NULL,
  `u_reg_game_version` varchar(58) NULL,
  `u_from_scene` varchar(58) NULL,
  `u_from_launch_from` varchar(58) NULL,
  `u_from_etag` varchar(58) NULL,
  `u_ip` text NULL,
  `u_city` text NULL,
  `u_country` text NULL
) ENGINE=OLAP
UNIQUE KEY(`u_openid`)
COMMENT "用户属性表"
DISTRIBUTED BY HASH(`u_openid`) BUCKETS 16
PROPERTIES (
"replication_allocation" = "tag.location.default: 3",
"in_memory" = "false",
"storage_format" = "V2"
);

-- 用户标签数据表 (Unique Key 模型)
-- 用于存放用户分群、标签等离线计算结果
CREATE TABLE IF NOT EXISTS `user_tag_data` (
  `ut_openid` varchar(100) NOT NULL COMMENT '用户唯一标识',
  `ut_tag_code` varchar(100) NOT NULL COMMENT '标签/分群编码',
  `ut_tag_value` varchar(100) NOT NULL COMMENT '标签值'
) ENGINE=OLAP
UNIQUE KEY(`ut_openid`, `ut_tag_code`)
COMMENT "用户标签数据表"
DISTRIBUTED BY HASH(`ut_openid`) BUCKETS 16
PROPERTIES (
"replication_allocation" = "tag.location.default: 3"
);

CREATE TABLE `tb_algo_info` (
                                `id` int unsigned NOT NULL COMMENT '算法ID, 主键',
                                `algo_name` varchar(32) DEFAULT NULL COMMENT '算法名称',
                                `provider` tinyint unsigned DEFAULT NULL COMMENT '算法厂商ID',
                                `provider_name` varchar(32) DEFAULT NULL COMMENT '算法厂商名称',
                                `uuser_id` int unsigned DEFAULT NULL COMMENT '算法厂商用户的ID',
                                `algorithm_type` smallint unsigned DEFAULT NULL COMMENT '算法类型, 算法厂商内部唯一',
                                `algorithm_type_name` varchar(16) DEFAULT NULL COMMENT '算法类型, 算法厂商内部使用',
                                `algorithm_status` tinyint unsigned DEFAULT NULL COMMENT '算法状态: 0可用 1不可用',
                                `parameter` varchar(2048) DEFAULT NULL COMMENT '算法所需参数',
                                `risk_group` int unsigned DEFAULT NULL COMMENT '算法风控组',
                                `create_time` bigint unsigned DEFAULT NULL COMMENT '创建时间戳,单位:秒',
                                `version` int unsigned DEFAULT NULL COMMENT '成交记录ID',
                                PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb3 COMMENT='算法信息表(MyISAM)'
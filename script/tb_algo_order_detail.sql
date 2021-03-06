CREATE TABLE `tb_algo_order_detail` (
                                        `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
                                        `child_order_id` bigint NOT NULL COMMENT '订单ID',
                                        `algorithm_type` int unsigned NOT NULL COMMENT '算法类型',
                                        `algorithm_id` int unsigned NOT NULL COMMENT '算法ID',
                                        `usecurity_id` int unsigned NOT NULL COMMENT '证券ID',
                                        `security_id` varchar(8) DEFAULT NULL COMMENT '证券代码',
                                        `order_qty` bigint NOT NULL COMMENT '委托订单数量',
                                        `price` bigint NOT NULL COMMENT '委托订单价格',
                                        `order_type` smallint unsigned DEFAULT NULL COMMENT '订单类型 ：1-限价委托 2-本方最优 3-对手方最优 4-市价立即成交剩余撤销 5-市价全额成交或撤销 6-市价最优五档全额成交剩余撤销 7-限价全额成交或撤销(期权用）',
                                        `last_px` bigint DEFAULT NULL COMMENT '成交价格',
                                        `last_qty` bigint DEFAULT NULL COMMENT '成交数量',
                                        `com_qty` bigint DEFAULT NULL COMMENT '累计成交数量',
                                        `arrived_price` bigint DEFAULT NULL COMMENT '到达价格',
                                        `ord_status` smallint unsigned DEFAULT NULL COMMENT '订单状态 1-新建 2-成交 3-撤销 4-拒绝',
                                        `transact_time` bigint NOT NULL COMMENT '交易时间',
                                        `transact_at` bigint DEFAULT NULL COMMENT '交易时间（精确到分钟）',
                                        `proc_status` smallint unsigned DEFAULT NULL COMMENT '处理状态   0-未处理   1-已处理',
                                        `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                        PRIMARY KEY (`id`),
                                        KEY `transact_at` (`transact_at`),
                                        KEY `algo_id,security_id` (`algorithm_id`,`usecurity_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1033 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='算法子单详情表'
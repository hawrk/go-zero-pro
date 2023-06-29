-- auto-generated definition
create table tb_algo_profile_his
(
    id                 bigint auto_increment comment '自增ID'
        primary key,
    date               int                                 not null comment '日期 (20220612)',
    batch_no           bigint                              null comment '批次号',
    account_id         varchar(45)                         not null comment '用户ID （交易账户ID）',
    account_name       varchar(45)                         null comment '账户名称',
    account_type       tinyint                             null comment '用户类型  1-普通用户，2-虚拟用户（汇总用户）',
    provider           varchar(45)                         null comment '厂商名称',
    algo_id            int                                 null comment '算法ID',
    algo_name          varchar(45)                         null comment '算法名称',
    algo_type          int                                 null comment '算法类型',
    sec_id             varchar(12)                         null comment '证券代码',
    sec_name           varchar(45)                         null comment '证券名称',
    algo_order_id      bigint                              null comment '母单ID',
    industry           varchar(45)                         null comment '行业类型',
    fund_type          tinyint                             null comment '市值类型    1- huge超大 2-big大 3-middle中等 4-small小',
    liquidity          tinyint                             null comment '流动性   1-高 2-中 3-低',
    trade_cost         bigint                              null comment '交易成本（ 成交价格* 成交数量）',
    total_trade_amount bigint                              null comment '双边总交易额， 总金额',
    total_trade_fee    bigint                              null comment '总手续费（券商手续旨，过户费，印花税）',
    cross_fee          bigint                              null comment '流量费',
    profit_amount      bigint                              null comment '盈亏金额',
    profit_rate        decimal(20, 4)                      null comment '收益率， 盈亏比率',
    cancel_rate        decimal(20, 4)                      null comment '撤单率',
    progress_rate      decimal(20, 4)                      null comment '完成度',
    mini_split_order   int                                 null comment '最小拆单单位',
    mini_joint_rate    decimal(20, 4)                      null comment '最小贴合度',
    withdraw_rate      decimal(20, 4)                      null comment '回撤比例',
    vwap_dev           decimal(20, 4)                      null comment 'vwap 滑点',
    entrust_vwap       decimal(20, 4)                      null comment '委托vwap( sum (委托价格*成交数量））',
    trade_count        int                                 null comment '交易次数（一个子单交易回执算一次交易次数，被 拒绝或撤销的不算， 用来算滑点平均值）',
    trade_count_plus   int                                 null comment '盈亏为正的交易次数，统计胜率',
    avg_deviation      decimal(20, 4)                      null comment '滑点平均值',
    standard_deviation decimal(20, 4)                      null comment '滑点标准差',
    pf_rate_std_dev    decimal(20, 4)                      null comment '收益率标准差',
    factor             decimal(20, 4)                      null comment '绩效因子',
    order_time         bigint                              null comment '订单开始时间（母单开始时间)',
    source_from        tinyint                             null comment '数据来源，1:总线推送，2:数据导入',
    create_time        timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    update_time        timestamp                           null comment '更新时间'
)
    collate = utf8mb4_bin;
-- 创建存储过程（查询tb_algo_profile表最近30天的数据并插入到tb_algo_profile_his中）
DELIMITER $$
CREATE PROCEDURE profile_backup()
BEGIN
REPLACE INTO tb_algo_profile_his  select * from tb_algo_profile  where IF(update_time,update_time,create_time) between DATE_SUB(NOW(),INTERVAL 30 day) and DATE_SUB(NOW(),INTERVAL 0 day);
DELETE FROM tb_algo_profile  where IF(update_time,update_time,create_time) between DATE_SUB(NOW(),INTERVAL 30 day) and DATE_SUB(NOW(),INTERVAL 0 day);
END $$
DELIMITER ;

--创建定时任务（执行备份的存储过程，每30天执行一次）
CREATE EVENT profile_event
ON SCHEDULE EVERY 30 DAY
on completion preserve ENABLE
do CALL profile_backup();
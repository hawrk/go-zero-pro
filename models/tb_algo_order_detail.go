// Package models
/*
 Author: hawrkchen
 Date: 2022/3/28 11:42
 Desc:
*/
package models

import (
	"time"
)

// 算法子单详情表
type TbAlgoOrderDetail struct {
	Id            uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:自增ID" json:"id"`
	Date          int       `gorm:"column:date;type:int(11);comment:时间" json:"date"`
	BatchNo       int64     `gorm:"column:batch_no;type:bigint(20);comment:批次号;NOT NULL" json:"batch_no"`
	BatchName     string    `gorm:"column:batch_name;type:varchar(45);comment:批次号名称" json:"batch_name"`
	AlgoOrderId   uint      `gorm:"column:algo_order_id;type:int(10) unsigned;comment:母单号;NOT NULL" json:"algo_order_id"`
	ChildOrderId  int64     `gorm:"column:child_order_id;type:bigint(20);comment:订单ID;NOT NULL" json:"child_order_id"`
	AlgorithmType uint      `gorm:"column:algorithm_type;type:int(10) unsigned;comment:算法类型;NOT NULL" json:"algorithm_type"`
	AlgorithmId   uint      `gorm:"column:algorithm_id;type:int(10) unsigned;comment:算法ID;NOT NULL" json:"algorithm_id"`
	UserId        string    `gorm:"column:user_id;type:varchar(45);comment:用户ID" json:"user_id"`
	UsecurityId   uint      `gorm:"column:usecurity_id;type:int(10) unsigned;comment:证券ID;NOT NULL" json:"usecurity_id"`
	SecurityId    string    `gorm:"column:security_id;type:varchar(8);comment:证券代码" json:"security_id"`
	TradeSide     int       `gorm:"column:trade_side;type:tinyint(4);comment:买卖方向" json:"trade_side"`
	OrderQty      int64     `gorm:"column:order_qty;type:bigint(20);comment:委托订单数量;NOT NULL" json:"order_qty"`
	Price         int64     `gorm:"column:price;type:bigint(20);comment:委托订单价格;NOT NULL" json:"price"`
	OrderType     uint      `gorm:"column:order_type;type:smallint(5) unsigned;comment:订单类型 ：1-限价委托 2-本方最优 3-对手方最优 4-市价立即成交剩余撤销 5-市价全额成交或撤销 6-市价最优五档全额成交剩余撤销 7-限价全额成交或撤销(期权用）" json:"order_type"`
	LastPx        int64     `gorm:"column:last_px;type:bigint(20);comment:成交价格" json:"last_px"`
	LastQty       int64     `gorm:"column:last_qty;type:bigint(20);comment:成交数量" json:"last_qty"`
	ComQty        int64     `gorm:"column:com_qty;type:bigint(20);comment:累计成交数量" json:"com_qty"`
	ArrivedPrice  int64     `gorm:"column:arrived_price;type:bigint(20);comment:到达价格" json:"arrived_price"`
	TotalFee      float64   `gorm:"column:total_fee;type:decimal(20,2);comment:手续费" json:"total_fee"`
	OrdStatus     uint      `gorm:"column:ord_status;type:smallint(5) unsigned;comment:订单状态 1-新建 2-成交 3-撤销 4-拒绝" json:"ord_status"`
	TransactTime  int64     `gorm:"column:transact_time;type:bigint(20);comment:交易时间;NOT NULL" json:"transact_time"`
	TransactAt    int64     `gorm:"column:transact_at;type:bigint(20);comment:交易时间（精确到分钟）" json:"transact_at"`
	ProcStatus    uint      `gorm:"column:proc_status;type:smallint(5) unsigned;comment:处理状态   0-未处理   1-已处理" json:"proc_status"`
	Source        int       `gorm:"column:source;type:tinyint(4);comment:数据来源 0-总线 1-数据修复 2-数据导入" json:"source"`
	CreateTime    time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间;NOT NULL" json:"create_time"`
}

func (m *TbAlgoOrderDetail) TableName() string {
	return "tb_algo_order_detail"
}

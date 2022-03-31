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
	Id            uint64       `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`        // 自增ID
	ChildOrderId  int64        `gorm:"column:child_order_id;type:bigint(20);NOT NULL" json:"child_order_id"`           // 订单ID
	AlgorithmType uint         `gorm:"column:algorithm_type;type:smallint(6) unsigned;NOT NULL" json:"algorithm_type"` // 算法类型
	AlgorithmId   uint         `gorm:"column:algorithm_id;type:int(11) unsigned;NOT NULL" json:"algorithm_id"`         // 算法ID
	UsecurityId   uint         `gorm:"column:usecurity_id;type:int(11) unsigned;NOT NULL" json:"usecurity_id"`         // 证券ID
	SecurityId    string       `gorm:"column:security_id;type:varchar(8)" json:"security_id"`                          // 证券代码
	OrderQty      uint         `gorm:"column:order_qty;type:int(11) unsigned;NOT NULL" json:"order_qty"`               // 委托订单数量
	Price         uint         `gorm:"column:price;type:int(11) unsigned;NOT NULL" json:"price"`                       // 委托订单价格
	OrderType     uint         `gorm:"column:order_type;type:smallint(6) unsigned" json:"order_type"`                  // 订单类型 ：1-限价委托 2-本方最优 3-对手方最优 4-市价立即成交剩余撤销 5-市价全额成交或撤销 6-市价最优五档全额成交剩余撤销 7-限价全额成交或撤销(期权用）
	LastPx        uint         `gorm:"column:last_px;type:int(11) unsigned" json:"last_px"`                            // 成交价格
	LastQty       uint         `gorm:"column:last_qty;type:int(11) unsigned" json:"last_qty"`                          // 成交数量
	ComQty        uint         `gorm:"column:com_qty;type:int(11) unsigned" json:"com_qty"`                            // 累计成交数量
	ArrivedPrice  uint         `gorm:"column:arrived_price;type:int(11) unsigned" json:"arrived_price"`                // 到达价格
	OrdStatus     uint         `gorm:"column:ord_status;type:smallint(6) unsigned" json:"ord_status"`                  // 订单状态 1-新建 2-成交 3-撤销 4-拒绝
	TransactTime  uint         `gorm:"column:transact_time;type:int(11) unsigned;NOT NULL" json:"transact_time"`       // 交易时间
	TransactAt    uint64       `gorm:"column:transact_at;type:bigint(20) unsigned" json:"transact_at"`                 // 交易时间（精确到分钟）
	ProcStatus    uint         `gorm:"column:proc_status;type:smallint(6) unsigned" json:"proc_status"`                // 处理状态   0-未处理   1-已处理
	CreateTime    time.Time    `gorm:"column:create_time;type:timestamp" json:"create_time"` // 创建时间
}

func (m *TbAlgoOrderDetail) TableName() string {
	return "tb_algo_order_detail"
}

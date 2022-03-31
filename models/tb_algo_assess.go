// Package models
/*
 Author: hawrkchen
 Date: 2022/3/28 11:44
 Desc:
*/
package models

import (
	"time"
)

// 算法绩效信息表
type TbAlgoAssess struct {
	Id                    uint64       `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`        // 自增ID
	AlgorithmType         uint         `gorm:"column:algorithm_type;type:int(11) unsigned;NOT NULL" json:"algorithm_type"`     // 算法类型
	AlgorithmId           uint         `gorm:"column:algorithm_id;type:int(11) unsigned;NOT NULL" json:"algorithm_id"`         // 算法ID
	UsecurityId           uint         `gorm:"column:usecurity_id;type:int(11) unsigned;NOT NULL" json:"usecurity_id"`         // 证券ID
	SecurityId            string       `gorm:"column:security_id;type:varchar(8)" json:"security_id"`                          // 证券代码
	TimeDimension         int          `gorm:"column:time_dimension;type:tinyint(4);NOT NULL" json:"time_dimension"`           // 时间维度  1-秒 2-分 3-小时 4-天 5-周 6-月
	TransactTime          int          `gorm:"column:transact_time;type:int(11);NOT NULL" json:"transact_time"`                // 交易时间
	ArrivedPrice          int          `gorm:"column:arrived_price;type:int(11)" json:"arrived_price"`                         // 交易最新价格(到达价格)
	Vwap                  int          `gorm:"column:vwap;type:int(11)" json:"vwap"`                                           // 成交量加权平均价
	DealRate              int          `gorm:"column:deal_rate;type:int(11)" json:"deal_rate"`                                 // 成交进度
	OrderQty              uint         `gorm:"column:order_qty;type:int(11) unsigned" json:"order_qty"`                        // 委托订单数量
	LastQty               uint         `gorm:"column:last_qty;type:int(11) unsigned" json:"last_qty"`                          // 成交数量
	CancelQty             int          `gorm:"column:cancel_qty;type:int(11)" json:"cancel_qty"`                               // 撤销数量
	RejectedQty           int          `gorm:"column:rejected_qty;type:int(11)" json:"rejected_qty"`                           // 拒绝数量
	MarketRate            int          `gorm:"column:market_rate;type:int(11)" json:"market_rate"`                             // 市场参与率
	VwapDeviation         int          `gorm:"column:vwap_deviation;type:int(11)" json:"vwap_deviation"`                       // vwap滑点
	ArrivedPriceDeviation int          `gorm:"column:arrived_price_deviation;type:int(11)" json:"arrived_price_deviation"`     // 到达价滑点
	CreateTime            time.Time    `gorm:"column:create_time;type:timestamp" json:"create_time"` // 创建时间
}

func (m *TbAlgoAssess) TableName() string {
	return "tb_algo_assess"
}
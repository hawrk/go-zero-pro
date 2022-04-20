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
	Id                    uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`          // 自增ID
	AlgorithmType         uint      `gorm:"column:algorithm_type;type:int(10) unsigned;NOT NULL" json:"algorithm_type"`       // 算法类型
	AlgorithmId           uint      `gorm:"column:algorithm_id;type:int(10) unsigned;NOT NULL" json:"algorithm_id"`           // 算法ID
	UsecurityId           uint      `gorm:"column:usecurity_id;type:int(10) unsigned;NOT NULL" json:"usecurity_id"`           // 证券ID
	SecurityId            string    `gorm:"column:security_id;type:varchar(8)" json:"security_id"`                            // 证券代码
	TimeDimension         int       `gorm:"column:time_dimension;type:tinyint(4);NOT NULL" json:"time_dimension"`             // 时间维度  1-秒 2-分 3-小时 4-天 5-周 6-月
	TransactTime          int64     `gorm:"column:transact_time;type:bigint(20);NOT NULL" json:"transact_time"`               // 交易时间(精确到分)
	ArrivedPrice          int64     `gorm:"column:arrived_price;type:bigint(20)" json:"arrived_price"`                        // 交易最新价格(到达价格)
	LastPrice             int64     `gorm:"column:last_price;type:bigint(20)" json:"last_price"`                              // 最新价格
	Vwap                  float64   `gorm:"column:vwap;type:decimal(10,4)" json:"vwap"`                                       // 成交量加权平均价
	DealRate              float64   `gorm:"column:deal_rate;type:decimal(10,4)" json:"deal_rate"`                             // 成交量比重
	OrderQty              int64     `gorm:"column:order_qty;type:bigint(20)" json:"order_qty"`                                // 委托订单数量
	LastQty               int64     `gorm:"column:last_qty;type:bigint(20)" json:"last_qty"`                                  // 成交数量
	CancelQty             int64     `gorm:"column:cancel_qty;type:bigint(20)" json:"cancel_qty"`                              // 撤销数量
	RejectedQty           int64     `gorm:"column:rejected_qty;type:bigint(20)" json:"rejected_qty"`                          // 拒绝数量
	MarketRate            float64   `gorm:"column:market_rate;type:decimal(10,4)" json:"market_rate"`                         // 市场参与率
	DealProgress          float64   `gorm:"column:deal_progress;type:decimal(10,4)" json:"deal_progress"`                     // 成交进度
	VwapDeviation         float64   `gorm:"column:vwap_deviation;type:decimal(10,4)" json:"vwap_deviation"`                   // vwap滑点
	ArrivedPriceDeviation float64   `gorm:"column:arrived_price_deviation;type:decimal(10,4)" json:"arrived_price_deviation"` // 到达价滑点
	CreateTime            time.Time `gorm:"column:create_time;type:timestamp" json:"create_time"`                             // 创建时间
}

func (m *TbAlgoAssess) TableName() string {
	return "tb_algo_assess"
}

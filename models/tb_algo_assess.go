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
	Id                    uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:自增ID" json:"id"`
	BatchNo               int64     `gorm:"column:batch_no;type:bigint(20);comment:批次号" json:"batch_no"`
	AlgorithmType         uint      `gorm:"column:algorithm_type;type:int(10) unsigned;comment:算法类型;NOT NULL" json:"algorithm_type"`
	AlgorithmId           uint      `gorm:"column:algorithm_id;type:int(10) unsigned;comment:算法ID;NOT NULL" json:"algorithm_id"`
	UsecurityId           uint      `gorm:"column:usecurity_id;type:int(10) unsigned;comment:证券ID;NOT NULL" json:"usecurity_id"`
	SecurityId            string    `gorm:"column:security_id;type:varchar(12);comment:证券代码" json:"security_id"`
	UserId                string    `gorm:"column:user_id;type:varchar(45);comment:交易账户ID" json:"user_id"`
	TimeDimension         int       `gorm:"column:time_dimension;type:tinyint(4);comment:时间维度  1-秒 2-分 3-小时 4-天 5-周 6-月;NOT NULL" json:"time_dimension"`
	TransactTime          int64     `gorm:"column:transact_time;type:bigint(20);comment:交易时间(精确到分);NOT NULL" json:"transact_time"`
	ArrivedPrice          int64     `gorm:"column:arrived_price;type:bigint(20);comment:交易最新价格(到达价格)" json:"arrived_price"`
	LastPrice             int64     `gorm:"column:last_price;type:bigint(20);comment:最新价格" json:"last_price"`
	Vwap                  float64   `gorm:"column:vwap;type:decimal(20,4);comment:成交量加权平均价" json:"vwap"`
	DealRate              float64   `gorm:"column:deal_rate;type:decimal(20,4);comment:成交量比重" json:"deal_rate"`
	OrderQty              int64     `gorm:"column:order_qty;type:bigint(20);comment:委托订单数量" json:"order_qty"`
	LastQty               int64     `gorm:"column:last_qty;type:bigint(20);comment:成交数量" json:"last_qty"`
	CancelQty             int64     `gorm:"column:cancel_qty;type:bigint(20);comment:撤销数量" json:"cancel_qty"`
	RejectedQty           int64     `gorm:"column:rejected_qty;type:bigint(20);comment:拒绝数量" json:"rejected_qty"`
	MarketRate            float64   `gorm:"column:market_rate;type:decimal(20,4);comment:市场参与率" json:"market_rate"`
	DealProgress          float64   `gorm:"column:deal_progress;type:decimal(20,4);comment:成交进度" json:"deal_progress"`
	VwapDeviation         float64   `gorm:"column:vwap_deviation;type:decimal(20,4);comment:vwap滑点" json:"vwap_deviation"`
	ArrivedPriceDeviation float64   `gorm:"column:arrived_price_deviation;type:decimal(20,4);comment:到达价滑点" json:"arrived_price_deviation"`
	SourceFrom            int       `gorm:"column:source_from;type:tinyint(4);comment:数据来源，1:总线推送，2:数据导入" json:"source_from"`
	CreateTime            time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	UpdateTime            time.Time `gorm:"column:update_time;type:timestamp;default:null;comment:更新时间" json:"update_time"`
}

func (m *TbAlgoAssess) TableName() string {
	return "tb_algo_assess"
}

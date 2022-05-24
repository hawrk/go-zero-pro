// Package models
/*
 Author: hawrkchen
 Date: 2022/4/24 18:06
 Desc:
*/
package models

// 深市十档行情数据
type TbSzQuoteLevel struct {
	Id            uint64  `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"` // 自增ID
	SeculityId    string  `gorm:"column:seculity_id;type:varchar(12);NOT NULL" json:"seculity_id"`         // 证券ID
	OrgiTime      int64   `gorm:"column:orgi_time;type:bigint(20);NOT NULL" json:"orgi_time"`              // 快照时间
	LastPrice     int64   `gorm:"column:last_price;type:bigint(20)" json:"last_price"`                     // 最新价
	AskPrice      string  `gorm:"column:ask_price;type:varchar(256)" json:"ask_price"`                     // 申卖价
	AskVol        string  `gorm:"column:ask_vol;type:varchar(256)" json:"ask_vol"`                         // 申卖量
	BidPrice      string  `gorm:"column:bid_price;type:varchar(256)" json:"bid_price"`                     // 申买价
	BidVol        string  `gorm:"column:bid_vol;type:varchar(256)" json:"bid_vol"`                         // 申买量
	TotalTradeVol int64   `gorm:"column:total_trade_vol;type:bigint(20)" json:"total_trade_vol"`           // 成交总量
	TotalAskVol   int64   `gorm:"column:total_ask_vol;type:bigint(20)" json:"total_ask_vol"`
	TotalBidVol   int64   `gorm:"column:total_bid_vol;type:bigint(20)" json:"total_bid_vol"`
	MkVwap        float64 `gorm:"column:mk_vwap;type:decimal(20,4)" json:"mk_vwap"`
}

func (m *TbSzQuoteLevel) TableName() string {
	return "tb_sz_quote_level"
}

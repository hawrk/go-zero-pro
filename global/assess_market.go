// Package global
/*
 Author: hawrkchen
 Date: 2022/5/11 9:55
 Desc:
*/
package global

// QuoteLevel2Data 深沪市十档行情数据转换后的数据
type QuoteLevel2Data struct {
	SecID         string  `json:"sec_id"`          // 证券ID， 去掉空格
	OrigTime      int64   `json:"orig_time"`       // 快照时间  格式: 202205100930
	LastPrice     int64   `json:"last_price"`      // 最新价格 以分为单位
	TotalTradeVol int64   `json:"total_trade_vol"` // 成交总量     实际成交数量
	AskPrice      string  `json:"ask_price"`       // 申卖价，  以元为单位的数组
	AskVol        string  `json:"ask_vol"`         // 申卖量，  实时成交数量的数组
	BidPrice      string  `json:"bid_price"`       // 申买价， 以元为单位的数组
	BidVol        string  `json:"bid_vol"`         // 申买量， 实时成交数量的数组
	TotalBidVol   int64   `json:"total_bid_vol"`   // 委托买入总量， 实际数量
	TotalAskVol   int64   `json:"total_ask_vol"`   // 委托卖出总量， 实际数量
	Vwap          float64 `json:"vwap"`            // 市场vwap 需计算
}

// Package repo
/*
 Author: hawrkchen
 Date: 2022/6/23 10:09
 Desc:
*/
package repo

// AllAlgoInfo dashboard 所有算法信息
type AllAlgoInfo struct {
	AlgoType     int32
	AlgoTypeName string
	Count        int32
}

// UserOrderSummary dashboard里所有汇总信息
type UserOrderSummary struct {
	UserNum     int64 `json:"user_num"`     // 用户数量
	TradeAmount int64 `json:"trade_amount"` // 总交易量
	OrderNum    int64 `json:"order_num"`    // 订单数量
	ProviderNum int64 `json:"provider_num"` // 厂商数量
	EntrustQty  int64 `json:"entrust_qty"`  // 总委托数量
	DealQty     int64 `json:"deal_qty"`     // 总成交数量
}

// AlgoOrderSummary  dashboard里算法类型查交易汇总列表信息
type AlgoOrderSummary struct {
	Provider    string  `json:"provider"`     // 厂商名称
	UserNum     int64   `json:"user_num"`     // 用户数量
	TradeAmount int64   `json:"trade_amount"` // 交易总金额
	ProfitRate  float64 `json:"profit_rate"`  // 收益率
	OrderNum    int64   `json:"order_num"`    // 订单数量
	BuyAmount   int64   `json:"buy_amount"`   // 买入总金额
	SellAmount  int64   `json:"sell_amount"`  // 卖出总金额
}

// ProviderAccountInfo 算法厂商的账户信息，用来根据算法id 映射出对应的算法厂商的账户ID
type ProviderAccountInfo struct {
	AlgoId   int    `json:"algo_id"`  // 算法ID
	Provider string `json:"provider"` // 算法厂商名称
	UserId   string `json:"user_id"`  // 算法厂商账户
}

// AvgSummary 对比分析里，跨天处理，需要求平均值
type AvgSummary struct {
	AlgoId        int64   `json:"algo_id"`
	AlgoName      string  `json:"algo_name"`
	AssessScore   float64 `json:"assess_score"`
	ProgressScore float64 `json:"progress_score"`
	StableScore   float64 `json:"stable_score"`
	RiskScore     float64 `json:"risk_score"`
	EconomyScore  float64 `json:"economy_score"`
	CumsumScore   float64 `json:"cumsum_score"`
}

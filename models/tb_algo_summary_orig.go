// Package models
/*
 Author: hawrkchen
 Date: 2022/7/18 16:42
 Desc:
*/
package models

import (
	"time"
)

// 算法分析汇总表
type TbAlgoSummaryOrig struct {
	Id              int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:自增ID" json:"id"`
	Date            int       `gorm:"column:date;type:int(11);comment:日期 (格式:20221612)" json:"date"`
	BatchNo         int64     `gorm:"column:batch_no;type:bigint(20);comment:批次号" json:"batch_no"`
	UserId          string    `gorm:"column:user_id;type:varchar(45);comment:用户ID" json:"user_id"`
	AccountType     int       `gorm:"column:account_type;type:tinyint(4);comment:用户类型  1-普通用户，2-虚拟用户（汇总用户）" json:"account_type"`
	AlgoId          int       `gorm:"column:algo_id;type:int(11);comment:算法ID" json:"algo_id"`
	AlgoName        string    `gorm:"column:algo_name;type:varchar(45);comment:算法名称" json:"algo_name"`
	AlgoType        int       `gorm:"column:algo_type;type:smallint(6);comment:算法类型" json:"algo_type"`
	Provider        string    `gorm:"column:provider;type:varchar(45);comment:算法厂商" json:"provider"`
	OrderNum        int       `gorm:"column:order_num;type:int(11);comment:订单个数，以篮子为单位" json:"order_num"`
	EntrustQty      int64     `gorm:"column:entrust_qty;type:bigint(20);comment:总委托数量" json:"entrust_qty"`
	DealQty         int64     `gorm:"column:deal_qty;type:bigint(20);comment:总成交数量" json:"deal_qty"`
	OrderAmount     int64     `gorm:"column:order_amount;type:bigint(20);comment:总交易金额" json:"order_amount"`
	BuyAmount       int64     `gorm:"column:buy_amount;type:bigint(20);comment:买入总金额" json:"buy_amount"`
	SellAmount      int64     `gorm:"column:sell_amount;type:bigint(20);comment:卖出总金额" json:"sell_amount"`
	Perfit          float64   `gorm:"column:perfit;type:decimal(20,2);comment:盈亏金额" json:"perfit"`
	PerfitRate      float64   `gorm:"column:perfit_rate;type:decimal(20,2);comment:盈亏百分比" json:"perfit_rate"`
	AssessScore     int       `gorm:"column:assess_score;type:int(11);comment:算法绩效评分" json:"assess_score"`
	ProgressScore   int       `gorm:"column:progress_score;type:int(11);comment:完成度评分" json:"progress_score"`
	StableScore     int       `gorm:"column:stable_score;type:int(11);comment:稳定性评分" json:"stable_score"`
	RiskScore       int       `gorm:"column:risk_score;type:int(11);comment:风险度评分" json:"risk_score"`
	EconomyScore    int       `gorm:"column:economy_score;type:int(11);comment:经济性评分" json:"economy_score"`
	CumsumScore     int       `gorm:"column:cumsum_score;type:int(11);comment:综合评分" json:"cumsum_score"`
	FundRateJson    string    `gorm:"column:fund_rate_json;type:varchar(1024);comment:资金占比Json串" json:"fund_rate_json"`
	TradeVolJson    string    `gorm:"column:trade_vol_json;type:varchar(1024);comment:交易量占比" json:"trade_vol_json"`
	StockTypeJson   string    `gorm:"column:stock_type_json;type:varchar(1024);comment:股价类型Json串" json:"stock_type_json"`
	TradeDirectJson string    `gorm:"column:trade_direct_json;type:varchar(1024);comment:买卖方向json" json:"trade_direct_json"`
	SourceFrom      int       `gorm:"column:source_from;type:tinyint(4);comment:数据来源，1:总线推送，2:数据导入" json:"source_from"`
	CreateTime      time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	UpdateTime      time.Time `gorm:"column:update_time;type:timestamp;default:null;comment:更新时间" json:"update_time"`
}

func (m *TbAlgoSummaryOrig) TableName() string {
	return "tb_algo_summary_orig"
}

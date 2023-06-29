// Package models
/*
 Author: hawrkchen
 Date: 2022/6/23 11:20
 Desc:
*/
package models

import (
	"time"
)

type TbAlgoProfileOrig struct {
	Id                int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:自增ID" json:"id"`
	Date              int       `gorm:"column:date;type:int(11);comment:日期 (20220612);NOT NULL" json:"date"`
	BatchNo           int64     `gorm:"column:batch_no;type:bigint(20);comment:批次号" json:"batch_no"`
	AccountId         string    `gorm:"column:account_id;type:varchar(45);comment:用户ID （交易账户ID）;NOT NULL" json:"account_id"`
	AccountName       string    `gorm:"column:account_name;type:varchar(45);comment:账户名称" json:"account_name"`
	AccountType       int       `gorm:"column:account_type;type:tinyint(4);comment:用户类型  1-普通用户，2-虚拟用户（汇总用户）" json:"account_type"`
	Provider          string    `gorm:"column:provider;type:varchar(45);comment:厂商名称" json:"provider"`
	AlgoId            int       `gorm:"column:algo_id;type:int(11);comment:算法ID" json:"algo_id"`
	AlgoName          string    `gorm:"column:algo_name;type:varchar(45);comment:算法名称" json:"algo_name"`
	AlgoType          int       `gorm:"column:algo_type;type:int(11);comment:算法类型" json:"algo_type"`
	SecId             string    `gorm:"column:sec_id;type:varchar(12);comment:证券代码" json:"sec_id"`
	SecName           string    `gorm:"column:sec_name;type:varchar(45);comment:证券名称" json:"sec_name"`
	AlgoOrderId       int64     `gorm:"column:algo_order_id;type:bigint(20);comment:母单ID" json:"algo_order_id"`
	Industry          string    `gorm:"column:industry;type:varchar(45);comment:行业类型" json:"industry"`
	FundType          int       `gorm:"column:fund_type;type:tinyint(4);comment:市值类型    1- huge超大 2-big大 3-middle中等 4-small小" json:"fund_type"`
	Liquidity         int       `gorm:"column:liquidity;type:tinyint(4);comment:流动性   1-高 2-中 3-低" json:"liquidity"`
	TradeCost         int64     `gorm:"column:trade_cost;type:bigint(20);comment:交易成本（ 成交价格* 成交数量）" json:"trade_cost"`
	TotalTradeAmount  int64     `gorm:"column:total_trade_amount;type:bigint(20);comment:双边总交易额， 总金额" json:"total_trade_amount"`
	TotalTradeFee     int64     `gorm:"column:total_trade_fee;type:bigint(20);comment:总手续费（券商手续旨，过户费，印花税）" json:"total_trade_fee"`
	CrossFee          int64     `gorm:"column:cross_fee;type:bigint(20);comment:流量费" json:"cross_fee"`
	AvgTradePrice     float64   `gorm:"column:avg_trade_price;type:decimal(20,4);comment:母单执行均价（成交均价)" json:"avg_trade_price"`
	AvgArrivePrice    float64   `gorm:"column:avg_arrive_price;type:decimal(20,4);comment:到达价均价" json:"avg_arrive_price"`
	Pwp               float64   `gorm:"column:pwp;type:decimal(20,4);comment:pwp价格" json:"pwp"`
	AlgoDuration      int64     `gorm:"column:algo_duration;type:bigint(20);comment:母单有效时长" json:"algo_duration"`
	Twap              float64   `gorm:"column:twap;type:decimal(20,4);comment:TWAP" json:"twap"`
	TwapDev           float64   `gorm:"column:twap_dev;type:decimal(20,4);comment:TWAP滑点" json:"twap_dev"`
	Vwap              float64   `gorm:"column:vwap;type:decimal(20,4);comment:vwap值" json:"vwap"`
	VwapDev           float64   `gorm:"column:vwap_dev;type:decimal(20,4);comment:vwap 滑点" json:"vwap_dev"`
	ProfitAmount      float64   `gorm:"column:profit_amount;type:decimal(20,4);comment:盈亏金额" json:"profit_amount"`
	ProfitRate        float64   `gorm:"column:profit_rate;type:decimal(20,4);comment:收益率， 盈亏比率" json:"profit_rate"`
	CancelRate        float64   `gorm:"column:cancel_rate;type:decimal(20,4);comment:撤单率" json:"cancel_rate"`
	ProgressRate      float64   `gorm:"column:progress_rate;type:decimal(20,4);comment:完成度" json:"progress_rate"`
	MiniSplitOrder    int       `gorm:"column:mini_split_order;type:int(11);comment:最小拆单单位" json:"mini_split_order"`
	MiniJointRate     float64   `gorm:"column:mini_joint_rate;type:decimal(20,4);comment:最小贴合度" json:"mini_joint_rate"`
	WithdrawRate      float64   `gorm:"column:withdraw_rate;type:decimal(20,4);comment:回撤比例" json:"withdraw_rate"`
	DealEffi          float64   `gorm:"column:deal_effi;type:decimal(20);comment:成交效率" json:"deal_effi"`
	AlgoOrderFit      float64   `gorm:"column:algo_order_fit;type:decimal(20);comment:母单贴合度" json:"algo_order_fit"`
	TradeVolFit       float64   `gorm:"column:trade_vol_fit;type:decimal(20);comment:成交量贴合度" json:"trade_vol_fit"`
	TradeCount        int       `gorm:"column:trade_count;type:int(11);comment:交易次数（一个子单交易回执算一次交易次数，被 拒绝或撤销的不算， 用来算滑点平均值）" json:"trade_count"`
	TradeCountPlus    int       `gorm:"column:trade_count_plus;type:int(11);comment:盈亏为正的交易次数，统计胜率" json:"trade_count_plus"`
	AvgDeviation      float64   `gorm:"column:avg_deviation;type:decimal(20,4);comment:滑点平均值" json:"avg_deviation"`
	StandardDeviation float64   `gorm:"column:standard_deviation;type:decimal(20,4);comment:滑点标准差" json:"standard_deviation"`
	PfRateStdDev      float64   `gorm:"column:pf_rate_std_dev;type:decimal(20,4);comment:收益率标准差" json:"pf_rate_std_dev"`
	Factor            float64   `gorm:"column:factor;type:decimal(20,4);comment:绩效收益率" json:"factor"`
	DealFitStdDev     float64   `gorm:"column:deal_fit_std_dev;type:decimal(20);comment:成交量贴合度标准差" json:"deal_fit_std_dev"`
	TimeFitStdDev     float64   `gorm:"column:time_fit_std_dev;type:decimal(20);comment:时间贴合度标准差" json:"time_fit_std_dev"`
	OrderTime         int64     `gorm:"column:order_time;type:bigint(20);comment:订单开始时间（母单开始时间)" json:"order_time"`
	SourceFrom        int       `gorm:"column:source_from;type:tinyint(4);comment:数据来源，1:总线推送，2:数据导入" json:"source_from"`
	CreateTime        time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间;NOT NULL" json:"create_time"`
	UpdateTime        time.Time `gorm:"column:update_time;type:timestamp;default:null;comment:更新时间" json:"update_time"`
}

func (m *TbAlgoProfileOrig) TableName() string {
	return "tb_algo_profile_orig"
}

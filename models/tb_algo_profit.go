// Package models
/*
 Author: hawrkchen
 Date: 2022/10/9 10:22
 Desc:
*/
package models

import (
	"time"
)

// 计算股票算法的盈亏信息
type TbAlgoProfit struct {
	Id           int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:自增ID" json:"id"`
	Date         int       `gorm:"column:date;type:int(12);comment:计算日期(20220612);NOT NULL" json:"date"`
	BatchNo      int64     `gorm:"column:batch_no;type:bigint(20);comment:批次号" json:"batch_no"`
	TransactTime int64     `gorm:"column:transact_time;type:bigint(20);comment:交易时间，精确到分(202206120955)" json:"transact_time"`
	AccountId    string    `gorm:"column:account_id;type:varchar(45);comment:用户ID(交易账户ID)" json:"account_id"`
	AccountName  string    `gorm:"column:account_name;type:varchar(45);comment:账户名称" json:"account_name"`
	AccountType  int       `gorm:"column:account_type;type:tinyint(4);comment:账户类型" json:"account_type"`
	Provider     string    `gorm:"column:provider;type:varchar(45);comment:算法厂商名称" json:"provider"`
	AlgoId       int       `gorm:"column:algo_id;type:int(8);comment:算法ID" json:"algo_id"`
	AlgoName     string    `gorm:"column:algo_name;type:varchar(45);comment:算法名称" json:"algo_name"`
	AlgoType     int       `gorm:"column:algo_type;type:tinyint(4)" json:"algo_type"`
	SecId        string    `gorm:"column:sec_id;type:varchar(45);comment:证券ID" json:"sec_id"`
	SecName      string    `gorm:"column:sec_name;type:varchar(45);comment:证券名称" json:"sec_name"`
	TradeCost    int64     `gorm:"column:trade_cost;type:bigint(20);comment:交易成本" json:"trade_cost"`
	TradeFee     int64     `gorm:"column:trade_fee;type:bigint(20);comment:总手续费" json:"trade_fee"`
	CrossFee     int64     `gorm:"column:cross_fee;type:bigint(20);comment:流量费" json:"cross_fee"`
	ProfitAmount int64     `gorm:"column:profit_amount;type:bigint(20);comment:盈亏金额" json:"profit_amount"`
	ProfitRate   float64   `gorm:"column:profit_rate;type:decimal(20,4);comment:收益率，盈亏比例" json:"profit_rate"`
	VwapStdDev   float64   `gorm:"column:vwap_std_dev;type:decimal(20,4);comment:vwap 滑点标准差" json:"vwap_std_dev"`
	PfRateStdDev float64   `gorm:"column:pf_rate_std_dev;type:decimal(20,4);comment:收益率标准差" json:"pf_rate_std_dev"`
	SourceFrom   int       `gorm:"column:source_from;type:tinyint(4);comment:数据来源，1:总线推送，2:数据导入" json:"source_from"`
	CreateTime   time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;type:timestamp;default:null;comment:更新时间" json:"update_time"`
}

func (m *TbAlgoProfit) TableName() string {
	return "tb_algo_profit"
}

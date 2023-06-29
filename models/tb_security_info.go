// Package models
/*
 Author: hawrkchen
 Date: 2022/7/6 15:50
 Desc:
*/
package models

import (
	"time"
)

// 证券基础信息表
type TbSecurityInfo struct {
	Id              int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:自增ID" json:"id"`
	SecurityId      string    `gorm:"column:security_id;type:varchar(8);comment:证券ID;NOT NULL" json:"security_id"`
	SecuritySource  string    `gorm:"column:security_source;type:varchar(4);comment:证券代码源  102-深市，103 沪市" json:"security_source"`
	SecurityName    string    `gorm:"column:security_name;type:varchar(45);comment:证券名称" json:"security_name"`
	PreClosePx      float64   `gorm:"column:pre_close_px;type:double;comment:前收盘价" json:"pre_close_px"`
	Status          int       `gorm:"column:status;type:tinyint(4);comment:证券状态0-正常 1-停盘 2-退市" json:"status"`
	IsPriceLimit    int       `gorm:"column:is_price_limit;type:tinyint(4);comment:涨停限制  -0-有 1-无" json:"is_price_limit"`
	LimtType        int       `gorm:"column:limt_type;type:tinyint(4);comment:涨跌限制类型 1-幅度 2-价格" json:"limt_type"`
	Property        int       `gorm:"column:property;type:tinyint(4);comment:板块属性" json:"property"`
	UpperLimitPrice int64     `gorm:"column:upper_limit_price;type:bigint(20);comment:涨停价格" json:"upper_limit_price"`
	LowerLimitPrice int64     `gorm:"column:lower_limit_price;type:bigint(20);comment:跌停价格" json:"lower_limit_price"`
	FundType        int       `gorm:"column:fund_type;type:tinyint(4);comment:市值类型    1- huge超大 2-big大 3-middle中等 4-small小" json:"fund_type"`
	StockType       int       `gorm:"column:stock_type;type:tinyint(4);comment:股价类型" json:"stock_type"`
	Liquidity       int       `gorm:"column:liquidity;type:tinyint(4);comment:流动性   1-高 2-中 3-低" json:"liquidity"`
	Industry        string    `gorm:"column:industry;type:varchar(45);comment:行业类型" json:"industry"`
	CreateTime      time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	UpdateTime      time.Time `gorm:"column:update_time;type:timestamp;default:null;comment:更新时间" json:"update_time"`
}

func (m *TbSecurityInfo) TableName() string {
	return "tb_security_info"
}

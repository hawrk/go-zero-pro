// Package busmodels
/*
 Author: hawrkchen
 Date: 2023/3/23 17:09
 Desc:
*/
package busmodels

// 证券信息表(MyISAM)
type TbSecurityInfo struct {
	Id                      uint    `gorm:"column:id;type:int(10) unsigned;primary_key;AUTO_INCREMENT;comment:[AUTO_INCREMENT]自增Id" json:"id"`
	SecurityId              string  `gorm:"column:security_id;type:varchar(8);comment:证券代码" json:"security_id"`
	SecurityIdSource        string  `gorm:"column:security_id_source;type:varchar(4);comment:证券代码源" json:"security_id_source"`
	SecurityName            string  `gorm:"column:security_name;type:varchar(40);comment:证券简称" json:"security_name"`
	PrevClosePx             float64 `gorm:"column:prev_close_px;type:double;comment:前收盘价" json:"prev_close_px"`
	BuyQtyUpperLimit        uint64  `gorm:"column:buy_qty_upper_limit;type:bigint(20) unsigned;comment:限价买数量上限" json:"buy_qty_upper_limit"`
	SellQtyUpperLimit       uint64  `gorm:"column:sell_qty_upper_limit;type:bigint(20) unsigned;comment:限价卖数量上限" json:"sell_qty_upper_limit"`
	MarketBuyQtyUpperLimit  uint64  `gorm:"column:market_buy_qty_upper_limit;type:bigint(20) unsigned;comment:市价买数量上限" json:"market_buy_qty_upper_limit"`
	MarketSellQtyUpperLimit uint64  `gorm:"column:market_sell_qty_upper_limit;type:bigint(20) unsigned;comment:市价卖数量上限" json:"market_sell_qty_upper_limit"`
	SecurityStatus          string  `gorm:"column:security_status;type:char(1);comment:证券状态;0正常, 1停盘, 2退市, 4-ST, 5-*ST,6-上市次日到5日(无涨跌幅), 7-删除" json:"security_status"`
	HasPriceLimit           string  `gorm:"column:has_price_limit;type:char(1);comment:是否有涨跌停价格限制,Y=是,N=否" json:"has_price_limit"`
	LimitType               string  `gorm:"column:limit_type;type:char(1);comment:涨跌限制类型,1=幅度(百分比),2=价格(绝对值)" json:"limit_type"`
	Property                string  `gorm:"column:property;type:char(1);comment:股票板块属性 1-主板 3-创业板" json:"property"`
	UpperLimitPrice         uint64  `gorm:"column:upper_limit_price;type:bigint(20) unsigned;comment:上涨限价" json:"upper_limit_price"`
	LowerLimitPrice         uint64  `gorm:"column:lower_limit_price;type:bigint(20) unsigned;comment:下跌限价" json:"lower_limit_price"`
	BuyQtyUnit              uint    `gorm:"column:buy_qty_unit;type:int(10) unsigned;comment:限价单位买数量" json:"buy_qty_unit"`
	SellQtyUnit             uint    `gorm:"column:sell_qty_unit;type:int(10) unsigned;comment:限价单位卖数量" json:"sell_qty_unit"`
	MarketBuyQtyUnit        uint    `gorm:"column:market_buy_qty_unit;type:int(10) unsigned;comment:市价单位买数量" json:"market_buy_qty_unit"`
	MarketSellQtyUnit       uint    `gorm:"column:market_sell_qty_unit;type:int(10) unsigned;comment:市价单位卖数量" json:"market_sell_qty_unit"`
	UpdateTime              uint64  `gorm:"column:update_time;type:bigint(20) unsigned;comment:更新时间" json:"update_time"`
	Version                 uint    `gorm:"column:version;type:int(10) unsigned;comment:版本号" json:"version"`
	CreditType              int     `gorm:"column:credit_type;type:tinyint(4);comment:信用类型 bit0-可融资 bit1-可融券 bit2-可作担保品 bit3-新股" json:"credit_type"`
	SecurityType            uint    `gorm:"column:security_type;type:smallint(6) unsigned;comment:证券类型: 0x00未知 0x01指数 0x02股票 0x04基金 0x08债券 0x10权证 0x20期权" json:"security_type"`
}

func (m *TbSecurityInfo) TableName() string {
	return "tb_security_info"
}

// Package logic
/*
 Author: hawrkchen
 Date: 2022/9/9 11:02
 Desc: 导入导出文件字段模板
*/
package logic

import "encoding/xml"

// ------------------证券信息导入导出模板  ----------------//

type XmlSecurity struct {
	XmlName      xml.Name       `xml:"SecurityInfo"`
	SecurityData []SecurityData `xml:"SecurityData"`
}

type SecurityData struct {
	SecurityId              string  `json:"security_id" xml:"SecurityId"`
	SecuritySource          string  `json:"security_source" xml:"SecurityIdSource"`
	SecurityName            string  `json:"security_name" xml:"SecurityName"`
	PreClosePx              float64 `json:"pre_close_px" xml:"PrevClosePx"`
	BuyQtyUpperLimit        int64   `json:"buy_qty_upper_limit" xml:"BuyQtyUpperLimit"`
	SellQtyUpperLimit       int64   `json:"sell_qty_upper_limit" xml:"SellQtyUpperLimit"`
	MarketBuyQtyUpperLimit  int64   `json:"market_buy_qty_upper_limit" xml:"MarketBuyQtyUpperLimit"`
	MarketSellQtyUpperLimit int64   `json:"market_sell_qty_upper_limit" xml:"MarketSellQtyUpperLimit"`
	SecurityStatus          int32   `json:"security_status" xml:"SecurityStatus"`
	HasPriceLimit           int32   `json:"has_price_limit" xml:"HasPriceLimit"`
	LimitType               int32   `json:"lim+it_type" xml:"LimitType"`
	Property                int32   `json:"property" xml:"Property"`
	UpperLimitPrice         int64   `json:"upper_limit_price" xml:"UpperLimitPrice"`
	LowerLimitPrice         int64   `json:"lower_limit_price" xml:"LowerLimitPrice"`
	BuyQtyUnit              int64   `json:"buy_qty_unit" xml:"BuyQtyUnit"`
	SellQtyUnit             int64   `json:"sell_qty_unit" xml:"SellQtyUnit"`
	MarketBuyQtyUnit        int64   `json:"market_buy_qty_unit" xml:"MarketBuyQtyUnit"`
	MarketSellQtyUnit       int64   `json:"market_sell_qty_unit" xml:"MarketSellQtyUnit"`
	FundType                int32   `json:"fund_type" xml:"FundType"`   // 绩效平台新增
	StockType               int32   `json:"stock_type" xml:"StockType"` // 绩效平台新增
	Liquidity               int32   `json:"liquidity" xml:"Liquidity"`  // 绩效平台新增
	Industry                string  `json:"industry" xml:"Industry"`
}

// ----------------- 用户信息导入导出模板   -------------//

type XmlUserInfo struct {
	XmlName  xml.Name   `xml:"UserInfo"`
	UserInfo []UserData `xml:"UserData"`
}

type UserData struct {
	UserId           string `json:"user_id" xml:"UserId"`
	UserName         string `json:"user_name" xml:"UserName"`
	UserPasswd       string `json:"user_passwd" xml:"UserPasswd"`
	PasswdEnCodeType int32  `json:"passwd_en_code_type" xml:"PasswdEnCodeType"`
	UserType         int32  `json:"user_type" xml:"UserType"`
	RiskGroup        int32  `json:"risk_group" xml:"RiskGroup"`
	UuserId          string `json:"uuser_id" xml:"UuserId"`
	UserGrade        string `json:"user_grade" xml:"UserGrade"` // 绩效平台新增
}

// Package busmodels
/*
 Author: hawrkchen
 Date: 2023/3/23 15:46
 Desc:
*/
package busmodels

// 资金信息表(MyISAM)
type TbAssetInfo struct {
	Id           uint    `gorm:"column:id;type:int(10) unsigned;primary_key;comment:用户ID" json:"id"`
	AssetAccount string  `gorm:"column:asset_account;type:varchar(32);comment:资金账户" json:"asset_account"`
	Balance      float64 `gorm:"column:balance;type:double;comment:日间余额" json:"balance"`
	Frozen       float64 `gorm:"column:frozen;type:double;comment:冻结资金" json:"frozen"`
	MarginAmount float64 `gorm:"column:margin_amount;type:double;comment:实时保证金" json:"margin_amount"`
	CreateTime   uint64  `gorm:"column:create_time;type:bigint(20) unsigned;comment:注册时间" json:"create_time"`
	UpdateTime   uint64  `gorm:"column:update_time;type:bigint(20) unsigned;comment:更新时间" json:"update_time"`
	Version      uint    `gorm:"column:version;type:int(10) unsigned;comment:版本号" json:"version"`
	CurrencyType uint    `gorm:"column:currency_type;type:tinyint(4) unsigned;comment:币种:1 人民币 2 港币 3 美元" json:"currency_type"`
	AccountType  uint    `gorm:"column:account_type;type:tinyint(4) unsigned;comment:账户类型:1股票 2期权 3融资融券 同BusinessType" json:"account_type"`
	CustOrgid    string  `gorm:"column:cust_orgid;type:varchar(16);comment:机构编码" json:"cust_orgid"`
	CustBranchid string  `gorm:"column:cust_branchid;type:varchar(16);comment:分支编码" json:"cust_branchid"`
}

func (m *TbAssetInfo) TableName() string {
	return "tb_asset_info"
}

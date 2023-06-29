// Package busmodels
/*
 Author: hawrkchen
 Date: 2023/3/23 15:48
 Desc:
*/
package busmodels

// 用户持仓表
type TbUserPosition struct {
	Id               uint    `gorm:"column:id;type:int(10) unsigned;primary_key;AUTO_INCREMENT;comment:[AUTO_INCREMENT]自增Id" json:"id"`
	UuserId          uint    `gorm:"column:uuser_id;type:int(10) unsigned;comment:用户ID" json:"uuser_id"`
	AccountId        string  `gorm:"column:account_id;type:varchar(12);comment:股东账户" json:"account_id"`
	SecurityId       string  `gorm:"column:security_id;type:varchar(8);comment:证券代码" json:"security_id"`
	SecurityIdSource string  `gorm:"column:security_id_source;type:varchar(4);comment:证券代码源" json:"security_id_source"`
	PositionType     uint    `gorm:"column:position_type;type:tinyint(4) unsigned;comment:持仓方向: 0现货 1权利持仓(期权) 2义务持仓(期权) 3备兑持仓(期权) 4两融多仓 5两融空仓" json:"position_type"`
	PositionQty      int64   `gorm:"column:position_qty;type:bigint(20);comment:持仓总量x100" json:"position_qty"`
	OriginQty        int64   `gorm:"column:origin_qty;type:bigint(20);comment:当前开盘前的原始仓位数量x100" json:"origin_qty"`
	OriginOpenPrice  int64   `gorm:"column:origin_open_price;type:bigint(20);comment:当天前的原始持仓的平均开仓价格x10000" json:"origin_open_price"`
	FreeQty          int64   `gorm:"column:free_qty;type:bigint(20);comment:可卖仓位数量x100" json:"free_qty"`
	FrozenQty        int64   `gorm:"column:frozen_qty;type:bigint(20);comment:冻结数量x100" json:"frozen_qty"`
	AvgPrice         int64   `gorm:"column:avg_price;type:bigint(20);comment:平均价格x10000" json:"avg_price"`
	ProfitAndLoss    float64 `gorm:"column:profit_and_loss;type:double;comment:盈亏" json:"profit_and_loss"`
	ProfitRate       float64 `gorm:"column:profit_rate;type:double;comment:收益率" json:"profit_rate"`
	PositionRatio    float64 `gorm:"column:position_ratio;type:double;comment:持仓比例" json:"position_ratio"`
	Version          uint    `gorm:"column:version;type:int(10) unsigned;comment:版本号" json:"version"`
	UpdateTime       uint64  `gorm:"column:update_time;type:bigint(20) unsigned;comment:更新时间" json:"update_time"`
	CounterVersion   uint    `gorm:"column:counter_version;type:int(10) unsigned;comment:柜台版本号" json:"counter_version"`
}

func (m *TbUserPosition) TableName() string {
	return "tb_user_position"
}

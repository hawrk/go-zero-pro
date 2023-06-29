// Package busmodels
/*
 Author: hawrkchen
 Date: 2022/9/27 10:16
 Desc:
*/
package busmodels

// TbUserInfo 用户信息表
type TbUserInfo struct {
	Id           uint   `gorm:"column:id;type:int(10) unsigned;primary_key;AUTO_INCREMENT;comment:[AUTO_INCREMENT]UUserId,用户[1]->股东账户[N]" json:"id"`
	UserId       string `gorm:"column:user_id;type:varchar(12);comment:用户ID" json:"user_id"`
	UserName     string `gorm:"column:user_name;type:varchar(32);comment:用户名" json:"user_name"`
	UserPasswd   string `gorm:"column:user_passwd;type:varchar(32);comment:用户密码" json:"user_passwd"`
	UserType     uint   `gorm:"column:user_type;type:tinyint(4) unsigned;comment:用户类型: 1个人用户 2算法厂商用户 3多用户管理员" json:"user_type"`
	RiskGroup    uint   `gorm:"column:risk_group;type:int(10) unsigned;comment:用户风控组" json:"risk_group"`
	UuserId      uint   `gorm:"column:uuser_id;type:int(10) unsigned;comment:自关联此表ID: 多用户管理员[1]->个人用户[N], 个人用户[1]->资金账号[1], 父级没有时为0" json:"uuser_id"`
	CreateTime   uint64 `gorm:"column:create_time;type:bigint(20) unsigned;comment:创建时间" json:"create_time"`
	Version      uint   `gorm:"column:version;type:int(10) unsigned;comment:版本号" json:"version"`
	AlgoGroup    uint   `gorm:"column:algo_group;type:int(10) unsigned;comment:算法可用组 对应表AlgoGroupInfoData->Id" json:"algo_group"`
	AlgoProperty string `gorm:"column:algo_property;type:varchar(32);comment:算法属性,使用二进制位标识权限:0不可用 1可用" json:"algo_property"`
}

func (m *TbUserInfo) TableName() string {
	return "tb_user_info"
}

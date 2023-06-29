// Package busmodels
/*
 Author: hawrkchen
 Date: 2022/9/27 10:07
 Desc:
*/
package busmodels

// TbAlgoGroupInfo 算法可用组信息表
type TbAlgoGroupInfo struct {
	Id           uint   `gorm:"column:id;type:int(10) unsigned;primary_key;AUTO_INCREMENT;comment:[AUTO_INCREMENT]" json:"id"`
	AlgoProperty string `gorm:"column:algo_property;type:varchar(32);comment:算法属性,使用二进制位标识权限:0不可用 1可用" json:"algo_property"`
	CreateTime   uint64 `gorm:"column:create_time;type:bigint(20) unsigned;comment:创建时间" json:"create_time"`
	UpdateTime   uint64 `gorm:"column:update_time;type:bigint(20) unsigned;comment:更新时间" json:"update_time"`
	GroupName    string `gorm:"column:group_name;type:varchar(32);comment:算法组名" json:"group_name"`
}

func (m *TbAlgoGroupInfo) TableName() string {
	return "tb_algo_group_info"
}

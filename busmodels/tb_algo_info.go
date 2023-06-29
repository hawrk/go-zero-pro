// Package busmodels
/*
 Author: hawrkchen
 Date: 2022/9/27 10:09
 Desc:
*/
package busmodels

// TbAlgoInfo 算法信息表(MyISAM)
type TbAlgoInfo struct {
	Id                uint   `gorm:"column:id;type:int(10) unsigned;primary_key;comment:算法ID, 主键" json:"id"`
	AlgoName          string `gorm:"column:algo_name;type:varchar(32);comment:算法名称" json:"algo_name"`
	ProviderName      string `gorm:"column:provider_name;type:varchar(32);comment:算法厂商名称" json:"provider_name"`
	UuserId           uint   `gorm:"column:uuser_id;type:int(10) unsigned;comment:算法厂商用户的ID" json:"uuser_id"`
	AlgorithmType     uint   `gorm:"column:algorithm_type;type:smallint(6) unsigned;comment:算法类型, 1:T0日内回转 2:智能委托 3:调仓" json:"algorithm_type"`
	AlgorithmTypeName string `gorm:"column:algorithm_type_name;type:varchar(16);comment:算法类型, 算法厂商内部使用" json:"algorithm_type_name"`
	AlgorithmStatus   uint   `gorm:"column:algorithm_status;type:tinyint(4) unsigned;comment:算法状态 —— bit0: 0-不显示 1-显示; bit1: 0-不可用 1-可用" json:"algorithm_status"`
	Parameter         string `gorm:"column:parameter;type:varchar(2048);comment:算法所需参数" json:"parameter"`
	RiskGroup         uint   `gorm:"column:risk_group;type:int(10) unsigned;comment:算法风控组" json:"risk_group"`
	CreateTime        uint64 `gorm:"column:create_time;type:bigint(20) unsigned;comment:创建时间戳,单位:秒" json:"create_time"`
	Version           uint   `gorm:"column:version;type:int(10) unsigned;comment:成交记录ID" json:"version"`
}

func (m *TbAlgoInfo) TableName() string {
	return "tb_algo_info"
}

// Package model
/*
 Author: hawrkchen
 Date: 2022/3/22 14:35
 Desc:
*/
package models

// 算法信息表(MyISAM)
type TbAlgoInfo struct {
	Id                uint   `gorm:"column:id;type:int(11) unsigned;primary_key" json:"id"`                    // 算法ID, 主键
	AlgoName          string `gorm:"column:algo_name;type:varchar(32)" json:"algo_name"`                       // 算法名称
	Provider          uint   `gorm:"column:provider;type:tinyint(4) unsigned" json:"provider"`                 // 算法厂商ID
	ProviderName      string `gorm:"column:provider_name;type:varchar(32)" json:"provider_name"`               // 算法厂商名称
	UuserId           uint   `gorm:"column:uuser_id;type:int(11) unsigned" json:"uuser_id"`                    // 算法厂商用户的ID
	AlgorithmType     uint   `gorm:"column:algorithm_type;type:smallint(6) unsigned" json:"algorithm_type"`    // 算法类型, 算法厂商内部唯一
	AlgorithmTypeName string `gorm:"column:algorithm_type_name;type:varchar(16)" json:"algorithm_type_name"`   // 算法类型, 算法厂商内部使用
	AlgorithmStatus   uint   `gorm:"column:algorithm_status;type:tinyint(4) unsigned" json:"algorithm_status"` // 算法状态: 0可用 1不可用
	Parameter         string `gorm:"column:parameter;type:varchar(2048)" json:"parameter"`                     // 算法所需参数
	RiskGroup         uint   `gorm:"column:risk_group;type:int(11) unsigned" json:"risk_group"`                // 算法风控组
	CreateTime        uint64 `gorm:"column:create_time;type:bigint(20) unsigned" json:"create_time"`           // 创建时间戳,单位:秒
	Version           uint   `gorm:"column:version;type:int(11) unsigned" json:"version"`                      // 成交记录ID
}

func (m *TbAlgoInfo) TableName() string {
	return "tb_algo_info"
}
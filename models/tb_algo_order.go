// Package models
/*
 Author: hawrkchen
 Date: 2022/4/26 10:54
 Desc:
*/
package models

// 母单落地表
type TbAlgoOrder struct {
	Id           int64  `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"` // 自增ID
	AlgoId       int    `gorm:"column:algo_id;type:int(11);NOT NULL" json:"algo_id"`            // 母单ID
	AlgorithmId  int    `gorm:"column:algorithm_id;type:int(11)" json:"algorithm_id"`           // 算法ID
	UsecId       int    `gorm:"column:usec_id;type:int(11)" json:"usec_id"`                     // 证券代码
	SecId        string `gorm:"column:sec_id;type:varchar(8)" json:"sec_id"`                    // 证券ID
	AlgoOrderQty int64  `gorm:"column:algo_order_qty;type:bigint(20)" json:"algo_order_qty"`    // 订单数量
	TransTime    int64  `gorm:"column:trans_time;type:bigint(20)" json:"trans_time"`            // 交易时间
}

func (m *TbAlgoOrder) TableName() string {
	return "tb_algo_order"
}

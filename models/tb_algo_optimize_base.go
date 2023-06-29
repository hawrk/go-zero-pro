// Package models
/*
 Author: yuzh
 Date: 2022/7/28
 Desc:
*/
package models

import (
	"time"
)

// TbAlgoOptimizeBase 一键优选基础信息表
type TbAlgoOptimizeBase struct {
	Id           int64     `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT;comment:主键id" json:"id"`
	ProviderId   int       `gorm:"column:provider_id;type:int;comment:算法厂商id;NOT NULL" json:"provider_id"`
	ProviderName string    `gorm:"column:provider_name;type:varchar(32);comment:算法厂商名称;" json:"provider_name"`
	SecId        string    `gorm:"column:sec_id;type:varchar(12);comment:证券ID;NOT NULL" json:"sec_id"`
	SecName      string    `gorm:"column:sec_name;type:varchar(45);comment:证券名称" json:"sec_name"`
	AlgoId       int       `gorm:"column:algo_id;type:int(10);comment:算法ID" json:"algo_id"`
	AlgoType     int       `gorm:"column:algo_type;type:tinyint(4);comment:算法类型" json:"algo_type"`
	AlgoName     string    `gorm:"column:algo_name;type:varchar(45);comment:算法名称" json:"algo_name"`
	OpenRate     float64   `gorm:"column:open_rate;type:double;comment:开仓率" json:"open_rate"`
	IncomeRate   float64   `gorm:"column:income_rate;type:double;comment:收益率" json:"income_rate"`
	BasisPoint   float64   `gorm:"column:basis_point;type:double;comment:基点" json:"basis_point"`
	CreateTime   time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;type:timestamp;comment:更新时间" json:"update_time"`
}

func (m *TbAlgoOptimizeBase) TableName() string {
	return "tb_algo_optimize_base"
}

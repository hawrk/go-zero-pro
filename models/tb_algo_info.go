// Package model
/*
 Author: hawrkchen
 Date: 2022/3/22 14:35
 Desc:
*/
package models

import (
	"time"
)

// 算法基础信息表
type TbAlgoInfo struct {
	Id           int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:自增ID" json:"id"`
	AlgoId       int       `gorm:"column:algo_id;type:int(11);comment:算法ID;NOT NULL" json:"algo_id"`
	AlgoName     string    `gorm:"column:algo_name;type:varchar(45);comment:算法名称" json:"algo_name"`
	AlgoType     int       `gorm:"column:algo_type;type:smallint(6);comment:算法类型" json:"algo_type"`
	AlgoTypeName string    `gorm:"column:algo_type_name;type:varchar(45);comment:算法类型名称" json:"algo_type_name"`
	Provider     string    `gorm:"column:provider;type:varchar(45);comment:算法厂商" json:"provider"`
	AlgoStatus   int       `gorm:"column:algo_status;type:smallint(6);comment:算法状态  (0-未启用 ， 1-启用)" json:"algo_status"`
	CreateTime   time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;type:timestamp;default:null;comment:更新时间" json:"update_time"`
}

func (m *TbAlgoInfo) TableName() string {
	return "tb_algo_info"
}

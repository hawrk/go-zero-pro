// Package models
/*
 Author: hawrkchen
 Date: 2022/7/14 16:04
 Desc:
*/
package models

import (
	"time"
)

// 一键优选
type TbAlgoOptimize struct {
	Id           int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:自增ID" json:"id"`
	ProviderId   int       `gorm:"column:provider_id;type:int(10);comment:算法ID" json:"provider_id"`
	ProviderName string    `gorm:"column:provider_name;type:varchar(32);comment:厂商名称;NOT NULL" json:"provider_name"`
	SecId        string    `gorm:"column:sec_id;type:varchar(12);comment:证券ID;NOT NULL" json:"sec_id"`
	SecName      string    `gorm:"column:sec_name;type:varchar(45);comment:证券名称" json:"sec_name"`
	AlgoId       int       `gorm:"column:algo_id;type:int(10);comment:算法ID" json:"algo_id"`
	AlgoType     int       `gorm:"column:algo_type;type:tinyint(4);comment:算法类型" json:"algo_type"`
	AlgoName     string    `gorm:"column:algo_name;type:varchar(45);comment:算法名称" json:"algo_name"`
	Score        float64   `gorm:"column:score;type:double;comment:综合分数" json:"score"`
	CreateTime   time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	UpateTime    time.Time `gorm:"column:upate_time;type:timestamp;comment:更新时间" json:"upate_time"`
}

func (m *TbAlgoOptimize) TableName() string {
	return "tb_algo_optimize"
}

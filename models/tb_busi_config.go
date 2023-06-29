// Package models
/*
 Author: hawrkchen
 Date: 2022/7/15 10:41
 Desc:
*/
package models

import (
	"time"
)

// 算法业务配置表
type TbBusiConfig struct {
	Id         int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:自增ID" json:"id"`
	BusiType   int       `gorm:"column:busi_type;type:smallint(6);comment:业务类型 1- 市值占比2-交易量占比3股价类型" json:"busi_type"`
	SecType    int       `gorm:"column:sec_type;type:smallint(6)" json:"sec_type"`
	UpperValue int       `gorm:"column:upper_value;type:int(11);comment:上限" json:"upper_value"`
	LowerValue int       `gorm:"column:lower_value;type:int(11);comment:下限" json:"lower_value"`
	Params     string    `gorm:"column:params;type:varchar(1024);comment:具体参数 ， 必须是json 格式化后的数据" json:"params"`
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:timestamp;comment:更新时间" json:"update_time"`
}

func (m *TbBusiConfig) TableName() string {
	return "tb_busi_config"
}

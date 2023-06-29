// Package models
/*
 Author: hawrkchen
 Date: 2022/7/18 15:54
 Desc:
*/
package models

import (
	"time"
)

// 算法绩效时间图
type TbAlgoTimeLine struct {
	Id           int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:自增ID" json:"id"`
	BatchNo      int64     `gorm:"column:batch_no;type:bigint(20);comment:批次号" json:"batch_no"`
	Date         int       `gorm:"column:date;type:int(11);comment:日期 (20220621);NOT NULL" json:"date"`
	AccountId    string    `gorm:"column:account_id;type:varchar(45);comment:用户ID" json:"account_id"`
	AccountType  int       `gorm:"column:account_type;type:tinyint(4);comment:用户类型  1-普通用户，2-虚拟用户（汇总用户）" json:"account_type"`
	TransactTime int64     `gorm:"column:transact_time;type:bigint(20);comment:交易时间，精确到分(202206120955)" json:"transact_time"`
	AlgoId       int       `gorm:"column:algo_id;type:int(11);comment:算法ID" json:"algo_id"`
	AlgoType     int       `gorm:"column:algo_type;type:int(11);comment:算法类型" json:"algo_type"`
	Provider     string    `gorm:"column:provider;type:varchar(45);comment:算法厂商" json:"provider"`
	AssessScore  float64   `gorm:"column:assess_score;type:decimal(20,2);comment:综合绩效评分" json:"assess_score"`
	RiskScore    float64   `gorm:"column:risk_score;type:decimal(20,2);comment:风险度评分" json:"risk_score"`
	Progress     float64   `gorm:"column:progress;type:decimal(20,2);comment:完成度" json:"progress"`
	SourceFrom   int       `gorm:"column:source_from;type:tinyint(4);comment:数据来源，1:总线推送，2:数据导入" json:"source_from"`
	CreateTime   time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;type:timestamp;default:null;comment:更新时间" json:"update_time"`
}

func (m *TbAlgoTimeLine) TableName() string {
	return "tb_algo_time_line"
}

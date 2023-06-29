// Package models
/*
 Author: hawrkchen
 Date: 2022/4/26 10:54
 Desc:
*/
package models

import (
	"time"
)

// 母单信息表
type TbAlgoOrder struct {
	Id            int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:自增ID" json:"id"`
	Date          int       `gorm:"column:date;type:int(11);comment:时间" json:"date"`
	BatchNo       int64     `gorm:"column:batch_no;type:bigint(20);comment:批次号;NOT NULL" json:"batch_no"`
	BatchName     string    `gorm:"column:batch_name;type:varchar(45);comment:批次号名称" json:"batch_name"`
	BasketId      int       `gorm:"column:basket_id;type:int(11);comment:篮子ID;NOT NULL" json:"basket_id"`
	UserId        string    `gorm:"column:user_id;type:varchar(45);comment:交易账户;NOT NULL" json:"user_id"`
	AlgoId        int       `gorm:"column:algo_id;type:int(11);comment:母单ID;NOT NULL" json:"algo_id"`
	AlgorithmId   int       `gorm:"column:algorithm_id;type:int(11);comment:算法ID" json:"algorithm_id"`
	AlgorithmType int       `gorm:"column:algorithm_type;type:int(10);comment:算法类型" json:"algorithm_type"`
	UsecId        int       `gorm:"column:usec_id;type:int(11);comment:证券代码" json:"usec_id"`
	SecId         string    `gorm:"column:sec_id;type:varchar(8);comment:证券ID;NOT NULL" json:"sec_id"`
	AlgoOrderQty  int64     `gorm:"column:algo_order_qty;type:bigint(20);comment:订单数量" json:"algo_order_qty"`
	UnixTime      int64     `gorm:"column:unix_time;type:bigint(20);comment:交易时间(时间戳)" json:"unix_time"`
	TransTime     int64     `gorm:"column:trans_time;type:bigint(20);comment:交易时间" json:"trans_time"`
	StartTime     int64     `gorm:"column:start_time;type:bigint(20)" json:"start_time"`
	EndTime       int64     `gorm:"column:end_time;type:bigint(20)" json:"end_time"`
	Source        int       `gorm:"column:source;type:tinyint(4);comment:数据来源 0-总线 1-数据修复 2-数据导入" json:"source"`
	Status        int       `gorm:"column:status;type:tinyint(4);comment:状态:0-未处理 1-已处理" json:"status"`
	CreateTime    time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP" json:"create_time"`
}

func (m *TbAlgoOrder) TableName() string {
	return "tb_algo_order"
}

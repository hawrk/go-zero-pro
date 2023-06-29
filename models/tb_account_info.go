// Package models
/*
 Author: hawrkchen
 Date: 2022/6/24 13:49
 Desc:
*/
package models

import (
	"time"
)

// 用户账户信息表
type TbAccountInfo struct {
	Id         int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:自增ID" json:"id"`
	AccountId  int       `gorm:"column:account_id;type:int(11);comment:账户ID;NOT NULL" json:"account_id"`
	UserId     string    `gorm:"column:user_id;type:varchar(45);comment:用户ID" json:"user_id"`
	UserName   string    `gorm:"column:user_name;type:varchar(45);comment:用户名称" json:"user_name"`
	UserPasswd string    `gorm:"column:user_passwd;type:varchar(45);comment:登陆密码" json:"user_passwd"`
	UserType   int       `gorm:"column:user_type;type:smallint(6);comment:用户类型1-普通用户，2-算法厂商用户，3 -管理员" json:"user_type"`
	UserGrade  string    `gorm:"column:user_grade;type:varchar(45);comment:用户级别： ABC" json:"user_grade"`
	GradeFixed int       `gorm:"column:grade_fixed;type:smallint(4);comment:是否固定级别 1-固定，其他-不固定" json:"grade_fixed"`
	RiskGroup  int       `gorm:"column:risk_group;type:int(11);comment:风控组" json:"risk_group"`
	ParUserId  string       `gorm:"column:par_user_id;type:varchar(45);comment:对应管理员ID （可以有多个,用逗号分隔)" json:"par_user_id"`
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:timestamp;default:null;comment:更新时间" json:"update_time"`
}

func (m *TbAccountInfo) TableName() string {
	return "tb_account_info"
}

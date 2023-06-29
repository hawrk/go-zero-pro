// Package models
/*
 Author: hawrkchen
 Date: 2022/10/13 10:15
 Desc:
*/
package models

import (
	"time"
)

// 角色权限表
type TbAuthRole struct {
	Id         int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	RoleId     uint      `gorm:"column:role_id;type:int(8) unsigned;comment:角色ID" json:"role_id"`
	RoleName   string    `gorm:"column:role_name;type:varchar(45);comment:角色名称" json:"role_name"`
	RoleAuth   string    `gorm:"column:role_auth;type:varchar(4096);comment:角色权限列表" json:"role_auth"`
	ChanType   int       `gorm:"column:chan_type;type:tinyint(4);comment:渠道类型 （1-绩效平台 2-算法管理平台）" json:"chan_type"`
	Status     int       `gorm:"column:status;type:tinyint(4);comment:状态" json:"status"`
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:timestamp;default:null;comment:更新时间" json:"update_time"`
}

func (m *TbAuthRole) TableName() string {
	return "tb_auth_role"
}

// Package models
/*
 Author: hawrkchen
 Date: 2022/10/12 15:06
 Desc:
*/
package models

import (
	"time"
)

// 系统权限用户表
type TbAuthUser struct {
	Id         int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:自增ID" json:"id"`
	UserId     string    `gorm:"column:user_id;type:varchar(45);comment:用户ID" json:"user_id"`
	UserName   string    `gorm:"column:user_name;type:varchar(45);comment:用户名称" json:"user_name"`
	UserPasswd string    `gorm:"column:user_passwd;type:varchar(45);comment:用户密码" json:"user_passwd"`
	RoleId     int       `gorm:"column:role_id;type:int(8);comment:角色ID" json:"role_id"`
	RoleName   string    `gorm:"column:role_name;type:varchar(45);comment:角色名称" json:"role_name"`
	ChanType   int       `gorm:"column:chan_type;type:tinyint(4);comment:渠道类型 （1-绩效平台 2-算法管理平台）" json:"chan_type"`
	UserType   int       `gorm:"column:user_type;type:tinyint(4);comment:用户类型： 1-超级管理员 2-普通用户" json:"user_type"`
	Status     int       `gorm:"column:status;type:tinyint(4);comment:状态   1-有效  2-删除" json:"status"`
	FirstLogin int       `gorm:"column:first_login;type:tinyint(4);comment:是否首次登陆 :0-首次登陆,1-非首次登陆  " json:"first_login"`
	ExpireTime int64     `gorm:"column:expire_time;type:bigint(20);comment:密码过期时间（时间戳，180天后过期）" json:"expire_time"`
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:timestamp;default:null;comment:更新时间" json:"update_time"`
}

func (m *TbAuthUser) TableName() string {
	return "tb_auth_user"
}

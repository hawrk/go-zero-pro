// Package models
/*
 Author: hawrkchen
 Date: 2022/10/17 10:17
 Desc:
*/
package models

import (
	"time"
)

// 权限菜单表
type TbAuthMenu struct {
	Id          int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:自增ID" json:"id"`
	MenuId      int       `gorm:"column:menu_id;type:int(8);comment:采单ID" json:"menu_id"`
	MenuName    string    `gorm:"column:menu_name;type:varchar(45);comment:菜单名称" json:"menu_name"`
	MenuType    int       `gorm:"column:menu_type;type:tinyint(4);comment:菜单类型(1-一级菜单2-二级菜单3-三级菜单 11-页面控件）" json:"menu_type"`
	ParManuId   int       `gorm:"column:par_manu_id;type:int(8);comment:父级菜单 (一级菜单父级菜单为0）" json:"par_manu_id"`
	ChanType    int       `gorm:"column:chan_type;type:tinyint(4);comment:渠道类型 （1-绩效平台 2-算法管理平台）" json:"chan_type"`
	CmptType    int       `gorm:"column:cmpt_type;type:tinyint(4);comment:控件类型:1-查询2-新增3-上传4-导出列表5-下载报告6-删除（只有menu_type=11时有值）" json:"cmpt_type"`
	Status      int       `gorm:"column:status;type:tinyint(4);comment:状态  1-正常，2-作废" json:"status"`
	Desc        string    `gorm:"column:desc;type:varchar(45)" json:"desc"`
	AuthDefault int       `gorm:"column:auth_default;type:tinyint(4);comment:是否默认权限    1-是  2-不是" json:"auth_default"`
	CreateTime  time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time;type:timestamp;default:null;comment:更新时间" json:"update_time"`
}

func (m *TbAuthMenu) TableName() string {
	return "tb_auth_menu"
}

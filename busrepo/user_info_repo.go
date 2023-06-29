// Package busrepo
/*
 Author: hawrkchen
 Date: 2022/9/27 10:21
 Desc:
*/
package busrepo

import (
	"algo_assess/busmodels"
	"context"
	"gorm.io/gorm"
)

type UserInfoRepo interface {
	GetGroupIdByUserId(ctx context.Context, userId string) (int, error)
	CheckLogin(ctx context.Context, userId string) ([]*busmodels.TbUserInfo, error)
	GetIdByUserId(ctx context.Context, userId string) (int64, error)
}

type DefaultUserInfo struct {
	DB *gorm.DB
}

func NewUserInfo(conn *gorm.DB) UserInfoRepo {
	return &DefaultUserInfo{
		DB: conn,
	}
}

func (d *DefaultUserInfo) GetGroupIdByUserId(ctx context.Context, userId string) (int, error) {
	var groupId int
	result := d.DB.Model(busmodels.TbUserInfo{}).Select("algo_group").Where("user_id like ?", userId+"%").Find(&groupId)
	if result.Error != nil {
		return 0, result.Error
	}
	return groupId, nil
}

func (d *DefaultUserInfo) CheckLogin(ctx context.Context, userId string) ([]*busmodels.TbUserInfo, error) {
	var account []*busmodels.TbUserInfo
	result := d.DB.Model(busmodels.TbUserInfo{}).Select("user_id, user_name,user_passwd,user_type").Where("user_id like ?", userId+"%").Find(&account)
	if result.Error != nil {
		return nil, result.Error
	}
	return account, nil
}

func (d *DefaultUserInfo) GetIdByUserId(ctx context.Context, userId string) (int64, error) {
	var id int64
	result := d.DB.Model(busmodels.TbUserInfo{}).Select("id").Where("user_id=?", userId).Find(&id)
	if result.Error != nil {
		return 0, result.Error
	}
	return id, nil
}

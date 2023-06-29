// Package repo
/*
 Author: hawrkchen
 Date: 2022/10/17 10:25
 Desc:
*/
package repo

import (
	"account-auth/account-auth-server/models"
	"context"
	"gorm.io/gorm"
)

type AuthMenuRepo interface {
	CreateAuthMenu(ctx context.Context, menuId int) error
	GetDefaultAuthMenu(ctx context.Context, chanType int) ([]models.TbAuthMenu, error)
}

type DefaultAuthMenu struct {
	DB *gorm.DB
}

func NewAuthMenu(conn *gorm.DB) AuthMenuRepo {
	return &DefaultAuthMenu{
		DB: conn,
	}
}

func (d *DefaultAuthMenu) CreateAuthMenu(ctx context.Context, menuId int) error {
	return nil
}

// GetDefaultAuthMenu 取默认菜单
func (d *DefaultAuthMenu) GetDefaultAuthMenu(ctx context.Context, chanType int) ([]models.TbAuthMenu, error) {
	var out []models.TbAuthMenu
	result := d.DB.Model(&models.TbAuthMenu{}).Select("menu_id,menu_name,menu_type,par_manu_id,cmpt_type,auth_default").
		Where("status=1 and chan_type=?", chanType).Order("menu_id").Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

// Package busrepo
/*
 Author: hawrkchen
 Date: 2023/3/23 17:26
 Desc:
*/
package busrepo

import (
	"algo_assess/busmodels"
	"context"
	"gorm.io/gorm"
)

type SecurityInfoRepo interface {
	GetSecurityByIds(ctx context.Context, secId []string) ([]*busmodels.TbSecurityInfo, error)
}

type DefaultSecurityInfo struct {
	DB *gorm.DB
}

func NewSecurityInfo(conn *gorm.DB) SecurityInfoRepo {
	return &DefaultSecurityInfo{
		DB: conn,
	}
}

func (d *DefaultSecurityInfo) GetSecurityByIds(ctx context.Context, secId []string) ([]*busmodels.TbSecurityInfo, error) {
	var out []*busmodels.TbSecurityInfo
	result := d.DB.Model(busmodels.TbSecurityInfo{}).Select("security_id,security_name").
		Where("security_id in ?", secId).Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

// Package repo
/*
 Author: hawrkchen
 Date: 2022/7/15 10:42
 Desc:
*/
package repo

import (
	"algo_assess/models"
	"context"
	"gorm.io/gorm"
	"time"
)

type BusiConfigRepo interface {
	CreateBusiConfig(ctx context.Context, profileType int32, param string) error
	GetAllBusiConfig(ctx context.Context) ([]*models.TbBusiConfig, error)
	GetBusiConfigByProfile(ctx context.Context, profileType int32) (int64, string, error)
	UpdateBusiParam(ctx context.Context, id int64, param string) error
}

type DefaultBusiConfig struct {
	DB *gorm.DB
}

func NewDefaultBusiConfig(conn *gorm.DB) BusiConfigRepo {
	return &DefaultBusiConfig{
		DB: conn,
	}
}

func (d *DefaultBusiConfig) CreateBusiConfig(ctx context.Context, profileType int32, param string) error {
	in := models.TbBusiConfig{
		BusiType:   4,
		SecType:    int(profileType),
		UpperValue: 0,
		LowerValue: 0,
		Params:     param,
	}
	result := d.DB.Create(&in)
	if result.Error != nil || result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

func (d *DefaultBusiConfig) GetAllBusiConfig(ctx context.Context) ([]*models.TbBusiConfig, error) {
	var out []*models.TbBusiConfig
	result := d.DB.Model(models.TbBusiConfig{}).Select("busi_type,sec_type,upper_value,lower_value,params").
		Where("1=1").Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

// GetBusiConfigId 根据profile类型找到该记录的ID
func (d *DefaultBusiConfig) GetBusiConfigByProfile(ctx context.Context, profileType int32) (int64, string, error) {
	var out models.TbBusiConfig
	result := d.DB.Model(models.TbBusiConfig{}).Select("id,params").Where("busi_type=4 and sec_type=?", profileType).Find(&out)
	if result.Error != nil {
		return 0, "", result.Error
	}
	return out.Id, out.Params, nil
}

// UpdateBusiParam 修改param 字段值
func (d *DefaultBusiConfig) UpdateBusiParam(ctx context.Context, id int64, param string) error {
	result := d.DB.Model(models.TbBusiConfig{}).Where("id=?", id).
		Updates(models.TbBusiConfig{
			Params:     param,
			UpdateTime: time.Now(),
		})
	return result.Error
}

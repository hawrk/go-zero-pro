// Package busrepo
/*
 Author: hawrkchen
 Date: 2022/9/27 10:49
 Desc:
*/
package busrepo

import (
	"algo_assess/busmodels"
	"context"
	"gorm.io/gorm"
)

type AlgoGroupRepo interface {
	GetAlgoPropertyById(ctx context.Context, id int) (string, error)
}

type DefaultAlgoGroup struct {
	DB *gorm.DB
}

func NewAlgoGroup(conn *gorm.DB) AlgoGroupRepo {
	return &DefaultAlgoGroup{
		DB: conn,
	}
}

func (d *DefaultAlgoGroup) GetAlgoPropertyById(ctx context.Context, id int) (string, error) {
	var property string
	result := d.DB.Model(busmodels.TbAlgoGroupInfo{}).Select("algo_property").Where("id=?", id).Find(&property)
	if result.Error != nil {
		return "", result.Error
	}
	return property, nil
}

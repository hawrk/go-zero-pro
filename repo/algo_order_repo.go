// Package repo
/*
 Author: hawrkchen
 Date: 2022/4/26 10:55
 Desc:
*/
package repo

import (
	"algo_assess/global"
	"algo_assess/models"
	"context"
	"gorm.io/gorm"
)

type AlgoOrderRepo interface {
	CreateAlgoOrder(ctx context.Context, data *global.MAlgoOrder) error
}

type DefaultAlgoOrder struct {
	DB *gorm.DB
}

func NewDefaultAlgoOrder(conn *gorm.DB) *DefaultAlgoOrder {
	return &DefaultAlgoOrder{
		DB: conn,
	}
}

func (a *DefaultAlgoOrder) CreateAlgoOrder(ctx context.Context, data *global.MAlgoOrder) error {
	algoOrder := &models.TbAlgoOrder{
		AlgoId:       data.AlgoId,
		AlgorithmId:  data.AlgorithmId,
		UsecId:       data.UsecId,
		SecId:        data.SecId,
		AlgoOrderQty: data.AlgoOrderQty,
		TransTime:    data.TransTime,
	}
	result := a.DB.Create(&algoOrder)
	if result.Error != nil || result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

// Package repo
/*
 Author: hawrkchen
 Date: 2022/3/28 16:14
 Desc:
*/
package repo

import (
	"algo_assess/assess-mq-server/proto/order"
	"algo_assess/global"
	"algo_assess/models"
	"context"
	"gorm.io/gorm"
)

type OrderAssessRepo interface {
	CreateOrderAssess(ctx context.Context, data *global.OrderAssess) error
	UpdateAlgoAssess(ctx context.Context, orders []*models.TbAlgoAssess, data *order.ChildOrderPerf) error
	GetAlgoAssess(ctx context.Context, algoId int32, usecId string, td int32, statusType int32, start, end int64) ([]*models.TbAlgoAssess, *gorm.DB)
}

type DefaultOrderAssess struct {
	DB *gorm.DB
}

func NewAlgoAssessRepo(conn *gorm.DB) OrderAssessRepo {
	return &DefaultOrderAssess{
		DB: conn,
	}
}

func (o *DefaultOrderAssess) GetAlgoAssess(ctx context.Context, algoId int32, usecId string, td int32, statusType int32,
	start, end int64) ([]*models.TbAlgoAssess, *gorm.DB) {
	var assess []*models.TbAlgoAssess
	result := o.DB.Select("transact_time, last_price, arrived_price, vwap, deal_rate, order_qty, last_qty,cancel_qty, rejected_qty,market_rate,"+
		" deal_progress, vwap_deviation, arrived_price_deviation").
		Where("algorithm_id = ? and security_id=? and time_dimension = ? and transact_time between ? and ? ",
			algoId, usecId, td, start, end).
		Find(&assess)
	if result.Error != nil {
		return nil, result
	}
	return assess, result
}

func (o *DefaultOrderAssess) CreateOrderAssess(ctx context.Context, data *global.OrderAssess) error {
	assess := &models.TbAlgoAssess{
		AlgorithmType:         data.AlgorithmType,
		AlgorithmId:           data.AlgorithmId,
		UsecurityId:           data.UsecurityId,
		SecurityId:            data.SecurityId,
		TimeDimension:         data.TimeDimension,
		TransactTime:          data.TransactAt,
		ArrivedPrice:          data.ArrivedPrice,
		LastPrice:             data.LastPrice,
		Vwap:                  data.Vwap,
		DealRate:              data.DealRate,
		OrderQty:              data.OrderQty,
		LastQty:               data.LastQty,
		CancelQty:             data.CancelQty,
		RejectedQty:           data.RejectedQty,
		MarketRate:            data.MarketRate,
		DealProgress:          data.DealProgress,
		VwapDeviation:         data.VwapDeviation,
		ArrivedPriceDeviation: data.ArrivedPriceDeviation,
		CreateTime:            data.CreateTime,
	}
	result := o.DB.Create(&assess)
	if result.Error != nil || result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

func (o *DefaultOrderAssess) UpdateAlgoAssess(ctx context.Context, orders []*models.TbAlgoAssess, data *order.ChildOrderPerf) error {
	return nil
}

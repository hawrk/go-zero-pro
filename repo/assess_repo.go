// Package repo
/*
 Author: hawrkchen
 Date: 2022/3/28 16:14
 Desc:
*/
package repo

import (
	"algo_assess/assess-rpc-server/proto"
	"algo_assess/global"
	"algo_assess/models"
	"context"
	"gorm.io/gorm"
	"time"
)

type OrderAssessRepo interface {
	CreateOrderAssess(ctx context.Context, data *global.OrderAssess) error
	UpdateAlgoAssess(ctx context.Context, data *global.OrderAssess) error
	GetAlgoAssess(ctx context.Context, req *proto.GeneralReq) ([]*models.TbAlgoAssess, *gorm.DB)
	GetAlgoAssessByKey(ctx context.Context, userId string, transTime int64, algoId int, secId string) (int64, error)
}

type DefaultOrderAssess struct {
	DB *gorm.DB
}

func NewAlgoAssessRepo(conn *gorm.DB) OrderAssessRepo {
	return &DefaultOrderAssess{
		DB: conn,
	}
}

func (o *DefaultOrderAssess) GetAlgoAssess(ctx context.Context, req *proto.GeneralReq) ([]*models.TbAlgoAssess, *gorm.DB) {
	var assess []*models.TbAlgoAssess
	result := o.DB.Select("transact_time, last_price, arrived_price, vwap, deal_rate, order_qty, last_qty,cancel_qty, rejected_qty,market_rate,"+
		" deal_progress, vwap_deviation, arrived_price_deviation").
		Where("user_id = ? and algorithm_id = ? and security_id=? and time_dimension = ? and transact_time between ? and ? ",
			req.GetUserId(), req.GetAlgoId(), req.GetSecId(), req.GetTimeDemension(), req.GetStartTime(), req.GetEndTime()).
		Find(&assess)
	if result.Error != nil {
		return nil, result
	}
	return assess, result
}

func (o *DefaultOrderAssess) CreateOrderAssess(ctx context.Context, data *global.OrderAssess) error {
	assess := &models.TbAlgoAssess{
		BatchNo:       data.BatchNo,
		AlgorithmType: uint(data.AlgorithmType),
		AlgorithmId:   uint(data.AlgorithmId),
		UsecurityId:   data.UsecurityId,
		SecurityId:    data.SecurityId,
		UserId:        data.UserId, // add 账户信息
		TimeDimension: data.TimeDimension,
		TransactTime:  data.TransactAt,
		ArrivedPrice:  data.ArrivedPrice,
		LastPrice:     data.LastPrice,
		//Vwap:                  data.Vwap,
		DealRate:              data.DealRate,
		OrderQty:              data.OrderQty,
		LastQty:               data.LastQty,
		CancelQty:             data.CancelQty,
		RejectedQty:           data.RejectedQty,
		MarketRate:            data.MarketRate,
		DealProgress:          data.DealProgress,
		VwapDeviation:         data.VwapDeviation,
		ArrivedPriceDeviation: data.ArrivedPriceDeviation,
		SourceFrom:            data.SourceFrom,
		CreateTime:            data.CreateTime,
	}
	result := o.DB.Create(&assess)
	if result.Error != nil || result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

func (o *DefaultOrderAssess) UpdateAlgoAssess(ctx context.Context, data *global.OrderAssess) error {
	result := o.DB.Model(models.TbAlgoAssess{}).Where("user_id=? and transact_time=? and algorithm_id=? and security_id=?",
		data.UserId, data.TransactAt, data.AlgorithmId, data.SecurityId).
		Updates(models.TbAlgoAssess{
			//AlgorithmType: uint(data.AlgorithmType),
			//AlgorithmId:   uint(data.AlgorithmId),
			//UsecurityId:   data.UsecurityId,
			//SecurityId:    data.SecurityId,
			//UserId:        data.UserId, // add 账户信息
			TimeDimension: data.TimeDimension,
			//TransactTime:  data.TransactAt,
			ArrivedPrice: data.ArrivedPrice,
			LastPrice:    data.LastPrice,
			//Vwap:                  data.Vwap,
			DealRate:              data.DealRate,
			OrderQty:              data.OrderQty,
			LastQty:               data.LastQty,
			CancelQty:             data.CancelQty,
			RejectedQty:           data.RejectedQty,
			MarketRate:            data.MarketRate,
			DealProgress:          data.DealProgress,
			VwapDeviation:         data.VwapDeviation,
			ArrivedPriceDeviation: data.ArrivedPriceDeviation,
			//CreateTime:            data.CreateTime,
			UpdateTime: time.Now(),
		})
	return result.Error
}

// GetAlgoAssessByKey 根据key查找是否已有记录   key: transacttime+algoId+secId
func (o *DefaultOrderAssess) GetAlgoAssessByKey(ctx context.Context, userId string, transTime int64, algoId int, secId string) (int64, error) {
	var id int64
	result := o.DB.Model(models.TbAlgoAssess{}).Select("id").Where("user_id=? and transact_time=? and algorithm_id=? and security_id=?",
		userId, transTime, algoId, secId).Find(&id)
	if result.Error != nil {
		return 0, result.Error
	}
	return id, nil
}

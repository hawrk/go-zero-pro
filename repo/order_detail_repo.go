// Package repo
/*
 Author: hawrkchen
 Date: 2022/3/28 17:23
 Desc:
*/
package repo

import (
	"algo_assess/global"
	"algo_assess/models"
	"context"
	"gorm.io/gorm"
	"time"
)

type OrderDetailRepo interface {
	CreateOrderDetail(ctx context.Context, data *global.ChildOrderData) error
	UpdateOrderDetail(transactAt int64) error
	QueryOrderDetail(t int64) (orders []*models.TbAlgoOrderDetail, err error)
}

type DefaultOrderDetail struct {
	DB *gorm.DB
}

func NewOrderDetailRepo(conn *gorm.DB) OrderDetailRepo {
	return &DefaultOrderDetail{
		DB: conn,
	}
}

func (d *DefaultOrderDetail) CreateOrderDetail(ctx context.Context, data *global.ChildOrderData) error {
	detail := &models.TbAlgoOrderDetail{
		ChildOrderId:  data.OrderId,
		AlgoOrderId:   uint(data.AlgoOrderId),
		AlgorithmType: data.AlgorithmType,
		AlgorithmId:   data.AlgoId,
		UsecurityId:   data.UsecId,
		SecurityId:    data.SecId,
		OrderQty:      data.OrderQty,
		Price:         data.Price,
		OrderType:     data.OrderType,
		LastPx:        data.LastPx,
		LastQty:       data.LastQty,
		ComQty:        data.ComQty,
		ArrivedPrice:  data.ArrivePrice,
		OrdStatus:     data.ChildOrderStatus,
		//TransactTime:  data.TransTime,
		TransactAt: data.TransTime,
		ProcStatus: 0,
		CreateTime: time.Now(),
	}
	result := d.DB.Create(detail)
	if result.Error != nil || result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

func (d *DefaultOrderDetail) UpdateOrderDetail(transactAt int64) error {
	result := d.DB.Model(models.TbAlgoOrderDetail{}).Where("transact_at = ?", transactAt).
		Updates(models.TbAlgoOrderDetail{ProcStatus: 1})
	return result.Error
}

func (d *DefaultOrderDetail) QueryOrderDetail(t int64) (orders []*models.TbAlgoOrderDetail, err error) {
	result := d.DB.Where("transact_at = ?", t).Find(&orders)
	if result.Error != nil {
		return nil, err
	}
	return orders, nil
}

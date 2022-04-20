// Package repo
/*
 Author: hawrkchen
 Date: 2022/3/28 17:23
 Desc:
*/
package repo

import (
	"algo_assess/assess-mq-server/proto/order"
	"algo_assess/models"
	"context"
	"gorm.io/gorm"
	"time"
)

type OrderDetailRepo interface {
	CreateOrderDetail(ctx context.Context, t int64, data *order.ChildOrderPerf) error
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

func (d *DefaultOrderDetail) CreateOrderDetail(ctx context.Context, t int64, data *order.ChildOrderPerf) error {
	detail := &models.TbAlgoOrderDetail{
		ChildOrderId:  int64(data.GetId()),
		AlgoOrderId:   uint(data.GetAlgoOrderId()),
		AlgorithmType: uint(data.GetAlgorithmType()),
		AlgorithmId:   uint(data.GetAlgorithmId()),
		UsecurityId:   uint(data.GetUSecurityId()),
		SecurityId:    data.GetSecurityId(),
		OrderQty:      int64(data.GetOrderQty()),
		Price:         int64(data.GetPrice()),
		OrderType:     uint(data.GetOrderType()),
		LastPx:        int64(data.GetLastPx()),
		LastQty:       int64(data.GetLastQty()),
		ComQty:        int64(data.GetCumQty()),
		ArrivedPrice:  int64(data.GetArrivedPrice()),
		OrdStatus:     uint(data.GetChildOrdStatus()),
		TransactTime:  int64(data.GetTransactTime()),
		TransactAt:    t,
		ProcStatus:    0,
		CreateTime:    time.Now(),
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

// Package repo
/*
 Author: hawrkchen
 Date: 2022/3/28 17:23
 Desc:
*/
package repo

import (
	"algo_assess/models"
	"algo_assess/mqueue/proto/order"
	"context"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"time"
)

type OrderDetailRepo interface {
	CreateOrderDetail(ctx context.Context, t uint64, data *order.ChildOrderPerf) error
	UpdateOrderDetail(transactAt uint64) error
}

type DefaultOrderDetail struct {
	DB *gorm.DB
}

func NewOrderDetailRepo(conn *gorm.DB) OrderDetailRepo {
	return &DefaultOrderDetail{
		DB: conn,
	}
}

func (d *DefaultOrderDetail) CreateOrderDetail(ctx context.Context, t uint64, data *order.ChildOrderPerf) error {
	detail := &models.TbAlgoOrderDetail{
		ChildOrderId:  cast.ToInt64(data.GetId()),
		AlgorithmType: cast.ToUint(data.GetAlgorithmType()),
		AlgorithmId:   cast.ToUint(data.GetAlgorithmId()),
		UsecurityId:   cast.ToUint(data.USecurityId),
		SecurityId:    data.SecurityId,
		OrderQty:      cast.ToUint(data.GetOrderQty()),
		Price:         cast.ToUint(data.GetPrice()),
		OrderType:     cast.ToUint(data.GetOrderType()),
		LastPx:        cast.ToUint(data.GetLastPx()),
		LastQty:       cast.ToUint(data.GetLastQty()),
		ComQty:        cast.ToUint(data.GetCumQty()),
		OrdStatus:     cast.ToUint(data.GetChildOrdStatus()),
		TransactTime:  cast.ToUint(data.GetTransactTime()),
		TransactAt:    t,
		ProcStatus: 0,
		CreateTime:    time.Now(),
	}
	result := d.DB.Create(detail)
	if result.Error != nil || result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

func (d *DefaultOrderDetail) UpdateOrderDetail(transactAt uint64) error {
	//TODO
	return nil
}
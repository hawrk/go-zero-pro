// Package repo
/*
 Author: hawrkchen
 Date: 2022/3/28 16:14
 Desc:
*/
package repo

import (
	"algo_assess/global"
	"algo_assess/models"
	"algo_assess/mqueue/proto/order"
	"context"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type OrderAssessRepo interface {
	CreateOrderAssess(ctx context.Context, data *global.OrderAssess) error
	UpdateOrder(ctx context.Context, orders []*models.TbAlgoAssess, data *order.ChildOrderPerf) error
	GetOrders(ctx context.Context, id uint32, start, end string) ([]*models.TbAlgoAssess, *gorm.DB)
}

type DefaultOrderAssess struct {
	DB *gorm.DB
}

func NewOrderRepo(conn *gorm.DB) OrderAssessRepo {
	return &DefaultOrderAssess{
		DB: conn,
	}
}

func (o *DefaultOrderAssess) GetOrders(ctx context.Context, id uint32, start, end string) ([]*models.TbAlgoAssess, *gorm.DB) {
	var orders []*models.TbAlgoAssess
	result := o.DB.Select("child_order_id, order_status_type").
		Where("child_order_id = ? and trade_time_at between ? and ? ", id, start, end).
		Find(&orders)
	if result.Error != nil {
		return nil, result
	}
	return orders, result
}

func (o *DefaultOrderAssess) CreateOrderAssess(ctx context.Context, data *global.OrderAssess) error {
	var assess models.TbAlgoAssess
	_ = copier.Copy(&assess, data)
	result := o.DB.Create(assess)
	if result.Error != nil || result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

func (o *DefaultOrderAssess) UpdateOrder(ctx context.Context, orders []*models.TbAlgoAssess, data *order.ChildOrderPerf) error {
	return nil
}

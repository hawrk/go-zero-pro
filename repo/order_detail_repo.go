// Package repo
/*
 Author: hawrkchen
 Date: 2022/3/28 17:23
 Desc:
*/
package repo

import (
	"algo_assess/assess-rpc-server/proto"
	"algo_assess/global"
	"algo_assess/models"
	"context"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"time"
)

type OrderDetailRepo interface {
	CreateOrderDetail(ctx context.Context, data *global.ChildOrderData) error
	UpdateOrderDetail(ctx context.Context, data *global.ChildOrderData) error
	UpdateOrderDetailStatus(transactAt int64) error
	QueryOrderDetail(t int64) (orders []*models.TbAlgoOrderDetail, err error)
	QueryChildOrder(ctx context.Context, in *proto.ReqQueryChildOrder, algoId int32) (int64, []*models.TbAlgoOrderDetail, error)
	GetChildByAlgorithms(ctx context.Context, algorithmIds []int32, date string) ([]*models.TbAlgoOrderDetail, error)
	GetChildOrder(orderType int32, in []int32) ([]*models.TbAlgoOrderDetail, error)
	QueryChildOrderByChildOrder(date int32, childOrderId int32) (models.TbAlgoOrderDetail, error)
	QueryChildOrderByDate(date int, algoId int) ([]*models.TbAlgoOrderDetail, error)
	QueryAllChildOrderIds(ctx context.Context) ([]int64, error)
	BatchQueryChildOrder(ctx context.Context, ids []int64) ([]int64, error)
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
		Date:          int(data.CurDate),
		BatchNo:       data.BatchNo,
		BatchName:     data.BatchName,
		ChildOrderId:  data.OrderId,
		AlgoOrderId:   uint(data.AlgoOrderId),
		AlgorithmType: uint(data.AlgorithmType),
		AlgorithmId:   uint(data.AlgoId),
		UserId:        data.UserId,
		UsecurityId:   data.UsecId,
		SecurityId:    data.SecId,
		TradeSide:     data.TradeSide,
		OrderQty:      data.OrderQty,
		Price:         data.Price,
		OrderType:     data.OrderType,
		LastPx:        data.LastPx,
		LastQty:       data.LastQty,
		ComQty:        data.ComQty,
		ArrivedPrice:  data.ArrivePrice,
		OrdStatus:     data.ChildOrderStatus,
		TotalFee:      float64(data.TotalFee),
		TransactTime:  data.UnixTimeMillSec,
		TransactAt:    data.TransTime,
		ProcStatus:    0,
		Source:        data.SourceFrom,
		//CreateTime: time.Now(),
	}
	result := d.DB.Create(detail)
	if result.Error != nil || result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

func (d *DefaultOrderDetail) UpdateOrderDetail(ctx context.Context, data *global.ChildOrderData) error {
	result := d.DB.Model(models.TbAlgoOrderDetail{}).Where("child_order_id=?", data.OrderId).Updates(
		map[string]interface{}{
			"trade_side":    data.TradeSide,
			"order_qty":     data.OrderQty,
			"price":         data.Price,
			"last_px":       data.LastPx,
			"last_qty":      data.LastQty,
			"com_qty":       data.ComQty,
			"arrived_price": data.ArrivePrice,
			"total_fee":     data.TotalFee,
			"ord_status":    data.ChildOrderStatus,
			"transact_time": data.UnixTimeMillSec,
			"transact_at":   data.TransTime,
			"source":        data.SourceFrom})
	//result := d.DB.Model(models.TbAlgoOrderDetail{}).Select("trade_side,order_qty,price,last_px,last_qty,com_qty,arrived_price,total_fee").
	//	Where("child_order_id=?", data.OrderId).Updates(models.TbAlgoOrderDetail{
	//	TradeSide:    data.TradeSide,
	//	OrderQty:     data.OrderQty,
	//	Price:        data.Price,
	//	LastPx:       data.LastPx,
	//	LastQty:      data.LastQty,
	//	ComQty:       data.ComQty,
	//	ArrivedPrice: data.ArrivePrice,
	//	TotalFee:     float64(data.TotalFee),
	//	OrdStatus:    data.ChildOrderStatus,
	//	TransactTime: data.UnixTime,
	//	TransactAt:   data.TransTime,
	//	Source:       data.SourceFrom,
	//})
	return result.Error
}

func (d *DefaultOrderDetail) UpdateOrderDetailStatus(transactAt int64) error {
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

func (d *DefaultOrderDetail) QueryChildOrder(ctx context.Context, in *proto.ReqQueryChildOrder, algoId int32) (int64, []*models.TbAlgoOrderDetail, error) {
	var infos []*models.TbAlgoOrderDetail
	var count int64
	model := d.DB.Model(&models.TbAlgoOrderDetail{})
	//if in.UserType != 1 {
	if in.GetUserId() != "" {
		model.Where("user_id", in.GetUserId())
	}
	//}
	if in.GetSecurityId() != "" {
		model.Where("security_id", in.GetSecurityId())
	}
	if in.GetChildOrderId() != 0 {
		model.Where("child_order_id", in.GetChildOrderId())
	}
	if in.GetAlgoOrderId() != 0 {
		model.Where("algo_order_id", in.GetAlgoOrderId())
	}
	if algoId != 0 {
		model.Where("algorithm_id", algoId) // 算法ID
	}
	if in.GetScene() == 1 {
		model.Where("source in(0,1)")
	}
	if in.GetStartTime() != 0 && in.GetEndTime() != 0 {
		model.Where("transact_at between ? and ?", cast.ToString(in.GetStartTime())+"0000", cast.ToString(in.GetEndTime())+"2359")
	}
	err := model.Order("create_time desc").Count(&count).Limit(int(in.GetPageNum())).Offset(int((in.GetPageId() - 1) * in.GetPageNum())).Find(&infos).Error
	if err != nil {
		return 0, nil, err
	}
	return count, infos, nil
}

func (d *DefaultOrderDetail) GetChildByAlgorithms(ctx context.Context, algorithmIds []int32, date string) ([]*models.TbAlgoOrderDetail, error) {
	var out []*models.TbAlgoOrderDetail
	result := d.DB.Model(&models.TbAlgoOrderDetail{}).Where("transact_at like ?  and algorithm_id in ?",
		date+"%", algorithmIds).Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

func (d *DefaultOrderDetail) GetChildOrder(orderType int32, in []int32) ([]*models.TbAlgoOrderDetail, error) {
	var infos []*models.TbAlgoOrderDetail
	model := d.DB.Model(&models.TbAlgoOrderDetail{})
	if orderType == 1 {
		model.Where("algo_order_id in ?", in)
	} else if orderType == 2 {
		model.Where("child_order_id in ?", in)
	}
	err := model.Find(&infos).Error
	if err != nil {
		return nil, err
	}
	return infos, nil
}

// QueryChildOrderByChildOrder 根据日期和子单号确定唯一性
func (d *DefaultOrderDetail) QueryChildOrderByChildOrder(date int32, childOrderId int32) (models.TbAlgoOrderDetail, error) {
	var out models.TbAlgoOrderDetail
	result := d.DB.Model(&models.TbAlgoOrderDetail{}).Select("id, date,algo_order_id,child_order_id,algorithm_type,algorithm_id,"+
		"user_id,usecurity_id,security_id,trade_side,order_qty,price,order_type,last_px,last_qty,com_qty,arrived_price,"+
		"total_fee,ord_status,transact_time,transact_at").
		Where("date=? and child_order_id=?", date, childOrderId).Find(&out)
	if result.Error != nil {
		return models.TbAlgoOrderDetail{}, nil
	}
	return out, nil
}

// QueryChildOrderByDate 根据日期和母单号查询子单
func (d *DefaultOrderDetail) QueryChildOrderByDate(date int, algoOrderId int) ([]*models.TbAlgoOrderDetail, error) {
	var out []*models.TbAlgoOrderDetail
	result := d.DB.Model(&models.TbAlgoOrderDetail{}).Select("algo_order_id,child_order_id,algorithm_type,algorithm_id,"+
		"user_id,usecurity_id,security_id,trade_side,order_qty,price,order_type,last_px,last_qty,com_qty,arrived_price,"+
		"total_fee,ord_status,transact_time,transact_at").
		Where("date=? and algo_order_id=?", date, algoOrderId).Find(&out)
	if result.Error != nil {
		return nil, nil
	}
	return out, nil
}

// QueryAllChildOrderIds 废弃
func (d *DefaultOrderDetail) QueryAllChildOrderIds(ctx context.Context) ([]int64, error) {
	var ids []int64
	t := time.Now().AddDate(0, 0, -7).Format(global.TimeFormat)
	result := d.DB.Model(&models.TbAlgoOrderDetail{}).Select("child_order_id").Where("create_time > ?", t).Find(&ids)
	if result.Error != nil {
		return nil, result.Error
	}
	return ids, nil
}

// BatchQueryChildOrder 数据修复用
func (d *DefaultOrderDetail) BatchQueryChildOrder(ctx context.Context, ids []int64) ([]int64, error) {
	var out []int64
	result := d.DB.Model(&models.TbAlgoOrderDetail{}).Select("child_order_id").Where("child_order_id in ?", ids).Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

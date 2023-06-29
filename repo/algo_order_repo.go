// Package repo
/*
 Author: hawrkchen
 Date: 2022/4/26 10:55
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

type AlgoOrderRepo interface {
	CreateAlgoOrder(ctx context.Context, data *global.MAlgoOrder) error
	UpdateAlgoOrder(ctx context.Context, data *global.MAlgoOrder) error
	QueryAlgoOrder(ctx context.Context, in *proto.ReqQueryAlgoOrder, algorithmId int32) (int64, []*models.TbAlgoOrder, error)
	GetAlgoByAlgorithms(ctx context.Context, algorithmIds []int32, date string) ([]*models.TbAlgoOrder, error)
	GetAlgoOrder(in []int32) ([]*models.TbAlgoOrder, error)
	QueryAlgoOrderByDate(date int32, id int32) (models.TbAlgoOrder, error)
	QueryAlgoOrderIds(ctx context.Context) ([]int64, error)
	BatchQueryAlgoOrder(ctx context.Context, ids []int) ([]int64, error)
}

type DefaultAlgoOrder struct {
	DB *gorm.DB
}

func NewDefaultAlgoOrder(conn *gorm.DB) AlgoOrderRepo {
	return &DefaultAlgoOrder{
		DB: conn,
	}
}

func (a *DefaultAlgoOrder) CreateAlgoOrder(ctx context.Context, data *global.MAlgoOrder) error {
	var Qty int64
	if data.AlgorithmType == global.AlgoTypeT0 { // T0前面已经*2，落表要还原
		Qty = data.AlgoOrderQty / 2
	} else {
		Qty = data.AlgoOrderQty
	}
	algoOrder := &models.TbAlgoOrder{
		Date:          cast.ToInt(cast.ToString(data.TransTime)[:8]),
		BatchNo:       data.BatchNo,
		BatchName:     data.BatchName,
		UserId:        data.UserId,
		AlgoId:        data.AlgoId,
		BasketId:      data.BasketId,
		AlgorithmId:   data.AlgorithmId,
		AlgorithmType: data.AlgorithmType,
		UsecId:        data.UsecId,
		SecId:         data.SecId,
		AlgoOrderQty:  Qty,
		UnixTime:      data.UnixTimeMillSec,
		TransTime:     data.TransTime,
		StartTime:     data.StartTime,
		EndTime:       data.EndTime,
		Source:        data.SourceFrom,
	}
	result := a.DB.Create(&algoOrder)
	if result.Error != nil || result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

func (a *DefaultAlgoOrder) UpdateAlgoOrder(ctx context.Context, data *global.MAlgoOrder) error {
	result := a.DB.Model(models.TbAlgoOrder{}).Where("algo_id=?", data.AlgoId).Updates(
		map[string]interface{}{
			"algo_order_qty": data.AlgoOrderQty,
			"unix_time":      data.UnixTimeMillSec,
			"start_time":     data.StartTime,
			"end_time":       data.EndTime,
			"source":         data.SourceFrom,
		})
	//result := a.DB.Model(models.TbAlgoOrder{}).Where("algo_id=?", data.AlgoId).Updates(models.TbAlgoOrder{
	//	AlgoOrderQty: data.AlgoOrderQty,
	//	UnixTime:     data.UnixTimeMillSec,
	//	StartTime:    data.StartTime,
	//	EndTime:      data.EndTime,
	//})
	return result.Error
}

func (a *DefaultAlgoOrder) QueryAlgoOrder(ctx context.Context, in *proto.ReqQueryAlgoOrder, algorithmId int32) (int64, []*models.TbAlgoOrder, error) {
	var infos []*models.TbAlgoOrder
	var count int64
	model := a.DB.Model(&models.TbAlgoOrder{})
	if in.GetUserId() != "" {
		model.Where("user_id", in.GetUserId())
	}
	//if in.UserType != 1 {
	//	if in.GetUserId() != "" {
	//		model.Where("user_id", in.GetUserId())
	//	}
	//}
	if in.GetAlgoId() != 0 {
		model.Where("algo_id", in.GetAlgoId()) // 母单号
	}
	if in.GetSecId() != "" {
		model.Where("sec_id", in.GetSecId()) // 证券代码
	}
	if algorithmId != 0 {
		model.Where("algorithm_id", algorithmId) // 算法ID
	}

	if in.GetScene() == 1 { // scene = 2绩效分析的话，则全取
		model.Where("source in(0,1)")
	}
	if in.GetStartTime() != 0 && in.GetEndTime() != 0 {
		model.Where("trans_time between  ? and ? ", cast.ToString(in.GetStartTime())+"0000", cast.ToString(in.GetEndTime())+"2359")
	}

	err := model.Order("create_time desc").Count(&count).Limit(int(in.GetPageNum())).Offset(int((in.GetPageId() - 1) * in.GetPageNum())).Find(&infos).Error
	if err != nil {
		return 0, nil, err
	}
	return count, infos, nil
}

func (a *DefaultAlgoOrder) GetAlgoOrder(in []int32) ([]*models.TbAlgoOrder, error) {
	var infos []*models.TbAlgoOrder
	model := a.DB.Model(&models.TbAlgoOrder{})
	model.Where("algo_id in ?", in)
	err := model.Find(&infos).Error
	if err != nil {
		return nil, err
	}
	return infos, nil
}

// QueryAlgoOrderByDate 根据日期和母单ID确认唯一性
func (a *DefaultAlgoOrder) QueryAlgoOrderByDate(date int32, id int32) (models.TbAlgoOrder, error) {
	var out models.TbAlgoOrder
	result := a.DB.Model(&models.TbAlgoOrder{}).Select("id,date,basket_id,user_id,algo_id,algorithm_id,"+
		"algorithm_type,usec_id,sec_id,algo_order_qty,unix_time,trans_time,start_time,end_time").
		Where("date= ? and algo_id=?", date, id).Find(&out)
	if result.Error != nil {
		return models.TbAlgoOrder{}, result.Error
	}
	return out, nil
}

// GetAlgoByAlgorithms 母单修复数据，根据算法ID列表查当天所有母单数据
func (a *DefaultAlgoOrder) GetAlgoByAlgorithms(ctx context.Context, algorithmIds []int32, date string) ([]*models.TbAlgoOrder, error) {
	var out []*models.TbAlgoOrder
	result := a.DB.Model(&models.TbAlgoOrder{}).Select("batch_no,batch_name,basket_id,user_id,algo_id,algorithm_id,algorithm_type,"+
		"usec_id,sec_id,algo_order_qty,unix_time,start_time,end_time,source").Where("trans_time like ?  and algorithm_id in ?",
		date+"%", algorithmIds).Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

// QueryAlgoOrderIds 查该母单表下所有的ID，用作数据修复校验用---废弃
func (a *DefaultAlgoOrder) QueryAlgoOrderIds(ctx context.Context) ([]int64, error) {
	var ids []int64
	// 只取七天内的数据
	t := time.Now().AddDate(0, 0, -7).Format(global.TimeFormat)
	result := a.DB.Model(&models.TbAlgoOrder{}).Select("algo_id").Where("create_time > ?", t).Find(&ids)
	if result.Error != nil {
		return nil, result.Error
	}
	return ids, nil
}

func (a *DefaultAlgoOrder) BatchQueryAlgoOrder(ctx context.Context, ids []int) ([]int64, error) {
	var out []int64
	result := a.DB.Model(&models.TbAlgoOrder{}).Select("algo_id").Where("algo_id in ?", ids).Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil

}

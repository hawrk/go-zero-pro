// Package repo
/*
 Author: hawrkchen
 Date: 2022/7/18 15:56
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

type AlgoTimeLineOrigRepo interface {
	CreateAlgoTimeLine(ctx context.Context, ps *global.ProfileSum) error
	UpdateAlgoTimeLine(ctx context.Context, ps *global.ProfileSum) error
	GetAlgoTimeLineByAllUser(ctx context.Context, date int64, algoId int, userId string, userType int32) ([]*models.TbAlgoTimeLineOrig, error)
	GetImportAlgoTimeLine(ctx context.Context, date int64, batchNo int64, userId string, userType int32) ([]*models.TbAlgoTimeLine, error)
	//GetAlgoTimeLineCrossDay(ctx context.Context, start, end int64, UserId string, algoIds int32) ([]*models.TbAlgoTimeLineOrig, error)
	GetMultiTimeLine(ctx context.Context, start, end int64, userId string, userType int32, algoId int) ([]*models.TbAlgoTimeLineOrig, error)
	GetMultiTimeLineBatch(ctx context.Context, batchId int64, userId string, userType int32, algoId int32) ([]*models.TbAlgoTimeLine, error)
	// reload
	ReloadTimeLine(ctx context.Context, date int64) ([]*models.TbAlgoTimeLineOrig, error)

	// exception process
	GetDataByTimeLineKey(ctx context.Context, date int64, userId string, algoId int) (models.TbAlgoTimeLineOrig, error)
}

type DefaultAlgoTimeLineOrig struct {
	DB *gorm.DB
}

func NewAlgoTimeLineOrig(conn *gorm.DB) AlgoTimeLineOrigRepo {
	return &DefaultAlgoTimeLineOrig{
		DB: conn,
	}
}

func (d *DefaultAlgoTimeLineOrig) CreateAlgoTimeLine(ctx context.Context, ps *global.ProfileSum) error {
	v := &models.TbAlgoTimeLineOrig{
		Date:         int(ps.Date),
		BatchNo:      ps.BatchNo,
		AccountId:    ps.AccountId,
		AccountType:  ps.AccountType,
		TransactTime: ps.TransAt,
		AlgoId:       ps.AlgoId,
		AlgoType:     ps.AlgoType,
		Provider:     ps.Provider,

		AssessScore: float64(ps.AssessScore),
		RiskScore:   float64(ps.RiskScore),
		Progress:    float64(ps.ProgressScore),
		SourceFrom:  ps.SourceFrom,
	}
	result := d.DB.Create(&v)
	if result.Error != nil || result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

func (d *DefaultAlgoTimeLineOrig) UpdateAlgoTimeLine(ctx context.Context, ps *global.ProfileSum) error {
	result := d.DB.Model(models.TbAlgoTimeLineOrig{}).
		Where("batch_no=? and account_id=? and transact_time=? and algo_id=?",
			ps.BatchNo, ps.AccountId, ps.TransAt, ps.AlgoId).
		Updates(models.TbAlgoTimeLineOrig{
			AssessScore: float64(ps.AssessScore),
			RiskScore:   float64(ps.RiskScore),
			Progress:    float64(ps.ProgressScore),
			UpdateTime:  time.Now(),
		})
	return result.Error
}

// GetAlgoTimeLineByAllUser 多日分析里的非跨天， 算法分析， dashboard场景
func (d *DefaultAlgoTimeLineOrig) GetAlgoTimeLineByAllUser(ctx context.Context, date int64, algoId int, userId string, userType int32) ([]*models.TbAlgoTimeLineOrig, error) {
	var out []*models.TbAlgoTimeLineOrig
	result := d.DB.Model(models.TbAlgoTimeLineOrig{}).Select("transact_time, assess_score,progress,risk_score")
	if userType == global.UserTypeAdmin {
		result.Where("date =? and account_type=4 and algo_id = ?", date, algoId)
	} else {
		result.Where("date =? and account_id=? and algo_id = ?", date, userId, algoId)
	}
	result.Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

// GetImportAlgoTimeLine 订单导入场景，根据批次号查询
func (d *DefaultAlgoTimeLineOrig) GetImportAlgoTimeLine(ctx context.Context, date int64, batchNo int64, userId string, userType int32) ([]*models.TbAlgoTimeLine, error) {
	var out []*models.TbAlgoTimeLine
	result := d.DB.Model(models.TbAlgoTimeLineOrig{}).Select("transact_time, assess_score,progress,risk_score")
	if userType == global.UserTypeAdmin {
		result.Where("account_type=4")
	} else {
		result.Where("account_id=?", userId)
	}
	if batchNo != 0 {
		result.Where("batch_no = ?", batchNo)
	}
	result.Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

//func (d *DefaultAlgoTimeLineOrig) GetAlgoTimeLineCrossDay(ctx context.Context, start, end int64, UserId string, algoIds int32) ([]*models.TbAlgoTimeLineOrig, error) {
// select algo_id, date,avg(assess_score) from tb_algo_time_line where account_id ='aUser0000229' and date between 20220920 and 20220930
//group by date, algo_id order by algo_id, date;
//}

// GetMultiTimeLine 求多日的分数信息   跨天
func (d *DefaultAlgoTimeLineOrig) GetMultiTimeLine(ctx context.Context, start, end int64, userId string, userType int32, algoId int) ([]*models.TbAlgoTimeLineOrig, error) {
	var out []*models.TbAlgoTimeLineOrig
	result := d.DB.Model(models.TbAlgoTimeLineOrig{}).Select("date, avg(assess_score) as assess_score, avg(progress) as progress," +
		"avg(risk_score) as risk_score")
	if userType == global.UserTypeAdmin {
		result.Where("date between ? and ? and account_type=4 and algo_id = ?", start, end, algoId)
	} else {
		result.Where("date between ? and ? and account_id=? and algo_id = ?", start, end, userId, algoId)
	}
	result.Group("date").Order("date").Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

// GetMultiTimeLineBatch 绩效分析 跨天分析，根据批次号查询
func (d *DefaultAlgoTimeLineOrig) GetMultiTimeLineBatch(ctx context.Context, batchNo int64, userId string, userType int32, algoId int32) ([]*models.TbAlgoTimeLine, error) {
	var out []*models.TbAlgoTimeLine
	result := d.DB.Model(models.TbAlgoTimeLineOrig{}).Select("date, avg(assess_score) as assess_score, avg(progress) as progress," +
		"avg(risk_score) as risk_score")
	if userType == global.UserTypeAdmin {
		result.Where("batch_no =? and account_type=4 and algo_id = ?", batchNo, algoId)
	} else {
		result.Where("batch_no =? and account_id=? and algo_id = ?", batchNo, userId, algoId)
	}
	result.Group("date").Order("date").Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

func (d *DefaultAlgoTimeLineOrig) ReloadTimeLine(ctx context.Context, date int64) ([]*models.TbAlgoTimeLineOrig, error) {
	var out []*models.TbAlgoTimeLineOrig
	result := d.DB.Model(models.TbAlgoTimeLineOrig{}).Select("date,account_id,transact_time,algo_id,assess_score,risk_score,progress").
		Where("date=?", date).Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

func (d *DefaultAlgoTimeLineOrig) GetDataByTimeLineKey(ctx context.Context, date int64, userId string, algoId int) (models.TbAlgoTimeLineOrig, error) {
	var out models.TbAlgoTimeLineOrig
	result := d.DB.Model(models.TbAlgoTimeLineOrig{}).Select("account_id,transact_time").
		Where("account_id=? and algo_id=? and transact_time=?", userId, algoId, date).
		Order("transact_time desc").Limit(1).Find(&out)
	if result.Error != nil {
		return models.TbAlgoTimeLineOrig{}, result.Error
	}
	return out, nil
}

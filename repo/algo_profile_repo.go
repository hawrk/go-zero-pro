// Package repo
/*
 Author: hawrkchen
 Date: 2022/6/23 11:22
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

type AlgoProfileRepo interface {
	CreateAlgoProfile(ctx context.Context, profile *global.Profile) error
	UpdateAlgoProfile(ctx context.Context, profile *global.Profile) error
	GetProfiles(ctx context.Context, in *proto.AlgoProfileReq, algoType int32) ([]*models.TbAlgoProfile, int64, error)
	GetWinRatioProfile(ctx context.Context, userId string, userType int32, algoId int32, start, end int64) ([]*models.TbAlgoProfile, error)
	GetDailyStock(ctx context.Context, userId string, userType int32, algoId int32, start, end int64) (int64, error)

	//reload
	ReloadProfiles(ctx context.Context, date int64) ([]*models.TbAlgoProfile, error)

	// exception process
	GetDataByProfileKey(ctx context.Context, date int64, userId string, algoId int, secId string, algoOrderId int64) (models.TbAlgoProfile, error)
}

type DefaultAlgoProfile struct {
	DB *gorm.DB
}

func NewAlgoProfile(conn *gorm.DB) AlgoProfileRepo {
	return &DefaultAlgoProfile{
		DB: conn,
	}
}

func (d *DefaultAlgoProfile) CreateAlgoProfile(ctx context.Context, profile *global.Profile) error {
	p := models.TbAlgoProfile{
		Date:        int(profile.Date),
		BatchNo:     profile.BatchNo,
		AccountId:   profile.AccountId,
		AccountName: profile.AccountName,
		AccountType: profile.AccountType,
		Provider:    profile.Provider,
		AlgoId:      profile.AlgoId,
		AlgoName:    profile.AlgoName,
		AlgoType:    profile.AlgoType,
		SecId:       profile.SecId,
		SecName:     profile.SecName,
		AlgoOrderId: profile.AlgoOrderId,
		Industry:    profile.Industry,
		FundType:    profile.FundType,
		Liquidity:   profile.Liquidity,

		TradeCost:        profile.TotalTradeVol, // 交易成本
		TotalTradeAmount: profile.TotalT0Cost,   // 双边总交易额
		TotalTradeFee:    profile.TotalCharge,
		CrossFee:         profile.TotalCrossFee,
		// 2023-06-14 add
		AvgTradePrice:  profile.AvgTradePrice,   // 执行均价
		AvgArrivePrice: profile.AvgEntrustPrice, // 到达均价
		Pwp:            profile.PWP,
		AlgoDuration:   profile.TwapTotalDur,
		Twap:           profile.TWAP,
		TwapDev:        profile.TwapDev,
		Vwap:           profile.VWAP,
		VwapDev:        profile.VwapDev,
		// end
		ProfitAmount:   profile.ProfitAmount,
		ProfitRate:     profile.ProfitRate,
		CancelRate:     profile.CancelRate,
		ProgressRate:   profile.Progress,
		MiniSplitOrder: int(profile.MiniSplitOrder),
		MiniJointRate:  profile.MinJointRate,
		WithdrawRate:   profile.WithdrawRate,

		//EntrustVwap:       float64(profile.VwapEntrust),
		DealEffi:       profile.DealEffi,     // 成交效率
		AlgoOrderFit:   profile.AlgoOrderFit, // 母单贴合度
		TradeVolFit:    profile.TradeVolFit,  // 成交量贴合度
		TradeCount:     int(profile.TradeCount),
		TradeCountPlus: int(profile.TradeCountPlus),
		//AvgDeviation:      profile.VwapDevAvg,
		StandardDeviation: profile.VwapStdDev,
		PfRateStdDev:      profile.PfRateStdDev,
		DealFitStdDev:     profile.TradeVolFitStdDev, // 成交量贴合度标准差
		TimeFitStdDev:     profile.TimeFitStdDev,     // 时间贴合度标准差
		Factor:            profile.AssessFactor,
		OrderTime:         profile.CreateTime,
		SourceFrom:        profile.SourceFrom,
	}
	result := d.DB.Create(&p)
	if result.Error != nil || result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

func (d *DefaultAlgoProfile) UpdateAlgoProfile(ctx context.Context, profile *global.Profile) error {
	result := d.DB.Model(models.TbAlgoProfile{}).Where("date=? and account_id=? and algo_id=? and sec_id=? and algo_order_id=?",
		profile.Date, profile.AccountId, profile.AlgoId, profile.SecId, profile.AlgoOrderId).
		Updates(map[string]interface{}{
			"trade_cost":         profile.TotalTradeVol,
			"total_trade_fee":    profile.TotalCharge,
			"cross_fee":          profile.TotalCrossFee,
			"profit_amount":      profile.ProfitAmount,
			"profit_rate":        profile.ProfitRate,
			"cancel_rate":        profile.CancelRate,
			"progress_rate":      profile.Progress,
			"mini_split_order":   profile.MiniSplitOrder,
			"mini_joint_rate":    profile.MinJointRate,
			"withdraw_rate":      profile.WithdrawRate,
			"avg_trade_price":    profile.AvgTradePrice,   // 执行均价
			"avg_arrive_price":   profile.AvgEntrustPrice, // 到达均价
			"pwp":                profile.PWP,
			"algo_duration":      profile.TwapTotalDur,
			"twap":               profile.TWAP,
			"twap_dev":           profile.TwapDev,
			"vwap":               profile.VWAP,
			"vwap_dev":           profile.VwapDev,
			"deal_effi":          profile.DealEffi,     // 成交效率
			"algo_order_fit":     profile.AlgoOrderFit, // 母单贴合度
			"trade_vol_fit":      profile.TradeVolFit,  // 成交量贴合度
			"trade_count":        profile.TradeCount,
			"trade_count_plus":   profile.TradeCountPlus,
			"standard_deviation": profile.VwapStdDev,
			"pf_rate_std_dev":    profile.PfRateStdDev,
			"deal_fit_std_dev":   profile.TradeVolFitStdDev,
			"time_fit_std_dev":   profile.TimeFitStdDev,
			"factor":             profile.AssessFactor, // 绩效收益率
			"source_from":        profile.SourceFrom,
			"update_time":        time.Now(),
		})
		//Updates(models.TbAlgoProfile{
		//	//AccountName:       profile.AccountName,
		//	//AccountType:       profile.AccountType,
		//	//AlgoName:          profile.AlgoName,
		//	//AlgoType:          profile.AlgoType,
		//	TradeCost:        profile.TotalTradeVol,
		//	TotalTradeAmount: 0,
		//	TotalTradeFee:    profile.TotalCharge,
		//	CrossFee:         profile.TotalCrossFee,
		//	ProfitAmount:     profile.ProfitAmount,
		//	ProfitRate:       profile.ProfitRate,
		//	CancelRate:       profile.CancelRate,
		//	ProgressRate:     profile.Progress,
		//	MiniSplitOrder:   int(profile.MiniSplitOrder),
		//	MiniJointRate:    profile.MinJointRate,
		//	WithdrawRate:     profile.WithdrawRate,
		//	VwapDev:          profile.VwapDev,
		//	//EntrustVwap:      float64(profile.VwapEntrust),
		//	DealEffi:     profile.DealEffi,     // 成交效率
		//	AlgoOrderFit: profile.AlgoOrderFit, // 母单贴合度
		//	TradeVolFit:  profile.TradeVolFit,  // 成交量贴合度
		//	//TradeCount:        int(profile.VwapDevCnt),
		//	TradeCountPlus: int(profile.TradeCountPlus),
		//	//AvgDeviation:      profile.VwapDevAvg,
		//	StandardDeviation: profile.VwapStdDev,
		//	PfRateStdDev:      profile.PfRateStdDev,
		//	DealFitStdDev:     profile.TradeVolFitStdDev,
		//	TimeFitStdDev:     profile.TimeFitStdDev,
		//	Factor:            profile.AssessFactor, // 绩效收益率
		//	SourceFrom:        profile.SourceFrom,
		//	UpdateTime:        time.Now(),
		//})
	return result.Error
}

// GetProfiles 明细返回
func (d *DefaultAlgoProfile) GetProfiles(ctx context.Context, in *proto.AlgoProfileReq, algoType int32) ([]*models.TbAlgoProfile, int64, error) {
	userType := in.GetUserType()
	userId := in.GetUserId()
	provider := in.GetProvider()
	algoId := in.GetAlgoId()
	startTime := in.StartTime
	endTime := in.EndTime
	page := int(in.Page)
	limit := int(in.Limit)
	var out []*models.TbAlgoProfile
	var count int64
	// 用魔法打败魔法
	result := d.DB.Model(&models.TbAlgoProfile{}).Select("batch_no, account_id, account_name, account_type,provider,algo_type, algo_id,algo_name," +
		"sec_id,sec_name,algo_order_id,industry,fund_type,liquidity,trade_cost," +
		"total_trade_fee, cross_fee,avg_trade_price,avg_arrive_price,pwp,algo_duration,twap,twap_dev,vwap,vwap_dev," +
		"profit_amount,profit_rate," +
		"cancel_rate,mini_split_order, progress_rate," +
		"deal_effi,algo_order_fit,trade_vol_fit," +
		"mini_joint_rate, withdraw_rate,vwap_dev, standard_deviation,pf_rate_std_dev,deal_fit_std_dev,time_fit_std_dev,factor,order_time ")
	if userType != global.UserTypeAdmin {
		result.Where("account_id=?", userId)
	}
	if provider != "" {
		result.Where("provider=?", provider)
	}
	if algoType != 0 {
		result.Where("algo_type=?", algoType)
	}
	if algoId != 0 {
		result.Where("algo_id=?", algoId)
	}
	//if a := in.SourceFrom; a != 0 {
	//	result.Where("source_from=?", a)
	//}
	//if in.SourceFrom == global.SourceFromImport && in.BatchNo != 0 {
	//	result.Where("batch_no=?", in.BatchNo)
	//}
	result.Where("order_time between ? and ?", startTime, endTime).Order("order_time desc").
		Offset((page - 1) * limit).Limit(limit).Find(&out).Offset(-1).Limit(-1).Count(&count)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return out, count, nil
}

// GetWinRatioProfile 胜率分析下的数据拉取
func (d *DefaultAlgoProfile) GetWinRatioProfile(ctx context.Context, userId string, userType int32, algoId int32,
	start, end int64) ([]*models.TbAlgoProfile, error) {
	var out []*models.TbAlgoProfile
	result := d.DB.Model(&models.TbAlgoProfile{}).Select("date,algo_name, avg(profit_rate) as profit_rate," +
		"sum(trade_cost) as trade_cost," +
		"sum(profit_amount) as profit_amount," +
		"avg(withdraw_rate) as withdraw_rate," +
		" avg(progress_rate) as progress_rate,sum(trade_count) as trade_count,sum(trade_count_plus) as trade_count_plus")
	if userType == global.UserTypeAdmin {
		result.Where("date  between ? and ? and account_type=4 and algo_id=?", start, end, algoId)
	} else {
		result.Where("date between ? and ? and account_id=? and algo_id=?", start, end, userId, algoId)
	}
	result.Group("date").Order("date").Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

// GetDailyStock 取参与交易的股票数量
func (d *DefaultAlgoProfile) GetDailyStock(ctx context.Context, userId string, userType int32, algoId int32, start, end int64) (int64, error) {
	var count int64
	result := d.DB.Model(&models.TbAlgoProfile{}).Select("count(distinct(sec_id)) as count")
	if userType == global.UserTypeAdmin {
		result.Where("date  between ? and ? and account_type=4 and algo_id=?", start, end, algoId)
	} else {
		result.Where("date between ? and ? and account_id=? and algo_id=?", start, end, userId, algoId)
	}
	result.Find(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil

}

func (d *DefaultAlgoProfile) ReloadProfiles(ctx context.Context, date int64) ([]*models.TbAlgoProfile, error) {
	var out []*models.TbAlgoProfile
	result := d.DB.Model(&models.TbAlgoProfile{}).Select("date,account_id,algo_id,trade_cost,total_trade_fee,cross_fee,profit_amount,profit_rate,"+
		"cancel_rate,progress_rate,mini_split_order,mini_joint_rate,withdraw_rate,vwap_dev,entrust_vwap,trade_count,avg_deviation,"+
		"standard_deviation,pf_rate_std_dev,factor,order_time").Where("date=?", date).Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

// GetDataByProfileKey 异常处理，从表中取一条最近时间的记录，与异常队列的数据比较
func (d *DefaultAlgoProfile) GetDataByProfileKey(ctx context.Context, date int64,
	userId string, algoId int, secId string, algoOrderId int64) (models.TbAlgoProfile, error) {
	var out models.TbAlgoProfile
	result := d.DB.Model(models.TbAlgoProfile{}).Select("account_id").
		Where("date=? and account_id=? and algo_id=? and sec_id=? and algo_order_id=?", date, userId, algoId, secId, algoOrderId).
		Order("update_time desc").Limit(1).Find(&out)
	if result.Error != nil {
		return models.TbAlgoProfile{}, result.Error
	}
	return out, nil
}

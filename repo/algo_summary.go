// Package repo
/*
 Author: hawrkchen
 Date: 2022/7/18 16:43
 Desc:
*/
package repo

import (
	"algo_assess/global"
	"algo_assess/models"
	"algo_assess/pkg/tools"
	"context"
	"gorm.io/gorm"
	"time"
)

type AlgoSummaryRepo interface {
	CreateAlgoSummary(ctx context.Context, profile *global.ProfileSum) error
	UpdateAlgoSummary(ctx context.Context, profile *global.ProfileSum) error
	GetAlgoSummary(ctx context.Context, userId string, userType int32, date int64, algoId int32) ([]*models.TbAlgoSummary, error)
	GetImportAlgoSummary(ctx context.Context, userId string, userType int32, date int64, batchNo int64) ([]*models.TbAlgoSummary, error)
	GetUserOrderSummary(ctx context.Context, date int64, userId string, userType int32) ([]*models.TbAlgoSummary, []*UserOrderSummary, error)
	GetAlgoOrderSummary(ctx context.Context, date int64, userId string, userType int32, algoType int32, page, limit int) ([]*AlgoOrderSummary, int64, error)
	GetTopRankAlgoSummary(ctx context.Context, date int64, userId string, userType int32, algoType int) ([]*models.TbAlgoSummary, error)
	GetAlgoSummaryByAlgoIds(ctx context.Context, algoId []int32, date int64, userId string, userType int32) ([]*models.TbAlgoSummary, error)
	GetCrossDaySummaryByAlgoIds(ctx context.Context, algoId []int32, start, end int64, userId string, userType int32) ([]*AvgSummary, error)
	GetCrossDaySummaryByAlgoId(ctx context.Context, algoId int32, start, end int64, userId string, userType int32) ([]*models.TbAlgoSummary, error)
	GetAlgoNameByRanking(ctx context.Context, date int64, userId string) ([]string, error)
	GetTotalScoreRanking(ctx context.Context, date int64, userId string, page, limit int) ([]*models.TbAlgoSummary, int64, error)
	GetAlgoSummaryCnt(ctx context.Context, date int64, userId, userName string, role int) (int32, int32, error)
	GetAlgoSummaryCntByAdmin(ctx context.Context, date int64) (int32, int32, error)
	GetCumsumList(ctx context.Context, date int64, role int) ([]int, error)
	GetCumsumListCrossDay(ctx context.Context, start, end int64, role int) ([]int, error)
	GetUserProfit(ctx context.Context, date int64, userId string, userType int32, algoId int) (models.TbAlgoSummary, error)
	GetSignalSummary(ctx context.Context, start, end int64, userId string, userType int32, algoId int) ([]*models.TbAlgoSummary, error)

	// reload
	ReloadSummary(ctx context.Context, date int64) ([]*models.TbAlgoSummary, error)
	GetDefaultAlgoId(ctx context.Context, date int64, userId string, userType int32) (models.TbAlgoSummary, error)
	GetAdvanceDefaultAlgo(ctx context.Context, start, end int64, userId string, userType int32) (models.TbAlgoSummary, error)
	// exception process perf fix
	GetDataBySummaryKey(ctx context.Context, date int64, userId string, algoId int) (models.TbAlgoSummary, error)
}

type DefaultAlgoSummary struct {
	DB *gorm.DB
}

func NewAlgoSummary(conn *gorm.DB) AlgoSummaryRepo {
	return &DefaultAlgoSummary{
		DB: conn,
	}
}

func (d *DefaultAlgoSummary) CreateAlgoSummary(ctx context.Context, profile *global.ProfileSum) error {
	v := &models.TbAlgoSummary{
		Date:            int(profile.Date),
		BatchNo:         profile.BatchNo,
		UserId:          profile.AccountId,
		AccountType:     profile.AccountType,
		AlgoId:          profile.AlgoId,
		AlgoName:        profile.AlgoName,
		AlgoType:        profile.AlgoType,
		Provider:        profile.Provider,
		OrderNum:        int(profile.OrderNum),
		EntrustQty:      profile.EntrustQty,
		DealQty:         profile.DealQty,
		OrderAmount:     profile.TotalLastCost,
		BuyAmount:       profile.TotalBuyVol,
		SellAmount:      profile.TotalSellVol,
		Perfit:          float64(profile.ProfitAmount),
		PerfitRate:      profile.ProfitRate,
		AssessScore:     profile.AssessScore,
		ProgressScore:   profile.ProgressScore,
		StableScore:     profile.StabilityScore,
		RiskScore:       profile.RiskScore,
		EconomyScore:    profile.EconomyScore,
		CumsumScore:     profile.TotalScore,
		FundRateJson:    tools.Json2Str(profile.FundPercent),
		TradeVolJson:    tools.Json2Str(profile.TradeVolVal),
		StockTypeJson:   tools.Json2Str(profile.StockTypeVal),
		TradeDirectJson: tools.Json2Str(profile.TradeVolDict),
		SourceFrom:      profile.SourceFrom,
	}
	result := d.DB.Create(&v)
	if result.Error != nil || result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

func (d *DefaultAlgoSummary) UpdateAlgoSummary(ctx context.Context, profile *global.ProfileSum) error {
	result := d.DB.Model(models.TbAlgoSummary{}).Where("date=? and user_id=? and algo_id=?", profile.Date, profile.AccountId, profile.AlgoId).
		Updates(map[string]interface{}{
			"order_num":         profile.OrderNum,
			"entrust_qty":       profile.EntrustQty,
			"deal_qty":          profile.DealQty,
			"order_amount":      profile.TotalLastCost,
			"buy_amount":        profile.TotalBuyVol,
			"sell_amount":       profile.TotalSellVol,
			"perfit":            profile.ProfitAmount,
			"perfit_rate":       profile.ProfitRate,
			"assess_score":      profile.AssessScore,
			"progress_score":    profile.ProgressScore,
			"stable_score":      profile.StabilityScore,
			"risk_score":        profile.RiskScore,
			"economy_score":     profile.EconomyScore,
			"cumsum_score":      profile.TotalScore,
			"fund_rate_json":    tools.Json2Str(profile.FundPercent),
			"trade_vol_json":    tools.Json2Str(profile.TradeVolVal),
			"stock_type_json":   tools.Json2Str(profile.StockTypeVal),
			"trade_direct_json": tools.Json2Str(profile.TradeVolDict),
			"source_from":       profile.SourceFrom,
			"update_time":       time.Now(),
		})
		//Updates(models.TbAlgoSummary{
		//	OrderNum:        int(profile.OrderNum),
		//	EntrustQty:      profile.EntrustQty,
		//	DealQty:         profile.DealQty,
		//	OrderAmount:     profile.TotalLastCost,
		//	BuyAmount:       profile.TotalBuyVol,
		//	SellAmount:      profile.TotalSellVol,
		//	Perfit:          float64(profile.ProfitAmount),
		//	PerfitRate:      profile.ProfitRate,
		//	AssessScore:     profile.AssessScore,
		//	ProgressScore:   profile.ProgressScore,
		//	StableScore:     profile.StabilityScore,
		//	RiskScore:       profile.RiskScore,
		//	EconomyScore:    profile.EconomyScore,
		//	CumsumScore:     profile.TotalScore,
		//	FundRateJson:    tools.Json2Str(profile.FundPercent),
		//	TradeVolJson:    tools.Json2Str(profile.TradeVolVal),
		//	StockTypeJson:   tools.Json2Str(profile.StockTypeVal),
		//	TradeDirectJson: tools.Json2Str(profile.TradeVolDict),
		//	SourceFrom:      profile.SourceFrom,
		//	UpdateTime:      time.Now(),
		//})
	return result.Error
}

// GetAlgoSummary 查评分的数据
func (d *DefaultAlgoSummary) GetAlgoSummary(ctx context.Context, userId string, userType int32, date int64, algoId int32) ([]*models.TbAlgoSummary, error) {
	var out []*models.TbAlgoSummary
	result := d.DB.Model(models.TbAlgoSummary{}).Select("assess_score,progress_score,stable_score,risk_score,economy_score,cumsum_score," +
		"fund_rate_json,trade_direct_json,stock_type_json,trade_vol_json")
	if userType == global.UserTypeAdmin {
		result.Where("date=?  and account_type=4 and algo_id=?", date, algoId)
	} else {
		result.Where("date=? and user_id=? and algo_id=?", date, userId, algoId)
	}
	result.Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

// GetImportAlgoSummary 取订单导入的算法动态信息，根据批次号获取
func (d *DefaultAlgoSummary) GetImportAlgoSummary(ctx context.Context, userId string, userType int32, date int64, batchNo int64) ([]*models.TbAlgoSummary, error) {
	var out []*models.TbAlgoSummary
	result := d.DB.Model(models.TbAlgoSummary{}).Select("assess_score,progress_score,stable_score,risk_score,economy_score,cumsum_score," +
		"fund_rate_json,trade_direct_json,stock_type_json,trade_vol_json")
	if userType == global.UserTypeAdmin {
		result.Where("date=?  and account_type=4 ", date)
	} else {
		result.Where("date=? and user_id=?", date, userId)
	}
	if batchNo != 0 {
		result.Where("batch_no =?", batchNo)
	}

	result.Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

// GetUserOrderSummary 查用户订单汇总信息
func (d *DefaultAlgoSummary) GetUserOrderSummary(ctx context.Context, date int64, userId string, userType int32) ([]*models.TbAlgoSummary, []*UserOrderSummary, error) {
	var out []*models.TbAlgoSummary
	result := d.DB.Model(models.TbAlgoSummary{}).Select("fund_rate_json,trade_direct_json").Where("date = ?", date)
	if userType == global.UserTypeAdmin {
		result.Where("user_id=?", global.AdminUserId)
	} else {
		result.Where("user_id=?", userId)
	}
	result.Find(&out)

	if result.Error != nil {
		return nil, nil, result.Error
	}
	var orderSummary []*UserOrderSummary
	orderResult := d.DB.Model(models.TbAlgoSummary{}).Select("count(distinct(user_id)) as user_num, " +
		"count(distinct(provider)) as provider_num, sum(order_amount) as trade_amount, sum(order_num) as order_num ," +
		"sum(entrust_qty) as entrust_qty, sum(deal_qty) as deal_qty")
	if userType == global.UserTypeAdmin { // 超级管理员时，对普通用户进行汇总统计
		orderResult.Where("date=? and account_type =1", date)
	} else { // TODO: 这里如果是算法厂商的话，可能会有问题
		orderResult.Where("date=? and user_id=?", date, userId)
	}
	orderResult.Find(&orderSummary)
	if orderResult.Error != nil {
		return nil, nil, orderResult.Error
	}
	return out, orderSummary, nil
}

// GetAlgoOrderSummary 根据算法类型查交易汇总信息
func (d *DefaultAlgoSummary) GetAlgoOrderSummary(ctx context.Context, date int64, userId string, userType int32,
	algoType int32, page, limit int) ([]*AlgoOrderSummary, int64, error) {
	var out []*AlgoOrderSummary
	var count int64
	result := d.DB.Model(models.TbAlgoSummary{}).Select("provider, count(distinct(user_id)) as user_num, " +
		"sum(order_amount) as trade_amount, sum(order_num) as order_num, avg(perfit_rate) as profit_rate," +
		"sum(buy_amount) as buy_amount, sum(sell_amount) as sell_amount ")
	if userType == global.UserTypeAdmin {
		result.Where("date=? and account_type=1 and algo_type =?", date, algoType)
	} else {
		result.Where("date=? and user_id=? and algo_type =?", date, userId, algoType)
	}
	result.Group("provider").Offset(page - 1).Limit(limit).Find(&out).Offset(-1).Limit(-1).Count(&count)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return out, count, nil
}

// GetTopRankAlgoSummary 取dashboard 里top 四条算法 ---需要区分算法类型
func (d *DefaultAlgoSummary) GetTopRankAlgoSummary(ctx context.Context, date int64, userId string, userType int32, algoType int) ([]*models.TbAlgoSummary, error) {
	var out []*models.TbAlgoSummary
	result := d.DB.Model(models.TbAlgoSummary{}).Select("user_id,algo_id, algo_name,assess_score,cumsum_score")
	if userType == global.UserTypeAdmin {
		result.Where("date = ? and account_type=4 and algo_type=?", date, algoType)
	} else {
		result.Where("date = ? and user_id=? and algo_type=?", date, userId, algoType)
	}
	result.Order("assess_score").Limit(4).Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

// GetAlgoSummaryByAlgoIds 根据 algo_id 取汇总信息，对比分析
func (d *DefaultAlgoSummary) GetAlgoSummaryByAlgoIds(ctx context.Context, algoId []int32, date int64,
	userId string, userType int32) ([]*models.TbAlgoSummary, error) {
	var out []*models.TbAlgoSummary
	result := d.DB.Model(models.TbAlgoSummary{}).Select("algo_id, algo_name,assess_score,progress_score,stable_score,risk_score," +
		"economy_score,cumsum_score")
	if userType == global.UserTypeAdmin {
		result.Where("date =? and account_type=4 and algo_id in ?", date, algoId)
	} else {
		result.Where("date =? and user_id=? and algo_id in ?", date, userId, algoId)
	}
	result.Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

// GetCrossDaySummaryByAlgoIds 根据算法ID列表查跨天的汇总数据,求平均值
func (d *DefaultAlgoSummary) GetCrossDaySummaryByAlgoIds(ctx context.Context, alogId []int32, start, end int64,
	userId string, userType int32) ([]*AvgSummary, error) {
	// DB 里分数字段都是整型，这时求平均数后都是浮点型，所以先用自定义的数据结构
	var out []*AvgSummary
	result := d.DB.Model(models.TbAlgoSummary{}).Select("algo_id, any_value(algo_name) as algo_name, avg(assess_score) as assess_score," +
		"avg(progress_score) as progress_score,avg(stable_score) as stable_score, avg(risk_score) as risk_score,avg(economy_score) as economy_score, " +
		"avg(cumsum_score) as cumsum_score")
	if userType == global.UserTypeAdmin {
		result.Where("date between ? and ? and account_type=4 and algo_id in(?)", start, end, alogId)
	} else {
		result.Where("date between ? and ? and user_id=? and algo_id in(?)", start, end, userId, alogId)
	}
	result.Group("algo_id").Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

// GetCrossDaySummaryByAlgoId 根据单个算法求跨天的汇总数据，需要查询Json字段无法使用avg,只能查多笔再计算（算法动态跨天场景）
func (d *DefaultAlgoSummary) GetCrossDaySummaryByAlgoId(ctx context.Context, algoId int32, start, end int64,
	userId string, userType int32) ([]*models.TbAlgoSummary, error) {
	var out []*models.TbAlgoSummary
	result := d.DB.Model(models.TbAlgoSummary{}).Select("assess_score,progress_score,stable_score,risk_score,economy_score,cumsum_score," +
		"fund_rate_json,trade_direct_json,stock_type_json,trade_vol_json")
	if userType == global.UserTypeAdmin {
		result.Where("date between ? and ? and account_type=4 and algo_id=?", start, end, algoId)
	} else {
		result.Where("date between ? and ? and user_id=? and algo_id=?", start, end, userId, algoId)
	}
	result.Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

// GetCrossDayGroupByAlgoIds 根据多日请求返回日期和算法的分组，与 GetCrossDaySummaryByAlgoIds 配套使用
//func (d *DefaultAlgoSummary) GetCrossDayGroupByAlgoIds(ctx context.Context, alogId []int32, start, end int64, userId string) ([]*models.TbAlgoSummary, error) {
//	var out []*models.TbAlgoSummary
//	result := d.DB.Model(models.TbAlgoSummary{}).Select("algo_id,date, any_value(algo_name) as algo_name").
//		Where("date between ? and ? and user_id=? and algo_id in(?)", start, end, userId, alogId).
//		Group("algo_id").Group("date").Find(&out)
//	if result.Error != nil {
//		return nil, result.Error
//	}
//	return out, nil
//}

// GetAlgoNameByRanking 根据日期取当天有交易的算法名称列表，按算法从高到低绩效排名
func (d *DefaultAlgoSummary) GetAlgoNameByRanking(ctx context.Context, date int64, userId string) ([]string, error) {
	var out []string
	result := d.DB.Model(models.TbAlgoSummary{}).Select("algo_name").Where("date=? and user_id=?", date, userId).
		Order("assess_score").Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

func (d *DefaultAlgoSummary) GetTotalScoreRanking(ctx context.Context, date int64, userId string, page, limit int) ([]*models.TbAlgoSummary, int64, error) {
	var out []*models.TbAlgoSummary
	var count int64
	result := d.DB.Model(models.TbAlgoSummary{}).Select("algo_name,cumsum_score").Where("date=? and user_id=?", date, userId).Order("cumsum_score desc").
		Offset(page - 1).Limit(limit).Find(&out).Offset(-1).Limit(-1).Count(&count)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return out, count, nil
}

func (d *DefaultAlgoSummary) GetAlgoSummaryCnt(ctx context.Context, date int64, userId, userName string, role int) (int32, int32, error) {
	type TradeCount struct {
		AlgoCnt     int32
		ProviderCnt int32
	}
	var cnt TradeCount
	// 需要根据角色权限限制查询条件
	var result *gorm.DB
	if role == 1 || role == 3 { // 普通用户
		result = d.DB.Model(models.TbAlgoSummary{}).Select("count(distinct(algo_id)) as algo_cnt ,count(distinct(provider)) as provider_cnt").
			Where("date=? and user_id=?", date, userId).Find(&cnt)
	} else if role == 2 { // 算法厂商
		result = d.DB.Model(models.TbAlgoSummary{}).Select("count(distinct(algo_id)) as algo_cnt ,count(distinct(provider)) as provider_cnt").
			Where("date=? and provider=? and account_type =1", date, userName).Find(&cnt)
	} else { // 如果都不是，那就是绩效这边管理菜单权限的用户
		return cnt.AlgoCnt, cnt.ProviderCnt, nil
	}
	//} else if role == 3 { // 管理员
	//	result = d.DB.Model(models.TbAlgoSummary{}).Select("count(distinct(algo_id)) as algo_cnt ,count(distinct(provider)) as provider_cnt").
	//		Where("date=? and account_type =1", date).Find(&cnt)
	//}
	if result.Error != nil {
		return 0, 0, result.Error
	}
	return cnt.AlgoCnt, cnt.ProviderCnt, nil
}

// GetAlgoSummaryCntByAdmin 根据管理员权限取动态算法统计数据
func (d *DefaultAlgoSummary) GetAlgoSummaryCntByAdmin(ctx context.Context, date int64) (int32, int32, error) {
	type TradeCount struct {
		AlgoCnt     int32
		ProviderCnt int32
	}
	var cnt TradeCount
	// 需要根据角色权限限制查询条件
	result := d.DB.Model(models.TbAlgoSummary{}).Select("count(distinct(algo_id)) as algo_cnt ,count(distinct(provider)) as provider_cnt").
		Where("date=? and account_type =1", date).Find(&cnt)
	if result.Error != nil {
		return 0, 0, result.Error
	}
	return cnt.AlgoCnt, cnt.ProviderCnt, nil

}

// GetCumsumList 取综合评分排名列表--分用户类型
func (d *DefaultAlgoSummary) GetCumsumList(ctx context.Context, date int64, role int) ([]int, error) {
	var out []int
	result := d.DB.Model(models.TbAlgoSummary{}).Select("distinct(cumsum_score) as cumsum_score").Where("account_type=? and date=?", role, date).
		Order("cumsum_score desc").Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

// GetCumsumListCrossDay 取综合评分排名列表 跨天
func (d *DefaultAlgoSummary) GetCumsumListCrossDay(ctx context.Context, start, end int64, role int) ([]int, error) {
	var out []int
	result := d.DB.Model(models.TbAlgoSummary{}).Select("cumsum_score").Where("account_type=? and date between ? and ?", role, start, end).
		Group("algo_id").Order("cumsum_score desc").Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

// GetSignalSummary 高阶分析，信号分析数据获取
func (d *DefaultAlgoSummary) GetSignalSummary(ctx context.Context, start, end int64, userId string, userType int32, algoId int) ([]*models.TbAlgoSummary, error) {
	var out []*models.TbAlgoSummary
	result := d.DB.Model(models.TbAlgoSummary{}).Select("date,entrust_qty,deal_qty,order_num")
	if userType == global.UserTypeAdmin {
		result.Where("date between ? and ? and account_type=4 and algo_id=?", start, end, algoId)
	} else {
		result.Where("date between ? and ? and user_id=? and algo_id=?", start, end, userId, algoId)
	}
	result.Order("date").Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

func (d *DefaultAlgoSummary) ReloadSummary(ctx context.Context, date int64) ([]*models.TbAlgoSummary, error) {
	var out []*models.TbAlgoSummary
	result := d.DB.Model(models.TbAlgoSummary{}).Select("date,user_id,account_type,algo_id,algo_name,algo_type,provider,"+
		"order_num,entrust_qty,deal_qty,order_amount,buy_amount,sell_amount,perfit,perfit_rate,assess_score,progress_score,"+
		"stable_score,risk_score,economy_score,cumsum_score,fund_rate_json,trade_vol_json,stock_type_json,trade_direct_json").
		Where("date=?", date).Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

func (d *DefaultAlgoSummary) GetDefaultAlgoId(ctx context.Context, date int64, userId string, userType int32) (models.TbAlgoSummary, error) {
	var out models.TbAlgoSummary
	result := d.DB.Model(models.TbAlgoSummary{}).Select("algo_id, algo_name,algo_type,provider")
	if userType == global.UserTypeAdmin {
		result.Where("date=?", date)
	} else {
		result.Where("date=? and user_id=?", date, userId)
	}
	result.Limit(1).Find(&out)
	if result.Error != nil {
		return models.TbAlgoSummary{}, result.Error
	}
	return out, nil
}

// GetAdvanceDefaultAlgo 需要跨天请求数据
func (d *DefaultAlgoSummary) GetAdvanceDefaultAlgo(ctx context.Context, start, end int64, userId string, userType int32) (models.TbAlgoSummary, error) {
	var out models.TbAlgoSummary
	result := d.DB.Model(models.TbAlgoSummary{}).Select("algo_id, algo_name,algo_type,provider")
	if userType == global.UserTypeAdmin {
		result.Where("date between ? and ?", start, end)
	} else {
		result.Where("date between ? and ? and user_id=?", start, end, userId)
	}
	result.Limit(1).Find(&out)
	if result.Error != nil {
		return models.TbAlgoSummary{}, result.Error
	}
	return out, nil
}

func (d *DefaultAlgoSummary) GetDataBySummaryKey(ctx context.Context, date int64, userId string, algoId int) (models.TbAlgoSummary, error) {
	var out models.TbAlgoSummary
	result := d.DB.Model(models.TbAlgoSummary{}).Select("user_id").
		Where("date=? and user_id=? and algo_id=?", date, userId, algoId).
		Order("update_time desc").Limit(1).Find(&out)
	if result.Error != nil {
		return models.TbAlgoSummary{}, result.Error
	}
	return out, nil
}

// GetUserProfit 取用户画像下页面的资金交易信息
func (d *DefaultAlgoSummary) GetUserProfit(ctx context.Context, date int64,
	userId string, userType int32, algoId int) (models.TbAlgoSummary, error) {
	var out models.TbAlgoSummary
	result := d.DB.Model(models.TbAlgoSummary{}).Select("order_amount,perfit,order_num,entrust_qty,deal_qty")
	if userType == global.UserTypeAdmin {
		result.Where("date=? and account_type=4 and algo_id=?", date, algoId)
	} else {
		result.Where("date=? and user_id=? and algo_id=?", date, userId, algoId)
	}
	result.Find(&out)

	if result.Error != nil {
		return models.TbAlgoSummary{}, result.Error
	}
	return out, nil
}

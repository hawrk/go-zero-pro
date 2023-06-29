// Package consumer
/*
 Author: hawrkchen
 Date: 2022/12/15 17:58
 Desc:
*/
package consumer

import (
	"algo_assess/assess-mq-server/internal/dao"
	"algo_assess/global"
	"algo_assess/pkg/tools"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"math"
)

// CalculateProfile 计算画像明细
func CalculateProfile(head *global.ProfileHead, in *global.Profile, data *global.ChildOrderData,
	qty *global.OrderTradeQty, redis *redis.Redis) {

	in.ProfileHead = *head
	in.LastQty += data.LastQty                     // 总成交数量
	in.TotalTradeVol += data.LastPx * data.LastQty // 交易量
	in.TotalCharge += data.TotalFee                // 总手续费
	in.TotalCrossFee += 100                        //TODO: 流量费，暂时先加0.01 (跟总线保持一致，除以10000，加100就是0.01元）
	if data.LastQty > 0 {
		in.DealCount++ // 算滑点个数
	}
	if data.TradeSide == global.TradeSideBuy {
		in.TotalBuyVol += data.LastPx * data.LastQty // 买入总量
	} else if data.TradeSide == global.TradeSideSell {
		in.TotalSellVol += data.LastPx * data.LastQty // 卖出总量
	}

	var profitRate float64                       // 收益率
	in.TradeCount++                              // 交易次数
	if data.AlgorithmType == global.AlgoTypeT0 { //TO算法
		//T0区分卖出和买入
		in.TotalT0Fee += data.TotalFee
		if data.TradeSide == global.TradeSideSell { // 卖出总成本
			in.TotalSellQty += data.LastQty // 卖出总量
			in.SellCost += data.LastPx * data.LastQty
			if in.TotalSellQty > 0 {
				in.AvgSellPrice = in.SellCost / in.TotalSellQty
			}
		} else if data.TradeSide == global.TradeSideBuy { // 买入总成本
			in.TotalBuyQty += data.LastQty // 买入总量
			in.BuyCost += data.LastPx * data.LastQty
			if in.TotalBuyQty > 0 {
				in.AvgBuyPrice = in.BuyCost / in.TotalBuyQty
			}
		}
		// 取当前 min 成交数量
		minDealQty := tools.MinInt(in.TotalSellQty, in.TotalBuyQty)
		// 算盈亏金额
		in.ProfitAmount = float64((in.AvgSellPrice-in.AvgBuyPrice)*minDealQty - in.TotalT0Fee)
		if in.ProfitAmount > 0 { //盈亏
			in.TradeCountPlus++
		}
		if data.MarketPrice != 0 { // 优先用本地的行情数据
			in.TotalEntrustCost += data.MarketPrice * data.LastQty // 盈亏，用到达价算
		} else {
			in.TotalEntrustCost += data.ArrivePrice * data.LastQty // 盈亏，用到达价算
		}
		in.AvgEntrustPrice = float64(in.TotalEntrustCost) / float64(in.LastQty) // 到达价均价
		//总交易成本
		in.TotalT0Cost = in.SellCost + in.BuyCost                        // 双边总交易额
		in.AvgTradePrice = float64(in.TotalT0Cost) / float64(in.LastQty) // T0 母单执行均价

		in.PWP = (float64(in.TotalT0Cost) - in.ProfitAmount) / float64(in.LastQty) // PWP 价格 = （执行成本-盈亏）/总成交数量
		//in.TotalTradeCost += data.LastPx * data.LastQty
		// 收益率
		if in.TotalT0Cost > 0 {
			profitRate = in.ProfitAmount / (float64(in.TotalT0Cost) / 2)
		}
		// 绩效收益率
		in.AssessFactor = 1 * profitRate // T0的绩效收益率就是盈亏收益率
		// 平均交易均价
		if in.DealCount > 0 {
			in.AvgCost = float64(in.TotalT0Cost) / float64(in.DealCount)
		}

	} else if data.AlgorithmType == global.AlgoTypeSplit { // 拆单算法
		in.TotalSplitFee += data.TotalFee
		//in.TotalEntrustCost += data.Price * data.LastQty
		// 拆单到达价取本系统行情市场价,先从Redis取，取不到再从Mysql取
		if data.MarketPrice != 0 { // 优先用本地的行情数据
			in.TotalEntrustCost += data.MarketPrice * data.LastQty // 盈亏，用到达价算
		} else {
			in.TotalEntrustCost += data.ArrivePrice * data.LastQty // 盈亏，用到达价算
		}
		//in.TotalTradeCost += data.LastPx * data.LastQty // 交易总价
		if in.LastQty > 0 {
			in.AvgEntrustPrice = float64(in.TotalEntrustCost) / float64(in.LastQty) // 普通交易均价  （到达价均价)
			in.AvgTradePrice = float64(in.TotalTradeVol) / float64(in.LastQty)      // 成交均价
		}
		in.ProfitAmount = (in.AvgEntrustPrice-in.AvgTradePrice)*float64(in.LastQty) - float64(in.TotalSplitFee)
		if in.TotalTradeVol > 0 {
			profitRate = float64(in.ProfitAmount) / float64(in.TotalTradeVol) // 收益率
		}
		// 平均交易均价
		if in.DealCount > 0 {
			in.AvgCost = float64(in.TotalTradeVol) / float64(in.DealCount)
		}
		// PWP
		in.PWP = (float64(in.TotalTradeVol) - in.ProfitAmount) / float64(in.LastQty)
		// 算最大亏损金额， 取总交易亏损的最大值
		// 算出当前交易的亏损金额
		t := data.Price*data.LastQty - data.LastPx*data.LastQty // 当前交易的盈亏
		tm := tools.MinFloat(float64(t), in.ProfitAmount)       // 当前交易盈亏和累加交易盈亏取最小值
		if in.MaxLossAmount > tm {                              // 严格来说，in.MaxLossAmount 是小于等于0的数
			in.MaxLossAmount = tm
		}
		// 回撤比例
		if in.TotalTradeVol > 0 {
			in.WithdrawRate = math.Abs(in.MaxLossAmount) / float64(in.TotalTradeVol) * 100
		}
		// 算胜率，区分买卖方向
		if data.TradeSide == global.TradeSideSell { // 卖 成交价大于市场价，则胜
			if data.MarketPrice > 0 {
				if data.LastPx > data.MarketPrice {
					in.TradeCountPlus++
				}
			} else if data.ArrivePrice > 0 {
				if data.LastPx > data.ArrivePrice {
					in.TradeCountPlus++
				}
			}
		} else if data.TradeSide == global.TradeSideBuy { // 买 成交价小于市价，则胜
			if data.MarketPrice > 0 {
				if data.LastPx < data.MarketPrice {
					in.TradeCountPlus++
				}
			} else if data.ArrivePrice > 0 {
				if data.LastPx < data.ArrivePrice {
					in.TradeCountPlus++
				}
			}
		}

	}
	in.EntrustQty = qty.EntrustQty
	in.DealQty = qty.DealQty
	in.CancelQty = qty.CancelQty
	if qty.EntrustQty > 0 {
		// 算撤单率
		in.CancelRate = float64(qty.CancelQty) / float64(qty.EntrustQty)
		// 算完成度
		in.Progress = (float64(qty.DealQty) / float64(qty.EntrustQty)) * 100
		// 算母单贴合度
		if in.MiniDealOrder > data.LastQty && data.LastQty > 0 {
			in.MiniDealOrder = data.LastQty
		}
		in.AlgoOrderFit = float64(in.MiniDealOrder) / float64(qty.EntrustQty)
		// 算成交量贴合度
		in.TradeVolFit = 1 - math.Abs(float64(in.LastQty-qty.EntrustQty))/float64(qty.EntrustQty)
	}

	// 算最小拆单单位
	if in.MiniSplitOrder > data.OrderQty {
		in.MiniSplitOrder = data.OrderQty
	}
	// 算最小贴合度
	c, enTime := GetMinJointRate(redis, in.Progress, data)
	if in.MinJointRate == 0.00 || in.MinJointRate > c {
		in.MinJointRate = c
	}
	// 算成交效率
	durTime := data.UnixTime - enTime // 当前子单最新交易时间-母单委托时间
	if durTime > 0 {
		in.DealEffi = (float64(in.TotalTradeVol) / 10000) / float64(durTime) // 除以10000，与金额单位为元保持一致
	}
	// 计算TWAP   t1, t2, t3 成交p1， p2, p3 ,则计算公式 (（t1-t0)*p1+(t2-t1)*p2+(t3-t2)*p3) / (t1-t0)+(t2-t1)+(t3-t2)
	// 其中第一笔t0为母单的委托时间， 时间单位为秒
	if in.TwapStartTimePoint == 0 { // 如果是第一笔，则开始时间点取母单的委托时间
		in.TwapStartTimePoint = enTime
	}
	in.TwapDurTime = int64(math.Abs(float64(data.UnixTime - in.TwapStartTimePoint)))
	// 计算成交的TWAP
	in.TwapTotalTrade += in.TwapDurTime * data.LastPx
	in.TwapTotalDur += in.TwapDurTime // 母单有效时长
	if in.TwapTotalDur > 0 {
		in.TWAP = float64(in.TwapTotalTrade) / float64(in.TwapTotalDur) / 10000
	}
	// 计算市场价的TWAP
	if data.MarketPrice != 0 {
		in.TwapTotalMarketTrade += in.TwapDurTime * data.MarketPrice
	} else if data.ArrivePrice != 0 {
		in.TwapTotalMarketTrade += in.TwapDurTime * data.ArrivePrice
	}
	if in.TwapTotalDur > 0 {
		in.TWAPMarket = float64(in.TwapTotalMarketTrade) / float64(in.TwapTotalDur) / 10000
	}
	// 计算TWAP滑点
	in.TwapDev = in.TWAP - in.TWAPMarket

	// 当前计算完成后，重置开始时间节点为当前时间点，供下一次计算用
	in.TwapStartTimePoint = data.UnixTime
	// TWAP 计算end

	in.ProfitRate = profitRate // 填充收益率
	// 交易次数统计一下
	if data.LastQty > 0 { // 有子单回执成功算一次交易，被拒绝或撤销的不算
		// 收益率标准差
		in.PfRateList = append(in.PfRateList, in.ProfitRate)
		in.PfRateSum += in.ProfitRate
		pfAvg := in.PfRateSum / float64(in.DealCount)
		logx.Info("get in.PfRateList:", in.PfRateList, ", PfRateAvg:", pfAvg)
		pfVariance := CalVariance(in.PfRateList, pfAvg)
		in.PfRateStdDev = math.Sqrt(pfVariance / float64(in.DealCount))
		//}
		// 算出滑点值
		in.VwapDeal += data.LastPx * data.LastQty
		in.VWAP = float64(in.VwapDeal) / float64(in.LastQty) / 10000 // vwap值
		//in.VwapEntrust += data.Price * data.LastQty
		in.VwapEntrust += data.MarketPrice * data.LastQty
		// vwapSlippage
		in.VwapDev = in.VWAP - float64(in.VwapEntrust)/float64(in.LastQty)/10000 // 当前点的滑点值
		// 算滑点标准差
		in.VwapDevList = append(in.VwapDevList, in.VwapDev) // 保存每笔交易的滑点值
		in.VwapDevSum += in.VwapDev                         // 算滑点总和
		vdAvg := in.VwapDevSum / float64(in.DealCount)      // 滑点平均值
		logx.Info("get in.VwapDevList:", in.VwapDevList, ", avg:", vdAvg)
		vdVariance := CalVariance(in.VwapDevList, vdAvg)
		in.VwapStdDev = math.Sqrt(vdVariance / float64(in.DealCount)) // 标准差 = sqrt(平均方差)

		// 成交量贴合度标准差
		in.TradeVolFitList = append(in.TradeVolFitList, in.TradeVolFit)
		in.TradeVolFitSum += in.TradeVolFit
		tvfAvg := in.TradeVolFitSum / float64(in.DealCount)
		tvfVariance := CalVariance(in.TradeVolFitList, tvfAvg)
		in.TradeVolFitStdDev = math.Sqrt(tvfVariance / float64(in.DealCount))

		// 时间贴合度标准差
		in.TimeFitList = append(in.TimeFitList, c)
		in.TimeFitSum += c
		tfAvg := in.TimeFitSum / float64(in.DealCount)
		tfVariance := CalVariance(in.TimeFitList, tfAvg)
		in.TimeFitStdDev = math.Sqrt(tfVariance / float64(in.DealCount))

	}
	// 计算绩效因子
	if data.AlgorithmType == global.AlgoTypeSplit {
		factor := CalFactor(in.VwapStdDev, float64(data.MarketPrice), in.AvgCost, data.TradeSide)
		in.AssessFactor = factor * in.ProfitRate // 绩效收益率
	}

	if math.IsNaN(in.VwapDev) {
		in.VwapDev = 0.00
	}
	if math.IsNaN(in.VwapDevSum) {
		in.VwapDevSum = 0.00
	}
	if math.IsNaN(in.VwapStdDev) {
		in.VwapStdDev = 0.00
	}
	if math.IsNaN(in.PfRateStdDev) {
		in.PfRateStdDev = 0.00
	}
	if math.IsNaN(in.AvgCost) {
		in.AvgCost = 0.00
	}
	in.IndexCount++ // 增加计算次数
}

// CalculateProfileSum 计算算法画像
func CalculateProfileSum(head *global.ProfileHead, profileSumKey string, in *global.ProfileSum, data *global.ChildOrderData, usertype int) {
	in.ProfileHead = *head
	// 取画像汇总信息
	pSum := GetProfileSum(profileSumKey, usertype)

	in.OrderNum = pSum.Count       // 订单数量 （一个母单算一个订单)
	in.TradeVol = pSum.TradeVol    // 交易量
	in.LastQty = pSum.LastQty      // 总成交量
	in.TotalFee = pSum.Charge      // 总手续费
	in.CrossFee = pSum.CrossFee    // 总流量费
	in.TotalBuyVol = pSum.BuyVol   // 买入总量
	in.TotalSellVol = pSum.SellVol // 卖出总量

	in.EntrustQty = pSum.EntrustQty
	in.DealQty = pSum.DealQty
	in.CancelQty = pSum.CancelQty

	in.ProfitAmount = pSum.ProfitAmount // 盈亏金额
	in.TotalLastCost = pSum.CostAmount  // 总交易金额
	// 日内回转和智能委托收益率计算有所区别， 日内回转需要拿 （卖出+买入）/2的成本价除盈亏金额，智能委托直接拿交易成本除盈亏金额
	if data.AlgorithmType == global.AlgoTypeT0 {
		if pSum.CostT0Amount > 0 {
			in.ProfitRate = float64(in.ProfitAmount) / float64(pSum.CostT0Amount) // 收益率
		}
	} else if data.AlgorithmType == global.AlgoTypeSplit {
		if in.TotalLastCost > 0 {
			in.ProfitRate = float64(in.ProfitAmount) / float64(in.TotalLastCost) // 收益率
		}
	}
	// vwap 滑点需要区分证券ID进行计算
	// 上面单独计算完每支证券的滑点标准差后，算汇总滑点标准差时，求平均值
	if pSum.Count > 0 {
		in.VwapDev = pSum.VwapSum / float64(pSum.Count)
		in.VwapDevAvg = in.VwapDev
		in.VwapStdDev = pSum.VwapStdDevSum / float64(pSum.Count)
		in.PfRateStdDev = pSum.PfRateStdDevSum / float64(pSum.Count)
		in.AssessProfitRate = pSum.AssessFactor / float64(pSum.Count)

		in.TradeVolFitStdDev = pSum.TradeVolFitStdSum / float64(pSum.Count)
		in.TimeFitStdDev = pSum.TimeFitStdSum / float64(pSum.Count)

		in.DealEffi = pSum.DealEffi / float64(pSum.Count)
		in.AlgoOrderFit = pSum.AlgoOrderFit / float64(pSum.Count)
		in.TradeVolFit = pSum.TradeVolFit / float64(pSum.Count)
	}
	if data.AlgorithmType == global.AlgoTypeSplit { // 拆单算法
		in.WithdrawRate = pSum.Withdraw
		in.AvgCost = float64(pSum.CostAmount) / float64(pSum.Count)
	} else {
		in.AvgCost = float64(pSum.CostT0Amount) / float64(pSum.Count)
	}

	if in.EntrustQty > 0 {
		// 算撤单率
		in.CancelRate = float64(in.CancelQty) / float64(in.EntrustQty)
		// 算完成度
		in.Progress = (float64(in.DealQty) / float64(in.EntrustQty)) * 100
	}
	// 算最小拆单单位
	if in.MiniSplitOrder > data.OrderQty {
		in.MiniSplitOrder = data.OrderQty
	}
	// 算法的最小贴合度直接从母单的取最小值
	//c := GetMinJointRate(redis, in.Progress, data)
	if in.MinJointRate == 0.00 || in.MinJointRate > pSum.MinJointRate {
		in.MinJointRate = pSum.MinJointRate
	}

	// 统计盈亏为正的次数
	in.TradeCountPlus = pSum.TradeCountPlus

	if math.IsNaN(in.ProfitRate) {
		in.ProfitRate = 0.00
	}
	if math.IsNaN(in.VwapDevAvg) {
		in.VwapDevAvg = 0.00
	}
	if math.IsNaN(in.VwapVariance) {
		in.VwapVariance = 0.00
	}
	if math.IsNaN(in.VwapStdDev) {
		in.VwapStdDev = 0.00
	}
	if math.IsNaN(in.PfRateStdDev) {
		in.PfRateStdDev = 0.00
	}
	if math.IsNaN(in.AvgCost) {
		in.AvgCost = 0.00
	}
	// 计算绩效因子
	//in.AssessFactor = CalFactor(in.VwapStdDev, float64(data.ArrivePrice), in.AvgCost, data.TradeSide)

	//return hasProfile
}

// CalculateDynamic 计算动态汇总
func CalculateDynamic(in *global.ProfileSum, data *global.ChildOrderData) {
	// 算经济性
	in.EconomyScore = dao.GScoreConf.GetEconomyScore(in.TradeVol, in.ProfitAmount, in.TotalFee, in.MiniSplitOrder, in.ProfitRate, in.CancelRate, in.DealEffi)
	// 算完成度
	in.ProgressScore = dao.GScoreConf.GetProgressScore(in.Progress, in.AlgoOrderFit, in.TradeVolFit)
	// 算风险度
	in.RiskScore = dao.GScoreConf.GetRiskScore(in.MinJointRate, in.ProfitRate, in.WithdrawRate)
	// 算绩效评分
	in.AssessScore = dao.GScoreConf.GetAssessScore(in.VwapDev, in.AssessProfitRate)
	//ra := rand.New(rand.NewSource(time.Now().UnixNano()))
	//in.AssessScore = ra.Intn(8) + 1
	// 算稳定性
	in.StabilityScore = dao.GScoreConf.GetStabilityScore(in.VwapStdDev, in.ProfitRate, in.MinJointRate, in.TradeVolFitStdDev, in.TimeFitStdDev)

	// s.Logger.Info("get economy:", in.EconomyScore, ", progress:", in.ProgressScore, ",risk:", in.RiskScore,
	// 	", assess:", in.AssessScore, ", stability:", in.StabilityScore)

	// 总评分
	in.TotalScore = (in.EconomyScore + in.ProgressScore + in.RiskScore + in.AssessScore + in.StabilityScore) * 2
	//s.Logger.Info("get total sum score :", in.TotalScore)

	// 先算出当前交易的金额
	amount := data.LastPx * data.LastQty
	GetTradeVolDirect(amount, &in.TradeVolDict, data) // 买卖方向
	GetTradeVolRate(amount, &in.TradeVolVal)          // 交易量

	dao.GSecurityMap.RWMutex.RLock()
	fundType := dao.GSecurityMap.SecurityBase[data.SecId].FundType
	stockType := dao.GSecurityMap.SecurityBase[data.SecId].StockType
	dao.GSecurityMap.RWMutex.RUnlock()
	GetFundRate(fundType, amount, &in.FundPercent, data) // 资金占比
	GetStockType(stockType, &in.StockTypeVal, data)      // 股价类型

	in.IndexCount++
}

func CalculateTimeLine(head *global.ProfileHead, pk string, in *global.ProfileSum, data *global.ChildOrderData, userType int) {
	in.ProfileHead = *head
	// 取画像汇总信息
	pf := GetProfileSum(pk, userType)

	in.LastQty = pf.LastQty
	// 算完成度时间线
	if pf.EntrustQty > 0 {
		in.Progress = float64(pf.DealQty) / float64(pf.EntrustQty) * 100 // 真实完成度
	}
	// 这里的逻辑与计算profile逻辑需要保持一致
	// 不按分钟计算其盈亏数据，是因为每笔交易都可能会跨分钟成交，特别是日内回转其盈亏数据浮动会很大
	// 所以这里 直接取profile的计算数据
	// 求滑点的标准差

	in.ProfitAmount = pf.ProfitAmount // 盈亏金额
	in.TotalLastCost = pf.CostAmount  // 总交易金额
	if pf.Count > 0 {
		in.VwapDev = pf.VwapSum / float64(pf.Count)
		in.AssessProfitRate = pf.AssessFactor / float64(pf.Count)

		in.AlgoOrderFit = pf.AlgoOrderFit / float64(pf.Count)
		in.TradeVolFit = pf.TradeVolFit / float64(pf.Count)
	}
	// 日内回转和智能委托收益率计算有所区别， 日内回转需要拿 （卖出+买入）/2的成本价除盈亏金额，智能委托直接拿交易成本除盈亏金额
	if data.AlgorithmType == global.AlgoTypeT0 {
		if pf.CostT0Amount > 0 {
			in.ProfitRate = float64(in.ProfitAmount) / float64(pf.CostT0Amount) // 收益率
		}
	} else if data.AlgorithmType == global.AlgoTypeSplit {
		if in.TotalLastCost > 0 {
			in.ProfitRate = float64(in.ProfitAmount) / float64(in.TotalLastCost) // 收益率
		}
		in.WithdrawRate = pf.Withdraw // 拆单算回撤比例
	}
	in.ProgressScore = dao.GScoreConf.GetProgressScore(in.Progress, in.AlgoOrderFit, in.TradeVolFit) // 完成度评分
	// 算法绩效评分
	in.AssessScore = dao.GScoreConf.GetAssessScore(in.VwapDev, in.AssessProfitRate)
	//ra := rand.New(rand.NewSource(time.Now().UnixNano()))
	//in.AssessScore = ra.Intn(8) + 1

	// 算法的最小贴合度直接从母单的取最小值
	//c := GetMinJointRate(redis, in.Progress, data)
	if in.MinJointRate == 0.00 || in.MinJointRate > pf.MinJointRate {
		in.MinJointRate = pf.MinJointRate
	}
	// 风险评分
	in.RiskScore = dao.GScoreConf.GetRiskScore(in.MinJointRate, in.ProfitRate, in.WithdrawRate)

	in.IndexCount++
}

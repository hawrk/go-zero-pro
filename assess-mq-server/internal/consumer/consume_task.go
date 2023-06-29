// Package consumer
/*
 Author: hawrkchen
 Date: 2023/5/18 17:57
 Desc:
*/
package consumer

import (
	"algo_assess/global"
)

func GetProfileSum(profileSumKey string, userType int) ProfileResults {
	var r ProfileResults
	if userType == global.AccountTypeNormal {
		for k, v := range global.GNorUserProfile.ProfileMap {
			CalProfileSum(k, profileSumKey, v, &r)
		}
	} else if userType == global.AccountTypeProvider {
		for k, v := range global.GProviderProfile.ProfileMap {
			CalProfileSum(k, profileSumKey, v, &r)
		}
	} else if userType == global.AccountTypeMngr {
		for k, v := range global.GMngrProfile.ProfileMap {
			CalProfileSum(k, profileSumKey, v, &r)
		}
	} else if userType == global.AccountTypeSuAdmin {
		for k, v := range global.GAdminProfile.ProfileMap {
			CalProfileSum(k, profileSumKey, v, &r)
		}
	}
	return r
}

func CalProfileSum(profileKey, profileSumKey string, v *global.Profile, r *ProfileResults) {
	//logx.Info("get profileKey:", profileKey, "profilesumKey:", profilesumKey)
	if ProfileContains(profileKey, profileSumKey) {
		r.TradeVol += v.TotalTradeVol
		r.LastQty += v.LastQty
		r.Charge += v.TotalCharge
		r.CrossFee += v.TotalCrossFee
		r.BuyVol += v.TotalBuyVol
		r.SellVol += v.TotalSellVol
		r.EntrustQty += v.EntrustQty
		r.DealQty += v.DealQty
		r.CancelQty += v.CancelQty

		r.ProfitAmount += v.ProfitAmount
		r.CostAmount += v.TotalTradeVol // 总成本
		r.CostT0Amount += v.TotalT0Cost // T0 总成本
		if r.MinJointRate == 0.00 || r.MinJointRate > v.MinJointRate {
			r.MinJointRate = v.MinJointRate
		}
		// 绩效收益率
		r.AssessFactor += v.AssessFactor
		//滑点
		r.VwapSum += v.VwapDev
		// 滑点标准差
		r.VwapStdDevSum += v.VwapStdDev

		r.PfRateStdDevSum += v.PfRateStdDev
		r.TradeVolFitStdSum += v.TradeVolFitStdDev
		r.TimeFitStdSum += v.TimeFitStdDev

		r.DealEffi += v.DealEffi
		r.AlgoOrderFit += v.AlgoOrderFit
		r.TradeVolFit += v.TradeVolFit

		r.Count++
		if r.ProfitAmount > 0 {
			r.TradeCountPlus++
		}
		// 当前股票的回撤比例
		if r.Withdraw == 0.00 || r.Withdraw < v.WithdrawRate { // 取最大的回撤比例
			r.Withdraw = v.WithdrawRate
		}
	}
}

// GetProfileResults 取收益相关计算结果
func GetNorProfileResults(profileSumKey string) ProfileResults {
	var r ProfileResults
	// 在调用函数加锁了
	for k, v := range global.GNorUserProfile.ProfileMap {
		//只能拿该profile key前缀的数据
		// 比如该profile key (profilesumKey) 为 20221009:aUser0001:1 , profit key (k)为 20221009:aUser0001:1:000001
		//logx.Info("get k:", k, "profilesumKey:", profilesumKey)
		if ProfileContains(k, profileSumKey) {
			//logx.Info("into sum:", r.EntrustQty)
			//s.Logger.Info("add profit,k:", k, ",pk:", pk, ",profitAmount:", v.ProfitAmount,
			//	",totalTradeCost:", v.TotalTradeCost, ",totalT0Cost:", v.TotalT0Cost)
			r.TradeVol += v.TotalTradeVol
			r.LastQty += v.LastQty
			r.Charge += v.TotalCharge
			r.CrossFee += v.TotalCrossFee
			r.BuyVol += v.TotalBuyVol
			r.SellVol += v.TotalSellVol
			r.EntrustQty += v.EntrustQty
			r.DealQty += v.DealQty
			r.CancelQty += v.CancelQty

			r.ProfitAmount += v.ProfitAmount
			r.CostAmount += v.TotalTradeVol // 总成本
			r.CostT0Amount += v.TotalT0Cost // T0 总成本
			if r.MinJointRate == 0.00 || r.MinJointRate > v.MinJointRate {
				r.MinJointRate = v.MinJointRate
			}
			// 绩效收益率
			r.AssessFactor += v.AssessFactor
			//滑点
			r.VwapSum += v.VwapDev
			// 滑点标准差
			r.VwapStdDevSum += v.VwapStdDev

			r.PfRateStdDevSum += v.PfRateStdDev
			r.TradeVolFitStdSum += v.TradeVolFitStdDev
			r.TimeFitStdSum += v.TimeFitStdDev

			r.DealEffi += v.DealEffi
			r.AlgoOrderFit += v.AlgoOrderFit
			r.TradeVolFit += v.TradeVolFit

			r.Count++
			if r.ProfitAmount > 0 {
				r.TradeCountPlus++
			}
			// 当前股票的回撤比例
			if r.Withdraw == 0.00 || r.Withdraw < v.WithdrawRate { // 取最大的回撤比例
				r.Withdraw = v.WithdrawRate
			}
		}
	}
	return r
}

// GetProfileResults 取收益相关计算结果
func GetMngrProfileResults(profileSumKey string) ProfileResults {
	var r ProfileResults
	// 在调用函数加锁了
	for k, v := range global.GMngrProfile.ProfileMap {
		//只能拿该profile key前缀的数据
		// 比如该profile key (profilesumKey) 为 20221009:aUser0001:1 , profit key (k)为 20221009:aUser0001:1:000001
		//logx.Info("get k:", k, "profilesumKey:", profilesumKey)
		if ProfileContains(k, profileSumKey) {
			//logx.Info("into sum:", r.EntrustQty)
			//s.Logger.Info("add profit,k:", k, ",pk:", pk, ",profitAmount:", v.ProfitAmount,
			//	",totalTradeCost:", v.TotalTradeCost, ",totalT0Cost:", v.TotalT0Cost)
			r.TradeVol += v.TotalTradeVol
			r.LastQty += v.LastQty
			r.Charge += v.TotalCharge
			r.CrossFee += v.TotalCrossFee
			r.BuyVol += v.TotalBuyVol
			r.SellVol += v.TotalSellVol
			r.EntrustQty += v.EntrustQty
			r.DealQty += v.DealQty
			r.CancelQty += v.CancelQty

			r.ProfitAmount += v.ProfitAmount
			r.CostAmount += v.TotalTradeVol // 总成本
			r.CostT0Amount += v.TotalT0Cost // T0 总成本
			if r.MinJointRate == 0.00 || r.MinJointRate > v.MinJointRate {
				r.MinJointRate = v.MinJointRate
			}
			// 绩效收益率
			r.AssessFactor += v.AssessFactor
			//滑点
			r.VwapSum += v.VwapDev
			// 滑点标准差
			r.VwapStdDevSum += v.VwapStdDev

			r.PfRateStdDevSum += v.PfRateStdDev
			r.TradeVolFitStdSum += v.TradeVolFitStdDev
			r.TimeFitStdSum += v.TimeFitStdDev

			r.DealEffi += v.DealEffi
			r.AlgoOrderFit += v.AlgoOrderFit
			r.TradeVolFit += v.TradeVolFit

			r.Count++
			if r.ProfitAmount > 0 {
				r.TradeCountPlus++
			}
			// 当前股票的回撤比例
			if r.Withdraw == 0.00 || r.Withdraw < v.WithdrawRate { // 取最大的回撤比例
				r.Withdraw = v.WithdrawRate
			}
		}
	}
	return r
}

// GetProfileResults 取收益相关计算结果
func GetProviderProfileResults(profileSumKey string) ProfileResults {
	var r ProfileResults
	// 在调用函数加锁了
	for k, v := range global.GProviderProfile.ProfileMap {
		//只能拿该profile key前缀的数据
		// 比如该profile key (profilesumKey) 为 20221009:aUser0001:1 , profit key (k)为 20221009:aUser0001:1:000001
		//logx.Info("get k:", k, "profilesumKey:", profilesumKey)
		if ProfileContains(k, profileSumKey) {
			//logx.Info("into sum:", r.EntrustQty)
			//s.Logger.Info("add profit,k:", k, ",pk:", pk, ",profitAmount:", v.ProfitAmount,
			//	",totalTradeCost:", v.TotalTradeCost, ",totalT0Cost:", v.TotalT0Cost)
			r.TradeVol += v.TotalTradeVol
			r.LastQty += v.LastQty
			r.Charge += v.TotalCharge
			r.CrossFee += v.TotalCrossFee
			r.BuyVol += v.TotalBuyVol
			r.SellVol += v.TotalSellVol
			r.EntrustQty += v.EntrustQty
			r.DealQty += v.DealQty
			r.CancelQty += v.CancelQty

			r.ProfitAmount += v.ProfitAmount
			r.CostAmount += v.TotalTradeVol // 总成本
			r.CostT0Amount += v.TotalT0Cost // T0 总成本
			if r.MinJointRate == 0.00 || r.MinJointRate > v.MinJointRate {
				r.MinJointRate = v.MinJointRate
			}
			// 绩效收益率
			r.AssessFactor += v.AssessFactor
			//滑点
			r.VwapSum += v.VwapDev
			// 滑点标准差
			r.VwapStdDevSum += v.VwapStdDev

			r.PfRateStdDevSum += v.PfRateStdDev
			r.TradeVolFitStdSum += v.TradeVolFitStdDev
			r.TimeFitStdSum += v.TimeFitStdDev

			r.DealEffi += v.DealEffi
			r.AlgoOrderFit += v.AlgoOrderFit
			r.TradeVolFit += v.TradeVolFit

			r.Count++
			if r.ProfitAmount > 0 {
				r.TradeCountPlus++
			}
			// 当前股票的回撤比例
			if r.Withdraw == 0.00 || r.Withdraw < v.WithdrawRate { // 取最大的回撤比例
				r.Withdraw = v.WithdrawRate
			}
		}
	}
	return r
}

// GetProfileResults 取收益相关计算结果
func GetAdminProfileResults(profileSumKey string) ProfileResults {
	var r ProfileResults
	// 在调用函数加锁了
	for k, v := range global.GAdminProfile.ProfileMap {
		//只能拿该profile key前缀的数据
		// 比如该profile key (profilesumKey) 为 20221009:aUser0001:1 , profit key (k)为 20221009:aUser0001:1:000001
		//logx.Info("get k:", k, "profilesumKey:", profilesumKey)
		if ProfileContains(k, profileSumKey) {
			//logx.Info("into sum:", r.EntrustQty)
			//s.Logger.Info("add profit,k:", k, ",pk:", pk, ",profitAmount:", v.ProfitAmount,
			//	",totalTradeCost:", v.TotalTradeCost, ",totalT0Cost:", v.TotalT0Cost)
			r.TradeVol += v.TotalTradeVol
			r.LastQty += v.LastQty
			r.Charge += v.TotalCharge
			r.CrossFee += v.TotalCrossFee
			r.BuyVol += v.TotalBuyVol
			r.SellVol += v.TotalSellVol
			r.EntrustQty += v.EntrustQty
			r.DealQty += v.DealQty
			r.CancelQty += v.CancelQty

			r.ProfitAmount += v.ProfitAmount
			r.CostAmount += v.TotalTradeVol // 总成本
			r.CostT0Amount += v.TotalT0Cost // T0 总成本
			if r.MinJointRate == 0.00 || r.MinJointRate > v.MinJointRate {
				r.MinJointRate = v.MinJointRate
			}
			// 绩效收益率
			r.AssessFactor += v.AssessFactor
			//滑点
			r.VwapSum += v.VwapDev
			// 滑点标准差
			r.VwapStdDevSum += v.VwapStdDev

			r.PfRateStdDevSum += v.PfRateStdDev
			r.TradeVolFitStdSum += v.TradeVolFitStdDev
			r.TimeFitStdSum += v.TimeFitStdDev

			r.DealEffi += v.DealEffi
			r.AlgoOrderFit += v.AlgoOrderFit
			r.TradeVolFit += v.TradeVolFit

			r.Count++
			if r.ProfitAmount > 0 {
				r.TradeCountPlus++
			}
			// 当前股票的回撤比例
			if r.Withdraw == 0.00 || r.Withdraw < v.WithdrawRate { // 取最大的回撤比例
				r.Withdraw = v.WithdrawRate
			}
		}
	}
	return r
}

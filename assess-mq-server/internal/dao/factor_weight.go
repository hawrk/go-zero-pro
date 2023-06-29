// Package dao
/*
 Author: hawrkchen
 Date: 2022/8/9 10:23
 Desc: 计算每个影响评分的因子总分
*/
package dao

import (
	"github.com/spf13/cast"
)

// ----------------  经济性评分因子   ----------------------

// GetIntergerScore 计算整数类型评分
// 顺序，in 值越大，分值越高
func GetIntergerScore(in int64, rg []Region) int {
	if len(rg) == 0 { // 没有配置分数范围时，默认为5分
		return 5
	}
	var min, max float64 // 下限值和上限值
	min = rg[0].Lower
	max = rg[len(rg)-1].Upper
	if min >= max { // 最小值等于最大值时，都是0或是非法数据
		return 5
	}
	if in < cast.ToInt64(min) { // 低于下限值时，返回1
		return 1
	} else if in >= cast.ToInt64(max) { // 高于上限值时，返回10
		return 10
	}
	// 在配置区间内
	for _, v := range rg {
		if in >= cast.ToInt64(v.Lower) && in < cast.ToInt64(v.Upper) {
			return v.Score
		}
		//return 5 // 都无法匹配到时，默认返回5
	}
	return 5
}

// GetFloatScore 倒序排序，in 值越高，则分值越低
func GetFloatScore(in float64, rg []Region) int {
	if len(rg) == 0 { // 没有配置分数范围时，默认为5分
		return 5
	}
	var min, max float64 // 下限值和上限值
	min = rg[0].Lower
	max = rg[len(rg)-1].Upper
	if min >= max { // 最小值等于最大值时，都是0或是非法数据
		return 5
	}
	if in < min { // 低于下限值时，返回10
		return 10
	} else if in >= max { // 高于上限值时，返回1
		return 1
	}
	// 在配置区间内
	for _, v := range rg {
		if in >= v.Lower && in < v.Upper {
			return v.Score
		}
		//return 5 // 都无法匹配到时，默认返回5
	}
	return 5
}

// GetFloatScoreAsc 顺序排序，in 值越高，则分值越高
func GetFloatScoreAsc(in float64, rg []Region) int {
	if len(rg) == 0 { // 没有配置分数范围时，默认为5分
		return 5
	}
	var min, max float64 // 下限值和上限值
	min = rg[0].Lower
	max = rg[len(rg)-1].Upper
	if min >= max { // 最小值等于最大值时，都是0或是非法数据
		return 5
	}
	if in < min { // 低于下限值时，返回1
		return 1
	} else if in >= max { // 高于上限值时，返回10
		return 10
	}
	// 在配置区间内
	for _, v := range rg {
		if in >= v.Lower && in < v.Upper {
			return v.Score
		}
		//return 5 // 都无法匹配到时，默认返回5
	}
	return 5
}

// GetTradeVolScore 计算交易量评分
func GetTradeVolScore(vol int64) int {
	switch {
	case vol < 100000:
		return 1
	case vol >= 100000 && vol < 1000000:
		return 3
	case vol >= 1000000 && vol < 10000000:
		return 5
	case vol >= 10000000 && vol < 100000000:
		return 7
	case vol >= 100000000 && vol < 500000000:
		return 8
	case vol >= 500000000 && vol < 1000000000:
		return 9
	case vol >= 1000000000:
		return 10
	}
	return 0
}

// GetProfitScore  计算收益评分
func GetProfitScore(profit int64) int {
	switch {
	case profit < 1000:
		return 1
	case profit >= 1000 && profit < 5000:
		return 3
	case profit >= 5000 && profit < 10000:
		return 5
	case profit >= 10000 && profit < 50000:
		return 7
	case profit >= 50000 && profit < 100000:
		return 8
	case profit >= 100000 && profit < 500000:
		return 9
	case profit >= 500000:
		return 10
	}
	return 0
}

// GetProfitRateScore 收益率评分
func GetProfitRateScore(profitRate float64) int {
	switch {
	case profitRate < 0.01:
		return 1
	case profitRate >= 0.01 && profitRate < 0.1:
		return 3
	case profitRate >= 0.1 && profitRate < 1:
		return 5
	case profitRate >= 1 && profitRate < 5:
		return 7
	case profitRate >= 5 && profitRate < 10:
		return 8
	case profitRate >= 10 && profitRate < 20:
		return 9
	case profitRate >= 20:
		return 10
	}
	return 0
}

// GetTotalFeeScore 手续费评分
func GetTotalFeeScore(totalFee int64) int {
	switch {
	case totalFee < 10:
		return 10
	case totalFee >= 10 && totalFee < 100:
		return 9
	case totalFee >= 100 && totalFee < 1000:
		return 8
	case totalFee >= 1000 && totalFee < 5000:
		return 7
	case totalFee >= 5000 && totalFee < 10000:
		return 5
	case totalFee >= 10000 && totalFee < 50000:
		return 3
	case totalFee >= 50000:
		return 1
	}
	return 0
}

// GetCrossFeeScore 流量费评分
func GetCrossFeeScore(crossFee int64) int {
	return 10
}

// GetMinSplitRateScore 最小拆单单位评分
func GetMinSplitRateScore(minSplitRate int64) int {
	switch {
	case minSplitRate < 10:
		return 10
	case minSplitRate >= 10 && minSplitRate < 100:
		return 9
	case minSplitRate >= 100 && minSplitRate < 1000:
		return 8
	case minSplitRate >= 1000 && minSplitRate < 5000:
		return 7
	case minSplitRate >= 5000 && minSplitRate < 10000:
		return 5
	case minSplitRate >= 10000 && minSplitRate < 50000:
		return 3
	case minSplitRate >= 50000:
		return 1
	}
	return 0
}

// GetCancelRateScore 撤单率评分
func GetCancelRateScore(cancelRate float64) int {
	switch {
	case cancelRate < 0.01:
		return 10
	case cancelRate >= 0.01 && cancelRate < 0.1:
		return 9
	case cancelRate >= 0.1 && cancelRate < 1:
		return 8
	case cancelRate >= 1 && cancelRate < 5:
		return 7
	case cancelRate >= 5 && cancelRate < 10:
		return 5
	case cancelRate >= 10 && cancelRate < 50:
		return 3
	case cancelRate >= 50:
		return 1
	}
	return 0
}

// ----------------- 完成度评分因子  ----------------------

// GetProScore 完成度评分
func GetProScore(progress float64) int {
	switch {
	case progress <= 0:
		return 0
	case progress > 0 && progress < 20:
		return 1
	case progress >= 20 && progress < 30:
		return 2
	case progress >= 30 && progress < 40:
		return 3
	case progress >= 40 && progress < 50:
		return 4
	case progress >= 50 && progress < 60:
		return 5
	case progress >= 60 && progress < 70:
		return 6
	case progress >= 70 && progress < 80:
		return 7
	case progress >= 80 && progress < 90:
		return 8
	case progress >= 90 && progress < 100:
		return 9
	case progress >= 100:
		return 10
	}
	return 0
}

// ----------------- 风险度评分因子  ----------------------

// GetMinJointRateScore 最小贴合度评分
func GetMinJointRateScore(minJointRate float64) int {
	switch {
	case minJointRate < 0.01:
		return 10
	case minJointRate >= 0.01 && minJointRate < 0.1:
		return 9
	case minJointRate >= 0.1 && minJointRate < 1:
		return 8
	case minJointRate >= 1 && minJointRate < 5:
		return 7
	case minJointRate >= 5 && minJointRate < 10:
		return 5
	case minJointRate >= 10 && minJointRate < 50:
		return 3
	case minJointRate >= 50:
		return 1
	}
	return 0
}

// GetWithdrawRateScore 回撤比例评分
func GetWithdrawRateScore(withdrawRate float64) int {
	switch {
	case withdrawRate < 0.01:
		return 10
	case withdrawRate >= 0.01 && withdrawRate < 0.1:
		return 9
	case withdrawRate >= 0.1 && withdrawRate < 1:
		return 8
	case withdrawRate >= 1 && withdrawRate < 5:
		return 7
	case withdrawRate >= 5 && withdrawRate < 10:
		return 5
	case withdrawRate >= 10 && withdrawRate < 50:
		return 3
	case withdrawRate >= 50:
		return 1
	}
	return 0
}

// GetRProfitRateScore 收益率评分
func GetRProfitRateScore(profitRate float64) int {
	switch {
	case profitRate < 0.01:
		return 1
	case profitRate >= 0.01 && profitRate < 0.1:
		return 3
	case profitRate >= 0.1 && profitRate < 1:
		return 5
	case profitRate >= 1 && profitRate < 5:
		return 7
	case profitRate >= 5 && profitRate < 10:
		return 8
	case profitRate >= 10 && profitRate < 20:
		return 9
	case profitRate >= 20:
		return 10
	}
	return 0
}

// ----------------- 绩效评分因子  ----------------------

// GetVwapDevScore vwap滑点值评分
func GetVwapDevScore(vwapDev float64) int {
	switch {
	case vwapDev < 0.01:
		return 10
	case vwapDev >= 0.01 && vwapDev < 0.1:
		return 9
	case vwapDev >= 0.1 && vwapDev < 1:
		return 8
	case vwapDev >= 1 && vwapDev < 5:
		return 7
	case vwapDev >= 5 && vwapDev < 10:
		return 5
	case vwapDev >= 10 && vwapDev < 20:
		return 3
	case vwapDev >= 20:
		return 1
	}
	return 0
}

// GetAProfitRateScore 收益率评分
func GetAProfitRateScore(profitRate float64) int {
	switch {
	case profitRate < 0.01:
		return 1
	case profitRate >= 0.01 && profitRate < 0.1:
		return 3
	case profitRate >= 0.1 && profitRate < 1:
		return 5
	case profitRate >= 1 && profitRate < 5:
		return 7
	case profitRate >= 5 && profitRate < 10:
		return 8
	case profitRate >= 10 && profitRate < 20:
		return 9
	case profitRate >= 20:
		return 10
	}
	return 0
}

// ----------------- 稳定性评分因子  ----------------------

// GetVwapDevStdScore vwap 滑点标准差评分
func GetVwapDevStdScore(vwapDevStd float64) int {
	switch {
	case vwapDevStd < 0.01:
		return 10
	case vwapDevStd >= 0.01 && vwapDevStd < 0.1:
		return 9
	case vwapDevStd >= 0.1 && vwapDevStd < 1:
		return 8
	case vwapDevStd >= 1 && vwapDevStd < 5:
		return 7
	case vwapDevStd >= 5 && vwapDevStd < 10:
		return 5
	case vwapDevStd >= 10 && vwapDevStd < 50:
		return 3
	case vwapDevStd >= 50:
		return 1
	}
	return 0
}

// GetProfitRateStdScore 收益率标准差评分
func GetProfitRateStdScore(profitRateStd float64) int {
	switch {
	case profitRateStd < 0.01:
		return 10
	case profitRateStd >= 0.01 && profitRateStd < 0.1:
		return 9
	case profitRateStd >= 0.1 && profitRateStd < 1:
		return 8
	case profitRateStd >= 1 && profitRateStd < 5:
		return 7
	case profitRateStd >= 5 && profitRateStd < 10:
		return 5
	case profitRateStd >= 10 && profitRateStd < 50:
		return 3
	case profitRateStd >= 50:
		return 1
	}
	return 0
}

// GetJointRateScore 贴合度评分
func GetJointRateScore(jointRate float64) int {
	switch {
	case jointRate < 0.01:
		return 10
	case jointRate >= 0.01 && jointRate < 0.1:
		return 9
	case jointRate >= 0.1 && jointRate < 1:
		return 8
	case jointRate >= 1 && jointRate < 5:
		return 7
	case jointRate >= 5 && jointRate < 10:
		return 5
	case jointRate >= 10 && jointRate < 50:
		return 3
	case jointRate >= 50:
		return 1
	}
	return 0
}

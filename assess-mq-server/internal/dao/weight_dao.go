// Package dao
/*
 Author: hawrkchen
 Date: 2022/7/15 17:03
 Desc:
*/
package dao

import (
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"math"
)

// 算法评分计算，每个交易因子的计算权重
const (
	TradeVolWeight   = iota + 1 //交易量
	ProfitWeight                // 盈亏收益
	ProfitRateWeight            //收益率
	ChargeWeight                // 手续费
	//CrossFeeWeight              // 流量费
	MinSplitWeight   // 最小拆单单位
	CancelRateWeight // 撤单率
	DealEffiWeight   // 成交效率

	ProgressWeight     // 完成度
	AlgoOrderFitWeight // 母单贴合度
	TradeVolFitWeight  // 成交量贴合度

	MinJoionWeight     // 最小贴合度
	WithDrawRateWeight // 回撤比例
	RProfitRateWeight  // 收益率

	VwapDevWeight     // vwap滑点值
	AProfitRateWeight // 收益率

	VwapDevStdWeight     // vwap滑点标准差
	ProfitRateStdWeight  // 收益率标准差
	JoinRateWeight       // 贴合度
	TradeVolFitStdWeight // 成交量贴合度标准差
	TimeFitStdWeight     // 时间贴合度标准差
)

//var GWeights Weights
var GScoreConf ScoreConf

// Weights 适配Json反序列化
type Weights struct {
	Ew []AlgoFactor // 经济性
	Pw []AlgoFactor // 完成度
	Rw []AlgoFactor // 风险度
	Aw []AlgoFactor // 绩效
	Sw []AlgoFactor // 稳定性
}

type Region struct { // 区间配置
	Lower float64 `json:"lower"`
	Upper float64 `json:"upper"`
	Score int     `json:"score"`
}

// AlgoFactor  算法因子配置
type AlgoFactor struct {
	FactorKey  string   `json:"factor_key"`
	FactorName string   `json:"factor_name"`
	Weight     int      `json:"weight"`
	Regions    []Region `json:"regions"`
}

type FactorConf struct {
	Weight  int      `json:"weight"`
	Regions []Region `json:"regions"`
}

type ScoreConf struct {
	Ew EconomyConf
	Pw ProgressConf
	Rw RiskConf
	Aw AssessConf
	Sw StableConf
}

// EconomyConf  经济性算法配置
type EconomyConf struct {
	TradeVol   FactorConf `json:"trade_vol"`
	Profit     FactorConf `json:"profit"`
	ProfitRate FactorConf `json:"profit_rate"`
	Charge     FactorConf `json:"charge"`
	//CrossFee   FactorConf `json:"cross_fee"`
	MinSplit   FactorConf `json:"min_split"`
	CancelRate FactorConf `json:"cancel_rate"`
	DealEffi   FactorConf `json:"deal_effi"` // 成交效率
}

// ProgressConf 完成度算法配置
type ProgressConf struct {
	Progress     FactorConf `json:"progress"`
	AlgoOrderFit FactorConf `json:"algo_order_fit"` // 母单贴合度
	TradeVolFit  FactorConf `json:"trade_vol_fit"`  // 成交量贴合度
}

// RiskConf 风险度算法配置
type RiskConf struct {
	MinJoion     FactorConf `json:"min_joion"`
	WithDrawRate FactorConf `json:"with_draw_rate"`
	ProfitRate   FactorConf `json:"profit_rate"`
}

// AssessConf 绩效算法配置
type AssessConf struct {
	VwapDev          FactorConf `json:"vwap_dev"`
	AssessProfitRate FactorConf `json:"assess_profitrate"`
}

// StableConf 稳定性算法配置
type StableConf struct {
	VwapDevStd    FactorConf `json:"vwap_dev_std"`
	ProfitRateStd FactorConf `json:"profit_rate_std"`
	JoinRate      FactorConf `json:"join_rate"`
	TradeVolStd   FactorConf `json:"trade_vol_std_dev"` // 成交量贴合标准差
	TimeStdDev    FactorConf `json:"time_std_dev"`      // 时间贴合度标准差

}

func NewWeights(m map[int]string) ScoreConf {
	var w Weights
	var out ScoreConf
	for k, v := range m {
		if k == 1 { // 经济性
			if err := json.Unmarshal([]byte(v), &w.Ew); err != nil {
				logx.Error("Unmarshal EconomyWeight param error:", err)
			}
			ParseBusiConfig(w.Ew, &out)
		} else if k == 2 { // 完成度
			if err := json.Unmarshal([]byte(v), &w.Pw); err != nil {
				logx.Error("Unmarshal ProgressWeight param error:", err)
			}
			ParseBusiConfig(w.Pw, &out)
		} else if k == 3 { // 风险度
			if err := json.Unmarshal([]byte(v), &w.Rw); err != nil {
				logx.Error("Unmarshal RiskWeight param error:", err)
			}
			ParseBusiConfig(w.Rw, &out)
		} else if k == 4 { // 绩效
			if err := json.Unmarshal([]byte(v), &w.Aw); err != nil {
				logx.Error("Unmarshal AssessWeight param error:", err)
			}
			ParseBusiConfig(w.Aw, &out)
		} else if k == 5 { // 稳定性
			if err := json.Unmarshal([]byte(v), &w.Sw); err != nil {
				logx.Error("Unmarshal StabilityWeight param error:", err)
			}
			ParseBusiConfig(w.Sw, &out)
		}
	}
	logx.Infof("get GScoreConf:%+v", out)
	//logx.Infof("get Ew:%+v", w.Ew)
	//logx.Infof("get Pw:%+v", w.Pw)
	//logx.Infof("get Rw:%+v", w.Rw)
	//logx.Infof("get Aw:%+v", w.Aw)
	//logx.Infof("get Sw:%+v", w.Sw)
	return out
}

// GetEconomyScore 计算经济性评分
// 经济性评分规则: 根据 交易量 收益 总手续费 流量费 最小拆单单位 收益率 回撤比例 七个维度进行计算
// 每个维度满分 10分，乘以每个维度的加权系数（加权系数在数据库中可配置），得到每个维度的最终分数 a1
// 所有维度相加的总分即为该经济性的最终评分

func (w *ScoreConf) GetEconomyScore(vol int64, profit float64, totalFee, MinSplitRate int64, profitRate, CancelRate, dealEffi float64) int {
	score := GetIntergerScore(vol, w.Ew.TradeVol.Regions)*GetWeight(w.Ew.TradeVol.Weight, TradeVolWeight) +
		GetFloatScoreAsc(profit, w.Ew.Profit.Regions)*GetWeight(w.Ew.Profit.Weight, ProfitWeight) +
		GetFloatScoreAsc(profitRate, w.Ew.ProfitRate.Regions)*GetWeight(w.Ew.ProfitRate.Weight, ProfitRateWeight) +
		GetIntergerScore(totalFee, w.Ew.Charge.Regions)*GetWeight(w.Ew.Charge.Weight, ChargeWeight) +
		//GetIntergerScore(crossFee, w.Ew.CrossFee.Regions)*GetWeight(w.Ew.CrossFee.Weight, CrossFeeWeight) +
		GetIntergerScore(MinSplitRate, w.Ew.MinSplit.Regions)*GetWeight(w.Ew.MinSplit.Weight, MinSplitWeight) +
		GetFloatScore(CancelRate, w.Ew.CancelRate.Regions)*GetWeight(w.Ew.CancelRate.Weight, CancelRateWeight) +
		GetFloatScoreAsc(dealEffi, w.Ew.DealEffi.Regions)*GetWeight(w.Ew.DealEffi.Weight, DealEffiWeight)
	return ScoreRound(score)
}

// GetProgressScore 计算完成度评分
func (w *ScoreConf) GetProgressScore(progress, algoOrderFit, tradeVolFit float64) int {
	//logx.Info("get weight:", w.Pw.Progress.Weight)
	//logx.Infof("get regin:%+v", w.Pw.Progress.Regions)
	score := GetFloatScoreAsc(progress, w.Pw.Progress.Regions)*GetWeight(w.Pw.Progress.Weight, ProgressWeight) +
		GetFloatScoreAsc(algoOrderFit, w.Pw.AlgoOrderFit.Regions)*GetWeight(w.Pw.AlgoOrderFit.Weight, AlgoOrderFitWeight) +
		GetFloatScoreAsc(tradeVolFit, w.Pw.TradeVolFit.Regions)*GetWeight(w.Pw.TradeVolFit.Weight, TradeVolFitWeight)
	return ScoreRound(score)

}

// GetRiskScore 计算风险度评分
func (w *ScoreConf) GetRiskScore(minJointRate, ProfitRate, withdrawRate float64) int {
	score := GetFloatScore(minJointRate, w.Rw.MinJoion.Regions)*GetWeight(w.Rw.MinJoion.Weight, MinJoionWeight) +
		GetFloatScoreAsc(ProfitRate, w.Rw.ProfitRate.Regions)*GetWeight(w.Rw.ProfitRate.Weight, WithDrawRateWeight) +
		GetFloatScore(withdrawRate, w.Rw.WithDrawRate.Regions)*GetWeight(w.Rw.WithDrawRate.Weight, RProfitRateWeight)
	return ScoreRound(score)

}

// GetAssessScore 计算算法绩效评分
func (w *ScoreConf) GetAssessScore(vwapDev, profitRate float64) int {
	score := GetFloatScore(vwapDev, w.Aw.VwapDev.Regions)*GetWeight(w.Aw.VwapDev.Weight, VwapDevWeight) +
		GetFloatScoreAsc(profitRate, w.Aw.AssessProfitRate.Regions)*GetWeight(w.Aw.AssessProfitRate.Weight, AProfitRateWeight)
	return ScoreRound(score)

}

// GetStabilityScore 计算稳定性评分
func (w *ScoreConf) GetStabilityScore(vwapStd, profitRateStd, jonitRate, tvfStd, timeFitStd float64) int {
	score := GetFloatScore(vwapStd, w.Sw.VwapDevStd.Regions)*GetWeight(w.Sw.VwapDevStd.Weight, VwapDevStdWeight) +
		GetFloatScore(profitRateStd, w.Sw.ProfitRateStd.Regions)*GetWeight(w.Sw.ProfitRateStd.Weight, ProfitRateStdWeight) +
		GetFloatScore(jonitRate, w.Sw.JoinRate.Regions)*GetWeight(w.Sw.JoinRate.Weight, JoinRateWeight) +
		GetFloatScoreAsc(tvfStd, w.Sw.TradeVolStd.Regions)*GetWeight(w.Sw.TradeVolStd.Weight, TradeVolWeight) +
		GetFloatScoreAsc(timeFitStd, w.Sw.TimeStdDev.Regions)*GetWeight(w.Sw.TimeStdDev.Weight, TimeFitStdWeight)
	return ScoreRound(score)
}

// ScoreRound 四舍五入
func ScoreRound(score int) int {
	return int(math.Round(float64(score) / 100))
}

//GetWeight 算法配置为空时，给默认值
func GetWeight(score int, t int) int {
	if score != 0 {
		return score
	}
	switch t {
	case TradeVolWeight:
		return 10
	case ProfitWeight:
		return 10
	case ProfitRateWeight:
		return 30
	case ChargeWeight:
		return 10
	//case CrossFeeWeight:
	//	return 5
	case MinSplitWeight:
		return 20
	case CancelRateWeight:
		return 10
	case DealEffiWeight:
		return 10

	case ProgressWeight:
		return 80
	case AlgoOrderFitWeight:
		return 10
	case TradeVolFitWeight:
		return 10

	case MinJoionWeight:
		return 35
	case WithDrawRateWeight:
		return 25
	case RProfitRateWeight:
		return 40

	case VwapDevWeight:
		return 40
	case AProfitRateWeight:
		return 60

	case VwapDevStdWeight:
		return 30
	case ProfitRateStdWeight:
		return 30
	case JoinRateWeight:
		return 10
	case TradeVolFitStdWeight:
		return 20
	case TimeFitStdWeight:
		return 10
	default:
		return 100
	}
}

func ParseBusiConfig(in []AlgoFactor, o *ScoreConf) {
	for _, v := range in {
		switch v.FactorKey {
		case "trade_vol":
			o.Ew.TradeVol.Weight = v.Weight
			o.Ew.TradeVol.Regions = v.Regions
		case "profit":
			o.Ew.Profit.Weight = v.Weight
			o.Ew.Profit.Regions = v.Regions
		case "profit_rate":
			o.Ew.ProfitRate.Weight = v.Weight
			o.Ew.ProfitRate.Regions = v.Regions
			//风险
			o.Rw.ProfitRate.Weight = v.Weight
			o.Rw.ProfitRate.Regions = v.Regions

		case "charge":
			o.Ew.Charge.Weight = v.Weight
			o.Ew.Charge.Regions = v.Regions
		//case "cross_fee":
		//	o.Ew.CrossFee.Weight = v.Weight
		//	o.Ew.CrossFee.Regions = v.Regions
		case "cancel_rate":
			o.Ew.CancelRate.Weight = v.Weight
			o.Ew.CancelRate.Regions = v.Regions
		case "min_split_order":
			o.Ew.MinSplit.Weight = v.Weight
			o.Ew.MinSplit.Regions = v.Regions
		case "deal_effi":
			o.Ew.DealEffi.Weight = v.Weight
			o.Ew.DealEffi.Regions = v.Regions
			// 完成度
		case "progress":
			o.Pw.Progress.Weight = v.Weight
			o.Pw.Progress.Regions = v.Regions
		case "algo_order_fit":
			o.Pw.AlgoOrderFit.Weight = v.Weight
			o.Pw.AlgoOrderFit.Regions = v.Regions
		case "trade_vol_fit":
			o.Pw.TradeVolFit.Weight = v.Weight
			o.Pw.TradeVolFit.Regions = v.Regions
			// 风险度
		case "min_jonit_rate":
			o.Rw.MinJoion.Weight = v.Weight
			o.Rw.MinJoion.Regions = v.Regions
		case "withdraw_rate":
			o.Rw.WithDrawRate.Weight = v.Weight
			o.Rw.WithDrawRate.Regions = v.Regions
			// 绩效
		case "vwap_dev":
			o.Aw.VwapDev.Weight = v.Weight
			o.Aw.VwapDev.Regions = v.Regions
		case "assess_profitrate":
			o.Aw.AssessProfitRate.Weight = v.Weight
			o.Aw.AssessProfitRate.Regions = v.Regions
			// 稳定
		case "vwap_std_dev":
			o.Sw.VwapDevStd.Weight = v.Weight
			o.Sw.VwapDevStd.Regions = v.Regions
		case "profit_rate_std":
			o.Sw.ProfitRateStd.Weight = v.Weight
			o.Sw.ProfitRateStd.Regions = v.Regions
		case "joint_rate":
			o.Sw.JoinRate.Weight = v.Weight
			o.Sw.JoinRate.Regions = v.Regions
		case "trade_vol_std_dev":
			o.Sw.TradeVolStd.Weight = v.Weight
			o.Sw.TradeVolStd.Regions = v.Regions
		case "time_std_dev":
			o.Sw.TimeStdDev.Weight = v.Weight
			o.Sw.TimeStdDev.Regions = v.Regions
		}

	}
}

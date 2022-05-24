// Package logic
/*
 Author: hawrkchen
 Date: 2022/4/25 14:52
 Desc:
*/
package logic

import (
	"algo_assess/assess-api-server/internal/types"
	mqservice "algo_assess/assess-mq-server/assessmqservice"
	"algo_assess/assess-rpc-server/assessservice"
	mkservice "algo_assess/market-mq-server/marketservice"
	"algo_assess/pkg/tools"
	"github.com/spf13/cast"
)

func GetTimePoint(s string) string {
	if len(s) < 12 {
		return ""
	}
	return s[8:10] + ":" + s[10:]
}

// MakeMarketInfo 拼组行情信息
func MakeMarketInfo(t int64, qRsp *mkservice.MarkRsp) (lastPrice, mVwap float64, tradeVol int64, askPrice, askVol, bidPrice, bidVol string) {
	if _, exist := qRsp.Attrs[t]; exist {
		lastPrice = tools.Fen2Yuan(qRsp.Attrs[t].GetLastPrice())
		tradeVol = qRsp.Attrs[t].GetTradeVol()
		askPrice = qRsp.Attrs[t].GetAskPrice()
		askVol = qRsp.Attrs[t].GetAskVol()
		bidPrice = qRsp.Attrs[t].GetBidPrice()
		bidVol = qRsp.Attrs[t].GetBidVol()
		mVwap = qRsp.Attrs[t].GetMarketVwap()
	}
	return
}

// MakeAssessMqRsp 拼组实时缓存绩效信息
func MakeAssessMqRsp(m map[string]types.GeneralData, mrsp *mqservice.GeneralRsp, qRsp *mkservice.MarkRsp) {
	for _, v := range mrsp.Info {
		timePoint := GetTimePoint(cast.ToString(v.TransactTime))
		lp, mvwap, tv, ap, av, bp, bv := MakeMarketInfo(v.TransactTime, qRsp)
		detail := types.GeneralData{
			TransTime:    timePoint,
			OrderQty:     v.OrderQty,
			LastQty:      v.LastQty,
			CancelQty:    v.CancelledQty,
			RejectQty:    v.RejectedQty,
			Vwap:         mvwap, // 取市场vwap
			VwapDev:      v.VwapDeviation,
			LastPrice:    tools.Fen2Yuan(v.LastPrice),
			ArriPrice:    tools.Fen2Yuan(v.ArrivedPrice),
			ArriPriceDev: v.ArrivedPriceDeviation,
			MarketRate:   v.MarketRate,
			DealRate:     v.DealRate,
			DealProgress: v.DealProgress,
			MkLastPrice:  lp,
			MkTradeVol:   tv,
			AskPrice:     ap,
			AskVol:       av,
			BidPrice:     bp,
			BidVol:       bv,
		}
		m[timePoint] = detail
	}
}

// MakeAssessRPCRsp 拼组落地DB绩效信息
func MakeAssessRPCRsp(m map[string]types.GeneralData, grsp *assessservice.GeneralRsp, qRsp *mkservice.MarkRsp) {
	for _, value := range grsp.Info {
		timePoint := GetTimePoint(cast.ToString(value.TransactTime))
		lp, mvwap, tv, ap, av, bp, bv := MakeMarketInfo(value.TransactTime, qRsp)
		detail := types.GeneralData{
			TransTime:    timePoint,
			OrderQty:     value.OrderQty,
			LastQty:      value.LastQty,
			CancelQty:    value.CancelledQty,
			RejectQty:    value.RejectedQty,
			Vwap:         mvwap,
			VwapDev:      value.VwapDeviation,
			LastPrice:    tools.Fen2Yuan(value.LastPrice),
			ArriPrice:    tools.Fen2Yuan(value.ArrivedPrice),
			ArriPriceDev: value.ArrivedPriceDeviation,
			MarketRate:   value.MarketRate,
			DealRate:     value.DealRate,
			DealProgress: value.DealProgress,
			MkLastPrice:  lp,
			MkTradeVol:   tv,
			AskPrice:     ap,
			AskVol:       av,
			BidPrice:     bp,
			BidVol:       bv,
		}
		m[timePoint] = detail
	}
}

func MakeEmptyDataRsp(start string, qRsp *mkservice.MarkRsp, progress float64) types.GeneralData {
	lp, mvwap, tv, ap, av, bp, bv := MakeMarketInfo(cast.ToInt64(start), qRsp)
	out := types.GeneralData{
		TransTime:    GetTimePoint(start),
		OrderQty:     0,
		LastQty:      0,
		CancelQty:    0,
		RejectQty:    0,
		Vwap:         mvwap,
		VwapDev:      0,
		LastPrice:    0,
		ArriPrice:    0,
		ArriPriceDev: 0,
		MarketRate:   0,
		DealRate:     0,
		DealProgress: progress,
		MkLastPrice:  lp,
		MkTradeVol:   tv,
		AskPrice:     ap,
		AskVol:       av,
		BidPrice:     bp,
		BidVol:       bv,
	}
	return out
}

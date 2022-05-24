// Package logic
/*
 Author: hawrkchen
 Date: 2022/4/14 14:38
 Desc:
*/
package logic

import (
	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/global"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strconv"
	"time"
)

// RealTimeCal 掌握核心技术
func RealTimeCal(svcCtx *svc.ServiceContext, in *global.OrderAssess, data *global.ChildOrderData) (*global.OrderAssess, error) {
	if in == nil {
		in = new(global.OrderAssess)
	}
	logx.Infof("begin cal, in:%+v", in)
	in.AlgorithmType = data.AlgorithmType
	in.AlgorithmId = data.AlgoId
	in.UsecurityId = data.UsecId
	in.SecurityId = data.SecId
	in.CreateTime = time.Now()
	in.TimeDimension = 1
	in.ArrivedPrice = data.ArrivePrice
	in.TransactAt = data.TransTime

	in.LastQty += data.LastQty

	in.SubVwap += data.LastPx * data.LastQty             // 成交价格*成交笔数
	in.SubVwapEntrust += data.Price * data.LastQty       // 委托价格*成交笔数
	in.SubVwapArrived += data.ArrivePrice * data.LastQty // 到达价格*成交笔数
	// 该子单算法证券状态变更时会推送多次，委托数量只需要统计一次，但是该子单从下单开始到成交，可能会跨越多个时间段，所以要加上当前交易时间
	orderKey := fmt.Sprintf("%d:%d", data.TransTime, data.OrderId)
	logx.Info("get orderKey:", orderKey)
	if _, orderKeyExist := global.GlobalAssess.OrderMap[orderKey]; !orderKeyExist {
		in.OrderQty += data.OrderQty
		global.GlobalAssess.OrderMap[orderKey] = struct{}{}
	}
	if data.ChildOrderStatus == global.OrderStatusCancel {
		in.CancelQty += data.OrderQty - data.ComQty
	}
	if data.ChildOrderStatus == global.OrderStatusApReject || data.ChildOrderStatus == global.OrderStatusCtReject ||
		data.ChildOrderStatus == global.OrderStatusTaReject {
		in.RejectedQty += data.OrderQty
	}
	// 计算成交量比重  市场参与率
	quoteKey := fmt.Sprintf("%s%s:%d", global.Level2RedisPrx, data.SecId, data.TransTime)
	logx.Info("get quoteKey:", quoteKey)
	in.LastPrice, in.MarketRate, in.DealRate = GetMarketRate(svcCtx.RedisClient, quoteKey, in.LastQty)
	// 取市场vwap   ---已经落表，不需要在这里取了
	//in.Vwap = GetMarketVwap(quoteKey)
	// 计算成交进度
	algoKey := fmt.Sprintf("%d:%s", data.AlgoId, data.SecId)
	in.DealProgress = GetDealProgress(algoKey, in.LastQty)

	if in.LastQty > 0 {
		in.TradeVwap = float64(in.SubVwap) / float64(in.LastQty)                                         // 计算vwap
		in.VwapDeviation = (in.TradeVwap - float64(in.SubVwapEntrust)/float64(in.LastQty)) / 100         // 计算vwap 滑点(计算是以分为单位，需除100
		in.ArrivedPriceDeviation = (in.TradeVwap - float64(in.SubVwapArrived)/float64(in.LastQty)) / 100 // 计算到达价滑点
	}
	logx.Infof("after real cal, out:%+v", in)
	return in, nil
}

// GetDealProgress 成交进度
func GetDealProgress(algoKey string, lastQty int64) (progress float64) {
	global.GlobalAlgoOrder.RWMutex.RLock()
	defer global.GlobalAlgoOrder.RWMutex.RUnlock()
	// 这里不能根据母单ID作为key是因为每支算法和证券ID可能有多个母单
	if global.GlobalAlgoOrder.AlgoOrder[algoKey] > 0 { // 母单成交数量/ 母单委托数量
		global.GlobalAlgoOrder.DealAlgoOrder[algoKey] += lastQty
		progress = (float64(global.GlobalAlgoOrder.DealAlgoOrder[algoKey]) / float64(global.GlobalAlgoOrder.AlgoOrder[algoKey])) * 100
	}
	return progress
}

// GetMarketVwap 市场vwap
func GetMarketVwap(quoteKey string) (vwap float64) {
	global.GlobalMarketVwap.RWMutex.RLock()
	defer global.GlobalMarketVwap.RWMutex.RUnlock()

	vwap = global.GlobalMarketVwap.MVwap[quoteKey]
	return vwap
}

// GetMarketRate 市场参与率， 成交量比重
func GetMarketRate(redis *redis.Redis, quoteKey string, lastQty int64) (marketLastPrice int64, MarketRate, DealRate float64) {
	if true { // 读redis的数据
		out, _ := redis.Hmget(quoteKey, "entrustvol", "tradevol", "lastprice")
		if len(out) < 3 {
			return
		}
		marketLastPrice = cast.ToInt64(out[2])
		marketQty := cast.ToInt64(out[0])
		if marketQty > 0 {
			MarketRate = (float64(lastQty) / float64(marketQty)) * 100
		}
		totalTradeVol := cast.ToInt64(out[1])
		if totalTradeVol > 0 {
			DealRate = (float64(lastQty) / float64(totalTradeVol)) * 100 // 成交量比重
		}
		return marketLastPrice, MarketRate, DealRate
	} else {
		global.GlobalMarketLevel2.RWMutex.RLock()
		defer global.GlobalMarketLevel2.RWMutex.RUnlock()

		marketLastPrice = global.GlobalMarketLevel2.LastPrice[quoteKey] // 该值是市场的最新价格，不是子单成交价格
		marketQty := global.GlobalMarketLevel2.EntrustVol[quoteKey]     // 市场委托量
		if marketQty > 0 {
			MarketRate = (float64(lastQty) / float64(marketQty)) * 100 // 比例0-100
		}
		totalTradeVol := global.GlobalMarketLevel2.TradeVol[quoteKey] // 成交量比重
		if totalTradeVol > 0 {
			DealRate = (float64(lastQty) / float64(totalTradeVol)) * 100
		}
		return marketLastPrice, MarketRate, DealRate
	}

}

//func BuildCurMarketTime(u uint64) (out string) {
//	date := time.Now().Format(global.TimeFormatDay)
//	str := strconv.FormatUint(u, 10)
//	if len(str) == 8 {
//		out = date + "0" + str[:3]
//	} else if len(str) == 9 {
//		out = date + str[:4]
//	}
//	return out
//}
// BuildLastMarketTime 前一分钟，入参： 202205061230 格式， 输出 ：202205061229
func BuildLastMarketTime(u int64) (out string) {
	//date := time.Now().Format(global.TimeFormatDay)
	//str := strconv.FormatUint(u, 10)
	//var t string
	//if len(str) == 8 {
	//	t = date + "0" + str[:3]
	//} else if len(str) == 9 {
	//	t = date + str[:4]
	//}
	//if t == "" {
	//	return ""
	//}
	t := cast.ToString(u)
	if len(t) < 12 {
		return ""
	}
	year, _ := strconv.Atoi(t[:4])
	mon, _ := strconv.Atoi(t[4:6])
	day, _ := strconv.Atoi(t[6:8])
	hour, _ := strconv.Atoi(t[8:10])
	min, _ := strconv.Atoi(t[10:])
	out = time.Date(year, time.Month(mon), day, hour, min, 0, 0, time.UTC).Add(-time.Minute * time.Duration(1)).Format(global.TimeFormatMinInt)
	return out
}

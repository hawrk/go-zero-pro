// Package logic
/*
 Author: hawrkchen
 Date: 2022/4/14 14:38
 Desc:
*/
package logic

import (
	pb "algo_assess/assess-mq-server/proto/order"
	"algo_assess/global"
	"fmt"
	"strconv"
	"time"
)

func RealTimeCal(transAt string, in *global.OrderAssess, data *pb.ChildOrderPerf) (*global.OrderAssess, error) {
	if in == nil {
		in = new(global.OrderAssess)
	}
	in.AlgorithmType = uint(data.GetAlgorithmType())
	in.AlgorithmId = uint(data.AlgorithmId)
	in.UsecurityId = uint(data.USecurityId)
	in.SecurityId = data.GetSecurityId()
	in.CreateTime = time.Now()
	in.TimeDimension = 1
	in.ArrivedPrice = int64(data.ArrivedPrice)
	in.TransactAt = TimeStr2TimeStamp(transAt)

	in.LastQty += data.GetLastQty()

	in.SubVwap += data.GetLastPx() * data.GetLastQty()              // 成交价格*成交笔数
	in.SubVwapEntrust += data.GetPrice() * data.GetLastQty()        // 委托价格*成交笔数
	in.SubVwapArrived += data.GetArrivedPrice() * data.GetLastQty() // 到达价格*成交笔数
	orderKey := fmt.Sprintf("%s:%d", transAt, data.Id)
	if _, orderKeyExist := global.GlobalAssess.OrderMap[orderKey]; !orderKeyExist { // 该算法股票不存在同一个订单号才加上委托数量
		in.OrderQty += data.GetOrderQty()
		global.GlobalAssess.OrderMap[orderKey] = struct{}{}
	}
	if data.GetChildOrdStatus() == global.OrderStatusCancel {
		in.CancelQty += (data.GetOrderQty() - data.GetCumQty())
	}
	if data.GetChildOrdStatus() == global.OrderStatusApReject || data.GetChildOrdStatus() == global.OrderStatusCtReject ||
		data.GetChildOrdStatus() == global.OrderStatusTaReject {
		in.RejectedQty += data.GetOrderQty()
	}
	// 计算成交量比重  市场参与率
	quoteKey := fmt.Sprintf("%s%d", data.USecurityId, transAt)
	global.GlobalMarketLevel2.RWMutex.RLock()
	in.LastPrice = global.GlobalMarketLevel2.LastPrice[quoteKey]
	marketQty := global.GlobalMarketLevel2.EntrustVol[quoteKey]
	if marketQty > 0 {
		in.MarketRate = float64(in.OrderQty) / float64(marketQty)
	}
	totalTradeVol := global.GlobalMarketLevel2.TradeVol[quoteKey]
	if totalTradeVol > 0 {
		in.DealRate = float64(in.LastQty) / float64(totalTradeVol)
	}
	global.GlobalMarketLevel2.RWMutex.RUnlock()
	// 计算成交进度
	algoOrderKey := fmt.Sprintf("%s:%d:%d", transAt, data.AlgorithmId, data.USecurityId)
	global.GlobalAlgoOrder.RWMutex.RLock()
	if global.GlobalAlgoOrder.AlgoOrder[algoOrderKey] > 0 {
		in.DealProgress = float64(in.LastQty) / float64(global.GlobalAlgoOrder.AlgoOrder[algoOrderKey])
	}
	global.GlobalAlgoOrder.RWMutex.RUnlock()

	if in.LastQty > 0 {
		in.Vwap = float64(in.SubVwap) / float64(in.LastQty)                                 // 计算vwap
		in.VwapDeviation = in.Vwap - float64(in.SubVwapEntrust)/float64(in.LastQty)         // 计算vwap 滑点
		in.ArrivedPriceDeviation = in.Vwap - float64(in.SubVwapArrived)/float64(in.LastQty) // 计算到达价滑点
	}
	return in, nil
}

func BuildCurMarketTime(u uint64) (out string) {
	date := time.Now().Format(global.TimeFormatDay)
	str := strconv.FormatUint(u, 10)
	if len(str) == 8 {
		out = date + "0" + str[:3]
	} else if len(str) == 9 {
		out = date + str[:4]
	}
	return out
}

func BuildLastMarketTime(u uint64) (out string) {
	date := time.Now().Format(global.TimeFormatDay)
	str := strconv.FormatUint(u, 10)
	var t string
	if len(str) == 8 {
		t = date + "0" + str[:3]
	} else if len(str) == 9 {
		t = date + str[:4]
	}
	year, _ := strconv.Atoi(t[:4])
	mon, _ := strconv.Atoi(t[4:6])
	day, _ := strconv.Atoi(t[6:8])
	hour, _ := strconv.Atoi(t[8:10])
	min, _ := strconv.Atoi(t[10:])
	out = time.Date(year, time.Month(mon), day, hour, min, 0, 0, time.UTC).Add(-time.Minute * time.Duration(1)).Format(global.TimeFormatMinInt)
	return out
}

func TimeStr2TimeStamp(s string) int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	tt, _ := time.ParseInLocation(global.TimeFormatMinInt, s, loc)
	return tt.Unix()
}

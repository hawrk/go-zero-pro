// Package consumer
/*
 Author: hawrkchen
 Date: 2022/4/18 17:08
 Desc:
*/
package consumer

import (
	"algo_assess/global"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// CalculateMarketVwap 这里的作用主要是计算出市场的Vwap，赋给data结构体中，然后序列化到Redis中，客户端查绩效时取这里的数据
func CalculateMarketVwap(data *global.QuoteLevel2Data) int64 {
	// 累计数量计算vwap,就是算一天的总量,不是一分钟的总量
	//logx.Infof("get data:%+v", data)
	dayKey := data.SecID + ":" + cast.ToString(data.OrigTime)[:8]
	//logx.Info("get dayKey:", dayKey)
	// 计算市场vwap
	// GlobalMarketVwap 这个结构体key为证券ID+日期，所以理论上一天的Key增量为5000左右
	global.GlobalMarketVwap.RWMutex.Lock()
	defer global.GlobalMarketVwap.RWMutex.Unlock()

	lastVol := global.GlobalMarketVwap.LastVol[dayKey] // 取上一个成交总量
	//logx.Info("debug core dump: lastVol:", lastVol)
	curVol := data.TotalTradeVol - lastVol // 当前时间点净成交量 =  当前成交总量 - 上一个成交总量
	//logx.Info("curVol:", curVol)

	global.GlobalMarketVwap.TotalPrxCal[dayKey] += curVol * data.LastPrice // ∑(订单成交数量 *成交价格)
	global.GlobalMarketVwap.TotalVol[dayKey] += curVol                     // ∑订单成交数量
	if global.GlobalMarketVwap.TotalPrxCal[dayKey] <= 0 || global.GlobalMarketVwap.TotalVol[dayKey] <= 0 {
		data.Vwap = 0.0000
	} else {
		fVwap := (float64(global.GlobalMarketVwap.TotalPrxCal[dayKey]) / float64(global.GlobalMarketVwap.TotalVol[dayKey])) / 100
		data.Vwap, _ = decimal.NewFromFloat(fVwap).Truncate(4).Float64()
	}
	global.GlobalMarketVwap.MVwap[dayKey] = data.Vwap
	global.GlobalMarketVwap.LastVol[dayKey] = data.TotalTradeVol // 填充当前成交总量
	//logx.Infof("after count get market vwap :%+v", global.GlobalMarketVwap)
	return curVol // 返回当前时间点的净增量，推给下游
}

// BuildLevel2Data 行情委托数量，净成交量，最新价格写入redis,计算绩效时用
func BuildLevel2Data(redis *redis.Redis, netTradeVol int64, data *global.QuoteLevel2Data) error {
	curKey := fmt.Sprintf("%s:%s:%d", global.Level2RedisPrx, data.SecID[3:], data.OrigTime)
	var totalVol int64
	global.GlobalMarketLevel2.RWMutex.Lock()
	global.GlobalMarketLevel2.TradeVol[curKey] += netTradeVol // 成交总增量, 上游已经计算好了
	totalVol = global.GlobalMarketLevel2.TradeVol[curKey]
	global.GlobalMarketLevel2.RWMutex.Unlock()

	m := map[string]string{
		"entrustvol": cast.ToString(data.TotalAskVol + data.TotalBidVol),
		"tradevol":   cast.ToString(totalVol),
		"lastprice":  cast.ToString(data.LastPrice),
	}
	if err := redis.Hmset(curKey, m); err != nil {
		return err
	}
	// 设置过期时间  --大流量会出现slow call,放在定时任务统一处理
	//if err := redis.Expire(curKey, 43200); err != nil {
	//	return err
	//}
	return nil
}

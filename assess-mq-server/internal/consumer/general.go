// Package consumer
/*
 Author: hawrkchen
 Date: 2022/6/21 14:39
 Desc:  绩效概况逻辑处理 （一期）
*/
package consumer

import (
	"algo_assess/global"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"time"
)

func (s *AlgoPlatformOrderTrade) AssessGeneral(orderData *global.ChildOrderData, u AssessUser) {
	// 并发量大时，普通用户和管理员的Dispatch 可单独开一个协程去并发计算
	// 计算普通用户维度的绩效
	if u.NormalUser != "" {
		algoKey := fmt.Sprintf("%s:%d:%s:%d:%s", orderData.SourcePrx, orderData.TransTime, u.NormalUser, orderData.AlgoId, orderData.SecId)
		s.Logger.Info("get nomal user algoKey:", algoKey)
		s.Dispatch(u.NormalUser, algoKey, orderData)
	}

	// 账户管理员维度的绩效
	for _, v := range u.MngrUser {
		algoKey := fmt.Sprintf("%s:%d:%s:%d:%s", orderData.SourcePrx, orderData.TransTime, v, orderData.AlgoId, orderData.SecId)
		s.Logger.Info("get manager user algoKey:", algoKey)
		s.Dispatch(v, algoKey, orderData)
	}

	// 算法厂商绩效
	if u.ProviderUser != "" {
		algoKey := fmt.Sprintf("%s:%d:%s:%d:%s", orderData.SourcePrx, orderData.TransTime, u.ProviderUser, orderData.AlgoId, orderData.SecId)
		s.Logger.Info("get provider user algoKey:", algoKey)
		s.Dispatch(u.ProviderUser, algoKey, orderData)
	}
	// 计算一下超级管理员的
	algoKey := fmt.Sprintf("%s:%d:%s:%d:%s", orderData.SourcePrx, orderData.TransTime, u.AdminUser, orderData.AlgoId, orderData.SecId)
	s.Logger.Info("get admin user algoKey:", algoKey)
	s.Dispatch(u.AdminUser, algoKey, orderData)

}

func (s *AlgoPlatformOrderTrade) Dispatch(user, algoKey string, orderData *global.ChildOrderData) {
	// 实时计算
	global.GlobalAssess.RWMutex.Lock()
	v := global.GlobalAssess.CalAlgo[algoKey] // 取旧的算法数据，如有则全取，如无则默认都是0
	out, err := s.GeneralRealTimeCal(user, v, orderData)
	if err != nil {
		s.Logger.Error("real time cal error:", err)
		return
	}
	global.GlobalAssess.CalAlgo[algoKey] = out
	// 这里计算之后的数据快照不需要写redis,每一分钟产生的交易直接落地，后面的数据都是单独计算，所以不需要reload了
	s.Logger.Info("get algo map len:", len(global.GlobalAssess.CalAlgo))
	global.GlobalAssess.RWMutex.Unlock()
	// 实时计算end
}

func (s *AlgoPlatformOrderTrade) GeneralRealTimeCal(user string,
	in *global.OrderAssess, data *global.ChildOrderData) (*global.OrderAssess, error) {
	if in == nil {
		in = new(global.OrderAssess)
	}
	s.Logger.Infof("begin cal, in:%+v", in)
	in.SourceFrom = data.SourceFrom
	in.BatchNo = data.BatchNo

	in.AlgorithmType = data.AlgorithmType
	in.AlgorithmId = data.AlgoId
	in.UserId = user
	in.UsecurityId = data.UsecId
	in.SecurityId = data.SecId
	in.CreateTime = time.Now()
	in.TimeDimension = 1
	in.ArrivedPrice = data.ArrivePrice
	in.TransactAt = data.TransTime

	in.LastQty += data.LastQty

	// 该子单算法证券状态变更时会推送多次，委托数量只需要统计一次，但是该子单从下单开始到成交，可能会跨越多个时间段，所以要加上当前交易时间
	// 因总线平台有补单操作，母单号不变，所以这里的orderQty直接取母单号的委托数量
	// 只统计当前时间点的子单委托数量，不累加
	orderKey := fmt.Sprintf("%s:%d:%s:%d", data.SourcePrx, data.TransTime, user, data.OrderId)
	s.Logger.Info("get orderKey:", orderKey)
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
	quoteKey := fmt.Sprintf("%s:%s:%d", global.Level2RedisPrx, data.SecId, data.TransTime)
	s.Logger.Info("get quoteKey:", quoteKey)
	in.LastPrice, in.MarketRate, in.DealRate = GetMarketRate(s.svcCtx.RedisClient, quoteKey, in.LastQty)
	// 取市场vwap   ---已经落表，不需要在这里取了
	//in.Vwap = GetMarketVwap(quoteKey)
	// 计算成交进度
	//in.DealProgress, in.OrderQty = GetDealProgress(s.svcCtx.RedisClient, user, data)
	in.DealProgress, _ = GetDealProgress(s.svcCtx.RedisClient, user, data)

	in.SubVwap += data.LastPx * data.LastQty     // 成交价格*成交笔数
	if data.OrderType == 1 && in.LastPrice > 0 { // 如果是限价委托，取当前行情的价格
		in.SubVwapEntrust += in.LastPrice * data.LastQty
	} else {
		in.SubVwapEntrust += data.Price * data.LastQty // 委托价格*成交笔数
	}
	in.SubVwapArrived += data.ArrivePrice * data.LastQty // 到达价格*成交笔数

	// 计算 滑点
	if in.LastQty > 0 {
		in.TradeVwap = float64(in.SubVwap) / float64(in.LastQty)                                           // 计算vwap
		in.VwapDeviation = (in.TradeVwap - float64(in.SubVwapEntrust)/float64(in.LastQty)) / 10000         // 计算vwap 滑点(计算是以分为单位，需除100
		in.ArrivedPriceDeviation = (in.TradeVwap - float64(in.SubVwapArrived)/float64(in.LastQty)) / 10000 // 计算到达价滑点
	}
	s.Logger.Infof("after real cal, out:%+v", in)

	return in, nil
}

// GetDealProgress 成交进度
func GetDealProgress(redis *redis.Redis, user string, data *global.ChildOrderData) (progress float64, eQty int64) {
	date := cast.ToString(data.TransTime)[:8]
	var algoKey, algoDealKey string
	algoKey = fmt.Sprintf("%s:%s:%s:%s:%d:%s", data.SourcePrx, global.AlgoEntrustPrx, date, user, data.AlgoId, data.SecId)  // 母单委托数量
	algoDealKey = fmt.Sprintf("%s:%s:%s:%s:%d:%s", data.SourcePrx, global.AlgoDealPrx, date, user, data.AlgoId, data.SecId) // 母单成交数量
	logx.Info("get algoEntrust:", algoKey, ", algoDealKey:", algoDealKey)

	entrustQty, err := redis.Get(algoKey) // 找到该用户下的母单委托数量
	if err != nil {
		logx.Error("get redis algo order entrust qty err:", err)
		return 0.00, 0
	}
	dealQty, err := redis.Incrby(algoDealKey, data.LastQty) // 累加当前成交量
	if err != nil {
		logx.Error("incrby deal order err:", err)
	}
	if err := redis.Expire(algoDealKey, global.RedisKeyExpireTime); err != nil {
		logx.Error("expire algoDealKey :", algoDealKey, " error:", err)
	}
	logx.Info("get algoEntrust Qty:", entrustQty, ", dealQty:", dealQty)
	if cast.ToInt64(entrustQty) > 0 { // 母单成交数量/ 母单委托数量
		progress = (float64(dealQty) / cast.ToFloat64(entrustQty)) * 100
	}

	return progress, cast.ToInt64(entrustQty)
}

// GetMarketRate 市场参与率， 成交量比重
func GetMarketRate(redis *redis.Redis, quoteKey string, lastQty int64) (marketLastPrice int64, MarketRate, DealRate float64) {
	// 读redis的数据
	out, _ := redis.Hmget(quoteKey, "entrustvol", "tradevol", "lastprice")
	logx.Info("get redis market data:", out)

	if len(out) < 3 { // 共用一份数据
		return
	}
	marketLastPrice = cast.ToInt64(out[2]) * 100 // 行情的价格是以分为单位，这里要再*100与总线保持一致
	marketQty := cast.ToInt64(out[0])
	if marketQty > 0 {
		MarketRate = (float64(lastQty) / float64(marketQty)) * 100
	}
	totalTradeVol := cast.ToInt64(out[1])
	if totalTradeVol > 0 {
		DealRate = (float64(lastQty) / float64(totalTradeVol)) * 100 // 成交量比重
	}
	return marketLastPrice, MarketRate, DealRate
}

// Package comsumer
/*
 Author: hawrkchen
 Date: 2022/4/12 14:35
 Desc: 母单信息下发 (主要作用是子单根据母单号查询母单委托数量）
*/
package consumer

import (
	"algo_assess/assess-mq-server/internal/svc"
	pb "algo_assess/assess-mq-server/proto/order"
	"algo_assess/global"
	"context"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

type AlgoOrderTrade struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlgoOrderTrade(ctx context.Context, svcCtx *svc.ServiceContext) *AlgoOrderTrade {
	return &AlgoOrderTrade{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (s *AlgoOrderTrade) Consume(key string, val string) error {
	if !CheckHeartBeat(s.svcCtx) {
		return nil
	}
	//start := time.Now()
	//time.Sleep(time.Second* 100)
	s.Logger.Info("-----------------algo order start------------------")
	data := pb.AlgoOrderPerf{}
	if err := proto.Unmarshal([]byte(val), &data); err != nil {
		s.Logger.Error("Unmarshal data fail:", err)
		return nil
	}
	s.Logger.Infof("母单原始数据:get algo order origin data:%+v", data)
	if err := CheckAlgoOrderParam(&data); err != nil {
		s.Logger.Error("check param error:", err)
		return nil
	}
	assessAlgoOrder := TransAlgoOrderData(&data)
	s.Logger.Infof("母单转换后数据:get algo order trans data：%+v", assessAlgoOrder)

	if data.SourceFrom == global.SourceFromBus { // 只有总线来源才会落库， 数据修复和数据导入已经先入库了
		global.OrderChan <- assessAlgoOrder
		//if err := s.svcCtx.AlgoOrderRepo.CreateAlgoOrder(s.ctx, &assessAlgoOrder); err != nil {
		//	s.Logger.Error("insert table fail:", err)
		//	return nil
		//}
	}

	// add 母单号ID，用以判断母单必须比子单先到--加上批次号
	algoOrderId := fmt.Sprintf("%s:%s:%d", assessAlgoOrder.SourcePrx, global.AlgoOrderIdPrx, assessAlgoOrder.AlgoId)
	m := make(map[string]string)
	m["basketId"] = cast.ToString(assessAlgoOrder.BasketId)
	m["startTime"] = cast.ToString(assessAlgoOrder.StartTime)
	m["endTime"] = cast.ToString(assessAlgoOrder.EndTime)
	m["algoOrderQty"] = cast.ToString(assessAlgoOrder.AlgoOrderQty)
	m["transTime"] = cast.ToString(assessAlgoOrder.UnixTime)

	if err := s.svcCtx.RedisClient.Hmset(algoOrderId, m); err != nil {
		s.Logger.Error(" set algo order id error:", err)
	}
	if err := s.svcCtx.RedisClient.Expire(algoOrderId, global.RedisKeyExpireTime); err != nil {
		s.Logger.Error("expire algoOrderKey:", algoOrderId, " error:", err)
	}

	// 加上当天日期
	date := cast.ToString(assessAlgoOrder.TransTime)[:8]

	u := GetAllUsers(assessAlgoOrder.UserId, assessAlgoOrder.AlgorithmId)

	//  计算一期成交量用的
	s.SaveAlgoEntrust1(u.NormalUser, u.ProviderUser, u.MngrUser, u.AdminUser, date, &assessAlgoOrder)
	// 二期，保存母单委托数据
	s.SaveAlgoEntrust2(u.NormalUser, u.ProviderUser, u.MngrUser, u.AdminUser, date, &assessAlgoOrder)

	// 保存篮子信息
	//s.SaveBasketInfo(normalUser, providerUser, mngrUser, adminUser, date, &assessAlgoOrder)

	//global.GlobalAlgoOrder.RWMutex.Lock()
	//global.GlobalAlgoOrder.AlgoOrder[algoOrderKey] += assessAlgoOrder.AlgoOrderQty
	////s.Logger.Infof("get map:%+v", global.GlobalAlgoOrder.AlgoOrder)
	//for k,v := range global.GlobalAlgoOrder.AlgoOrder {
	//	s.Logger.Info("get key:", k, ", get val:", v)
	//}
	////s.Logger.Info("get map len:", len(global.GlobalAlgoOrder.AlgoOrder))
	//global.GlobalAlgoOrder.RWMutex.Unlock()

	//time.Sleep(time.Second * 10)
	//s.svcCtx.BigCache.Set(cast.ToString(data.Id), tools.IntToBytes(data.AlgoOrderQty))
	//s.Logger.Error("time.Since(start):", time.Since(start))
	return nil
}

//SaveAlgoEntrust1 保存母单委托数量，算完成时用到，注意这里要区分普通用户，管理员用户和汇总用户
func (s *AlgoOrderTrade) SaveAlgoEntrust1(norUser, proUser string, mngrUser []string, adminUser string, date string, assessAlgoOrder *global.MAlgoOrder) {
	if norUser != "" {
		s.Save2Redis(date, norUser, assessAlgoOrder)
	}
	for _, v := range mngrUser { // 普通用户所属的券商管理员mngr类型的用户
		s.Save2Redis(date, v, assessAlgoOrder)
	}
	if proUser != "" {
		s.Save2Redis(date, proUser, assessAlgoOrder)
	}
	// 超级管理员，汇总所有数据
	s.Save2Redis(date, adminUser, assessAlgoOrder)
}

func (s *AlgoOrderTrade) Save2Redis(date string, user string, assessAlgoOrder *global.MAlgoOrder) {
	algoOrderKey := fmt.Sprintf("%s:%s:%s:%s:%d:%s",
		assessAlgoOrder.SourcePrx, global.AlgoEntrustPrx, date, user, assessAlgoOrder.AlgorithmId, assessAlgoOrder.SecId)
	// 入Redis
	enQty, err := s.svcCtx.RedisClient.Incrby(algoOrderKey, assessAlgoOrder.AlgoOrderQty)
	if err != nil {
		s.Logger.Error("incrby order entrust err:", err)
	}
	if err := s.svcCtx.RedisClient.Expire(algoOrderKey, global.RedisKeyExpireTime); err != nil {
		s.Logger.Error("expire algoOrderKey:", algoOrderKey, " error:", err)
	}
	s.Logger.Info("一期母单 get algoOrderKey:", algoOrderKey, ", entrust qty:", enQty)
}

// SaveAlgoEntrust2 保存母单委托数量，在子单绩效计算的时候用到，需要根据单个用户和所有用户的委托数量保存
// 二期指标计算用
func (s *AlgoOrderTrade) SaveAlgoEntrust2(norUser, proUser string, mngrUser []string, adminUser string, date string, assessAlgoOrder *global.MAlgoOrder) {
	// 普通用户维度保存
	if norUser != "" {
		s.BuildRedisData(norUser, date, assessAlgoOrder)
	}
	// 管理员用户， 根据算法维度计算保存
	for _, v := range mngrUser {
		s.BuildRedisData(v, date, assessAlgoOrder)
	}

	// add 加上算法厂商用户
	if proUser != "" {
		s.BuildRedisData(proUser, date, assessAlgoOrder)
	}
	// 超级管理员
	s.BuildRedisData(adminUser, date, assessAlgoOrder)
}

func (s *AlgoOrderTrade) BuildRedisData(userId string, date string, algo *global.MAlgoOrder) {
	//userAlgoKey := fmt.Sprintf("%s:%s:%d:%s:%d",
	//	global.UserAlgoEntrust, date, batchNo, userId, algoId)
	userAlgoKey := fmt.Sprintf("%s:%s:%s:%s:%d:%s:%d",
		algo.SourcePrx, global.UserAlgoEntrust, date, userId, algo.AlgorithmId, algo.SecId, algo.AlgoId)
	userEnQty, err := s.svcCtx.RedisClient.Incrby(userAlgoKey, algo.AlgoOrderQty)
	if err != nil {
		s.Logger.Error("incrby algo user order entrust fail:", err)
	}
	if err := s.svcCtx.RedisClient.Expire(userAlgoKey, global.RedisKeyExpireTime); err != nil {
		s.Logger.Error("expire userAlgoKey fail:", err)
	}
	s.Logger.Info("二期用户 get userAlgoKey:", userAlgoKey, ",user entrust Qty: ", userEnQty)
}

// SaveBasketInfo 保存篮子信息，用于统计交易订单数量
func (s *AlgoOrderTrade) SaveBasketInfo(norUser, proUser string, mngrUser []string, adminUser string, date string, assessAlgoOrder *global.MAlgoOrder) {
	// 普通用户维度保存
	if norUser != "" {
		s.BuildBasketData(norUser, date, assessAlgoOrder.BatchNo, assessAlgoOrder.AlgorithmId, assessAlgoOrder.BasketId)
	}
	// 管理员用户， 根据算法维度计算保存
	for _, v := range mngrUser {
		s.BuildBasketData(v, date, assessAlgoOrder.BatchNo, assessAlgoOrder.AlgorithmId, assessAlgoOrder.BasketId)
	}

	// add 加上算法厂商用户
	if proUser != "" {
		s.BuildBasketData(proUser, date, assessAlgoOrder.BatchNo, assessAlgoOrder.AlgorithmId, assessAlgoOrder.BasketId)
	}
	// 超级管理员
	s.BuildBasketData(adminUser, date, assessAlgoOrder.BatchNo, assessAlgoOrder.AlgorithmId, assessAlgoOrder.BasketId)
}

func (s *AlgoOrderTrade) BuildBasketData(userId string, date string, batchNo int64, algoId int, baseketId int) {
	//basketKey := fmt.Sprintf("%s:%s:%d:%s:%d",
	//	global.BasketsPrx, date, batchNo,userId, algoId)
	basketKey := fmt.Sprintf("%s:%s:%s:%d",
		global.BasketsPrx, date, userId, algoId)
	s.Logger.Info("basketKey:", basketKey)
	// 转存redis
	// 存储Set类型
	_, err := s.svcCtx.RedisClient.Sadd(basketKey, baseketId)
	if err != nil {
		s.Logger.Error("set basket info error:", err)
	}
	if err := s.svcCtx.RedisClient.Expire(basketKey, global.RedisKeyExpireTime); err != nil {
		s.Logger.Error("expire basketKey fail:", err)
	}
	/*
		global.GAlgoOrderBasket.RWMutex.Lock()
		baskets := global.GAlgoOrderBasket.Baskets[basketKey]
		if !tools.InSlice(baskets, assessAlgoOrder.BasketId) {
			baskets = append(baskets, assessAlgoOrder.BasketId)
		}
		global.GAlgoOrderBasket.Baskets[basketKey] = baskets
		s.Logger.Info("get baskets order len:", len(baskets))
		global.GAlgoOrderBasket.RWMutex.Unlock()
	*/
}

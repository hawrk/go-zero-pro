// Package kafkaq
/*
 Author: hawrkchen
 Date: 2022/3/24 15:04
 Desc:  订单交易信息接收
*/
package consumer

import (
	"algo_assess/assess-mq-server/internal/dao"
	"algo_assess/assess-mq-server/internal/svc"
	pb "algo_assess/assess-mq-server/proto/order"
	"algo_assess/global"
	"algo_assess/pkg/tools"
	"context"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
	"strings"
	"time"
)

type AlgoPlatformOrderTrade struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlgoPlatformOrderTrade(ctx context.Context, svcCtx *svc.ServiceContext) *AlgoPlatformOrderTrade {
	return &AlgoPlatformOrderTrade{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 子单版本号
var childOrderVersion = 0

// CheckHeartBeat 心跳检测，用于在集群或主备环境中，只能有一个服务在计算
// 返回true 表示需要计算，  返回false 表示不需要计算
func CheckHeartBeat(s *svc.ServiceContext) bool {
	if !s.Config.WorkProcesser.EnableCheckHeartBeat { // need check heart beat
		return true
	}
	v, err := s.RedisClient.Get(global.HeartBeatKey)
	if err != nil { // 取不到时，默认也返回true
		logx.Error("get heart beat key error:", err)
		return true
	}
	if v == "" {
		return true
	}
	// 3.有数据了，先判断一下是不是自已的
	sli := strings.Split(v, ":")
	if len(sli) < 2 {
		return true
	}
	//1. 是自已的,直接返回true
	if tools.GetLocalIP() == sli[0] {
		return true
	} else { // 不是自已的话，取当前时间戳，然后比较时间戳， 如果大于90秒，则表示另外一个机器可能挂了
		now := time.Now().Unix()
		if now-cast.ToInt64(sli[1]) > 90 {
			// reload 快照数据，然后接管该计算任务
			logx.Info("master machine has down, current machine cover...")
			dao.Reload(s.Config, s)
			// reload 完后，需要更新redis 同步数据
			now := time.Now().Unix()
			ip := tools.GetLocalIP()
			v := fmt.Sprintf("%s:%d", ip, now)
			if err := s.RedisClient.Setex(global.HeartBeatKey, v, 60*60); err != nil {
				logx.Error("set heart beat key error:", err)
			}
			return true
		} else { // 小于90秒时，表示有另外的主机在处理，这里就不用处理了
			return false
		}
	}
}

func (s *AlgoPlatformOrderTrade) Consume(_ string, val string) error {
	if !CheckHeartBeat(s.svcCtx) {
		return nil
	}
	//start := time.Now()
	data := pb.ChildOrderPerf{}
	if err := proto.Unmarshal([]byte(val), &data); err != nil {
		s.Logger.Error("Unmarshal data fail:", err)
		return nil
	}
	s.Logger.Info("-----------------child order start------------------")
	// TODO:  控制协程并发数量
	//threading.GoSafe(func() {
	//s.Logger.Infof("get data:%+v", data)
	s.Logger.Infof(" 子单原始数据:get origin data:%+v", data)
	if err := CheckOrderParam(&data); err != nil {
		s.Logger.Info("skipping process, reason:", err)
		return nil
	}
	// 结构体转换
	orderData := TransChildOrderData(&data)
	// check一下母单信息是否已到达
	algoOrderIdKey := fmt.Sprintf("%s:%s:%d", orderData.SourcePrx, global.AlgoOrderIdPrx, orderData.AlgoOrderId)
	if !CheckAlgoOrder(algoOrderIdKey, s.svcCtx.RedisClient, &data) {
		s.Logger.Info("algo order not found, return...")
		return nil
	}
	// 落地DB
	if orderData.SourceFrom == global.SourceFromBus {
		global.OrderDetailChan <- orderData
		//if err := s.svcCtx.OrderDetailRepo.CreateOrderDetail(s.ctx, &orderData); err != nil {
		//	s.Logger.Error("insert into order detail fail:", err)
		//	return nil // 写表失败，有可能是主键重复，不再处理
		//}
	}
	// 取行情市场价格
	marketPrice := GetArrviPrice(s.svcCtx.SzMarketRepo, s.svcCtx.ShMarketRepo, s.svcCtx.RedisClient, orderData.SecId, orderData.TransTime)
	orderData.MarketPrice = marketPrice // 补一下行情市价格
	s.Logger.Infof("子单转换后数据:get order trans Data:%+v", orderData)

	// 时间戳落Redis
	s.svcCtx.RedisClient.Sadd(global.AssessTimeSetKey, orderData.TransTime)
	// 取所有计算用户
	u := GetAllUsers(orderData.UserId, orderData.AlgoId)
	//var wg sync.WaitGroup
	// 计算一期绩效汇总
	if s.svcCtx.Config.WorkProcesser.EnableFirstPhase {
		s.Logger.Info("process first phase....................")
		//wg.Add(1)
		//threading.GoSafe(func() {
		//defer wg.Done()
		s.AssessGeneral(&orderData, u)
		//})
	}

	// 二期算法分析   （dashboard,画像）
	if s.svcCtx.Config.WorkProcesser.EnableSecondPhase {
		s.Logger.Info("process second phase.....................")
		//wg.Add(1)
		//threading.GoSafe(func() {
		//defer wg.Done()
		if s.svcCtx.Config.WorkProcesser.EnableConcurrency {
			s.AlgoAnalysisConcurrent(&orderData, u)
		} else {
			s.AlgoAnalysis(&orderData, u)
		}
		//})
	}
	//wg.Wait()
	//s.Logger.Error("time.Since(start):", time.Since(start))

	return nil
}

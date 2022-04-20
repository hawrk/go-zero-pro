// Package kafkaq
/*
 Author: hawrkchen
 Date: 2022/3/24 15:04
 Desc:  订单交易信息接收
*/
package consumer

import (
	"algo_assess/assess-mq-server/internal/logic"
	"algo_assess/assess-mq-server/internal/svc"
	pb "algo_assess/assess-mq-server/proto/order"
	"algo_assess/global"
	"context"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
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

func (s *AlgoPlatformOrderTrade) Consume(_ string, val string) error {
	//defer func() {
	//	if err := recover(); err != nil {
	//		fmt.Println("recover:", err)
	//	}
	//}()
	// 消费
	time.Sleep(time.Second * 100)
	s.Logger.Info("-----------------start------------------")
	data := pb.ChildOrderPerf{}
	if err := proto.Unmarshal([]byte(val), &data); err != nil {
		s.Logger.Error("Unmarshal data fail:", err)
		return err
	}

	s.Logger.Infof("get data:%+v", data)
	data.Id = 130
	data.AlgorithmId = 2560
	data.USecurityId = 79
	//TODO:
	//if data.GetAlgoOrderId() == 0 || data.AlgorithmId == 0 {  // 过滤普通单，没有绑定算法的不参与计算
	//	s.Logger.Info("normal order:", data.GetId())
	//	return nil
	//}
	if data.GetTransactTime() <= 0 {
		s.Logger.Info("TransactTime empty", data.GetId())
		return nil
	}
	if data.GetOrderQty() <= 0 {
		s.Logger.Error(" order qty invalid:", data.GetId())
		return nil
	}
	if data.GetChildOrdStatus() == global.OrderStatusApAccept || data.GetChildOrdStatus() == global.OrderStatusCtAccept {
		s.Logger.Info("state access, continue")
		return nil
	}

	// 交易时间计算到分钟
	transactAt := time.UnixMicro(int64(data.GetTransactTime())).Format(global.TimeFormatMinInt)
	s.Logger.Info("get transactAt:", transactAt)
	transact := cast.ToInt64(transactAt)
	// 时间戳落Redis
	s.svcCtx.RedisClient.Sadd(global.AssessTimeSetKey, transact)
	// 落地DB
	if err := s.svcCtx.OrderDetailRepo.CreateOrderDetail(s.ctx, transact, &data); err != nil {
		s.Logger.Error("insert into order detail fail:", err)
		//return err
	}

	algoId := fmt.Sprintf("%s:%d:%d", transactAt, data.AlgorithmId, data.USecurityId)
	// 定时计算时，直接落缓存落表
	//global.GlobalOrders.RWMutex.Lock()
	//global.GlobalOrders.CalOrders[transact] = append(global.GlobalOrders.CalOrders[transact], &data)
	//global.GlobalOrders.RWMutex.Unlock()
	//s.Logger.Info("get map len:", len(global.GlobalOrders.CalOrders), ",slice: ", transactAt, ",len:", len(global.GlobalOrders.CalOrders[transact]))

	// 实时计算
	global.GlobalAssess.RWMutex.Lock()
	v := global.GlobalAssess.CalAlgo[algoId] // 取旧的算法数据，如有则全取，如无则默认都是0
	out, err := logic.RealTimeCal(transactAt, v, &data)
	if err != nil {
		s.Logger.Error("divisor zero, invalid order,get key:", algoId)
		return nil
	}
	global.GlobalAssess.CalAlgo[algoId] = out
	s.Logger.Info("get map key:", algoId)
	s.Logger.Infof("get map value:%+v", *global.GlobalAssess.CalAlgo[algoId])
	global.GlobalAssess.RWMutex.Unlock()

	s.Logger.Info("get algo map len:", len(global.GlobalAssess.CalAlgo))
	// 实时计算end

	time.Sleep(time.Second * 10)
	return nil
}

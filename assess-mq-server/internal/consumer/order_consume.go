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
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
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
	//time.Sleep(time.Second * 1)
	s.Logger.Info("-----------------child order start------------------")
	data := pb.ChildOrderPerf{}
	if err := proto.Unmarshal([]byte(val), &data); err != nil {
		s.Logger.Error("Unmarshal data fail:", err)
		return err
	}
	//s.Logger.Infof("get data:%+v", data)
	if err := checkOrderParam(&data); err != nil {
		s.Logger.Error("skipping process, err:", err)
		return nil
	}

	// 结构体转换
	orderData := TransChildOrderData(&data)
	s.Logger.Infof("get orderData:%+v", orderData)

	// 时间戳落Redis
	s.svcCtx.RedisClient.Sadd(global.AssessTimeSetKey, orderData.TransTime)
	// 落地DB
	if err := s.svcCtx.OrderDetailRepo.CreateOrderDetail(s.ctx, &orderData); err != nil {
		s.Logger.Error("insert into order detail fail:", err)
		//return err
	}
	algoKey := fmt.Sprintf("%d:%d:%s", orderData.TransTime, orderData.AlgoId, orderData.SecId)
	s.Logger.Info("get algoKey:", algoKey)

	// 定时计算时，直接落缓存落表
	//global.GlobalOrders.RWMutex.Lock()
	//global.GlobalOrders.CalOrders[transact] = append(global.GlobalOrders.CalOrders[transact], &data)
	//global.GlobalOrders.RWMutex.Unlock()
	//s.Logger.Info("get map len:", len(global.GlobalOrders.CalOrders), ",slice: ", transactAt, ",len:", len(global.GlobalOrders.CalOrders[transact]))
	// 实时计算
	global.GlobalAssess.RWMutex.Lock()
	v := global.GlobalAssess.CalAlgo[algoKey] // 取旧的算法数据，如有则全取，如无则默认都是0
	out, err := logic.RealTimeCal(s.svcCtx, v, &orderData)
	if err != nil {
		s.Logger.Error("real time cal error:", err)
		return nil
	}
	global.GlobalAssess.CalAlgo[algoKey] = out
	//s.Logger.Info("get map key:", orderKey)
	//s.Logger.Infof("get map value:%+v", *global.GlobalAssess.CalAlgo[orderKey])
	s.Logger.Info("get algo map len:", len(global.GlobalAssess.CalAlgo))
	global.GlobalAssess.RWMutex.Unlock()
	// 实时计算end

	//time.Sleep(time.Second * 10)
	return nil
}

func checkOrderParam(v *pb.ChildOrderPerf) error {
	if v.GetAlgoOrderId() == 0 || v.AlgorithmId == 0 { // 过滤普通单，没有绑定算法的不参与计算
		return errors.New("normal order, continue")
	}
	if v.GetTransactTime() <= 0 {
		return errors.New("error field transTime")
	}
	if v.GetOrderQty() <= 0 {
		return errors.New("error field orderQty")
	}
	if v.GetChildOrdStatus() == global.OrderStatusApAccept || v.GetChildOrdStatus() == global.OrderStatusCtAccept {
		return errors.New("accept status, continue")
	}
	return nil
}

// Package kafkaq
/*
 Author: hawrkchen
 Date: 2022/3/24 15:04
 Desc:  订单交易信息接收
*/
package logic

import (
	"algo_assess/global"
	"algo_assess/mqueue/internal/svc"
	pb "algo_assess/mqueue/proto/order"
	"context"
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

func (s *AlgoPlatformOrderTrade) Consume(key string, val string) error {
	// 消费
	// s.Logger.Info("get order trade consume:key :", key, ", val:", val)
	data := pb.ChildOrderPerf{}
	if err := proto.Unmarshal([]byte(val), &data); err != nil {
		s.Logger.Error("Unmarshal data fail:", err)
		return err
	}
	s.Logger.Info("get data:", data)
	// 交易时间计算到分钟
	transactAt := time.Unix(cast.ToInt64(data.GetTransactTime()), 0).Format(global.TimeFormatMinInt)
	s.Logger.Info("get transactAt:", transactAt)
	transact := cast.ToUint64(transactAt)

	global.GlobalOrders.RWMutex.Lock()
	global.GlobalOrders.CalOrders[transact] = append(global.GlobalOrders.CalOrders[transact], data)
	global.GlobalOrders.RWMutex.Unlock()

	// s.Logger.Info("get map len:", len(global.GlobalOrders.CalOrders), ",slice: ", transactAt,",len:", len(global.GlobalOrders.CalOrders[transact]))
	// 落地DB
	if err := s.svcCtx.OrderDetail.CreateOrderDetail(s.ctx, transact, &data); err != nil {
		s.Logger.Error("insert into order detail fail:", err)
		return err
	}

	time.Sleep(30 * time.Second)
	return nil
}

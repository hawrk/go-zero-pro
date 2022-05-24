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
	//time.Sleep(time.Second* 100)
	s.Logger.Info("-----------------algo order start------------------")
	data := pb.AlgoOrderPerf{}
	if err := proto.Unmarshal([]byte(val), &data); err != nil {
		s.Logger.Error("Unmarshal data fail:", err)
		return err
	}
	//s.Logger.Infof("get data:%+v", data)
	assessAlgoOrder := TransAlgoOrderData(&data)

	if err := s.svcCtx.AlgoOrderRepo.CreateAlgoOrder(s.ctx, &assessAlgoOrder); err != nil {
		s.Logger.Error("insert table fail:", err)
		return nil
	}
	algoOrderKey := fmt.Sprintf("%d:%s", data.AlgorithmId, assessAlgoOrder.SecId)
	s.Logger.Info("get algoOrderKey:", algoOrderKey)
	global.GlobalAlgoOrder.RWMutex.Lock()
	global.GlobalAlgoOrder.AlgoOrder[algoOrderKey] += assessAlgoOrder.AlgoOrderQty
	s.Logger.Infof("get map:%v", global.GlobalAlgoOrder.AlgoOrder)
	//s.Logger.Info("get map len:", len(global.GlobalAlgoOrder.AlgoOrder))
	global.GlobalAlgoOrder.RWMutex.Unlock()

	//time.Sleep(time.Second * 10)
	//s.svcCtx.BigCache.Set(cast.ToString(data.Id), tools.IntToBytes(data.AlgoOrderQty))
	return nil
}

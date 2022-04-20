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
	"time"
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

func (s *AlgoOrderTrade) Consume(_ string, val string) error {
	// s.Logger.Info("get algo order trade consume:key :", key, ", val:", val)
	data := pb.AlgoOrderPerf{}
	if err := proto.Unmarshal([]byte(val), &data); err != nil {
		s.Logger.Error("Unmarshal data fail:", err)
		return err
	}
	s.Logger.Info("get data:", data)
	transactAt := time.UnixMicro(int64(data.GetTransactTime())).Format(global.TimeFormatMinInt)
	key := fmt.Sprintf("%s:%d:%d", transactAt, data.AlgorithmId, data.USecurityId)

	global.GlobalAlgoOrder.RWMutex.Lock()
	global.GlobalAlgoOrder.AlgoOrder[key] += data.AlgoOrderQty
	global.GlobalAlgoOrder.RWMutex.Unlock()

	//s.svcCtx.BigCache.Set(cast.ToString(data.Id), tools.IntToBytes(data.AlgoOrderQty))

	return nil
}

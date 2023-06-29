// Package consumer
/*
 Author: hawrkchen
 Date: 2023/5/20 10:22
 Desc:
*/
package consumer

import (
	mqpb "algo_assess/assess-mq-server/proto"
	"algo_assess/mq-router-server/internal/svc"
	pb "algo_assess/mq-router-server/ordproto/asorder"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
	"time"
)

type ChildOrderConsume struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChildOrderConsume(ctx context.Context, svcCtx *svc.ServiceContext) *ChildOrderConsume {
	return &ChildOrderConsume{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (s *ChildOrderConsume) Consume(key string, val string) error {

	data := pb.ChildOrderPerf{}
	if err := proto.Unmarshal([]byte(val), &data); err != nil {
		s.Logger.Error("Unmarshal data fail:", err)
		return nil
	}
	//s.Logger.Info("-----------------child order start------------------")
	//s.Logger.Infof(" 子单原始数据:get origin data:%+v", data)

	req := &mqpb.ApiAlgoOrderReq{
		Value: []byte(val),
	}

	if data.AlgorithmId%2 == 0 {
		st := time.Now()
		_, err := s.svcCtx.AssessMQMClient.SendChildOrder(s.ctx, req)
		if err != nil {
			s.Logger.Error("rpc SendChildOrder error:", err)
		}
		s.Logger.Info("child time latency:", time.Since(st))
	} else {
		st := time.Now()
		_, err := s.svcCtx.AssessMQSClient.SendChildOrder(s.ctx, req)
		if err != nil {
			s.Logger.Error("rpc SendChildOrder error:", err)
		}
		s.Logger.Info("child time latency:", time.Since(st))
	}

	return nil
}

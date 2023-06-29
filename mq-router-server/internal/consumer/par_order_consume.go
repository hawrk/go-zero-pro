// Package consumer
/*
 Author: hawrkchen
 Date: 2023/5/20 10:22
 Desc:
*/
package consumer

import (
	mqpb "algo_assess/assess-mq-server/proto"
	"algo_assess/mq-router-server/internal/job"
	"algo_assess/mq-router-server/internal/svc"
	pb "algo_assess/mq-router-server/ordproto/asorder"
	"context"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
	"time"
)

type AlgoOrderConsume struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlgoOrderConsume(ctx context.Context, svcCtx *svc.ServiceContext) *AlgoOrderConsume {
	return &AlgoOrderConsume{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (s *AlgoOrderConsume) Consume(key string, val string) error {

	data := pb.AlgoOrderPerf{}
	if err := proto.Unmarshal([]byte(val), &data); err != nil {
		s.Logger.Error("Unmarshal data fail:", err)
		return nil
	}
	//s.Logger.Info("-----------------algo order start------------------")
	//s.Logger.Infof("母单原始数据:get algo order origin data:%+v", data)
	// 调用assessmq server
	req := &mqpb.ApiAlgoOrderReq{
		Value: []byte(val),
	}

	node := job.ConsisHash.GetNode(cast.ToString(data.Id))
	s.Logger.Info("get mq-server node:", node)
	// 初始化客户端调用对象

	return nil

	if data.AlgorithmId%2 == 0 {
		st := time.Now()
		_, err := s.svcCtx.AssessMQMClient.SendAlgoOrder(s.ctx, req)
		if err != nil {
			s.Logger.Error("rpc SendAlgoOrder error:", err)
			return nil
		}
		s.Logger.Info("algo time latency:", time.Since(st))
	} else {
		st := time.Now()
		_, err := s.svcCtx.AssessMQSClient.SendAlgoOrder(s.ctx, req)
		if err != nil {
			s.Logger.Error("rpc SendAlgoOrder error:", err)
			return nil
		}
		s.Logger.Info("algo time latency:", time.Since(st))
	}

	return nil
}

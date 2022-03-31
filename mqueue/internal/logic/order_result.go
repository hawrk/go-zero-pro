// Package logic
/*
 Author: hawrkchen
 Date: 2022/3/24 15:10
 Desc: 订单交易成交结果回执
*/
package logic

import (
	"algo_assess/mqueue/internal/svc"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type AlgoPlatformOrderResult struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlgoPlatformOrderResult(ctx context.Context, svcCtx *svc.ServiceContext) *AlgoPlatformOrderResult {
	return &AlgoPlatformOrderResult{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (s *AlgoPlatformOrderResult) Consume(key string, val string) error {
	s.Logger.Info("consume trade result |k:", key, ", value:", val)

	return nil
}

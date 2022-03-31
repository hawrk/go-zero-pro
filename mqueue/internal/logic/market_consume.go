// Package logic
/*
 Author: hawrkchen
 Date: 2022/3/24 15:14
 Desc:
*/
package logic

import (
	"algo_assess/mqueue/internal/svc"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type AlgoPlatformMarketInfo struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlgoPlatformMarketInfo(ctx context.Context, svcCtx *svc.ServiceContext) *AlgoPlatformMarketInfo {
	return &AlgoPlatformMarketInfo{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (s *AlgoPlatformMarketInfo) Consume(key string, val string) error {
	s.Logger.Info("cumsume market info|key:", key, ", value:", val)

	return nil
}

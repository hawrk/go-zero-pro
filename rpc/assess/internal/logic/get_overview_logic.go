package logic

import (
	"context"

	"algo_assess/rpc/assess/internal/svc"
	"algo_assess/rpc/assess/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOverviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOverviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOverviewLogic {
	return &GetOverviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOverviewLogic) GetOverview(in *proto.OverviewReq) (*proto.OverVewRsp, error) {
	// todo: add your logic here and delete this line

	return &proto.OverVewRsp{}, nil
}

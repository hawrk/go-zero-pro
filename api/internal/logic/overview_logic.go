package logic

import (
	"context"

	"algo_assess/api/internal/svc"
	"algo_assess/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OverviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOverviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OverviewLogic {
	return &OverviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OverviewLogic) Overview(req *types.OverviewReq) (resp *types.OverviewRsp, err error) {
	// todo: add your logic here and delete this line

	return
}

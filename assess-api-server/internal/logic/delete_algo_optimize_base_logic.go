package logic

import (
	"algo_assess/assess-rpc-server/assessservice"
	"context"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAlgoOptimizeBaseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteAlgoOptimizeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAlgoOptimizeBaseLogic {
	return &DeleteAlgoOptimizeBaseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAlgoOptimizeBaseLogic) DeleteAlgoOptimizeBase(req *types.DeleteOptimizeBaseReq) (resp *types.OptimizeBaseRsp, err error) {
	rsp, err := l.svcCtx.AssessClient.DeleteOptimizeBase(l.ctx, &assessservice.DeleteOptimizeBaseReq{
		Id: req.Id,
	})
	p := &types.OptimizeBaseRsp{
		Code: int(rsp.GetCode()),
		Msg:  rsp.GetMsg(),
	}
	return p, err
}

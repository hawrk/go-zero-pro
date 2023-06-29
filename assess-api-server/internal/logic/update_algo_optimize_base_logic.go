package logic

import (
	"algo_assess/assess-rpc-server/assessservice"
	"context"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAlgoOptimizeBaseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAlgoOptimizeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAlgoOptimizeBaseLogic {
	return &UpdateAlgoOptimizeBaseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAlgoOptimizeBaseLogic) UpdateAlgoOptimizeBase(req *types.UpdateOptimizeBaseReq) (resp *types.OptimizeBaseRsp, err error) {
	rsp, err := l.svcCtx.AssessClient.UpdateOptimizeBase(l.ctx, &assessservice.UpdateOptimizeBaseReq{
		Id:           req.Id,
		ProviderId:   req.ProviderId,
		ProviderName: req.ProviderName,
		SecId:        req.SecId,
		SecName:      req.SecName,
		AlgoId:       req.AlgoId,
		AlgoType:     req.AlgoType,
		AlgoName:     req.AlgoName,
		OpenRate:     req.OpenRate,
		IncomeRate:   req.IncomeRate,
		BasisPoint:   req.BasisPoint,
	})
	p := &types.OptimizeBaseRsp{
		Code: int(rsp.GetCode()),
		Msg:  rsp.GetMsg(),
	}
	return p, err
}
